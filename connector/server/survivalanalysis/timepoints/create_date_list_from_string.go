package timepoints

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// createDateListFromString parses a list of times and returns the date
func createDateListFromString(t *testing.T, dateStrings []string) (timeList []time.Time) {
	timeList = make([]time.Time, len(dateStrings))

	for i, dateString := range dateStrings {
		date, parseErr := time.Parse(sqlDateFormat, dateString)
		assert.NoError(t, parseErr)
		timeList[i] = date
	}
	return
}
