package querytoolsclient

import (
	"github.com/ldsec/medco/connector/restapi/client"
	"github.com/ldsec/medco/connector/restapi/client/medco_node"

	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	utilclient "github.com/ldsec/medco/connector/util/client"
	utilcommon "github.com/ldsec/medco/connector/util/common"

	httptransport "github.com/go-openapi/runtime/client"

	"github.com/sirupsen/logrus"
)

// GetCohorts is a MedCo client query to get saved cohorts
type GetCohorts struct {
	// httpMedCoClients are the HTTP clients for the MedCo connectors
	httpMedCoClients []*client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string
}

const timeOutInSeconds = 30

// NewGetCohorts creates a new get cohorts query
func NewGetCohorts(token string, disableTLSCheck bool) (getCohorts *GetCohorts, err error) {
	getCohorts = &GetCohorts{
		authToken: token,
	}

	getMetadataResp, err := utilclient.MetaData(token, disableTLSCheck)

	getCohorts.httpMedCoClients = make([]*client.MedcoCli, len(getMetadataResp.Payload.Nodes))
	for _, node := range getMetadataResp.Payload.Nodes {
		if getCohorts.httpMedCoClients[*node.Index] != nil {
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
		getCohorts.httpMedCoClients[*node.Index] = client.New(nodeTransport, nil)
	}
	return
}

type nodeResult = struct {
	cohorts   []*medco_node.GetCohortsOKBodyItems0
	nodeIndex int
}

// Execute executes the get cohorts query
func (getCohorts *GetCohorts) Execute() (results [][]utilcommon.Cohort, err error) {
	nOfNodes := len(getCohorts.httpMedCoClients)
	errChan := make(chan error)
	resultChan := make(chan nodeResult, nOfNodes)
	results = make([][]utilcommon.Cohort, nOfNodes)

	for idx := 0; idx < nOfNodes; idx++ {

		go func(idx int) {
			res, Error := getCohorts.submitToNode(idx)
			if Error != nil {
				logrus.Errorf("Survival analysis exection error : %s", Error)
				errChan <- Error
			} else {

				resultChan <- nodeResult{cohorts: res, nodeIndex: idx}
			}
		}(idx)
	}

	timeout := time.After(time.Duration(timeOutInSeconds) * time.Second)
	for idx := 0; idx < nOfNodes; idx++ {
		select {
		case err = <-errChan:
			return
		case <-timeout:
			err = fmt.Errorf("Timeout %d seconds elapsed", timeOutInSeconds)
			logrus.Error(err)
			return
		case nodeRes := <-resultChan:
			results[nodeRes.nodeIndex], err = convertCohort(nodeRes.cohorts)
			if err != nil {
				err = fmt.Errorf("cohort format in results: %s", err.Error())
				logrus.Error(err)
				return
			}

		}
	}
	return
}

func (getCohorts *GetCohorts) submitToNode(nodeIdx int) ([]*medco_node.GetCohortsOKBodyItems0, error) {
	//TODO timeout
	params := medco_node.NewGetCohortsParamsWithTimeout(time.Duration(utilclient.GenomicAnnotationsQueryTimeoutSeconds) * time.Second)
	response, err := getCohorts.httpMedCoClients[nodeIdx].MedcoNode.GetCohorts(params, httptransport.BearerToken(getCohorts.authToken))
	if err != nil {
		return nil, err
	}
	return response.GetPayload(), nil

}

//TODO int64
func convertCohort(apiRes []*medco_node.GetCohortsOKBodyItems0) (res []utilcommon.Cohort, err error) {
	res = make([]utilcommon.Cohort, len(apiRes))
	for i, apiCohort := range apiRes {
		res[i].CohortID = int(apiCohort.CohortID)
		res[i].CohortName = apiCohort.CohortName
		res[i].CreationDate, err = time.Parse(time.RFC3339, apiCohort.CreationDate)
		if err != nil {
			return
		}
		res[i].UpdateDate, err = time.Parse(time.RFC3339, apiCohort.UpdateDate)
		if err != nil {
			return
		}
		res[i].QueryID = int(apiCohort.QueryID)
	}
	return
}
