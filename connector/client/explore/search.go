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

// ExploreSearchConceptChildren is a MedCo client explore concept children search
type ExploreSearchConceptChildren struct {

	// httpMedCoClients is the HTTP client for the MedCo connectors
	httpMedCoClient *client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

	conceptPath string
}

// ExploreSearchModifierChildren is a MedCo client explore modifier children search
type ExploreSearchModifierChildren struct {

	// httpMedCoClients is the HTTP client for the MedCo connectors
	httpMedCoClient *client.MedcoCli
	// authToken is the OIDC authentication JWT
	authToken string

	modifierPath   string
	appliedPath    string
	appliedConcept string
}

// ExploreSearchConceptInfo is a MedCo client explore concept info search
type ExploreSearchConceptInfo ExploreSearchConceptChildren

// ExploreSearchModifierInfo is a MedCo client explore modifier info search
type ExploreSearchModifierInfo ExploreSearchModifierChildren

// NewExploreSearchConceptChildren creates a new MedCo client explore concept children search
func NewExploreSearchConceptChildren(authToken, conceptPath string, disableTLSCheck bool) (scc *ExploreSearchConceptChildren, err error) {

	scc = &ExploreSearchConceptChildren{
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
	scc.httpMedCoClient = client.New(transport, nil)

	return
}

// NewExploreSearchModifierChildren creates a new MedCo client explore modifier children search
func NewExploreSearchModifierChildren(authToken, modifierPath, appliedPath, appliedConcept string, disableTLSCheck bool) (smc *ExploreSearchModifierChildren, err error) {

	smc = &ExploreSearchModifierChildren{
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
	smc.httpMedCoClient = client.New(transport, nil)

	return
}

// NewExploreSearchConceptInfo creates a new MedCo client explore concept info search
func NewExploreSearchConceptInfo(authToken, conceptPath string, disableTLSCheck bool) (sci *ExploreSearchConceptInfo, err error) {

	scc, err := NewExploreSearchConceptChildren(authToken, conceptPath, disableTLSCheck)
	sci = (*ExploreSearchConceptInfo)(scc)
	return

}

// NewExploreSearchModifierInfo creates a new MedCo client explore modifier info search
func NewExploreSearchModifierInfo(authToken, modifierPath, appliedPath, appliedConcept string, disableTLSCheck bool) (smi *ExploreSearchModifierInfo, err error) {

	smc, err := NewExploreSearchModifierChildren(authToken, modifierPath, appliedPath, "", disableTLSCheck)
	smi = (*ExploreSearchModifierInfo)(smc)
	return

}

// Execute executes the MedCo client concept children search
func (exploreSearchConceptChildren *ExploreSearchConceptChildren) Execute() (*medco_node.ExploreSearchConceptChildrenOK, error) {

	params := medco_node.NewExploreSearchConceptChildrenParamsWithTimeout(time.Duration(utilclient.SearchTimeoutSeconds) * time.Second)
	params.SearchConceptChildrenRequest = &models.ExploreSearchConceptChildren{Path: &exploreSearchConceptChildren.conceptPath}

	response, err := exploreSearchConceptChildren.httpMedCoClient.MedcoNode.ExploreSearchConceptChildren(params, httptransport.BearerToken(exploreSearchConceptChildren.authToken))

	if err != nil {
		logrus.Error("Explore Search Concept Children error: ", err)
		return nil, err
	}

	return response, nil

}

// Execute executes the MedCo client modifier children search
func (exploreSearchModifierChildren *ExploreSearchModifierChildren) Execute() (*medco_node.ExploreSearchModifierChildrenOK, error) {

	params := medco_node.NewExploreSearchModifierChildrenParamsWithTimeout(time.Duration(utilclient.SearchTimeoutSeconds) * time.Second)
	params.SearchModifierChildrenRequest = &models.ExploreSearchModifierChildren{
		Path:           &exploreSearchModifierChildren.modifierPath,
		AppliedPath:    &exploreSearchModifierChildren.appliedPath,
		AppliedConcept: &exploreSearchModifierChildren.appliedConcept,
	}

	response, err := exploreSearchModifierChildren.httpMedCoClient.MedcoNode.ExploreSearchModifierChildren(params, httptransport.BearerToken(exploreSearchModifierChildren.authToken))

	if err != nil {
		logrus.Error("Explore Search Modifier Children error: ", err)
		return nil, err
	}

	return response, nil

}

// Execute executes the MedCo client concept info search
func (exploreSearchConceptInfo *ExploreSearchConceptInfo) Execute() (*medco_node.ExploreSearchConceptInfoOK, error) {

	params := medco_node.NewExploreSearchConceptInfoParamsWithTimeout(time.Duration(utilclient.SearchTimeoutSeconds) * time.Second)
	params.SearchConceptInfoRequest = &models.ExploreSearchConceptInfo{Path: &exploreSearchConceptInfo.conceptPath}

	response, err := exploreSearchConceptInfo.httpMedCoClient.MedcoNode.ExploreSearchConceptInfo(params, httptransport.BearerToken(exploreSearchConceptInfo.authToken))

	if err != nil {
		logrus.Error("Explore Search Concept Info error: ", err)
		return nil, err
	}

	return response, nil

}

// Execute executes the MedCo client modifier info search
func (exploreSearchModifierInfo *ExploreSearchModifierInfo) Execute() (*medco_node.ExploreSearchModifierInfoOK, error) {

	params := medco_node.NewExploreSearchModifierInfoParamsWithTimeout(time.Duration(utilclient.SearchTimeoutSeconds) * time.Second)
	params.SearchModifierInfoRequest = &models.ExploreSearchModifierInfo{
		Path:        &exploreSearchModifierInfo.modifierPath,
		AppliedPath: &exploreSearchModifierInfo.appliedPath,
	}

	response, err := exploreSearchModifierInfo.httpMedCoClient.MedcoNode.ExploreSearchModifierInfo(params, httptransport.BearerToken(exploreSearchModifierInfo.authToken))

	if err != nil {
		logrus.Error("Explore Search Modifier Children error: ", err)
		return nil, err
	}

	return response, nil

}
