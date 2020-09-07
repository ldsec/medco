package survivalserver

import (
	"fmt"
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
		}{}}
	res.Result.Timers = make(map[string]time.Duration)
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
	/*
		if q.SubGroupDefinitions == nil || len(q.SubGroupDefinitions) == 0 {
			eventGroups = append(eventGroups, &EventGroup{GroupID: q.QueryName + "_FULL_COHORT"})
			panels := [][]string{{q.StartConcept}}
			logrus.Debug(q.StartConcept, panels[0][0])
			not := []bool{false}
			timer = time.Now()
			initCount, patientList, err := SubGroupExplore(q.QueryName, 0, panels, not)
			q.addTimers("medco-connector-i2b2-query-group0", timer)
			if err != nil {
				return err
			}
			initialCounts = append(initialCounts, initCount)
			timer = time.Now()
			initialCountEncrypt, err := unlynx.EncryptWithCothorityKey(initCount)
			if err != nil {
				return err
			}
			q.addTimers("medco-connector-encrypt-init-count-group0", timer)
			logrus.Debug("initialcount ", initialCountEncrypt)
			eventGroups[0].EncInitialCount = initialCountEncrypt
			patientLists = append(patientLists, Intersect(cohort, patientList))
		}
	*/

	for i, definition := range definitions {
		eventGroups = append(eventGroups, &EventGroup{GroupID: q.QueryName + fmt.Sprintf("_GROUP_%d", i)})
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
		patientLists = append(patientLists, Intersect(cohort, patientList))
		initialCounts = append(initialCounts, initialCount)
		logrus.Debug("Initial Counts", initialCounts)
		if err != nil {
			return err
		}
		timer = time.Now()
		initialCountEncrypt, err := unlynx.EncryptWithCothorityKey(initialCount)
		q.addTimers(fmt.Sprintf("medco-connector-encrypt-init-count-group%d", i), timer)
		if err != nil {
			return err
		}
		logrus.Debug("initialcount ", initialCountEncrypt)
		eventGroups[i].EncInitialCount = initialCountEncrypt
	}

	//get timepoints count for events of interest and censoring events

	//TODO: pad for zero time
	logrus.Debug("patient lists", len(patientLists))
	timer = time.Now()
	for i, patientList := range patientLists {
		logrus.Debug(i)
		logrus.Debug(patientList)
		logrus.Debug(startConceptCode)
		logrus.Debug(q.StartModifier)
		logrus.Debug(q.EndConcept)
		logrus.Debug(q.EndColumn)
		logrus.Debug(q.EndModifier)

		sqlTimePoints, err := BuildTimePoints(DirectI2B2,
			patientList,
			startConceptCode,
			q.StartColumn,
			q.StartModifier,
			q.EndConcept,
			q.EndColumn,
			q.EndModifier,
		)
		if err != nil {
			logrus.Error("error while getting building time points", timer)
			return err
		}
		logrus.Debugf("got %d time points", len(sqlTimePoints))
		//locally encrypt
		timePointsSet := make(map[int](struct{}))
		for _, sqlTimePoint := range sqlTimePoints {
			if sqlTimePoint.timePoint <= q.TimeLimit {
				timePointsSet[sqlTimePoint.timePoint] = struct{}{}
				localEventEncryption, err := unlynx.EncryptWithCothorityKey(int64(sqlTimePoint.localEventAggregate))
				if err != nil {
					return err
				}
				localCensoringEncryption, err := unlynx.EncryptWithCothorityKey(int64(sqlTimePoint.localCensoringAggrete))
				if err != nil {
					return err
				}

				eventGroups[i].TimePointResults = append(eventGroups[i].TimePointResults, &TimePointResult{
					TimePoint: sqlTimePoint.timePoint,
					Result: Result{
						EventValueAgg:     localEventEncryption,
						CensoringValueAgg: localCensoringEncryption,
					}})
			}
		}

		//TODO this put a full vector from 0 to time limit, with a lot of zero encrypted points
		for j := 0; j < q.TimeLimit; j++ {
			if _, isIn := timePointsSet[j]; !isIn {

				zeroEncrypt, err := unlynx.EncryptWithCothorityKey(int64(0))
				if err != nil {
					return err
				}
				zeroEncrypt1, err := unlynx.EncryptWithCothorityKey(int64(0))
				if err != nil {
					return err
				}
				eventGroups[i].TimePointResults = append(eventGroups[i].TimePointResults, &TimePointResult{
					TimePoint: j,
					Result: Result{
						EventValueAgg:     zeroEncrypt,
						CensoringValueAgg: zeroEncrypt1,
					}})
			}
		}

	}
	q.addTimers("medco-connector-timepoint-queries", timer)

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
		q.Result.Timers[timerName] = time.Since(since)
	}
}

func (q *Query) PrintTimers() {
	logrus.Debug("timer, duration:")
	for timerName, duration := range q.Result.Timers {
		logrus.Debug(timerName, " , ", duration.Milliseconds())
	}

}
