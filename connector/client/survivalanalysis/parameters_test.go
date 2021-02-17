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
		GroupName string "yaml:\"group_name\""
		Panels    []*struct {
			Not   bool "yaml:\"not\""
			Items []*struct {
				Path     string    `yaml:"path"`
				Modifier *modifier `yaml:"modifier,omitempty"`
			} "yaml:\"items\""
		} "yaml:\"panels\""
	}{
		{
			GroupName: "AAA",
			Panels: []*struct {
				Not   bool "yaml:\"not\""
				Items []*struct {
					Path     string    `yaml:"path"`
					Modifier *modifier `yaml:"modifier,omitempty"`
				} "yaml:\"items\""
			}{
				{
					Not: false,
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
					Not: true,
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
			GroupName: "BBB",
			Panels: []*struct {
				Not   bool "yaml:\"not\""
				Items []*struct {
					Path     string    `yaml:"path"`
					Modifier *modifier `yaml:"modifier,omitempty"`
				} "yaml:\"items\""
			}{
				{
					Not: false,
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
