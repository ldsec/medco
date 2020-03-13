package survivalserver

import (
	"database/sql"
	"fmt"
	"os"
	"strconv"

	utilserver "github.com/ldsec/medco-connector/util/server"
	"github.com/sirupsen/logrus"

	"strings"

	"time"

	"github.com/ldsec/medco-connector/wrappers/unlynx"
)

const (
	schema              string    = "i2b2demodata_i2b2"
	table               string    = "observation_fact"
	blobColumn          string    = "observation_blob"
	timeConceptColumn   string    = "concept_cd"
	patientIDColumn     string    = "patient_num"
	naiveVersion        bool      = false
	interBlob           string    = "," //blobconcept does not contain any comma
	eventFlagSeparator  string    = ` `
	timeConceptRootPath string    = `/SurvivalAnalysis/`
	errLogTrace         int       = 10
	available           lockState = 0
	locked              lockState = 1
	conceptTable        string    = "concept_dimension"
	pathColumn          string    = "concept_path"
	tagQuery            string    = `SELECT concept_path,concept_cd FROM i2b2demodata_i2b2.` + conceptTable
	errorChanMaxSize    int       = 1024
)

// EncryptedEncID is  ElGamal encryption
type EncryptedEncID string

// TagID is  deterministic tag
type TagID string

var directAccessDB *sql.DB
var dbName string

func init() {

	host := os.Getenv("DIRECT_ACCESS_DB_HOST")
	port, err := strconv.ParseInt(os.Getenv("DIRECT_ACCESS_DB_PORT"), 10, 64)
	if err != nil || port < 0 || port > 65535 {
		logrus.Warn("Invalid port, defaulted")
		port = 5432
	}
	name := os.Getenv("DIRECT_ACCESS_DB_NAME")
	dbName = name
	loginUser := os.Getenv("DIRECT_ACCESS_DB_USER")
	loginPw := os.Getenv("DIRECT_ACCESS_DB_PW")

	directAccessDB, err = utilserver.InitializeConnectionToDB(host, int(port), name, loginUser, loginPw)
	if err != nil {
		logrus.Error("Unable to connect database for direct access to I2B2")
		return
	}
	logrus.Info("Connected I2B2 DB for direct access")
	return
}

// Result holds information about time point, events, execution times and error
type Result struct {
	EventValue     string //or kyber.point
	CensoringValue string //or kyber.point
	Delay          time.Duration
	Error          error
}

func breakBlob(blobValue string) (eventOfInterest, censoringEvent string, err error) {
	res := strings.Split(blobValue, eventFlagSeparator)
	//TODO magic
	if len(res) != 2 {
		err = fmt.Errorf(`Blob value %s is ill-formed, it should be  to base64 encoded sequence separated by "%s" (without quotes)`, blobValue, eventFlagSeparator)
		return
	}
	eventOfInterest = res[0]
	censoringEvent = res[1]
	return
}

func zeroPointEncryption() (res string, err error) {

	res, err = unlynx.EncryptWithCothorityKey(int64(0))
	return
}

func unlynxRequestName(queryName string, timecode TagID) string {
	return queryName + string(timecode)
}
