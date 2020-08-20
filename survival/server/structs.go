package survivalserver

import (
	"database/sql"
	"errors"
	"os"
	"strconv"

	utilserver "github.com/ldsec/medco-connector/util/server"
	"github.com/sirupsen/logrus"

	"github.com/ldsec/medco-connector/wrappers/unlynx"
)

const (
	schema                string = "i2b2demodata_i2b2"
	table                 string = "observation_fact"
	blobColumn            string = "observation_blob"
	timeConceptColumn     string = "time_code"
	patientNumColumn      string = "patient_num"
	eventOfInterestColumn string = "event_of_interest"
	censoringEventColumn  string = "censoring_event"
	naiveVersion          bool   = false
	interBlob             string = "," //blobconcept does not contain any comma
	eventFlagSeparator    string = ` `
	timeConceptRootPath   string = `/SurvivalAnalysis/`
	errLogTrace           int    = 10
	timeTable             string = "time_dimension"
	pathColumn            string = "concept_path"
	tagQuery              string = `SELECT time_path,time_code FROM survivaldemodata_i2b2.` + timeTable
	errorChanMaxSize      int    = 1024
)

// EncryptedEncID is  ElGamal encryption
type EncryptedEncID string

// TagID is  deterministic tag
type TagID string
type PatientNum string

var (
	directAccessDB *sql.DB
)

func NewSqlAccess(envPrefix string) (access *sql.DB, err error) {
	host := os.Getenv(envPrefix + "HOST")
	port, err := strconv.ParseInt(os.Getenv(envPrefix+"PORT"), 10, 64)
	if err != nil || port < 0 || port > 65535 {
		logrus.Warn("Invalid port, defaulted")
		port = 5432
	}
	dbName := os.Getenv(envPrefix + "NAME")
	if dbName == "" {
		err = errors.New("dbname variable not found for envPrefix : " + envPrefix)
		return
	}

	loginUser := os.Getenv(envPrefix + "USER")
	if loginUser == "" {
		err = errors.New("user  variable not found for envPrefix : " + envPrefix)
		return
	}
	loginPw := os.Getenv(envPrefix + "PW")
	if loginPw == "" {
		err = errors.New("password variable not found for envPrefix : " + envPrefix)
		return
	}

	access, err = utilserver.InitializeConnectionToDB(host, int(port), dbName, loginUser, loginPw)
	if err != nil {
		logrus.Errorf("Unable to connect database for direct access to I2B2 for host : %s, for port %d", host, port)
		return
	}
	logrus.Info("Connected I2B2 DB for direct access")
	return
}

func init() {
	var err error
	directAccessDB, err = NewSqlAccess("DIRECT_ACCESS_DB_")
	if err != nil {
		logrus.Error("Unable to connect to database for direct access")
	}

	return
}

// Result holds information about time point, events, execution times and error
type Result struct {
	EventValueAgg     string //or kyber.point
	CensoringValueAgg string //or kyber.point

}

type TimePointResult struct {
	TimePoint string
	Result    Result
}

type EventGroup struct {
	GroupID          string
	EncInitialCount  string
	TimePointResults []*TimePointResult
}
type EventGroups []*EventGroup

func (eventGroups EventGroups) Len() int {
	return len(eventGroups)
}

func (eventGroups EventGroups) Less(i, j int) bool {
	return eventGroups[i].GroupID < eventGroups[j].GroupID
}

func (eventGroups EventGroups) Swap(i, j int) {
	eventGroups[i], eventGroups[j] = eventGroups[j], eventGroups[i]
}

func (eventGroup EventGroup) Len() int {
	return len(eventGroup.TimePointResults)
}

func (eventGroup EventGroup) Less(i, j int) bool {
	return eventGroup.TimePointResults[i].TimePoint < eventGroup.TimePointResults[j].TimePoint
}

func (eventGroup EventGroup) Swap(i, j int) {
	eventGroup.TimePointResults[i], eventGroup.TimePointResults[j] = eventGroup.TimePointResults[j], eventGroup.TimePointResults[i]
}

func zeroPointEncryption() (res string, err error) {

	res, err = unlynx.EncryptWithCothorityKey(int64(0))
	return
}

func unlynxRequestName(queryName string, timecode TagID) string {
	return queryName + string(timecode)
}
