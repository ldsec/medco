package querytoolsserver

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	utilcommon "github.com/ldsec/medco-connector/util/common"

	"github.com/sirupsen/logrus"
)

// GetPatientList runs a SQL query on db and returns the list of patient IDs for given queryID and userID
func GetPatientList(db *sql.DB, userID string, resultInstanceID int64) (patientNums []int64, err error) {
	row := db.QueryRow(getPatientList, userID, resultInstanceID)
	patientNumsString := new(string)
	err = row.Scan(patientNumsString)
	var pNum int64
	logrus.Tracef("Got response %s", *patientNumsString)
	for _, pID := range strings.Split(strings.Trim(*patientNumsString, "{}"), ",") {

		pNum, err = strconv.ParseInt(pID, 10, 64)
		if err != nil {
			return
		}
		patientNums = append(patientNums, pNum)
	}

	return
}

// GetSavedCohorts runs a SQL query on db and returns the list of saved cohorts for given queryID and userID
func GetSavedCohorts(db *sql.DB, userID string) ([]utilcommon.Cohort, error) {
	rows, err := db.Query(getCohorts, userID)
	if err != nil {
		return nil, err
	}
	var id int
	var qid int
	var name string
	var createDateString string
	var createDate time.Time
	var updateDateString string
	var updateDate time.Time
	var cohorts = make([]utilcommon.Cohort, 0)
	for rows.Next() {
		err = rows.Scan(&id, &qid, &name, &createDateString, &updateDateString)
		if err != nil {
			return nil, err
		}
		createDate, err = time.Parse(time.RFC3339, createDateString)
		if err != nil {
			return nil, err
		}
		updateDate, err = time.Parse(time.RFC3339, updateDateString)
		if err != nil {
			return nil, err
		}
		cohorts = append(cohorts, utilcommon.Cohort{
			CohortID:     id,
			QueryID:      qid,
			CohortName:   name,
			CreationDate: createDate,
			UpdateDate:   updateDate,
		})
	}
	err = rows.Close()
	if err != nil {
		return nil, err
	}

	logrus.Infof("Got %d cohorts", len(cohorts))
	return cohorts, nil
}

// GetDate runs a SQL query on db and returns the update date of cohort corresponding to  cohortID
func GetDate(db *sql.DB, userID string, cohortID int) (time.Time, error) {
	row := db.QueryRow(getDate, userID, cohortID)
	timeString := new(string)
	err := row.Scan(timeString)
	if err != nil {
		return time.Now(), err
	}

	timeParsed, err := time.Parse(time.RFC3339, *timeString)

	return timeParsed, err

}

// InsertCohort runs a SQL query to either insert a new cohort or update an existing one
func InsertCohort(db *sql.DB, userID string, queryID int, cohortName string, createDate, updateDate time.Time) (int, error) {
	row := db.QueryRow(insertCohort, userID, queryID, cohortName, createDate, updateDate)
	res := new(string)
	err := row.Scan(res)
	if err != nil {
		return -1, err
	}
	cohortID, err := strconv.Atoi(*res)
	return cohortID, err
}

const insertCohort string = `
INSERT INTO query_tools.saved_cohorts(user_id,query_id,cohort_name,create_date,update_date)
VALUES ($1,$2,$3,$4,$5)
ON CONFLICT (user_id,cohort_name) DO UPDATE SET query_id = $2, update_date=$5
RETURNING cohort_id
`

const updateCohort string = `
UPDATE query_tools.saved_cohorts
SET query_id=$3, update_date= $4
WHERE cohort_id = $1 AND user_id = $2
`

const getCohorts string = `
SELECT cohort_id, query_id, cohort_name, create_date, update_date FROM query_tools.saved_cohorts
WHERE user_id = $1
`

const getDate string = `
SELECT update_date FROM query_tools.saved_cohorts
WHERE user_id =$1 and cohort_id=$2
`
