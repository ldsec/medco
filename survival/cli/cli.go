package survivalcli

import (
	"errors"
	"sync"
	"time"

	medcoclient "github.com/ldsec/medco-connector/client"
	survivalclient "github.com/ldsec/medco-connector/survival/client"
	"github.com/ldsec/medco-connector/wrappers/unlynx"
	"github.com/sirupsen/logrus"
)

type timeCodesMaps [2]map[string]string

type ClientResultElement struct {
	ClearTimePoint     string
	EncEventOfInterest string
	EncCensoringEvent  string
}

func ClientSurvival(token string, granularity string, limit int64, queryString, username, password string, disableTLSCheck bool) (err error) {

	accessToken, err := survivalclient.GetToken(token, username, password, disableTLSCheck)
	if err != nil {
		return
	}

	var patientSetIDs map[int]string
	var encTimeCodesMap map[string]string
	var encTimeCodesInverseMap map[string]string

	errChan := make(chan error)
	patientSetIDsChan := make(chan map[int]string)
	encTimeCodesMapChan := make(chan timeCodesMaps)
	endChan := make(chan struct{})

	var waitGroup sync.WaitGroup
	waitGroup.Add(2)

	go func() {
		defer waitGroup.Done()

		encPanels, panelsIsNot, err := ParseAndEncryptQueryString(queryString)
		if err != nil {
			errChan <- err
			return
		}
		patientSetIDs, err := survivalclient.GetPatientList(accessToken, encPanels, panelsIsNot, disableTLSCheck)
		if err != nil {
			errChan <- err
			return
		}
		patientSetIDsChan <- patientSetIDs
	}()
	go func() {
		defer waitGroup.Done()
		timeCodes, err := survivalclient.GetTimeCodes(accessToken, granularity, limit, disableTLSCheck)
		if err != nil {
			errChan <- err
			return
		}
		encTimeCodesMap, encTimeCodesInverseMap, err := survivalclient.EncryptTimeCodes(timeCodes)
		if err != nil {
			errChan <- err
			return
		}
		encTimeCodesMapChan <- timeCodesMaps{encTimeCodesMap, encTimeCodesInverseMap}
	}()

	go func() {
		waitGroup.Wait()
		endChan <- struct{}{}
	}()

	select {
	case <-endChan:
		select {
		case patientSetIDs = <-patientSetIDsChan:
		case <-time.After(time.Duration(300) * time.Second):
			err = errors.New("Unexpected empty channel")
			logrus.Error(err)
		}
		select {
		case maps := <-encTimeCodesMapChan:
			encTimeCodesMap = maps[0]
			encTimeCodesInverseMap = maps[1]
		case <-time.After(time.Duration(300) * time.Second):
			err = errors.New("Unexpected empty channel")
			logrus.Error(err)

		}
	case err = <-errChan:
		return
	}

	var encTimeCodes []string
	for _, encTimeCode := range encTimeCodesMap {
		encTimeCodes = append(encTimeCodes, encTimeCode)
	}

	survivalAnalysis, err := survivalclient.NewSurvivalAnalysis(accessToken, patientSetIDs, encTimeCodes, disableTLSCheck)
	if err != nil {
		return
	}
	survResults, err := survivalAnalysis.Execute()
	if err != nil {
		return
	}
	var results []ClientResultElement

	for _, encElement := range survResults {
		clearTimePoint := encTimeCodesInverseMap[encElement.TimePoint]
		results = append(results, ClientResultElement{ClearTimePoint: clearTimePoint,
			EncEventOfInterest: encElement.Events.EventsOfInterest,
			EncCensoringEvent:  encElement.Events.CensoringEvents,
		})
	}
	return

}

func ParseAndEncryptQueryString(queryString string) (encPanels [][]string, panelsIsNot []bool, err error) {
	panels, panelsIsNot, err := medcoclient.ParseQueryString(queryString)
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
