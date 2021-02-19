package explorestatisticsclient

import (
	"fmt"
	"os"

	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
)

// Parameters holds parameters to build an explore statistics query
type Parameters struct {
	CohortDefinition []*models.Panel `yaml:"cohortDefinition"`
	ConceptsPaths    []string        `yaml:"concepts"`
	Modifiers        []*modifier     `yaml:"modifiers,omitempty"`
	nbBuckets        int64           `yaml:"numberOfBuckets"`
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

	modifiersString := ""
	for _, m := range p.Modifiers {
		modifiersString = modifiersString + fmt.Sprintf("{ModifierKey:%s AppliedPath:%s}", m.ModifierKey, m.AppliedPath)
	}

	return fmt.Sprintf("{Concepts:%s Modifiers:%s}",
		p.ConceptsPaths,
		modifiersString)
}
