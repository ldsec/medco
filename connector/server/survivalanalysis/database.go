package survivalserver

import (
	"fmt"
	"strconv"
	"strings"

	medcomodels "github.com/ldsec/medco/connector/models"

	utilserver "github.com/ldsec/medco/connector/util/server"

	"github.com/sirupsen/logrus"
)

// buildTimePoints execute a SQL query that returns event counts per time point, for given input patient set, start and end  concept codes and modifiers
func buildTimePoints(
	patientList []int64,
	startConceptCodes []string,
	startConceptModifiers []string,
	startEarliest bool,
	endConceptCodes []string,
	endConceptModifiers []string,
	endEarliest bool,
	timeLimit int,
) (timePoints medcomodels.TimePoints, err error) {

	pList := make([]string, len(patientList))
	for i, pNum := range patientList {
		pList[i] = strconv.FormatInt(pNum, 10)
	}
	patients := "{" + strings.Join(pList, ",") + "}"
	startConcepts := "{" + strings.Join(startConceptCodes, ",") + "}"
	endConcepts := "{" + strings.Join(endConceptCodes, ",") + "}"
	startModifiers := "{" + strings.Join(startConceptModifiers, ",") + "}"
	endModifiers := "{" + strings.Join(endConceptModifiers, ",") + "}"
	description := fmt.Sprintf(
		"Build time points (patient list :%s, start concept codes: %s, start modifier codes: %s, start earliest: %t,"+
			" end concept codes: %s, end modifier codes: %s, end earliest: %t, time limit: %d), procedure: %s",
		patients,
		startConcepts,
		startModifiers,
		startEarliest,
		endConcepts,
		endModifiers,
		endEarliest,
		timeLimit,
		"i2b2demodata_i2b2.build_timepoints",
	)
	logrus.Debugf("running: %s", description)
	rows, err := utilserver.I2B2DBConnection.Query(
		"SELECT i2b2demodata_i2b2.build_timepoints($1, $2, $3, $4, $5, $6, $7, $8)",
		patients,
		startConcepts,
		startModifiers,
		startEarliest,
		endConcepts,
		endModifiers,
		endEarliest,
		timeLimit,
	)
	if err != nil {
		err = fmt.Errorf("while execution SQL query: %s, DB operation: %s", err.Error(), description)
		return
	}
	logrus.Debug("successfully selected")

	// initialize the response
	allTimePoints := make(medcomodels.TimePoints, timeLimit)
	for i := range allTimePoints {
		allTimePoints[i].Time = i
	}

	record := new(string)

	var timePoint int64
	var counts int64
	var eventType int

	for rows.Next() {

		scanErr := rows.Scan(record)
		if scanErr != nil {
			err = scanErr
			err = fmt.Errorf("while scanning SQL record: %s", err.Error())
			return
		}
		logrus.Tracef("Record: %s", *record)
		cells := strings.Split(strings.Trim(*record, "()"), ",")
		timePoint, err = strconv.ParseInt(cells[0], 10, 64)
		if err != nil {
			err = fmt.Errorf("while scanning parsing integer string (relative time) \"%s\": %s", cells[0], err.Error())
			return
		}
		eventType, err = strconv.Atoi(cells[1])
		if err != nil {
			err = fmt.Errorf("while scanning parsing integer string (type of event) \"%s\": %s", cells[2], err.Error())
			return
		}
		counts, err = strconv.ParseInt(cells[2], 10, 64)
		if err != nil {
			err = fmt.Errorf("while scanning parsing integer string (number of events) \"%s\": %s", cells[1], err.Error())
			return
		}

		// 1 is for event of interest, 0 for censoring event
		if timePoint >= int64(timeLimit) {
			err = fmt.Errorf("Unexpected time point code %d, must be smaller than time limit %d", timePoint, timeLimit)
			return
		}
		switch eventType {
		case 0:
			allTimePoints[timePoint].Events.CensoringEvents = counts
			break
		case 1:
			allTimePoints[timePoint].Events.EventsOfInterest = counts
			break
		default:
			err = fmt.Errorf("Unexpected envent type code %d, must be either 0 (event of interest) or 1 (censoring event)", eventType)
			return
		}
	}

	// filter out empty time point

	timePoints = make(medcomodels.TimePoints, 0)
	for _, timePoint := range allTimePoints {
		if (timePoint.Events.CensoringEvents != 0) || (timePoint.Events.EventsOfInterest != 0) {
			logrus.Tracef("New timepoint added %+v", timePoint)
			timePoints = append(timePoints, timePoint)
		}
	}

	return

}
