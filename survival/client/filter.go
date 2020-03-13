package survivalclient

import (
	"fmt"
	"strconv"

	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/sirupsen/logrus"
)

const (
	basePath   string = "/i2b2_SRVA/SurvivalAnalysis/"
	searchType string = ""
)

var fromRequestToI2b2 map[string]string = map[string]string{"day": "Day", "week": "Week", "month": "Month", "year": "Year"}

// GetTimeCodes execute the explore search to find the available time points for a given granularity and retrieves their related integer identifier
func GetTimeCodes(accessToken, granularity string, limit int64, disableTLS bool) (timeCodes map[string]int64, err error) {

	gran, ok := fromRequestToI2b2[granularity]
	if !ok {
		err = fmt.Errorf("Time resolution %s not found in available granularities", granularity)
		return
	}
	path := basePath + gran + "/"
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
