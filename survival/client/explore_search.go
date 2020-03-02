package survivalclient

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"

	httptransport "github.com/go-openapi/runtime/client"
	"github.com/ldsec/medco-connector/restapi/client"
	"github.com/ldsec/medco-connector/restapi/client/medco_network"
	"github.com/ldsec/medco-connector/restapi/client/medco_node"
	"github.com/ldsec/medco-connector/restapi/models"
	utilclient "github.com/ldsec/medco-connector/util/client"
	"github.com/sirupsen/logrus"
)

type ExploreSearch struct {
	httpMedCoClients []*client.MedcoCli
	AuthToken        string
	Path             string
	Type             string
}

func NewExploreSearch(accessToken, path, Type string, disableTLSCheck bool) (search *ExploreSearch, err error) {
	search = &ExploreSearch{
		AuthToken: accessToken,
		Path:      path,
		Type:      Type,
	}

	// retrieve network information
	parsedURL, err := url.Parse(utilclient.MedCoConnectorURL)
	if err != nil {
		logrus.Error("cannot parse MedCo connector URL: ", err)
		return
	}

	transport := httptransport.New(parsedURL.Host, parsedURL.Path, []string{parsedURL.Scheme})
	transport.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: disableTLSCheck}

	getMetadataResp, err := client.New(transport, nil).MedcoNetwork.GetMetadata(
		medco_network.NewGetMetadataParamsWithTimeout(30*time.Second),
		httptransport.BearerToken(accessToken),
	)
	if err != nil {
		logrus.Error("get network metadata request error: ", err)
		return
	}

	// parse network information
	search.httpMedCoClients = make([]*client.MedcoCli, len(getMetadataResp.Payload.Nodes))
	for _, node := range getMetadataResp.Payload.Nodes {
		if search.httpMedCoClients[*node.Index] != nil {
			err = errors.New("duplicated node index in network metadata")
			logrus.Error(err)
			return
		}

		nodeURL, err := url.Parse(node.URL)
		if err != nil {
			logrus.Error("cannot parse MedCo connector URL: ", err)
			return nil, err
		}

		nodeTransport := httptransport.New(nodeURL.Host, nodeURL.Path, []string{nodeURL.Scheme})
		nodeTransport.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: disableTLSCheck}
		search.httpMedCoClients[*node.Index] = client.New(nodeTransport, nil)
	}

	return
}

func (search *ExploreSearch) Execute() (searchResult *ExploreSearchResult, err error) {
	searchResultChan := make(chan []*models.ExploreSearchResultElement)
	searchErrChan := make(chan error)
	//all nodes should have the same ontology
	idxOnlyOne := 0
	go func() {
		routineResponse, routineErr := search.submitToNode(idxOnlyOne)
		if routineErr != nil {
			searchErrChan <- routineErr
		} else {
			searchResultChan <- routineResponse
		}

	}()
	//TODO magic number
	timeout := time.After(time.Duration(300) * time.Second)

	select {
	case err = <-searchErrChan:
		logrus.Error("Search execution error : ", err)
		logrus.Debug(fmt.Sprintf("Search is %+v", search))
		return
	case result := <-searchResultChan:
		searchResult = NewExploreSearchResultFromModels(result)
		return
	case <-timeout:
		err = errors.New("timeout")
		logrus.Error("Search execution error: ", err)
		return

	}
}

func (search *ExploreSearch) submitToNode(idx int) (result []*models.ExploreSearchResultElement, err error) {
	logrus.Debug("Submitting search to node # ", idx)
	//TODO magic number
	params := medco_node.NewExploreSearchParamsWithTimeout(time.Duration(300) * time.Second)
	searchRequest := &models.ExploreSearch{
		Path: &search.Path,
		Type: search.Type,
	}
	params.SetSearchRequest(searchRequest)

	response, err := search.httpMedCoClients[idx].MedcoNode.ExploreSearch(params, httptransport.BearerToken(search.AuthToken))
	if err != nil {
		logrus.Error("explore search  error: ", err)
		return
	}

	return response.Payload.Results, nil

}

type ExploreSearchResult struct {
	Elements []*models.ExploreSearchResultElement
}

func NewExploreSearchResultFromModels(resultElements []*models.ExploreSearchResultElement) *ExploreSearchResult {
	return &ExploreSearchResult{Elements: resultElements}
}
