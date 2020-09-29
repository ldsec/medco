package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/go-openapi/runtime/security"

	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/ldsec/medco/connector/util/server"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/onet/v3/log"

	"github.com/ldsec/medco/connector/restapi/server"
	"github.com/ldsec/medco/connector/restapi/server/operations"
	"github.com/ldsec/medco/connector/restapi/server/operations/genomic_annotations"
	"github.com/ldsec/medco/connector/restapi/server/operations/medco_network"
	"github.com/ldsec/medco/connector/restapi/server/operations/medco_node"
)

func TestNetwork(t *testing.T) {
	ctx, req := getContextRequest(t, "GET", "/network", "")

	ri, rCtx, ok := ctx.RouteInfo(req)
	require.True(t, ok)
	req = rCtx
	err := ctx.BindValidRequest(req, ri, &medco_network.GetMetadataParams{})

	require.NoError(t, err)
}

func TestExploreSearch(t *testing.T) {
	for _, test := range []struct {
		ok      bool
		reqType string
		reqPath string
	}{
		{false, "child", ""},
		{true, "children", "/"},
		{false, "children", ""},
		{true, "children", "/abc/def"},
		{true, "children", "/abc/def/"},
		{true, "children", "/abc/def/asdasdas"},
		{false, "children", "abc/def/"},
		{false, "children", "//def/"},
		{false, "children", "///"},
	} {
		body := fmt.Sprintf(`{"type":"%s", "path":"%s"}`,
			test.reqType, test.reqPath)
		ctx, req := getContextRequest(t, "POST", "/node/explore/search",
			body)

		ri, rCtx, ok := ctx.RouteInfo(req)
		require.True(t, ok)
		req = rCtx

		log.Lvlf2("checking for %t with body: %s", test.ok, body)
		err := ctx.BindValidRequest(req, ri, &medco_node.ExploreSearchParams{})
		require.Equal(t, test.ok, err == nil, "wrong result for %+v: %s",
			test, err)
	}
}

type exploreQueryRequest struct {
	ID    string       `json:"id"`
	Query exploreQuery `json:"query"`
}

type exploreQuery struct {
	Type    string  `json:"type"`
	UserPub string  `json:"userPublicKey"`
	Panels  []panel `json:"panels"`
}

type panel struct {
	Not   bool   `json:"not"`
	Items []item `json:"items"`
}

type item struct {
	QueryTerm string `json:"queryTerm"`
	Operator  string `json:"operator"`
	Value     string `json:"value"`
	Encrypted bool   `json:"encrypted"`
}

type teqTests struct {
	ok    bool
	query exploreQueryRequest
}

func eqValid() exploreQueryRequest {
	return exploreQueryRequest{
		"id", exploreQuery{
			"patient_list", "userPub", []panel{
				{false, []item{
					{"queryTerm", "exists", "", false},
				}},
			},
		}}
}

func TestExploreQuery(t *testing.T) {
	tests := []teqTests{{true, eqValid()}}
	for i := 0; i < 9; i++ {
		tests = append(tests, teqTests{false, eqValid()})
	}
	tests[1].query.ID = "123@"
	tests[2].query.Query.UserPub = "123@"
	tests[3].query.Query.Type = "non-enum"
	tests[4].query.Query.Panels[0].Items[0].Value = "something"
	tests[5].query.Query.Panels[0].Items[0].Operator = "non-enum"
	tests[6].query.Query.Panels[0].Items[0].QueryTerm = "word@"
	tests[7].query.Query.Panels[0].Items[0].QueryTerm = "abc/def"
	tests[8].query.Query.Panels[0].Items[0].QueryTerm = "abc/def/"
	tests[9].query.Query.Panels[0].Items[0].QueryTerm = "/abc/def//"
	for i := 0; i < 3; i++ {
		tests = append(tests, teqTests{true, eqValid()})
	}
	tests[10].query.Query.Panels[0].Items[0].QueryTerm = "word=-word"
	tests[11].query.Query.Panels[0].Items[0].QueryTerm = "/abc/def"
	tests[12].query.Query.Panels[0].Items[0].QueryTerm = "/abc123@/def123@/"

	for _, test := range tests {
		body, err := json.Marshal(test.query)
		require.NoError(t, err)
		ctx, req := getContextRequest(t, "POST", "/node/explore/query",
			string(body))
		// sync is evaluated to true by looking it up in a map. All elements not
		// in the map resolve to false.
		// So "fjdkls" is parsed as a valid false boolean...
		req.URL.RawQuery = "sync=fjklds"
		ri, rCtx, ok := ctx.RouteInfo(req)
		require.True(t, ok)
		req = rCtx

		log.Lvlf2("checking for %t with body: %s", test.ok, body)
		err = ctx.BindValidRequest(req, ri, &medco_node.ExploreQueryParams{})
		require.Equal(t, test.ok, err == nil, "wrong result for %+v: %s",
			test, err)
	}
}

func TestGenomicAnnotations(t *testing.T) {
	for _, test := range []struct {
		ok    bool
		annot string
		value string
		limit int
	}{
		{true, "annotation", "abc", 10},
		{true, "a", "@", 0},
		{false, "", "@", 0},
		{false, "@", "@", 0},
		{false, "a", "@", -1},
	} {
		log.Lvlf2("checking test %+v", test)
		ctx, req := getContextRequest(t, "GET",
			fmt.Sprintf("/genomic-annotations/%s", test.annot), `{}`)
		req.URL.RawQuery = fmt.Sprintf("value=%s&limit=%d", test.value,
			test.limit)
		ri, rCtx, ok := ctx.RouteInfo(req)
		if !ok {
			if test.ok {
				require.Fail(t, "should test OK for OK test")
			}
			continue
		}
		req = rCtx

		log.Lvlf2("url: %s", req.URL)
		err := ctx.BindValidRequest(req, ri, &genomic_annotations.GetValuesParams{})
		require.Equal(t, test.ok, err == nil, "wrong result for %+v: %s",
			test, err)
	}
}

func TestGenomicVariants(t *testing.T) {
	for _, test := range []struct {
		ok       bool
		annot    string
		value    string
		zygosity []string
	}{
		{true, "annot", "@", []string{"heterozygous"}},
		{true, "annot", "@", []string{"homozygous"}},
		{true, "annot", "@", []string{"unknown"}},
		{true, "annot", "@", []string{"heterozygous", "unknown"}},
		{false, "", "@", []string{"unknown"}},
		{false, "annot", "", []string{"unknown"}},
		{false, "annot", "@", []string{"jfklads"}},
		{false, "annot", "@", []string{"unknown", "fasd"}},
	} {
		log.Lvlf2("checking test %+v", test)
		ctx, req := getContextRequest(t, "GET",
			fmt.Sprintf("/genomic-annotations/%s/%s", test.annot, test.value),
			`{}`)
		req.URL.RawQuery = fmt.Sprintf("zygosity=%s",
			strings.Join(test.zygosity, ","))
		ri, rCtx, ok := ctx.RouteInfo(req)
		if !ok {
			if test.ok {
				require.Fail(t, "should test OK for OK test")
			}
			continue
		}
		req = rCtx

		log.Lvlf2("url: %s", req.URL)
		err := ctx.BindValidRequest(req, ri, &genomic_annotations.GetVariantsParams{})
		require.Equal(t, test.ok, err == nil, "wrong result for %+v: %s",
			test, err)
	}
}

func TestAuthorizations(t *testing.T) {
	spec, api := getApi()
	api.Init()
	var authorized bool

	// Shortcut jwt authorization method to not authenticate.
	api.MedcoJwtAuth = func(token string, requiredAuthorizations []string) (principal *models.User, err error) {
		authorized = false
		// Don't check authentication, only authorization.
		rapia := models.RestAPIAuthorization(token)
		principal = &models.User{
			ID:    "userID",
			Token: "some_token",
			Authorizations: &models.UserAuthorizations{
				RestAPI: []models.RestAPIAuthorization{rapia}}}

		// check rest api authorizations
		for _, requiredAuthorization := range requiredAuthorizations {
			err = utilserver.AuthorizeRestAPIEndpoint(principal, models.RestAPIAuthorization(requiredAuthorization))
			if err != nil {
				return
			}
		}
		authorized = true
		return nil, errors.New("only checking authorizations")
	}

	// Replace BearerAuth to only take the proposed authentication.
	api.BearerAuthenticator = func(name string,
		authenticate security.ScopedTokenAuthentication) runtime.Authenticator {
		return security.ScopedAuthenticator(func(r *security.
			ScopedAuthRequest) (bool, interface{}, error) {
			token := r.Request.Header.Get("Authorization")
			p, err := authenticate(token, r.RequiredScopes)
			return true, p, err
		})
	}

	for _, test := range []struct {
		authorized bool
		method     string
		path       string
		restAPI    models.RestAPIAuthorization
	}{
		{true, "", "/network", models.RestAPIAuthorizationMedcoNetwork},
		{false, "", "/network", models.RestAPIAuthorizationMedcoExplore},
		{false, "", "/network", models.RestAPIAuthorizationMedcoGenomicAnnotations},
		{false, "POST", "/node/explore/search",
			models.RestAPIAuthorizationMedcoNetwork},
		{true, "POST", "/node/explore/search",
			models.RestAPIAuthorizationMedcoExplore},
		{false, "POST", "/node/explore/search",
			models.RestAPIAuthorizationMedcoGenomicAnnotations},
		{false, "", "/genomic-annotations/abc",
			models.RestAPIAuthorizationMedcoNetwork},
		{false, "", "/genomic-annotations/abc",
			models.RestAPIAuthorizationMedcoExplore},
		{true, "", "/genomic-annotations/abc",
			models.RestAPIAuthorizationMedcoGenomicAnnotations},
		{false, "", "/genomic-annotations/abc/123",
			models.RestAPIAuthorizationMedcoNetwork},
		{false, "", "/genomic-annotations/abc/123",
			models.RestAPIAuthorizationMedcoExplore},
		{true, "", "/genomic-annotations/abc/123",
			models.RestAPIAuthorizationMedcoGenomicAnnotations},
	} {
		ctx, req := getContextRequestFromApi(t, spec, api,
			test.method, test.path, "")
		req.Header.Set("Authorization", string(test.restAPI))
		route, ok := ctx.LookupRoute(req)
		require.True(t, ok)
		_, _, err := ctx.Authorize(req, route)
		require.Error(t, err)
		require.Equal(t, test.authorized, authorized)
	}
}

func getContextRequest(t *testing.T,
	method, p, str string) (*middleware.Context, *http.Request) {
	spec, api := getApi()
	api.Init()
	return getContextRequestFromApi(t, spec, api, method, p, str)
}

func getContextRequestFromApi(t *testing.T, spec *loads.Document,
	api *operations.MedcoConnectorAPI, method, p,
	str string) (*middleware.Context, *http.Request) {
	ctx := middleware.NewContext(spec, nil,
		middleware.DefaultRouter(spec, api))
	req, err := http.NewRequest(method,
		path.Join(spec.BasePath(), p),
		strings.NewReader(str))
	require.NoError(t, err)
	req.Header.Set("Content-Type", runtime.JSONMime)

	return ctx, req
}

func getApi() (*loads.Document, *operations.MedcoConnectorAPI) {
	swaggerSpec, err := loads.Embedded(server.SwaggerJSON, server.FlatSwaggerJSON)
	if err != nil {
		log.Fatal(err)
	}

	return swaggerSpec, operations.NewMedcoConnectorAPI(swaggerSpec)
}
