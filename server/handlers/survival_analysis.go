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

	survivalAnalysisQuery := survivalserver.NewQuery(param.Body.ID, param.Body.UserPublicKey, param.Body.PatientSetID, param.Body.TimeCodes)

	batchNumber := 1 //solution for higher is not implemented yet
	if err := survivalAnalysisQuery.Execute(batchNumber); err != nil {
		logrus.Error(fmt.Sprintf("Query execution error : %s", err))
		return survival_analysis.NewGetSurvivalAnalysisDefault(500).WithPayload(&survival_analysis.GetSurvivalAnalysisDefaultBody{Message: err.Error()})
	}
	results := survivalAnalysisQuery.Result
	if results == nil {
		logrus.Panic("Unexpected nil results")

		//return survival_analysis.NewGetSurvivalAnalysisDefault(500).WithPayload(&survival_analysis.GetSurvivalAnalysisDefaultBody{Message: "Query execution error : result pointer is nil"})
	}
	var resultList []*survival_analysis.GetSurvivalAnalysisOKBodyItems0
	for key, val := range results.EncEvents {
		timePoint := key
		event := val[0]
		censoring := val[1]
		events := &survival_analysis.GetSurvivalAnalysisOKBodyItems0Events{Eventofinterest: event, Censoringevent: censoring}
		resultList = append(resultList, &survival_analysis.GetSurvivalAnalysisOKBodyItems0{Timepoint: timePoint, Events: events})
	}

	return survival_analysis.NewGetSurvivalAnalysisOK().WithPayload(resultList)

}
