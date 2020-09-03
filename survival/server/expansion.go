package survivalserver

func Expansion(timePoints []SqlTimePoint, timeLimit int) []SqlTimePoint {
	res := make([]SqlTimePoint, timeLimit)
	availableTimePoints := make(map[int][2]int, len(timePoints))
	for _, timePoint := range timePoints {

		availableTimePoints[timePoint.timePoint] = [2]int{timePoint.localEventAggregate, timePoint.localCensoringAggrete}
	}
	for i := 0; i < timeLimit; i++ {
		if events, ok := availableTimePoints[i]; ok {
			res[i] = SqlTimePoint{
				timePoint:             i,
				localEventAggregate:   events[0],
				localCensoringAggrete: events[1],
			}
		} else {
			res[i] = SqlTimePoint{
				timePoint:             i,
				localEventAggregate:   0,
				localCensoringAggrete: 0,
			}
		}

	}
	return res
}
