package handlers

import (
	"fmt"
	"time"

	medcoserver "github.com/ldsec/medco-connector/server/explore"
	querytoolsserver "github.com/ldsec/medco-connector/server/querytools"
	"github.com/sirupsen/logrus"

	"github.com/go-openapi/runtime/middleware"

	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/ldsec/medco-connector/restapi/server/operations/medco_node"

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
	query, err := medcoserver.NewExploreQuery(params.QueryRequest.ID, params.QueryRequest.Query, principal.ID)
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
			PatientSetID:         int64(query.Result.PatientSetID),
			Timers:               timers,
			Status:               models.ExploreQueryResultElementStatusAvailable,
		}})
}

// MedCoNodeGetCohortsHandler handles GET /medco/node/explore/cohorts  API endpoint
func MedCoNodeGetCohortsHandler(params medco_node.GetCohortsParams, principal *models.User) middleware.Responder {
	userID := principal.ID
	cohorts, err := querytoolsserver.GetSavedCohorts(utilserver.DBConnection, userID)
	if err != nil {
		medco_node.NewGetCohortsDefault(500).WithPayload(&medco_node.GetCohortsDefaultBody{
			Message: "Get cohort execution error: " + err.Error(),
		})
	}
	payload := &medco_node.GetCohortsOK{}
	for _, cohort := range cohorts {
		payload.Payload = append(payload.Payload,
			&medco_node.GetCohortsOKBodyItems0{
				CohortName:   cohort.CohortName,
				CohortID:     int64(cohort.CohortID),
				QueryID:      int64(cohort.QueryID),
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

	hasID, err := querytoolsserver.CheckQueryID(utilserver.DBConnection, principal.ID, int(cohort.PatientSetID))
	if err != nil {
		return medco_node.NewPostCohortsDefault(500).WithPayload(&medco_node.PostCohortsDefaultBody{
			Message: fmt.Sprintf("During execution of CheckQueryID"),
		})
	}
	logrus.Trace("has ID", hasID)

	if !hasID {
		return medco_node.NewPostCohortsDefault(400).WithPayload(&medco_node.PostCohortsDefaultBody{
			Message: fmt.Sprintf("User does not have a stored query result with ID: %d", cohort.PatientSetID),
		})
	}

	creationDate, err := time.Parse(time.RFC3339, cohort.CreationDate)
	if err != nil {
		return medco_node.NewPostCohortsDefault(400).WithPayload(&medco_node.PostCohortsDefaultBody{
			Message: fmt.Sprintf("String %s is not a date with RF3339 layout", cohort.CreationDate),
		})
	}
	updateDate, err := time.Parse(time.RFC3339, cohort.UpdateDate)
	if err != nil {
		return medco_node.NewPostCohortsDefault(400).WithPayload(&medco_node.PostCohortsDefaultBody{
			Message: fmt.Sprintf("String %s is not a date with RF3339 layout", cohort.UpdateDate),
		})
	}
	cohorts, err := querytoolsserver.GetSavedCohorts(utilserver.DBConnection, principal.ID)
	if err != nil {
		return medco_node.NewPostCohortsDefault(500).WithPayload(&medco_node.PostCohortsDefaultBody{
			Message: "Get cohort execution error: " + err.Error(),
		})
	}
	for _, existingCohort := range cohorts {
		if existingCohort.CohortName == cohort.CohortName {
			if existingCohort.UpdateDate.After(updateDate) {
				return medco_node.NewPostCohortsDefault(400).WithPayload(&medco_node.PostCohortsDefaultBody{
					Message: fmt.Sprintf(
						"Cohort %s  has a more recent date in DB %s, provided %s",
						cohort.CohortName,
						cohort.UpdateDate,
						existingCohort.UpdateDate.Format(time.RFC3339)),
				})
			}
			break
		}
	}
	querytoolsserver.InsertCohort(utilserver.DBConnection, principal.ID, int(cohort.PatientSetID), cohort.CohortName, creationDate, updateDate)

	return medco_node.NewPostCohortsOK()
}

// MedCoNodeDeleteCohortsHandler handles DELETE /medco/node/explore/cohorts  API endpoint
func MedCoNodeDeleteCohortsHandler(params medco_node.DeleteCohortsParams, principal *models.User) middleware.Responder {
	cohortName := params.Body
	user := principal.ID

	// check if cohort exists
	hasCohort, err := querytoolsserver.DoesCohortExist(utilserver.DBConnection, user, cohortName)
	if err != nil {
		return medco_node.NewDeleteCohortsDefault(500).WithPayload(&medco_node.DeleteCohortsDefaultBody{
			Message: "Delete cohort execution error: " + err.Error(),
		})
	}
	logrus.Trace("hasCohort", hasCohort)
	if !hasCohort {
		return medco_node.NewDeleteCohortsDefault(400).WithPayload(&medco_node.DeleteCohortsDefaultBody{
			Message: "Cohort does not exist",
		})
	}

	// delete the cohorts
	err = querytoolsserver.RemoveCohort(utilserver.DBConnection, user, cohortName)
	if err != nil {
		return medco_node.NewDeleteCohortsDefault(500).WithPayload(&medco_node.DeleteCohortsDefaultBody{
			Message: "Delete cohort execution error: " + err.Error(),
		})
	}

	return medco_node.NewDeleteCohortsOK()

}
