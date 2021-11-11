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
	TimeResolution   string      `yaml:"time_resolution"`
	TimeLimit        int         `yaml:"time_limit"`
	CohortName       string      `yaml:"cohort_name"`
	StartConceptPath string      `yaml:"start_concept_path"`
	StartModifier    *modifier   `yaml:"start_modifier,omitempty"`
	StartsWhen       string      `yaml:"starts_when"`
	EndConceptPath   string      `yaml:"end_concept_path"`
	EndModifier      *modifier   `yaml:"end_modifier,omitempty"`
	EndsWhen         string      `yaml:"ends_when"`
	SubGroups        []*subGroup `yaml:"sub_groups,omitempty"`
}

type conceptItem struct {
	Path     string    `yaml:"path"`
	Modifier *modifier `yaml:"modifier,omitempty"`
	Operator string    `yaml:"operator,omitempty"`
	Type     string    `yaml:"type,omitempty"`
	Value    string    `yaml:"value,omitempty"`
}

type panel struct {
	Not          bool           `yaml:"not"`
	ConceptItems []*conceptItem `yaml:"concept_items"`
	CohortItems  []string       `yaml:"cohort_items"`
	PanelTiming  string         `yaml:"panel_timing"`
}

type subGroup struct {
	GroupName        string             `yaml:"group_name"`
	GroupTiming      string             `yaml:"group_timing"`
	SelectionPanels  []*panel           `yaml:"selection_panels"`
	SequentialPanels []*panel           `yaml:"sequential_panels"`
	SequenceOfEvents []*sequenceElement `yaml:"sequence_of_events,omitempty"`
}

type modifier struct {
	ModifierKey string `yaml:"modifier_key"`
	AppliedPath string `yaml:"applied_path"`
}

type sequenceElement struct {
	WhichDateFirst         string      `yaml:"which_date_first,omitempty"`
	WhichDateSecond        string      `yaml:"which_date_second,omitempty"`
	WhichObservationFirst  string      `yaml:"which_observation_first,omitempty"`
	WhichObservationSecond string      `yaml:"which_observation_second,omitempty"`
	When                   string      `yaml:"when,omitempty"`
	Spans                  []*timeSpan `yaml:"spans,omitempty"`
}

type timeSpan struct {
	Operator string `yaml:"operator"`
	Value    int64  `yaml:"value"`
	Units    string `yaml:"units"`
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
		selectionPanelStrings := make([]string, 0, len(subGroup.SelectionPanels))
		for _, selectionPanel := range subGroup.SelectionPanels {
			selectionPanelStrings = append(selectionPanelStrings, fmt.Sprintf("%+v", selectionPanel))
		}
		sequentialPanelStrings := make([]string, 0, len(subGroup.SequentialPanels))
		for _, sequentialPanel := range subGroup.SequentialPanels {
			sequentialPanelStrings = append(sequentialPanelStrings, fmt.Sprintf("%+v", sequentialPanel))
		}
		selectionPanelArray := "[" + strings.Join(selectionPanelStrings, " ") + "]"
		sequentialPanelArray := "[" + strings.Join(sequentialPanelStrings, " ") + "]"
		subGroupStrings = append(subGroupStrings,
			fmt.Sprintf("{GroupName:%s SelectionPanels:%s SequentialPanels:%s}", subGroup.GroupName, selectionPanelArray, sequentialPanelArray))
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
