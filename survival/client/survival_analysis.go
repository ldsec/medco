package survivalclient

//for client !!!
import (
	"crypto/tls"
	"encoding/csv"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/ldsec/medco-connector/restapi/client/survival_analysis"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/ldsec/medco-connector/restapi/client"
	"github.com/ldsec/medco-connector/restapi/client/medco_network"
	utilclient "github.com/ldsec/medco-connector/util/client"
	"github.com/ldsec/medco-connector/wrappers/unlynx"
	"github.com/sirupsen/logrus"
)

//SurvivalAnalysis represent a survival analysis requeset
type SurvivalAnalysis struct {
	// httpMedCoClient is the HTTP client for the MedCo connector
	httpMedCoClients []*client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

	id string

	patientSetID int

	startConcept        string
	startColumn         string
	endConcept          string
	endColumn           string
	subGroupDefinitions []*survival_analysis.SubGroupDefinitionsItems0

	userPublicKey string

	userPrivateKey string

	formats strfmt.Registry

	timers map[string]time.Duration

	profilingBuffer BufferToPrint

	profiling *csv.Writer
}

// NewSurvivalAnalysis constructor for survival analysis request
func NewSurvivalAnalysis(token string, patientSetID int, subGroupDefinitions []*survival_analysis.SubGroupDefinitionsItems0, startConcept, startColumn, endConcept, endColumn string, disableTLSCheck bool) (q *SurvivalAnalysis, err error) {
	q = &SurvivalAnalysis{
		authToken:           token,
		id:                  "MedCo_Survival_Analysis" + time.Now().Format(time.RFC3339),
		patientSetID:        patientSetID,
		subGroupDefinitions: subGroupDefinitions,
		startConcept:        startConcept,
		startColumn:         startColumn,
		endConcept:          endConcept,
		endColumn:           endColumn,
		formats:             strfmt.Default,
		timers:              make(map[string]time.Duration),
	}

	q.profiling = csv.NewWriter(&q.profilingBuffer)
	q.profiling.Write([]string{"label", "value_in_seconds", "node_index"})

	parsedURL, err := url.Parse(utilclient.MedCoConnectorURL)
	if err != nil {
		logrus.Error("cannot parse MedCo connector URL: ", err)
		return
	}

	transport := httptransport.New(parsedURL.Host, parsedURL.Path, []string{parsedURL.Scheme})
	transport.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: disableTLSCheck}

	//q.httpMedCoClient = client.New(transport, nil)

	getMetadataResp, err := client.New(transport, nil).MedcoNetwork.GetMetadata(
		medco_network.NewGetMetadataParamsWithTimeout(30*time.Second),
		httptransport.BearerToken(token),
	)
	if err != nil {
		logrus.Error("get network metadata request error: ", err)
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

// EncryptedResults holds a TimePoint and the corresponding encrypted events
type EncryptedResults struct {
	TimePoint string
	Events    Events
}

//Events holds the number of events at a given timepoint in base64 strings
type Events struct {
	EventsOfInterest string
	CensoringEvents  string
}

//Decrypt deciphers a value that is expected to be encrypted under user public key
func (clientSurvivalAnalysis *SurvivalAnalysis) Decrypt(value string) (int64, error) {
	return unlynx.Decrypt(value, clientSurvivalAnalysis.userPrivateKey)
}

type NodeResult struct {
	Body      *survival_analysis.SurvivalAnalysisOKBody
	NodeIndex int
}

//Execute is the main function that sends the request and waits
func (clientSurvivalAnalysis *SurvivalAnalysis) Execute() (results map[string][]*EncryptedResults, err error) {
	totalTimer := time.Now()
	defer func(since time.Time) {
		err = clientSurvivalAnalysis.addTimer("time for the total execution ", since)
	}(totalTimer)
	errChan := make(chan error)
	resultChan := make(chan NodeResult)

	for idx := range clientSurvivalAnalysis.httpMedCoClients {

		go func(idx int) {
			res, Error := clientSurvivalAnalysis.submitToNode(idx)
			if Error != nil {
				logrus.Errorf("Survival analysis exection error : %s", Error)
				errChan <- Error
			} else {

				resultChan <- NodeResult{Body: res, NodeIndex: idx}
			}
		}(idx)
	}
	//TODO magic number
	timeout := time.After(time.Duration(3000) * time.Second)

	results = make(map[string][]*EncryptedResults)
nodeLoop:
	for idx := range clientSurvivalAnalysis.httpMedCoClients {
		select {
		case nodeLoopErr := <-errChan:
			if err != nil {
				err = fmt.Errorf("Node %d threw %s : %s", idx, nodeLoopErr, err)
			} else {
				err = fmt.Errorf("Node %d threw %s", idx, nodeLoopErr)
			}
		case nodeLoopRes := <-resultChan:
			for label, value := range nodeLoopRes.Body.Timers {

				timerErr := clientSurvivalAnalysis.profiling.Write([]string{label, strconv.FormatFloat(value, 'f', -1, 64), strconv.Itoa(nodeLoopRes.NodeIndex)})
				if timerErr != nil {
					err = fmt.Errorf("Node %d threw %s : %s", idx, timerErr, err)
					return
				}

			}
			clientSurvivalAnalysis.profiling.Flush()

			for _, groups := range nodeLoopRes.Body.Results {

				var innerList []*EncryptedResults
				for _, val := range groups.GroupResults {
					innerList = append(innerList, &EncryptedResults{TimePoint: val.Timepoint,
						Events: Events{EventsOfInterest: val.Events.Eventofinterest,
							CensoringEvents: val.Events.Censoringevent,
						}})
				}

				results[groups.GroupID] = append(results[groups.GroupID], innerList...)

			}

		case <-timeout:
			err = fmt.Errorf(" Timeout : %s", err.Error())
			break nodeLoop

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

//TODO test if it copies or not the map, I think it does not copy

// GetTimers returns the timers of the SurvivalAnalysis
func (clientSurvivalAnalysis *SurvivalAnalysis) GetTimers() map[string]time.Duration {
	return clientSurvivalAnalysis.timers
}

// PrintTimers prints the timers in the standard output if debug is enabled
func (clientSurvivalAnalysis *SurvivalAnalysis) PrintTimers() {
	for label, duration := range clientSurvivalAnalysis.timers {
		logrus.Debug(label + duration.String())
	}
}

func (clientSurvivalAnalysis *SurvivalAnalysis) DumpTimers() error {
	for label, value := range clientSurvivalAnalysis.timers {
		durationSeconds := value.Seconds()
		err := clientSurvivalAnalysis.profiling.Write([]string{label, strconv.FormatFloat(durationSeconds, 'f', -1, 64), "-1"})
		if err != nil {
			return err
		}
	}
	clientSurvivalAnalysis.profiling.Flush()

	return nil

}
