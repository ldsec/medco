package utilclient

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ldsec/medco-connector/restapi/client"
	"github.com/ldsec/medco-connector/restapi/client/medco_network"

	"github.com/sirupsen/logrus"
)

// MetaData creates a request get metadata and returns the response
func MetaData(token string, disableTLSCheck bool) (*medco_network.GetMetadataOK, error) {
	parsedURL, err := url.Parse(MedCoConnectorURL)

	if err != nil {
		err = fmt.Errorf("cannot parse MedCo connector URL: %s", err.Error())
		logrus.Error(err)
		return nil, err
	}

	transport := httptransport.New(parsedURL.Host, parsedURL.Path, []string{parsedURL.Scheme})
	transport.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: disableTLSCheck}

	getMetadataResp, err := client.New(transport, nil).MedcoNetwork.GetMetadata(
		medco_network.NewGetMetadataParamsWithTimeout(30*time.Second),
		httptransport.BearerToken(token),
	)
	if err != nil {
		err = fmt.Errorf("get network metadata request error: %s", err.Error())
		logrus.Error(err)
	}
	return getMetadataResp, err

}
