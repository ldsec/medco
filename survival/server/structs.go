package survivalserver

import (
	"fmt"

	"strconv"
	"strings"
	"sync"

	"time"

	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/ldsec/medco-connector/wrappers/i2b2"
	"github.com/ldsec/medco-connector/wrappers/unlynx"
)

type Result struct {
	EventValue     string //or kyber.point
	CensoringValue string //or kyber.point
	Delay          time.Duration
	Error          error
}

const eventFlagSeparator = ` `

//type Map sync.Map

type PatientID string

func (str PatientID) String() string {
	return string(str)
}

var ResultMap sync.Map

func BreakBlob(blobValue string) (eventOfInterest, censoringEvent string, err error) {
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

const timeConceptRootPath = `/SurvivalAnalysis/`

//TODO another function already exists in unlynx wrapper
func ZeroPointEncryption() (res string, err error) {

	events, err := unlynx.EncryptWithCothorityKey(int64(0))
	if err != nil {
		return
	}
	censoringEvents, err := unlynx.EncryptWithCothorityKey(int64(0))
	if err != nil {
		return
	}
	res = events + eventFlagSeparator + censoringEvents
	return
}

func UnlynxRequestName(queryName, timecode string) string {
	return queryName + timecode
}

//this redundance is a quick solutiio to avoid cyclic imports with medcoservers
//TODO expected to become ghost code
type ExploreQuery struct {
	ID     string
	Query  *models.ExploreQuery
	Result struct {
		EncCount       string
		EncPatientList []string
		Timers         map[string]time.Duration
		EncEvents      map[string][2]string
	}
	ExecuteCallback func(*ExploreQuery, []string, []string, int) error
}

var ExecCallback func(*ExploreQuery, []string, []string, int) error

func (q *ExploreQuery) Execute(patientIDs []string) error {
	timeCodes, err := GetTimeCodes()
	if err != nil {
		return err
	}

	return q.ExecuteCallback(q, patientIDs, timeCodes, 1)
}

func GetTimeCodes() (timeCode []string, err error) {
	results, err := i2b2.GetOntologyChildren(timeConceptRootPath)
	if err != nil {
		return
	}
	for _, result := range results {
		timeCode = append(timeCode, strconv.FormatInt(*(result.MedcoEncryption.ID), 10))
	}
	return

}
