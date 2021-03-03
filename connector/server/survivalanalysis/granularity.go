package survivalserver

import (
	"fmt"
	"math"

	medcomodels "github.com/ldsec/medco/connector/models"
)

const (
	dInWeek  = 7
	dInMonth = 30
	dInYear  = 365
)

var granularityValues = map[string]int{
	"day":   1,
	"week":  dInWeek,
	"month": dInMonth,
	"year":  dInYear,
}

var granularityFunctions = map[string]func(int64) int64{
	"day":   func(x int64) int64 { return x },
	"week":  week,
	"month": month,
	"year":  year,
}

func granularity(points medcomodels.TimePoints, granularity string) (medcomodels.TimePoints, error) {
	if granFunction, isIn := granularityFunctions[granularity]; isIn {
		return binTimePoint(points, granFunction), nil
	}
	return nil, fmt.Errorf("granularity %s is not implemented: should be one of year, month, week, day", granularity)

}

func ceil(val int, granularity int) int {
	return int(math.Ceil(float64(val) / float64(granularity)))
}

func week(val int64) int64 {
	return int64(ceil(int(val), dInWeek))
}

func month(val int64) int64 {
	return int64(ceil(int(val), dInMonth))
}

func year(val int64) int64 {
	return int64(ceil(int(val), dInYear))
}

func binTimePoint(timePoints medcomodels.TimePoints, groupingFunction func(int64) int64) medcomodels.TimePoints {
	bins := make(map[int64]struct {
		EventsOfInterest int64
		CensoringEvents  int64
	})
	var ceiled int64
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

	newSQLTimePoints := make(medcomodels.TimePoints, 0)
	for time, agg := range bins {
		newSQLTimePoints = append(newSQLTimePoints, medcomodels.TimePoint{
			Time: int64(time),
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
