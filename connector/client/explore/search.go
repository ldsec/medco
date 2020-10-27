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

// ExploreSearchConcept is a MedCo client explore concept search
type ExploreSearchConcept struct {

	// httpMedCoClients is the HTTP client for the MedCo connectors
	httpMedCoClient *client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

	conceptPath string
}

// ExploreSearchModifier is a MedCo client explore modifier search
type ExploreSearchModifier struct {

	// httpMedCoClients is the HTTP client for the MedCo connectors
	httpMedCoClient *client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

	modifierPath   string
	appliedPath    string
	appliedConcept string
}

// NewExploreSearchConcept creates a new MedCo client explore concept search
func NewExploreSearchConcept(authToken, conceptPath string, disableTLSCheck bool) (sc *ExploreSearchConcept, err error) {

	sc = &ExploreSearchConcept{
		authToken:   authToken,
		conceptPath: conceptPath,
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
	sc.httpMedCoClient = client.New(transport, nil)

	return
}

// NewExploreSearchModifier creates a new MedCo client explore modifier search
func NewExploreSearchModifier(authToken, modifierPath, appliedPath, appliedConcept string, disableTLSCheck bool) (sm *ExploreSearchModifier, err error) {

	sm = &ExploreSearchModifier{
		authToken:      authToken,
		modifierPath:   modifierPath,
		appliedPath:    appliedPath,
		appliedConcept: appliedConcept,
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
	sm.httpMedCoClient = client.New(transport, nil)

	return
}

// Execute executes the MedCo client concept search
func (exploreSearchConcept *ExploreSearchConcept) Execute() (*medco_node.ExploreSearchConceptOK, error) {

	params := medco_node.NewExploreSearchConceptParamsWithTimeout(time.Duration(utilclient.SearchTimeoutSeconds) * time.Second)
	params.SearchConceptRequest = &models.ExploreSearchConcept{Path: &exploreSearchConcept.conceptPath}

	response, err := exploreSearchConcept.httpMedCoClient.MedcoNode.ExploreSearchConcept(params, httptransport.BearerToken(exploreSearchConcept.authToken))

	if err != nil {
		logrus.Error("Explore Search Concept error: ", err)
		return nil, err
	}

	return response, nil

}

// Execute executes the MedCo client modifier search
func (exploreSearchModifier *ExploreSearchModifier) Execute() (*medco_node.ExploreSearchModifierOK, error) {

	params := medco_node.NewExploreSearchModifierParamsWithTimeout(time.Duration(utilclient.SearchTimeoutSeconds) * time.Second)
	params.SearchModifierRequest = &models.ExploreSearchModifier{
		Path:           &exploreSearchModifier.modifierPath,
		AppliedPath:    &exploreSearchModifier.appliedPath,
		AppliedConcept: &exploreSearchModifier.appliedConcept,
	}

	response, err := exploreSearchModifier.httpMedCoClient.MedcoNode.ExploreSearchModifier(params, httptransport.BearerToken(exploreSearchModifier.authToken))

	if err != nil {
		logrus.Error("Explore Search Modifier error: ", err)
		return nil, err
	}

	return response, nil

}
