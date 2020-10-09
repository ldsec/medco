package survivalserver

import (
	"database/sql"
	"strconv"
	"strings"

	utilcommon "github.com/ldsec/medco/connector/util/common"

	"github.com/sirupsen/logrus"
)

// buildTimePoints execute a SQL query that returns event counts per time point, for given input patient set, start and end  concept codes and modifiers
func buildTimePoints(db *sql.DB, patientList []int64, startConceptCode string, startConceptModifier string, endConceptCode string, endConceptModifier string, timeLimit int) (timePoints utilcommon.TimePoints, err error) {
	logrus.Debug("SQL query : " + sql6)
	pList := make([]string, len(patientList))
	for i, pNum := range patientList {
		pList[i] = strconv.FormatInt(pNum, 10)
	}
	patients := "{" + strings.Join(pList, ",") + "}"
	rows, err := db.Query(sql6, startConceptCode, startConceptModifier, patients, endConceptCode, endConceptModifier, timeLimit)
	if err != nil {
		return
	}
	timePointString := new(string)
	eventsString := new(string)
	censoringString := new(string)
	for rows.Next() {
		sqlTimePoint := utilcommon.TimePoint{}
		scanErr := rows.Scan(timePointString, eventsString, censoringString)
		if scanErr != nil {
			err = scanErr
			return
		}
		sqlTimePoint.Events.EventsOfInterest, err = strconv.ParseInt(*eventsString, 10, 64)
		if err != nil {

			return
		}
		sqlTimePoint.Events.CensoringEvents, err = strconv.ParseInt(*censoringString, 10, 64)
		if err != nil {

			return
		}
		sqlTimePoint.Time, err = strconv.Atoi(*timePointString)
		if err != nil {

			return
		}
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

const sql1 string = `
SELECT patient_num,start_date 
FROM i2b2demodata_i2b2.observation_fact
WHERE concept_cd = $1 and modifier_cd = $2 and patient_num = ANY($3::integer[])
`
const sql2 string = `
SELECT patient_num,end_date
FROM i2b2demodata_i2b2.observation_fact
WHERE concept_cd = $4 and modifier_cd = $5 and patient_num = ANY($3::integer[])
`
const sql3 string = `
SELECT DATE_PART('day',end_date::timestamp - start_date::timestamp) AS timepoint, COUNT(*) AS event_count
FROM (` + sql1 + `) AS x
INNER JOIN  (` + sql2 + `) AS y
ON x.patient_num = y.patient_num
GROUP BY timepoint
`

const sql4 string = `
SELECT patient_num, MAX(end_date) AS end_date
FROM i2b2demodata_i2b2.observation_fact
WHERE patient_num = ANY($3::integer[]) AND patient_num NOT IN (SELECT patient_num FROM (` + sql2 + `) AS patients_with_events)
GROUP BY patient_num
`

const sql5 string = `
SELECT * FROM (SELECT DATE_PART('day', end_date::timestamp - start_date::timestamp) AS timepoint, COUNT(*) AS censoring_count
FROM (` + sql4 + `) AS x
INNER JOIN  (` + sql1 + `) AS y
ON (x.patient_num = y.patient_num)
GROUP BY timepoint) AS z
WHERE timepoint < $6
`

const sql6 string = `
SELECT COALESCE(xx.timepoint,yy.timepoint) AS timepoint , COALESCE(event_count,0) AS event_count, COALESCE(censoring_count,0) AS censoring_count FROM (` + sql3 + `) AS xx  FULL JOIN (` + sql5 + `) AS yy
ON xx.timepoint = yy.timepoint
`
