package survivalclient

import (
	medcoclient "github.com/ldsec/medco-connector/client"
	"github.com/ldsec/medco-connector/restapi/models"
)

const queryType models.ExploreQueryType = models.ExploreQueryTypePatientSet

// GetPatientSetIDs executes the explore query to retrieve the patient set
func GetPatientSetIDs(accessToken string, panels [][]string, panelsIsNot []bool, disableTLSCheck bool) (patientSetIDs map[int]string, err error) {

	clientQuery, err := medcoclient.NewExploreQuery(accessToken, queryType, panels, panelsIsNot, disableTLSCheck)
	if err != nil {
		return
	}
	nodesResult, err := clientQuery.Execute()
	if err != nil {
		return
	}
	patientSetIDs = make(map[int]string)
	for nodeIdx, result := range nodesResult {
		patientSetIDs[nodeIdx] = result.PatientSetID
	}
	return

}
