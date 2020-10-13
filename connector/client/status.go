package medcoclient

import (
	"crypto/tls"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ldsec/medco/connector/restapi/client"
	"github.com/ldsec/medco/connector/restapi/client/medco_node"
	utilclient "github.com/ldsec/medco/connector/util/client"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

// GetNodeStatus is a GET request to the medco/node/status endpoint
type GetNodeStatus struct {

	// httpMedCoClient is the HTTP client for the MedCo connector
	httpMedCoClient *client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string
}

// NewGetNodeStatus creates a new GetNodeStatus request
func NewGetNodeStatus(authToken, connectorURL string, disableTLSCheck bool) (q *GetNodeStatus, err error) {

	q = &GetNodeStatus{
		authToken: authToken,
	}

	var parsedURL *url.URL

	if connectorURL == "" {
		parsedURL, err = url.Parse(utilclient.MedCoConnectorURL)
	} else {
		parsedURL, err = url.Parse(connectorURL)
	}
	if err != nil {
		logrus.Error("cannot parse MedCo connector URL: ", err)
		return
	}

	transport := httptransport.New(parsedURL.Host, parsedURL.Path, []string{parsedURL.Scheme})
	transport.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: disableTLSCheck}

	q.httpMedCoClient = client.New(transport, nil)

	return

}

// Execute executes the GetNodeStatus request
func (clientGetStatus *GetNodeStatus) Execute() (result *medco_node.GetNodeStatusOKBody, err error) {

	result, err = clientGetStatus.submitToNode()
	return

}

func (clientGetStatus *GetNodeStatus) submitToNode() (result *medco_node.GetNodeStatusOKBody, err error) {

	params := medco_node.NewGetNodeStatusParamsWithTimeout(time.Duration(utilclient.GetNodeStatusTimeoutSeconds) * time.Second)

	response, err := clientGetStatus.httpMedCoClient.MedcoNode.GetNodeStatus(params, httptransport.BearerToken(clientGetStatus.authToken))

	if err != nil {
		logrus.Error("Get node status error: ", err)
		return nil, err
	}

	return response.Payload, nil

}
