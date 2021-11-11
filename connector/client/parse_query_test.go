// +build unit_test

package medcoclient

import (
	"testing"

	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/stretchr/testify/assert"
)

func TestQueryParsing(t *testing.T) {

	// correct queries
	testQuery0 := "a OR NOT b WITH c THEN d"
	testQuery1 := "a OR NOT b WITH c AND d"
	testQuery2 := "a AND b WITH c THEN d"

	// incorrect queries
	testQuery3 := "a AND b WITH c THEN d WITH e"
	testQuery4 := "a AND b WITH c AND d"
	testQuery5 := "a THEN b WITH c THEN d"
	testQuery6 := "a THEN b AND c"

	selection, sequential, err := ParseQueryString(testQuery0)
	assert.NoError(t, err)
	assert.ElementsMatch(t, selection, []string{"a OR NOT b"})
	assert.ElementsMatch(t, sequential, []string{"c", "d"})

	selection, sequential, err = ParseQueryString(testQuery1)
	assert.NoError(t, err)
	assert.ElementsMatch(t, selection, []string{"c", "d"})
	assert.ElementsMatch(t, sequential, []string{"a OR NOT b"})

	selection, sequential, err = ParseQueryString(testQuery2)
	assert.NoError(t, err)
	assert.ElementsMatch(t, selection, []string{"a", "b"})
	assert.ElementsMatch(t, sequential, []string{"c", "d"})

	_, _, err = ParseQueryString(testQuery3)
	assert.Error(t, err)

	_, _, err = ParseQueryString(testQuery4)
	assert.Error(t, err)

	_, _, err = ParseQueryString(testQuery5)
	assert.Error(t, err)

	_, _, err = ParseQueryString(testQuery6)
	assert.Error(t, err)

}
func TestTimeSequenceParsing(t *testing.T) {

	testQuery0 := "thisissurelywrong,first,enddate,last,startdate:beforeorsametime,any,enddate,last,startdate"

	testQuery1 := "before,first,enddate,last,startdate:beforeorsametime,any,enddate,last,startdate"

	sequences, err := ParseSequences(testQuery0)
	assert.Error(t, err)

	sequences, err = ParseSequences(testQuery1)
	assert.NoError(t, err)
	assert.NotEmpty(t, sequences)

	assert.Equal(t, *sequences[1].When, models.TimingSequenceInfoWhenLESSEQUAL)
	assert.Equal(t, *sequences[1].WhichObservationFirst, models.TimingSequenceInfoWhichObservationFirstANY)
	assert.Equal(t, *sequences[1].WhichDateFirst, models.TimingSequenceInfoWhichDateFirstENDDATE)
	assert.Equal(t, *sequences[1].WhichObservationSecond, models.TimingSequenceInfoWhichObservationFirstLAST)
	assert.Equal(t, *sequences[1].WhichDateSecond, models.TimingSequenceInfoWhichDateFirstSTARTDATE)
}
