package common

//for client !!!
import (
	"crypto/tls"
	"net/http"
	"net/url"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/ldsec/medco-connector/restapi/client"
	utilclient "github.com/ldsec/medco-connector/util/client"
	"github.com/sirupsen/logrus"
)

//for medco cli
type SurvivalAnalysis struct {

	// httpMedCoClient is the HTTP client for the MedCo connector
	httpMedCoClient *client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

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

	q.httpMedCoClient = client.New(transport, nil)
	return
}

func (clientSurvivalAnalysis *SurvivalAnalysis) Execute() (results []string, err error) {
	results, err = clientSurvivalAnalysis.submitToNode()
	return
}
