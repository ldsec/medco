package survivalserver

import (
	"testing"

	"github.com/ldsec/medco/connector/util"
	"github.com/stretchr/testify/assert"
)

func TestGranularity(t *testing.T) {
	_, err := granularity(bigTimePoints, "notATimeResolution")
	assert.Error(t, err, "should be an error")
	day, err := granularity(bigTimePoints, "day")
	if err != nil {
		t.Error(err)
	}
	assert.ElementsMatch(t, bigTimePoints, day)
	week, err := granularity(bigTimePoints, "week")
	if err != nil {
		t.Error(err)
	}
	assert.ElementsMatch(t, bigTimePointsWeek, week)
	month, err := granularity(bigTimePoints, "month")
	if err != nil {
		t.Error(err)
	}
	assert.ElementsMatch(t, bigTimePointsMonth, month)
	year, err := granularity(bigTimePoints, "year")
	if err != nil {
		t.Error(err)
	}
	assert.ElementsMatch(t, bigTimePointsYear, year)
}

var bigTimePointsWeek = util.TimePointsFromTable([][]int{{14, 4, 1}, {78, 0, 1}, {21, 5, 0}, {28, 1, 3}, {50, 3, 0}, {54, 0, 1}, {58, 0, 1}, {62, 3, 0}, {64, 2, 1}, {104, 1, 0}, {117, 1, 0}, {120, 0, 1}, {9, 5, 0}, {18, 1, 0}, {38, 0, 1}, {48, 0, 1}, {52, 4, 1}, {88, 1, 0}, {39, 4, 2}, {81, 1, 0}, {84, 1, 1}, {5, 2, 0}, {19, 3, 0}, {20, 1, 0}, {41, 5, 1}, {68, 1, 0}, {105, 2, 0}, {110, 1, 0}, {36, 1, 1}, {46, 1, 0}, {69, 1, 0}, {138, 0, 1}, {15, 1, 1}, {16, 3, 0}, {30, 3, 0}, {40, 0, 2}, {53, 2, 0}, {90, 1, 0}, {92, 2, 0}, {45, 2, 1}, {57, 1, 0}, {11, 1, 0}, {24, 6, 0}, {35, 3, 2}, {82, 1, 0}, {106, 0, 1}, {116, 0, 1}, {146, 0, 1}, {3, 1, 0}, {27, 3, 2}, {49, 2, 0}, {75, 4, 0}, {99, 2, 0}, {113, 1, 0}, {4, 1, 0}, {10, 2, 0}, {31, 1, 1}, {33, 3, 2}, {37, 0, 1}, {44, 3, 1}, {55, 0, 2}, {118, 0, 1}, {23, 2, 0}, {29, 5, 3}, {32, 3, 3}, {42, 3, 2}, {61, 1, 0}, {76, 0, 1}, {80, 1, 1}, {51, 3, 1}, {1, 1, 0}, {8, 3, 0}, {12, 3, 0}, {13, 2, 0}, {25, 2, 3}, {34, 0, 2}, {47, 1, 0}, {59, 0, 1}, {65, 2, 0}, {66, 2, 1}, {94, 2, 0}, {17, 2, 0}, {26, 8, 1}, {43, 1, 3}, {56, 2, 0}, {79, 1, 1}, {101, 2, 0}, {2, 6, 0}, {22, 1, 0}, {77, 1, 0}, {127, 1, 0}, {145, 0, 1}, {73, 0, 2}})
var bigTimePointsMonth = util.TimePointsFromTable([][]int{{22, 4, 0}, {1, 10, 0}, {6, 16, 4}, {34, 0, 1}, {16, 5, 1}, {3, 10, 0}, {11, 8, 3}, {30, 1, 0}, {5, 10, 0}, {10, 8, 8}, {20, 2, 1}, {21, 2, 0}, {35, 0, 1}, {9, 6, 5}, {12, 8, 2}, {25, 3, 1}, {26, 1, 0}, {19, 3, 3}, {23, 2, 0}, {8, 9, 9}, {13, 8, 4}, {14, 1, 2}, {33, 0, 1}, {4, 10, 2}, {7, 15, 8}, {18, 5, 3}, {24, 2, 0}, {27, 1, 1}, {28, 1, 2}, {2, 7, 0}, {15, 7, 1}})
var bigTimePointsYear = util.TimePointsFromTable([][]int{{1, 121, 42}, {2, 38, 14}, {3, 6, 7}})
