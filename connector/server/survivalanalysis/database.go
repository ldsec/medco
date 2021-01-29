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
	startConceptEarliest bool,
	endConceptCodes []string,
	endConceptModifiers []string,
	endConceptEarliest bool,
	timeLimit int,
) (timePoints medcomodels.TimePoints, err error) {

	pList := make([]string, len(patientList))
	for i, pNum := range patientList {
		pList[i] = strconv.FormatInt(pNum, 10)
	}
	sql6Instance := sql6(startConceptEarliest, endConceptEarliest)
	patients := "{" + strings.Join(pList, ",") + "}"
	startConcepts := "{" + strings.Join(startConceptCodes, ",") + "}"
	endConcepts := "{" + strings.Join(endConceptCodes, ",") + "}"
	startModifiers := "{" + strings.Join(startConceptModifiers, ",") + "}"
	endModifiers := "{" + strings.Join(endConceptModifiers, ",") + "}"
	description := fmt.Sprintf("selecting start concept code %s, start concept modifier %s, patients list %s, end concept code %s, end concept modifier %s, time limit %d: SQL %s", startConcepts, startModifiers, patients, endConcepts, endModifiers, timeLimit, sql6Instance)
	logrus.Debugf("running: %s", description)
	rows, err := utilserver.I2B2DBConnection.Query(sql6Instance, startConcepts, startModifiers, patients, endConcepts, endModifiers, timeLimit)
	if err != nil {
		err = fmt.Errorf("while execution SQL query: %s, DB operation: %s", err.Error(), description)
		return
	}
	logrus.Debug("successfully selected")
	timePointString := new(string)
	eventsString := new(string)
	censoringString := new(string)
	for rows.Next() {
		sqlTimePoint := medcomodels.TimePoint{}
		scanErr := rows.Scan(timePointString, eventsString, censoringString)
		if scanErr != nil {
			err = scanErr
			err = fmt.Errorf("while scanning SQL record: %s", err.Error())
			return
		}
		sqlTimePoint.Events.EventsOfInterest, err = strconv.ParseInt(*eventsString, 10, 64)
		if err != nil {
			err = fmt.Errorf("while scanning parsing integer string (number of events) \"%s\": %s", *eventsString, err.Error())
			return
		}
		sqlTimePoint.Events.CensoringEvents, err = strconv.ParseInt(*censoringString, 10, 64)
		if err != nil {
			err = fmt.Errorf("while scanning parsing integer string (number of censoring) \"%s\": %s", *censoringString, err.Error())
			return
		}
		sqlTimePoint.Time, err = strconv.Atoi(*timePointString)
		if err != nil {
			err = fmt.Errorf("while scanning parsing integer string (relative time) \"%s\": %s", *timePointString, err.Error())
			return
		}
		logrus.Tracef("new time point: %+v", sqlTimePoint)
		timePoints = append(timePoints, sqlTimePoint)
	}
	return

}

/*

Recap of arguments:

$1 start event concept code

$2 start event modifier code

$3 list of patient in cohort or sub group

$4 end event concept code

$5 end event modifier code

$6 max time limit in day

*/

// TODO find better names or com
// prepare those functions in DB

const sql1Earliest string = `
SELECT patient_num, MIN(start_date) AS start_date
FROM i2b2demodata_i2b2.observation_fact
WHERE concept_cd = ANY ($1::varchar[]) and modifier_cd = ANY ($2::varchar[]) and patient_num = ANY($3::integer[])
GROUP BY patient_num
`
const sql1Latest string = `
SELECT patient_num, MAX(start_date) AS start_date 
FROM i2b2demodata_i2b2.observation_fact
WHERE concept_cd = ANY ($1::varchar[]) and modifier_cd = ANY ($2::varchar[]) and patient_num = ANY($3::integer[])
GROUP BY patient_num
`
const sql2Earliest string = `
SELECT patient_num, MIN(end_date) AS end_date
FROM i2b2demodata_i2b2.observation_fact
WHERE concept_cd = ANY ($4::varchar[]) and modifier_cd = ANY  ($5::varchar[]) and patient_num = ANY($3::integer[])
GROUP BY patient_num
`
const sql2Latest string = `
SELECT patient_num, MAX(end_date) AS end_date
FROM i2b2demodata_i2b2.observation_fact
WHERE concept_cd = ANY ($4::varchar[]) and modifier_cd = ANY ($5::varchar[]) and patient_num = ANY($3::integer[])
GROUP BY patient_num
`

func sql3(startConceptEarliest, endConceptEarliest bool) string {
	var sql1, sql2 string
	if startConceptEarliest {
		sql1 = sql1Earliest
	} else {
		sql1 = sql1Latest
	}

	if endConceptEarliest {
		sql2 = sql2Earliest
	} else {
		sql2 = sql2Latest
	}

	return `
	SELECT DATE_PART('day',end_date::timestamp - start_date::timestamp) AS timepoint, COUNT(*) AS event_count
	FROM (` + sql1 + `) AS x
	INNER JOIN  (` + sql2 + `) AS y
	ON x.patient_num = y.patient_num
	GROUP BY timepoint
	`
}

func sql4(endConceptEarliest bool) string {
	var sql2 string
	if endConceptEarliest {
		sql2 = sql2Earliest
	} else {
		sql2 = sql2Latest
	}
	return `
SELECT patient_num, MAX(end_date) AS end_date
FROM i2b2demodata_i2b2.observation_fact
WHERE patient_num = ANY($3::integer[]) AND patient_num NOT IN (SELECT patient_num FROM (` + sql2 + `) AS patients_with_events)
GROUP BY patient_num
`
}

func sql5(startConceptEarliest, endConceptEarliest bool) string {
	var sql1 string
	if startConceptEarliest {
		sql1 = sql1Earliest
	} else {
		sql1 = sql1Latest
	}
	return `
SELECT * FROM (SELECT DATE_PART('day', end_date::timestamp - start_date::timestamp) AS timepoint, COUNT(*) AS censoring_count
FROM (` + sql4(endConceptEarliest) + `) AS x
INNER JOIN  (` + sql1 + `) AS y
ON (x.patient_num = y.patient_num)
GROUP BY timepoint) AS z
WHERE timepoint < $6
`
}

func sql6(startConceptEarliest, endConceptEarliest bool) string {
	return `
SELECT COALESCE(xx.timepoint,yy.timepoint) AS timepoint , COALESCE(event_count,0) AS event_count, COALESCE(censoring_count,0) AS censoring_count FROM (` + sql3(startConceptEarliest, endConceptEarliest) + `) AS xx  FULL JOIN (` + sql5(startConceptEarliest, endConceptEarliest) + `) AS yy
ON xx.timepoint = yy.timepoint
`
}
