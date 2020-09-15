package handlers

import (
	"fmt"
	"time"

	querytools "github.com/ldsec/medco-connector/queryTools"

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
	timers := query.Result.Timers.TimersToAPIModel()

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

// MedCoNodeGetCohortsHandler handles GET /medco/node/explore/cohorts  API endpoint
func MedCoNodeGetCohortsHandler(params medco_node.GetCohortsParams, principal *models.User) middleware.Responder {
	userID := principal.ID
	cohorts, err := querytools.GetSavedCohorts(querytools.ConnectorDB, userID)
	if err != nil {
		medco_node.NewGetCohortsDefault(500).WithPayload(&medco_node.GetCohortsDefaultBody{
			Message: err.Error(),
		})
	}
	if len(cohorts) == 0 {
		return medco_node.NewGetCohortsNotFound()
	}
	payload := &medco_node.GetCohortsOK{}
	for _, cohort := range cohorts {
		payload.Payload = append(payload.Payload,
			&medco_node.GetCohortsOKBodyItems0{
				CohortName:   cohort.CohortName,
				CohortID:     float64(cohort.CohortID),
				QueryID:      float64(cohort.QueryID),
				CreationDate: cohort.CreationDate.Format(time.RFC3339),
				UpdateDate:   cohort.UpdateDate.Format(time.RFC3339),
			},
		)
	}

	return medco_node.NewGetCohortsOK().WithPayload(payload.Payload)

}

// MedCoNodePostCohortsHandler handles POST /medco/node/explore/cohorts  API endpoint
func MedCoNodePostCohortsHandler(params medco_node.PostCohortsParams, principal *models.User) middleware.Responder {

	cohort := params.Body

	creationDate, err := time.Parse(time.RFC3339, cohort.CreationDate)
	if err != nil {
		return medco_node.NewPostCohortsDefault(500).WithPayload(&medco_node.PostCohortsDefaultBody{
			Message: fmt.Sprintf("String %s is not a date with RF3339 layout", cohort.CreationDate),
		})
	}
	updateDate, err := time.Parse(time.RFC3339, cohort.UpdateDate)
	if err != nil {
		return medco_node.NewPostCohortsDefault(500).WithPayload(&medco_node.PostCohortsDefaultBody{
			Message: fmt.Sprintf("String %s is not a date with RF3339 layout", cohort.UpdateDate),
		})
	}
	cohorts, err := querytools.GetSavedCohorts(querytools.ConnectorDB, principal.ID)
	if err != nil {
		return medco_node.NewPostCohortsDefault(500).WithPayload(&medco_node.PostCohortsDefaultBody{
			Message: err.Error(),
		})
	}
	for _, existingCohort := range cohorts {
		if existingCohort.CohortName == cohort.CohortName {
			lastUpdate, _ := time.Parse(time.RFC3339, cohort.UpdateDate)
			if lastUpdate.After(updateDate) {
				return medco_node.NewPostCohortsInternalServerError()
			}
			break
		}
	}
	querytools.InsertCohort(querytools.ConnectorDB, principal.ID, int(cohort.PatientSetID), cohort.CohortName, creationDate, updateDate)

	return medco_node.NewPostCohortsOK().WithPayload("cohorts successfully updated")
}
