// This file is safe to edit. Once it exists it will not be overwritten

package server

import (
	"crypto/tls"
	"net/http"

	"github.com/ldsec/medco/connector/server/handlers"
	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/sirupsen/logrus"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/ldsec/medco/connector/restapi/server/operations"
	"github.com/ldsec/medco/connector/restapi/server/operations/genomic_annotations"
	"github.com/ldsec/medco/connector/restapi/server/operations/medco_network"
	"github.com/ldsec/medco/connector/restapi/server/operations/medco_node"
	"github.com/ldsec/medco/connector/restapi/server/operations/survival_analysis"

	"github.com/ldsec/medco/connector/restapi/models"
)

//go:generate swagger generate server --target ../../../medco-connector --name MedcoConnector --spec ../../swagger/medco-connector.yml --model-package restapi/models --server-package restapi/server --principal models.User

func configureFlags(api *operations.MedcoConnectorAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MedcoConnectorAPI) http.Handler {
	api.ServeError = errors.ServeError
	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()
	api.Logger = logrus.Printf

	// validate identity and generate principal, check endpoint-based authorizations
	api.MedcoJwtAuth = func(token string, requiredAuthorizations []string) (principal *models.User, err error) {

		// authenticate user
		principal, err = utilserver.AuthenticateUser(token)
		if err != nil {
			return
		}

		// check rest api authorizations
		for _, requiredAuthorization := range requiredAuthorizations {
			err = utilserver.AuthorizeRestAPIEndpoint(principal, models.RestAPIAuthorization(requiredAuthorization))
			if err != nil {
				return
			}
		}

		return
	}

	// /medco/network
	api.MedcoNetworkGetMetadataHandler = medco_network.GetMetadataHandlerFunc(handlers.MedCoNetworkGetMetadataHandler)

	// /medco/node/explore/search/concept
	api.MedcoNodeExploreSearchConceptHandler = medco_node.ExploreSearchConceptHandlerFunc(handlers.MedCoNodeExploreSearchConceptHandler)

	// /medco/node/explore/search/modifier
	api.MedcoNodeExploreSearchModifierHandler = medco_node.ExploreSearchModifierHandlerFunc(handlers.MedCoNodeExploreSearchModifierHandler)

	// /medco/node/explore/query
	api.MedcoNodeExploreQueryHandler = medco_node.ExploreQueryHandlerFunc(handlers.MedCoNodeExploreQueryHandler)

	// /medco/node/explore/query/{queryId}
	api.MedcoNodeGetExploreQueryHandler = medco_node.GetExploreQueryHandlerFunc(func(params medco_node.GetExploreQueryParams, principal *models.User) middleware.Responder {
		return middleware.NotImplemented("operation medco_node.GetQueryResult has not yet been implemented")
	})

	// /medco/genomic-annotations/{annotation}
	api.GenomicAnnotationsGetValuesHandler = genomic_annotations.GetValuesHandlerFunc(handlers.MedCoGenomicAnnotationsGetValuesHandler)

	// /genomic-annotations/{annotation}/{value}
	api.GenomicAnnotationsGetVariantsHandler = genomic_annotations.GetVariantsHandlerFunc(handlers.MedCoGenomicAnnotationsGetVariantsHandler)

	// /node/explore/cohorts
	api.MedcoNodeGetCohortsHandler = medco_node.GetCohortsHandlerFunc(handlers.MedCoNodeGetCohortsHandler)

	// /node/explore/cohorts
	api.MedcoNodePostCohortsHandler = medco_node.PostCohortsHandlerFunc(handlers.MedCoNodePostCohortsHandler)

	// /node/explore/cohorts
	api.MedcoNodePutCohortsHandler = medco_node.PutCohortsHandlerFunc(handlers.MedCoNodePutCohortsHandler)

	// /node/explore/cohorts
	api.MedcoNodeDeleteCohortsHandler = medco_node.DeleteCohortsHandlerFunc(handlers.MedCoNodeDeleteCohortsHandler)

	// /node/explore/cohorts/patientList
	api.MedcoNodePostCohortsPatientListHandler = medco_node.PostCohortsPatientListHandlerFunc(handlers.MedCoNodePostCohortsPatientListHandler)

	// /node/analysis/survival/query
	api.SurvivalAnalysisSurvivalAnalysisHandler = survival_analysis.SurvivalAnalysisHandlerFunc(handlers.MedCoSurvivalAnalysisHandler)

	api.ServerShutdown = func() {}

	return setupGlobalMiddleware(api.Serve(setupMiddlewares))
}

// The TLS configuration before HTTPS server starts.
func configureTLS(tlsConfig *tls.Config) {
	// Make all necessary changes to the TLS configuration here.
}

// As soon as server is initialized but not run yet, this function will be called.
// If you need to modify a config, store server instance to stop it individually later, this is the place.
// This function can be called multiple times, depending on the number of serving schemes.
// scheme value will be set accordingly: "http", "https" or "unix"
func configureServer(s *http.Server, scheme, addr string) {
}

// The middleware configuration is for the handler executors. These do not apply to the swagger.json document.
// The middleware executes after routing but before authentication, binding and validation
func setupMiddlewares(handler http.Handler) http.Handler {
	return handler
}

// The middleware configuration happens before anything, this middleware also applies to serving the swagger.json document.
// So this is a good place to plug in a panic handling middleware, logging and metrics
func setupGlobalMiddleware(handler http.Handler) http.Handler {
	return handler
}
