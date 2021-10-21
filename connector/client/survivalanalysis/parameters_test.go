// +build unit_test

package survivalclient

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testParamFile string = "../../../test/survival_unit_test_parameters.yaml"

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
	TimeResolution:   "day",
	TimeLimit:        19,
	CohortName:       "anyName",
	StartConceptPath: "/any/start/path/",
	StartModifier: &modifier{
		ModifierKey: "/any/start/modifier/key/",
		AppliedPath: "/any/start/path/%",
	},
	StartsWhen:     "earliest",
	EndConceptPath: "/any/end/path/",
	EndModifier: &modifier{
		ModifierKey: "/any/end/modifier/key/",
		AppliedPath: "/any/end/path/%",
	},
	EndsWhen: "earliest",
	SubGroups: []*subGroup{
		{
			GroupName: "AAA",
			Panels: []*panel{
				{
					PanelTiming: "any",
					Not:         false,
					Items: []*item{
						{
							Path: "/path/1/",
						},
						{
							Path: "/path/2/",
							Modifier: &modifier{
								ModifierKey: "/key1/",
								AppliedPath: "/appliedpath1/",
							},
						},
					},
				},
				{
					PanelTiming: "sameinstancenum",
					Not:         true,
					Items: []*item{
						{
							Path: "/path/3/",
						},
					},
				},
			},
			SequenceOfEvents: []*sequenceElement{
				{
					WhichDateFirst:         "startdate",
					WhichDateSecond:        "enddate",
					WhichObservationFirst:  "any",
					WhichObservationSecond: "last",
					When:                   "sametime",
					Spans: []*timeSpan{
						{Operator: "before",
							Value: 34,
							Units: "years"},
						{Operator: "sametime",
							Value: 21,
							Units: "days"},
					},
				},
			},
		},
		{
			GroupName: "BBB",
			Panels: []*panel{
				{
					Not:         false,
					PanelTiming: "any",
					Items: []*item{
						{
							Path: "/path/4/",
						},
					},
				},
				{
					Not:         true,
					PanelTiming: "sameinstancenum",
					Items: []*item{
						{
							Path: "/path/3/",
						},
					},
				},
			},
			SequenceOfEvents: []*sequenceElement{
				{},
			},
		},
		{
			GroupName:   "CCC",
			GroupTiming: "sameinstancenum",
			Panels: []*panel{
				{
					Not:         false,
					PanelTiming: "any",
					Items: []*item{
						{
							Path: "/path/4/",
						},
					},
				},
			},
		},
	},
}
