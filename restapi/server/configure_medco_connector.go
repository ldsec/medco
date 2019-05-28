// This file is safe to edit. Once it exists it will not be overwritten

package server

import (
	"crypto/tls"
	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/runtime/security"
	"github.com/lca1/medco-connector/medco"
	"github.com/lca1/medco-connector/restapi/models"
	"github.com/lca1/medco-connector/restapi/server/operations"
	"github.com/lca1/medco-connector/restapi/server/operations/picsure2"
	"github.com/sirupsen/logrus"
	"net/http"
)

//go:generate swagger generate server --target ../../swagger --name MedcoConnector --spec ../swagger.yml

func configureFlags(api *operations.MedcoConnectorAPI) {
	// api.CommandLineOptionsGroups = []swag.CommandLineOptionsGroup{ ... }
}

func configureAPI(api *operations.MedcoConnectorAPI) http.Handler {
	api.ServeError = errors.ServeError
	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()
	api.Logger = logrus.Printf

	// override goswagger authentication methods, to be handled manually in requests
	// api.MedCoTokenAuth = ...
	// api.APIAuthorizer = security.Authorized()
	api.BearerAuthenticator = func(name string, authenticate security.ScopedTokenAuthentication) runtime.Authenticator {
		return security.ScopedAuthenticator(func(r *security.ScopedAuthRequest) (bool, interface{}, error) {
			return true, &models.User{}, nil
		})
	}

	// /medco/picsure2/info
	api.Picsure2GetInfoHandler = picsure2.GetInfoHandlerFunc(medco.GetInfoHandlerFunc)

	// /medco/picsure2/query
	api.Picsure2QueryHandler = picsure2.QueryHandlerFunc(func(params picsure2.QueryParams, principal *models.User) middleware.Responder {
		return middleware.NotImplemented("operation picsure2.Query has not yet been implemented")
	})

	// /medco/picsure2/query/{id}/result
	api.Picsure2QueryResultHandler = picsure2.QueryResultHandlerFunc(func(params picsure2.QueryResultParams, principal *models.User) middleware.Responder {
		return middleware.NotImplemented("operation picsure2.QueryResult has not yet been implemented")
	})

	// /medco/picsure2/query/{id}/status
	api.Picsure2QueryStatusHandler = picsure2.QueryStatusHandlerFunc(func(params picsure2.QueryStatusParams, principal *models.User) middleware.Responder {
		return middleware.NotImplemented("operation picsure2.QueryStatus has not yet been implemented")
	})

	// /medco/picsure2/query/sync
	api.Picsure2QuerySyncHandler = picsure2.QuerySyncHandlerFunc(medco.QuerySyncHandlerFunc)

	// /medco/picsure2/search
	api.Picsure2SearchHandler = picsure2.SearchHandlerFunc(medco.SearchHandlerFunc)

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
