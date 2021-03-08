package timepoints

import (
	"time"

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
	timers medcomodels.Timers,
	err error,
) {
	timers = make(medcomodels.Timers)

	timer := time.Now()
	patientsToStartEvent, patientWithoutStartEvent, err := startEvent(patientSet, startConceptCodes, startModifierCodes, startEarliest)
	if err != nil {
		return
	}
	timers.AddTimers("build-time-points-start-event", timer, nil)

	timer = time.Now()
	patientsToEndEvents, err := endEvents(patientsToStartEvent, endConceptCodes, endModifierCodes)
	if err != nil {
		return
	}
	timers.AddTimers("build-time-points-end-events", timer, nil)

	timer = time.Now()
	patientsWithoutEnd, startToEndEvent, err := patientAndEndEvents(patientsToStartEvent, patientsToEndEvents, endEarliest)
	if err != nil {
		return
	}
	timers.AddTimers("build-time-points-sequential-data-event-of-interest", timer, nil)

	timer = time.Now()
	patientsToCensoringEvent, patientWithoutAnyEndEvent, err := censoringEvent(patientsToStartEvent, patientsWithoutEnd, endConceptCodes, endModifierCodes)
	if err != nil {
		return
	}
	timers.AddTimers("build-time-points-censoring-event", timer, nil)

	timer = time.Now()
	startToCensoringEvent, err := patientAndCensoring(patientsToStartEvent, patientsWithoutEnd, patientsToCensoringEvent)
	if err != nil {
		return
	}
	timers.AddTimers("build-time-points-sequential-data-censoring-event", timer, nil)

	timer = time.Now()
	eventAggregates, err = compileTimePoints(startToEndEvent, startToCensoringEvent, maxLimit)
	if err != nil {
		return
	}
	timers.AddTimers("build-time-points-aggregate-sequential-data", timer, nil)

	return
}
