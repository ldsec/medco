package handlers

import (
	"time"

	cohortsserver "github.com/ldsec/medco-connector/cohorts/server"

	"github.com/go-openapi/runtime/middleware"
	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/ldsec/medco-connector/restapi/server/operations/medco_node"
	medcoserver "github.com/ldsec/medco-connector/server"
	utilserver "github.com/ldsec/medco-connector/util/server"
	"github.com/ldsec/medco-connector/wrappers/i2b2"
)

// MedCoNodeExploreSearchHandler handles /medco/node/explore/search API endpoint
func MedCoNodeExploreSearchHandler(params medco_node.ExploreSearchParams, principal *models.User) middleware.Responder {

	searchResult, err := i2b2.GetOntologyChildren(*params.SearchRequest.Path)
	if err != nil {
		return medco_node.NewExploreSearchDefault(500).WithPayload(&medco_node.ExploreSearchDefaultBody{
			Message: err.Error(),
		})
	}

	return medco_node.NewExploreSearchOK().WithPayload(&medco_node.ExploreSearchOKBody{
		Search:  params.SearchRequest,
		Results: searchResult,
	})
}

// MedCoNodeExploreQueryHandler handles /medco/node/explore/query API endpoint
func MedCoNodeExploreQueryHandler(params medco_node.ExploreQueryParams, principal *models.User) middleware.Responder {

	// authorizations of query
	err := utilserver.AuthorizeExploreQueryType(principal, params.QueryRequest.Query.Type)
	if err != nil {
		return medco_node.NewExploreQueryDefault(403).WithPayload(&medco_node.ExploreQueryDefaultBody{
			Message: "Authorization of query failed: " + err.Error(),
		})
	}

	// create query
	query, err := medcoserver.NewExploreQuery(params.QueryRequest.ID, params.QueryRequest.Query)
	if err != nil {
		return medco_node.NewExploreQueryDefault(400).WithPayload(&medco_node.ExploreQueryDefaultBody{
			Message: "Bad query: " + err.Error(),
		})
	}

	// parse query type
	queryType, err := medcoserver.NewExploreQueryType(params.QueryRequest.Query.Type)
	if err != nil {
		return medco_node.NewExploreQueryDefault(400).WithPayload(&medco_node.ExploreQueryDefaultBody{
			Message: "Bad query type: " + err.Error(),
		})
	}

	// execute query
	err = query.Execute(queryType)
	if err != nil {
		return medco_node.NewExploreQueryDefault(500).WithPayload(&medco_node.ExploreQueryDefaultBody{
			Message: "Query execution error: " + err.Error(),
		})
	}

	// parse timers
	timers := make([]*models.ExploreQueryResultElementTimersItems0, 0)
	for timerName, timerDuration := range query.Result.Timers {
		milliseconds := int64(timerDuration / time.Millisecond)
		timers = append(timers, &models.ExploreQueryResultElementTimersItems0{
			Name:         timerName,
			Milliseconds: &milliseconds,
		})
	}

	return medco_node.NewExploreQueryOK().WithPayload(&medco_node.ExploreQueryOKBody{
		ID:    query.ID,
		Query: params.QueryRequest.Query,
		Result: &models.ExploreQueryResultElement{
			EncryptedCount:       query.Result.EncCount,
			EncryptedPatientList: query.Result.EncPatientList,
			PatientSetID:         float64(query.Result.PatientSetID),
			Timers:               timers,
			Status:               models.ExploreQueryResultElementStatusAvailable,
		}})
}

// MedCoNodeExploreQueryHandler handles /medco/node/explore/query API endpoint
func MedCoNodeGetCohortsHandler(params medco_node.GetCohortsParams, principal *models.User) middleware.Responder {
	userID := principal.ID
	cohorts, err := cohortsserver.GetCohorts(userID)
	if err != nil {
		medco_node.NewGetCohortsDefault(500).WithPayload(&medco_node.GetCohortsDefaultBody{
			Message: err.Error(),
		})
	}
	payload := &medco_node.GetCohortsOK{}
	for _, cohort := range cohorts {
		payload.Payload = append(payload.Payload,
			&medco_node.GetCohortsOKBodyItems0{
				CohortName:   cohort.CohortId,
				PatientSetID: float64(cohort.ResultInstanceID),
				CreationDate: float64(cohort.CreationDate),
				UpdateDate:   float64(cohort.UpdateDate),
			},
		)
	}

	return medco_node.NewGetCohortsOK().WithPayload(payload.Payload)

}

func MedCoNodePostCohortsHandler(params medco_node.PostCohortsParams, principal *models.User) middleware.Responder {
	cohort := params.Body.Cohort
	err := cohortsserver.InsertCohorts(principal.ID, cohort.CohortName, int(cohort.PatientSetID), int64(cohort.CreationDate), int64(cohort.UpdateDate))

	if err != nil {
		return medco_node.NewPostCohortsDefault(500).WithPayload(&medco_node.PostCohortsDefaultBody{
			Message: err.Error(),
		})
	}

	return medco_node.NewPostCohortsOK().WithPayload("cohort successfully updated")
}
