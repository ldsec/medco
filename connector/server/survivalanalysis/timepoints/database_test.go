package timepoints

import (
	"testing"
	"time"

	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/stretchr/testify/assert"
)

func init() {
	utilserver.SetForTesting()
}

func TestStartEvent(t *testing.T) {
	utilserver.TestI2B2DBConnection(t)

	// test empty, it should not throw an error
	emptyResult, err := startEvent([]int64{}, []string{"A168", "A125"}, []string{"@"}, true)
	assert.NoError(t, err)
	assert.Empty(t, emptyResult)

	emptyResult, err = startEvent([]int64{1137, 1138}, []string{}, []string{"@"}, true)
	assert.NoError(t, err)
	assert.Empty(t, emptyResult)

	emptyResult, err = startEvent([]int64{1137, 1138}, []string{"A168", "A125"}, []string{}, true)
	assert.NoError(t, err)
	assert.Empty(t, emptyResult)

	// test with correct parameters
	result, err := startEvent([]int64{1137, 1138}, []string{"A168", "A125"}, []string{"@"}, true)
	assert.NoError(t, err)
	expectedFirstTime, err := time.Parse(SQLDateFormat, "1971-04-15")
	expectedSecondTime, err := time.Parse(SQLDateFormat, "1970-03-14")
	assert.NoError(t, err)

	firstTime, isIn := result[1137]
	assert.True(t, isIn)
	assert.Equal(t, expectedFirstTime, firstTime)

	secondTime, isIn := result[1138]
	assert.True(t, isIn)
	assert.Equal(t, expectedSecondTime, secondTime)

	// another test with latest instead of earliest
	result, err = startEvent([]int64{1137, 1138}, []string{"A168", "A125"}, []string{"@"}, false)
	assert.NoError(t, err)
	expectedFirstTime, err = time.Parse(SQLDateFormat, "1972-02-15")
	expectedSecondTime, err = time.Parse(SQLDateFormat, "1971-06-12")
	assert.NoError(t, err)

	firstTime, isIn = result[1137]
	assert.True(t, isIn)
	assert.Equal(t, expectedFirstTime, firstTime)

	secondTime, isIn = result[1138]
	assert.True(t, isIn)
	assert.Equal(t, expectedSecondTime, secondTime)

}

func TestEndEvents(t *testing.T) {
	utilserver.TestI2B2DBConnection(t)

	absoluteEarliest, err := time.Parse(SQLDateFormat, "1970-03-13")
	assert.NoError(t, err)

	fullStartEventMap := map[int64]time.Time{
		1137: absoluteEarliest,
		1138: absoluteEarliest,
	}

	// test empty, it should not throw an error
	emptyResult, err := endEvents(map[int64]time.Time{}, []string{"A168", "A125"}, []string{"@"})
	assert.NoError(t, err)
	assert.Empty(t, emptyResult)

	emptyResult, err = endEvents(fullStartEventMap, []string{}, []string{"@"})
	assert.NoError(t, err)
	assert.Empty(t, emptyResult)

	emptyResult, err = endEvents(fullStartEventMap, []string{"A168", "A125"}, []string{})
	assert.NoError(t, err)
	assert.Empty(t, emptyResult)

	// expect all results
	result, err := endEvents(fullStartEventMap, []string{"A168", "A125"}, []string{"@"})
	assert.NoError(t, err)

	expectedFirstList := make([]time.Time, 0)
	for _, dateString := range []string{"1971-04-15", "1972-02-15"} {
		date, parseErr := time.Parse(SQLDateFormat, dateString)
		assert.NoError(t, parseErr)
		expectedFirstList = append(expectedFirstList, date)
	}
	firstList, isIn := result[1137]
	assert.True(t, isIn)
	assert.ElementsMatch(t, expectedFirstList, firstList)

	expectedSecondList := make([]time.Time, 0)
	for _, dateString := range []string{"1970-03-14", "1971-06-12"} {
		date, parseErr := time.Parse(SQLDateFormat, dateString)
		assert.NoError(t, parseErr)
		expectedSecondList = append(expectedSecondList, date)
	}
	secondList, isIn := result[1138]
	assert.True(t, isIn)
	assert.ElementsMatch(t, expectedSecondList, secondList)

	// expect shorter list if the start date is equal or bigger
	collidingEarliest, err := time.Parse(SQLDateFormat, "1970-03-14")
	assert.NoError(t, err)

	oneCollisionStartEventMap := map[int64]time.Time{
		1137: collidingEarliest,
		1138: collidingEarliest,
	}
	result, err = endEvents(oneCollisionStartEventMap, []string{"A168", "A125"}, []string{"@"})
	assert.NoError(t, err)

	expectedList := make([]time.Time, 0)
	for _, dateString := range []string{"1971-06-12"} {
		date, parseErr := time.Parse(SQLDateFormat, dateString)
		assert.NoError(t, parseErr)
		expectedList = append(expectedList, date)
	}
	list, isIn := result[1138]
	assert.True(t, isIn)
	assert.ElementsMatch(t, expectedList, list)

	// expect

}
