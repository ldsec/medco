//go:build unit_test
// +build unit_test

package timepoints

import (
	"testing"
	"time"

	medcomodels "github.com/ldsec/medco/connector/models"
	"github.com/stretchr/testify/assert"
)

func TestPatientAndEndEvents(t *testing.T) {
	someDates := (createDateListFromString(t, []string{
		"1970-09-01",
		"1970-09-02",
		"1970-09-03",
		"1970-09-04",
		"1970-09-05",
		"1970-09-06",
		"1970-09-07"}))

	// test earliest
	withoutEndEvents, relativeTime, err := patientAndEndEvents(map[int64]time.Time{
		0: someDates[0],
		1: someDates[1],
	}, map[int64][]time.Time{
		0: {someDates[2], someDates[3]},
		1: {someDates[4], someDates[5]},
	},
		true,
	)

	assert.NoError(t, err)
	assert.Empty(t, withoutEndEvents)
	count1, isIn := relativeTime[0]
	assert.True(t, isIn)
	assert.Equal(t, int64(2), count1)
	count2, isIn := relativeTime[1]
	assert.True(t, isIn)
	assert.Equal(t, int64(3), count2)

	// test latest
	withoutEndEvents, relativeTime, err = patientAndEndEvents(map[int64]time.Time{
		0: someDates[0],
		1: someDates[1],
	}, map[int64][]time.Time{
		0: {someDates[2], someDates[3]},
		1: {someDates[4], someDates[5]},
	},
		false,
	)

	assert.NoError(t, err)
	assert.Empty(t, withoutEndEvents)
	count1, isIn = relativeTime[0]
	assert.True(t, isIn)
	assert.Equal(t, int64(3), count1)
	count2, isIn = relativeTime[1]
	assert.True(t, isIn)
	assert.Equal(t, int64(4), count2)

	// test patient without end events
	withoutEndEvents, relativeTime, err = patientAndEndEvents(map[int64]time.Time{
		0: someDates[0],
		1: someDates[1],
	}, map[int64][]time.Time{
		0: {someDates[3]},
	},
		false,
	)

	assert.NoError(t, err)
	_, isIn = withoutEndEvents[1]
	assert.True(t, isIn)
	count1, isIn = relativeTime[0]
	assert.True(t, isIn)
	assert.Equal(t, int64(3), count1)
	_, isIn = relativeTime[1]
	assert.False(t, isIn)

	// test wrong data, end event occuring after
	withoutEndEvents, relativeTime, err = patientAndEndEvents(map[int64]time.Time{
		0: someDates[6],
	}, map[int64][]time.Time{
		0: {someDates[0]},
	},
		false,
	)

	assert.Error(t, err)

	// test wrong data, end event is the same

	withoutEndEvents, relativeTime, err = patientAndEndEvents(map[int64]time.Time{
		0: someDates[0],
	}, map[int64][]time.Time{
		0: {someDates[0]},
	},
		false,
	)

	assert.Error(t, err)

	// test wrong data, empty list
	withoutEndEvents, relativeTime, err = patientAndEndEvents(map[int64]time.Time{
		0: someDates[0],
	}, map[int64][]time.Time{
		0: {},
	},
		false,
	)

	assert.Error(t, err)

	// test wrong data, nil
	withoutEndEvents, relativeTime, err = patientAndEndEvents(map[int64]time.Time{
		0: someDates[0],
	}, map[int64][]time.Time{
		0: nil,
	},
		false,
	)

	assert.Error(t, err)

}

func TestPatientAndCensoring(t *testing.T) {
	someDates := (createDateListFromString(t, []string{
		"1970-09-01",
		"1970-09-02",
		"1970-09-03",
		"1970-09-04",
		"1970-09-05",
		"1970-09-06",
		"1970-09-07"}))

	// test full set, the extra data in end events is silently ignored

	relativeTime, err := patientAndCensoring(map[int64]time.Time{
		0: someDates[0],
		1: someDates[1],
	},
		map[int64]struct{}{
			0: {},
			1: {},
		},
		map[int64]time.Time{
			0: someDates[2],
			1: someDates[4],
			2: someDates[5],
		},
	)
	assert.NoError(t, err)

	count1, isIn := relativeTime[0]
	assert.True(t, isIn)
	assert.Equal(t, int64(2), count1)
	count2, isIn := relativeTime[1]
	assert.True(t, isIn)
	assert.Equal(t, int64(3), count2)
	_, isIn = relativeTime[2]
	assert.False(t, isIn)

	// test one patient missing

	relativeTime, err = patientAndCensoring(map[int64]time.Time{
		0: someDates[0],
		1: someDates[1],
	},
		map[int64]struct{}{
			0: {},
		},
		map[int64]time.Time{
			0: someDates[2],
		},
	)
	assert.NoError(t, err)

	count1, isIn = relativeTime[0]
	assert.True(t, isIn)
	assert.Equal(t, int64(2), count1)
	_, isIn = relativeTime[1]
	assert.False(t, isIn)

	// test wrong data, extra data in patients-without-end-data

	relativeTime, err = patientAndCensoring(map[int64]time.Time{
		0: someDates[0],
	},
		map[int64]struct{}{
			0: {},
			1: {},
		},
		map[int64]time.Time{
			0: someDates[2],
			1: someDates[2],
		},
	)
	assert.Error(t, err)

	// test wrong data, censoring date before

	relativeTime, err = patientAndCensoring(map[int64]time.Time{
		0: someDates[4],
	},
		map[int64]struct{}{
			0: {},
		},
		map[int64]time.Time{
			0: someDates[0],
		},
	)
	assert.Error(t, err)

	// test wrong data, censoring same date

	relativeTime, err = patientAndCensoring(map[int64]time.Time{
		0: someDates[0],
	},
		map[int64]struct{}{
			0: {},
		},
		map[int64]time.Time{
			0: someDates[0],
		},
	)
	assert.Error(t, err)

}

func TestCompileTimePoints(t *testing.T) {
	patientsWithEndEvent := map[int64]int64{
		0: 1,
		1: 1,
		2: 3,
		3: 4,
	}

	patientsWithCensoring := map[int64]int64{
		4: 1,
		5: 2,
		6: 3,
		7: 3,
		8: 2,
	}
	expectedEvents := map[int64]*medcomodels.Events{
		1: {
			EventsOfInterest: 2,
			CensoringEvents:  1,
		},
		2: {
			EventsOfInterest: 0,
			CensoringEvents:  2,
		},
		3: {
			EventsOfInterest: 1,
			CensoringEvents:  2,
		},
		4: {
			EventsOfInterest: 1,
			CensoringEvents:  0,
		},
	}

	// test full events
	events, err := compileTimePoints(patientsWithEndEvent, patientsWithCensoring, int64(4))
	assert.NoError(t, err)

	for relativeTime, expectedEvent := range expectedEvents {
		event, isIn := events[relativeTime]
		assert.True(t, isIn, "event aggregates for relative time %d not found", relativeTime)

		assert.Equal(t, expectedEvent, event, "event aggregates for relative time %d are not the same as expected ones", relativeTime)
	}

	// test limit parameter
	expectedEventsLimited := make(map[int64]*medcomodels.Events, 3)
	for key, value := range expectedEvents {
		if key != 4 {
			expectedEventsLimited[key] = value
		}
	}

	events, err = compileTimePoints(patientsWithEndEvent, patientsWithCensoring, int64(3))
	assert.NoError(t, err)
	_, isIn := events[4]
	assert.False(t, isIn)

	for relativeTime, expectedEvent := range expectedEventsLimited {
		event, isIn := events[relativeTime]
		assert.True(t, isIn, "event aggregates for relative time %d not found", relativeTime)

		assert.Equal(t, expectedEvent, event, "event aggregates for relative time %d are not the same as expected ones", relativeTime)
	}

	// test wrong data, bad relative time
	_, err = compileTimePoints(map[int64]int64{0: 1, 1: 0}, patientsWithCensoring, int64(4))
	assert.Error(t, err)
	_, err = compileTimePoints(patientsWithEndEvent, map[int64]int64{6: 1, 7: 0}, int64(4))
	assert.Error(t, err)

	// test wrong data, bad time limit
	_, err = compileTimePoints(patientsWithEndEvent, patientsWithCensoring, int64(0))
	assert.Error(t, err)

}
