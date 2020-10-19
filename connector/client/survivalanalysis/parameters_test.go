package survivalclient

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testParamFile string = "testParams.yaml"

func TestNewParametersFromFile(t *testing.T) {
	file, err := os.Open(testParamFile)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("test parameter files does not exiss")
		} else {
			t.Error(err)
		}
	}
	err = file.Close()
	assert.NoError(t, err)

	_, err = NewParametersFromFile("./thisFileDoesNotExist")
	assert.Error(t, err)

	params, err := NewParametersFromFile(testParamFile)
	assert.NoError(t, err)
	assert.Equal(t, parameters, params)

}

var parameters = &Parameters{
	TimeResolution:       "day",
	TimeLimit:            19,
	CohortName:           "anyName",
	StartConceptPath:     "/any/start/path/",
	StartConceptModifier: "anyStartMCode",
	EndConceptPath:       "/any/end/path/",
	EndConceptModifier:   "anyEndMCode",
	Cohorts: []*struct {
		Panels []*struct {
			Not   bool     "yaml:\"not\""
			Paths []string "yaml:\"paths\""
		} "yaml:\"panels\""
	}{
		{
			Panels: []*struct {
				Not   bool     "yaml:\"not\""
				Paths []string "yaml:\"paths\""
			}{
				{
					Not:   false,
					Paths: []string{"/path/1/", "/path/2/"},
				},
				{
					Not:   true,
					Paths: []string{"/path/3/"},
				},
			},
		},
		{
			Panels: []*struct {
				Not   bool     "yaml:\"not\""
				Paths []string "yaml:\"paths\""
			}{
				{
					Not:   false,
					Paths: []string{"/path/4/"},
				},
			},
		},
	},
}
