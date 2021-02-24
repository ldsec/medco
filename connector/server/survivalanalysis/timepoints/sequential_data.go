package timepoints

import "time"

func PatientAndEndEvents(startEvent map[int64]time.Time, endEvents []struct {
	patient int64
	date    time.Time
}, endEarliest bool) (map[int64]struct{}, map[int64]struct {
	startDate time.Time
	endDate   time.Time
}, error)

func PatientAndCensoring(startEvent map[int64]time.Time, patientWithCesoring map[int64]time.Time) (map[int64]struct{}, map[int64]struct {
	startDate time.Time
	endDate   time.Time
}, error)

func BuildTimePoints(patientWithEndEvents, patientWithCensoringEvents map[int64]struct {
	startDate time.Time
	endDate   time.Time
}) (map[int64]int64, error)
