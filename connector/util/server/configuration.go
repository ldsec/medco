package utilserver

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	onetLog "go.dedis.ch/onet/v3/log"
	"os"
	"strconv"
	"strings"
	"time"
)

// LogLevel is the log level, assuming the same convention as the cothority / unlynx log levels:
// TRACE(5), DEBUG(4), INFO(3), WARNING(2), ERROR(1), FATAL(0)
var LogLevel int

// MedCoNodesURL is the slice of the URL of all the MedCo nodes, with the order matching the position in the slice
var MedCoNodesURL []string

// MedCoNodeIdx is the index of this node in the list of nodes
var MedCoNodeIdx int

// UnlynxGroupFilePath is the path of the unlynx group file from which is derived the cothority public key
var UnlynxGroupFilePath string

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

// I2b2WaitTimeSeconds is the i2b2 timeout (seconds)
var I2b2WaitTimeSeconds int

// OidcProviders are the OIDC providers this node trusts
var OidcProviders []*oidcProvider

// JwksTimeout is the timeout for JWKS requests made to OIDC providers
var JwksTimeout time.Duration

// MedCoObfuscationMin is the minimum variance passed to the random distribution for the obfuscation
var MedCoObfuscationMin int

// DBHost is the host of the DB
var DBHost string

// DBPort is the number of the port used by the DB
var DBPort int

// DBName is the name of the database
var DBName string

// DBLoginUser is the database login user
var DBLoginUser string

// DBLoginPassword is the database login password
var DBLoginPassword string

// DBConnection is the connection to the database
var DBConnection *sql.DB

// I2B2DBHost is the host of the I2B2 DB
var I2B2DBHost string

// I2B2DBPort is the number of the port used by the I2B2 DB
var I2B2DBPort int

// I2B2DBName is the name of the I2B2 database
var I2B2DBName string

// I2B2DBLoginUser is the I2B2 database login user
var I2B2DBLoginUser string

// I2B2DBLoginPassword is the I2B2 database login password
var I2B2DBLoginPassword string

// I2B2DBConnection is the connection to the i2b2 database
var I2B2DBConnection *sql.DB

func init() {
	SetLogLevel(os.Getenv("LOG_LEVEL"))

	MedCoNodesURL = strings.Split(os.Getenv("MEDCO_NODES_URL"), ",")

	idx, err := strconv.ParseInt(os.Getenv("MEDCO_NODE_IDX"), 10, 64)
	if err != nil || idx < 0 {
		logrus.Warn("invalid MedCoNodeIdx")
		idx = 0
	}
	MedCoNodeIdx = int(idx)

	UnlynxGroupFilePath = os.Getenv("UNLYNX_GROUP_FILE_PATH")

	unlynxto, err := strconv.ParseInt(os.Getenv("UNLYNX_TIMEOUT_SECONDS"), 10, 64)
	if err != nil || unlynxto < 0 {
		logrus.Warn("invalid UnlynxTimeoutSeconds, defaulted")
		unlynxto = 3 * 60 // 3 minutes
	}
	UnlynxTimeoutSeconds = int(unlynxto)

	I2b2HiveURL = os.Getenv("I2B2_HIVE_URL")
	I2b2LoginDomain = os.Getenv("I2B2_LOGIN_DOMAIN")
	I2b2LoginProject = os.Getenv("I2B2_LOGIN_PROJECT")
	I2b2LoginUser = os.Getenv("I2B2_LOGIN_USER")
	I2b2LoginPassword = os.Getenv("I2B2_LOGIN_PASSWORD")

	i2b2to, err := strconv.ParseInt(os.Getenv("I2B2_WAIT_TIME_SECONDS"), 10, 64)
	if err != nil || i2b2to < 0 {
		logrus.Warn("invalid I2b2WaitTimeSeconds, defaulted")
		i2b2to = 3 * 60 // 3 minutes
	}
	I2b2WaitTimeSeconds = int(i2b2to)

	parseOidcProviders()

	jwksto, err := strconv.ParseInt(os.Getenv("OIDC_JWKS_TIMEOUT_SECONDS"), 10, 64)
	if err != nil || jwksto < 0 {
		logrus.Warn("invalid JwksTimeoutSeconds, defaulted")
		i2b2to = 30 // 30 seconds
	}
	JwksTimeout = time.Duration(jwksto) * time.Second

	obf, err := strconv.ParseInt(os.Getenv("MEDCO_OBFUSCATION_MIN"), 10, 64)
	if err != nil || obf < 0 {
		logrus.Warn("invalid MedCoObfuscationMin, defaulted")
		obf = 5
	}
	MedCoObfuscationMin = int(obf)

	DBHost = os.Getenv("MC_DB_HOST")
	DBName = os.Getenv("MC_DB_NAME")
	DBLoginUser = os.Getenv("MC_DB_USER")
	DBLoginPassword = os.Getenv("MC_DB_PASSWORD")

	dbPort, err := strconv.ParseInt(os.Getenv("MC_DB_PORT"), 10, 64)
	if err != nil || dbPort < 0 || dbPort > 65535 {
		logrus.Warn("invalid DB port, defaulted")
		dbPort = 5432
	}
	DBPort = int(dbPort)

	DBConnection, err = InitializeConnectionToDB(DBHost, DBPort, DBName, DBLoginUser, DBLoginPassword)
	if err != nil {
		logrus.Error("Impossible to initialize connection to DB")
		return
	}

	I2B2DBHost = os.Getenv("I2B2_DB_HOST")
	I2B2DBName = os.Getenv("I2B2_DB_NAME")
	I2B2DBLoginUser = os.Getenv("I2B2_DB_USER")
	I2B2DBLoginPassword = os.Getenv("I2B2_DB_PW")

	I2B2dbPort, err := strconv.ParseInt(os.Getenv("I2B2_DB_PORT"), 10, 64)
	if err != nil || I2B2dbPort < 0 || I2B2dbPort > 65535 {
		logrus.Warn("invalid DB port, defaulted")
		I2B2dbPort = 5432
	}
	I2B2DBPort = int(I2B2dbPort)

	I2B2DBConnection, err = InitializeConnectionToDB(I2B2DBHost, I2B2DBPort, I2B2DBName, I2B2DBLoginUser, I2B2DBLoginPassword)
	if err != nil {
		logrus.Error("Impossible to initialize connection to I2B2 DB")
		return
	}
}

func parseOidcProviders() {

	jwksUrls := strings.Split(os.Getenv("OIDC_JWKS_URLS"), ",")
	jwtIssuers := strings.Split(os.Getenv("OIDC_JWT_ISSUERS"), ",")
	clientIds := strings.Split(os.Getenv("OIDC_CLIENT_IDS"), ",")
	jwtUserIDClaims := strings.Split(os.Getenv("OIDC_JWT_USER_ID_CLAIMS"), ",")

	if len(jwksUrls) != len(jwtIssuers) || len(jwksUrls) != len(clientIds) || len(jwksUrls) != len(jwtUserIDClaims) {
		logrus.Fatal("inconsistent OIDC configuration")
	}

	for i := range jwksUrls {
		OidcProviders = append(OidcProviders, &oidcProvider{
			JwksURL:           jwksUrls[i],
			JwtIssuer:         jwtIssuers[i],
			ClientID:          clientIds[i],
			JwtUserIDClaim:    jwtUserIDClaims[i],
			JwksTTL:           time.Hour,
			JwtAcceptableSkew: 30 * time.Second,
		})
	}
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
