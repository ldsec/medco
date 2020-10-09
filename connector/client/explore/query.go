package medcoclient

import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ldsec/medco/connector/restapi/client"
	"github.com/ldsec/medco/connector/restapi/client/medco_node"
	"github.com/ldsec/medco/connector/restapi/models"
	utilclient "github.com/ldsec/medco/connector/util/client"
	"github.com/ldsec/medco/connector/wrappers/unlynx"
	"github.com/sirupsen/logrus"
)

// ExploreQuery is a MedCo client explore query
type ExploreQuery struct {

	// httpMedCoClients are the HTTP clients for the MedCo connectors
	httpMedCoClients []*client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

	// userPublicKey is the user public key
	userPublicKey string
	// userPrivateKey is the user private key
	userPrivateKey string

	// queryType is the type of explore query requested
	queryType models.ExploreQueryType
	// encPanelsItemKeys is part of the query: contains the encrypted item keys organized by panel
	encPanelsItemKeys [][]string
	// panelsIsNot is part of the query: indicates which panels are negated
	panelsIsNot []bool
}

// NewExploreQuery creates a new MedCo client query
func NewExploreQuery(authToken string, queryType models.ExploreQueryType, encPanelsItemKeys [][]string, panelsIsNot []bool, disableTLSCheck bool) (q *ExploreQuery, err error) {

	q = &ExploreQuery{
		authToken:         authToken,
		queryType:         queryType,
		encPanelsItemKeys: encPanelsItemKeys,
		panelsIsNot:       panelsIsNot,
	}

	// retrieve network information
	getMetadataResp, err := utilclient.MetaData(authToken, disableTLSCheck)
	if err != nil {
		logrus.Error(err)
		return
	}

	// parse network information
	q.httpMedCoClients = make([]*client.MedcoCli, len(getMetadataResp.Payload.Nodes))
	for _, node := range getMetadataResp.Payload.Nodes {
		if q.httpMedCoClients[*node.Index] != nil {
			err = errors.New("duplicated node index in network metadata")
			logrus.Error(err)
			return
		}

		nodeURL, err := url.Parse(node.URL)
		if err != nil {
			logrus.Error("cannot parse MedCo connector URL: ", err)
			return nil, err
		}

		nodeTransport := httptransport.New(nodeURL.Host, nodeURL.Path, []string{nodeURL.Scheme})
		nodeTransport.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: disableTLSCheck}
		q.httpMedCoClients[*node.Index] = client.New(nodeTransport, nil)
	}

	// generate ephemeral pair of user keys
	q.userPublicKey, q.userPrivateKey, err = unlynx.GenerateKeyPair()
	if err != nil {
		return
	}

	return
}

type nodeResult struct {
	Body      *models.ExploreQueryResultElement
	NodeIndex int
}

// Execute executes the MedCo client query synchronously on all the nodes
func (clientQuery *ExploreQuery) Execute() (nodesResult map[int]*ExploreQueryResult, err error) {

	queryResultsChan := make(chan nodeResult)
	queryErrChan := make(chan error)

	// execute requests on all nodes
	for idx := range clientQuery.httpMedCoClients {
		go func(idxLocal int) {

			queryResult, err := clientQuery.submitToNode(idxLocal)
			if err != nil {
				queryErrChan <- err
			} else {
				queryResultsChan <- nodeResult{Body: queryResult, NodeIndex: idxLocal}
			}

		}(idx)
	}

	// parse the results as they come, or interrupt if one of them errors, or if a timeout occurs
	timeout := time.After(time.Duration(utilclient.QueryTimeoutSeconds) * time.Second)
	nodesResult = make(map[int]*ExploreQueryResult)
forLoop:
	for range clientQuery.httpMedCoClients {
		select {
		case queryResult := <-queryResultsChan:
			parsedQueryResult, err := newQueryResult(queryResult.Body, clientQuery.userPrivateKey)
			if err != nil {
				return nil, err
			}

			nodesResult[queryResult.NodeIndex] = parsedQueryResult
			logrus.Info("MedCo client explore query successful for node ", queryResult.NodeIndex)

			if len(nodesResult) == len(clientQuery.httpMedCoClients) {
				logrus.Info("MedCo client explore query successful for all resources")
				return nodesResult, nil
			}

		case err = <-queryErrChan:
			logrus.Error("MedCo client explore query error: ", err)
			break forLoop

		case <-timeout:
			err = errors.New("MedCo client explore query timeout")
			logrus.Error(err)
			break forLoop
		}
	}

	// if execution reaches that stage, there was an error
	if err == nil {
		// this should not happen
		err = errors.New("inconsistent state")
		logrus.Error(err)
	}
	return nil, err
}

// submitToNode sends a query to a node of the network, from the list of PIC-SURE resources
func (clientQuery *ExploreQuery) submitToNode(nodeIdx int) (result *models.ExploreQueryResultElement, err error) {
	logrus.Debug("Submitting query to node #", nodeIdx)

	params := medco_node.NewExploreQueryParamsWithTimeout(time.Duration(utilclient.QueryTimeoutSeconds) * time.Second)
	params.QueryRequest = medco_node.ExploreQueryBody{
		ID:    "MedCo_CLI_Query_" + time.Now().Format(time.RFC3339),
		Query: clientQuery.generateModel(),
	}

	response, err := clientQuery.httpMedCoClients[nodeIdx].MedcoNode.ExploreQuery(params, httptransport.BearerToken(clientQuery.authToken))
	if err != nil {
		logrus.Error("explore query error: ", err)
		return
	}

	return response.Payload.Result, nil
}

// generateModel parses the query terms and generate the model to be sent
func (clientQuery *ExploreQuery) generateModel() (queryModel *models.ExploreQuery) {

	// query model
	queryModel = &models.ExploreQuery{
		Type:          clientQuery.queryType,
		UserPublicKey: clientQuery.userPublicKey,
		Panels:        []*models.Panel{},
	}

	// query terms
	true := true
	for panelIdx, panel := range clientQuery.encPanelsItemKeys {

		panelModel := &models.Panel{
			Items: []*models.PanelItemsItems0{},
			Not:   &clientQuery.panelsIsNot[panelIdx],
		}

		for _, encItem := range panel {
			encrypted := new(bool)
			*encrypted = true
			queryTerm := new(string)
			*queryTerm = encItem
			panelModel.Items = append(panelModel.Items, &models.PanelItemsItems0{
				Encrypted: encrypted,
				Operator:  "exists",
				QueryTerm: queryTerm,
				Value:     "",
			})

		}

		queryModel.Panels = append(queryModel.Panels, panelModel)
	}

	return
}
