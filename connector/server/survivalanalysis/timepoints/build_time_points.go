package timepoints

import medcomodels "github.com/ldsec/medco/connector/models"

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
	patientWithoutAnyEvent map[int64]struct{},
	err error,
) {

	patientsToStartEvent, err := startEvent(patientSet, startConceptCodes, startModifierCodes, startEarliest)
	patientsToEndEvents, err := endEvents(patientsToStartEvent, endConceptCodes, endModifierCodes)

	patientsWithoutEnd, startToEndEvent, err := patientAndEndEvents(patientsToStartEvent, patientsToEndEvents, endEarliest)

	patientsToCensoringEvent, patientWithoutAnyEvent, err := censoringEvents(patientsToStartEvent, patientsWithoutEnd, endConceptCodes, endModifierCodes)

	startToCensoringEvent, err := patientAndCensoring(patientsToStartEvent, patientsWithoutEnd, patientsToCensoringEvent)
	eventAggregates = compileTimePoints(startToEndEvent, startToCensoringEvent, maxLimit)

	return
}
