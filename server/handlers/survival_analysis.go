package handlers

import (
	"fmt"

	survivalserver "github.com/ldsec/medco-connector/survival/server"

	"github.com/go-openapi/runtime/middleware"
	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/ldsec/medco-connector/restapi/server/operations/survival_analysis"
	"github.com/sirupsen/logrus"
)

// MedCoSurvivalAnalysisGetSurvivalAnalysisHandler handles /survival-analysis API endpoint
func MedCoSurvivalAnalysisGetSurvivalAnalysisHandler(param survival_analysis.SurvivalAnalysisParams, principal *models.User) middleware.Responder {

	survivalAnalysisQuery := survivalserver.NewQuery(param.Body.UserPublicKey, int(param.Body.SetID), param.Body.SubGroupDefinitions, int(param.Body.TimeLimit), param.Body.TimeGranularity, param.Body.StartConcept, param.Body.StartColumn, param.Body.EndConcept, param.Body.EndColumn)

	if err := survivalAnalysisQuery.Execute(); err != nil {
		logrus.Error(fmt.Sprintf("Query execution error : %s", err.Error()))
		return survival_analysis.NewSurvivalAnalysisDefault(500).WithPayload(&survival_analysis.SurvivalAnalysisDefaultBody{Message: err.Error()})
	}
	survivalAnalysisQuery.PrintTimers()
	results := survivalAnalysisQuery.Result

	timers := make(map[string]float64, len(results.Timers))
	for timerKey, timerValue := range results.Timers {
		if _, exists := timers[timerKey]; exists {
			logrus.Warn("timer for " + timerKey + " already exists, previous value will be lost")
		}
		timers[timerKey] = timerValue.Seconds()
	}
	resultList := make([]*survival_analysis.ResultsItems0, 0)
	logrus.Debugf("Shiba %v: ", survivalAnalysisQuery.Result.EncEvents)
	for _, group := range survivalAnalysisQuery.Result.EncEvents {
		logrus.Debugf("Inu %v: ", group)

		timePoints := make([]*survival_analysis.ResultsItems0GroupResultsItems0, 0)
		for _, timePoint := range group.TimePointResults {
			timePoints = append(timePoints, &survival_analysis.ResultsItems0GroupResultsItems0{Timepoint: timePoint.TimePoint,
				Events: &survival_analysis.ResultsItems0GroupResultsItems0Events{
					Eventofinterest: timePoint.Result.EventValueAgg,
					Censoringevent:  timePoint.Result.CensoringValueAgg,
				}})
		}
		resultList = append(resultList, &survival_analysis.ResultsItems0{
			GroupID:      group.GroupID,
			GroupResults: timePoints,
		})
	}
	requestResult := &survival_analysis.SurvivalAnalysisOKBody{Results: resultList, Timers: timers}

	return survival_analysis.NewSurvivalAnalysisOK().WithPayload(requestResult)

}
