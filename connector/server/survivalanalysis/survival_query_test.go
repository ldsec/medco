package survivalserver

import (
	"testing"

	utilcommon "github.com/ldsec/medco/connector/util/common"
	"github.com/stretchr/testify/assert"
)

func TestExpansion(t *testing.T) {
	res, err := expansion(timePoints, 15, "notATimeResolution")
	assert.Error(t, err)
	res, err = expansion(timePoints, 15, "day")
	if err != nil {
		t.Error(err)
	}
	assert.ElementsMatch(t, timePointsExpanded, res)

}

var timePoints = utilcommon.TimePointsFromTable([][]int{{1, 1, 0}, {5, 1, 0}, {13, 2, 5}})
var timePointsExpanded = utilcommon.TimePointsFromTable([][]int{{0, 0, 0}, {1, 1, 0}, {2, 0, 0}, {3, 0, 0}, {4, 0, 0}, {5, 1, 0}, {6, 0, 0}, {7, 0, 0}, {8, 0, 0}, {9, 0, 0}, {10, 0, 0}, {11, 0, 0}, {12, 0, 0}, {13, 2, 5}, {14, 0, 0}})
