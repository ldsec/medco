package survivalserver

import (
	"strconv"

	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/ldsec/medco/connector/wrappers/i2b2"
)

// SubGroupExplore executes an I2B2 Explore query with panels
func SubGroupExplore(queryName string, subGroupIndex int, panels []*models.Panel, groupTiming models.Timing) ([]int64, error) {

	_, patientSetID, err := i2b2.ExecutePsmQuery(queryName+"_SUBGROUP_"+strconv.Itoa(subGroupIndex), panels, groupTiming)
	if err != nil {
		return nil, err
	}
	patientIDs, _, err := i2b2.GetPatientSet(patientSetID, false)
	if err != nil {
		return nil, err
	}

	pIDs := make([]int64, len(patientIDs))

	for i, pID := range patientIDs {
		id, err := strconv.ParseInt(pID, 10, 64)
		if err != nil {
			return nil, err
		}
		pIDs[i] = id
	}

	return pIDs, nil

}

// intersect returns the intersection of two sets of int64
func intersect(x []int64, y []int64) []int64 {
	set := make(map[int64]struct{})
	for _, elm := range x {
		set[elm] = struct{}{}
	}

	result := make([]int64, 0)
	for _, elm := range y {
		if _, ok := set[elm]; ok {
			result = append(result, elm)
		}
	}

	return result
}
