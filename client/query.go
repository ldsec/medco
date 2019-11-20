package medcoclient

import (
	"crypto/tls"
	"errors"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ldsec/medco-connector/restapi/client"
	"github.com/ldsec/medco-connector/restapi/client/medco_network"
	"github.com/ldsec/medco-connector/restapi/client/medco_node"
	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/ldsec/medco-connector/util/client"
	"github.com/ldsec/medco-connector/wrappers/unlynx"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
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
	panelsIsNot       []bool
}

// NewExploreQuery creates a new MedCo client query
func NewExploreQuery(authToken string, queryType models.ExploreQueryType, encPanelsItemKeys [][]string, panelsIsNot []bool, disableTLSCheck bool) (q *ExploreQuery, err error) {

	q = &ExploreQuery{
		authToken:            authToken,
		queryType: queryType,
		encPanelsItemKeys: encPanelsItemKeys,
		panelsIsNot: panelsIsNot,
	}

	// retrieve network information
	parsedUrl, err := url.Parse(utilclient.MedCoConnectorURL)
	if err != nil {
		logrus.Error("cannot parse MedCo connector URL: ", err)
		return
	}

	transport := httptransport.New(parsedUrl.Host, parsedUrl.Path, []string{parsedUrl.Scheme})
	transport.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: disableTLSCheck}

	getMetadataResp, err := client.New(transport, nil).MedcoNetwork.GetMetadata(
		medco_network.NewGetMetadataParamsWithTimeout(30 * time.Second),
		httptransport.BearerToken(authToken),
	)
	if err != nil {
		logrus.Error("get network metadata request error: ", err)
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

		nodeUrl, err := url.Parse(node.URL)
		if err != nil {
			logrus.Error("cannot parse MedCo connector URL: ", err)
			return nil, err
		}

		nodeTransport := httptransport.New(nodeUrl.Host, nodeUrl.Path, []string{nodeUrl.Scheme})
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

// Execute executes the MedCo client query synchronously on all the nodes
func (clientQuery *ExploreQuery) Execute() (nodesResult map [int]*ExploreQueryResult, err error) {

	queryResultsChan := make(chan *models.ExploreQueryResultElement)
	queryErrChan := make(chan error)

	// execute requests on all nodes
	for idx := range clientQuery.httpMedCoClients {
		idxLocal := idx
		go func() {

			queryResult, err := clientQuery.submitToNode(idxLocal)
			if err != nil {
				queryErrChan <- err
			} else {
				queryResultsChan <- queryResult
			}

		}()
	}

	// parse the results as they come, or interrupt if one of them errors, or if a timeout occurs
	timeout := time.After(time.Duration(utilclient.QueryTimeoutSeconds) * time.Second)
	nodesResult = make(map [int]*ExploreQueryResult)
	forLoop: for nodeIdx := range clientQuery.httpMedCoClients {
		select {
		case queryResult := <-queryResultsChan:
			parsedQueryResult, err := newQueryResult(queryResult, clientQuery.userPrivateKey)
			if err != nil {
				return nil, err
			}

			nodesResult[nodeIdx] = parsedQueryResult
			logrus.Info("MedCo client explore query successful for node ", nodeIdx)

			if len(nodesResult) == len(clientQuery.httpMedCoClients) {
				logrus.Info("MedCo client explore query successful for all resources")
				return nodesResult, nil
			}

		case err = <- queryErrChan:
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
		ID	 : "MedCo_CLI_Query_" + time.Now().Format(time.RFC3339),
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
		Type: clientQuery.queryType,
		UserPublicKey: clientQuery.userPublicKey,
		Panels: []*models.ExploreQueryPanelsItems0{},
	}

	// query terms
	true := true
	for panelIdx, panel := range clientQuery.encPanelsItemKeys {

		panelModel := &models.ExploreQueryPanelsItems0{
			Items: []*models.ExploreQueryPanelsItems0ItemsItems0{},
			Not: &clientQuery.panelsIsNot[panelIdx],
		}

		for _, encItem := range panel {
			panelModel.Items = append(panelModel.Items, &models.ExploreQueryPanelsItems0ItemsItems0{
				Encrypted: &true,
				Operator:  "exists",
				QueryTerm: encItem,
				Value:     "",
			})

		}

		queryModel.Panels = append(queryModel.Panels, panelModel)
	}

	return
}
