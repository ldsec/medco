package survivalclient

import (
	"context"
	"net/http"
	"time"

	"github.com/sirupsen/logrus"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ldsec/medco-connector/restapi/client/survival_analysis"
	utilclient "github.com/ldsec/medco-connector/util/client"
)

//for client !!

func (clientSurvivalAnalysis *SurvivalAnalysis) submitToNode(nodeIdx int) (results *survival_analysis.GetSurvivalAnalysisOKBody, err error) {
	//magicNumber
	params := survival_analysis.NewGetSurvivalAnalysisParamsWithTimeout(time.Duration(utilclient.QueryTimeoutSeconds) * time.Second)

	body := &survival_analysis.GetSurvivalAnalysisBody{
		ID:                  clientSurvivalAnalysis.id,
		UserPublicKey:       clientSurvivalAnalysis.userPublicKey,
		SetID:               float64(clientSurvivalAnalysis.patientSetID),
		SubGroupDefinitions: clientSurvivalAnalysis.subGroupDefinitions,
		StartConcept:        clientSurvivalAnalysis.startConcept,
		StartColumn:         clientSurvivalAnalysis.startColumn,
		EndConcept:          clientSurvivalAnalysis.endConcept,
		EndColumn:           clientSurvivalAnalysis.endColumn,
	}
	params.SetBody(*body)
	response, err := clientSurvivalAnalysis.httpMedCoClients[nodeIdx].SurvivalAnalysis.GetSurvivalAnalysis(params, httptransport.BearerToken(clientSurvivalAnalysis.authToken))

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
