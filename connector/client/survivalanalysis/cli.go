package survivalclient

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	medcoclient "github.com/ldsec/medco/connector/client"
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
func ExecuteClientSurvival(token, parameterFileURL, username, password string, disableTLSCheck bool, resultFile string, timerFile string, limit int, cohortName string, granularity, startItem, startsWhen, endItem, endsWhen string) (err error) {

	err = inputValidation(parameterFileURL, limit, cohortName, startItem, endItem)
	if err != nil {
		logrus.Error(err)
		return
	}

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
			var startPanel []*models.PanelItemsItems0
			var endPanel []*models.PanelItemsItems0
			startPanel, err = medcoclient.ParseQueryItem(startItem)
			if err != nil {
				logrus.Error("while parsing start item: ", err.Error())
				return
			}

			err = panelValidation(startPanel)
			if err != nil {
				logrus.Error("while validating start item", err)
				return
			}
			startConcept := *(startPanel[0].QueryTerm)

			endPanel, err = medcoclient.ParseQueryItem(endItem)
			if err != nil {
				logrus.Error("while parsing end item", err)
				return
			}

			err = panelValidation(endPanel)

			if err != nil {
				logrus.Error("while validating start item", err)
				return
			}
			endConcept := *(endPanel[0].QueryTerm)
			var startModifier *modifier
			var endModifier *modifier
			if startMod := startPanel[0].Modifier; startMod != nil {
				startModifier = &modifier{
					ModifierKey: *startMod.ModifierKey,
					AppliedPath: *startMod.AppliedPath,
				}
			}
			if endMod := endPanel[0].Modifier; endMod != nil {
				endModifier = &modifier{
					ModifierKey: *endMod.ModifierKey,
					AppliedPath: *endMod.AppliedPath,
				}
			}

			parameters = &Parameters{
				granularity,
				limit,
				cohortName,
				startConcept,
				startModifier,
				startsWhen,
				endConcept,
				endModifier,
				endsWhen,
				nil,
			}
		}
	}

	// --- convert panels
	timer := time.Now()
	logrus.Info("converting panels")
	panels := convertPanel(parameters)
	logrus.Info("panels converted")
	for _, panel := range panels {
		logrus.Trace(modelPanelsToString(panel))
	}

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
	var APIstartModifier *survival_analysis.SurvivalAnalysisParamsBodyStartModifier
	var APIendModifier *survival_analysis.SurvivalAnalysisParamsBodyEndModifier

	if startMod := parameters.StartModifier; startMod != nil {
		logrus.Debug("start modifier provided")
		APIstartModifier = &survival_analysis.SurvivalAnalysisParamsBodyStartModifier{
			ModifierKey: new(string),
			AppliedPath: new(string),
		}
		*APIstartModifier.ModifierKey = startMod.ModifierKey
		*APIstartModifier.AppliedPath = startMod.AppliedPath
	}

	if endMod := parameters.EndModifier; endMod != nil {
		APIendModifier = &survival_analysis.SurvivalAnalysisParamsBodyEndModifier{
			ModifierKey: new(string),
			AppliedPath: new(string),
		}
		*APIendModifier.ModifierKey = endMod.ModifierKey
		*APIendModifier.AppliedPath = endMod.AppliedPath
	}

	startsWhen, err := parseStartsEndsWhen(parameters.StartsWhen)
	if err != nil {
		logrus.Errorf("when parsing startsWhen argument: %s", err.Error())
		return
	}

	endsWhen, err := parseStartsEndsWhen(parameters.EndsWhen)
	if err != nil {
		logrus.Errorf("when parsing endsWhen argument: %s", err.Error())
		return
	}

	query, err := NewSurvivalAnalysis(
		accessToken,
		parameters.CohortName,
		panels,
		parameters.TimeLimit,
		parameters.TimeResolution,
		parameters.StartConceptPath,
		APIstartModifier,
		startsWhen,
		parameters.EndConceptPath,
		APIendModifier,
		endsWhen,
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
	panels := make([]*survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0, len(parameters.SubGroups))
	for i, selection := range parameters.SubGroups {
		newSelection := &survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0{}
		newSelection.GroupName = fmt.Sprintf(selection.GroupName)
		newPanels := make([]*models.Panel, len(selection.Panels))
		for j, panel := range selection.Panels {
			newPanel := &models.Panel{}
			newPanel.Not = new(bool)
			*newPanel.Not = panel.Not
			newItems := make([]*models.PanelItemsItems0, len(panel.Items))
			for k, item := range panel.Items {
				encrypted := new(bool)
				itemString := new(string)
				*encrypted = false
				*itemString = item.Path
				var modifier *models.PanelItemsItems0Modifier
				if mod := item.Modifier; mod != nil {
					modifier = &models.PanelItemsItems0Modifier{
						ModifierKey: new(string),
						AppliedPath: new(string),
					}
					*modifier.ModifierKey = mod.ModifierKey
					*modifier.AppliedPath = mod.AppliedPath

				}
				newItems[k] = &models.PanelItemsItems0{
					Encrypted: encrypted,
					QueryTerm: itemString,
					Modifier:  modifier,
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
					group.GroupID,
					strconv.FormatInt(group.InitialCount, 10),
					strconv.FormatInt(timePoint.Time, 10),
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

	err = medcoclient.DumpTimers(timerFile, timers, clientTimers)
	if err != nil {
		err = fmt.Errorf("while dumping timers: %s", err.Error())
		logrus.Error(err)
		return
	}

	logrus.Info("timers dumped")
	return

}

func inputValidation(parameterFileURL string, limit int, cohortName, startConcept, endConcept string) error {
	if parameterFileURL == "" {
		if limit == 0 {
			return fmt.Errorf("Limit -l is not set")
		}
		if cohortName == "" {
			return fmt.Errorf("Cohort name -c is not set")
		}
		if startConcept == "" {
			return fmt.Errorf("Start concept path -s is not set")
		}
		if endConcept == "" {
			return fmt.Errorf("End concept path -e is not set")
		}
	}
	return nil
}

func panelValidation(panel []*models.PanelItemsItems0) (err error) {
	if len(panel) == 0 {
		err = fmt.Errorf("panels are empty. Was the start item string empry ?")
		logrus.Error(err)
		return
	}
	if len(panel) > 1 {
		err = fmt.Errorf("multiple items retrieved from the item string. Was a file provided ? Only clear concept should be used")
		logrus.Error(err)
		return
	}
	if *(panel[0].Encrypted) {
		err = fmt.Errorf("encrypted concept found, only clear concept should be used")
		logrus.Error(err)
	}
	return
}

func modelPanelsToString(subGroup *survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0) string {

	panelStrings := make([]string, 0, len(subGroup.Panels))
	for _, panel := range subGroup.Panels {
		itemStrings := make([]string, 0, len(panel.Items))
		for _, item := range panel.Items {
			itemStrings = append(itemStrings, fmt.Sprintf("{Encrypted:%t Modifier:%v Operator:%s QueryTerm:%s Value:%s}",
				*item.Encrypted,
				item.Modifier,
				item.Operator,
				*item.QueryTerm,
				item.Value))
		}
		itemArray := "[" + strings.Join(itemStrings, " ") + "]"
		panelStrings = append(panelStrings, fmt.Sprintf("{Items:%s Not:%t}", itemArray, *panel.Not))
	}
	panelArray := "[" + strings.Join(panelStrings, " ") + "]"
	return fmt.Sprintf("{GroupName:%s Panels:%s", subGroup.GroupName, panelArray)
}
