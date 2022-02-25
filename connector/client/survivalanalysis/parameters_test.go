//go:build unit_test
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
	SubGroups: []*struct {
		GroupName   string "yaml:\"group_name\""
		GroupTiming string "yaml:\"group_timing\""
		Panels      []*struct {
			Not   bool "yaml:\"not\""
			Items []*struct {
				Path     string    `yaml:"path"`
				Modifier *modifier `yaml:"modifier,omitempty"`
			} "yaml:\"items\""
			PanelTiming string "yaml:\"panel_timing\""
		} "yaml:\"panels\""
	}{
		{
			GroupTiming: "any",
			GroupName:   "AAA",
			Panels: []*struct {
				Not   bool "yaml:\"not\""
				Items []*struct {
					Path     string    `yaml:"path"`
					Modifier *modifier `yaml:"modifier,omitempty"`
				} "yaml:\"items\""
				PanelTiming string "yaml:\"panel_timing\""
			}{
				{
					PanelTiming: "any",
					Not:         false,
					Items: []*struct {
						Path     string    `yaml:"path"`
						Modifier *modifier `yaml:"modifier,omitempty"`
					}{
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
					Items: []*struct {
						Path     string    `yaml:"path"`
						Modifier *modifier `yaml:"modifier,omitempty"`
					}{
						{
							Path: "/path/3/",
						},
					},
				},
			},
		},
		{
			GroupName:   "BBB",
			GroupTiming: "sameinstancenum",
			Panels: []*struct {
				Not   bool "yaml:\"not\""
				Items []*struct {
					Path     string    `yaml:"path"`
					Modifier *modifier `yaml:"modifier,omitempty"`
				} "yaml:\"items\""
				PanelTiming string "yaml:\"panel_timing\""
			}{
				{
					Not:         false,
					PanelTiming: "any",
					Items: []*struct {
						Path     string    `yaml:"path"`
						Modifier *modifier `yaml:"modifier,omitempty"`
					}{
						{
							Path: "/path/4/",
						},
					},
				},
			},
		},
	},
}
