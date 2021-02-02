package querytoolsserver

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	medcomodels "github.com/ldsec/medco/connector/models"

	utilserver "github.com/ldsec/medco/connector/util/server"

	"github.com/sirupsen/logrus"
)

// GetPatientList runs a SQL query on db and returns the list of patient IDs for given queryID and userID
func GetPatientList(userID string, cohortName string) (patientNums []int64, err error) {
	description := fmt.Sprintf("GetPatientList (ID %s, cohort name ID %s), SQL: %s", userID, cohortName, getPatientList)
	logrus.Debugf("running: %s", description)
	row := utilserver.DBConnection.QueryRow(getPatientList, userID, cohortName)
	patientNumsString := new(string)
	err = row.Scan(patientNumsString)
	if err != nil {
		err = fmt.Errorf("while scanning SQL record: %s, DB operation: %s", err.Error(), description)
		return
	}
	var pNum int64
	logrus.Tracef("Got response %s", *patientNumsString)
	patientListString := strings.Trim(*patientNumsString, "{}")
	if patientListString == "" {
		logrus.Debugf("empty patient list, DB operation: %s", description)
		return
	}
	for _, pID := range strings.Split(patientListString, ",") {

		pNum, err = strconv.ParseInt(pID, 10, 64)
		if err != nil {
			err = fmt.Errorf("while parsing patient ID \"%s\": %s, DB operation: %s", pID, err.Error(), description)
			return
		}
		patientNums = append(patientNums, pNum)
	}

	logrus.Debugf("successfully retrieved %d patients, DB operation: %s", len(patientNums), description)
	return
}

// GetSavedCohorts runs a SQL query on db and returns the list of saved cohorts for given queryID and userID
func GetSavedCohorts(userID string, limit int) ([]medcomodels.Cohort, error) {
	var description string
	var rows *sql.Rows
	var err error
	if limit > 0 {
		description = fmt.Sprintf("GetSavedCohorts(user ID %s, limit %d), SQL: %s", userID, limit, getCohorts)
		logrus.Debugf("running: %s", description)
		rows, err = utilserver.DBConnection.Query(getCohorts, userID, limit)

	} else {
		description = fmt.Sprintf("GetSavedCohorts(user ID %s), SQL: %s", userID, getCohortsNoLimit)
		logrus.Debugf("running: %s", description)
		rows, err = utilserver.DBConnection.Query(getCohortsNoLimit, userID)

	}

	if err != nil {
		err = fmt.Errorf("while executing SQL: %s, DB operation: %s", err.Error(), description)
		return nil, err
	}
	var id int
	var qid int
	var name string
	var createDateString string
	var createDate time.Time
	var updateDateString string
	var updateDate time.Time
	var cohorts = make([]medcomodels.Cohort, 0)
	for rows.Next() {
		err = rows.Scan(&id, &qid, &name, &createDateString, &updateDateString)
		if err != nil {
			err = fmt.Errorf("while scanning SQL record: %s, DB operation: %s", err.Error(), description)
			return nil, err
		}
		createDate, err = time.Parse(time.RFC3339, createDateString)
		if err != nil {
			err = fmt.Errorf("while parsing create date string \"%s\": %s, DB operation: %s", createDateString, err.Error(), description)
			return nil, err
		}
		updateDate, err = time.Parse(time.RFC3339, updateDateString)
		if err != nil {
			err = fmt.Errorf("while parsing update date string \"%s\": %s, DB operation: %s", updateDateString, err.Error(), description)
			return nil, err
		}
		newCohort := medcomodels.Cohort{
			CohortID:     id,
			QueryID:      qid,
			CohortName:   name,
			CreationDate: createDate,
			UpdateDate:   updateDate,
		}
		logrus.Tracef("got cohort %+v", newCohort)
		cohorts = append(cohorts, newCohort)
	}
	err = rows.Close()
	if err != nil {
		err = fmt.Errorf("while closing SQL record stream: %s, DB operation :%s", err.Error(), description)
		return nil, err
	}

	logrus.Debugf("successfully retrieved %d cohorts, DB operation: %s", len(cohorts), description)
	return cohorts, nil
}

// GetDate runs a SQL query on db and returns the update date of cohort corresponding to  cohortID
func GetDate(userID string, cohortID int) (time.Time, error) {
	description := fmt.Sprintf("GetDate (user ID %s, cohort ID %d), SQL: %s", userID, cohortID, getDate)
	logrus.Debugf("running: %s", description)
	row := utilserver.DBConnection.QueryRow(getDate, userID, cohortID)
	timeString := new(string)
	err := row.Scan(timeString)
	if err != nil {
		err = fmt.Errorf("while scanning SQL record: %s, DB operation: %s", err.Error(), description)
		return time.Now(), err
	}

	timeParsed, err := time.Parse(time.RFC3339, *timeString)

	if err != nil {
		err = fmt.Errorf("while parsing time string: %s, DB operation: %s", err.Error(), description)
		return time.Now(), err
	}

	logrus.Debugf("successfully retrieved date: %s, DB operation: %s", *timeString, description)

	return timeParsed, nil

}

// InsertCohort runs a SQL query to either insert a new cohort or update an existing one
func InsertCohort(userID string, queryID int, cohortName string, createDate, updateDate time.Time) (int, error) {
	description := fmt.Sprintf(
		"InsertCohort (user ID: %s, query ID: %d, cohort name: %s, create date: %s, update date: %s), SQL: %s",
		userID,
		queryID,
		cohortName,
		createDate.Format(time.RFC3339),
		updateDate.Format(time.RFC3339),
		insertCohort,
	)
	logrus.Debugf("running: %s", description)
	row := utilserver.DBConnection.QueryRow(insertCohort, userID, queryID, cohortName, createDate, updateDate)
	res := new(string)
	err := row.Scan(res)
	if err != nil {
		err = fmt.Errorf("while scanning SQL record: %s, DB operation: %s", err.Error(), description)
		return -1, err
	}
	cohortID, err := strconv.Atoi(*res)
	if err != nil {
		err = fmt.Errorf("while parsing integer string \"%s\": %s, DB operation: %s", *res, err.Error(), description)
		return -1, err
	}
	logrus.Debugf("successfully inserted cohort, cohort ID: %d, DB operation: %s", cohortID, description)

	return cohortID, err
}

// UpdateCohort runs a SQL query to either insert a new cohort or update an existing one
func UpdateCohort(cohortName, userID string, queryID int, updateDate time.Time) (int, error) {
	description := fmt.Sprintf("UpdateCohort (cohort name: %s, user ID: %s, query ID: %d, update time: %s), SQL: %s", cohortName, userID, queryID, updateDate.Format(time.RFC3339), updateCohort)
	logrus.Debugf("running: %s", description)
	row := utilserver.DBConnection.QueryRow(updateCohort, cohortName, userID, queryID, updateDate)
	res := new(string)
	err := row.Scan(res)
	if err != nil {
		err = fmt.Errorf("while scanning SQL record: %s", err.Error())
		return -1, err
	}
	cohortID, err := strconv.Atoi(*res)
	if err != nil {
		err = fmt.Errorf("while parsing integer string \"%s\": %s, DB operation: %s", *res, err.Error(), description)
		return -1, err
	}
	logrus.Debugf("successful cohort update, cohort ID: %d, DB operation: %s", cohortID, description)

	return cohortID, err
}

// DoesCohortExist check whether a cohort exists for provided user ID and a cohort name.
func DoesCohortExist(userID, cohortName string) (bool, error) {
	description := fmt.Sprintf("DoesCohortExist (user ID: %s, cohort name: %s), SQL: %s", userID, cohortName, doesCohortExist)
	logrus.Debugf("running: %s", description)
	row := utilserver.DBConnection.QueryRow(doesCohortExist, userID, cohortName)
	res := new(string)
	err := row.Scan(res)
	if err != nil {
		err = fmt.Errorf("while scanning SQL record: %s, DB operation: %s", err.Error(), description)
		return false, err
	}

	cohortNumber, err := strconv.Atoi(*res)
	retValue := cohortNumber > 0
	logrus.Debugf("successful cohort existence check: %t, DB operation: %s", retValue, description)
	return retValue, err
}

// RemoveCohort deletes cohort
func RemoveCohort(userID, cohortName string) error {
	description := fmt.Sprintf("RemoveCohort (deleting user ID: %s, cohort name: %s), SQL: %s", userID, cohortName, removeCohort)
	logrus.Debugf("running: %s", description)
	_, err := utilserver.DBConnection.Exec(removeCohort, userID, cohortName)
	if err != nil {
		err = fmt.Errorf("while executing SQL: %s, DB operation: %s", err.Error(), description)
		return err
	}
	logrus.Debugf("successfully deleted, DB operation: %s", description)
	return nil
}

const insertCohort string = `
INSERT INTO query_tools.saved_cohorts(user_id,query_id,cohort_name,create_date,update_date)
VALUES ($1,$2,$3,$4,$5)
RETURNING cohort_id
`

const updateCohort string = `
UPDATE query_tools.saved_cohorts
SET query_id=$3, update_date= $4
WHERE cohort_name = $1 AND user_id = $2
RETURNING cohort_id
`

const getCohorts string = `
SELECT cohort_id, query_id, cohort_name, create_date, update_date FROM query_tools.saved_cohorts
WHERE user_id = $1
ORDER BY cohort_name
LIMIT $2
`

const getCohortsNoLimit string = `
SELECT cohort_id, query_id, cohort_name, create_date, update_date FROM query_tools.saved_cohorts
WHERE user_id = $1
`

const getDate string = `
SELECT update_date FROM query_tools.saved_cohorts
WHERE user_id =$1 and cohort_id=$2
`

const doesCohortExist string = `
SELECT COUNT(cohort_id) FROM query_tools.saved_cohorts
WHERE user_id = $1 and cohort_name = $2
`

const removeCohort string = `
DELETE FROM query_tools.saved_cohorts
WHERE user_id = $1 AND cohort_name = $2
`
