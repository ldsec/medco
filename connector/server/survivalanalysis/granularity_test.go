//go:build integration_test
// +build integration_test

package survivalserver

import (
	"testing"

	medcomodels "github.com/ldsec/medco/connector/models"

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

var bigTimePointsWeek = medcomodels.TimePointsFromTable([][]int64{{14, 4, 1}, {78, 0, 1}, {21, 5, 0}, {28, 1, 3}, {50, 3, 0}, {54, 0, 1}, {58, 0, 1}, {62, 3, 0}, {64, 2, 1}, {104, 1, 0}, {117, 1, 0}, {120, 0, 1}, {9, 5, 0}, {18, 1, 0}, {38, 0, 1}, {48, 0, 1}, {52, 4, 1}, {88, 1, 0}, {39, 4, 2}, {81, 1, 0}, {84, 1, 1}, {5, 2, 0}, {19, 3, 0}, {20, 1, 0}, {41, 5, 1}, {68, 1, 0}, {105, 2, 0}, {110, 1, 0}, {36, 1, 1}, {46, 1, 0}, {69, 1, 0}, {138, 0, 1}, {15, 1, 1}, {16, 3, 0}, {30, 3, 0}, {40, 0, 2}, {53, 2, 0}, {90, 1, 0}, {92, 2, 0}, {45, 2, 1}, {57, 1, 0}, {11, 1, 0}, {24, 6, 0}, {35, 3, 2}, {82, 1, 0}, {106, 0, 1}, {116, 0, 1}, {146, 0, 1}, {3, 1, 0}, {27, 3, 2}, {49, 2, 0}, {75, 4, 0}, {99, 2, 0}, {113, 1, 0}, {4, 1, 0}, {10, 2, 0}, {31, 1, 1}, {33, 3, 2}, {37, 0, 1}, {44, 3, 1}, {55, 0, 2}, {118, 0, 1}, {23, 2, 0}, {29, 5, 3}, {32, 3, 3}, {42, 3, 2}, {61, 1, 0}, {76, 0, 1}, {80, 1, 1}, {51, 3, 1}, {1, 1, 0}, {8, 3, 0}, {12, 3, 0}, {13, 2, 0}, {25, 2, 3}, {34, 0, 2}, {47, 1, 0}, {59, 0, 1}, {65, 2, 0}, {66, 2, 1}, {94, 2, 0}, {17, 2, 0}, {26, 8, 1}, {43, 1, 3}, {56, 2, 0}, {79, 1, 1}, {101, 2, 0}, {2, 6, 0}, {22, 1, 0}, {77, 1, 0}, {127, 1, 0}, {145, 0, 1}, {73, 0, 2}})
var bigTimePointsMonth = medcomodels.TimePointsFromTable([][]int64{{22, 4, 0}, {1, 10, 0}, {6, 16, 4}, {34, 0, 1}, {16, 5, 1}, {3, 10, 0}, {11, 8, 3}, {30, 1, 0}, {5, 10, 0}, {10, 8, 8}, {20, 2, 1}, {21, 2, 0}, {35, 0, 1}, {9, 6, 5}, {12, 8, 2}, {25, 3, 1}, {26, 1, 0}, {19, 3, 3}, {23, 2, 0}, {8, 9, 9}, {13, 8, 4}, {14, 1, 2}, {33, 0, 1}, {4, 10, 2}, {7, 15, 8}, {18, 5, 3}, {24, 2, 0}, {27, 1, 1}, {28, 1, 2}, {2, 7, 0}, {15, 7, 1}})
var bigTimePointsYear = medcomodels.TimePointsFromTable([][]int64{{1, 121, 42}, {2, 38, 14}, {3, 6, 7}})
var bigTimePoints = medcomodels.TimePointsFromTable([][]int64{{5, 1, 0}, {11, 3, 0}, {12, 1, 0}, {13, 2, 0}, {15, 1, 0}, {26, 1, 0}, {30, 1, 0}, {31, 1, 0}, {53, 2, 0}, {54, 1, 0}, {59, 1, 0}, {60, 2, 0}, {61, 1, 0}, {62, 1, 0}, {65, 2, 0}, {71, 1, 0}, {79, 1, 0}, {81, 2, 0}, {88, 2, 0}, {92, 1, 1}, {93, 1, 0}, {95, 2, 0}, {105, 1, 1}, {107, 2, 0}, {110, 1, 0}, {116, 1, 0}, {118, 1, 0}, {122, 1, 0}, {131, 1, 0}, {132, 2, 0}, {135, 1, 0}, {142, 1, 0}, {144, 1, 0}, {145, 2, 0}, {147, 1, 0}, {153, 1, 0}, {156, 2, 0}, {163, 3, 0}, {166, 2, 0}, {167, 1, 0}, {170, 1, 0}, {173, 0, 1}, {174, 0, 1}, {175, 1, 1}, {176, 1, 0}, {177, 1, 1}, {179, 2, 0}, {180, 1, 0}, {181, 2, 0}, {182, 1, 0}, {183, 1, 0}, {185, 0, 1}, {186, 1, 0}, {188, 0, 1}, {189, 1, 0}, {191, 0, 1}, {192, 0, 1}, {194, 1, 0}, {196, 0, 1}, {197, 1, 1}, {199, 1, 0}, {201, 2, 0}, {202, 1, 1}, {203, 0, 1}, {207, 1, 0}, {208, 1, 0}, {210, 1, 0}, {211, 0, 1}, {212, 1, 0}, {218, 1, 0}, {221, 0, 1}, {222, 1, 1}, {223, 1, 0}, {224, 0, 1}, {225, 0, 2}, {226, 1, 0}, {229, 1, 0}, {230, 1, 0}, {235, 0, 1}, {237, 0, 1}, {239, 2, 0}, {240, 0, 1}, {243, 0, 1}, {245, 1, 0}, {246, 1, 0}, {252, 0, 1}, {259, 0, 1}, {266, 0, 1}, {267, 1, 0}, {268, 1, 0}, {269, 1, 1}, {270, 1, 0}, {272, 0, 1}, {276, 0, 1}, {279, 0, 1}, {283, 1, 0}, {284, 1, 1}, {285, 2, 0}, {286, 1, 0}, {288, 1, 0}, {291, 1, 0}, {292, 0, 2}, {293, 1, 0}, {296, 0, 1}, {300, 0, 1}, {301, 1, 1}, {303, 1, 1}, {305, 1, 0}, {306, 1, 0}, {310, 2, 0}, {315, 0, 1}, {320, 1, 0}, {329, 1, 0}, {332, 0, 1}, {337, 1, 0}, {340, 1, 0}, {345, 1, 0}, {348, 1, 0}, {350, 1, 0}, {351, 1, 0}, {353, 2, 0}, {356, 0, 1}, {361, 1, 0}, {363, 2, 0}, {364, 1, 1}, {371, 2, 0}, {376, 0, 1}, {382, 0, 1}, {384, 0, 1}, {387, 1, 0}, {390, 1, 0}, {394, 1, 0}, {404, 0, 1}, {413, 0, 1}, {426, 1, 0}, {428, 1, 0}, {429, 1, 0}, {433, 1, 0}, {442, 1, 0}, {444, 1, 1}, {450, 1, 0}, {455, 1, 0}, {457, 1, 0}, {458, 0, 1}, {460, 1, 0}, {473, 1, 0}, {477, 1, 0}, {511, 0, 2}, {519, 1, 0}, {520, 1, 0}, {524, 2, 0}, {529, 0, 1}, {533, 1, 0}, {543, 0, 1}, {550, 1, 0}, {551, 0, 1}, {558, 1, 0}, {559, 0, 1}, {567, 1, 0}, {574, 1, 0}, {583, 1, 0}, {588, 0, 1}, {613, 1, 0}, {624, 1, 0}, {641, 1, 0}, {643, 1, 0}, {654, 1, 0}, {655, 1, 0}, {687, 1, 0}, {689, 1, 0}, {705, 1, 0}, {707, 1, 0}, {728, 1, 0}, {731, 1, 0}, {735, 1, 0}, {740, 0, 1}, {765, 1, 0}, {791, 1, 0}, {806, 0, 1}, {814, 1, 0}, {821, 0, 1}, {840, 0, 1}, {883, 1, 0}, {965, 0, 1}, {1010, 0, 1}, {1022, 0, 1}})
