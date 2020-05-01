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
func MedCoSurvivalAnalysisGetSurvivalAnalysisHandler(param survival_analysis.GetSurvivalAnalysisParams, principal *models.User) middleware.Responder {

	survivalAnalysisQuery := survivalserver.NewQuery(param.Body.ID, param.Body.UserPublicKey, param.Body.PatientSetID, param.Body.PatientGroupIDs, param.Body.TimeCodes)

	if err := survivalAnalysisQuery.Execute(); err != nil {
		logrus.Error(fmt.Sprintf("Query execution error : %s", err.Error()))
		return survival_analysis.NewGetSurvivalAnalysisDefault(500).WithPayload(&survival_analysis.GetSurvivalAnalysisDefaultBody{Message: err.Error()})
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
	//logrus.Panicf("hohohoho : %d", len(resultList))

	requestResult := &survival_analysis.GetSurvivalAnalysisOKBody{Results: resultList, Timers: timers}

	return survival_analysis.NewGetSurvivalAnalysisOK().WithPayload(requestResult)

}
