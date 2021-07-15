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

// PostCohorts is a MedCo client query to update saved cohorts
type PostCohorts struct {
	// httpMedCoClients are the HTTP clients for the MedCo connectors
	httpMedCoClients []*client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

	patientSetID []int
	cohortName   string
}

// NewPostCohorts creates a new post cohorts query
func NewPostCohorts(token string, patientSetID []int, cohortName string, disableTLSCheck bool) (postCohorts *PostCohorts, err error) {
	postCohorts = &PostCohorts{
		authToken:    token,
		cohortName:   cohortName,
		patientSetID: patientSetID,
	}

	getMetadataResp, err := medcoclient.MetaData(token, disableTLSCheck)
	if err != nil {
		logrus.Error(err)
		return
	}

	nofNodes := len(getMetadataResp.Payload.Nodes)
	if len(patientSetID) != nofNodes {
		err = fmt.Errorf("number of provided patient set IDs must be the same as that of MedCo nodes: provided %d, connected nodes %d", len(patientSetID), nofNodes)
		logrus.Error(err)
		return
	}
	postCohorts.httpMedCoClients = make([]*client.MedcoCli, nofNodes)
	for _, node := range getMetadataResp.Payload.Nodes {
		if postCohorts.httpMedCoClients[*node.Index] != nil {
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
		postCohorts.httpMedCoClients[*node.Index] = client.New(nodeTransport, nil)
	}
	return
}

// Execute executes the post cohorts query
func (postCohorts *PostCohorts) Execute() (err error) {
	nOfNodes := len(postCohorts.httpMedCoClients)
	errChan := make(chan error)
	resultChan := make(chan *medco_node.PostCohortsOK, nOfNodes)
	logrus.Infof("There are %d nodes", nOfNodes)

	for idx := 0; idx < nOfNodes; idx++ {

		go func(idx int) {
			logrus.Infof("Submitting to node %d", idx)
			res, Error := postCohorts.submitToNode(idx)
			if Error != nil {
				logrus.Errorf("Query tool execution error : %s", Error)
				errChan <- Error
			} else {
				logrus.Infof("Node %d successfully posted cohort", idx)

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
			logrus.Infof("Node %d succesfully updated cohort", idx)
		}

	}
	logrus.Info("Operation completed")

	return
}

func (postCohorts *PostCohorts) submitToNode(nodeIdx int) (*medco_node.PostCohortsOK, error) {
	creationDate := time.Now()
	updateDate := time.Now()
	params := medco_node.NewPostCohortsParamsWithTimeout(time.Duration(utilclient.QueryTimeoutSeconds) * time.Second)
	body := medco_node.PostCohortsBody{
		CreationDate: new(string),
		UpdateDate:   new(string),
		QueryID:      new(int64),
	}

	*body.CreationDate = creationDate.Format(time.RFC3339)
	*body.UpdateDate = updateDate.Format(time.RFC3339)
	*body.QueryID = int64(postCohorts.patientSetID[nodeIdx])

	params.SetCohortsRequest(body)
	params.SetName(postCohorts.cohortName)

	response, err := postCohorts.httpMedCoClients[nodeIdx].MedcoNode.PostCohorts(params, httptransport.BearerToken(postCohorts.authToken))
	if err != nil {
		return nil, err
	}
	return response, nil

}
