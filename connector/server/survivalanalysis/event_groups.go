package survivalserver

import (
	"fmt"
)

// Result holds information about time point, events, execution times and error
type Result struct {
	EventValueAgg     string
	CensoringValueAgg string
}

// TimePointResult holds information about time point, events, execution times and error
type TimePointResult struct {
	TimePoint int64
	Result    Result
}

// String implements Stringer interface for TimePointResult pointers
func (t *TimePointResult) String() string {
	return fmt.Sprintf("{%d,%v}", t.TimePoint, t.Result)
}

// EventGroup implements sort.Interface interface. It holds the group name, the encryption of the inital group size and the list of encrypted timepoints
type EventGroup struct {
	GroupID          string
	EncInitialCount  string
	TimePointResults []*TimePointResult
}

// String implements Stringer interface for EventGroup pointers
func (eventGroup *EventGroup) String() string {
	return fmt.Sprintf("{%s,%s,%v}", eventGroup.GroupID, eventGroup.EncInitialCount, eventGroup.TimePointResults)
}

// EventGroups implements sort.Interface interface. It holds a list of EventGroup instances.
type EventGroups []*EventGroup

// Len returns the number of elements
func (eventGroups EventGroups) Len() int {
	return len(eventGroups)
}

// Less returns true if element at position i is smaller than element at position j
func (eventGroups EventGroups) Less(i, j int) bool {
	return eventGroups[i].GroupID < eventGroups[j].GroupID
}

// Swap exchanges elements at positions i and j
func (eventGroups EventGroups) Swap(i, j int) {
	eventGroups[i], eventGroups[j] = eventGroups[j], eventGroups[i]
}

// Len returns the number of elements
func (eventGroup EventGroup) Len() int {
	return len(eventGroup.TimePointResults)
}

// Less returns true if element at position i is smaller than element at position j
func (eventGroup EventGroup) Less(i, j int) bool {
	return eventGroup.TimePointResults[i].TimePoint < eventGroup.TimePointResults[j].TimePoint
}

// Swap exchanges elements at positions i and j
func (eventGroup EventGroup) Swap(i, j int) {
	eventGroup.TimePointResults[i], eventGroup.TimePointResults[j] = eventGroup.TimePointResults[j], eventGroup.TimePointResults[i]
}
