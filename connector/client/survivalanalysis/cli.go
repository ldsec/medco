package survivalclient

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/ldsec/medco/connector/restapi/client/survival_analysis"
	"github.com/ldsec/medco/connector/restapi/models"
	utilclient "github.com/ldsec/medco/connector/util/client"
	utilcommon "github.com/ldsec/medco/connector/util/common"
	"github.com/sirupsen/logrus"
)

// time out in seconds for  access token read parameter file, seconds
const accessToketimeout time.Duration = 300

// time out in seconds for survival analysis result, seconds
const resultTimeout time.Duration = 900

// print ticking while waiting for the results, seconds
const tick time.Duration = 5

// ClientResultElement holds the information for the CLI whole susrvival analysis loop
type ClientResultElement struct {
	ClearTimePoint     string
	EncEventOfInterest string
	EncCensoringEvent  string
}

// ExecuteClientSurvival creates a survival analysis form parameters given in parameter file, and makes a call to the API to executes this query
func ExecuteClientSurvival(token, parameterFileURL, username, password string, disableTLSCheck bool, resultFile string, timerFile string, limit int, cohortName string, granularity, startConcept, startModifier, endConcept, endModifier string) (err error) {

	//initialize onjects and channels
	clientTimers := utilcommon.NewTimers()
	var accessToken string
	var parameters *Parameters
	tokenChan := make(chan string, 1)
	parametersChan := make(chan *Parameters, 1)
	errChan := make(chan error)
	signal := make(chan struct{})
	resultChan := make(chan struct {
		Results []EncryptedResults
		Timers  []utilcommon.Timers
	})
	wait := &sync.WaitGroup{}
	wait.Add(2)
	go func() {
		wait.Wait()
		signal <- struct{}{}
	}()

	// --- get token
	logrus.Debug("requesting access token")
	go func() {
		defer wait.Done()
		accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
		if err != nil {
			errChan <- err
			return
		}
		tokenChan <- accessToken
		logrus.Debug("access token received")
		logrus.Tracef("token %s", accessToken)
		return

	}()

	// --- get parameters
	if parameterFileURL != "" {
		logrus.Debugf("reading parameters")
		go func() {
			defer wait.Done()
			parameters, err := NewParametersFromFile(parameterFileURL)
			if err != nil {
				errChan <- err
				return
			}
			parametersChan <- parameters
			logrus.Debug("parameters read")
			logrus.Tracef("parameters %+v", parameters)
			return
		}()
	} else {
		wait.Done()
	}

	select {
	case err = <-errChan:
		logrus.Error(err)
		return
	case <-time.After(accessToketimeout * time.Second):
		err = fmt.Errorf("timeout %d seconds", accessToketimeout)
		logrus.Error(err)
		return

	case <-signal:
		accessToken = <-tokenChan
		if parameterFileURL != "" {
			parameters = <-parametersChan
		} else {
			parameters = &Parameters{
				granularity,
				limit,
				cohortName,
				startConcept,
				startModifier,
				endConcept,
				endModifier,
				nil,
			}

		}
	}

	// --- convert panels
	timer := time.Now()
	logrus.Debug("converting panels")
	panels := make([]*survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0, len(parameters.Cohorts))
	for i, selection := range parameters.Cohorts {
		newSelection := &survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0{}
		newSelection.CohortName = fmt.Sprintf("SUB_GROUP_%d", i)
		newPanels := make([]*models.Panel, len(selection.Panels))
		for j, panel := range selection.Panels {
			newPanel := &models.Panel{}
			newPanel.Not = new(bool)
			*newPanel.Not = panel.Not
			newItems := make([]*models.PanelItemsItems0, len(panel.Paths))
			for k, item := range panel.Paths {
				encrypted := new(bool)
				itemString := new(string)
				*encrypted = false
				*itemString = item
				newItems[k] = &models.PanelItemsItems0{
					Encrypted: encrypted,
					Operator:  "exists",
					QueryTerm: itemString,
				}
			}
			newPanel.Items = newItems
			newPanels[j] = newPanel
		}
		newSelection.Panels = newPanels
		panels[i] = newSelection
	}
	clientTimers.AddTimers("medco-connector-panel-conversion", timer, nil)
	logrus.Debug("panels converted")
	logrus.Tracef("panels: %+v", panels)

	// --- execute query
	timer = time.Now()
	logrus.Debug("executing query")
	query, err := NewSurvivalAnalysis(
		accessToken,
		parameters.CohortName,
		panels,
		parameters.TimeLimit,
		parameters.TimeResolution,
		parameters.StartConceptPath,
		parameters.StartConceptModifier,
		parameters.EndConceptPath,
		parameters.EndConceptModifier,
		disableTLSCheck,
	)
	if err != nil {
		return
	}

	resTimeout := time.After(resultTimeout * time.Second)
	resultTicks := time.Tick(tick * time.Second)
	go func() {
		results, timers, err := query.Execute()
		if err != nil {
			logrus.Error(err)
			errChan <- err
			return
		}
		resultChan <- struct {
			Results []EncryptedResults
			Timers  []utilcommon.Timers
		}{Results: results, Timers: timers}
		return

	}()
	var results []EncryptedResults
	var timers []utilcommon.Timers

	tickTime := 0
resLoop:
	for {
		select {
		case <-resTimeout:
			err = fmt.Errorf("Timeout %d", resultTimeout)
		case err = <-errChan:
			logrus.Error(err)
			return
		case res := <-resultChan:
			results = res.Results
			timers = res.Timers
			break resLoop
		case <-resultTicks:
			tickTime += int(tick)
			logrus.Debugf("waiting for response (%d seconds)", tickTime)
		}
	}
	clientTimers.AddTimers("medco-connector-survival-query-remote-execution", timer, nil)
	logrus.Debug("query executed")
	logrus.Tracef("encrypted results: %+v", results)
	logrus.Tracef("timers: %v", timers)

	// --- decrypt result
	timer = time.Now()
	logrus.Debug("decrypting results")
	clearResults := make([]ClearResults, len(results))
	for idx, encryptedResults := range results {
		clearResults[idx], err = encryptedResults.Decrypt(query.userPrivateKey)
		if err != nil {
			err = fmt.Errorf("while decrypting survival analysis results: %s", err.Error())
			logrus.Error(err)
			return
		}
	}
	clientTimers.AddTimers("medco-connector-decryptions", timer, nil)
	logrus.Debug("results decrypted")
	logrus.Tracef("clear results: %+v", clearResults)

	// --- printing results
	logrus.Debug("printing results")
	csv, err := utilclient.NewCSV(resultFile)
	if err != nil {
		err = fmt.Errorf("while creating CSV file handler: %s", err)
		logrus.Error(err)
		return
	}
	err = csv.Write([]string{"time_granularity", "node_index", "group_id", "initial_count", "time_point", "event_of_interest_count", "censoring_event_count"})
	if err != nil {
		err = fmt.Errorf("while writing result headers:%s", err.Error())
		logrus.Error(err)
		return
	}
	for nodeIdx := range clearResults {
		sort.Sort(clearResults[nodeIdx])
		for groupIdx := range clearResults[nodeIdx] {

			sort.Sort(clearResults[nodeIdx][groupIdx].TimePoints)
			var group = clearResults[nodeIdx][groupIdx]
			for _, timePoint := range group.TimePoints {
				csv.Write([]string{
					parameters.TimeResolution,
					strconv.Itoa(nodeIdx),
					strconv.Itoa(groupIdx),
					strconv.FormatInt(group.InitialCount, 10),
					strconv.Itoa(timePoint.Time),
					strconv.FormatInt(timePoint.Events.EventsOfInterest, 10),
					strconv.FormatInt(timePoint.Events.CensoringEvents, 10),
				})
				if err != nil {
					err = fmt.Errorf("while writing a record: %s", err.Error())
					logrus.Error(err)
					return
				}
			}

		}

	}
	err = csv.Flush()
	if err != nil {
		err = fmt.Errorf("while flushing buffer to result file: %s", err.Error())
		logrus.Error(err)
		return
	}

	err = csv.Close()
	if err != nil {
		err = fmt.Errorf("while closing result file: %s", err.Error())
		logrus.Error(err)
		return
	}
	logrus.Debug("results printed")

	// print timers
	logrus.Debug("dumping timers")
	dumpCSV, err := utilclient.NewCSV(timerFile)
	if err != nil {
		err = fmt.Errorf("while creating CSV file handler: %s", err)
		logrus.Error(err)
		return
	}
	dumpCSV.Write([]string{"node_index", "timer_description", "duration_milliseconds"})
	if err != nil {
		err = fmt.Errorf("while writing headers for timer file: %s", err)
		logrus.Error(err)
		return
	}
	// each remote time profilings
	for nodeIdx, nodeTimers := range timers {
		sortedTimers := utilclient.SortTimers(nodeTimers)
		for _, duration := range sortedTimers {
			dumpCSV.Write([]string{
				strconv.Itoa(nodeIdx),
				duration[0],
				duration[1],
			})
			if err != nil {
				err = fmt.Errorf("while writing record for timer file: %s", err)
				logrus.Error(err)
				return
			}
		}

	}
	// and local
	localSortedTimers := utilclient.SortTimers(clientTimers)
	for _, duration := range localSortedTimers {
		dumpCSV.Write([]string{
			"client",
			duration[0],
			duration[1],
		})
		if err != nil {
			err = fmt.Errorf("while writing record for timer file: %s", err)
			logrus.Error(err)
			return
		}
	}

	err = dumpCSV.Flush()
	if err != nil {
		err = fmt.Errorf("while flushing timer file: %s", err)
		logrus.Error(err)
		return
	}
	err = dumpCSV.Close()
	if err != nil {
		err = fmt.Errorf("while closing timer file: %s", err)
		logrus.Error(err)
		return
	}

	logrus.Debug("timers dumped")

	return

}
