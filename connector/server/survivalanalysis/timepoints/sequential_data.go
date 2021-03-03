package timepoints

import (
	"fmt"
	"sort"
	"time"

	medcomodels "github.com/ldsec/medco/connector/models"
	"github.com/sirupsen/logrus"
)

func patientAndEndEvents(startEvent map[int64]time.Time, endEvents map[int64][]time.Time, endEarliest bool) (map[int64]struct{}, map[int64]int64, error) {

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

func patientAndCensoring(startEvent map[int64]time.Time, patientsWithoutEndEvent map[int64]struct{}, patientWithCensoring map[int64]time.Time) (map[int64]int64, error) {
	patientsWithStartAndEndEvents := make(map[int64]int64, len(startEvent))
	logrus.Infof("tremeta %d", len(patientsWithoutEndEvent))
	for patientID := range patientsWithoutEndEvent {
		if endDate, isIn := patientWithCensoring[patientID]; isIn {
			startDate := startEvent[patientID]

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
			patientsWithStartAndEndEvents[patientID] = numberInDays
		}
	}
	return patientsWithStartAndEndEvents, nil
}

func compileTimePoints(patientWithEndEvents, patientWithCensoringEvents map[int64]int64, maxLimit int64) map[int64]*medcomodels.Events {
	timePointTable := make(map[int64]*medcomodels.Events, int(maxLimit)+1)
	for _, timePoint := range patientWithEndEvents {
		if timePoint > maxLimit {
			continue
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
			continue
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
	return timePointTable
}
