package handlers

import (
	"fmt"
	"time"

	medcoserver "github.com/ldsec/medco/connector/server/explore"

	querytoolsserver "github.com/ldsec/medco/connector/server/querytools"

	"github.com/go-openapi/runtime/middleware"
	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/ldsec/medco/connector/restapi/server/operations/medco_node"
	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/ldsec/medco/connector/wrappers/i2b2"
	"github.com/sirupsen/logrus"
)

// MedCoNodeExploreSearchHandler handles the /medco/node/explore/search API endpoint.
func MedCoNodeExploreSearchHandler(params medco_node.ExploreSearchParams, principal *models.User) middleware.Responder {

	logrus.Infof("received request at /medco/node/explore/search with parameters: %s %d", *params.SearchRequest.SearchString, params.SearchRequest.Limit)

	var searchResult []*models.ExploreSearchResultElement

	searchResult, err := i2b2.GetOntologyElements(*params.SearchRequest.SearchString, params.SearchRequest.Limit)
	if err != nil {
		return medco_node.NewExploreSearchDefault(500).WithPayload(&medco_node.ExploreSearchDefaultBody{
			Message: "Search execution error: " + err.Error(),
		})
	}

	return medco_node.NewExploreSearchOK().WithPayload(&medco_node.ExploreSearchOKBody{
		Search:  params.SearchRequest,
		Results: searchResult,
	})

}

// MedCoNodeExploreSearchConceptHandler handles the /medco/node/explore/search/concept API endpoint
func MedCoNodeExploreSearchConceptHandler(params medco_node.ExploreSearchConceptParams, principal *models.User) middleware.Responder {

	var searchResult []*models.ExploreSearchResultElement

	if *params.SearchConceptRequest.Operation == models.ExploreSearchConceptOperationChildren {

		searchResult1, err := i2b2.GetOntologyConceptChildren(*params.SearchConceptRequest.Path)
		if err != nil {
			return medco_node.NewExploreSearchConceptDefault(500).WithPayload(&medco_node.ExploreSearchConceptDefaultBody{
				Message: err.Error(),
			})
		}

		var searchResult2 []*models.ExploreSearchResultElement

		if *params.SearchConceptRequest.Path != "/" {
			searchResult2, err = i2b2.GetOntologyModifiers(*params.SearchConceptRequest.Path)
			if err != nil {
				return medco_node.NewExploreSearchConceptDefault(500).WithPayload(&medco_node.ExploreSearchConceptDefaultBody{
					Message: err.Error(),
				})
			}
		}

		searchResult = append(searchResult1, searchResult2...)

	} else if *params.SearchConceptRequest.Operation == models.ExploreSearchConceptOperationInfo {

		var err error
		searchResult, err = i2b2.GetOntologyConceptInfo(*params.SearchConceptRequest.Path)
		if err != nil {
			return medco_node.NewExploreSearchConceptDefault(500).WithPayload(&medco_node.ExploreSearchConceptDefaultBody{
				Message: err.Error(),
			})
		}
	}

	return medco_node.NewExploreSearchConceptOK().WithPayload(&medco_node.ExploreSearchConceptOKBody{
		Search:  params.SearchConceptRequest,
		Results: searchResult,
	})

}

// MedCoNodeExploreSearchModifierHandler handles the /medco/node/explore/search/modifier API endpoint
func MedCoNodeExploreSearchModifierHandler(params medco_node.ExploreSearchModifierParams, principal *models.User) middleware.Responder {

	var searchResult []*models.ExploreSearchResultElement
	var err error

	if *params.SearchModifierRequest.Operation == models.ExploreSearchModifierOperationChildren {
		searchResult, err = i2b2.GetOntologyModifierChildren(*params.SearchModifierRequest.Path, *params.SearchModifierRequest.AppliedPath, *params.SearchModifierRequest.AppliedConcept)
		if err != nil {
			return medco_node.NewExploreSearchModifierDefault(500).WithPayload(&medco_node.ExploreSearchModifierDefaultBody{
				Message: err.Error(),
			})
		}
	} else if *params.SearchModifierRequest.Operation == models.ExploreSearchModifierOperationInfo {
		searchResult, err = i2b2.GetOntologyModifierInfo(*params.SearchModifierRequest.Path, *params.SearchModifierRequest.AppliedPath)
		if err != nil {
			return medco_node.NewExploreSearchModifierDefault(500).WithPayload(&medco_node.ExploreSearchModifierDefaultBody{
				Message: err.Error(),
			})
		}
	}

	return medco_node.NewExploreSearchModifierOK().WithPayload(&medco_node.ExploreSearchModifierOKBody{
		Search:  params.SearchModifierRequest,
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
			QueryID:              int64(query.Result.QueryID),
		}})
}

// MedCoNodeGetCohortsHandler handles GET /medco/node/explore/cohorts  API endpoint
func MedCoNodeGetCohortsHandler(params medco_node.GetCohortsParams, principal *models.User) middleware.Responder {
	userID := principal.ID
	cohorts, err := querytoolsserver.GetSavedCohorts(userID, int(*params.Limit))
	if err != nil {
		return medco_node.NewGetCohortsDefault(500).WithPayload(&medco_node.GetCohortsDefaultBody{
			Message: "Get cohort execution error: " + err.Error(),
		})
	}
	payload := &medco_node.GetCohortsOK{}
	for _, cohort := range cohorts {
		queryID := int64(cohort.QueryID)
		var queryDefinition *medco_node.GetCohortsOKBodyItems0QueryDefinition
		queryDefinitionString, err := querytoolsserver.GetQueryDefinition(int(queryID))
		if err != nil {
			medco_node.NewGetCohortsDefault(500).WithPayload(&medco_node.GetCohortsDefaultBody{
				Message: "Get cohort execution error: " + err.Error(),
			})
		}
		queryModel := &models.ExploreQuery{}
		err = queryModel.UnmarshalBinary([]byte(queryDefinitionString))
		if err != nil {
			logrus.Warnf("skipping query definition unmarshalling with string %s: %s", queryDefinitionString, err.Error())
		} else {
			queryDefinition = &medco_node.GetCohortsOKBodyItems0QueryDefinition{
				Panels:      queryModel.Panels,
				QueryTiming: queryModel.QueryTiming,
			}
		}

		payload.Payload = append(payload.Payload,
			&medco_node.GetCohortsOKBodyItems0{
				CohortName:      cohort.CohortName,
				CohortID:        int64(cohort.CohortID),
				QueryID:         queryID,
				CreationDate:    cohort.CreationDate.Format(time.RFC3339Nano),
				UpdateDate:      cohort.UpdateDate.Format(time.RFC3339Nano),
				QueryDefinition: queryDefinition,
			},
		)
	}

	return medco_node.NewGetCohortsOK().WithPayload(payload.Payload)

}

// MedCoNodePostCohortsPatientListHandler handles POST /medco/node/explore/cohorts/patientList  API endpoint
func MedCoNodePostCohortsPatientListHandler(params medco_node.PostCohortsPatientListParams, principal *models.User) middleware.Responder {
	body := params.CohortsPatientListRequest
	authorizedQueryType, err := utilserver.FetchAuthorizedExploreQueryType(principal)
	if err != nil {
		return medco_node.NewPostCohortsPatientListForbidden().WithPayload(&medco_node.PostCohortsPatientListForbiddenBody{
			Message: fmt.Sprintf("User %s is not authorized to use MedCo Explore", principal.ID),
		})
	}

	// if the user is authorized to get the patient list in Explore Query, it is the same for getting list from saved cohorts
	queryType, err := medcoserver.NewExploreQueryType(authorizedQueryType)
	if err != nil {
		return medco_node.NewPostCohortsPatientListNotFound().WithPayload(&medco_node.PostCohortsPatientListNotFoundBody{
			Message: "Bad query type: " + err.Error(),
		})
	}

	if !queryType.PatientList {
		return medco_node.NewPostCohortsPatientListForbidden().WithPayload(&medco_node.PostCohortsPatientListForbiddenBody{
			Message: fmt.Sprintf("User %s is not authorized to query patient lists", principal.ID),
		})
	}

	// check if the cohort exists
	doesExist, err := querytoolsserver.DoesCohortExist(principal.ID, *body.CohortName)
	if err != nil {
		return medco_node.NewPostCohortsPatientListDefault(500).WithPayload(&medco_node.PostCohortsPatientListDefaultBody{
			Message: "DB error while checking if the cohort exists: " + err.Error(),
		})
	}
	if !doesExist {
		return medco_node.NewPostCohortsPatientListNotFound().WithPayload(&medco_node.PostCohortsPatientListNotFoundBody{
			Message: fmt.Sprintf("No cohort named %s found for user %s", *body.CohortName, principal.ID),
		})
	}

	// user is authorized and cohort exists
	cohortsPatientList := querytoolsserver.NewCohortsPatientList(body.ID, *body.UserPublicKey, *body.CohortName, principal)

	err = cohortsPatientList.Execute()
	if err != nil {
		return medco_node.NewPostCohortsPatientListDefault(500).WithPayload(&medco_node.PostCohortsPatientListDefaultBody{
			Message: "During execution of cohorts patient list: " + err.Error(),
		})
	}

	payload := &medco_node.PostCohortsPatientListOKBody{
		Results: cohortsPatientList.Result.EncPatientList,
		Timers:  cohortsPatientList.Result.Timers.TimersToAPIModel(),
	}
	return medco_node.NewPostCohortsPatientListOK().WithPayload(payload)

}

// MedCoNodePostCohortsHandler handles POST /medco/node/explore/cohorts  API endpoint
func MedCoNodePostCohortsHandler(params medco_node.PostCohortsParams, principal *models.User) middleware.Responder {

	cohort := params.CohortsRequest
	cohortName := params.Name

	hasID, err := querytoolsserver.CheckQueryID(principal.ID, int(*cohort.QueryID))
	if err != nil {
		return medco_node.NewPostCohortsDefault(500).WithPayload(&medco_node.PostCohortsDefaultBody{
			Message: fmt.Sprintf("During execution of CheckQueryID"),
		})
	}
	logrus.Trace("has ID", hasID)

	if !hasID {

		return medco_node.NewPostCohortsNotFound().WithPayload(&medco_node.PostCohortsNotFoundBody{
			Message: fmt.Sprintf("User %s does not have a stored query result with ID: %d", principal.ID, *cohort.QueryID),
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
				Message: fmt.Sprintf("Cohort %s already exists. Try update-saved-cohorts instead of add-saved-cohorts", cohortName),
			})
		}
	}
	querytoolsserver.InsertCohort(principal.ID, int(*cohort.QueryID), cohortName, creationDate, updateDate)

	return medco_node.NewPostCohortsOK()
}

// MedCoNodePutCohortsHandler handles PUT /medco/node/explore/cohorts  API endpoint
func MedCoNodePutCohortsHandler(params medco_node.PutCohortsParams, principal *models.User) middleware.Responder {

	cohort := params.CohortsRequest
	cohortName := params.Name

	hasID, err := querytoolsserver.CheckQueryID(principal.ID, int(*cohort.QueryID))
	if err != nil {
		return medco_node.NewPutCohortsDefault(500).WithPayload(&medco_node.PutCohortsDefaultBody{
			Message: fmt.Sprintf("User does not have a stored query result with ID: %d", *cohort.QueryID),
		})
	}
	logrus.Trace("has ID", hasID)

	if !hasID {
		return medco_node.NewPutCohortsNotFound().WithPayload(&medco_node.PutCohortsNotFoundBody{
			Message: fmt.Sprintf("There is no result instance with id %d", int(*cohort.QueryID)),
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
						existingCohort.UpdateDate.Format(time.RFC3339Nano),
						updateDate.Format(time.RFC3339Nano),
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

	querytoolsserver.UpdateCohort(cohortName, principal.ID, int(*cohort.QueryID), updateDate)

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
