package survivalserver

import (
	"fmt"
	"time"

	utilserver "github.com/ldsec/medco-connector/util/server"

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
		return err
	}

	if q.SubGroupDefinitions == nil || len(q.SubGroupDefinitions) == 0 {
		eventGroups = append(eventGroups, &EventGroup{GroupID: q.QueryName + "_FULL_COHORT"})
		panels := [][]string{{q.StartConcept}}
		not := []bool{false}
		initCount, patientList, err := SubGroupExplore("", 0, panels, not)
		if err != nil {
			return err
		}
		initialCounts = append(initialCounts, initCount)
		initialCountEncrypt, err := unlynx.EncryptWithCothorityKey(initCount)
		if err != nil {
			return err
		}
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
			eventGroups[i].EncInitialCount = initialCountEncrypt
		}
	}

	//get timepoints count for events of interest and censoring events

	//TODO: pad for zero time

	for i, patientList := range patientLists {

		sqlTimePoints, err := BuildTimePoints(utilserver.DBConnection,
			patientList,
			q.StartConcept,
			q.StartColumn,
			q.StartModifier,
			q.EndConcept,
			q.EndColumn,
			q.EndModifier,
		)
		if err != nil {
			return err
		}
		timePointResults := eventGroups[i].TimePointResults
		//locally encrypt
		for _, sqlTimePoint := range sqlTimePoints {
			localEventEncryption, err := unlynx.EncryptWithCothorityKey(int64(sqlTimePoint.localEventAggregate))
			if err != nil {
				return err
			}
			localCensoringEncryption, err := unlynx.EncryptWithCothorityKey(int64(sqlTimePoint.localEventAggregate))
			if err != nil {
				return err
			}

			timePointResults = append(timePointResults, &TimePointResult{
				TimePoint: sqlTimePoint.timePoint,
				Result: Result{
					EventValueAgg:     localEventEncryption,
					CensoringValueAgg: localCensoringEncryption,
				}})
		}

	}

	//Key Switch !!

	q.Result.EncEvents, _, err = AKSgroups(q.QueryName+"_AGG_AND_KEYSWITCH", eventGroups, q.UserPublicKey)

	return err

}

func (q *Query) PrintTimers() {

}
