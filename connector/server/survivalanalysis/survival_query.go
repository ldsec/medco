package survivalserver

import (
	"encoding/base64"
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

	medcomodels "github.com/ldsec/medco/connector/models"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	querytoolsserver "github.com/ldsec/medco/connector/server/querytools"
	"github.com/ldsec/medco/connector/server/survivalanalysis/timepoints"
	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/ldsec/medco/connector/wrappers/i2b2"
	"github.com/ldsec/medco/connector/wrappers/unlynx"

	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/ldsec/medco/connector/restapi/server/operations/survival_analysis"
)

// Query holds the ID of the survival analysis, its parameters and a pointer to its results
type Query struct {
	UserID              string
	UserPublicKey       string
	QueryName           string
	CohortName          string
	SubGroupDefinitions []*survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0
	TimeLimit           int
	TimeGranularity     string
	StartConcept        string
	StartModifier       *survival_analysis.SurvivalAnalysisParamsBodyStartModifier
	StartsWhen          *string
	EndConcept          string
	EndModifier         *survival_analysis.SurvivalAnalysisParamsBodyEndModifier
	EndsWhen            *string
	Result              *struct {
		Timers    medcomodels.Timers
		EncEvents EventGroups
	}
}

// NewQuery query constructor
func NewQuery(UserID,
	QueryName,
	UserPublicKey string,
	CohortName string,
	SubGroupDefinitions []*survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0,
	TimeLimit int,
	TimeGranularity string,
	StartConcept string,
	StartModifier *survival_analysis.SurvivalAnalysisParamsBodyStartModifier,
	StartsWhen *string,
	EndConcept string,
	EndModifier *survival_analysis.SurvivalAnalysisParamsBodyEndModifier,
	EndsWhen *string) *Query {
	res := &Query{
		UserPublicKey:       UserPublicKey,
		UserID:              UserID,
		QueryName:           QueryName,
		CohortName:          CohortName,
		SubGroupDefinitions: SubGroupDefinitions,
		TimeLimit:           TimeLimit,
		TimeGranularity:     TimeGranularity,
		StartConcept:        StartConcept,
		StartModifier:       StartModifier,
		StartsWhen:          StartsWhen,
		EndConcept:          EndConcept,
		EndModifier:         EndModifier,
		EndsWhen:            EndsWhen,
		Result: &struct {
			Timers    medcomodels.Timers
			EncEvents EventGroups
		}{}}
	if StartModifier != nil {
		logrus.Debugf("Survival analysis: provided start modifier with key %s and applied path %s", *StartModifier.ModifierKey, *StartModifier.AppliedPath)
	} else {
		logrus.Debug("Survival analysis: no start modifier provided")
	}
	if EndModifier != nil {
		logrus.Debugf("Survival analysis: provided end modifier with key %s and applied path %s", *EndModifier.ModifierKey, *EndModifier.AppliedPath)
	} else {
		logrus.Debug("Survival analysis: no end modifier provided")
	}
	res.Result.Timers = make(map[string]time.Duration)

	return res
}

// Execute runs the survival analysis query
func (q *Query) Execute() error {

	eventGroups := make(EventGroups, 0)
	timeLimitInDays := q.TimeLimit * granularityValues[q.TimeGranularity]
	timer := time.Now()

	startConceptCodes, startModifierCodes, endConceptCodes, endModifierCodes, cohort, timers, err := prepareArguments(q.UserID, q.CohortName, q.StartConcept, q.StartModifier, q.EndConcept, q.EndModifier)
	if err != nil {
		err = fmt.Errorf("while retrieving concept codes and patient indices: %s", err.Error())
		return err
	}
	q.Result.Timers.AddTimers("", timer, timers)

	// --- build sub groups

	definitions := q.SubGroupDefinitions
	if q.SubGroupDefinitions == nil || len(q.SubGroupDefinitions) == 0 {
		definitions = fullCohort(q.StartConcept)
	}
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(len(definitions))
	channels := make([]chan struct {
		*EventGroup
		medcomodels.Timers
	}, len(definitions))
	errChan := make(chan error, len(definitions))
	signal := make(chan struct{})

	for i, definition := range definitions {
		channels[i] = make(chan struct {
			*EventGroup
			medcomodels.Timers
		}, 1)
		go func(i int, definition *survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0) {
			defer waitGroup.Done()
			timers := medcomodels.NewTimers()

			newEventGroup := &EventGroup{GroupID: definition.GroupName}

			not := false
			item := &models.PanelConceptItemsItems0{
				QueryTerm: &q.StartConcept,
			}
			newPanel := &models.Panel{
				ConceptItems: []*models.PanelConceptItemsItems0{
					item,
				},
				Not:         &not,
				PanelTiming: models.TimingAny,
			}
			panels := append(definition.Panels, newPanel)

			timer = time.Now()
			logrus.Infof("Survival analysis: I2B2 explore for subgroup %d", i)
			logrus.Tracef("Survival analysis: panels %+v", panels)
			patientList, err := SubGroupExplore(q.QueryName, i, panels, definition.SubGroupTiming)
			if err != nil {
				returnedErr := fmt.Errorf("during subgroup explore procedure")
				logrus.Errorf("%s: %s", returnedErr.Error(), err.Error())
				errChan <- returnedErr
				return
			}
			logrus.Infof("Survival analysis: successful I2B2 explore query %d", i)
			timers.AddTimers(fmt.Sprintf("medco-connector-i2b2-query-group%d", i), timer, nil)
			patientList = intersect(cohort, patientList)
			logrus.Debugf("Survival analysis: there are %d patients in the subgroup", len(patientList))

			// --- build time points

			timepointsEventsMap, patientWithoutStartEvent, patientWithoutEndEvent, timepointTimers, err := timepoints.BuildTimePoints(
				patientList,
				startConceptCodes,
				startModifierCodes,
				*q.StartsWhen == survival_analysis.SurvivalAnalysisBodyStartsWhenEarliest,
				endConceptCodes,
				endModifierCodes,
				*q.EndsWhen == survival_analysis.SurvivalAnalysisBodyEndsWhenEarliest,
				int64(timeLimitInDays),
			)

			if err != nil {
				// error details dumped in server node console, but survival request does not send them to the client
				logrus.Errorf("error while getting building time points: %s", err.Error())
				err = fmt.Errorf("error while getting building time points")
				errChan <- err
				return
			}
			timers.AddTimers("", time.Now(), timepointTimers)
			logrus.Debugf("Survival analysis: found %d patients without the start event", len(patientWithoutStartEvent))
			logrus.Debugf("Survival analysis: found %d patients without the end (censoring or of interest) event", len(patientWithoutEndEvent))
			timePoints := timePointMapToList(timepointsEventsMap)

			// --- initial count
			if len(patientList) < len(patientWithoutStartEvent) {
				logrus.Errorf("length of the patient list %d cannot be smaller than this of patients without start event %d", len(patientList), len(patientWithoutStartEvent))
				err = fmt.Errorf("while computing initial count")
				errChan <- err
				return
			}

			timer = time.Now()
			initialCount := int64(len(patientList) - len(patientWithoutStartEvent))
			initialCountEncrypt, err := unlynx.EncryptWithCothorityKey(initialCount)
			timers.AddTimers(fmt.Sprintf("medco-connector-encrypt-init-count-group%d", i), timer, nil)
			if err != nil {
				logrus.Errorf("while encrypting initial count %d: %s", initialCount, err.Error())
				err = fmt.Errorf("while encrypting initial count")
				errChan <- err
				return
			}
			logrus.Debug("Survival analysis: initialcount ", initialCountEncrypt)
			newEventGroup.EncInitialCount = initialCountEncrypt

			timer = time.Now()

			// change time granularity, fill zeros in arrays and encrypt group results
			processGroupResultTimers := processGroupResult(errChan, newEventGroup, timePoints, timeLimitInDays, q.TimeGranularity, i)
			q.Result.Timers.AddTimers("", timer, processGroupResultTimers)

			logrus.Tracef("Survival analysis: event groups %v", newEventGroup)
			timers.AddTimers(fmt.Sprintf("medco-connector-local-encryption%d", i), timer, nil)
			channels[i] <- struct {
				*EventGroup
				medcomodels.Timers
			}{newEventGroup, timers}
		}(i, definition)

	}
	go func() {
		waitGroup.Wait()
		signal <- struct{}{}
	}()
	select {
	case err := <-errChan:
		return err
	case <-signal:
		break
	}
	for _, channel := range channels {
		chanResult := <-channel

		eventGroups = append(eventGroups, chanResult.EventGroup)
		q.Result.Timers.AddTimers("", timer, chanResult.Timers)
	}

	// aggregate and key switch locally encrypted results

	for _, group := range eventGroups {
		logrus.Tracef("Survival analysis: eventGroup %v", group)
	}
	timer = time.Now()
	var aksTimers medcomodels.Timers
	q.Result.EncEvents, aksTimers, err = AKSgroups(q.QueryName+"_AGG_AND_KEYSWITCH", eventGroups, q.UserPublicKey)
	q.Result.Timers.AddTimers("medco-connector-aggregate-and-key-switch", timer, aksTimers)
	if err != nil {
		logrus.Errorf("during aggregation and keyswitch: %s", err.Error())
		err = fmt.Errorf("during aggregation and keyswitch")
	}
	return err
}

// Validate checks members of a Query instance for early error detection.
// Heading and trailing spaces are silently trimmed. Granularity string is silently written in lower case.
// If any other wrong member can be defaulted, a warning message is printed, otherwise an error is returned.
func (q *Query) Validate() error {

	q.QueryName = strings.TrimSpace(q.QueryName)
	if q.QueryName == "" {
		return fmt.Errorf("empty query name")
	}

	q.StartConcept = strings.TrimSpace(q.StartConcept)
	if q.StartConcept == "" {
		return fmt.Errorf("emtpy start concept path, queryID: %s", q.QueryName)
	}
	if q.StartModifier != nil {
		*q.StartModifier.ModifierKey = strings.TrimSpace(*q.StartModifier.ModifierKey)
		if *q.StartModifier.ModifierKey == "" {
			return fmt.Errorf("empty start modifier key, queryID: %s, start concept: %s", q.QueryName, q.StartConcept)
		}
		*q.StartModifier.AppliedPath = strings.TrimSpace(*q.StartModifier.AppliedPath)
		if *q.StartModifier.AppliedPath == "" {
			return fmt.Errorf(
				"empty start modifier applied path, queryID: %s, start concept: %s, start modifier key: %s",
				q.QueryName, q.StartConcept,
				*q.StartModifier.ModifierKey,
			)
		}
	}

	q.EndConcept = strings.TrimSpace(q.EndConcept)
	if q.EndConcept == "" {
		return fmt.Errorf("empty end concept path, queryID: %s", q.QueryName)
	}
	if q.EndModifier != nil {
		*q.EndModifier.ModifierKey = strings.TrimSpace(*q.EndModifier.ModifierKey)
		if *q.EndModifier.ModifierKey == "" {
			return fmt.Errorf("empty end modifier key, queryID: %s, end concept: %s", q.QueryName, q.EndConcept)
		}
		*q.EndModifier.AppliedPath = strings.TrimSpace(*q.EndModifier.AppliedPath)
		if *q.EndModifier.AppliedPath == "" {
			return fmt.Errorf(
				"empty end modifier applied path, queryID: %s, end concept: %s, end modifier key: %s",
				q.QueryName, q.EndConcept,
				*q.EndModifier.ModifierKey,
			)
		}
	}

	q.UserID = strings.TrimSpace(q.UserID)
	if q.UserID == "" {
		return fmt.Errorf("empty user name, queryID: %s", q.QueryName)
	}
	q.TimeGranularity = strings.ToLower(strings.TrimSpace(q.TimeGranularity))
	if q.TimeGranularity == "" {
		return fmt.Errorf("empty granularityqueryID: %s", q.QueryName)
	}
	if _, isIn := granularityFunctions[q.TimeGranularity]; !isIn {
		granularities := make([]string, 0, len(granularityFunctions))
		for name := range granularityFunctions {
			granularities = append(granularities, name)
		}
		return fmt.Errorf("granularity %s not implemented, must be one of %v; queryID: %s", q.TimeGranularity, granularities, q.QueryName)
	}
	q.UserPublicKey = strings.TrimSpace(q.UserPublicKey)
	if q.UserPublicKey == "" {
		return fmt.Errorf("empty user public keyqueryID: %s", q.QueryName)
	}
	_, err := base64.URLEncoding.DecodeString(q.UserPublicKey)
	if err != nil {
		return fmt.Errorf("user public key is not valid against the alternate RFC4648 base64 for URL: %s; queryID: %s", err.Error(), q.QueryName)
	}
	return nil

}

// prepareArguments retrieves concept codes and patients that will be used as the arguments of direct SQL call
func prepareArguments(
	userID,
	cohortName, startConcept string,
	startModifier *survival_analysis.SurvivalAnalysisParamsBodyStartModifier,
	endConcept string,
	endModifier *survival_analysis.SurvivalAnalysisParamsBodyEndModifier,
) (
	startConceptCodes,
	startModifierCodes,
	endConceptCodes,
	endModifierCodes []string,
	cohort []int64, timers medcomodels.Timers,
	err error,
) {
	timers = make(map[string]time.Duration)
	// --- cohort patient list
	timer := time.Now()
	logrus.Info("Survival analysis: get patients")
	cohort, err = querytoolsserver.GetPatientList(userID, cohortName)

	if err != nil {
		logrus.Error("error while getting patient list")
		return
	}

	timers.AddTimers("medco-connector-get-patient-list", timer, nil)
	logrus.Info("Survival analysis: got patients")

	// --- get concept and modifier codes from the ontology
	logrus.Info("Survival analysis: get concept and modifier codes")
	err = utilserver.I2B2DBConnection.Ping()
	if err != nil {
		err = fmt.Errorf("while connecting to clear project database: %s", err.Error())
		return
	}
	startConceptCodes, err = getConceptCodes(startConcept)
	if err != nil {
		err = fmt.Errorf("while retrieving start concept code: %s", err.Error())
		return
	}
	if startModifier == nil {
		startModifierCodes = []string{"@"}
	} else {
		startModifierCodes, err = getModifierCodes(*startModifier.ModifierKey, *startModifier.AppliedPath)
	}
	if err != nil {
		err = fmt.Errorf("while retrieving start modifier code: %s", err.Error())
		return
	}
	endConceptCodes, err = getConceptCodes(endConcept)
	if err != nil {
		err = fmt.Errorf("while retrieving end concept code: %s", err.Error())
		return
	}
	if endModifier == nil {
		endModifierCodes = []string{"@"}
	} else {
		endModifierCodes, err = getModifierCodes(*endModifier.ModifierKey, *endModifier.AppliedPath)
	}
	if err != nil {
		err = fmt.Errorf("while retrieving end modifier code: %s", err.Error())
		return
	}
	logrus.Info("Survival analysis: got concept and modifier codes")
	return
}

// expansion takes a slice of SQLTimepoints and add encryption of zeros for events of interest and censoring events for each missing relative time from 0 to timeLimit.
// Relative times greater than timeLimit are discarded.
// Note that the time limit unit for this function is day.
func expansion(timePoints medcomodels.TimePoints, timeLimitDay int, granularity string) (medcomodels.TimePoints, error) {
	var timeLimit int64
	if granFunction, isIn := granularityFunctions[granularity]; isIn {
		timeLimit = granFunction(int64(timeLimitDay))
	} else {
		return nil, fmt.Errorf("granularity %s is not implemented", granularity)
	}

	res := make(medcomodels.TimePoints, timeLimit)
	availableTimePoints := make(map[int64]struct {
		EventsOfInterest int64
		CensoringEvents  int64
	}, len(timePoints))
	for _, timePoint := range timePoints {

		availableTimePoints[timePoint.Time] = timePoint.Events
	}
	for i := int64(0); i < timeLimit; i++ {
		if events, ok := availableTimePoints[int64(i)]; ok {
			res[i] = medcomodels.TimePoint{
				Time:   int64(i),
				Events: events,
			}
		} else {
			res[i] = medcomodels.TimePoint{
				Time: int64(i),
				Events: struct {
					EventsOfInterest int64
					CensoringEvents  int64
				}{
					EventsOfInterest: 0,
					CensoringEvents:  0,
				},
			}
		}

	}
	return res, nil
}

// getCode takes the full path of a I2B2 concept and returns its code
func getCode(path string) (string, error) {
	logrus.Debugf("Survival analysis: get code concept path %s", path)
	res, err := i2b2.GetOntologyConceptInfo(path)
	if err != nil {
		return "", err
	}
	if len(res) != 1 {
		return "", errors.Errorf("Result length of GetOntologyConceptInfo is expected to be 1. Got: %d", len(res))
	}

	if res[0].Code == "" {
		return "", errors.New("Code is empty")
	}
	logrus.Debugf("Survival analysis: got concept code %s", res[0].Code)

	return res[0].Code, nil

}

// getModifierPath takes the full path of a I2B2 modifier and its applied paht and returns its code
func getModifierCode(path string, appliedPath string) (string, error) {
	logrus.Debugf("Survival analysis: get modifier code modifier path %s applied path %s", path, appliedPath)
	res, err := i2b2.GetOntologyModifierInfo(path, appliedPath)
	if err != nil {
		return "", err
	}

	if len(res) != 1 {
		return "", errors.Errorf("Result length of GetOntologyTermInfo is expected to be 1. Got: %d. "+
			"Is applied path %s available for modifier key %s ?", len(res), appliedPath, path)
	}
	if res[0].Code == "" {
		return "", errors.New("Code is empty")
	}
	if res[0].AppliedPath != appliedPath {
		return "", fmt.Errorf("applied paths don't match. Is applied path %s available for modifier key %s ?", appliedPath, path)
	}
	logrus.Debugf("Survival analysis: got modifier code %s", res[0].Code)

	return res[0].Code, nil
}

// processGroupResult change resolution, expand and encrypt group result
func processGroupResult(errChan chan error, newEventGroup *EventGroup, sqlTimePoints medcomodels.TimePoints, timeLimitInDays int, timeGranularity string, index int) (timers medcomodels.Timers) {
	timers = make(map[string]time.Duration)

	// --- change time resolution
	timer := time.Now()
	sqlTimePoints, err := granularity(sqlTimePoints, timeGranularity)
	if err != nil {
		logrus.Error("Error while changing granularity")
		errChan <- err
		return
	}
	timers.AddTimers(fmt.Sprintf("medco-connector-change-timepoints-to-new-resolution%d", index), timer, nil)
	logrus.Debugf("Survival analysis: changed resolution for %s,  got %d timepoints", timeGranularity, len(sqlTimePoints))
	logrus.Tracef("Survival analysis: time points with resolution %s %+v", timeGranularity, sqlTimePoints)

	// --- expand
	timer = time.Now()
	sqlTimePoints, err = expansion(sqlTimePoints, timeLimitInDays, timeGranularity)
	if err != nil {
		err = fmt.Errorf("while expanding: %s", err.Error())
		errChan <- err
		return
	}
	timers.AddTimers(fmt.Sprintf("medco-connector-expansion %d", index), timer, nil)
	logrus.Debugf("Survival analysis: expanded to %d timepoints", len(sqlTimePoints))
	logrus.Tracef("Survival analysis: expanded time points %v", sqlTimePoints)

	// --- locally encrypt
	timer = time.Now()
	for _, sqlTimePoint := range sqlTimePoints {

		localEventEncryption, err := unlynx.EncryptWithCothorityKey(sqlTimePoint.Events.EventsOfInterest)
		if err != nil {
			err = fmt.Errorf("while encrypting event count: %s", err.Error())
			errChan <- err
			return
		}
		localCensoringEncryption, err := unlynx.EncryptWithCothorityKey(sqlTimePoint.Events.CensoringEvents)
		if err != nil {
			err = fmt.Errorf("while encrypting censoring count: %s", err.Error())
			errChan <- err
		}

		newEventGroup.TimePointResults = append(newEventGroup.TimePointResults, &TimePointResult{
			TimePoint: sqlTimePoint.Time,
			Result: Result{
				EventValueAgg:     localEventEncryption,
				CensoringValueAgg: localCensoringEncryption,
			}})

	}
	return timers

}

// fullCohort is called to build an explore definition when no subgroups are provided
func fullCohort(startConcept string) []*survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0 {
	newItems := make([]*models.PanelConceptItemsItems0, 1)
	encrypted := new(bool)
	not := new(bool)
	*encrypted = false
	*not = false
	term := new(string)
	*term = startConcept
	newItems[0] = &models.PanelConceptItemsItems0{
		Encrypted: encrypted,
		QueryTerm: term,
	}
	newPanels := make([]*models.Panel, 1)
	newPanels[0] = &models.Panel{
		ConceptItems: newItems,
		Not:          not,
		PanelTiming:  models.TimingAny,
	}

	return []*survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0{
		{
			GroupName:      "Full cohort",
			Panels:         newPanels,
			SubGroupTiming: models.TimingAny,
		},
	}
}

// timePointMapToList takes the relative-time-to-event-aggregates map and put its content in a list of TimePoint.
// A TimePoint structure contains the same information as the event aggregates plus a field referring to the relative time.
func timePointMapToList(timePointsMap map[int64]*medcomodels.Events) (list medcomodels.TimePoints) {
	list = make(medcomodels.TimePoints, 0, len(timePointsMap))
	for relativeTime, event := range timePointsMap {
		list = append(list, medcomodels.TimePoint{
			Time: relativeTime,
			Events: medcomodels.Events{
				EventsOfInterest: event.EventsOfInterest,
				CensoringEvents:  event.CensoringEvents,
			},
		})
	}

	sort.Slice(list, func(i, j int) bool {
		return list[i].Time < list[j].Time
	})

	return
}
