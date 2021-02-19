// +build unit_test

package medcomodels

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSortTimePoints(t *testing.T) {
	times := timePoints
	sort.Sort(times)
	assert.Equal(t, sortedTimePoints, times)

}

var timePoints = TimePointsFromTable([][]int{{13, 2, 5}, {1, 1, 0}, {5, 1, 0}})
var sortedTimePoints = TimePointsFromTable([][]int{{1, 1, 0}, {5, 1, 0}, {13, 2, 5}})
