package querytoolsserver

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

// InsertExploreResultInstance is called when the query is created. A new row is inserted in explore query results table with status 'running'.
func InsertExploreResultInstance(db *sql.DB, userID, queryName, queryDefinition string) (int, error) {
	row := db.QueryRow(insertExploreResultInstance, userID, queryName, queryDefinition)
	res := new(string)
	err := row.Scan(res)
	if err != nil {
		return 0, err
	}
	queryID, err := strconv.Atoi(*res)
	if err != nil {
		return 0, err
	}
	return queryID, nil

}

// UpdateExploreResultInstance updates the instance corresponding to the given queryID. Its status is changed to 'completed'.
// I2b2 encrypted patient set ID should be the result instance id returned by I2B2 after a successful query in project for encrypted data.
// I2b2 non encrypted patient set ID should be the result instance id returned by I2B2 after a successful query in project for non encrypted data.
// Null pointer for i2b2 (non) encrypted patient set ID is used to update NULL value in the table.
func UpdateExploreResultInstance(db *sql.DB, queryID int, clearResultSetSize int, clearResultSet []int, i2b2EncryptedPatientSetID, i2b2NonEncryptedPatientSetID *int) error {
	var res sql.Result
	var err error
	setStrings := make([]string, len(clearResultSet))
	for i, patient := range clearResultSet {
		setStrings[i] = strconv.Itoa(patient)
	}
	setDefinition := "{" + strings.Join(setStrings, ",") + "}"
	if i2b2EncryptedPatientSetID == nil && i2b2NonEncryptedPatientSetID == nil {
		err = fmt.Errorf("I2B2 patient set is undefined for both non encrypted and encrypted projects")
		return err
	}
	if i2b2EncryptedPatientSetID == nil {
		res, err = db.Exec(updateExploreResultInstanceOnlyClear, queryID, clearResultSetSize, setDefinition, *i2b2NonEncryptedPatientSetID)
	} else if i2b2NonEncryptedPatientSetID == nil {
		res, err = db.Exec(updateExploreResultInstanceOnlyClear, queryID, clearResultSetSize, setDefinition, *i2b2EncryptedPatientSetID)
	} else {
		res, err = db.Exec(updateExploreResultInstanceBoth, queryID, clearResultSetSize, setDefinition, *i2b2EncryptedPatientSetID, *i2b2NonEncryptedPatientSetID)
	}
	if err != nil {
		return err
	}
	logrus.Tracef("sql execution result %+v", res)
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		err = fmt.Errorf("nothing updated")
		return err
	}

	return nil

}

// UpdateErrorExploreResultInstance updates the instance corresponding to the given queryID. Its status is changed to 'error'.
// UpdateErrorExploreResultInstance should be called whenever any I2B2 project throws an error while executing a query.
func UpdateErrorExploreResultInstance(db *sql.DB, queryID int) error {
	res, err := db.Exec(updateErrorExploreQueryInstance, queryID)
	logrus.Tracef("sql execution result %+v", res)
	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		err = fmt.Errorf("nothing updated")
	}
	return err
}

// CheckQueryID checks whether the user really has a query before inserting a new cohort defined by that query's id
func CheckQueryID(db *sql.DB, userID string, queryID int) (bool, error) {
	row := db.QueryRow(checkQueryID, userID, queryID)
	res := new(string)
	err := row.Scan(res)
	if err != nil {
		return false, err
	}
	count, err := strconv.Atoi(*res)
	if err != nil {
		return false, err
	}
	return (count > 0), err

}

const getPatientList string = `
SELECT clear_result_set FROM query_tools.explore_query_results
WHERE query_id = (SELECT query_id FROM query_tools.saved_cohorts WHERE user_id = $1 AND cohort_id = $2 AND query_status = 'completed');
`

const insertExploreResultInstance string = `
INSERT INTO query_tools.explore_query_results(user_id,query_name, query_status,query_definition)
VALUES ($1,$2,'running',$3)
RETURNING query_id
`

const updateExploreResultInstanceBoth string = `
UPDATE query_tools.explore_query_results
SET clear_result_set_size=$2, clear_result_set=$3, query_status='completed' , i2b2_encrypted_patient_set_id=$4, i2b2_non_encrypted_patient_set_id=$5
WHERE query_id = $1 AND status = 'running'
`
const updateExploreResultInstanceOnlyClear string = `
UPDATE query_tools.explore_query_results
SET clear_result_set_size=$2, clear_result_set=$3, query_status='completed' ,i2b2_non_encrypted_patient_set_id=$4
WHERE query_id = $1 AND query_status = 'running'
`
const updateExploreResultInstanceOnlyEncrypted string = `
UPDATE query_tools.explore_query_results
SET clear_result_set_size=$2, clear_result_set=$3, query_status='completed' ,i2b2_encrypted_patient_set_id=$4
WHERE query_id = $1 AND query_status = 'running'
`

const updateErrorExploreQueryInstance string = `
UPDATE query_tools.explore_query_results
SET query_status='error'
WHERE query_id = $1 AND query_status = 'running'
`

const checkQueryID string = `
SELECT COUNT(query_id) FROM query_tools.explore_query_results
WHERE user_id = $1 AND query_id = $2
`
