package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"strings"
	"testing"

	"github.com/go-openapi/loads"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/onet/v3/log"

	"github.com/ldsec/medco-connector/restapi/server"
	"github.com/ldsec/medco-connector/restapi/server/operations"
	"github.com/ldsec/medco-connector/restapi/server/operations/genomic_annotations"
	"github.com/ldsec/medco-connector/restapi/server/operations/medco_network"
	"github.com/ldsec/medco-connector/restapi/server/operations/medco_node"
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
		{true, "children", ""},
		{true, "children", "/abc/def"},
		{true, "children", "/abc/def/"},
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

func getContextRequest(t *testing.T,
	method, p, str string) (*middleware.Context, *http.Request) {
	spec, api := getApi()
	api.Init()
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
