//go:build integration_test
// +build integration_test

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
	emptyResult, patientsWithoutStartEvent, err := startEvent([]int64{}, []string{"A168", "A125"}, []string{"@"}, true)
	assert.NoError(t, err)
	assert.Empty(t, emptyResult)
	assert.Empty(t, patientsWithoutStartEvent)

	emptyResult, patientsWithoutStartEvent, err = startEvent([]int64{1137, 1138}, []string{}, []string{"@"}, true)
	assert.NoError(t, err)
	assert.Empty(t, emptyResult)
	assert.NotEmpty(t, patientsWithoutStartEvent)

	emptyResult, patientsWithoutStartEvent, err = startEvent([]int64{1137, 1138}, []string{"A168", "A125"}, []string{}, true)
	assert.NoError(t, err)
	assert.Empty(t, emptyResult)
	assert.NotEmpty(t, patientsWithoutStartEvent)

	// test with correct parameters, and an extra patient
	result, patientsWithoutStartEvent, err := startEvent([]int64{1137, 1138, 9999999}, []string{"A168", "A125"}, []string{"@"}, true)
	assert.NoError(t, err)
	expectedFirstTime, err := time.Parse(sqlDateFormat, "1971-04-15")
	expectedSecondTime, err := time.Parse(sqlDateFormat, "1970-03-14")
	assert.NoError(t, err)
	_, isIn := patientsWithoutStartEvent[9999999]
	assert.True(t, isIn)

	firstTime, isIn := result[1137]
	assert.True(t, isIn)
	assert.Equal(t, expectedFirstTime, firstTime)

	secondTime, isIn := result[1138]
	assert.True(t, isIn)
	assert.Equal(t, expectedSecondTime, secondTime)

	// another test with latest instead of earliest
	result, patientsWithoutStartEvent, err = startEvent([]int64{1137, 1138}, []string{"A168", "A125"}, []string{"@"}, false)
	assert.NoError(t, err)
	expectedFirstTime, err = time.Parse(sqlDateFormat, "1972-02-15")
	expectedSecondTime, err = time.Parse(sqlDateFormat, "1971-06-12")
	assert.NoError(t, err)
	assert.Empty(t, patientsWithoutStartEvent)

	firstTime, isIn = result[1137]
	assert.True(t, isIn)
	assert.Equal(t, expectedFirstTime, firstTime)

	secondTime, isIn = result[1138]
	assert.True(t, isIn)
	assert.Equal(t, expectedSecondTime, secondTime)

}

func TestEndEvents(t *testing.T) {
	utilserver.TestI2B2DBConnection(t)

	absoluteEarliest, err := time.Parse(sqlDateFormat, "1970-03-13")
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

	expectedFirstList := createDateListFromString(t, []string{"1971-04-15", "1972-02-15"})

	firstList, isIn := result[1137]
	assert.True(t, isIn)
	assert.ElementsMatch(t, expectedFirstList, firstList)

	expectedSecondList := createDateListFromString(t, []string{"1970-03-14", "1971-06-12"})

	secondList, isIn := result[1138]
	assert.True(t, isIn)
	assert.ElementsMatch(t, expectedSecondList, secondList)

	// expect shorter list if the start date is equal or bigger
	collidingEarliest, err := time.Parse(sqlDateFormat, "1970-03-14")
	assert.NoError(t, err)

	oneCollisionStartEventMap := map[int64]time.Time{
		1137: collidingEarliest,
		1138: collidingEarliest,
	}
	result, err = endEvents(oneCollisionStartEventMap, []string{"A168", "A125"}, []string{"@"})
	assert.NoError(t, err)

	expectedList := createDateListFromString(t, []string{"1971-06-12"})

	list, isIn := result[1138]
	assert.True(t, isIn)
	assert.ElementsMatch(t, expectedList, list)

	// expect empty results
	latest, err := time.Parse(sqlDateFormat, "1972-02-15")
	assert.NoError(t, err)

	latestStartEventMap := map[int64]time.Time{
		1137: latest,
		1138: latest,
	}
	result, err = endEvents(latestStartEventMap, []string{"A168", "A125"}, []string{"@"})
	assert.NoError(t, err)

	_, isIn = result[1138]
	assert.False(t, isIn)

}

func TestCensoringEvent(t *testing.T) {
	utilserver.TestI2B2DBConnection(t)

	absoluteEarliest, err := time.Parse(sqlDateFormat, "1970-03-13")
	assert.NoError(t, err)

	fullStartEventMap := map[int64]time.Time{
		1137: absoluteEarliest,
		1138: absoluteEarliest,
	}

	patientsNoEndEvent := map[int64]struct{}{
		1137: {},
		1138: {},
	}

	// the second argument is a subset of the first argument, an error is expected
	emptyResult, patientWithoutCensoring, err := censoringEvent(map[int64]time.Time{}, patientsNoEndEvent, []string{"A168", "A125"}, []string{"@"})
	assert.Error(t, err)

	emptyResult, patientWithoutCensoring, err = censoringEvent(fullStartEventMap, map[int64]struct{}{}, []string{"A168", "A125"}, []string{"@"})
	assert.NoError(t, err)
	assert.Empty(t, emptyResult)
	assert.Empty(t, patientWithoutCensoring)

	timeStrings := createDateListFromString(t, []string{"1972-02-15", "1971-06-12"})
	expectedCensoring := map[int64]time.Time{
		1137: timeStrings[0],
		1138: timeStrings[1],
	}

	expectedCensoringAuxiliary := func(t *testing.T, patientWithoutCensoring map[int64]struct{}, results map[int64]time.Time) {
		assert.NoError(t, err)
		assert.Empty(t, patientWithoutCensoring)
		firstTime, isIn := results[1137]
		assert.True(t, isIn)
		assert.Equal(t, expectedCensoring[1137], firstTime)
		secondTime, isIn := results[1138]
		assert.True(t, isIn)
		assert.Equal(t, expectedCensoring[1138], secondTime)
	}

	results, patientWithoutCensoring, err := censoringEvent(fullStartEventMap, patientsNoEndEvent, []string{}, []string{"@"})
	expectedCensoringAuxiliary(t, patientWithoutCensoring, results)

	results, patientWithoutCensoring, err = censoringEvent(fullStartEventMap, patientsNoEndEvent, []string{"A168", "A125"}, []string{})
	expectedCensoringAuxiliary(t, patientWithoutCensoring, results)

	results, patientWithoutCensoring, err = censoringEvent(fullStartEventMap, patientsNoEndEvent, []string{"A168", "A125"}, []string{"@"})
	expectedCensoringAuxiliary(t, patientWithoutCensoring, results)

	// put all possible concept and modifier codes, expecting empty results, but no error
	emptyResult, patientWithoutCensoring, err = censoringEvent(fullStartEventMap, patientsNoEndEvent, []string{"A168", "A125", "DEM|SEX:f"}, []string{"@", "126:1", "171:0"})
	assert.NoError(t, err)
	assert.Empty(t, emptyResult)
	assert.Empty(t, patientWithoutCensoring)

	// put start events that do not occur before any other events
	absoluteLatest, err := time.Parse(sqlDateFormat, "1972-02-15")
	assert.NoError(t, err)

	lateStartEventMap := map[int64]time.Time{
		1137: absoluteLatest,
		1138: absoluteLatest,
	}

	emptyResult, patientWithoutCensoring, err = censoringEvent(lateStartEventMap, patientsNoEndEvent, []string{}, []string{})
	assert.NoError(t, err)
	assert.Empty(t, emptyResult)
	_, isIn := patientWithoutCensoring[1137]
	assert.True(t, isIn)
	_, isIn = patientWithoutCensoring[1138]
	assert.True(t, isIn)

}
