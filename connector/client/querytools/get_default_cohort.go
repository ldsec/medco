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

// GetDefaultCohort is a MedCo client query for setting the default cohort
type GetDefaultCohort struct {
	// httpMedCoClients are the HTTP clients for the MedCo connectors
	httpMedCoClients []*client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string
}

// NewGetDefaultCohort creates a new put default cohort query
func NewGetDefaultCohort(token string, disableTLSCheck bool) (getDefaultCohort *GetDefaultCohort, err error) {
	getDefaultCohort = &GetDefaultCohort{
		authToken: token,
	}

	getMetadataResp, err := medcoclient.MetaData(token, disableTLSCheck)
	if err != nil {
		logrus.Error(err)
		return
	}

	// TODO this portion of code is redundant in all functions that create an API request, a generic function (type) should be created
	getDefaultCohort.httpMedCoClients = make([]*client.MedcoCli, len(getMetadataResp.Payload.Nodes))
	for _, node := range getMetadataResp.Payload.Nodes {
		if getDefaultCohort.httpMedCoClients[*node.Index] != nil {
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
		getDefaultCohort.httpMedCoClients[*node.Index] = client.New(nodeTransport, nil)
	}
	return
}

type getDefaultCohortNodeResult = struct {
	defaultCohort *medco_node.GetDefaultCohortOK
	nodeIndex     int
}

// Execute executes the get default cohort request
func (getDefaultCohort *GetDefaultCohort) Execute() (defaultCohortNames []string, err error) {
	nOfNodes := len(getDefaultCohort.httpMedCoClients)
	defaultCohortNames = make([]string, nOfNodes)
	errChan := make(chan error)
	resultChan := make(chan getDefaultCohortNodeResult, nOfNodes)
	logrus.Infof("There are %d nodes", nOfNodes)

	for idx := 0; idx < nOfNodes; idx++ {

		go func(idx int) {
			logrus.Infof("Submitting to node %d", idx)
			res, Error := getDefaultCohort.submitToNode(idx)
			if Error != nil {
				logrus.Errorf("Query tool execution error : %s", Error)
				errChan <- Error
			} else {
				logrus.Infof("Node %d successfully put default cohort", idx)

				resultChan <- getDefaultCohortNodeResult{res, idx}

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
		case nodeRes := <-resultChan:
			logrus.Infof("Node %d successfully got default cohort", idx)
			defaultCohortNames[nodeRes.nodeIndex] = nodeRes.defaultCohort.Payload
		}

	}

	logrus.Info("Operation completed")

	return
}

func (getDefaultCohort *GetDefaultCohort) submitToNode(nodeIdx int) (*medco_node.GetDefaultCohortOK, error) {
	params := medco_node.NewGetDefaultCohortParamsWithTimeout(time.Duration(utilclient.QueryTimeoutSeconds) * time.Second)

	response, err := getDefaultCohort.httpMedCoClients[nodeIdx].MedcoNode.GetDefaultCohort(params, httptransport.BearerToken(getDefaultCohort.authToken))
	if err != nil {
		return nil, err
	}
	return response, nil

}
