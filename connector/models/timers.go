package medcomodels

import (
	"sort"
	"strconv"
	"time"

	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/sirupsen/logrus"
)

// Timers stores execution times by string label. Execution times are in nanoseconds
type Timers map[string]time.Duration

// NewTimers return a Timers instance
func NewTimers() Timers {
	return make(map[string]time.Duration)
}

// AddTimers add a new recorded execution time and its name to a Timers instance. It also merges results from additional timers. This can be used for secondary timers created by go routine.
// If timerName is empty, the timers instance is not updated. If timerName is already stored, its value gets updated.
func (timers Timers) AddTimers(timerName string, since time.Time, additionalTimers Timers) {
	if timerName != "" {
		if _, isIn := timers[timerName]; isIn {
			logrus.Warnf("timer label %s already in recorded timers, overwriting previous value", timerName)
		}
		timers[timerName] = time.Since(since)
	} else {
		logrus.Warn("ignoring timer with empty label string")
	}

	if additionalTimers != nil {
		for k, v := range additionalTimers {
			timers[k] = v
		}
	}
}

// TimersToAPIModel converts a timer instance into timers from the API models. Execution times are converted from nanoseconds to milliseconds
func (timers Timers) TimersToAPIModel() models.Timers {
	res := make(models.Timers, 0)
	for timerName, timerDuration := range timers {
		milliseconds := new(int64)
		*milliseconds = int64(timerDuration / time.Millisecond)
		res = append(res, &models.TimersItems0{
			Name:         timerName,
			Milliseconds: milliseconds,
		})
	}
	return res
}

// SortTimers takes a Timers instance, whichi is a Golang map, and output a sorted 2D string array . This is useful for deterministic output as Golang maps are not deterministic. Stored times are converted from nanoseconds to milliseconds.
func (timers Timers) SortTimers() [][]string {
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

// NewTimersFromAPIModel converts the API model for timers into a Timers instance. Execution times are converted from nilliseconds to nanoseconds
func NewTimersFromAPIModel(APITimers models.Timers) Timers {
	res := NewTimers()
	for _, timer := range APITimers {
		res[timer.Name] = time.Duration(*timer.Milliseconds) * time.Millisecond
	}
	return res
}
