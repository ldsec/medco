package timepoints

import (
	medcomodels "github.com/ldsec/medco/connector/models"
)

// BuildTimePoints runs the SQL queries, process their results to build sequential data and aggregate them
func BuildTimePoints(
	patientSet []int64,
	startConceptCodes []string,
	startModifierCodes []string,
	startEarliest bool,
	endConceptCodes []string,
	endModifierCodes []string,
	endEarliest bool,
	maxLimit int64,
) (
	eventAggregates map[int64]*medcomodels.Events,
	patientWithoutStartEvent map[int64]struct{},
	patientWithoutAnyEndEvent map[int64]struct{},
	err error,
) {

	patientsToStartEvent, patientWithoutStartEvent, err := startEvent(patientSet, startConceptCodes, startModifierCodes, startEarliest)
	if err != nil {
		return
	}
	patientsToEndEvents, err := endEvents(patientsToStartEvent, endConceptCodes, endModifierCodes)
	if err != nil {
		return
	}
	patientsWithoutEnd, startToEndEvent, err := patientAndEndEvents(patientsToStartEvent, patientsToEndEvents, endEarliest)
	if err != nil {
		return
	}
	patientsToCensoringEvent, patientWithoutAnyEndEvent, err := censoringEvents(patientsToStartEvent, patientsWithoutEnd, endConceptCodes, endModifierCodes)
	if err != nil {
		return
	}
	startToCensoringEvent, err := patientAndCensoring(patientsToStartEvent, patientsWithoutEnd, patientsToCensoringEvent)
	if err != nil {
		return
	}
	eventAggregates = compileTimePoints(startToEndEvent, startToCensoringEvent, maxLimit)
	if err != nil {
		return
	}
	return
}
