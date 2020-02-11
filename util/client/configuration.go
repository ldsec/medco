package utilclient

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

// QueryTimeoutSeconds is the timeout for the client query in seconds (default to 3 minutes)
var QueryTimeoutSeconds int64

// GenomicAnnotationsQueryTimeoutSeconds is the timeout for the client query in seconds (default to 10 seconds)
var GenomicAnnotationsQueryTimeoutSeconds int64

// MedCoConnectorURL is the URL of the MedCo connector this client is attached to
var MedCoConnectorURL string

// OidcReqTokenURL is the URL from which the JWT is retrieved
var OidcReqTokenURL string

// OidcReqTokenClientID is the client ID used to retrieve the JWT
var OidcReqTokenClientID string

func init() {
	var err error

	QueryTimeoutSeconds, err = strconv.ParseInt(os.Getenv("CLIENT_QUERY_TIMEOUT_SECONDS"), 10, 64)
	if err != nil || QueryTimeoutSeconds < 0 {
		logrus.Warn("invalid client query timeout")
		QueryTimeoutSeconds = 3 * 60
	}

	GenomicAnnotationsQueryTimeoutSeconds, err = strconv.ParseInt(os.Getenv("CLIENT_GENOMIC_ANNOTATIONS_QUERY_TIMEOUT_SECONDS"), 10, 64)
	if err != nil || GenomicAnnotationsQueryTimeoutSeconds < 0 {
		logrus.Warn("invalid client genomic annotations query timeout")
		GenomicAnnotationsQueryTimeoutSeconds = 10
	}

	MedCoConnectorURL = os.Getenv("MEDCO_CONNECTOR_URL")

	OidcReqTokenURL = os.Getenv("OIDC_REQ_TOKEN_URL")

	OidcReqTokenClientID = os.Getenv("OIDC_REQ_TOKEN_CLIENT_ID")
}
