package handlers

import (
	"fmt"

	survivalserver "github.com/ldsec/medco-connector/server/survivalanalysis"

	"github.com/go-openapi/runtime/middleware"
	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/ldsec/medco-connector/restapi/server/operations/survival_analysis"
	"github.com/sirupsen/logrus"
)

// MedCoSurvivalAnalysisGetSurvivalAnalysisHandler handles /survival-analysis API endpoint
func MedCoSurvivalAnalysisGetSurvivalAnalysisHandler(param survival_analysis.SurvivalAnalysisParams, principal *models.User) middleware.Responder {

	survivalAnalysisQuery := survivalserver.NewQuery(
		principal.ID,
		param.Body.ID,
		param.Body.UserPublicKey,
		int(param.Body.SetID),
		param.Body.SubGroupDefinitions,
		int(param.Body.TimeLimit),
		param.Body.TimeGranularity,
		param.Body.StartConcept,
		param.Body.StartModifier,
		param.Body.EndConcept,
		param.Body.EndModifier,
	)
	logrus.Debug("survivalAnalysis: ", survivalAnalysisQuery)
	err := survivalAnalysisQuery.Validate()
	if err != nil {
		err = fmt.Errorf("query validation error: %s", err.Error())
		logrus.Error(err)
		return survival_analysis.NewSurvivalAnalysisDefault(500).WithPayload(&survival_analysis.SurvivalAnalysisDefaultBody{Message: err.Error()})
	}
	err = survivalAnalysisQuery.Execute()

	if err != nil {
		err = fmt.Errorf("query execution error: %s", err.Error())
		logrus.Error(err)
		return survival_analysis.NewSurvivalAnalysisDefault(500).WithPayload(&survival_analysis.SurvivalAnalysisDefaultBody{Message: err.Error()})
	}
	results := survivalAnalysisQuery.Result

	resultList := make([]*survival_analysis.ResultsItems0, 0)
	for _, group := range survivalAnalysisQuery.Result.EncEvents {

		timePoints := make([]*survival_analysis.ResultsItems0GroupResultsItems0, 0)
		for _, timePoint := range group.TimePointResults {
			timePoints = append(timePoints, &survival_analysis.ResultsItems0GroupResultsItems0{Timepoint: int64(timePoint.TimePoint),
				Events: &survival_analysis.ResultsItems0GroupResultsItems0Events{
					Eventofinterest: timePoint.Result.EventValueAgg,
					Censoringevent:  timePoint.Result.CensoringValueAgg,
				}})
		}
		resultList = append(resultList, &survival_analysis.ResultsItems0{
			GroupID:      group.GroupID,
			InitialCount: group.EncInitialCount,
			GroupResults: timePoints,
		})
	}

	//parse timers
	modelsTimers := results.Timers.TimersToAPIModel()

	requestResult := &survival_analysis.SurvivalAnalysisOKBody{Results: resultList, Timers: modelsTimers}

	return survival_analysis.NewSurvivalAnalysisOK().WithPayload(requestResult)

}
