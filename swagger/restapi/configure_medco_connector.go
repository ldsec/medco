// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/lca1/medco-connector/util"
	"log"
	"net/http"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"

	"github.com/lca1/medco-connector/swagger/restapi/operations"
	"github.com/lca1/medco-connector/swagger/restapi/operations/picsure2"
)

//go:generate swagger generate server --target ../../swagger --name MedcoConnector --spec ../swagger.yml

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

	api.JSONConsumer = runtime.JSONConsumer()

	api.JSONProducer = runtime.JSONProducer()

	// Applies when the "Authorization" header is set
	api.PICSURE2ResourceTokenAuth = func(token string, scopes []string) (interface{}, error) {
		log.Printf("Authenticating token %v", token)
		if util.ValidatePICSURE2InternalToken(token) {
			return "", nil
		} else {
			return nil, errors.New(401, "PICSURE2 internal authentication failed (invalid token)")
		}
	}

	// Set your custom authorizer if needed. Default one is security.Authorized()
	// Expected interface runtime.Authorizer
	//
	// Example:
	// api.APIAuthorizer = security.Authorized()
	// todo: use it to auth the user (access to the principal: populate it, set the model for the principal in swagger etc.)

	api.Picsure2GetInfoHandler = picsure2.GetInfoHandlerFunc(func(params picsure2.GetInfoParams, principal interface{}) middleware.Responder {

		return picsure2.NewGetInfoOK().WithPayload(&picsure2.GetInfoOKBody{
			ID: "",
			Name: "MedCo Connector (i2b2: " + util.I2b2CRCCellURL() + ")",
			QueryFormats: []*picsure2.QueryFormatsItems0{
				{
					Name:           "MedCo Query",
					Description:    "Execute a federated MedCo query",
					Examples:       nil,
					Specifications: nil,
				},
			},
		})
	})

	api.Picsure2QueryHandler = picsure2.QueryHandlerFunc(func(params picsure2.QueryParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation picsure2.Query has not yet been implemented")
	})
	api.Picsure2QueryResultHandler = picsure2.QueryResultHandlerFunc(func(params picsure2.QueryResultParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation picsure2.QueryResult has not yet been implemented")
	})
	api.Picsure2QueryStatusHandler = picsure2.QueryStatusHandlerFunc(func(params picsure2.QueryStatusParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation picsure2.QueryStatus has not yet been implemented")
	})
	api.Picsure2QuerySyncHandler = picsure2.QuerySyncHandlerFunc(func(params picsure2.QuerySyncParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation picsure2.QuerySync has not yet been implemented")
	})
	api.Picsure2SearchHandler = picsure2.SearchHandlerFunc(func(params picsure2.SearchParams, principal interface{}) middleware.Responder {
		// todo: better
		//params.Body.Query.Path

		//return models.SearchResult{}
		return picsure2.NewSearchOK().WithPayload(&picsure2.SearchOKBody{Results: nil, SearchQuery: ""})
		//return middleware.NotImplemented("operation picsure2.Search has not yet been implemented")
	})

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
