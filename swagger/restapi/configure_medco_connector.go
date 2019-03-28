// This file is safe to edit. Once it exists it will not be overwritten

package restapi

import (
	"crypto/tls"
	"github.com/lca1/medco-connector/i2b2"
	"github.com/lca1/medco-connector/medco"
	"github.com/lca1/medco-connector/util"
	"github.com/sirupsen/logrus"
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

	api.ServeError = errors.ServeError
	api.JSONConsumer = runtime.JSONConsumer()
	api.JSONProducer = runtime.JSONProducer()
	api.Logger = logrus.Printf

	// Applies when the "Authorization" header is set ; enforce layer of security to talk
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
	// todo: and redo with having a model for the principal

	// /medco/picsure2/info
	api.Picsure2GetInfoHandler = picsure2.GetInfoHandlerFunc(func(params picsure2.GetInfoParams, principal interface{}) middleware.Responder {

		return picsure2.NewGetInfoOK().WithPayload(&picsure2.GetInfoOKBody{
			ID: "",
			Name: "MedCo Connector (i2b2: " + util.I2b2HiveURL + ")",
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

	// /medco/picsure2/query
	api.Picsure2QueryHandler = picsure2.QueryHandlerFunc(func(params picsure2.QueryParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation picsure2.Query has not yet been implemented")
	})

	// /medco/picsure2/query/{id}/result
	api.Picsure2QueryResultHandler = picsure2.QueryResultHandlerFunc(func(params picsure2.QueryResultParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation picsure2.QueryResult has not yet been implemented")
	})

	// /medco/picsure2/query/{id}/status
	api.Picsure2QueryStatusHandler = picsure2.QueryStatusHandlerFunc(func(params picsure2.QueryStatusParams, principal interface{}) middleware.Responder {
		return middleware.NotImplemented("operation picsure2.QueryStatus has not yet been implemented")
	})

	// /medco/picsure2/query/sync
	api.Picsure2QuerySyncHandler = picsure2.QuerySyncHandlerFunc(func(params picsure2.QuerySyncParams, principal interface{}) middleware.Responder {

		queryResult, err := medco.I2b2MedCoQuery(params.Body.Query.I2b2Medco)
		if err != nil {
			return picsure2.NewQuerySyncDefault(500).WithPayload(&picsure2.QuerySyncDefaultBody{
				Message: err.Error(),
			})
		}

		return picsure2.NewQuerySyncOK().WithPayload(queryResult)
	})

	// /medco/picsure2/search
	api.Picsure2SearchHandler = picsure2.SearchHandlerFunc(func(params picsure2.SearchParams, principal interface{}) middleware.Responder {

		searchResult, err := i2b2.GetOntologyChildren(params.Body.Query.Path)
		if err != nil {
			return picsure2.NewSearchDefault(500).WithPayload(&picsure2.SearchDefaultBody{
				Message: err.Error(),
			})
		}

		return picsure2.NewSearchOK().WithPayload(&picsure2.SearchOKBody{
			Results: searchResult,
			SearchQuery: params.Body.Query.Path,
		})
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
