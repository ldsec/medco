package survivalclient

import (
	medcoclient "github.com/ldsec/medco-connector/client"
	"github.com/ldsec/medco-connector/restapi/models"
)

const queryType models.ExploreQueryType = models.ExploreQueryTypePatientSet

//fast solution to avoid cyclic import with medco-connector/client

func GetPatientList(accessToken string, panels [][]string, panelsIsNot []bool, disableTLSCheck bool) (patientSetIDs map[int]string, err error) {

	clientQuery, err := medcoclient.NewExploreQuery(accessToken, queryType, panels, panelsIsNot, disableTLSCheck)
	if err != nil {
		return
	}
	nodesResult, err := clientQuery.Execute()
	if err != nil {
		return
	}

	for nodeIdx, result := range nodesResult {
		patientSetIDs[nodeIdx] = result.PatientSetID
	}
	//TODO global count or persite ???
	return

}
