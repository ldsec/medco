package medcoclient

import (
	"crypto/tls"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ldsec/medco-connector/restapi/client"
	"github.com/ldsec/medco-connector/restapi/client/medco_node"
	utilclient "github.com/ldsec/medco-connector/util/client"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

// GetStatus is a MedCo client get-status request
type GetStatus struct {

	// httpMedCoClient is the HTTP client for the MedCo connector
	httpMedCoClient *client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string
}

// NewGetStatus creates a new MedCo client get-status request
func NewGetStatus(authToken string, disableTLSCheck bool) (q *GetStatus, err error) {

	q = &GetStatus{
		authToken: authToken,
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

// Execute executes the MedCo client get-status request
func (clientGetStatus *GetStatus) Execute() (result *medco_node.GetStatusOKBody, err error) {

	result, err = clientGetStatus.submitToNode()
	return

}

func (clientGetStatus *GetStatus) submitToNode() (result *medco_node.GetStatusOKBody, err error) {

	params := medco_node.NewGetStatusParamsWithTimeout(time.Duration(utilclient.GetStatusTimeoutSeconds) * time.Second)

	response, err := clientGetStatus.httpMedCoClient.MedcoNode.GetStatus(params, httptransport.BearerToken(clientGetStatus.authToken))

	if err != nil {
		logrus.Error("Get status error: ", err)
		return nil, err
	}

	return response.Payload, nil

}
