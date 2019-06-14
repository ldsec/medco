package medcoclient

import (
	"crypto/tls"
	"errors"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/lca1/medco-connector/restapi/client"
	"github.com/lca1/medco-connector/restapi/client/picsure2"
	"github.com/lca1/medco-connector/restapi/models"
	"github.com/lca1/medco-connector/unlynx"
	utilclient "github.com/lca1/medco-connector/util/client"
	"github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// todo: times should be parsed

// Query is a MedCo client query
type Query struct {

	// httpMedcoClient is the HTTP client for the MedCo connectors through PIC-SURE
	httpMedcoClient *client.MedcoCli
	// picsureResourceNames is the list of PIC-SURE resources names corresponding to the MedCo connector of each node
	picsureResourceNames []string
	// picsureResourceUUIDs is the list of PIC-SURE resources UUIDs corresponding to the MedCo connector of each node
	picsureResourceUUIDs []string
	// authToken is the OIDC authentication JWT
	authToken string

	// userPublicKey is the user public key
	userPublicKey string
	// userPrivateKey is the user private key
	userPrivateKey string

	// queryType is the type of query requested
	queryType models.QueryType
	// encPanelsItemKeys is part of the query: contains the encrypted item keys organized by panel
	encPanelsItemKeys [][]string
	// panelsIsNot is part of the query: indicates which panels are negated
	panelsIsNot       []bool
}

// NewQuery creates a new MedCo client query
func NewQuery(authToken string, queryType models.QueryType, encPanelsItemKeys [][]string, panelsIsNot []bool, disableTLSCheck bool) (q *Query, err error) {

	transport := httptransport.New(utilclient.Picsure2APIHost, utilclient.Picsure2APIBasePath, []string{utilclient.Picsure2APIScheme})
	transport.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: disableTLSCheck}

	q = &Query{

		httpMedcoClient: client.New(transport, nil),
		picsureResourceNames: utilclient.Picsure2Resources,
		picsureResourceUUIDs: []string{},
		authToken: authToken,

		queryType: queryType,
		encPanelsItemKeys: encPanelsItemKeys,
		panelsIsNot: panelsIsNot,
	}

	// retrieve resource UUIDs
	err = q.loadResourceUUIDs()
	if err != nil {
		return
	}

	// generate ephemeral pair of user keys
	q.userPublicKey, q.userPrivateKey, err = unlynx.GenerateKeyPair()
	if err != nil {
		return
	}

	return
}

// Execute executes the MedCo client query synchronously on all the nodes through PIC-SURE
func (clientQuery *Query) Execute() (nodesResult map [string]*QueryResult, err error) {

	queryResultsChan := make(chan *models.QueryResultElement)
	queryErrChan := make(chan error)

	// execute requests on all nodes
	for idx := range clientQuery.picsureResourceUUIDs {
		go func() {

			queryResult, err := clientQuery.submitToNode(idx)
			if err != nil {
				queryErrChan <- err
			} else {
				queryResultsChan <- queryResult
			}

		}()
	}

	// parse the results as they come, or interrupt if one of them errors, or if a timeout occurs
	timeout := time.After(time.Duration(utilclient.QueryTimeoutSeconds) * time.Second)
	nodesResult = make(map [string]*QueryResult)
	forLoop: for _, picsureResourceName := range clientQuery.picsureResourceNames {
		select {
		case queryResult := <-queryResultsChan:
			parsedQueryResult, err := newQueryResult(queryResult, clientQuery.userPrivateKey)
			if err != nil {
				return nil, err
			}

			nodesResult[picsureResourceName] = parsedQueryResult
			logrus.Info("MedCo client query successful for resource ", picsureResourceName)

			if len(nodesResult) == len(clientQuery.picsureResourceUUIDs) {
				logrus.Info("MedCo client query successful for all resources")
				return nodesResult, nil
			}

		case err = <- queryErrChan:
			logrus.Error("MedCo client query error: ", err)
			break forLoop

		case <-timeout:
			err = errors.New("MedCo client query timeout")
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

// loadResourceUUIDs requests the list of PIC-SURE resources, and retrieves their UUID from the names configured
func (clientQuery *Query) loadResourceUUIDs() (err error) {
	getResourcesParams := picsure2.NewGetResourcesParamsWithTimeout(30 * time.Second)
	resp, err := clientQuery.httpMedcoClient.Picsure2.GetResources(getResourcesParams, httptransport.BearerToken(clientQuery.authToken))
	if err != nil {
		logrus.Error("query error: ", err)
		return
	}

	// add in the same order the resources
	for _, resourceName := range clientQuery.picsureResourceNames {
		for _, respResource := range resp.Payload {
			if resourceName == respResource.Name {
				clientQuery.picsureResourceUUIDs = append(clientQuery.picsureResourceUUIDs, respResource.UUID)
			}
		}
	}

	// check all resources were indeed present
	if len(clientQuery.picsureResourceUUIDs) != len(clientQuery.picsureResourceNames) {
		err = errors.New("some resources were not found")
		logrus.Error(err)
		return
	}

	return
}

// submitToNode sends a query to a node of the network, from the list of PIC-SURE resources
func (clientQuery *Query) submitToNode(picsureResourceUUIDIdx int) (result *models.QueryResultElement, err error) {

	queryParams := picsure2.NewQuerySyncParamsWithTimeout(time.Duration(utilclient.QueryTimeoutSeconds) * time.Second)
	queryParams.Body = picsure2.QuerySyncBody{
		ResourceCredentials: &models.ResourceCredentials{
			MEDCOTOKEN: clientQuery.authToken,
		},
		ResourceUUID: clientQuery.picsureResourceUUIDs[picsureResourceUUIDIdx],
		Query: clientQuery.generateModel(),
	}

	resp, err := clientQuery.httpMedcoClient.Picsure2.QuerySync(queryParams, httptransport.BearerToken(clientQuery.authToken))
	if err != nil {
		logrus.Error("query error: ", err)
		return
	}

	return resp.Payload, nil
}

// generateModel parses the query terms and generate the model to be sent
func (clientQuery *Query) generateModel() (queryModel *models.Query) {

	// query model
	queryModel = &models.Query{
		Name: "MedCo_CLI_Query_" + time.Now().Format(time.RFC3339),
		I2b2Medco: &models.QueryI2b2Medco{
			UserPublicKey: clientQuery.userPublicKey,
			QueryType: clientQuery.queryType,
			Panels: []*models.QueryI2b2MedcoPanelsItems0{},
		},
	}

	// query terms
	for panelIdx, panel := range clientQuery.encPanelsItemKeys {

		panelModel := &models.QueryI2b2MedcoPanelsItems0{
			Items: []*models.QueryI2b2MedcoPanelsItems0ItemsItems0{},
			Not: &clientQuery.panelsIsNot[panelIdx],
		}

		true := true
		for _, encItem := range panel {
			panelModel.Items = append(panelModel.Items, &models.QueryI2b2MedcoPanelsItems0ItemsItems0{
				Encrypted: &true,
				Operator:  "exists",
				QueryTerm: encItem,
				Value:     "",
			})

		}

		queryModel.I2b2Medco.Panels = append(queryModel.I2b2Medco.Panels, panelModel)
	}

	return
}
