package querytoolsclient

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	medcomodels "github.com/ldsec/medco/connector/models"
	"github.com/ldsec/medco/connector/wrappers/unlynx"

	httptransport "github.com/go-openapi/runtime/client"
	medcoclient "github.com/ldsec/medco/connector/client"
	"github.com/ldsec/medco/connector/restapi/client"
	"github.com/ldsec/medco/connector/restapi/client/medco_node"
	utilclient "github.com/ldsec/medco/connector/util/client"
	"github.com/sirupsen/logrus"
)

// CohortsPatientList is a MedCo client query to get saved cohorts
type CohortsPatientList struct {
	// httpMedCoClients are the HTTP clients for the MedCo connectors
	httpMedCoClients []*client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

	id string

	cohortName string

	userPublicKey  string
	userPrivateKey string
}

// NewCohortsPatientList returns a new cohorts patient list
func NewCohortsPatientList(token, cohortName string, disableTLSCheck bool) (cohortsPatientList *CohortsPatientList, err error) {
	cohortsPatientList = &CohortsPatientList{
		authToken:  token,
		cohortName: cohortName,
		id:         "MedCo_Cohrts_Patient_List" + time.Now().Format(time.RFC3339),
	}

	getMetadataResp, err := medcoclient.MetaData(token, disableTLSCheck)
	if err != nil {
		logrus.Error(err)
		return
	}
	cohortsPatientList.httpMedCoClients = make([]*client.MedcoCli, len(getMetadataResp.Payload.Nodes))
	for _, node := range getMetadataResp.Payload.Nodes {
		if cohortsPatientList.httpMedCoClients[*node.Index] != nil {
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
		cohortsPatientList.httpMedCoClients[*node.Index] = client.New(nodeTransport, nil)
	}
	cohortsPatientList.userPublicKey, cohortsPatientList.userPrivateKey, err = unlynx.GenerateKeyPair()
	if err != nil {
		logrus.Error("while generating key pair: %s", err.Error())
		return
	}
	return
}

type patientListNodeResult = struct {
	patientList *medco_node.PostCohortsPatientListOKBody
	nodeIndex   int
}

// Execute executes the patient list retrieval and decryption
func (pl *CohortsPatientList) Execute() (patientLists [][]int64, nodeTimers []medcomodels.Timers, localTimers medcomodels.Timers, err error) {
	nOfNodes := len(pl.httpMedCoClients)
	errChan := make(chan error)
	nodeTimers = make([]medcomodels.Timers, nOfNodes)
	resultChan := make(chan patientListNodeResult, nOfNodes)
	cipherPatientLists := make([][]string, nOfNodes)
	patientLists = make([][]int64, nOfNodes)
	logrus.Infof("There are %d nodes", nOfNodes)

	for i := 0; i < nOfNodes; i++ {
		go func(idx int) {
			resultBody, err := pl.submitToNode(idx)
			if err != nil {
				errChan <- err
				return
			}
			resultChan <- patientListNodeResult{resultBody, idx}
			return
		}(i)
	}
	timeout := time.After(time.Duration(utilclient.QueryToolsTimeoutSeconds) * time.Second)
	for i := 0; i < nOfNodes; i++ {
		select {
		case err = <-errChan:
			return
		case <-timeout:
			err = fmt.Errorf("Timeout %d seconds elapsed", utilclient.QueryToolsTimeoutSeconds)
			logrus.Error(err)
			return
		case nodeRes := <-resultChan:
			logrus.Infof("Operation successfully returns in node %d.", nodeRes.nodeIndex)
			cipherPatientLists[nodeRes.nodeIndex] = nodeRes.patientList.Results
			nodeTimers[nodeRes.nodeIndex] = medcomodels.NewTimersFromAPIModel(nodeRes.patientList.Timers)
		}
	}

	//Decryption

	logrus.Info("All lists receiveds. Decrypting")
	localTimers = medcomodels.NewTimers()
	timer := time.Now()
	for i, list := range cipherPatientLists {
		patientLists[i] = make([]int64, len(list))
		for j, encNum := range list {
			patientLists[i][j], err = unlynx.Decrypt(encNum, pl.userPrivateKey)
			if err != nil {
				err = fmt.Errorf("during local decryption: %s", err.Error())
				logrus.Error(err)
				return
			}
		}
	}
	localTimers.AddTimers("local-decryption-of-patient-lists", timer, nil)
	return

}

func (pl *CohortsPatientList) submitToNode(nodeIdx int) (*medco_node.PostCohortsPatientListOKBody, error) {

	params := medco_node.NewPostCohortsPatientListParamsWithTimeout(time.Duration(utilclient.QueryTimeoutSeconds) * time.Second)
	body := medco_node.PostCohortsPatientListBody{
		CohortName:    &pl.cohortName,
		ID:            pl.id,
		UserPublicKey: &pl.userPublicKey,
	}
	params.SetCohortsPatientListRequest(body)

	response, err := pl.httpMedCoClients[nodeIdx].MedcoNode.PostCohortsPatientList(params, httptransport.BearerToken(pl.authToken))
	if err != nil {
		return nil, err
	}
	return response.GetPayload(), nil

}
