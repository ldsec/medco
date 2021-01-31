package handlers

import (
	"fmt"

	querytoolsserver "github.com/ldsec/medco/connector/server/querytools"
	survivalserver "github.com/ldsec/medco/connector/server/survivalanalysis"

	"github.com/go-openapi/runtime/middleware"
	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/ldsec/medco/connector/restapi/server/operations/survival_analysis"
	"github.com/sirupsen/logrus"
)

// MedCoSurvivalAnalysisHandler handles /survival-analysis API endpoint
func MedCoSurvivalAnalysisHandler(param survival_analysis.SurvivalAnalysisParams, principal *models.User) middleware.Responder {

	survivalAnalysisQuery := survivalserver.NewQuery(
		principal.ID,
		*param.Body.ID,
		*param.Body.UserPublicKey,
		*param.Body.CohortName,
		param.Body.SubGroupDefinitions,
		int(*param.Body.TimeLimit),
		*param.Body.TimeGranularity,
		*param.Body.StartConcept,
		param.Body.StartModifier,
		*param.Body.EndConcept,
		param.Body.EndModifier,
	)
	logrus.Debug("survivalAnalysis: ", survivalAnalysisQuery)
	err := survivalAnalysisQuery.Validate()
	if err != nil {
		logrus.Error(err)
		return survival_analysis.NewSurvivalAnalysisBadRequest().WithPayload(
			&survival_analysis.SurvivalAnalysisBadRequestBody{
				Message: "Survival query validation error:" + err.Error()})
	}
	found := false
	logrus.Info("checking cohort's existence")
	found, err = querytoolsserver.DoesCohortExist(principal.ID, *param.Body.CohortName)
	if err != nil {
		logrus.Error(err)
		return survival_analysis.NewSurvivalAnalysisDefault(500).WithPayload(
			&survival_analysis.SurvivalAnalysisDefaultBody{
				Message: "Survival query execution error:" + err.Error()})
	}

	if !found {
		logrus.Info("cohort not found")
		return survival_analysis.NewSurvivalAnalysisNotFound().WithPayload(
			&survival_analysis.SurvivalAnalysisNotFoundBody{
				Message: fmt.Sprintf("Cohort %s not found", *param.Body.CohortName),
			},
		)
	}
	logrus.Info("cohort found")

	err = survivalAnalysisQuery.Execute()

	if err != nil {
		err = fmt.Errorf("queryID: %s, error: %s", survivalAnalysisQuery.QueryName, err.Error())
		logrus.Error(err)
		return survival_analysis.NewSurvivalAnalysisDefault(500).WithPayload(
			&survival_analysis.SurvivalAnalysisDefaultBody{
				Message: "Survival query execution error:" + err.Error()})

	}
	results := survivalAnalysisQuery.Result

	resultList := make([]*survival_analysis.SurvivalAnalysisOKBodyResultsItems0, 0)
	for _, group := range survivalAnalysisQuery.Result.EncEvents {

		timePoints := make([]*survival_analysis.SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0, 0)
		for _, timePoint := range group.TimePointResults {
			timePoints = append(timePoints, &survival_analysis.SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0{Timepoint: int64(timePoint.TimePoint),
				Events: &survival_analysis.SurvivalAnalysisOKBodyResultsItems0GroupResultsItems0Events{
					Eventofinterest: timePoint.Result.EventValueAgg,
					Censoringevent:  timePoint.Result.CensoringValueAgg,
				}})
		}
		resultList = append(resultList, &survival_analysis.SurvivalAnalysisOKBodyResultsItems0{
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
