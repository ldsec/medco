package survivalclient

import (
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	medcomodels "github.com/ldsec/medco/connector/models"

	"github.com/ldsec/medco/connector/restapi/client/survival_analysis"
	"github.com/ldsec/medco/connector/restapi/models"
	utilclient "github.com/ldsec/medco/connector/util/client"
	"github.com/sirupsen/logrus"
)

// ClientResultElement holds the information for the CLI whole susrvival analysis loop
type ClientResultElement struct {
	ClearTimePoint     string
	EncEventOfInterest string
	EncCensoringEvent  string
}

// ExecuteClientSurvival creates a survival analysis form parameters given in parameter file, and makes a call to the API to executes this query
func ExecuteClientSurvival(token, parameterFileURL, username, password string, disableTLSCheck bool, resultFile string, timerFile string, limit int, cohortName string, granularity, startConcept, startModifier, endConcept, endModifier string) (err error) {

	//initialize objects and channels
	clientTimers := medcomodels.NewTimers()
	var accessToken string
	var parameters *Parameters
	tokenChan := make(chan string, 1)
	parametersChan := make(chan *Parameters, 1)
	errChan := make(chan error)
	signal := make(chan struct{})
	wait := &sync.WaitGroup{}
	wait.Add(2)
	go func() {
		wait.Wait()
		signal <- struct{}{}
	}()

	// --- get token
	logrus.Info("requesting access token")
	go func() {
		defer wait.Done()
		accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
		if err != nil {
			errChan <- err
			return
		}
		tokenChan <- accessToken
		logrus.Info("access token received")
		logrus.Tracef("token %s", accessToken)
		return

	}()

	// --- get parameters
	if parameterFileURL != "" {
		logrus.Info("reading parameters")
		go func() {
			defer wait.Done()
			parameters, err := NewParametersFromFile(parameterFileURL)
			if err != nil {
				errChan <- err
				return
			}
			parametersChan <- parameters
			logrus.Info("parameters read")
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
	case <-time.After(time.Duration(utilclient.TokenTimeoutSeconds) * time.Second):
		err = fmt.Errorf("timeout %d seconds", utilclient.TokenTimeoutSeconds)
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
	logrus.Info("converting panels")
	panels := convertPanel(parameters)
	logrus.Info("panels converted")
	logrus.Tracef("panels: %+v", panels)

	// --- execute query
	timer = time.Now()
	logrus.Info("executing query")
	results, timers, userPrivateKey, err := executeQuery(accessToken, panels, parameters, disableTLSCheck)
	if err != nil {
		err = fmt.Errorf("while executing survival analysis results: %s", err.Error())
		logrus.Error(err)
		return
	}
	clientTimers.AddTimers("medco-connector-survival-query-remote-execution", timer, nil)
	logrus.Info("query executed")
	logrus.Tracef("encrypted results: %+v", results)
	logrus.Tracef("timers: %v", timers)

	// --- decrypt result
	timer = time.Now()
	logrus.Info("decrypting results")
	clearResults := make([]ClearResults, len(results))
	for idx, encryptedResults := range results {
		clearResults[idx], err = encryptedResults.Decrypt(userPrivateKey)
		if err != nil {
			err = fmt.Errorf("while decrypting survival analysis results: %s", err.Error())
			logrus.Error(err)
			return
		}
	}
	clientTimers.AddTimers("medco-connector-decryptions", timer, nil)
	logrus.Info("results decrypted")
	logrus.Tracef("clear results: %+v", clearResults)

	// --- printing results
	printResults(clearResults, timers, clientTimers, parameters.TimeResolution, resultFile, timerFile)

	logrus.Info("Operation completed")
	return

}

func executeQuery(accessToken string, panels []*survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0, parameters *Parameters, disableTLSCheck bool) (results []EncryptedResults, timers []medcomodels.Timers, userPrivateKey string, err error) {
	errChan := make(chan error)
	resultChan := make(chan struct {
		Results []EncryptedResults
		Timers  []medcomodels.Timers
	})
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
	userPrivateKey = query.userPrivateKey
	if err != nil {
		return
	}

	resTimeout := time.After(time.Duration(utilclient.SurvivalAnalysisTimeoutSeconds) * time.Second)
	resultTicks := time.Tick(time.Duration(utilclient.WaitTickSeconds) * time.Second)
	go func() {
		results, timers, err := query.Execute()
		if err != nil {
			logrus.Error(err)
			errChan <- err
			return
		}
		resultChan <- struct {
			Results []EncryptedResults
			Timers  []medcomodels.Timers
		}{Results: results, Timers: timers}
		return

	}()

	tickTime := 0
resLoop:
	for {
		select {
		case <-resTimeout:
			err = fmt.Errorf("Timeout %d", utilclient.SurvivalAnalysisTimeoutSeconds)
			return
		case err = <-errChan:
			logrus.Error(err)
			return
		case res := <-resultChan:
			results = res.Results
			timers = res.Timers
			break resLoop
		case <-resultTicks:
			tickTime += int(utilclient.WaitTickSeconds)
			logrus.Infof("waiting for response (%d seconds)", tickTime)
		}
	}

	return

}

func convertPanel(parameters *Parameters) []*survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0 {
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
					Operator:  models.PanelItemsItems0OperatorExists,
					QueryTerm: itemString,
				}
			}
			newPanel.Items = newItems
			newPanels[j] = newPanel
		}
		newSelection.Panels = newPanels
		panels[i] = newSelection
	}
	return panels
}

func printResults(clearResults []ClearResults, timers []medcomodels.Timers, clientTimers medcomodels.Timers, timeResolution, resultFile, timerFile string) (err error) {
	logrus.Info("printing results")
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
					timeResolution,
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
	logrus.Info("results printed")

	// print timers
	logrus.Info("dumping timers")
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
		sortedTimers := nodeTimers.SortTimers()
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
	localSortedTimers := clientTimers.SortTimers()
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
	logrus.Info()
	err = dumpCSV.Close()
	if err != nil {
		err = fmt.Errorf("while closing timer file: %s", err)
		logrus.Error(err)
		return
	}

	logrus.Info("timers dumped")
	return

}
