package handlers

import (
	"strconv"

	utilserver "github.com/ldsec/medco/connector/util/server"
	"go.dedis.ch/onet/v3/log"

	medcoserver "github.com/ldsec/medco/connector/server/explore"

	"github.com/go-openapi/runtime/middleware"
	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/ldsec/medco/connector/restapi/server/operations/explore_statistics"
	referenceintervalserver "github.com/ldsec/medco/connector/server/reference_intervals"
	"github.com/sirupsen/logrus"
)

//ExploreStatisticsHandler handles /survival-analysis API endpoint. It is the entrypoint for creating a histogram about a concept's or modifier's observations.
func ExploreStatisticsHandler(param explore_statistics.ExploreStatisticsParams, principal *models.User) middleware.Responder {
	/*
	* The histogram is created by cutting the space of observations in equal sized interval. The number of interval is defined by the parameters of the query.
	* The observations are then fetched from the database and classified in each interval. Then the count of observations is determined per interval.
	* Those counts are then aggregated between all nodes. And the result is returned to the user.
	 */
	//check users authorisation for global count
	authorizedQueryType, err := utilserver.FetchAuthorizedExploreQueryType(principal)
	logrus.Debug("Authorized query type user ", authorizedQueryType)

	if err != nil {
		logrus.Error(err)
		explore_statistics.NewExploreStatisticsBadRequest().WithPayload(
			&explore_statistics.ExploreStatisticsBadRequestBody{
				Message: "Explore statistics query creation error:" + err.Error()})
	}

	queryType, err := medcoserver.NewExploreQueryType(authorizedQueryType)
	if err != nil {
		logrus.Error(err)
		explore_statistics.NewExploreStatisticsBadRequest().WithPayload(
			&explore_statistics.ExploreStatisticsBadRequestBody{
				Message: "Explore statistics query creation error:" + err.Error()})
	}

	if queryType.Obfuscated {
		explore_statistics.NewExploreStatisticsDefault(401).WithPayload(
			&explore_statistics.ExploreStatisticsDefaultBody{
				Message: "No authorization to perform such request. Only authorized to see obsfuscated result."})
	}

	//the user is authorized to do global count queries. We create a new query
	query, err := referenceintervalserver.NewQuery(
		principal.ID,
		param.Body)

	query.QueryType = queryType

	if err != nil {
		logrus.Error(err)
		explore_statistics.NewExploreStatisticsBadRequest().WithPayload(
			&explore_statistics.ExploreStatisticsBadRequestBody{
				Message: "Explore statistics query creation error:" + err.Error()})

	}

	logrus.Debug("explore statistics: ", query)
	err = query.Validate()
	if err != nil {
		logrus.Error(err)
		explore_statistics.NewExploreStatisticsBadRequest().WithPayload(
			&explore_statistics.ExploreStatisticsBadRequestBody{
				Message: "Explore statistics query validation error:" + err.Error()})

	}

	err = query.Execute(principal)

	if err != nil {
		logrus.Error(err)
		explore_statistics.NewExploreStatisticsDefault(500).WithPayload(
			&explore_statistics.ExploreStatisticsDefaultBody{
				Message: "Explore statistics query execution error:" + err.Error()})
	}

	//parse timers
	//TODO use global timers
	modelsTimers := query.Response.GlobalTimers.TimersToAPIModel()

	requestResult := explore_statistics.ExploreStatisticsOKBody{
		GlobalTimers: modelsTimers,
	}

	requestResult.CohortQueryID = int64(query.Response.QueryID)
	if queryType.PatientList {
		requestResult.EncryptedPatientList = query.Response.EncPatientList
		requestResult.PatientSetID = int64(query.Response.PatientSetID)
		requestResult.EncryptedCohortCount = query.Response.EncCount
	}

	for _, result := range query.Response.Results {
		apiIntervals := make([]*models.IntervalBucket, 0)
		for _, bucket := range result.Intervals {
			higherBound := strconv.FormatFloat(bucket.HigherBound, 'f', 5, 64)
			lowerBound := strconv.FormatFloat(bucket.LowerBound, 'f', 5, 64)

			apiBucket := models.IntervalBucket{
				EncCount:    &bucket.EncCount,
				HigherBound: &higherBound,
				LowerBound:  &lowerBound,
			}

			apiIntervals = append(apiIntervals, &apiBucket)
		}

		// Each analyte given as input to the http handler is associated to a result.
		requestResult.Results = append(requestResult.Results, &explore_statistics.ExploreStatisticsOKBodyResultsItems0{
			Timers:      result.Timers.TimersToAPIModel(),
			AnalyteName: result.AnalyteName,
			Intervals:   apiIntervals,
			Unit:        result.Unit,
		})

		intervalStr := ""
		for _, interval := range apiIntervals {
			intervalStr += "[" + *interval.LowerBound + ", " + *interval.HigherBound + "]"
		}
		log.Infof("Intervals for analyte %s, intervals: %s, unit: %s", result.AnalyteName, intervalStr, result.Unit)
	}

	return explore_statistics.NewExploreStatisticsOK().WithPayload(&requestResult)

}
