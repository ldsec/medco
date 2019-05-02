package medco
import(
	"github.com/go-openapi/runtime/middleware"
	"github.com/lca1/medco-connector/i2b2"
	"github.com/lca1/medco-connector/swagger/models"
	"github.com/lca1/medco-connector/swagger/restapi/operations/picsure2"
	"github.com/lca1/medco-connector/util"
)

// GetInfoHandlerFunc handles /info API endpoint
func GetInfoHandlerFunc(params picsure2.GetInfoParams, principal *models.User) middleware.Responder {

	err := util.Authorize(params.Body.ResourceCredentials, principal)
	if err != nil {
		return picsure2.NewGetInfoDefault(401).WithPayload(&picsure2.GetInfoDefaultBody{
			Message: "Authorization failed: " + err.Error(),
		})
	}

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
}

// SearchHandlerFunc handles /search API endpoint
func SearchHandlerFunc(params picsure2.SearchParams, principal *models.User) middleware.Responder {

	err := util.Authorize(params.Body.ResourceCredentials, principal)
	if err != nil {
		return picsure2.NewSearchDefault(401).WithPayload(&picsure2.SearchDefaultBody{
			Message: "Authorization failed: " + err.Error(),
		})
	}

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
}

// QuerySyncHandlerFunc handles /query/sync API endpoint
func QuerySyncHandlerFunc(params picsure2.QuerySyncParams, principal *models.User) middleware.Responder {

	err := util.Authorize(params.Body.ResourceCredentials, principal)
	if err != nil {
		return picsure2.NewQuerySyncDefault(401).WithPayload(&picsure2.QuerySyncDefaultBody{
			Message: "Authorization failed: " + err.Error(),
		})
	}

	queryResult, err := I2b2MedCoQuery(params.Body.Query.Name, params.Body.Query.I2b2Medco)
	if err != nil {
		return picsure2.NewQuerySyncDefault(500).WithPayload(&picsure2.QuerySyncDefaultBody{
			Message: err.Error(),
		})
	}

	return picsure2.NewQuerySyncOK().WithPayload(queryResult)
}