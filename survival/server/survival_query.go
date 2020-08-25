package survivalserver

import (
	"time"

	utilserver "github.com/ldsec/medco-connector/util/server"

	"github.com/ldsec/medco-connector/restapi/server/operations/survival_analysis"
)

// Query holds the ID of the survival analysis, its parameters and a pointer to its results
type Query struct {
	UserId              string
	UserPublicKey       string
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
		UserId:              UserId,
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

	cohort, err := GetPatientList(utilserver.DBConnection, int64(q.SetID), q.UserId)
	if err != nil {
		return err
	}

	if q.SubGroupDefinitions == nil || len(q.SubGroupDefinitions) == 0 {
		panels := [][]string{{q.StartConcept}}
		not := []bool{false}
		initCount, patientList, err := SubGroupExplore("", 0, panels, not)
		if err != nil {
			return err
		}
		initialCounts = append(initialCounts, initCount)
		patientLists = append(patientLists, Intersect(cohort, patientList))
	} else {
		for i, definition := range q.SubGroupDefinitions {
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
		}
	}
	return nil

}

func (q *Query) PrintTimers() {

}
