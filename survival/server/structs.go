package survivalserver

import (
	"fmt"

	"strings"

	"time"

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

const TimeConceptRootPath = `/SurvivalAnalysis/`

//TODO another function already exists in unlynx wrapper
func ZeroPointEncryption() (res string, err error) {

	res, err = unlynx.EncryptWithCothorityKey(int64(0))
	return
}

func UnlynxRequestName(queryName, timecode string) string {
	return queryName + timecode
}
