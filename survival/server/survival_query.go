package survivalserver

import (
	"fmt"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	querytools "github.com/ldsec/medco-connector/queryTools"
	"github.com/ldsec/medco-connector/wrappers/unlynx"

	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/ldsec/medco-connector/restapi/server/operations/survival_analysis"
)

// Query holds the ID of the survival analysis, its parameters and a pointer to its results
type Query struct {
	UserId              string
	UserPublicKey       string
	QueryName           string
	SetID               int
	SubGroupDefinitions []*survival_analysis.SubGroupDefinitionsItems0
	TimeLimit           int
	TimeGranularity     string
	StartConcept        string
	StartColumn         string
	StartModifier       string
	EndConcept          string
	EndColumn           string
	EndModifier         string
	Result              *struct {
		Timers    map[string]time.Duration
		EncEvents EventGroups
		timerLock *sync.Mutex
	}
}

// NewQuery query constructor
func NewQuery(UserId,
	QueryName,
	UserPublicKey string,
	SetID int,
	SubGroupDefinitions []*survival_analysis.SubGroupDefinitionsItems0,
	TimeLimit int,
	TimeGranularity string,
	StartConcept string,
	StartColumn string,
	StartModifier string,
	EndConcept string,
	EndColumn string,
	EndModifier string) *Query {
	res := &Query{
		UserPublicKey:       UserPublicKey,
		UserId:              UserId,
		QueryName:           QueryName,
		SetID:               SetID,
		SubGroupDefinitions: SubGroupDefinitions,
		TimeLimit:           TimeLimit,
		TimeGranularity:     TimeGranularity,
		StartConcept:        StartConcept,
		StartColumn:         StartColumn,
		StartModifier:       StartModifier,
		EndConcept:          EndConcept,
		EndColumn:           EndColumn,
		EndModifier:         EndModifier,
		Result: &struct {
			Timers    map[string]time.Duration
			EncEvents EventGroups
			timerLock *sync.Mutex
		}{}}
	res.Result.Timers = make(map[string]time.Duration)
	res.Result.timerLock = &sync.Mutex{}
	return res
}

func (q *Query) Execute() error {

	patientLists := make([][]int64, 0)
	initialCounts := make([]int64, 0)
	eventGroups := make(EventGroups, 0)

	//build subgroups

	timer := time.Now()
	cohort, err := GetPatientList(querytools.ConnectorDB, int64(q.SetID), q.UserId)
	q.addTimers("medco-connector-get-patient-list", timer)
	logrus.Debug("got patients")

	if err != nil {
		logrus.Error("error while getting patient list")
		return err
	}

	//get concept and modifier codes from the ontology
	err = DirectI2B2.Ping()

	if err != nil {
		logrus.Error("Unable to connect clear project database, ", err)
		return err
	}
	startConceptCode, err := GetCode(q.StartConcept)
	if err != nil {
		logrus.Error("Error while retrieving concept code, ", err)
		return err
	}

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
	channels := make([]chan *EventGroup, len(definitions))
	errChan := make(chan error, len(definitions))
	signal := make(chan struct{})

	for i, definition := range definitions {
		channels[i] = make(chan *EventGroup, 1)
		go func(i int, definition *survival_analysis.SubGroupDefinitionsItems0) {
			defer waitGroup.Done()

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
			q.addTimers(fmt.Sprintf("medco-connector-i2b2-query-group%d", i), timer)
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
			q.addTimers(fmt.Sprintf("medco-connector-encrypt-init-count-group%d", i), timer)
			if err != nil {
				errChan <- err
				return
			}
			logrus.Debug("initialcount ", initialCountEncrypt)
			newEventGroup.EncInitialCount = initialCountEncrypt
			timer = time.Now()
			sqlTimePoints, err := BuildTimePoints(DirectI2B2,
				patientList,
				startConceptCode,
				q.StartColumn,
				q.StartModifier,
				q.EndConcept,
				q.EndColumn,
				q.EndModifier,
			)
			q.addTimers(fmt.Sprintf("medco-connector-build-timepoints%d", i), timer)
			if err != nil {
				logrus.Error("error while getting building time points", timer)
				errChan <- err
				return
			}
			logrus.Debugf("got %d time points", len(sqlTimePoints))
			//locally encrypt
			timePointsSet := make(map[int](struct{}))
			for _, sqlTimePoint := range sqlTimePoints {
				if sqlTimePoint.timePoint <= q.TimeLimit {
					timePointsSet[sqlTimePoint.timePoint] = struct{}{}
					localEventEncryption, err := unlynx.EncryptWithCothorityKey(int64(sqlTimePoint.localEventAggregate))
					if err != nil {
						errChan <- err
						return
					}
					localCensoringEncryption, err := unlynx.EncryptWithCothorityKey(int64(sqlTimePoint.localCensoringAggrete))
					if err != nil {
						errChan <- err
					}

					newEventGroup.TimePointResults = append(newEventGroup.TimePointResults, &TimePointResult{
						TimePoint: sqlTimePoint.timePoint,
						Result: Result{
							EventValueAgg:     localEventEncryption,
							CensoringValueAgg: localCensoringEncryption,
						}})
				}
			}

			//get timepoints count for events of interest and censoring events

			//TODO: pad for zero time
			timer = time.Now()

			//TODO this put a full vector from 0 to time limit, with a lot of zero encrypted points
			for j := 0; j < q.TimeLimit; j++ {
				if _, isIn := timePointsSet[j]; !isIn {

					zeroEncrypt, err := unlynx.EncryptWithCothorityKey(int64(0))
					if err != nil {
						errChan <- err
						return
					}
					zeroEncrypt1, err := unlynx.EncryptWithCothorityKey(int64(0))
					if err != nil {
						errChan <- err
						return
					}
					newEventGroup.TimePointResults = append(newEventGroup.TimePointResults, &TimePointResult{
						TimePoint: j,
						Result: Result{
							EventValueAgg:     zeroEncrypt,
							CensoringValueAgg: zeroEncrypt1,
						}})
				}
			}
			q.addTimers(fmt.Sprintf("encrypt-zero-events%d", i), timer)
			channels[i] <- newEventGroup
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

		eventGroups = append(eventGroups, <-channel)
	}

	//TODO this put a full vector from 0 to time limit, with a lot of zero encrypted points

	//Key Switch !!

	for _, group := range eventGroups {
		logrus.Trace("eventGroup", *group)
	}
	timer = time.Now()
	q.Result.EncEvents, _, err = AKSgroups(q.QueryName+"_AGG_AND_KEYSWITCH", eventGroups, q.UserPublicKey)
	q.addTimers("medco-connector-aggregate-and-key-switch", timer)
	return err

}

// addTimers adds timers to the query results
func (q *Query) addTimers(timerName string, since time.Time) {
	if timerName != "" {
		q.Result.timerLock.Lock()
		q.Result.Timers[timerName] = time.Since(since)
		q.Result.timerLock.Unlock()
	}
}

func (q *Query) PrintTimers() {
	logrus.Debug("timer, duration:")
	q.Result.timerLock.Lock()
	for timerName, duration := range q.Result.Timers {
		logrus.Debug(timerName, " , ", duration.Milliseconds())
	}
	q.Result.timerLock.Unlock()

}
