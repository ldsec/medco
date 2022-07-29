package querytoolsclient

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	medcoclient "github.com/ldsec/medco/connector/client"
	"github.com/ldsec/medco/connector/restapi/client"
	"github.com/ldsec/medco/connector/restapi/client/medco_node"
	utilclient "github.com/ldsec/medco/connector/util/client"
	"github.com/sirupsen/logrus"
)

// PutDefaultCohort is a MedCo client query for setting the default cohort
type PutDefaultCohort struct {
	// httpMedCoClients are the HTTP clients for the MedCo connectors
	httpMedCoClients []*client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

	cohortName string
}

// NewPutDefaultCohort creates a new put default cohort query
func NewPutDefaultCohort(token, cohortName string, disableTLSCheck bool) (putDefaultCohort *PutDefaultCohort, err error) {
	putDefaultCohort = &PutDefaultCohort{
		authToken:  token,
		cohortName: cohortName,
	}

	getMetadataResp, err := medcoclient.MetaData(token, disableTLSCheck)
	if err != nil {
		logrus.Error(err)
		return
	}

	putDefaultCohort.httpMedCoClients = make([]*client.MedcoCli, len(getMetadataResp.Payload.Nodes))
	for _, node := range getMetadataResp.Payload.Nodes {
		if putDefaultCohort.httpMedCoClients[*node.Index] != nil {
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
		putDefaultCohort.httpMedCoClients[*node.Index] = client.New(nodeTransport, nil)
	}
	return
}

// Execute executes the put default cohort request
func (putDefaultCohort *PutDefaultCohort) Execute() (err error) {
	nOfNodes := len(putDefaultCohort.httpMedCoClients)
	errChan := make(chan error)
	resultChan := make(chan *medco_node.PutDefaultCohortOK, nOfNodes)
	logrus.Infof("There are %d nodes", nOfNodes)

	for idx := 0; idx < nOfNodes; idx++ {

		go func(idx int) {
			logrus.Infof("Submitting to node %d", idx)
			res, Error := putDefaultCohort.submitToNode(idx)
			if Error != nil {
				logrus.Errorf("Query tool execution error : %s", Error)
				errChan <- Error
			} else {
				logrus.Infof("Node %d successfully put default cohort", idx)

				resultChan <- res
			}
		}(idx)
	}

	timeout := time.After(time.Duration(utilclient.QueryToolsTimeoutSeconds) * time.Second)
	for idx := 0; idx < nOfNodes; idx++ {
		select {
		case err = <-errChan:
			return
		case <-timeout:
			err = fmt.Errorf("Timeout %d seconds elapsed", utilclient.QueryToolsTimeoutSeconds)
			logrus.Error(err)
			return
		case <-resultChan:
			logrus.Infof("Node %d successfully updated default cohort", idx)
		}

	}

	logrus.Info("Operation completed")

	return
}

func (putDefaultCohort *PutDefaultCohort) submitToNode(nodeIdx int) (*medco_node.PutDefaultCohortOK, error) {
	params := medco_node.NewPutDefaultCohortParamsWithTimeout(time.Duration(utilclient.QueryTimeoutSeconds) * time.Second)

	params.SetName(putDefaultCohort.cohortName)

	response, err := putDefaultCohort.httpMedCoClients[nodeIdx].MedcoNode.PutDefaultCohort(params, httptransport.BearerToken(putDefaultCohort.authToken))
	if err != nil {
		return nil, err
	}
	return response, nil

}
