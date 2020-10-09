package utilclient

import (
	"sort"
	"strconv"

	utilcommon "github.com/ldsec/medco/connector/util/common"
)

// SortTimers takes a Timers instance, whichi is a Golang map, and output a sorted 2D string array . This is useful for deterministic output as Golang maps are not deterministic. Stored times are converted from nanoseconds to milliseconds.
func SortTimers(timers utilcommon.Timers) [][]string {
	names := make([]string, 0, len(timers))
	res := make([][]string, len(timers))
	for name := range timers {
		names = append(names, name)
	}
	sort.Strings(names)

	for i, name := range names {
		res[i] = append([]string{name}, strconv.FormatInt(timers[name].Milliseconds(), 10))
	}
	return res
}
