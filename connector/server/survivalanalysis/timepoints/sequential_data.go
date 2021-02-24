package timepoints

import (
	"fmt"
	"sort"
	"time"
)

func PatientAndEndEvents(startEvent map[int64]time.Time, endEvents map[int64][]time.Time, endEarliest bool) (map[int64]struct{}, map[int64]int64, error) {

	patientsWithoutEndEvent := make(map[int64]struct{}, len(startEvent))
	patientsWithStartAndEndEvents := make(map[int64]int64, len(startEvent))
	for patientID, startDate := range startEvent {
		if endDates, isIn := endEvents[patientID]; isIn {
			sort.Slice(endDates, func(i, j int) bool {
				return endDates[i].Before(endDates[j])
			})

			var endDate time.Time
			if endEarliest {
				endDate = endDates[0]
			} else {
				endDate = endDates[len(endDates)-1]
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

func PatientAndCensoring(startEvent map[int64]time.Time, patientsWithoutEndEvent map[int64]struct{}, patientWithCesoring map[int64]time.Time) (map[int64]struct{}, map[int64]int64, error)

func BuildTimePoints(patientWithEndEvents, patientWithCensoringEvents map[int64]struct {
	startDate time.Time
	endDate   time.Time
}) (map[int64]int64, error)
