package survivalserver

import (
	"fmt"

	survivalcommon "github.com/ldsec/medco-connector/survival/common"
)

// Expansion takes a slice of SQLTimepoints and add encryption of zeros for events of interest and censoring events for each missing relative time from 0 to timeLimit.
// Relative times greater than timeLimit are discarded.
func Expansion(timePoints survivalcommon.TimePoints, timeLimitDay int, granularity string) (survivalcommon.TimePoints, error) {
	var timeLimit int
	if granFunction, isIn := granularityFunctions[granularity]; isIn {
		timeLimit = granFunction(timeLimitDay)
	} else {
		return nil, fmt.Errorf("granularity %s is not implemented", granularity)
	}

	res := make(survivalcommon.TimePoints, timeLimit)
	availableTimePoints := make(map[int]struct {
		EventsOfInterest int64
		CensoringEvents  int64
	}, len(timePoints))
	for _, timePoint := range timePoints {

		availableTimePoints[timePoint.Time] = timePoint.Events
	}
	for i := 0; i < timeLimit; i++ {
		if events, ok := availableTimePoints[i]; ok {
			res[i] = survivalcommon.TimePoint{
				Time:   i,
				Events: events,
			}
		} else {
			res[i] = survivalcommon.TimePoint{
				Time: i,
				Events: struct {
					EventsOfInterest int64
					CensoringEvents  int64
				}{
					EventsOfInterest: 0,
					CensoringEvents:  0,
				},
			}
		}

	}
	return res, nil
}
