package survivalserver

import (
	"sort"
	"time"

	"github.com/ldsec/medco-connector/wrappers/unlynx"
	"github.com/sirupsen/logrus"
)

type sortedMapElm struct {
	key   string
	value [2]string
}

type sortedMap []*sortedMapElm

func (sMap sortedMap) Less(i, j int) bool {
	return sMap[i].key < sMap[j].key

}

func (sMap sortedMap) Swap(i, j int) {
	sMap[i], sMap[j] = sMap[j], sMap[i]
}

func (sMap sortedMap) Len() int {
	return len(sMap)
}

type sortedMapMapElm struct {
	key   string
	value sortedMap
}

type sortedMapMap []*sortedMapMapElm

func (smapmap sortedMapMap) Less(i, j int) bool {
	return smapmap[i].key < smapmap[j].key
}

func (smapmap sortedMapMap) Swap(i, j int) {
	smapmap[i], smapmap[j] = smapmap[j], smapmap[i]
}

func (smapmap sortedMapMap) Len() int {
	return len(smapmap)
}

func AKSgroups(queryID string, eventGroups map[string]map[string][2]string, targetPubKey string) (aggEventGroups map[string]map[string][2]string, time map[string]time.Duration, err error) {
	nofEntries := 0
	sortedEventGroups := sortedMapMap(make(sortedMapMap, 0))
	for groupID, timePointToEvent := range eventGroups {
		sortedEvents := sortedMap(make(sortedMap, 0))
		for timePoint, events := range timePointToEvent {
			sortedEvents = append(sortedEvents, &sortedMapElm{key: timePoint, value: events})
			nofEntries += 2
		}
		sort.Sort(sortedEvents)
		sortedEventGroups = append(sortedEventGroups, &sortedMapMapElm{key: groupID, value: sortedEvents})
	}
	sort.Sort(sortedEventGroups)

	values := make([]string, 0)
	for _, sortedEventGroup := range sortedEventGroups {
		for _, value := range sortedEventGroup.value {
			values = append(values, value.value[0])
			values = append(values, value.value[1])
		}

	}
	if len(values) != nofEntries {
		logrus.Panic("something went wrong while flattening the associative arrays")
	}
	aggEventGroups = make(map[string]map[string][2]string)
	var aggValues []string
	aggValues, time, err = unlynx.AggregateAndKeySwitchValues(queryID, values, targetPubKey)
	if err != nil {
		return
	}

	position := 0

	for _, sortedEventGroup := range sortedEventGroups {
		aggEvents := make(map[string][2]string)
		for _, value := range sortedEventGroup.value {
			aggEvents[value.key] = [2]string{aggValues[position], aggValues[position+1]}
			position += 2

		}
		aggEventGroups[sortedEventGroup.key] = aggEvents
	}

	return

}
