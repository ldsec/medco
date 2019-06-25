package util

import (
	"github.com/sirupsen/logrus"
	onetLog "go.dedis.ch/onet/v3/log"
	"os"
	"strconv"
)

// LogLevel is the log level, assuming the same convention as the cothority / unlynx log levels:
// TRACE(5), DEBUG(4), INFO(3), WARNING(2), ERROR(1), FATAL(0)
var LogLevel int

// UnlynxGroupFilePath is the path of the unlynx group file from which is derived the cothority public key
var UnlynxGroupFilePath string
// UnlynxGroupFileIdx is the index (in the group file) of this node
var UnlynxGroupFileIdx int
// UnlynxTimeoutSeconds is the unlynx communication timeout (seconds)
var UnlynxTimeoutSeconds int

// I2b2HiveURL is the URL of the i2b2 hive this connector is using
var I2b2HiveURL string
// I2b2LoginDomain is the i2b2 login domain
var I2b2LoginDomain string
// I2b2LoginProject is the i2b2 login project
var I2b2LoginProject string
// I2b2LoginUser is the i2b2 login user
var I2b2LoginUser string
// I2b2LoginPassword is the i2b2 login password
var I2b2LoginPassword string
// I2b2TimeoutSeconds is the i2b2 timeout (seconds)
var I2b2TimeoutSeconds int

// JwksURL is the URL from which the JWT signing keys are retrieved
var JwksURL string
// JwksTTLSeconds is the TTL of JWKS requests
var JwksTTLSeconds int64
// OidcJwtIssuer is the token issuer (for JWT validation)
var OidcJwtIssuer string
// OidcClientID is the OIDC client ID
var OidcClientID string
// OidcJwtUserIDClaim is the JWT claim containing the user ID
var OidcJwtUserIDClaim string

// MedCoObfuscationMin is the minimum variance passed to the random distribution for the obfuscation
var MedCoObfuscationMin int

func init() {
	SetLogLevel(os.Getenv("LOG_LEVEL"))

	UnlynxGroupFilePath = os.Getenv("UNLYNX_GROUP_FILE_PATH")
	UnlynxTimeoutSeconds = 3 * 60 // 3 minutes

	idx, err := strconv.ParseInt(os.Getenv("UNLYNX_GROUP_FILE_IDX"), 10, 64)
	if err != nil || idx < 0 {
		logrus.Warn("invalid UnlynxGroupFileIdx")
		idx = 0
	}
	UnlynxGroupFileIdx = int(idx)

	I2b2HiveURL = os.Getenv("I2B2_HIVE_URL")
	I2b2LoginDomain = os.Getenv("I2B2_LOGIN_DOMAIN")
	I2b2LoginProject = os.Getenv("I2B2_LOGIN_PROJECT")
	I2b2LoginUser = os.Getenv("I2B2_LOGIN_USER")
	I2b2LoginPassword = os.Getenv("I2B2_LOGIN_PASSWORD")
	I2b2TimeoutSeconds = 3 * 60 // 3 minutes

	JwksURL = os.Getenv("JWKS_URL")
	JwksTTLSeconds = 60 * 60 // 1 hour
	OidcJwtIssuer = os.Getenv("OIDC_JWT_ISSUER")
	OidcClientID = os.Getenv("OIDC_CLIENT_ID")
	OidcJwtUserIDClaim = os.Getenv("OIDC_JWT_USER_ID_CLAIM")

	obf, err := strconv.ParseInt(os.Getenv("MEDCO_OBFUSCATION_MIN"), 10, 64)
	if err != nil || obf < 0 {
		logrus.Warn("invalid MedCoObfuscationMin, defaulted")
		obf = 5
	}
	MedCoObfuscationMin = int(obf)

}

// SetLogLevel initializes the log levels of all loggers
func SetLogLevel(lvl string) {
	// formatting
	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})

	intLvl, err := strconv.ParseInt(lvl, 10, 64)
	if err != nil || intLvl < 0 || intLvl > 5 {
		logrus.Warn("invalid LogLevel, defaulted")
		intLvl = 3
	}
	LogLevel = int(intLvl)
	logrus.SetLevel(logrus.Level(LogLevel + 1))
	onetLog.SetDebugVisible(LogLevel)
}
