package survivalclient

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/sirupsen/logrus"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ldsec/medco-connector/restapi/client/survival_analysis"
	utilclient "github.com/ldsec/medco-connector/util/client"
)

func (clientSurvivalAnalysis *SurvivalAnalysis) submitToNode(nodeIdx int) (results *survival_analysis.SurvivalAnalysisOKBody, err error) {

	params := survival_analysis.NewSurvivalAnalysisParamsWithTimeout(time.Duration(utilclient.QueryTimeoutSeconds) * time.Second)

	body := &survival_analysis.SurvivalAnalysisBody{
		ID:                  clientSurvivalAnalysis.id,
		UserPublicKey:       clientSurvivalAnalysis.userPublicKey,
		SetID:               int64(clientSurvivalAnalysis.patientSetID),
		SubGroupDefinitions: clientSurvivalAnalysis.subGroupDefinitions,
		StartConcept:        clientSurvivalAnalysis.startConceptPath,
		StartModifier:       clientSurvivalAnalysis.startModifierCode,
		EndConcept:          clientSurvivalAnalysis.endConceptPath,
		EndModifier:         clientSurvivalAnalysis.endModifierCode,
		TimeGranularity:     strings.ToLower(clientSurvivalAnalysis.granularity),
		TimeLimit:           int64(clientSurvivalAnalysis.limit),
	}
	params.SetBody(*body)
	response, err := clientSurvivalAnalysis.httpMedCoClients[nodeIdx].SurvivalAnalysis.SurvivalAnalysis(params, httptransport.BearerToken(clientSurvivalAnalysis.authToken))

	if err != nil {
		logrus.Error("survival analysis error: ", err)
		return
	}
	results = response.GetPayload()

	return
}

// GetSurvivalAnalysisParameter holds the information for interaction with the REST API
type GetSurvivalAnalysisParameter struct {
	Command    *SurvivalAnalysis
	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}
