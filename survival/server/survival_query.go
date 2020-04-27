package survivalserver

import (
	"errors"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

// Query holds the ID of the survival analysis, its parameters and a pointer to its results
type Query struct {
	ID            string
	UserPublicKey string
	PatientSetID  string
	TimeCodes     []EncryptedEncID

	//TODO also hide that
	Result struct {
		Timers    map[string]time.Duration
		EncEvents map[string]map[string][2]string
	}

	spin *Spin
}

// NewQuery query constructor
func NewQuery(qID, pubKey, patientSetID string, timeCodes []string) *Query {
	encryptedEncIDs := make([]EncryptedEncID, len(timeCodes))
	for idx, timeCode := range timeCodes {
		encryptedEncIDs[idx] = EncryptedEncID(timeCode)
	}
	res := &Query{ID: qID, UserPublicKey: pubKey,
		PatientSetID: patientSetID,
		TimeCodes:    encryptedEncIDs,
		spin:         NewSpin()}
	res.Result.EncEvents = make(map[string]map[string][2]string)
	res.Result.Timers = make(map[string]time.Duration)
	return res
}

func (q *Query) addTimer(label string, since time.Time) (err error) {
	if _, exists := q.Result.Timers[label]; exists {
		err = errors.New("label " + label + " already exists in timers")
		return
	}
	q.Result.Timers[label] = time.Since(since)
	return
}
func (q *Query) addTimers(timers map[string]time.Duration) (err error) {
	for label, duration := range timers {
		if _, exists := q.Result.Timers[label]; exists {
			err = errors.New("label " + label + " already exists in timers")
			return
		}
		q.Result.Timers[label] = duration
	}
	return
}

// PrintTimers print the execution time information if the log level is Debug
func (q *Query) PrintTimers() {
	for label, duration := range q.Result.Timers {
		logrus.Debug(label + duration.String())
	}
}

//SetResultMap sets the result map structure to return
func (q *Query) SetResultMap(resultMap map[string]map[string][2]string) error {
	q.spin.Lock()
	defer q.spin.Unlock()

	q.Result.EncEvents = resultMap
	return nil

}

// Execute is the function that translates a medco survival query in multiple calls on psql and unlynx
func (q *Query) Execute(batchNumber int) (err error) {
	err = errors.New("TODO")
	return
}

type unlynxResult struct {
	Key   *EncryptedEncID
	Value string
	Error error
}

func buildSingleQuery(patients []string, timeCode string) string {

	patientAcc := stringMapAndAdd(patients)
	res := `SELECT ` + blobColumn + ` FROM ` + schema + `.` + table + ` `
	res += `WHERE ` + timeConceptColumn + ` = '` + timeCode + `' AND ` + patientIDColumn + ` IN (` + patientAcc + `)`
	return res

}
func buildBatchQuery(patients []string, timeCodes []TagID) string {

	patientAcc := stringMapAndAdd(patients)
	timeCodesString := make([]string, len(timeCodes))
	for idx, timeCode := range timeCodes {
		timeCodesString[idx] = string(timeCode)
	}
	timeAcc := stringMapAndAdd(timeCodesString)
	res := `SELECT ` + timeConceptColumn + `, STRING_AGG(` + blobColumn + `,'` + interBlob + `') FROM ` + schema + `.` + table + ` `
	res += `WHERE ` + timeConceptColumn + ` IN (` + timeAcc + `) AND ` + patientIDColumn + ` IN (` + patientAcc + `)`
	res += `GROUP BY (` + timeConceptColumn + `)`
	return res

}

func stringMapAndAdd(inputList []string) string {
	outputList := make([]string, len(inputList))
	for i, str := range inputList {
		outputList[i] = `'` + str + `'`
	}
	return strings.Join(outputList, `,`)

}
