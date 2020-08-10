package survivalclient

import (
	"errors"
	"fmt"
	"sync"
	"time"

	medcoclient "github.com/ldsec/medco-connector/client"

	"github.com/ldsec/medco-connector/wrappers/unlynx"
	"github.com/sirupsen/logrus"
)

type timeCodesMaps [2]map[string]string

//ClientResultElement holds the information for the CLI whole susrvival analysis loop
type ClientResultElement struct {
	ClearTimePoint     string
	EncEventOfInterest string
	EncCensoringEvent  string
}

//ClientSurvival represents the whole survival analysis loop: it gets the time codes and the patient set, the request for the aggregates for the survival analysis and deciphers them
func ClientSurvival(token string, patientSetID int, startConcept, startColumn, endConcept, endColumn, username, password string, disableTLSCheck bool) (err error) {

	accessToken, err := GetToken(token, username, password, disableTLSCheck)
	if err != nil {
		return
	}

	var encTimeCodesMap map[string]string
	var encTimeCodesInverseMap map[string]string

	encTimeCodesMapChan := make(chan timeCodesMaps)

	var barrier sync.WaitGroup
	barrier.Add(2)

	encTimeCodesMaps := <-encTimeCodesMapChan

	encTimeCodesMap = encTimeCodesMaps[0]
	encTimeCodesInverseMap = encTimeCodesMaps[1]

	logrus.Info("Creating survival analysis request")

	var encTimeCodes []string
	for _, encTimeCode := range encTimeCodesMap {
		encTimeCodes = append(encTimeCodes, encTimeCode)
	}
	logrus.Debug("Mapping between duration in clear text and encryption of the integer identifier")
	logrus.Debug(fmt.Sprint(encTimeCodesMap))

	survivalAnalysis, err := NewSurvivalAnalysis(accessToken, patientSetID, nil, startConcept, startColumn, endConcept, endColumn, disableTLSCheck)
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	logrus.Info("Survival analysis request created. Executing")
	survResults, err := survivalAnalysis.Execute()
	if err != nil {
		logrus.Error(err.Error())
		return
	}
	survivalAnalysis.PrintTimers()

	logrus.Infof("obtained %d results", len(survResults))
	sequentialDecryptionTimer := time.Now()
	for groupID, encryptedResults := range survResults {
		for _, cipherEvents := range encryptedResults {
			var clearEventOfInterest int64
			var clearCensoringEvent int64
			clearTimePoint, ok := encTimeCodesInverseMap[cipherEvents.TimePoint]
			if !ok {
				err = errors.New("unexpected encrypted time point")
				return
			}
			clearEventOfInterest, err = survivalAnalysis.Decrypt(cipherEvents.Events.EventsOfInterest)
			if err != nil {
				return
			}
			clearCensoringEvent, err = survivalAnalysis.Decrypt(cipherEvents.Events.CensoringEvents)
			if err != nil {
				return
			}

			logrus.Infof("Group %s at timepoint %s has %d events of interest and %d censoring events", groupID, clearTimePoint, clearEventOfInterest, clearCensoringEvent)
		}
	}

	survivalAnalysis.addTimer("time for sequential decryption", sequentialDecryptionTimer)
	logrus.Info("Survival analysis request done")
	err = survivalAnalysis.DumpTimers()
	if err != nil {
		return
	}
	//TODO method in profiling.go
	fmt.Println("EXECUTION TIMES")
	fmt.Print(string(survivalAnalysis.profilingBuffer.textBuffer))
	return

}

func parseAndEncryptQueryString(queryString string) (encPanels [][]string, panelsIsNot []bool, err error) {
	panels, panelsIsNot, err := medcoclient.ParseQueryString(queryString)
	if err != nil {
		return
	}
	//from medco-connector client.go
	encPanels = make([][]string, 0)
	for _, panel := range panels {
		encItemKeys := make([]string, 0)
		for _, itemKey := range panel {
			encrypted, innerErr := unlynx.EncryptWithCothorityKey(itemKey)
			if innerErr != nil {
				err = innerErr
				return
			}
			encItemKeys = append(encItemKeys, encrypted)
		}
		encPanels = append(encPanels, encItemKeys)
	}
	return
}

func validateIntermediateResults(patientSetIDs map[int]string, timeCodes []string) (err error) {
	var str string
	if len(patientSetIDs) == 0 {
		str += "empty patient set ID map\n"
	} else {
		for nodeIdx, patientSetID := range patientSetIDs {
			if patientSetID == "" {
				str += fmt.Sprintf("Empty patient set ID for node index %d\n", nodeIdx)
			}
		}
	}

	if len(timeCodes) == 0 {
		str += "empty time code list\n"
	} else {
		for pos, timeCode := range timeCodes {
			if timeCode == "" {
				str += fmt.Sprintf("Empty time code at position %d in time code array\n", pos)
			}
		}
	}
	if str == "" {
		return nil
	}
	return errors.New(str)

}
