package medcoclient

import (
	"crypto/tls"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ldsec/medco-connector/restapi/client"
	"github.com/ldsec/medco-connector/restapi/client/genomic_annotations"
	utilclient "github.com/ldsec/medco-connector/util/client"
	"github.com/sirupsen/logrus"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// GenomicAnnotationsGetValues is a MedCo client genomic-annotations get-values request
type GenomicAnnotationsGetValues struct {

	// httpMedCoClient is the HTTP client for the MedCo connector
	httpMedCoClient *client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

	annotation string

	value string

	limit *int64
}

// GenomicAnnotationsGetVariants is a MedCo client genomic-annotations get-variants request
type GenomicAnnotationsGetVariants struct {

	// httpMedCoClient is the HTTP client for the MedCo connector
	httpMedCoClient *client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

	annotation string

	value string

	zygosity []string

	encrypted *bool
}

// NewGenomicAnnotationsGetValues creates a new MedCo client genomic-annotations get-values request
func NewGenomicAnnotationsGetValues(authToken, annotation, value string, limit *int64, disableTLSCheck bool) (q *GenomicAnnotationsGetValues, err error) {

	q = &GenomicAnnotationsGetValues{
		authToken:  authToken,
		annotation: annotation,
		value:      value,
		limit:      limit,
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

// NewGenomicAnnotationsGetVariants creates a new MedCo client genomic-annotations get-variants request
func NewGenomicAnnotationsGetVariants(authToken, annotation, value string, zygosity string, encrypted *bool, disableTLSCheck bool) (q *GenomicAnnotationsGetVariants, err error) {

	q = &GenomicAnnotationsGetVariants{
		authToken:  authToken,
		annotation: annotation,
		value:      value,
		zygosity:   strings.Split(zygosity, "|"),
		encrypted:  encrypted,
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

// Execute executes the MedCo client get-values query
func (clientGenomicAnnotationsGetValues *GenomicAnnotationsGetValues) Execute() (result []string, err error) {

	result, err = clientGenomicAnnotationsGetValues.submitToNode()
	return

}

// Execute executes the MedCo client get-variants query
func (clientGenomicAnnotationsGetVariants *GenomicAnnotationsGetVariants) Execute() (result []string, err error) {

	result, err = clientGenomicAnnotationsGetVariants.submitToNode()
	return

}

func (clientGenomicAnnotationsGetValues *GenomicAnnotationsGetValues) submitToNode() (result []string, err error) {

	params := genomic_annotations.NewGetValuesParamsWithTimeout(time.Duration(utilclient.GenomicAnnotationsQueryTimeoutSeconds) * time.Second)
	params.Annotation = clientGenomicAnnotationsGetValues.annotation
	params.Value = clientGenomicAnnotationsGetValues.value
	if *clientGenomicAnnotationsGetValues.limit != 0 {
		params.Limit = clientGenomicAnnotationsGetValues.limit
	}

	response, err := clientGenomicAnnotationsGetValues.httpMedCoClient.GenomicAnnotations.GetValues(params, httptransport.BearerToken(clientGenomicAnnotationsGetValues.authToken))

	if err != nil {
		logrus.Error("Genomic annotations get values error: ", err)
		return nil, err
	}

	return response.Payload, nil

}

func (clientGenomicAnnotationsGetVariants *GenomicAnnotationsGetVariants) submitToNode() (result []string, err error) {

	params := genomic_annotations.NewGetVariantsParamsWithTimeout(time.Duration(utilclient.GenomicAnnotationsQueryTimeoutSeconds) * time.Second)
	params.Annotation = clientGenomicAnnotationsGetVariants.annotation
	params.Value = clientGenomicAnnotationsGetVariants.value
	params.Zygosity = clientGenomicAnnotationsGetVariants.zygosity
	params.Encrypted = clientGenomicAnnotationsGetVariants.encrypted

	response, err := clientGenomicAnnotationsGetVariants.httpMedCoClient.GenomicAnnotations.GetVariants(params, httptransport.BearerToken(clientGenomicAnnotationsGetVariants.authToken))

	if err != nil {
		logrus.Error("Genomic annotations get variants error: ", err)
		return nil, err
	}

	return response.Payload, nil

}
