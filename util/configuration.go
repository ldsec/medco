package util

import (
	"github.com/sirupsen/logrus"
	onetLog "go.dedis.ch/onet/v3/log"
	"os"
	"strconv"
)

// LogLevel returns the log level, assuming the same convention as the cothority / unlynx log levels:
// TRACE(5), DEBUG(4), INFO(3), WARNING(2), ERROR(1), FATAL(0)
var LogLevel int

var UnlynxGroupFilePath string
var UnlynxGroupFileIdx int
var UnlynxTimeoutSeconds int

// URL of the i2b2 hive this connector is using
var I2b2HiveURL string
// i2b2 login domain
var I2b2LoginDomain string
// i2b2 login project
var I2b2LoginProject string
// i2b2 login user
var I2b2LoginUser string
// i2b2 login password
var I2b2LoginPassword string
// i2b2 timeout (seconds)
var I2b2TimeoutSeconds int

// token (shared secret) used for internal PICSURE 2 authorization
var picsure2InternalToken string

func init() {
	SetLogLevel(os.Getenv("LOG_LEVEL"))

	UnlynxGroupFilePath = os.Getenv("UNLYNX_GROUP_FILE_PATH")
	UnlynxTimeoutSeconds = 180

	UnlynxGroupFileIdx, err := strconv.ParseInt(os.Getenv("UNLYNX_GROUP_FILE_IDX"), 10, 64)
	if err != nil || UnlynxGroupFileIdx < 0 {
		logrus.Warn("invalid UnlynxGroupFileIdx")
	}

	I2b2HiveURL = os.Getenv("I2B2_HIVE_URL")
	I2b2LoginDomain = os.Getenv("I2B2_LOGIN_DOMAIN")
	I2b2LoginProject = os.Getenv("I2B2_LOGIN_PROJECT")
	I2b2LoginUser = os.Getenv("I2B2_LOGIN_USER")
	I2b2LoginPassword = os.Getenv("I2B2_LOGIN_PASSWORD")
	I2b2TimeoutSeconds = 180

	picsure2InternalToken = os.Getenv("PICSURE2_INTERNAL_TOKEN")
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
