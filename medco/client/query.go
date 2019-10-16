package medcoclient

import (
	"crypto/tls"
	"errors"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ldsec/medco-connector/restapi/client"
	"github.com/ldsec/medco-connector/restapi/client/picsure2"
	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/ldsec/medco-connector/unlynx"
	utilclient "github.com/ldsec/medco-connector/util/client"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

// Query is a MedCo client query
type Query struct {

	// httpMedCoClients are the HTTP clients for the MedCo connectors
	httpMedCoClients []*client.MedcoCli
	// httpPicsureClient is the HTTP client for the MedCo connectors through PIC-SURE
	httpPicsureClient *client.MedcoCli
	// picsureResourceNames is the list of PIC-SURE resources names corresponding to the MedCo connector of each node
	picsureResourceNames []string
	// picsureResourceUUIDs is the list of PIC-SURE resources UUIDs corresponding to the MedCo connector of each node
	picsureResourceUUIDs []string
	// authToken is the OIDC authentication JWT
	authToken string
	// bypassPicsure instructs to query directly directly the medco connectors
	bypassPicsure bool

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
func NewQuery(authToken string, queryType models.QueryType, encPanelsItemKeys [][]string, panelsIsNot []bool, disableTLSCheck bool, bypassPicsure bool) (q *Query, err error) {

	transport := httptransport.New(utilclient.Picsure2APIHost, utilclient.Picsure2APIBasePath, []string{utilclient.Picsure2APIScheme})
	transport.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: disableTLSCheck}

	q = &Query{

		httpPicsureClient:    client.New(transport, nil),
		picsureResourceNames: utilclient.Picsure2Resources,
		picsureResourceUUIDs: []string{},
		authToken:            authToken,
		bypassPicsure:	      bypassPicsure,

		queryType: queryType,
		encPanelsItemKeys: encPanelsItemKeys,
		panelsIsNot: panelsIsNot,
	}

	// retrieve resources information
	err = q.loadResources(disableTLSCheck)
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

// loadResources requests the list of PIC-SURE resources, and retrieves their UUID from the names configured, and their URL
func (clientQuery *Query) loadResources(disableTLSCheck bool) (err error) {
	getResourcesParams := picsure2.NewGetResourcesParamsWithTimeout(30 * time.Second)
	resp, err := clientQuery.httpPicsureClient.Picsure2.GetResources(getResourcesParams, httptransport.BearerToken(clientQuery.authToken))
	if err != nil {
		logrus.Error("query error: ", err)
		return
	}

	// add in the same order the resources
	for _, resourceName := range clientQuery.picsureResourceNames {
		for _, respResource := range resp.Payload {
			if resourceName == respResource.Name {
				clientQuery.picsureResourceUUIDs = append(clientQuery.picsureResourceUUIDs, respResource.UUID)

				resourceUrl, err := url.Parse(respResource.ResourceRSPath)
				if err != nil {
					logrus.Error("error parsing URL: ", err)
					return err
				}

				resourceTransport := httptransport.New(resourceUrl.Host, resourceUrl.Path, []string{resourceUrl.Scheme})
				resourceTransport.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: disableTLSCheck}
				clientQuery.httpMedCoClients = append(clientQuery.httpMedCoClients, client.New(resourceTransport, nil))
			}
		}
	}

	// check all resources were indeed present
	if len(clientQuery.picsureResourceUUIDs) != len(clientQuery.picsureResourceNames) ||
		len(clientQuery.httpMedCoClients) != len(clientQuery.picsureResourceNames) {
		err = errors.New("some resources were not found")
		logrus.Error(err)
		return
	}

	return
}

// submitToNode sends a query to a node of the network, from the list of PIC-SURE resources
func (clientQuery *Query) submitToNode(picsureResourceIdx int) (result *models.QueryResultElement, err error) {
	logrus.Debug("Submitting query to resource ", clientQuery.picsureResourceNames[picsureResourceIdx],
		", bypass PIC-SURE: ", clientQuery.bypassPicsure)

	queryParams := picsure2.NewQuerySyncParamsWithTimeout(time.Duration(utilclient.QueryTimeoutSeconds) * time.Second)
	queryParams.Body = picsure2.QuerySyncBody{
		ResourceCredentials: &models.ResourceCredentials{
			MEDCOTOKEN: clientQuery.authToken,
		},
		ResourceUUID: clientQuery.picsureResourceUUIDs[picsureResourceIdx],
		Query: clientQuery.generateModel(),
	}

	var response *picsure2.QuerySyncOK
	if clientQuery.bypassPicsure {
		response, err = clientQuery.httpMedCoClients[picsureResourceIdx].Picsure2.QuerySync(queryParams, httptransport.BearerToken(clientQuery.authToken))
	} else {
		response, err = clientQuery.httpPicsureClient.Picsure2.QuerySync(queryParams, httptransport.BearerToken(clientQuery.authToken))
	}

	if err != nil {
		logrus.Error("query error: ", err)
		return
	}

	return response.Payload, nil
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
