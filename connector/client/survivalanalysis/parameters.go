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
	TimeResolution       string `yaml:"time_resolution"`
	TimeLimit            int    `yaml:"time_limit"`
	CohortName           string `yaml:"cohort_name"`
	StartConceptPath     string `yaml:"start_concept_path"`
	StartConceptModifier string `yaml:"start_concept_modifier"`
	EndConceptPath       string `yaml:"end_concept_path"`
	EndConceptModifier   string `yaml:"end_concept_modifier"`
	SubGroups            []*struct {
		GroupName string `yaml:"group_name"`
		Panels    []*struct {
			Not   bool     `yaml:"not"`
			Paths []string `yaml:"paths"`
		} `yaml:"panels"`
	} `yaml:"sub_groups,omitempty"`
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

	return fmt.Sprintf("{TimeResolution:%s TimeLimit:%d CohortName:%s StartConceptPath:%s StartConceptModifier:%s EndConceptPath:%s EndConceptModifier:%s SubGroups:%s}",
		p.TimeResolution,
		p.TimeLimit,
		p.CohortName,
		p.StartConceptPath,
		p.StartConceptModifier,
		p.EndConceptPath,
		p.EndConceptModifier,
		subGroupArray)
}
