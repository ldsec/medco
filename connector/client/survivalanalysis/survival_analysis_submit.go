package survivalclient

import (
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ldsec/medco/connector/restapi/client/survival_analysis"
	utilclient "github.com/ldsec/medco/connector/util/client"
)

func (clientSurvivalAnalysis *SurvivalAnalysis) submitToNode(nodeIdx int) (results *survival_analysis.SurvivalAnalysisOKBody, err error) {

	params := survival_analysis.NewSurvivalAnalysisParamsWithTimeout(time.Duration(utilclient.SurvivalAnalysisTimeoutSeconds) * time.Second)

	body := survival_analysis.SurvivalAnalysisBody{
		ID:                  new(string),
		UserPublicKey:       new(string),
		CohortName:          new(string),
		SubGroupDefinitions: clientSurvivalAnalysis.subGroupDefinitions,
		StartConcept:        new(string),
		StartModifier:       clientSurvivalAnalysis.startModifier,
		StartsWhen:          new(string),
		EndConcept:          new(string),
		EndModifier:         clientSurvivalAnalysis.endModifier,
		EndsWhen:            new(string),
		TimeGranularity:     new(string),
		TimeLimit:           new(int64),
		CensoringFrom:       new(string),
	}

	*body.ID = clientSurvivalAnalysis.id
	*body.UserPublicKey = clientSurvivalAnalysis.userPublicKey
	*body.CohortName = clientSurvivalAnalysis.cohortName
	*body.StartConcept = clientSurvivalAnalysis.startConceptPath
	*body.StartsWhen = clientSurvivalAnalysis.startsWhen
	*body.EndConcept = clientSurvivalAnalysis.endConceptPath
	*body.EndsWhen = clientSurvivalAnalysis.endsWhen
	*body.TimeGranularity = strings.ToLower(clientSurvivalAnalysis.granularity)
	*body.TimeLimit = int64(clientSurvivalAnalysis.limit)
	*body.CensoringFrom = clientSurvivalAnalysis.censoringFrom

	params.SetBody(body)
	response, err := clientSurvivalAnalysis.httpMedCoClients[nodeIdx].SurvivalAnalysis.SurvivalAnalysis(params, httptransport.BearerToken(clientSurvivalAnalysis.authToken))

	if err != nil {
		logrus.Error("survival analysis error: ", err)
		return
	}
	results = response.GetPayload()

	return
}
