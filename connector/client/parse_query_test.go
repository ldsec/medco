package medcoclient

import (
	"testing"

	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/stretchr/testify/assert"
)

func TestTimeSequenceParsing(t *testing.T) {

	testQuery := "aaa bbb seq--before,first,enddate,last,startdate--beforeorsametime,any,enddate,last,startdate"

	stringWithoutSequence, sequences, err := parseSequences(testQuery)
	assert.NoError(t, err)
	assert.Equal(t, "aaa bbb", stringWithoutSequence)
	assert.NotEmpty(t, sequences)

	assert.Equal(t, *sequences[1].When, models.TimingSequenceInfoWhenBEFOREORSAMETIME)
	assert.Equal(t, *sequences[1].WhichObservationFirst, models.TimingSequenceInfoWhichObservationFirstANY)
	assert.Equal(t, *sequences[1].WhichDateFirst, models.TimingSequenceInfoWhichDateFirstENDDATE)
	assert.Equal(t, *sequences[1].WhichObservationSecond, models.TimingSequenceInfoWhichObservationFirstLAST)
	assert.Equal(t, *sequences[1].WhichDateSecond, models.TimingSequenceInfoWhichDateFirstSTARTDATE)
}
