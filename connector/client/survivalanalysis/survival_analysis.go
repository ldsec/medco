package survivalclient

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/ldsec/medco/connector/restapi/client/survival_analysis"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/ldsec/medco/connector/restapi/client"
	utilclient "github.com/ldsec/medco/connector/util/client"
	utilcommon "github.com/ldsec/medco/connector/util/common"
	"github.com/ldsec/medco/connector/wrappers/unlynx"
	"github.com/sirupsen/logrus"
)

const timeOutInSeconds int64 = 3000

//SurvivalAnalysis represent a survival analysis requeset
type SurvivalAnalysis struct {
	// httpMedCoClient is the HTTP client for the MedCo connector
	httpMedCoClients []*client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

	id string

	patientSetID int

	startConceptPath    string
	startModifierCode   string
	endConceptPath      string
	endModifierCode     string
	subGroupDefinitions []*survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0

	limit       int
	granularity string

	userPublicKey  string
	userPrivateKey string

	formats strfmt.Registry

	timers map[string]time.Duration
}

// NewSurvivalAnalysis constructor for survival analysis request
func NewSurvivalAnalysis(token string, patientSetID int, subGroupDefinitions []*survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0, limit int, granularity, startConcept, startModifier, endConcept, endModifier string, disableTLSCheck bool) (q *SurvivalAnalysis, err error) {
	q = &SurvivalAnalysis{
		authToken:           token,
		id:                  "MedCo_Survival_Analysis" + time.Now().Format(time.RFC3339),
		patientSetID:        patientSetID,
		subGroupDefinitions: subGroupDefinitions,
		startConceptPath:    startConcept,
		startModifierCode:   startModifier,
		endConceptPath:      endConcept,
		endModifierCode:     endModifier,
		limit:               limit,
		granularity:         granularity,
		formats:             strfmt.Default,
		timers:              make(map[string]time.Duration),
	}

	getMetadataResp, err := utilclient.MetaData(token, disableTLSCheck)
	if err != nil {
		logrus.Error(err)
		return
	}

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

	q.userPublicKey, q.userPrivateKey, err = unlynx.GenerateKeyPair()
	if err != nil {
		return
	}

	return

}

//Decrypt deciphers a value that is expected to be encrypted under user public key
func (clientSurvivalAnalysis *SurvivalAnalysis) Decrypt(value string) (int64, error) {
	return unlynx.Decrypt(value, clientSurvivalAnalysis.userPrivateKey)
}

type nodeResult struct {
	Body      *survival_analysis.SurvivalAnalysisOKBody
	NodeIndex int
}

//Execute makes a call to API for survival analysis,
func (clientSurvivalAnalysis *SurvivalAnalysis) Execute() (results []EncryptedResults, nodeTimers []utilcommon.Timers, err error) {

	var nOfNodes = len(clientSurvivalAnalysis.httpMedCoClients)
	errChan := make(chan error)
	resultChan := make(chan nodeResult, nOfNodes)
	results = make([]EncryptedResults, nOfNodes)
	nodeTimers = make([]utilcommon.Timers, nOfNodes)

	for idx := 0; idx < nOfNodes; idx++ {

		go func(idx int) {
			res, Error := clientSurvivalAnalysis.submitToNode(idx)
			if Error != nil {
				logrus.Errorf("Survival analysis exection error : %s", Error)
				errChan <- Error
			} else {

				resultChan <- nodeResult{Body: res, NodeIndex: idx}
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
			results[nodeRes.NodeIndex] = encryptedResultsFromAPIResponse(nodeRes.Body.Results)
			nodeTimers[nodeRes.NodeIndex] = utilcommon.APIModelToTimers(nodeRes.Body.Timers)
		}
	}

	return
}

func (clientSurvivalAnalysis *SurvivalAnalysis) addTimer(label string, since time.Time) (err error) {
	if _, exists := clientSurvivalAnalysis.timers[label]; exists {
		err = errors.New("Timer label " + label + " already exists")
		return
	}
	clientSurvivalAnalysis.timers[label] = time.Since(since)
	return

}

func (clientSurvivalAnalysis *SurvivalAnalysis) addTimers(timers map[string]time.Duration) (err error) {
	for label, duration := range timers {
		if _, exists := clientSurvivalAnalysis.timers[label]; exists {
			err = errors.New("Timer label " + label + " already exists")
			return
		}
		clientSurvivalAnalysis.timers[label] = duration
	}
	return

}
