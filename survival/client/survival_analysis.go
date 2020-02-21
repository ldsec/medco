package survivalclient

//for client !!!
import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
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

//for medco cli
type SurvivalAnalysis struct {

	// httpMedCoClient is the HTTP client for the MedCo connector
	httpMedCoClients []*client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

	patientSetIDs map[int]string

	timeCodes []string

	userPublicKey string

	userPrivateKey string

	formats strfmt.Registry
}

func NewSurvivalAnalysis(token string, patientSetIDs map[int]string, timeCodes []string, disableTLSCheck bool) (q *SurvivalAnalysis, err error) {
	q = &SurvivalAnalysis{
		authToken:     token,
		patientSetIDs: patientSetIDs,
		timeCodes:     timeCodes,
		formats:       strfmt.Default,
	}

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

type EncryptedResults struct {
	TimePoint string
	Events    Events
}

type Events struct {
	EventsOfInterest string
	CensoringEvents  string
}

func (clientSurvivalAnalysis *SurvivalAnalysis) Execute() (results []*EncryptedResults, err error) {

	errChan := make(chan error)
	resultChan := make(chan []*survival_analysis.GetSurvivalAnalysisOKBodyItems0)

	for idx := range clientSurvivalAnalysis.httpMedCoClients {

		go func(idx int) {
			res, Error := clientSurvivalAnalysis.submitToNode(idx)
			if Error != nil {
				errChan <- Error
			}
			resultChan <- res
		}(idx)
	}
	//TODO magic number
	timeout := time.After(time.Duration(3000) * time.Second)

nodeLoop:
	for idx := range clientSurvivalAnalysis.httpMedCoClients {
		select {
		case nodeLoopRes := <-resultChan:
			concatEncryptedResults(results, nodeLoopRes)

		case nodeLoopErr := <-errChan:
			err = fmt.Errorf("Node %d threw %s : %s", idx, nodeLoopErr.Error(), err.Error())
		case <-timeout:
			err = fmt.Errorf(" Timeout : %s", err.Error())
			break nodeLoop

		}
	}

	return
}

func concatEncryptedResults(target []*EncryptedResults, toExctract []*survival_analysis.GetSurvivalAnalysisOKBodyItems0) []*EncryptedResults {
	var res = target
	for _, bodyItem := range toExctract {
		result := &EncryptedResults{
			TimePoint: bodyItem.Timepoint,
			Events: Events{
				EventsOfInterest: bodyItem.Events.Eventofinterest,
				CensoringEvents:  bodyItem.Events.Censoringevent,
			},
		}
		res = append(res, result)
	}

	return res
}
