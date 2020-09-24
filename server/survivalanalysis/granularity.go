package survivalserver

import (
	"fmt"
	"math"

	utilcommon "github.com/ldsec/medco-connector/util/common"
)

const (
	dInWeek  = 7
	dInMonth = 30
	dInYear  = 365
)

var granularityFunctions = map[string]func(int) int{
	"day":   func(x int) int { return x },
	"week":  week,
	"month": month,
	"year":  year,
}

func granularity(points utilcommon.TimePoints, granularity string) (utilcommon.TimePoints, error) {
	if granFunction, isIn := granularityFunctions[granularity]; isIn {
		return binTimePoint(points, granFunction), nil
	}
	return nil, fmt.Errorf("granularity %s is not implemented: should be one of year, month, week, day", granularity)

}

func ceil(val int, granularity int) int {
	return int(math.Ceil(float64(val) / float64(granularity)))
}

func week(val int) int {
	return ceil(val, dInWeek)
}

func month(val int) int {
	return ceil(val, dInMonth)
}

func year(val int) int {
	return ceil(val, dInYear)
}

func binTimePoint(timePoints utilcommon.TimePoints, groupingFunction func(int) int) utilcommon.TimePoints {
	bins := make(map[int]struct {
		EventsOfInterest int64
		CensoringEvents  int64
	})
	var ceiled int
	for _, tp := range timePoints {
		ceiled = groupingFunction(tp.Time)
		if val, isInside := bins[ceiled]; isInside {
			bins[ceiled] = struct {
				EventsOfInterest int64
				CensoringEvents  int64
			}{
				EventsOfInterest: val.EventsOfInterest + tp.Events.EventsOfInterest,
				CensoringEvents:  val.CensoringEvents + tp.Events.CensoringEvents,
			}
		} else {
			bins[ceiled] = struct {
				EventsOfInterest int64
				CensoringEvents  int64
			}{
				EventsOfInterest: tp.Events.EventsOfInterest,
				CensoringEvents:  tp.Events.CensoringEvents,
			}
		}
	}

	newSQLTimePoints := make(utilcommon.TimePoints, 0)
	for time, agg := range bins {
		newSQLTimePoints = append(newSQLTimePoints, utilcommon.TimePoint{
			Time: time,
			Events: struct {
				EventsOfInterest int64
				CensoringEvents  int64
			}{
				EventsOfInterest: agg.EventsOfInterest,
				CensoringEvents:  agg.CensoringEvents,
			},
		})
	}
	return newSQLTimePoints
}
