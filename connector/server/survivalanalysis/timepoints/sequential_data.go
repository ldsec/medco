package timepoints

import (
	"fmt"
	"sort"
	"time"

	medcomodels "github.com/ldsec/medco/connector/models"
	"github.com/sirupsen/logrus"
)

// patientAndEndEvents take the patient-to-start-event map and the patient-to-end-event-candidates.
// For each patient in the in the first map, it checks its presence in the second one. If the latter contains data.
// endEarliest defines if it must take the earliest or the latest among candidates. Candidates must occur striclty after the start event, an error is thrown otherwise.
// The list of candidate events is not expected to be empty, an error is thrown if it is the case.
// The patient-to-difference-in-day map is returned alongside the list of patients present in the patient-to-start-event map and absent from patient-to-end-event.
func patientAndEndEvents(startEvent map[int64]time.Time, endEvents map[int64][]time.Time, endEarliest bool) (map[int64]struct{}, map[int64]int64, error) {

	patientsWithoutEndEvent := make(map[int64]struct{}, len(startEvent))
	patientsWithStartAndEndEvents := make(map[int64]int64, len(startEvent))
	for patientID, startDate := range startEvent {
		if endDates, isIn := endEvents[patientID]; isIn {
			if endDates == nil {
				err := fmt.Errorf("unexpected nil end-date list for patient %d", patientID)
				return nil, nil, err
			}
			nofEndDates := len(endDates)
			if nofEndDates == 0 {
				err := fmt.Errorf("unexpected empty end-date list for patient %d", patientID)
				return nil, nil, err
			}
			sort.Slice(endDates, func(i, j int) bool {
				return endDates[i].Before(endDates[j])
			})

			var endDate time.Time
			if endEarliest {
				endDate = endDates[0]
			} else {
				endDate = endDates[nofEndDates-1]
			}

			diffInHours := endDate.Sub(startDate).Hours()
			truncatedDiff := int64(diffInHours)
			if remaining := truncatedDiff % 24; remaining != 0 {
				err := fmt.Errorf("the remaining of the time difference must be divisible by 24, the remaining is actually %d", remaining)
				return nil, nil, err
			}
			numberInDays := truncatedDiff / 24

			if numberInDays <= 0 {
				err := fmt.Errorf("the difference is expected to be strictly greater than 0, actually got %d", numberInDays)
				return nil, nil, err
			}
			patientsWithStartAndEndEvents[patientID] = numberInDays

		} else {
			patientsWithoutEndEvent[patientID] = struct{}{}
		}
	}
	return patientsWithoutEndEvent, patientsWithStartAndEndEvents, nil
}

// patientAndCensoring takes the patient-to-start-event, the patient-without-end-event set and the patient-to-censoring map
// and compute the difference in day for each patients in the paitent-without-end-event between the censoring time taken from the second map and
// and the start time taken from the first map. The set of patients without end event is expected to be a subset of the patient-to-start-event keys and
// censoring events must happen strictly after the start event, an error is thrown otherwise.
// The patient-to-difference-in-day (for censoring events) is returned:
func patientAndCensoring(startEvent map[int64]time.Time, patientsWithoutEndEvent map[int64]struct{}, patientWithCensoring map[int64]time.Time) (map[int64]int64, error) {
	patientsWithStartAndCensoring := make(map[int64]int64, len(startEvent))
	for patientID := range patientsWithoutEndEvent {
		if endDate, isIn := patientWithCensoring[patientID]; isIn {
			startDate, isFound := startEvent[patientID]
			if !isFound {
				err := fmt.Errorf("the set of patients without the end event of interest must be a subset of the start-event keys: patient %d found in patients without events of interest, but is not a start-event key", patientID)
				return nil, err
			}

			diffInHours := endDate.Sub(startDate).Hours()
			truncatedDiff := int64(diffInHours)
			if remaining := truncatedDiff % 24; remaining != 0 {
				err := fmt.Errorf("the remaining of the time difference must be divisible by 24, the remaining is actually %d", remaining)
				return nil, err
			}
			numberInDays := truncatedDiff / 24

			if numberInDays <= 0 {
				err := fmt.Errorf("the difference is expected to be strictly greater than 0, actually got %d", numberInDays)
				return nil, err
			}
			patientsWithStartAndCensoring[patientID] = numberInDays
		}
	}
	return patientsWithStartAndCensoring, nil
}

// compileTimePoints takes the patient-to-end-event and the patient-to-censoring-event maps and aggregates te number of events, grouped by difference in days (aka relative times).
// If a relative time is strictly bigger than the max limit defined by the user, it is ignored. If the relative time or the maximum limit is smaller or equal to  zero, an error is thrown.
func compileTimePoints(patientWithEndEvents, patientWithCensoringEvents map[int64]int64, maxLimit int64) (map[int64]*medcomodels.Events, error) {
	if maxLimit <= 0 {
		err := fmt.Errorf("user-defined maximum limit %d must be strictly greater than 0", maxLimit)
		return nil, err
	}
	timePointTable := make(map[int64]*medcomodels.Events, int(maxLimit))
	for _, timePoint := range patientWithEndEvents {
		if timePoint > maxLimit {
			logrus.Tracef("Survival analysis: timepoint: timepoint %d beyond user-defined limit %d; dropped", timePoint, maxLimit)
			continue
		}
		if timePoint <= 0 {
			err := fmt.Errorf("while computing events aggregates: relative time in patients with end event must be strictly greater than 0, got %d", timePoint)
			return nil, err
		}
		if _, isIn := timePointTable[timePoint]; !isIn {
			timePointTable[timePoint] = &medcomodels.Events{
				EventsOfInterest: 1,
				CensoringEvents:  0,
			}
		} else {
			elm := timePointTable[timePoint]
			elm.EventsOfInterest++
		}
	}

	for _, timePoint := range patientWithCensoringEvents {
		if timePoint > maxLimit {
			logrus.Tracef("Survival analysis: timepoint: timepoint %d beyond user-defined limit %d; dropped", timePoint, maxLimit)
			continue
		}
		if timePoint <= 0 {
			err := fmt.Errorf("while computing events aggregates: relative time in patients with censoring event must be strictly greater than 0, got %d", timePoint)
			return nil, err
		}
		if _, isIn := timePointTable[timePoint]; !isIn {
			timePointTable[timePoint] = &medcomodels.Events{
				EventsOfInterest: 0,
				CensoringEvents:  1,
			}
		} else {
			elm := timePointTable[timePoint]
			elm.CensoringEvents++
		}
	}
	return timePointTable, nil
}
