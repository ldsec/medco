package explorestatisticsclient

import (
	"time"

	"github.com/sirupsen/logrus"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ldsec/medco/connector/restapi/client/explore_statistics"
	utilclient "github.com/ldsec/medco/connector/util/client"
)

func (exploreStatistics *ExploreStatistics) submitToNode(nodeIdx int) (results *explore_statistics.ExploreStatisticsOKBody, err error) {

	params := explore_statistics.NewExploreStatisticsParamsWithTimeout(time.Duration(utilclient.ExploreStatisticsTimeoutSeconds) * time.Second)

	var cohortDef = &explore_statistics.ExploreStatisticsParamsBodyCohortDefinition{
		Panels:      exploreStatistics.cohortDefinition.Panels,
		QueryTiming: exploreStatistics.cohortDefinition.QueryTiming,
	}

	var modifiers []*explore_statistics.ExploreStatisticsParamsBodyModifiersItems0

	for _, m := range exploreStatistics.modifiers {
		modifiers = append(modifiers, &explore_statistics.ExploreStatisticsParamsBodyModifiersItems0{
			AppliedPath: m.AppliedPath,
			ModifierKey: m.ModifierKey,
		})
	}
	body := explore_statistics.ExploreStatisticsBody{
		ID:               exploreStatistics.id,
		CohortDefinition: cohortDef,
		Concepts:         exploreStatistics.conceptsPaths,
		Modifiers:        modifiers,
		BucketSize:       1, //TODO change that to something that makes sense
		UserPublicKey:    exploreStatistics.userPublicKey,
	}

	logrus.Info("Submitting to node with index ", nodeIdx, " the query with body ", body)

	params.SetBody(body)
	response, err := exploreStatistics.httpMedCoClients[nodeIdx].ExploreStatistics.ExploreStatistics(params, httptransport.BearerToken(exploreStatistics.authToken))

	if err != nil {
		logrus.Error("Explore statistics error: ", err)
		return
	}
	results = response.GetPayload()

	return
}
