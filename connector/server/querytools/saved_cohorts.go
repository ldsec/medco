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
	description := fmt.Sprintf("GetPatientList (ID %s, cohort name ID %s), procedure: %s", userID, cohortName, "query_tools.get_patient_list")
	logrus.Debugf("running: %s", description)
	row := utilserver.DBConnection.QueryRow("SELECT query_tools.get_patient_list($1 ,$2);", userID, cohortName)
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

	description = fmt.Sprintf("GetSavedCohorts(user ID %s, limit %d), procedure: %s", userID, limit, "query_tools.get_cohorts")
	logrus.Debugf("running: %s", description)
	rows, err = utilserver.DBConnection.Query("SELECT query_tools.get_cohorts($1, $2);", userID, limit)

	if err != nil {
		err = fmt.Errorf("while executing procedure: %s, DB operation: %s", err.Error(), description)
		return nil, err
	}
	var id int
	var qid int
	var name string
	var createDate time.Time
	var updateDate time.Time
	var cohorts = make([]medcomodels.Cohort, 0)
	record := new(string)
	for rows.Next() {
		err = rows.Scan(record)
		if err != nil {
			err = fmt.Errorf("while scanning SQL record: %s, DB operation: %s", err.Error(), description)
			return nil, err
		}

		cells := strings.Split(strings.Trim(*record, "()"), ",")

		id, err = strconv.Atoi(cells[0])
		if err != nil {
			err = fmt.Errorf("while parsing cohort ID string \"%s\": %s, DB operation: %s", cells[0], err.Error(), description)
			return nil, err
		}
		name = cells[2]
		qid, err = strconv.Atoi(cells[1])
		if err != nil {
			err = fmt.Errorf("while parsing query ID string \"%s\": %s, DB operation: %s", cells[1], err.Error(), description)
			return nil, err
		}

		createDate, err = time.Parse("2006-01-02 15:04:05", strings.Trim(cells[3], `"`))
		if err != nil {
			err = fmt.Errorf("while parsing create date string \"%s\": %s, DB operation: %s", cells[3], err.Error(), description)
			return nil, err
		}
		updateDate, err = time.Parse("2006-01-02 15:04:05", strings.Trim(cells[4], `"`))
		if err != nil {
			err = fmt.Errorf("while parsing update date string \"%s\": %s, DB operation: %s", cells[4], err.Error(), description)
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
	description := fmt.Sprintf("GetDate (user ID %s, cohort ID %d), procedure: %s", userID, cohortID, "query_tools.get_date")
	logrus.Debugf("running: %s", description)
	row := utilserver.DBConnection.QueryRow("SELECT query_tools.get_date($1, $2);", userID, cohortID)
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
		"InsertCohort (user ID: %s, query ID: %d, cohort name: %s, create date: %s, update date: %s), procedure: %s",
		userID,
		queryID,
		cohortName,
		createDate.Format(time.RFC3339),
		updateDate.Format(time.RFC3339),
		"query_tools.insert_cohort",
	)
	logrus.Debugf("running: %s", description)
	row := utilserver.DBConnection.QueryRow("SELECT query_tools.insert_cohort($1, $2, $3, $4, $5)", userID, queryID, cohortName, createDate, updateDate)
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
	description := fmt.Sprintf("UpdateCohort (cohort name: %s, user ID: %s, query ID: %d, update time: %s), procedure: %s", cohortName, userID, queryID, updateDate.Format(time.RFC3339), "query_tools.update_cohort")
	logrus.Debugf("running: %s", description)
	row := utilserver.DBConnection.QueryRow("SELECT query_tools.update_cohort($1, $2, $3, $4);", cohortName, userID, queryID, updateDate)
	var cohortID sql.NullInt32
	err := row.Scan(&cohortID)

	if err != nil {
		err = fmt.Errorf("during cohort update: %s, DB operation: %s", err.Error(), description)
		return -1, err
	}
	if !cohortID.Valid {
		err = fmt.Errorf("nothing was updated, DB operation: %s", description)
		return -1, err
	}
	logrus.Debugf("successful cohort update, cohort ID: %d, DB operation: %s", cohortID.Int32, description)

	return int(cohortID.Int32), err
}

// DoesCohortExist check whether a cohort exists for provided user ID and a cohort name.
func DoesCohortExist(userID, cohortName string) (bool, error) {
	description := fmt.Sprintf("DoesCohortExist (user ID: %s, cohort name: %s), procedure: %s", userID, cohortName, "does_cohort_exist")
	logrus.Debugf("running: %s", description)
	row := utilserver.DBConnection.QueryRow("SELECT query_tools.does_cohort_exist($1, $2)", userID, cohortName)
	res := new(bool)
	err := row.Scan(res)
	if err != nil {
		err = fmt.Errorf("while scanning SQL record: %s, DB operation: %s", err.Error(), description)
		return false, err
	}

	logrus.Debugf("successful cohort existence check: %t, DB operation: %s", *res, description)
	return *res, err
}

// RemoveCohort deletes cohort
func RemoveCohort(userID, cohortName string) error {
	description := fmt.Sprintf("RemoveCohort (deleting user ID: %s, cohort name: %s), procedure: %s", userID, cohortName, "query_tools.remove_cohort")
	logrus.Debugf("running: %s", description)
	res := utilserver.DBConnection.QueryRow("SELECT query_tools.remove_cohort($1, $2);", userID, cohortName)
	var cohortID sql.NullInt32
	err := res.Scan(&cohortID)

	if err != nil {
		err = fmt.Errorf("while executing procedure: %s, DB operation: %s", err.Error(), description)
		return err
	}

	if !cohortID.Valid {
		err = fmt.Errorf("cohort to be removed was not found, DB operation: %s", description)
		return err
	}
	logrus.Debugf("successfully deleted, DB operation: %s", description)
	return nil
}
