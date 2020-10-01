package survivalserver

import (
	"encoding/base64"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	querytoolsserver "github.com/ldsec/medco-connector/server/querytools"
	utilcommon "github.com/ldsec/medco-connector/util/common"
	utilserver "github.com/ldsec/medco-connector/util/server"
	"github.com/ldsec/medco-connector/wrappers/i2b2"
	"github.com/ldsec/medco-connector/wrappers/unlynx"

	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/ldsec/medco-connector/restapi/server/operations/survival_analysis"
)

// Query holds the ID of the survival analysis, its parameters and a pointer to its results
type Query struct {
	UserID              string
	UserPublicKey       string
	QueryName           string
	SetID               int
	SubGroupDefinitions []*survival_analysis.SubGroupDefinitionsItems0
	TimeLimit           int
	TimeGranularity     string
	StartConcept        string
	StartModifier       string
	EndConcept          string
	EndModifier         string
	Result              *struct {
		Timers    utilcommon.Timers
		EncEvents EventGroups
	}
}

// NewQuery query constructor
func NewQuery(UserID,
	QueryName,
	UserPublicKey string,
	SetID int,
	SubGroupDefinitions []*survival_analysis.SubGroupDefinitionsItems0,
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
		SetID:               SetID,
		SubGroupDefinitions: SubGroupDefinitions,
		TimeLimit:           TimeLimit,
		TimeGranularity:     TimeGranularity,
		StartConcept:        StartConcept,
		StartModifier:       StartModifier,
		EndConcept:          EndConcept,
		EndModifier:         EndModifier,
		Result: &struct {
			Timers    utilcommon.Timers
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

	// --- cohort patient list

	timer := time.Now()
	cohort, err := querytoolsserver.GetPatientList(utilserver.DBConnection, q.UserID, int64(q.SetID))
	q.Result.Timers.AddTimers("medco-connector-get-patient-list", timer, nil)
	logrus.Debug("got patients")

	if err != nil {
		logrus.Error("error while getting patient list")
		return err
	}

	// --- get concept and modifier codes from the ontology
	err = utilserver.I2B2DBConnection.Ping()

	if err != nil {
		logrus.Error("Unable to connect clear project database, ", err)
		return err
	}
	startConceptCode, err := getCode(q.StartConcept)
	if err != nil {
		logrus.Error("Error while retrieving concept code, ", err)
		return err
	}
	endConceptCode, err := getCode(q.EndConcept)
	if err != nil {
		logrus.Error("Error while retrieving concept code, ", err)
		return err
	}

	// --- build sub groups

	definitions := q.SubGroupDefinitions
	if q.SubGroupDefinitions == nil || len(q.SubGroupDefinitions) == 0 {
		newItems := make([]*models.PanelItemsItems0, 1)
		encrypted := new(bool)
		not := new(bool)
		*encrypted = false
		*not = false
		term := new(string)
		*term = q.StartConcept
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

		definitions = []*survival_analysis.SubGroupDefinitionsItems0{
			{
				CohortName: "FULL_COHORT",
				Panels:     newPanels,
			},
		}

	}
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(len(definitions))
	channels := make([]chan struct {
		*EventGroup
		utilcommon.Timers
	}, len(definitions))
	errChan := make(chan error, len(definitions))
	signal := make(chan struct{})

	for i, definition := range definitions {
		channels[i] = make(chan struct {
			*EventGroup
			utilcommon.Timers
		}, 1)
		go func(i int, definition *survival_analysis.SubGroupDefinitionsItems0) {
			defer waitGroup.Done()
			timers := utilcommon.NewTimers()

			newEventGroup := &EventGroup{GroupID: q.QueryName + fmt.Sprintf("_GROUP_%d", i)}
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
			initialCount, patientList, err := SubGroupExplore(q.QueryName, i, panels, not)
			timers.AddTimers(fmt.Sprintf("medco-connector-i2b2-query-group%d", i), timer, nil)
			patientList = Intersect(cohort, patientList)
			patientLists = append(patientLists, patientList)
			initialCounts = append(initialCounts, initialCount)
			logrus.Debug("Initial Counts", initialCounts)
			if err != nil {
				errChan <- err
				return
			}
			timer = time.Now()
			initialCountEncrypt, err := unlynx.EncryptWithCothorityKey(initialCount)
			timers.AddTimers(fmt.Sprintf("medco-connector-encrypt-init-count-group%d", i), timer, nil)
			if err != nil {
				errChan <- err
				return
			}
			logrus.Debug("initialcount ", initialCountEncrypt)
			newEventGroup.EncInitialCount = initialCountEncrypt
			timer = time.Now()

			//  --- sql query on observation fact table
			sqlTimePoints, err := buildTimePoints(utilserver.I2B2DBConnection,
				patientList,
				startConceptCode,
				q.StartModifier,
				endConceptCode,
				q.EndModifier,
				q.TimeLimit,
			)
			timers.AddTimers(fmt.Sprintf("medco-connector-build-timepoints%d", i), timer, nil)
			if err != nil {
				logrus.Error("error while getting building time points", timer)
				errChan <- err
				return
			}
			logrus.Debugf("got %d time points", len(sqlTimePoints))
			logrus.Tracef("%+v", sqlTimePoints)

			// --- change time resolution
			timer = time.Now()
			sqlTimePoints, err = granularity(sqlTimePoints, q.TimeGranularity)
			if err != nil {
				logrus.Error("Error while changing granularity")
				errChan <- err
				return
			}
			timers.AddTimers(fmt.Sprintf("medco-connector-change-timepoints-to-new-resolution%d", i), timer, nil)
			logrus.Debugf("changed resolution for %s,  got %d timepoints", q.TimeGranularity, len(sqlTimePoints))
			logrus.Tracef("time points with resolution %s %+v", q.TimeGranularity, sqlTimePoints)

			// --- expand
			timer = time.Now()
			sqlTimePoints, err = expansion(sqlTimePoints, q.TimeLimit, q.TimeGranularity)
			if err != nil {
				logrus.Error("Error while expanding")
				errChan <- err
				return
			}
			timers.AddTimers(fmt.Sprintf("medco-connector-expansion%d", i), timer, nil)
			logrus.Debugf("expanded to %d timepoints", len(sqlTimePoints))
			logrus.Tracef("expanded time points %v", sqlTimePoints)

			//locally encrypt
			timer = time.Now()
			for _, sqlTimePoint := range sqlTimePoints {

				localEventEncryption, err := unlynx.EncryptWithCothorityKey(sqlTimePoint.Events.EventsOfInterest)
				if err != nil {
					errChan <- err
					return
				}
				localCensoringEncryption, err := unlynx.EncryptWithCothorityKey(sqlTimePoint.Events.CensoringEvents)
				if err != nil {
					errChan <- err
				}

				newEventGroup.TimePointResults = append(newEventGroup.TimePointResults, &TimePointResult{
					TimePoint: sqlTimePoint.Time,
					Result: Result{
						EventValueAgg:     localEventEncryption,
						CensoringValueAgg: localCensoringEncryption,
					}})

			}
			logrus.Tracef("Eventt groups %v", newEventGroup)
			timers.AddTimers(fmt.Sprintf("medco-connector-local-encryption%d", i), timer, nil)
			channels[i] <- struct {
				*EventGroup
				utilcommon.Timers
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
	var aksTimers utilcommon.Timers
	q.Result.EncEvents, aksTimers, err = AKSgroups(q.QueryName+"_AGG_AND_KEYSWITCH", eventGroups, q.UserPublicKey)
	q.Result.Timers.AddTimers("medco-connector-aggregate-and-key-switch", timer, aksTimers)
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

// expansion takes a slice of SQLTimepoints and add encryption of zeros for events of interest and censoring events for each missing relative time from 0 to timeLimit.
// Relative times greater than timeLimit are discarded.
func expansion(timePoints utilcommon.TimePoints, timeLimitDay int, granularity string) (utilcommon.TimePoints, error) {
	var timeLimit int
	if granFunction, isIn := granularityFunctions[granularity]; isIn {
		timeLimit = granFunction(timeLimitDay)
	} else {
		return nil, fmt.Errorf("granularity %s is not implemented", granularity)
	}

	res := make(utilcommon.TimePoints, timeLimit)
	availableTimePoints := make(map[int]struct {
		EventsOfInterest int64
		CensoringEvents  int64
	}, len(timePoints))
	for _, timePoint := range timePoints {

		availableTimePoints[timePoint.Time] = timePoint.Events
	}
	for i := 0; i < timeLimit; i++ {
		if events, ok := availableTimePoints[i]; ok {
			res[i] = utilcommon.TimePoint{
				Time:   i,
				Events: events,
			}
		} else {
			res[i] = utilcommon.TimePoint{
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
