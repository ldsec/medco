package survivalserver

import (
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	querytools "github.com/ldsec/medco-connector/queryTools"
	"github.com/ldsec/medco-connector/wrappers/unlynx"

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
	return &Query{
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
		EndModifier:         EndModifier}
}

func (q *Query) Execute() error {

	patientLists := make([][]int64, 0)
	initialCounts := make([]int64, 0)
	eventGroups := make(EventGroups, 0)

	//build subgroups

	cohort, err := GetPatientList(querytools.ConnectorDB, int64(q.SetID), q.UserId)

	if err != nil {
		logrus.Error("error while getting patient list")
		return err
	}

	if q.SubGroupDefinitions == nil || len(q.SubGroupDefinitions) == 0 {
		eventGroups = append(eventGroups, &EventGroup{GroupID: q.QueryName + "_FULL_COHORT"})
		panels := [][]string{{q.StartConcept}}
		logrus.Debug(q.StartConcept, panels[0][0])
		not := []bool{false}
		initCount, patientList, err := SubGroupExplore(q.QueryName, 0, panels, not)
		if err != nil {
			return err
		}
		initialCounts = append(initialCounts, initCount)
		initialCountEncrypt, err := unlynx.EncryptWithCothorityKey(initCount)
		if err != nil {
			return err
		}
		logrus.Debug("initialcount ", initialCountEncrypt)
		eventGroups[0].EncInitialCount = initialCountEncrypt
		patientLists = append(patientLists, Intersect(cohort, patientList))
	} else {
		for i, definition := range q.SubGroupDefinitions {
			eventGroups = append(eventGroups, &EventGroup{GroupID: q.QueryName + fmt.Sprintf("_GROUP_%d", i)})
			panels := make([][]string, 0)
			not := make([]bool, 0)
			for _, panel := range definition.Panels {
				terms := make([]string, 0)

				negation := *panel.Not
				for _, term := range panel.Items {
					terms = append(terms, *term.QueryTerm)
				}

				panels = append(panels, terms)
				not = append(not, negation)
			}

			initialCount, patientList, err := SubGroupExplore("", i, panels, not)
			patientLists = append(patientLists, Intersect(cohort, patientList))
			initialCounts = append(initialCounts, initialCount)
			if err != nil {
				return err
			}

			initialCountEncrypt, err := unlynx.EncryptWithCothorityKey(initialCount)
			if err != nil {
				return err
			}
			logrus.Debug("initialcount ", initialCountEncrypt)
			eventGroups[i].EncInitialCount = initialCountEncrypt
		}
	}

	//get concept and modifier codes from the ontology
	err = DirectI2B2.Ping()

	if err != nil {
		logrus.Error("Unable to connect clear project database, ", err)
		return err
	}
	startConceptCode, err := GetCode(DirectI2B2, q.StartConcept)
	if err != nil {
		logrus.Error("Error while retrieving concept code, ", err)
		return err
	}
	//get timepoints count for events of interest and censoring events

	//TODO: pad for zero time
	logrus.Debug("patient lists", len(patientLists))
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
			logrus.Error("error while getting building time points")
			return err
		}
		logrus.Debugf("got %d time points", len(sqlTimePoints))
		//locally encrypt
		for _, sqlTimePoint := range sqlTimePoints {
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

	//Key Switch !!
	q.Result = &struct {
		Timers    map[string]time.Duration
		EncEvents EventGroups
	}{}
	for _, group := range eventGroups {
		logrus.Trace("eventGroup", *group)
	}
	q.Result.EncEvents, _, err = AKSgroups(q.QueryName+"_AGG_AND_KEYSWITCH", eventGroups, q.UserPublicKey)

	return err

}

func (q *Query) PrintTimers() {

}
