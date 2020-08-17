// This file is safe to edit. Once it exists it will not be overwritten

package server

import (
	"crypto/tls"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/ldsec/medco-connector/models"
	"github.com/ldsec/medco-connector/restapi/server/operations"
	"github.com/ldsec/medco-connector/restapi/server/operations/genomic_annotations"
	"github.com/ldsec/medco-connector/restapi/server/operations/medco_network"
	"github.com/ldsec/medco-connector/restapi/server/operations/medco_node"
	"github.com/ldsec/medco-connector/restapi/server/operations/survival_analysis"
)

//go:generate swagger generate server --target ../../../medco-connector --name MedcoConnector --spec ../../swagger/medco-connector.yml --model-package restapi/models --server-package restapi/server --principal models.User

func configureFlags(api *operations.MedcoConnectorAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MedcoConnectorAPI) http.Handler {
	// configure the api here
	api.ServeError = errors.ServeError

	// Set your custom logger if needed. Default one is log.Printf
	// Expected interface func(string, ...interface{})
	//
	// Example:
	// api.Logger = log.Printf

	api.UseSwaggerUI()
	// To continue using redoc as your UI, uncomment the following line
	// api.UseRedoc()

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	if api.MedcoJwtAuth == nil {
		api.MedcoJwtAuth = func(token string, scopes []string) (*models.User, error) {
			return nil, errors.NotImplemented("oauth2 bearer auth (medco-jwt) has not yet been implemented")
		}
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	if api.MedcoNodeExploreQueryHandler == nil {
		api.MedcoNodeExploreQueryHandler = medco_node.ExploreQueryHandlerFunc(func(params medco_node.ExploreQueryParams, principal *models.User) middleware.Responder {
			return middleware.NotImplemented("operation medco_node.ExploreQuery has not yet been implemented")
		})
	}
	if api.MedcoNodeExploreSearchHandler == nil {
		api.MedcoNodeExploreSearchHandler = medco_node.ExploreSearchHandlerFunc(func(params medco_node.ExploreSearchParams, principal *models.User) middleware.Responder {
			return middleware.NotImplemented("operation medco_node.ExploreSearch has not yet been implemented")
		})
	}
	if api.MedcoNodeGetCohortsHandler == nil {
		api.MedcoNodeGetCohortsHandler = medco_node.GetCohortsHandlerFunc(func(params medco_node.GetCohortsParams, principal *models.User) middleware.Responder {
			return middleware.NotImplemented("operation medco_node.GetCohorts has not yet been implemented")
		})
	}
	if api.MedcoNodeGetExploreQueryHandler == nil {
		api.MedcoNodeGetExploreQueryHandler = medco_node.GetExploreQueryHandlerFunc(func(params medco_node.GetExploreQueryParams, principal *models.User) middleware.Responder {
			return middleware.NotImplemented("operation medco_node.GetExploreQuery has not yet been implemented")
		})
	}
	if api.MedcoNetworkGetMetadataHandler == nil {
		api.MedcoNetworkGetMetadataHandler = medco_network.GetMetadataHandlerFunc(func(params medco_network.GetMetadataParams, principal *models.User) middleware.Responder {
			return middleware.NotImplemented("operation medco_network.GetMetadata has not yet been implemented")
		})
	}
	if api.GenomicAnnotationsGetValuesHandler == nil {
		api.GenomicAnnotationsGetValuesHandler = genomic_annotations.GetValuesHandlerFunc(func(params genomic_annotations.GetValuesParams, principal *models.User) middleware.Responder {
			return middleware.NotImplemented("operation genomic_annotations.GetValues has not yet been implemented")
		})
	}
	if api.GenomicAnnotationsGetVariantsHandler == nil {
		api.GenomicAnnotationsGetVariantsHandler = genomic_annotations.GetVariantsHandlerFunc(func(params genomic_annotations.GetVariantsParams, principal *models.User) middleware.Responder {
			return middleware.NotImplemented("operation genomic_annotations.GetVariants has not yet been implemented")
		})
	}
	if api.MedcoNodePostCohortsHandler == nil {
		api.MedcoNodePostCohortsHandler = medco_node.PostCohortsHandlerFunc(func(params medco_node.PostCohortsParams, principal *models.User) middleware.Responder {
			return middleware.NotImplemented("operation medco_node.PostCohorts has not yet been implemented")
		})
	}
	if api.SurvivalAnalysisSurvivalAnalysisHandler == nil {
		api.SurvivalAnalysisSurvivalAnalysisHandler = survival_analysis.SurvivalAnalysisHandlerFunc(func(params survival_analysis.SurvivalAnalysisParams, principal *models.User) middleware.Responder {
			return middleware.NotImplemented("operation survival_analysis.SurvivalAnalysis has not yet been implemented")
		})
	}

	api.PreServerShutdown = func() {}

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
