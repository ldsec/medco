package medcoclient

import (
	"crypto/tls"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ldsec/medco/connector/restapi/client"
	"github.com/ldsec/medco/connector/restapi/client/medco_network"
	utilclient "github.com/ldsec/medco/connector/util/client"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"time"
)

// GetNetwork is a GET request to the /medco/network endpoint
type GetNetwork struct {

	// httpMedCoClient is the HTTP client for the MedCo connector
	httpMedCoClient *client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string
}

// NewGetNetwork creates a new GetNetwork request
func NewGetNetwork(authToken string, disableTLSCheck bool) (q *GetNetwork, err error) {

	q = &GetNetwork{
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

// Execute executes the GetNetwork request
func (clientGetNetwork *GetNetwork) Execute() (result *medco_network.GetMetadataOKBody, err error) {

	params := medco_network.NewGetMetadataParamsWithTimeout(time.Duration(utilclient.GetNetworkTimeoutSeconds) * time.Second)

	response, err := clientGetNetwork.httpMedCoClient.MedcoNetwork.GetMetadata(params, httptransport.BearerToken(clientGetNetwork.authToken))

	if err != nil {
		logrus.Error("Get network error: ", err)
		return nil, err
	}

	return response.Payload, nil

}
