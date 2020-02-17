package survival

import (
	survivalserver "github.com/ldsec/medco-connector/survival/server"
	"github.com/ldsec/medco-connector/survival/server/directaccess"
)

//TODO remove this magic
const queryname = `I'm the survival query !`

func Init() {
	survivalserver.ExecCallback = directaccess.QueryTimePoints
}

//like similar to explore_query_logic, but is not necessery to share the patient list,only the count

//this one must be run only by one node

var GlobalTimeCodes []string

/*
func getConceptInteger{

}


func getPatient(){
	queryType:=
	queryString:=
}

//get patientList from pannels

func GetPatientList(token, username, password, queryType, queryString string, disableTLSCheck bool) (patientList []string, err error) {
	// get token
	var accessToken string
	if len(token) > 0 {
		accessToken = token
	} else {
		logrus.Debug("No token provided, requesting token for user ", username, ", disable TLS check: ", disableTLSCheck)
		accessToken, err = utilclient.RetrieveAccessToken(username, password, disableTLSCheck)
		if err != nil {
			return
		}
	}

	// parse query type
	queryTypeParsed := models.ExploreQueryType(queryType)
	err = queryTypeParsed.Validate(nil)
	if err != nil {
		logrus.Error("invalid query type")
		return
	}

	// parse query string
	panelsItemKeys, panelsIsNot, err := medcoclient.ParseQueryString(queryString)
	if err != nil {
		return
	}

	// encrypt the item keys
	encPanelsItemKeys := make([][]string, 0)
	for _, panel := range panelsItemKeys {
		encItemKeys := make([]string, 0)
		for _, itemKey := range panel {
			var encrypted string
			encrypted, err = unlynx.EncryptWithCothorityKey(itemKey)
			if err != nil {
				return
			}
			encItemKeys = append(encItemKeys, encrypted)
		}
		encPanelsItemKeys = append(encPanelsItemKeys, encItemKeys)
	}

	// execute query
	clientQuery, err := medcoclient.NewExploreQuery(accessToken, queryTypeParsed, encPanelsItemKeys, panelsIsNot, disableTLSCheck)
	if err != nil {
		return
	}
}
*/

//get conceptList for time ID
