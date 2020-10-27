package querytoolsclient

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	medcoclient "github.com/ldsec/medco/connector/client"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ldsec/medco/connector/restapi/client"
	"github.com/ldsec/medco/connector/restapi/client/medco_node"
	utilclient "github.com/ldsec/medco/connector/util/client"
	"github.com/sirupsen/logrus"
)

// RemoveCohorts is a MedCo client query to remove a saved cohort
type RemoveCohorts struct {
	// httpMedCoClients are the HTTP clients for the MedCo connectors
	httpMedCoClients []*client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

	cohortName string
}

// NewRemoveCohorts creates a new post cohorts query
func NewRemoveCohorts(token string, cohortName string, disableTLSCheck bool) (removeCohorts *RemoveCohorts, err error) {
	removeCohorts = &RemoveCohorts{
		authToken:  token,
		cohortName: cohortName,
	}

	getMetadataResp, err := medcoclient.MetaData(token, disableTLSCheck)
	if err != nil {
		logrus.Error(err)
		return
	}

	removeCohorts.httpMedCoClients = make([]*client.MedcoCli, len(getMetadataResp.Payload.Nodes))
	for _, node := range getMetadataResp.Payload.Nodes {
		if removeCohorts.httpMedCoClients[*node.Index] != nil {
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
		removeCohorts.httpMedCoClients[*node.Index] = client.New(nodeTransport, nil)
	}
	return
}

// Execute executes the remove cohorts query
func (removeCohorts *RemoveCohorts) Execute() (err error) {
	nOfNodes := len(removeCohorts.httpMedCoClients)
	errChan := make(chan error)
	resultChan := make(chan *medco_node.DeleteCohortsOK, nOfNodes)
	logrus.Infof("There are %d nodes", nOfNodes)

	for idx := 0; idx < nOfNodes; idx++ {

		go func(idx int) {
			logrus.Infof("Submitting to node %d", idx)
			res, Error := removeCohorts.submitToNode(idx)
			if Error != nil {
				logrus.Errorf("Query tool execution error : %s", Error)
				errChan <- Error
			} else {
				logrus.Infof("Node %d successfully removed cohort", idx)
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
			logrus.Infof("Node %d successfully deleted cohort", idx)
		}

	}

	logrus.Info("Operation completed")

	return

}

func (removeCohorts *RemoveCohorts) submitToNode(nodeIdx int) (*medco_node.DeleteCohortsOK, error) {

	params := medco_node.NewDeleteCohortsParamsWithTimeout(time.Duration(utilclient.QueryTimeoutSeconds) * time.Second)

	params.Name = removeCohorts.cohortName

	response, err := removeCohorts.httpMedCoClients[nodeIdx].MedcoNode.DeleteCohorts(params, httptransport.BearerToken(removeCohorts.authToken))
	if err != nil {
		return nil, err
	}
	return response, nil

}
