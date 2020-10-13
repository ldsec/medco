package utilclient

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
)

// SearchTimeoutSeconds is the timeout for the client query in seconds (default to 3 minutes)
var SearchTimeoutSeconds int64

// QueryTimeoutSeconds is the timeout for the client query in seconds (default to 3 minutes)
var QueryTimeoutSeconds int64

// GenomicAnnotationsQueryTimeoutSeconds is the timeout for the client query in seconds (default to 10 seconds)
var GenomicAnnotationsQueryTimeoutSeconds int64

// GetNodeStatusTimeoutSeconds is the timeout for the GetNodeStatus in seconds (default to 60 seconds)
var GetNodeStatusTimeoutSeconds int64

// GetNetworkTimeoutSeconds is the timeout for the GetNetwork request in seconds (default to 60 seconds)
var GetNetworkTimeoutSeconds int64

// MedCoConnectorURL is the URL of the MedCo connector this client is attached to
var MedCoConnectorURL string

// OidcReqTokenURL is the URL from which the JWT is retrieved
var OidcReqTokenURL string

// OidcReqTokenClientID is the client ID used to retrieve the JWT
var OidcReqTokenClientID string

func init() {
	var err error

	SearchTimeoutSeconds, err = strconv.ParseInt(os.Getenv("CLIENT_SEARCH_TIMEOUT_SECONDS"), 10, 64)
	if err != nil || QueryTimeoutSeconds < 0 {
		logrus.Warn("invalid client search timeout")
		QueryTimeoutSeconds = 10
	}

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

	GetNodeStatusTimeoutSeconds, err = strconv.ParseInt(os.Getenv("CLIENT_GET_NODE_STATUS_TIMEOUT_SECONDS"), 10, 64)
	if err != nil || GetNodeStatusTimeoutSeconds < 0 {
		logrus.Warn("invalid client get node status timeout")
		GetNodeStatusTimeoutSeconds = 60
	}

	GetNetworkTimeoutSeconds, err = strconv.ParseInt(os.Getenv("CLIENT_GET_NETWORK_TIMEOUT_SECONDS"), 10, 64)
	if err != nil || GetNetworkTimeoutSeconds < 0 {
		logrus.Warn("invalid client get network timeout")
		GetNetworkTimeoutSeconds = 60
	}

	MedCoConnectorURL = os.Getenv("MEDCO_CONNECTOR_URL")

	OidcReqTokenURL = os.Getenv("OIDC_REQ_TOKEN_URL")

	OidcReqTokenClientID = os.Getenv("OIDC_REQ_TOKEN_CLIENT_ID")
}
