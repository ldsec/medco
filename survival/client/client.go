package survivalclient

//for client !!!
import (
	"crypto/tls"
	"errors"
	"net/http"
	"net/url"
	"time"

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

	userPublicKey string

	userPrivateKey string

	granularity string

	limit int64
	//TODO patient list []

	formats strfmt.Registry
}

func NewSurvivalAnalysis(token, granularity string, limit int64, disableTLSCheck bool) (q *SurvivalAnalysis, err error) {
	q = &SurvivalAnalysis{
		authToken:   token,
		granularity: granularity,
		limit:       limit,
		formats:     strfmt.Default,
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

func (clientSurvivalAnalysis *SurvivalAnalysis) Execute() (results []string, err error) {
	for idx := range clientSurvivalAnalysis.httpMedCoClients {
		clientSurvivalAnalysis.submitToNode(idx)
	}
	err = errors.New("TODO finihssh to impleent the results for the client side ")
	return
}
