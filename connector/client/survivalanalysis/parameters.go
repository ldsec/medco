package survivalclient

import (
	"fmt"
	"os"

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
	Cohorts              []*struct {
		Panels []*struct {
			Not   bool     `yaml:"not"`
			Paths []string `yaml:"paths"`
		} `yaml:"panels"`
	} `yaml:"cohorts,omitempty"`
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
