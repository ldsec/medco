package survivalclient

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/sirupsen/logrus"
)

const (
	timePath   string = "/survival_TIME/SurvivalAnalysis/Time/"
	typePath   string = "/survival_TYPE/SurvivalAnalysis/Type/"
	searchType string = ""
)

var fromRequestToI2b2 map[string]string = map[string]string{"day": "Day", "week": "Week", "month": "Month", "season": "Season", "year": "Year", "tumor-progression-free": "TumorProgressionFree"}

// GetTimeCodes execute the explore search to find the available time points for a given granularity and retrieves their related integer identifier
func GetTimeCodes(accessToken, granularity string, limit int64, disableTLS bool) (timeCodes map[string]int64, err error) {

	gran, ok := fromRequestToI2b2[granularity]
	if !ok {
		err = fmt.Errorf("Time resolution %s not found in available granularities", granularity)
		return
	}
	path := timePath + gran + "/"
	logrus.Debugf("survival time :%s  path:%s", granularity, path)
	exploreSearch, err := NewExploreSearch(accessToken, path, searchType, disableTLS)
	if err != nil {
		return
	}
	searchResults, err := exploreSearch.Execute()
	if err != nil {
		return
	}
	var skipped []*models.ExploreSearchResultElement
	var recordSkipped func(*models.ExploreSearchResultElement)
	if logrus.IsLevelEnabled(logrus.DebugLevel) {
		recordSkipped = func(skipResult *models.ExploreSearchResultElement) {
			skipped = append(skipped, skipResult)
		}
	} else {
		recordSkipped = func(skipResult *models.ExploreSearchResultElement) {
		}
	}
	timeCodes = make(map[string]int64)

	for _, result := range searchResults.Elements {
		//the Leaf nature of the concept was not return in getontolgoy children
		if *result.MedcoEncryption.Encrypted /* && *result.Leaf*/ {
			if value, isNotValidInt := strconv.ParseInt(result.Name, 10, 64); isNotValidInt == nil && value < limit {

				timeCodes[result.Name] = *result.MedcoEncryption.ID

			} else {
				recordSkipped(result)
			}
		} else {
			recordSkipped(result)
		}
	}

	if length := len(skipped); length != 0 {
		logrus.Debug(fmt.Sprintf("Skipped %d concepts", length))
		for _, skipConcept := range skipped {
			logrus.Debug(*skipConcept)
		}
		if err != nil {
			return
		}
	}

	return

}

func GetTypeCode(accessToken, survivalType string, disableTLS bool) (typeCode int64, err error) {
	survType, ok := fromRequestToI2b2[survivalType]
	if !ok {
		err = fmt.Errorf("Time resolution %s not found in available granularities", survivalType)
	}
	path := typePath
	//TODO what is searchType ???
	exploreSearch, err := NewExploreSearch(accessToken, path, searchType, disableTLS)
	if err != nil {
		return
	}
	searchResults, err := exploreSearch.Execute()
	if err != nil {
		return
	}
	logrus.Debugf("survival type :%s  path:%s", survivalType, path)
	if len(searchResults.Elements) == 0 {

		err = errors.New("the element was not found")
		return
	}

	for _, resultElm := range searchResults.Elements {
		if resultElm.Name == survType {
			typeCode = *resultElm.MedcoEncryption.ID
			return
		}
	}
	err = errors.New("element not found")
	return
}
