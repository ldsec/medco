package medcoclient

import (
	"testing"

	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/stretchr/testify/assert"
)

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
