package survivalserver

import (
	"errors"
	"sort"
	"time"

	"github.com/ldsec/medco-connector/wrappers/unlynx"
)

func AKSgroups(queryID string, eventGroups EventGroups, targetPubKey string) (aksEventGroups EventGroups, time map[string]time.Duration, err error) {

	if len(eventGroups) == 0 {
		err = errors.New("no group")
		return
	}

	var cumulativeLength int

	//-------- deep copy and sorting by keys
	aksEventGroups = EventGroups{}

	for _, group := range eventGroups {
		timePointResults := []*TimePointResult{}
		for _, res := range group.TimePointResults {
			cumulativeLength++
			timePointResults = append(timePointResults, &TimePointResult{
				TimePoint: res.TimePoint,
				Result:    res.Result,
			})
		}
		eventGroup := &EventGroup{GroupID: group.GroupID, TimePointResults: timePointResults}
		sort.Sort(eventGroup)
		aksEventGroups = append(aksEventGroups, eventGroup)
	}

	if cumulativeLength == 0 {
		err = errors.New("all groups are empty")
		return
	}

	sort.Sort(aksEventGroups)

	// ---------  flattening
	var flatInputs []string
	for _, group := range aksEventGroups {
		for _, timePoint := range group.TimePointResults {
			flatInputs = append(flatInputs, timePoint.Result.EventValueAgg)
			flatInputs = append(flatInputs, timePoint.Result.CensoringValueAgg)
		}
	}
	//TODO already tested
	if len(flatInputs) == 0 {
		err = errors.New("no data to aggregate")
		return
	}

	var flatOutputs []string
	flatOutputs, time, err = unlynx.AggregateAndKeySwitchValues(queryID, flatInputs, targetPubKey)
	if err != nil {
		return
	}
	//logrus.Panicf("aasasdasfasfsafdafsafsa %d", len(flatOutputs))

	position := 0

	for _, aksEventGroup := range aksEventGroups {
		for _, timePoint := range aksEventGroup.TimePointResults {

			timePoint.Result.EventValueAgg = flatOutputs[position]
			timePoint.Result.CensoringValueAgg = flatOutputs[position+1]
			position += 2

		}
	}
	//err = fmt.Errorf(" meeeeee %v mmmmmoooooo  %s maaaaaaaaaaa %s    ", aksEventGroups, aksEventGroups[0].GroupID, aksEventGroups[0].TimePointResults[0].TimePoint)

	return

}
