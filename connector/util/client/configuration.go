package utilclient

import (
	"os"
	"strconv"

	"github.com/sirupsen/logrus"
)

// SearchTimeoutSeconds is the timeout for the client search in seconds (default to 3 minutes)
var SearchTimeoutSeconds int64

// QueryTimeoutSeconds is the timeout for the client query in seconds (default to 3 minutes)
var QueryTimeoutSeconds int64

// GenomicAnnotationsQueryTimeoutSeconds is the timeout for the client genomic annotations in seconds (default to 10 seconds)
var GenomicAnnotationsQueryTimeoutSeconds int64

// QueryToolsTimeoutSeconds is the timeout for the client query tools in seconds (default to 10 seconds)
var QueryToolsTimeoutSeconds int64

// SurvivalAnalysisTimeoutSeconds is the timeout for the client survival analysis in seconds (default to 5 minutes)
var SurvivalAnalysisTimeoutSeconds int64

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
		SearchTimeoutSeconds = 10
	}

	QueryTimeoutSeconds, err = strconv.ParseInt(os.Getenv("CLIENT_QUERY_TIMEOUT_SECONDS"), 10, 64)
	if err != nil || QueryTimeoutSeconds < 0 {
		logrus.Warn("invalid client query timeout")
		QueryTimeoutSeconds = 3 * 60
	}

	QueryToolsTimeoutSeconds, err = strconv.ParseInt(os.Getenv("QUERY_TOOLS_TIMEOUT_SECONDS"), 10, 64)
	if err != nil || QueryToolsTimeoutSeconds < 0 {
		logrus.Warn("invalid client query tools timeout")
		QueryToolsTimeoutSeconds = 10
	}

	SurvivalAnalysisTimeoutSeconds, err = strconv.ParseInt(os.Getenv("SURVIVAL_ANALYSIS_TIMEOUT_SECONDS"), 10, 64)
	if err != nil || SurvivalAnalysisTimeoutSeconds < 0 {
		logrus.Warn("invalid client survival analysis timeout")
		SurvivalAnalysisTimeoutSeconds = 5 * 60
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
