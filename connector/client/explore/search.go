package exploreclient

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ldsec/medco/connector/restapi/client"
	"github.com/ldsec/medco/connector/restapi/client/medco_node"
	"github.com/ldsec/medco/connector/restapi/models"
	utilclient "github.com/ldsec/medco/connector/util/client"
	"github.com/sirupsen/logrus"
)

// ExploreSearch is a MedCo client explore search.
type ExploreSearch struct {

	// httpMedCoClient is the HTTP client for the MedCo connectors
	httpMedCoClient *client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

	searchString string
}

// ExploreSearchConcept is a MedCo client explore concept search.
type ExploreSearchConcept struct {

	// httpMedCoClient is the HTTP client for the MedCo connectors
	httpMedCoClient *client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

	conceptPath string

	operation string
}

// ExploreSearchModifier is a MedCo client explore modifier search.
type ExploreSearchModifier struct {

	// httpMedCoClient is the HTTP client for the MedCo connectors
	httpMedCoClient *client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

	modifierPath   string
	appliedPath    string
	appliedConcept string

	operation string
}

// NewExploreSearch creates a new MedCo client explore search.
func NewExploreSearch(authToken, searchString string, disableTLSCheck bool) (es *ExploreSearch, err error) {

	es = &ExploreSearch{
		authToken:    authToken,
		searchString: searchString,
	}

	// retrieve network information
	parsedURL, err := url.Parse(utilclient.MedCoConnectorURL)
	if err != nil {
		logrus.Error("cannot parse MedCo connector URL: ", err)
		return
	}

	transport := httptransport.New(parsedURL.Host, parsedURL.Path, []string{parsedURL.Scheme})
	transport.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: disableTLSCheck}

	// parse network information
	es.httpMedCoClient = client.New(transport, nil)

	return
}

// NewExploreSearchConcept creates a new MedCo client explore concept search.
func NewExploreSearchConcept(authToken, conceptPath, operation string, disableTLSCheck bool) (scc *ExploreSearchConcept, err error) {

	scc = &ExploreSearchConcept{
		authToken:   authToken,
		conceptPath: conceptPath,
		operation:   operation,
	}

	// retrieve network information
	parsedURL, err := url.Parse(utilclient.MedCoConnectorURL)
	if err != nil {
		logrus.Error("cannot parse MedCo connector URL: ", err)
		return
	}

	transport := httptransport.New(parsedURL.Host, parsedURL.Path, []string{parsedURL.Scheme})
	transport.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: disableTLSCheck}

	// parse network information
	scc.httpMedCoClient = client.New(transport, nil)

	return
}

// NewExploreSearchModifier creates a new MedCo client explore modifier search.
func NewExploreSearchModifier(authToken, modifierPath, appliedPath, appliedConcept, operation string, disableTLSCheck bool) (smc *ExploreSearchModifier, err error) {

	smc = &ExploreSearchModifier{
		authToken:      authToken,
		modifierPath:   modifierPath,
		appliedPath:    appliedPath,
		appliedConcept: appliedConcept,
		operation:      operation,
	}

	// retrieve network information
	parsedURL, err := url.Parse(utilclient.MedCoConnectorURL)
	if err != nil {
		logrus.Error("cannot parse MedCo connector URL: ", err)
		return
	}

	transport := httptransport.New(parsedURL.Host, parsedURL.Path, []string{parsedURL.Scheme})
	transport.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: disableTLSCheck}

	// parse network information
	smc.httpMedCoClient = client.New(transport, nil)

	return
}

// Execute executes the MedCo client search.
func (exploreSearch *ExploreSearch) Execute() (*medco_node.ExploreSearchOK, error) {

	params := medco_node.NewExploreSearchParamsWithTimeout(time.Duration(utilclient.SearchTimeoutSeconds) * time.Second)
	params.SearchRequest = &models.ExploreSearch{SearchString: &exploreSearch.searchString}

	response, err := exploreSearch.httpMedCoClient.MedcoNode.ExploreSearch(params, httptransport.BearerToken(exploreSearch.authToken))

	if err != nil {
		logrus.Error("Explore Search error: ", err)
		return nil, err
	}

	return response, nil

}

// Execute executes the MedCo client concept search.
func (exploreSearchConcept *ExploreSearchConcept) Execute() (*medco_node.ExploreSearchConceptOK, error) {

	params := medco_node.NewExploreSearchConceptParamsWithTimeout(time.Duration(utilclient.SearchTimeoutSeconds) * time.Second)
	params.SearchConceptRequest = &models.ExploreSearchConcept{Path: &exploreSearchConcept.conceptPath, Operation: &exploreSearchConcept.operation}

	response, err := exploreSearchConcept.httpMedCoClient.MedcoNode.ExploreSearchConcept(params, httptransport.BearerToken(exploreSearchConcept.authToken))

	if err != nil {
		logrus.Error("Explore Search Concept Children error: ", err)
		return nil, err
	}

	return response, nil

}

// Execute executes the MedCo client modifier search.
func (exploreSearchModifier *ExploreSearchModifier) Execute() (*medco_node.ExploreSearchModifierOK, error) {

	params := medco_node.NewExploreSearchModifierParamsWithTimeout(time.Duration(utilclient.SearchTimeoutSeconds) * time.Second)
	params.SearchModifierRequest = &models.ExploreSearchModifier{
		Path:           &exploreSearchModifier.modifierPath,
		AppliedPath:    &exploreSearchModifier.appliedPath,
		AppliedConcept: &exploreSearchModifier.appliedConcept,
		Operation:      &exploreSearchModifier.operation,
	}

	response, err := exploreSearchModifier.httpMedCoClient.MedcoNode.ExploreSearchModifier(params, httptransport.BearerToken(exploreSearchModifier.authToken))

	if err != nil {
		logrus.Error("Explore Search Modifier Children error: ", err)
		return nil, err
	}

	return response, nil

}
