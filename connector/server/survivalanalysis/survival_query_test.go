//go:build integration_test
// +build integration_test

package survivalserver

import (
	"testing"

	medcomodels "github.com/ldsec/medco/connector/models"
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

var timePoints = medcomodels.TimePointsFromTable([][]int64{{1, 1, 0}, {5, 1, 0}, {13, 2, 5}})
var timePointsExpanded = medcomodels.TimePointsFromTable([][]int64{{0, 0, 0}, {1, 1, 0}, {2, 0, 0}, {3, 0, 0}, {4, 0, 0}, {5, 1, 0}, {6, 0, 0}, {7, 0, 0}, {8, 0, 0}, {9, 0, 0}, {10, 0, 0}, {11, 0, 0}, {12, 0, 0}, {13, 2, 5}, {14, 0, 0}})
