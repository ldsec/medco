package handlers

import (
	"fmt"
	"time"

	querytoolsserver "github.com/ldsec/medco/connector/server/querytools"

	"github.com/go-openapi/runtime/middleware"
	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/ldsec/medco/connector/restapi/server/operations/medco_node"
	medcoserver "github.com/ldsec/medco/connector/server/explore"
	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/ldsec/medco/connector/wrappers/i2b2"
	"github.com/sirupsen/logrus"
)

// MedCoNodeExploreSearchConceptChildrenHandler handles the /medco/node/explore/search/concept-children API endpoint
func MedCoNodeExploreSearchConceptChildrenHandler(params medco_node.ExploreSearchConceptChildrenParams, principal *models.User) middleware.Responder {

	searchResult1, err := i2b2.GetOntologyConceptChildren(*params.SearchConceptChildrenRequest.Path)
	if err != nil {
		return medco_node.NewExploreSearchConceptChildrenDefault(500).WithPayload(&medco_node.ExploreSearchConceptChildrenDefaultBody{
			Message: err.Error(),
		})
	}

	var searchResult2 []*models.ExploreSearchResultElement

	if *params.SearchConceptChildrenRequest.Path != "/" {
		searchResult2, err = i2b2.GetOntologyModifiers(*params.SearchConceptChildrenRequest.Path)
		if err != nil {
			return medco_node.NewExploreSearchConceptChildrenDefault(500).WithPayload(&medco_node.ExploreSearchConceptChildrenDefaultBody{
				Message: err.Error(),
			})
		}
	}

	return medco_node.NewExploreSearchConceptChildrenOK().WithPayload(&medco_node.ExploreSearchConceptChildrenOKBody{
		Search:  params.SearchConceptChildrenRequest,
		Results: append(searchResult1, searchResult2...),
	})
}

// MedCoNodeExploreSearchModifierChildrenHandler handles the /medco/node/explore/search/modifier-children API endpoint
func MedCoNodeExploreSearchModifierChildrenHandler(params medco_node.ExploreSearchModifierChildrenParams, principal *models.User) middleware.Responder {

	searchResult, err := i2b2.GetOntologyModifierChildren(*params.SearchModifierChildrenRequest.Path, *params.SearchModifierChildrenRequest.AppliedPath, *params.SearchModifierChildrenRequest.AppliedConcept)
	if err != nil {
		return medco_node.NewExploreSearchModifierChildrenDefault(500).WithPayload(&medco_node.ExploreSearchModifierChildrenDefaultBody{
			Message: err.Error(),
		})
	}

	return medco_node.NewExploreSearchModifierChildrenOK().WithPayload(&medco_node.ExploreSearchModifierChildrenOKBody{
		Search:  params.SearchModifierChildrenRequest,
		Results: searchResult,
	})
}

// MedCoNodeExploreSearchConceptInfoHandler handles the /medco/node/explore/search/concept-info API endpoint
func MedCoNodeExploreSearchConceptInfoHandler(params medco_node.ExploreSearchConceptInfoParams, principal *models.User) middleware.Responder {

	searchResult, err := i2b2.GetOntologyConceptInfo(*params.SearchConceptInfoRequest.Path)
	if err != nil {
		return medco_node.NewExploreSearchConceptInfoDefault(500).WithPayload(&medco_node.ExploreSearchConceptInfoDefaultBody{
			Message: err.Error(),
		})
	}

	return medco_node.NewExploreSearchConceptInfoOK().WithPayload(&medco_node.ExploreSearchConceptInfoOKBody{
		Search:  params.SearchConceptInfoRequest,
		Results: searchResult,
	})
}

// MedCoNodeExploreSearchModifierInfoHandler handles the /medco/node/explore/search/modifier-info API endpoint
func MedCoNodeExploreSearchModifierInfoHandler(params medco_node.ExploreSearchModifierInfoParams, principal *models.User) middleware.Responder {

	searchResult, err := i2b2.GetOntologyModifierInfo(*params.SearchModifierInfoRequest.Path, *params.SearchModifierInfoRequest.AppliedPath)
	if err != nil {
		return medco_node.NewExploreSearchModifierInfoDefault(500).WithPayload(&medco_node.ExploreSearchModifierInfoDefaultBody{
			Message: err.Error(),
		})
	}

	return medco_node.NewExploreSearchModifierInfoOK().WithPayload(&medco_node.ExploreSearchModifierInfoOKBody{
		Search:  params.SearchModifierInfoRequest,
		Results: searchResult,
	})
}

// MedCoNodeExploreQueryHandler handles /medco/node/explore/query API endpoint
func MedCoNodeExploreQueryHandler(params medco_node.ExploreQueryParams, principal *models.User) middleware.Responder {

	// authorizations of query
	authorizedQueryType, err := utilserver.FetchAuthorizedExploreQueryType(principal)
	if err != nil {
		return medco_node.NewExploreQueryDefault(403).WithPayload(&medco_node.ExploreQueryDefaultBody{
			Message: "Authorization of query failed: " + err.Error(),
		})
	}

	// create query
	query, err := medcoserver.NewExploreQuery(params.QueryRequest.ID, params.QueryRequest.Query, principal)
	if err != nil {
		return medco_node.NewExploreQueryDefault(400).WithPayload(&medco_node.ExploreQueryDefaultBody{
			Message: "Bad query: " + err.Error(),
		})
	}

	// parse query type
	queryType, err := medcoserver.NewExploreQueryType(authorizedQueryType)
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
			Timers:               timers,
			Status:               models.ExploreQueryResultElementStatusAvailable,
			PatientSetID:         int64(query.Result.PatientSetID),
		}})
}

// MedCoNodeGetCohortsHandler handles GET /medco/node/explore/cohorts  API endpoint
func MedCoNodeGetCohortsHandler(params medco_node.GetCohortsParams, principal *models.User) middleware.Responder {
	userID := principal.ID
	cohorts, err := querytoolsserver.GetSavedCohorts(userID, int(*params.Limit))
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

	cohort := params.CohortRequest
	cohortName := params.Name

	hasID, err := querytoolsserver.CheckQueryID(principal.ID, int(*cohort.PatientSetID))
	if err != nil {
		return medco_node.NewPostCohortsDefault(500).WithPayload(&medco_node.PostCohortsDefaultBody{
			Message: fmt.Sprintf("During execution of CheckQueryID"),
		})
	}
	logrus.Trace("has ID", hasID)

	if !hasID {

		return medco_node.NewPostCohortsNotFound().WithPayload(&medco_node.PostCohortsNotFoundBody{
			Message: fmt.Sprintf("User does not have a stored query result with ID: %d", *cohort.PatientSetID),
		})
	}

	creationDate, err := time.Parse(time.RFC3339, *cohort.CreationDate)
	if err != nil {
		return medco_node.NewPostCohortsBadRequest().WithPayload(&medco_node.PostCohortsBadRequestBody{
			Message: fmt.Sprintf("String %s is not a date with RFC3339 layout", *cohort.CreationDate),
		})
	}
	updateDate, err := time.Parse(time.RFC3339, *cohort.UpdateDate)
	if err != nil {
		return medco_node.NewPostCohortsBadRequest().WithPayload(&medco_node.PostCohortsBadRequestBody{
			Message: fmt.Sprintf("String %s is not a date with RFC3339 layout", *cohort.UpdateDate),
		})
	}
	cohorts, err := querytoolsserver.GetSavedCohorts(principal.ID, 0)
	if err != nil {
		return medco_node.NewPostCohortsDefault(500).WithPayload(&medco_node.PostCohortsDefaultBody{
			Message: "Get cohort execution error: " + err.Error(),
		})
	}
	for _, existingCohort := range cohorts {
		if existingCohort.CohortName == cohortName {
			return medco_node.NewPostCohortsConflict().WithPayload(&medco_node.PostCohortsConflictBody{
				Message: "Cohort %s already exists. Try update-saved-cohorts instead of add-saved-cohorts",
			})
		}
	}
	querytoolsserver.InsertCohort(principal.ID, int(*cohort.PatientSetID), cohortName, creationDate, updateDate)

	return medco_node.NewPostCohortsOK()
}

// MedCoNodePutCohortsHandler handles PUT /medco/node/explore/cohorts  API endpoint
func MedCoNodePutCohortsHandler(params medco_node.PutCohortsParams, principal *models.User) middleware.Responder {

	cohort := params.CohortRequest
	cohortName := params.Name

	hasID, err := querytoolsserver.CheckQueryID(principal.ID, int(*cohort.PatientSetID))
	if err != nil {
		return medco_node.NewPutCohortsDefault(500).WithPayload(&medco_node.PutCohortsDefaultBody{
			Message: fmt.Sprintf("User does not have a stored query result with ID: %d", *cohort.PatientSetID),
		})
	}
	logrus.Trace("has ID", hasID)

	if !hasID {
		return medco_node.NewPutCohortsNotFound().WithPayload(&medco_node.PutCohortsNotFoundBody{
			Message: fmt.Sprintf("There is no result instance with id %d", int(*cohort.PatientSetID)),
		})
	}

	updateDate, err := time.Parse(time.RFC3339, *cohort.UpdateDate)
	if err != nil {
		return medco_node.NewPutCohortsBadRequest().WithPayload(&medco_node.PutCohortsBadRequestBody{
			Message: fmt.Sprintf("String %s is not a date with RF3339 layout", *cohort.UpdateDate),
		})
	}
	cohorts, err := querytoolsserver.GetSavedCohorts(principal.ID, 0)
	if err != nil {
		return medco_node.NewPutCohortsDefault(500).WithPayload(&medco_node.PutCohortsDefaultBody{
			Message: "Get cohort execution error: " + err.Error(),
		})
	}
	found := false
	for _, existingCohort := range cohorts {
		if existingCohort.CohortName == cohortName {
			if existingCohort.UpdateDate.After(updateDate) {
				return medco_node.NewPutCohortsConflict().WithPayload(&medco_node.PutCohortsConflictBody{
					Message: fmt.Sprintf("The cohort update date is more recent in server DB than the date provided by client. Server: %s, client: %s.",
						existingCohort.UpdateDate.Format(time.RFC3339),
						updateDate.Format(time.RFC3339),
					),
				})
			}
			found = true
			break
		}
	}
	if !found {
		return medco_node.NewPutCohortsNotFound()
	}

	querytoolsserver.UpdateCohort(cohortName, principal.ID, int(*cohort.PatientSetID), updateDate)

	return medco_node.NewPutCohortsOK()
}

// MedCoNodeDeleteCohortsHandler handles DELETE /medco/node/explore/cohorts  API endpoint
func MedCoNodeDeleteCohortsHandler(params medco_node.DeleteCohortsParams, principal *models.User) middleware.Responder {
	cohortName := params.Name
	user := principal.ID

	// check if cohort exists
	hasCohort, err := querytoolsserver.DoesCohortExist(user, cohortName)
	if err != nil {
		return medco_node.NewDeleteCohortsDefault(500).WithPayload(&medco_node.DeleteCohortsDefaultBody{
			Message: "Delete cohort execution error: " + err.Error(),
		})
	}
	logrus.Trace("hasCohort", hasCohort)
	if !hasCohort {
		return medco_node.NewDeleteCohortsNotFound().WithPayload(&medco_node.DeleteCohortsNotFoundBody{
			Message: fmt.Sprintf("Cohort %s not found", cohortName),
		})
	}

	// delete the cohorts
	err = querytoolsserver.RemoveCohort(user, cohortName)
	if err != nil {
		return medco_node.NewDeleteCohortsDefault(500).WithPayload(&medco_node.DeleteCohortsDefaultBody{
			Message: "Delete cohort execution error: " + err.Error(),
		})
	}

	return medco_node.NewDeleteCohortsOK()

}
