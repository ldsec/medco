package medcoclient

import (
	"fmt"
	"strconv"

	medcomodels "github.com/ldsec/medco/connector/models"
	utilclient "github.com/ldsec/medco/connector/util/client"
	"github.com/sirupsen/logrus"
)

// DumpTimers writes timers in a file if timerFile string is not empty. Print it in standard output otherwise.
func DumpTimers(timerFile string, timers []medcomodels.Timers, clientTimers medcomodels.Timers) error {
	// print timers
	logrus.Info("dumping timers")
	dumpCSV, err := utilclient.NewCSV(timerFile)
	if err != nil {
		err = fmt.Errorf("while creating CSV file handler: %s", err)
		logrus.Error(err)
		return err
	}
	dumpCSV.Write([]string{"node_index", "timer_description", "duration_milliseconds"})
	if err != nil {
		err = fmt.Errorf("while writing headers for timer file: %s", err)
		logrus.Error(err)
		return err
	}
	// each remote time profilings
	for nodeIdx, nodeTimers := range timers {
		sortedTimers := nodeTimers.SortTimers()
		for _, duration := range sortedTimers {
			dumpCSV.Write([]string{
				strconv.Itoa(nodeIdx),
				duration[0],
				duration[1],
			})
			if err != nil {
				err = fmt.Errorf("while writing record for timer file: %s", err)
				logrus.Error(err)
				return err
			}
		}

	}
	// and local
	localSortedTimers := clientTimers.SortTimers()
	for _, duration := range localSortedTimers {
		dumpCSV.Write([]string{
			"client",
			duration[0],
			duration[1],
		})
		if err != nil {
			err = fmt.Errorf("while writing record for timer file: %s", err)
			logrus.Error(err)
			return err
		}
	}

	err = dumpCSV.Flush()
	if err != nil {
		err = fmt.Errorf("while flushing timer file: %s", err)
		logrus.Error(err)
		return err
	}
	logrus.Info()
	err = dumpCSV.Close()
	if err != nil {
		err = fmt.Errorf("while closing timer file: %s", err)
		logrus.Error(err)
		return err
	}
	return nil
}
