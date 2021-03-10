package medcomodels

// TimePoint contains a relative time, the numbers of events of interest and censoring events occured at that time
type TimePoint struct {
	Time   int64
	Events Events
}

// Events contains the number of events of interest and censoring events occuring at the same relative time
type Events struct {
	EventsOfInterest int64
	CensoringEvents  int64
}

// TimePoints is a slice containing the time points and respectif counts of censoring events and events of interest.
// TimePoints implements sort.Interface interface.
type TimePoints []TimePoint

// Len implements Len method for sort.Interface interface
func (points TimePoints) Len() int {
	return len(points)
}

// Less implements Less method for sort.Interface interface
func (points TimePoints) Less(i, j int) bool {
	return points[i].Time < points[j].Time
}

//Swap implements Swap method for sort.Interface interface
func (points TimePoints) Swap(i, j int) {
	points[i], points[j] = points[j], points[i]
}

// TimePointsFromTable is a testing helper function that builds TimePoints instance from a 2D int array
func TimePointsFromTable(array [][]int64) TimePoints {
	res := make(TimePoints, len(array))
	for i, point := range array {
		res[i] =
			TimePoint{Time: point[0],
				Events: struct {
					EventsOfInterest int64
					CensoringEvents  int64
				}{
					EventsOfInterest: int64(point[1]),
					CensoringEvents:  int64(point[2]),
				},
			}

	}
	return res
}
