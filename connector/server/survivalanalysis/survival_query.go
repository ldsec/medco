package survivalserver

import (
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/ldsec/medco/connector/util"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	querytoolsserver "github.com/ldsec/medco/connector/server/querytools"
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
	StartModifier       string
	EndConcept          string
	EndModifier         string
	Result              *struct {
		Timers    util.Timers
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
	StartModifier string,
	EndConcept string,
	EndModifier string) *Query {
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
		EndConcept:          EndConcept,
		EndModifier:         EndModifier,
		Result: &struct {
			Timers    util.Timers
			EncEvents EventGroups
		}{}}
	res.Result.Timers = make(map[string]time.Duration)

	return res
}

// Execute runs the survival analysis query
func (q *Query) Execute() error {

	patientLists := make([][]int64, 0)
	initialCounts := make([]int64, 0)
	eventGroups := make(EventGroups, 0)
	timeLimitInDays := q.TimeLimit * granularityValues[q.TimeGranularity]
	timer := time.Now()

	startConceptCode, endConceptCode, cohort, timers, err := prepareArguments(q.UserID, q.CohortName, q.StartConcept, q.EndConcept)

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
		util.Timers
	}, len(definitions))
	errChan := make(chan error, len(definitions))
	signal := make(chan struct{})

	for i, definition := range definitions {
		channels[i] = make(chan struct {
			*EventGroup
			util.Timers
		}, 1)
		go func(i int, definition *survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0) {
			defer waitGroup.Done()
			timers := util.NewTimers()

			newEventGroup := &EventGroup{GroupID: definition.CohortName}
			panels := make([][]string, 0)
			not := make([]bool, 0)
			panels = append(panels, []string{q.StartConcept})
			not = append(not, false)
			for _, panel := range definition.Panels {
				terms := make([]string, 0)

				negation := *panel.Not

				for _, term := range panel.Items {
					terms = append(terms, *term.QueryTerm)
				}

				panels = append(panels, terms)
				not = append(not, negation)
			}

			timer = time.Now()
			logrus.Infof("I2B2 explore for subgroup %d", i)
			logrus.Tracef("panels %+v", panels)
			initialCount, patientList, err := SubGroupExplore(q.QueryName, i, panels, not)
			if err != nil {
				err = fmt.Errorf("during subgroup explore procedure: %s", err.Error())
				errChan <- err
				return
			}
			logrus.Infof("successful I2B2 explore query %d", i)
			timers.AddTimers(fmt.Sprintf("medco-connector-i2b2-query-group%d", i), timer, nil)
			patientList = intersect(cohort, patientList)
			patientLists = append(patientLists, patientList)
			initialCounts = append(initialCounts, initialCount)
			logrus.Debug("Initial Counts", initialCounts)

			timer = time.Now()
			initialCountEncrypt, err := unlynx.EncryptWithCothorityKey(initialCount)
			timers.AddTimers(fmt.Sprintf("medco-connector-encrypt-init-count-group%d", i), timer, nil)
			if err != nil {
				err = fmt.Errorf("while encrypting initial count: %s", err.Error())
				errChan <- err
				return
			}
			logrus.Debug("initialcount ", initialCountEncrypt)
			newEventGroup.EncInitialCount = initialCountEncrypt
			timer = time.Now()

			//  --- sql query on observation fact table
			sqlTimePoints, err := buildTimePoints(
				patientList,
				startConceptCode,
				q.StartModifier,
				endConceptCode,
				q.EndModifier,
				timeLimitInDays,
			)
			timers.AddTimers(fmt.Sprintf("medco-connector-build-timepoints%d", i), timer, nil)
			if err != nil {
				err = fmt.Errorf("error while getting building time points: %s", err.Error())
				errChan <- err
				return
			}
			logrus.Debugf("got %d time points", len(sqlTimePoints))
			logrus.Tracef("%+v", sqlTimePoints)
			processGroupResultTimers := processGroupResult(errChan, newEventGroup, sqlTimePoints, timeLimitInDays, q.TimeGranularity, i)
			q.Result.Timers.AddTimers("", timer, processGroupResultTimers)

			logrus.Tracef("Event groups %v", newEventGroup)
			timers.AddTimers(fmt.Sprintf("medco-connector-local-encryption%d", i), timer, nil)
			channels[i] <- struct {
				*EventGroup
				util.Timers
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
		logrus.Tracef("eventGroup %v", group)
	}
	timer = time.Now()
	var aksTimers util.Timers
	q.Result.EncEvents, aksTimers, err = AKSgroups(q.QueryName+"_AGG_AND_KEYSWITCH", eventGroups, q.UserPublicKey)
	q.Result.Timers.AddTimers("medco-connector-aggregate-and-key-switch", timer, aksTimers)
	if err != nil {
		err = fmt.Errorf("during aggregation and keyswitch: %s", err.Error())
	}
	return err
}

// Validate checks members of a Query instance for early error detection.
// Heading and trailing spaces are silently trimmed. Granularity string is silently written in lower case.
// If any other wrong member can be defaulted, a warning message is printed, otherwise an error is returned.
func (q *Query) Validate() error {
	q.StartConcept = strings.TrimSpace(q.StartConcept)
	if q.StartConcept == "" {
		return fmt.Errorf("emtpy start concept path")
	}
	q.StartModifier = strings.TrimSpace(q.StartModifier)
	if q.StartModifier == "" {
		logrus.Warn("empty start concept, defaulte to \"@\"")
	}
	q.QueryName = strings.TrimSpace(q.QueryName)
	if q.QueryName == "" {
		return fmt.Errorf("empty query name")
	}
	q.EndConcept = strings.TrimSpace(q.EndConcept)
	if q.EndConcept == "" {
		return fmt.Errorf("empty end concept path")
	}
	q.EndModifier = strings.TrimSpace(q.EndModifier)
	if q.EndModifier == "" {
		logrus.Warn("empty end modifier, defaulted to \"@\"")
	}
	q.UserID = strings.TrimSpace(q.UserID)
	if q.UserID == "" {
		return fmt.Errorf("empty user name")
	}
	q.TimeGranularity = strings.ToLower(strings.TrimSpace(q.TimeGranularity))
	if q.TimeGranularity == "" {
		return fmt.Errorf("empty granularity")
	}
	if _, isIn := granularityFunctions[q.TimeGranularity]; !isIn {
		granularities := make([]string, 0, len(granularityFunctions))
		for name := range granularityFunctions {
			granularities = append(granularities, name)
		}
		return fmt.Errorf("granularity %s not implemented, must be one of %v", q.TimeGranularity, granularities)
	}
	q.UserPublicKey = strings.TrimSpace(q.UserPublicKey)
	if q.UserPublicKey == "" {
		return fmt.Errorf("empty user public key")
	}
	_, err := base64.URLEncoding.DecodeString(q.UserPublicKey)
	if err != nil {
		return fmt.Errorf("user public key is not valid against the alternate RFC4648 base64 for URL: %s", err.Error())
	}
	return nil

}

// prepareArguments retrieves concept codes and patients that will be used as the arguments of direct SQL call
func prepareArguments(userID, cohortName, startConcept, endConcept string) (startConceptCode, endConceptCode string, cohort []int64, timers util.Timers, err error) {
	timers = make(map[string]time.Duration)
	// --- cohort patient list
	timer := time.Now()
	logrus.Info("get patients")
	cohort, err = querytoolsserver.GetPatientList(userID, cohortName)

	if err != nil {
		logrus.Error("error while getting patient list")
		return
	}

	timers.AddTimers("medco-connector-get-patient-list", timer, nil)
	logrus.Info("got patients")

	// --- get concept and modifier codes from the ontology
	logrus.Info("get concept and modifier codes")
	err = utilserver.I2B2DBConnection.Ping()
	if err != nil {
		err = fmt.Errorf("while connecting to clear project database: %s", err.Error())
		return
	}
	startConceptCode, err = getCode(startConcept)
	if err != nil {
		err = fmt.Errorf("while retrieving start concept code: %s", err.Error())
		return
	}
	endConceptCode, err = getCode(endConcept)
	if err != nil {
		err = fmt.Errorf("while retrieving end concept code: %s", err.Error())
		return
	}
	logrus.Info("got concept and modifier codes")
	return
}

// expansion takes a slice of SQLTimepoints and add encryption of zeros for events of interest and censoring events for each missing relative time from 0 to timeLimit.
// Relative times greater than timeLimit are discarded.
// Note that the time limit unit for this function is day.
func expansion(timePoints util.TimePoints, timeLimitDay int, granularity string) (util.TimePoints, error) {
	var timeLimit int
	if granFunction, isIn := granularityFunctions[granularity]; isIn {
		timeLimit = granFunction(timeLimitDay)
	} else {
		return nil, fmt.Errorf("granularity %s is not implemented", granularity)
	}

	res := make(util.TimePoints, timeLimit)
	availableTimePoints := make(map[int]struct {
		EventsOfInterest int64
		CensoringEvents  int64
	}, len(timePoints))
	for _, timePoint := range timePoints {

		availableTimePoints[timePoint.Time] = timePoint.Events
	}
	for i := 0; i < timeLimit; i++ {
		if events, ok := availableTimePoints[i]; ok {
			res[i] = util.TimePoint{
				Time:   i,
				Events: events,
			}
		} else {
			res[i] = util.TimePoint{
				Time: i,
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
	res, err := i2b2.GetOntologyTermInfo(path)
	if err != nil {
		return "", err
	}
	if len(res) != 1 {
		return "", errors.Errorf("Result length of GetOntologyTermInfo is expected to be 1. Got: %d", len(res))
	}

	if res[0].Code == "" {
		return "", errors.New("Code is empty")
	}

	return res[0].Code, nil

}

// processGroupResult change resolution, expand and encrypt group result
func processGroupResult(errChan chan error, newEventGroup *EventGroup, sqlTimePoints util.TimePoints, timeLimitInDays int, timeGranularity string, index int) (timers util.Timers) {
	timers = make(map[string]time.Duration)
	// --- expand
	timer := time.Now()
	sqlTimePoints, err := expansion(sqlTimePoints, timeLimitInDays, timeGranularity)
	if err != nil {
		err = fmt.Errorf("while expanding: %s", err.Error())
		errChan <- err
		return
	}
	timers.AddTimers(fmt.Sprintf("medco-connector-expansion %d", index), timer, nil)
	logrus.Debugf("expanded to %d timepoints", len(sqlTimePoints))
	logrus.Tracef("expanded time points %v", sqlTimePoints)

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
	newItems := make([]*models.PanelItemsItems0, 1)
	encrypted := new(bool)
	not := new(bool)
	*encrypted = false
	*not = false
	term := new(string)
	*term = startConcept
	newItems[0] = &models.PanelItemsItems0{
		Encrypted: encrypted,
		Operator:  "equals",
		QueryTerm: term,
	}
	newPanels := make([]*models.Panel, 1)
	newPanels[0] = &models.Panel{
		Items: newItems,
		Not:   not,
	}

	return []*survival_analysis.SurvivalAnalysisParamsBodySubGroupDefinitionsItems0{
		{
			CohortName: "Full cohort",
			Panels:     newPanels,
		},
	}
}
