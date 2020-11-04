package survivalserver

import (
	"strconv"
	"strings"

	"github.com/ldsec/medco/connector/wrappers/i2b2"
)

// SubGroupExplore executes an I2B2 Explore query with panelsItemKeys and isNot as definition
func SubGroupExplore(queryName string, subGroupIndex int, panelsItemKeys [][]string, isNot []bool) (int64, []int64, error) {

	patientCount, patientSetID, err := i2b2.ExecutePsmQuery(queryName+"_SUBGROUP_"+strconv.Itoa(subGroupIndex), backSlashFormat(panelsItemKeys), isNot)
	if err != nil {
		return 0, nil, err
	}
	patientIDs, _, err := i2b2.GetPatientSet(patientSetID, false)
	if err != nil {
		return 0, nil, err
	}
	pCount, err := strconv.ParseInt(patientCount, 10, 64)
	if err != nil {
		return 0, nil, err
	}
	pIDs := make([]int64, len(patientIDs))

	for i, pID := range patientIDs {
		id, err := strconv.ParseInt(pID, 10, 64)
		if err != nil {
			return 0, nil, err
		}
		pIDs[i] = id
	}

	return pCount, pIDs, nil

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

func backSlashFormat(panelsItemKeys [][]string) [][]string {
	//table access start with two backslashs
	resPanels := panelsItemKeys
	for i, panel := range panelsItemKeys {
		for j, item := range panel {
			resPanels[i][j] = `\` + strings.Join(strings.Split(item, "/"), `\`)
		}
	}
	return resPanels
}
