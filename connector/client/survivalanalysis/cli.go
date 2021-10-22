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

const (
	defaultTiming = models.TimingAny
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
	logrus.Info("Survival analysis: requesting access token")
	go func() {
		defer wait.Done()
		accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
		if err != nil {
			errChan <- err
			return
		}
		tokenChan <- accessToken
		logrus.Info("Survival analysis: access token received")
		logrus.Tracef("Survival analysis: token %s", accessToken)
		return

	}()

	// --- get parameters
	if parameterFileURL != "" {
		logrus.Info("Survival analysis: reading parameters")
		go func() {
			defer wait.Done()
			parameters, err := NewParametersFromFile(parameterFileURL)
			if err != nil {
				errChan <- err
				return
			}
			err = validateUserIntputSequenceOfEvents(parameters)
			if err != nil {
				errChan <- err
				return
			}
			parametersChan <- parameters
			logrus.Info("Survival analysis: parameters read")
			logrus.Tracef("Survival analysis: parameters %+v", parameters)
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
			var startPanel []*models.PanelConceptItemsItems0
			var endPanel []*models.PanelConceptItemsItems0
			startPanel, _, err = medcoclient.ParseQueryItem(startItem)
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

			endPanel, _, err = medcoclient.ParseQueryItem(endItem)
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
	logrus.Info("Survival analysis: converting panels")
	panels, err := convertParametersToSubGroupDefinition(parameters)
	if err != nil {
		err = fmt.Errorf("while converting panels: %s", err.Error())
		logrus.Error(err)
		return
	}
	logrus.Info("Survival analysis: panels converted")
	for _, panel := range panels {
		logrus.Trace(modelPanelsToString(panel))
	}

	// --- execute query
	timer = time.Now()
	logrus.Info("Survival analysis: executing query")
	results, timers, userPrivateKey, err := executeQuery(accessToken, panels, parameters, disableTLSCheck)
	if err != nil {
		err = fmt.Errorf("while executing survival analysis results: %s", err.Error())
		logrus.Error(err)
		return
	}
	clientTimers.AddTimers("medco-connector-survival-query-remote-execution", timer, nil)
	logrus.Info("Survival analysis: query executed")
	logrus.Tracef("Survival analysis: encrypted results: %+v", results)
	logrus.Tracef("Survival analysis: timers: %v", timers)

	// --- decrypt result
	timer = time.Now()
	logrus.Info("Survival analysis: decrypting results")
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
	logrus.Info("Survival analysis: results decrypted")
	logrus.Tracef("Survival analysis: clear results: %+v", clearResults)

	// --- printing results
	printResults(clearResults, timers, clientTimers, parameters.TimeResolution, resultFile, timerFile)

	logrus.Info("Survival analysis: operation completed")
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
		logrus.Debug("Survival analysis: start modifier provided")
		APIstartModifier = &survival_analysis.SurvivalAnalysisParamsBodyStartModifier{
			ModifierKey: new(string),
			AppliedPath: new(string),
		}
		*APIstartModifier.ModifierKey = startMod.ModifierKey
		*APIstartModifier.AppliedPath = startMod.AppliedPath
	}

	if endMod := parameters.EndModifier; endMod != nil {
		logrus.Debug("Survival analysis: end modifier provided")
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

func validateUserIntputSequenceOfEvents(parameters *Parameters) error {
	for _, subGroup := range parameters.SubGroups {
		err := validateSequenceOfEvents(subGroup.SequenceOfEvents, len(subGroup.Panels))
		if err != nil {
			return err
		}
	}
	return nil
}

func convertParametersToSubGroupDefinition(parameters *Parameters) ([]*survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0, error) {
	panels := make([]*survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0, len(parameters.SubGroups))
	var err error
	for i, selection := range parameters.SubGroups {
		newSelection := &survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0{}
		newSelection.GroupName = selection.GroupName
		newSelection.SubGroupTiming, err = timingFromStringToModel(selection.GroupTiming)
		if err != nil {
			err = fmt.Errorf("while parsing sub group timing: %s", err.Error())
			return nil, err
		}
		newPanels := make([]*models.Panel, len(selection.Panels))
		for j, panel := range selection.Panels {
			newPanel := &models.Panel{}
			newPanel.PanelTiming, err = timingFromStringToModel(panel.PanelTiming)
			if err != nil {
				err = fmt.Errorf("while parsing panel timing: %s", err.Error())
				return nil, err
			}
			newPanel.Not = new(bool)
			*newPanel.Not = panel.Not
			newConceptItems := make([]*models.PanelConceptItemsItems0, len(panel.ConceptItems))
			for k, conceptItem := range panel.ConceptItems {
				encrypted := new(bool)
				itemString := new(string)
				*encrypted = false
				*itemString = conceptItem.Path
				var modifier *models.PanelConceptItemsItems0Modifier
				if mod := conceptItem.Modifier; mod != nil {
					modifier = &models.PanelConceptItemsItems0Modifier{
						ModifierKey: new(string),
						AppliedPath: new(string),
					}
					*modifier.ModifierKey = mod.ModifierKey
					*modifier.AppliedPath = mod.AppliedPath

				}
				newConceptItems[k] = &models.PanelConceptItemsItems0{
					Encrypted: encrypted,
					QueryTerm: itemString,
					Modifier:  modifier,
					Operator:  conceptItem.Operator,
					Type:      conceptItem.Type,
					Value:     conceptItem.Value,
				}
			}

			newPanel.CohortItems = append(newPanel.CohortItems, panel.CohortItems...)

			newPanel.ConceptItems = newConceptItems

			newPanels[j] = newPanel
		}
		newSelection.Panels = newPanels
		sequenceOfEvents := defaultedSequenceOfEvents(selection.SequenceOfEvents, len(selection.Panels))
		newSelection.QueryTimingSequence, err = convertParametersToSequenceInfo(sequenceOfEvents)
		if err != nil {
			err = fmt.Errorf("while parsing temporal sequence info: %s", err.Error())
			return nil, err
		}
		panels[i] = newSelection
	}
	return panels, nil
}

func printResults(clearResults []ClearResults, timers []medcomodels.Timers, clientTimers medcomodels.Timers, timeResolution, resultFile, timerFile string) (err error) {
	logrus.Info("Survival analysis: printing results")
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
	logrus.Info("Survival analysis: results printed")

	err = medcoclient.DumpTimers(timerFile, timers, clientTimers)
	if err != nil {
		err = fmt.Errorf("while dumping timers: %s", err.Error())
		logrus.Error(err)
		return
	}

	logrus.Info("Survival analysis: timers dumped")
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

func panelValidation(panel []*models.PanelConceptItemsItems0) (err error) {
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
		itemStrings := make([]string, 0, len(panel.ConceptItems))
		for _, item := range panel.ConceptItems {
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
	return fmt.Sprintf("{GroupName:%s QueryTiming:%s Panels:%s", subGroup.GroupName, subGroup.SubGroupTiming, panelArray)
}

func timingFromStringToModel(timingString string) (models.Timing, error) {
	switch candidate := strings.ToLower(strings.TrimSpace(timingString)); candidate {
	case string(models.TimingAny):
		return models.TimingAny, nil
	case "":
		return defaultTiming, nil
	case string(models.TimingSameinstancenum):
		return models.TimingSameinstancenum, nil
	case string(models.TimingSamevisit):
		return models.TimingSamevisit, nil
	default:
		return "", fmt.Errorf("candidate %s is not implemented, must be one of any, sameinstancenum, samevisit", timingString)
	}
}
