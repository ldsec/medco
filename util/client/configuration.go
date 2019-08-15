package utilclient

import (
	"github.com/sirupsen/logrus"
	"os"
	"strconv"
	"strings"
)

// QueryTimeoutSeconds is the timeout for the client query in seconds (default to 3 minutes)
var QueryTimeoutSeconds int64
// Picsure2APIHost is the PIC-SURE host that broadcasts the query to the MedCo connectors
var Picsure2APIHost string
// Picsure2APIBasePath is the PIC-SURE hosts that broadcasts the query to the MedCo connectors
var Picsure2APIBasePath string
// Picsure2APIScheme is the PIC-SURE hosts that broadcasts the query to the MedCo connectors
var Picsure2APIScheme string
// Picsure2Resources are the resources to be queried by the client, corresponding to the medco connectors
var Picsure2Resources []string

// OidcReqTokenURL is the URL from which the JWT is retrieved
var OidcReqTokenURL string

func init() {
	var err error

	QueryTimeoutSeconds, err = strconv.ParseInt(os.Getenv("CLIENT_QUERY_TIMEOUT_SECONDS"), 10, 64)
	if err != nil || QueryTimeoutSeconds < 0 {
		logrus.Warn("invalid client query timeout")
		QueryTimeoutSeconds = 3 * 60
	}

	Picsure2APIHost = os.Getenv("PICSURE2_API_HOST")
	Picsure2APIBasePath = os.Getenv("PICSURE2_API_BASE_PATH")
	Picsure2APIScheme = os.Getenv("PICSURE2_API_SCHEME")
	Picsure2Resources = strings.Split(os.Getenv("PICSURE2_RESOURCES"), ",")

	OidcReqTokenURL = os.Getenv("OIDC_REQ_TOKEN_URL")
}
