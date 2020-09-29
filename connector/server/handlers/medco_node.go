package handlers

import (
	"github.com/go-openapi/runtime/middleware"
	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/ldsec/medco/connector/restapi/server/operations/medco_node"
	"github.com/ldsec/medco/connector/server"
	"github.com/ldsec/medco/connector/util/server"
	"github.com/ldsec/medco/connector/wrappers/i2b2"
	"time"
)

// MedCoNodeExploreSearchHandler handles the /medco/node/explore/search/concept API endpoint
func MedCoNodeExploreSearchConceptHandler(params medco_node.ExploreSearchConceptParams, principal *models.User) middleware.Responder {

	searchResult, err := i2b2.GetOntologyChildren(*params.SearchConceptRequest.Path)
	if err != nil {
		return medco_node.NewExploreSearchConceptDefault(500).WithPayload(&medco_node.ExploreSearchConceptDefaultBody{
			Message: err.Error(),
		})
	}

	return medco_node.NewExploreSearchConceptOK().WithPayload(&medco_node.ExploreSearchConceptOKBody{
		Search:  params.SearchConceptRequest,
		Results: searchResult,
	})
}

// MedCoNodeExploreSearchModifierHandler handles the /medco/node/explore/search/modifier API endpoint
func MedCoNodeExploreSearchModifierHandler(params medco_node.ExploreSearchModifierParams, principal *models.User) middleware.Responder {

	searchResult, err := i2b2.GetOntologyModifierChildren(*params.SearchModifierRequest.Path, *params.SearchModifierRequest.AppliedPath, *params.SearchModifierRequest.AppliedConcept)
	if err != nil {
		return medco_node.NewExploreSearchModifierDefault(500).WithPayload(&medco_node.ExploreSearchModifierDefaultBody{
			Message: err.Error(),
		})
	}

	return medco_node.NewExploreSearchModifierOK().WithPayload(&medco_node.ExploreSearchModifierOKBody{
		Search:  params.SearchModifierRequest,
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
			Timers:               timers,
			Status:               models.ExploreQueryResultElementStatusAvailable,
		}})
}
