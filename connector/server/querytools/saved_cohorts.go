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
	logrus.Debugf("selecting user ID %s, cohort name ID %s", userID, cohortName)
	logrus.Debugf("SQL: %s", getPatientList)
	row := utilserver.DBConnection.QueryRow(getPatientList, userID, cohortName)
	patientNumsString := new(string)
	err = row.Scan(patientNumsString)
	if err != nil {
		err = fmt.Errorf("while scanning SQL record: %s", err.Error())
		return
	}
	logrus.Debug("successfully selected")
	var pNum int64
	logrus.Tracef("Got response %s", *patientNumsString)
	patientListString := strings.Trim(*patientNumsString, "{}")
	if patientListString == "" {
		logrus.Debug("empty patient list")
		return
	}
	for _, pID := range strings.Split(patientListString, ",") {

		pNum, err = strconv.ParseInt(pID, 10, 64)
		if err != nil {
			err = fmt.Errorf("while parsing patient ID \"%s\": %s", pID, err.Error())
			return
		}
		patientNums = append(patientNums, pNum)
	}

	return
}

// GetI2b2NonEncryptedSetID runs a SQL query on db and returns the list of patient IDs for given queryID and userID
func GetI2b2NonEncryptedSetID(userID string, cohortName string) (i2b2SetID int64, err error) {
	logrus.Debugf("selecting user ID %s, cohort name ID %s", userID, cohortName)
	logrus.Debugf("SQL: %s", getI2b2NonEncryptedSetID)
	row := utilserver.DBConnection.QueryRow(getI2b2NonEncryptedSetID, userID, cohortName)
	i2b2SetNumString := new(string)
	err = row.Scan(i2b2SetNumString)
	if err != nil {
		err = fmt.Errorf("while scanning SQL record: %s", err.Error())
		return
	}
	logrus.Debug("successfully selected")
	logrus.Tracef("Got response %s", *i2b2SetNumString)

	i2b2SetID, err = strconv.ParseInt(*i2b2SetNumString, 10, 64)
	if err != nil {
		err = fmt.Errorf("while parsing patient ID \"%s\": %s", *i2b2SetNumString, err.Error())
		return
	}
	return
}

// GetSavedCohorts runs a SQL query on db and returns the list of saved cohorts for given queryID and userID
func GetSavedCohorts(userID string, limit int) ([]medcomodels.Cohort, error) {
	var rows *sql.Rows
	var err error
	if limit > 0 {
		logrus.Debugf("selecting user ID %s, limit %d", userID, limit)
		logrus.Debugf("SQL: %s", getCohorts)
		rows, err = utilserver.DBConnection.Query(getCohorts, userID, limit)

	} else {
		logrus.Debugf("selecting user ID %s", userID)
		logrus.Debugf("SQL: %s", getCohortsNoLimit)
		rows, err = utilserver.DBConnection.Query(getCohortsNoLimit, userID)

	}

	if err != nil {
		err = fmt.Errorf("while executing SQL: %s", err.Error())
		return nil, err
	}
	logrus.Debug("successfully selected")
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
			err = fmt.Errorf("while scanning SQL record: %s", err.Error())
			return nil, err
		}
		createDate, err = time.Parse(time.RFC3339, createDateString)
		if err != nil {
			err = fmt.Errorf("while parsing create date string \"%s\": %s", createDateString, err.Error())
			return nil, err
		}
		updateDate, err = time.Parse(time.RFC3339, updateDateString)
		if err != nil {
			err = fmt.Errorf("while parsing update date string \"%s\": %s", updateDateString, err.Error())
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
		err = fmt.Errorf("while closing SQL record stream: %s", err.Error())
		return nil, err
	}

	logrus.Debugf("got %d cohorts", len(cohorts))
	return cohorts, nil
}

// GetDate runs a SQL query on db and returns the update date of cohort corresponding to  cohortID
func GetDate(userID string, cohortID int) (time.Time, error) {
	logrus.Debugf("selecting user ID %s, cohort ID %d", userID, cohortID)
	logrus.Debugf("SQL: %s", getDate)
	row := utilserver.DBConnection.QueryRow(getDate, userID, cohortID)
	timeString := new(string)
	err := row.Scan(timeString)
	if err != nil {
		err = fmt.Errorf("while scanning SQL record: %s", err.Error())
		return time.Now(), err
	}
	logrus.Debug("successfully selected")

	timeParsed, err := time.Parse(time.RFC3339, *timeString)

	return timeParsed, err

}

// InsertCohort runs a SQL query to either insert a new cohort or update an existing one
func InsertCohort(userID string, queryID int, cohortName string, createDate, updateDate time.Time) (int, error) {
	logrus.Debugf("inserting %s, query ID %d, cohort name %s, create date %s, update date %s", userID, queryID, cohortName, createDate.Format(time.RFC3339), updateDate.Format(time.RFC3339))
	logrus.Debugf("SQL: %s", insertCohort)
	row := utilserver.DBConnection.QueryRow(insertCohort, userID, queryID, cohortName, createDate, updateDate)
	res := new(string)
	err := row.Scan(res)
	if err != nil {
		err = fmt.Errorf("while scanning SQL record: %s", err.Error())
		return -1, err
	}
	logrus.Debug("successfully inserted")
	cohortID, err := strconv.Atoi(*res)
	if err != nil {
		err = fmt.Errorf("while parsing integer string \"%s\": %s", *res, err.Error())
		return -1, err
	}
	return cohortID, err
}

// UpdateCohort runs a SQL query to either insert a new cohort or update an existing one
func UpdateCohort(cohortName, userID string, queryID int, updateDate time.Time) (int, error) {
	logrus.Debugf("updating user ID %s, cohort name %s, query ID %d, update time %s", userID, cohortName, queryID, updateDate.Format(time.RFC3339))
	logrus.Debugf("SQL: %s", updateCohort)
	row := utilserver.DBConnection.QueryRow(updateCohort, cohortName, userID, queryID, updateDate)
	res := new(string)
	err := row.Scan(res)
	if err != nil {
		err = fmt.Errorf("while scanning SQL record: %s", err.Error())
		return -1, err
	}
	logrus.Debug("successfully updated")
	cohortID, err := strconv.Atoi(*res)
	if err != nil {
		err = fmt.Errorf("while parsing integer string \"%s\": %s", *res, err.Error())
		return -1, err
	}
	return cohortID, err
}

// DoesCohortExist check whether a cohort exists for provided user ID and a cohort name.
func DoesCohortExist(userID, cohortName string) (bool, error) {
	logrus.Debugf("selecting user ID %s, cohort name %s", userID, cohortName)
	logrus.Debugf("SQL: %s", doesCohortExist)
	row := utilserver.DBConnection.QueryRow(doesCohortExist, userID, cohortName)
	res := new(string)
	err := row.Scan(res)
	if err != nil {
		err = fmt.Errorf("while scanning SQL record: %s", err.Error())
		return false, err
	}
	logrus.Debug("successfully selected")
	cohortNumber, err := strconv.Atoi(*res)
	return cohortNumber > 0, err
}

// RemoveCohort deletes cohort
func RemoveCohort(userID, cohortName string) error {
	logrus.Debugf("deleting user ID %s, cohort name %s", userID, cohortName)
	logrus.Debugf("SQL: %s", removeCohort)
	_, err := utilserver.DBConnection.Exec(removeCohort, userID, cohortName)
	if err != nil {
		err = fmt.Errorf("while executing SQL: %s", err.Error())
		return err
	}
	logrus.Debug("successfully deleted")
	return nil
}

const getPatientList string = `
SELECT clear_result_set FROM query_tools.explore_query_results
WHERE query_id = (SELECT query_id FROM query_tools.saved_cohorts WHERE user_id = $1 AND cohort_name = $2 AND query_status = 'completed');
`

const getI2b2NonEncryptedSetID string = `
SELECT i2b2_non_encrypted_patient_set_id FROM query_tools.explore_query_results
WHERE query_id = (SELECT query_id FROM query_tools.saved_cohorts WHERE user_id = $1 AND cohort_name = $2 AND query_status = 'completed');
`

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
