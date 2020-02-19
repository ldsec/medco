package handlers

import (
	"fmt"

	"github.com/go-openapi/runtime/middleware"
	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/ldsec/medco-connector/restapi/server/operations/survival_analysis"
	"github.com/ldsec/medco-connector/survival"
	"github.com/sirupsen/logrus"
)

func MedCoSurvivalAnalysisGetSurvivalAnalysisHandler(param survival_analysis.GetSurvivalAnalysisParams, principal *models.User) middleware.Responder {

	survivalAnalysisQuery := survival.NewQuery("TODO find me an ID", param.UserPublicKeyAndPanels.UserPublicKey, param.UserPublicKeyAndPanels.Panels)
	patientErrChan := make(chan error)
	timeCodeErrChan := make(chan error)
	go func() {
		patientErrChan <- survivalAnalysisQuery.LoadPatients()
	}()
	go func() {
		timeCodeErrChan <- survivalAnalysisQuery.LoadTimeCodes(param.Granularity)
	}()

	if patientErr := <-patientErrChan; patientErr != nil {
		logrus.Error(fmt.Sprintf("Query execution error : %s", patientErr.Error()))
		return survival_analysis.NewGetSurvivalAnalysisNotFound()

	}
	if timeCodeErr := <-patientErrChan; timeCodeErr != nil {
		logrus.Error(fmt.Sprintf("Query execution error : %s", timeCodeErr.Error()))
		return survival_analysis.NewGetSurvivalAnalysisNotFound()

	}

	if err := survivalAnalysisQuery.Execute(); err != nil {
		logrus.Error(fmt.Sprintf("Query execution error : %s", err.Error()))
		return survival_analysis.NewGetSurvivalAnalysisDefault(500).WithPayload(&survival_analysis.GetSurvivalAnalysisDefaultBody{Message: fmt.Sprintf("Query execution error : %s", err.Error())})
	}
	results := survivalAnalysisQuery.Result
	if results == nil {
		logrus.Error("Unexpected nil results")
		return survival_analysis.NewGetSurvivalAnalysisDefault(500).WithPayload(&survival_analysis.GetSurvivalAnalysisDefaultBody{Message: "Query execution error : result pointer is nil"})
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
