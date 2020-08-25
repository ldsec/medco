package survivalserver

import (
	"time"

	"github.com/ldsec/medco-connector/restapi/server/operations/survival_analysis"

	querytools "github.com/ldsec/medco-connector/queryTools"
	"github.com/sirupsen/logrus"
)

// Query holds the ID of the survival analysis, its parameters and a pointer to its results
type Query struct {
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
func NewQuery(UserPublicKey string,
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
	return &Query{SetID: SetID,
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

	var groupPatientSetIds []int

	if len(q.SubGroupDefinitions) > 0 {
		/*
			the i2b2 explore queries will go here:
			1)retrieve patient set coll ids from i2b2 qt tables (not from  medco connector)
			2)execute the query, add the coll ids to the clear text query
			3) redo the instersection if only crypto queries occured

		*/

	} else {
		groupPatientSetIds = append(groupPatientSetIds, q.SetID)
	}
	timePointsMap := make(map[int][]SqlTimePoint)
	for _, groupID := range groupPatientSetIds {
		patientNumbers, qtErr := querytools.GetPatientList(groupID)
		if qtErr != nil {
			return qtErr
		}
		if len(patientNumbers) == 0 {
			logrus.Debugf("Result instance ID %d is empty", groupID)
			timePointsMap[groupID] = make([]SqlTimePoint, 0)
			continue
		}

		timePoints, tpError := BuildTimePoints(patientNumbers, q.StartConcept, q.StartColumn, q.EndConcept, q.EndColumn)
		if tpError != nil {
			return tpError
		}
		timePointsMap[groupID] = timePoints
	}
	/*
		4) add zeros
	*/
	return nil

}

func (q *Query) PrintTimers() {

}
