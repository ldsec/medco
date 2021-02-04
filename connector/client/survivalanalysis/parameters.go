package survivalclient

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Parameters holds parameters to build a survival query
type Parameters struct {
	TimeResolution   string    `yaml:"time_resolution"`
	TimeLimit        int       `yaml:"time_limit"`
	CohortName       string    `yaml:"cohort_name"`
	StartConceptPath string    `yaml:"start_concept_path"`
	StartModifier    *modifier `yaml:"start_modifier,omitempty"`
	StartsWhen       string    `yaml:"starts_when"`
	EndConceptPath   string    `yaml:"end_concept_path"`
	EndModifier      *modifier `yaml:"end_modifier,omitempty"`
	EndsWhen         string    `yaml:"ends_when"`
	SubGroups        []*struct {
		GroupName string `yaml:"group_name"`
		Panels    []*struct {
			Not   bool `yaml:"not"`
			Items []*struct {
				Path     string    `yaml:"path"`
				Modifier *modifier `yaml:"modifier,omitempty"`
			} `yaml:"items"`
		} `yaml:"panels"`
	} `yaml:"sub_groups,omitempty"`
}

type modifier struct {
	ModifierKey string `yaml:"modifier_key"`
	AppliedPath string `yaml:"applied_path"`
}

// NewParametersFromFile builds a Parameters instance from YAML file
func NewParametersFromFile(fileURL string) (*Parameters, error) {
	logrus.Trace(fileURL)
	logrus.Trace(os.Getenv("PWD"))
	file, err := os.Open(fileURL)
	if err != nil {
		return nil, fmt.Errorf("while opening parameter file %s: %s", fileURL, err.Error())
	}
	decoder := yaml.NewDecoder(file)
	params := &Parameters{}
	err = decoder.Decode(params)
	if err != nil {
		return nil, fmt.Errorf("while decoding parameter file: %s", err.Error())
	}

	err = file.Close()
	if err != nil {
		return nil, fmt.Errorf("while closing parameter file: %s", err.Error())
	}
	return params, nil
}

// String inplements the Stringer interface
func (p *Parameters) String() string {

	subGroupStrings := make([]string, 0, len(p.SubGroups))
	for _, subGroup := range p.SubGroups {
		panelStrings := make([]string, 0, len(subGroup.Panels))
		for _, panel := range subGroup.Panels {
			panelStrings = append(panelStrings, fmt.Sprintf("%+v", panel))
		}
		panelArray := "[" + strings.Join(panelStrings, " ") + "]"
		subGroupStrings = append(subGroupStrings,
			fmt.Sprintf("{GroupName:%s Panels:%s}", subGroup.GroupName, panelArray))
	}
	subGroupArray := "[" + strings.Join(subGroupStrings, " ") + "]"

	startModifierString := ""
	if startMod := p.StartModifier; startMod != nil {
		startModifierString = fmt.Sprintf("{ModifierKey:%s AppliedPath:%s}", startMod.ModifierKey, startMod.AppliedPath)
	}

	endModifierString := ""
	if endMod := p.EndModifier; endMod != nil {
		endModifierString = fmt.Sprintf("{ModifierKey:%s AppliedPath:%s}", endMod.ModifierKey, endMod.AppliedPath)
	}

	return fmt.Sprintf("{TimeResolution:%s TimeLimit:%d CohortName:%s StartConceptPath:%s StartModifier:%s StartsWhen:%s EndConceptPath:%s EndModifier:%s EndsWhen:%s SubGroups:%s}",
		p.TimeResolution,
		p.TimeLimit,
		p.CohortName,
		p.StartConceptPath,
		startModifierString,
		p.StartsWhen,
		p.EndConceptPath,
		endModifierString,
		p.EndsWhen,
		subGroupArray)
}
