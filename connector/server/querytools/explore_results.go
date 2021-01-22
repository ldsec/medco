package querytoolsserver

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	utilserver "github.com/ldsec/medco/connector/util/server"

	"github.com/sirupsen/logrus"
)

// InsertExploreResultInstance is called when the query is created. A new row is inserted in explore query results table with status 'running'.
func InsertExploreResultInstance(userID, queryName, queryDefinition string) (int, error) {
	description := fmt.Sprintf("InsertExploreResultInstance (user ID %s, query name %s, query definition %s), SQL: %s", userID, queryName, queryDefinition, insertExploreResultInstance)
	logrus.Debugf("running: %s", description)
	row := utilserver.DBConnection.QueryRow(insertExploreResultInstance, userID, queryName, queryDefinition)
	res := new(string)
	err := row.Scan(res)
	if err != nil {
		err = fmt.Errorf("while scanning SQL record: %s, in DB operation :%s", err.Error(), description)
		return 0, err
	}
	logrus.Debugf("DB operation successful result %s, operation: %s", *res, description)
	queryID, err := strconv.Atoi(*res)
	if err != nil {
		err = fmt.Errorf("while parsing integer string \"%s\": %s, in DB operation :%s,", *res, err.Error(), description)
		return 0, err
	}
	return queryID, nil

}

// UpdateExploreResultInstance updates the instance corresponding to the given queryID. Its status is changed to 'completed'.
// I2b2 encrypted patient set ID should be the result instance id returned by I2B2 after a successful query in project for encrypted data.
// I2b2 non encrypted patient set ID should be the result instance id returned by I2B2 after a successful query in project for non encrypted data.
// Null pointer for i2b2 (non) encrypted patient set ID is used to update NULL value in the table.
func UpdateExploreResultInstance(queryID int, clearResultSetSize int, clearResultSet []int, i2b2EncryptedPatientSetID, i2b2NonEncryptedPatientSetID *int) error {
	var description string
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
		description = fmt.Sprintf("UpdateExploreResultInstanceOnlyClear (query ID: %d, clear patient set size %d, set definition: %s,clear patient set %d): SQL: %s", queryID, clearResultSetSize, setDefinition, *i2b2NonEncryptedPatientSetID, updateExploreResultInstanceOnlyClear)
		logrus.Debugf("running: %s", description)
		res, err = utilserver.DBConnection.Exec(updateExploreResultInstanceOnlyClear, queryID, clearResultSetSize, setDefinition, *i2b2NonEncryptedPatientSetID)
	} else if i2b2NonEncryptedPatientSetID == nil {
		description = fmt.Sprintf("UpdateExploreResultInstanceOnlyEncrypted (query ID: %d, clear patient set size %d, set definition: %s, encrypted patient set: %d): SQL: %s", queryID, clearResultSetSize, setDefinition, *i2b2EncryptedPatientSetID, updateExploreResultInstanceOnlyEncrypted)
		logrus.Debugf("running: %s", description)
		res, err = utilserver.DBConnection.Exec(updateExploreResultInstanceOnlyEncrypted, queryID, clearResultSetSize, setDefinition, *i2b2EncryptedPatientSetID)
	} else {
		description = fmt.Sprintf(
			"UpdateUpdateExploreResultInstanceBoth (query ID: %d, clear patient set size %d, set definition: %s, clear patient set: %d, encrypted patient set: %d), SQL: %s",
			queryID,
			clearResultSetSize,
			setDefinition,
			*i2b2NonEncryptedPatientSetID,
			*i2b2EncryptedPatientSetID,
			updateExploreResultInstanceOnlyEncrypted,
		)
		logrus.Debugf("running: %s", description)
		res, err = utilserver.DBConnection.Exec(updateExploreResultInstanceBoth, queryID, clearResultSetSize, setDefinition, *i2b2EncryptedPatientSetID, *i2b2NonEncryptedPatientSetID)
	}
	if err != nil {
		err = fmt.Errorf("while executing SQL: %s, DB operation: %s", err.Error(), description)
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		err = fmt.Errorf("while checking SQL row affected: %s, SQL object result: %+v, DB operation: %s", err.Error(), res, description)
		return err
	}
	if affected == 0 {
		err = fmt.Errorf("no error with DB, but nothing was updated, DB operation: %s", description)
		return err
	}

	logrus.Debugf("DB operation successful, number of affected rows %d, operation: %s", affected, description)

	return nil

}

// UpdateErrorExploreResultInstance updates the instance corresponding to the given queryID. Its status is changed to 'error'.
// UpdateErrorExploreResultInstance should be called whenever any I2B2 project throws an error while executing a query.
func UpdateErrorExploreResultInstance(queryID int) error {
	description := fmt.Sprintf("UpdateErrorExploreResultInstance (query ID: %d), SQL: %s", queryID, updateErrorExploreQueryInstance)
	logrus.Debugf("running: %s", description)
	res, err := utilserver.DBConnection.Exec(updateErrorExploreQueryInstance, queryID)
	if err != nil {
		err = fmt.Errorf("during DB operation execution: %s, DB operation: %s", err.Error(), description)
		return err
	}
	affected, err := res.RowsAffected()
	if err != nil {
		err = fmt.Errorf("while checking SQL row affected: %s, DB operation: %s", err.Error(), description)
		return err
	}
	if affected == 0 {
		err = fmt.Errorf("no error with DB, but nothing was updated, DB operation: %s", description)
		return err
	}
	logrus.Debugf("DB operation successful, number of affected rows %d, operation: %s", affected, description)

	return err
}

// CheckQueryID checks whether the user really has a query before inserting a new cohort defined by that query's id
func CheckQueryID(userID string, queryID int) (bool, error) {
	description := fmt.Sprintf("CheckQueryID (user ID: %s, query ID: %d), SQL: %s", userID, queryID, checkQueryID)
	logrus.Debugf("running: %s", description)
	row := utilserver.DBConnection.QueryRow(checkQueryID, userID, queryID)
	res := new(string)
	err := row.Scan(res)
	if err != nil {
		err = fmt.Errorf("while scanning SQL record: %s, DB operation: %s", err.Error(), description)
		return false, err
	}
	count, err := strconv.Atoi(*res)
	if err != nil {
		err = fmt.Errorf("while parsing integer string \"%s\": %s, Db operation: %s", *res, err.Error(), description)
		return false, err
	}
	retValue := count > 0

	logrus.Debugf("DB operation successful, result: %t, DB operation: %s", retValue, description)

	return retValue, err

}

// GetQueryDefinition is called when the query is created. A new row is inserted in explore query results table with status 'running'.
func GetQueryDefinition(queryID int) (string, error) {

	description := fmt.Sprintf("GetQueryDefinition (ID %d), SQL: %s", queryID, getQueryDefinition)
	logrus.Debugf("running: %s", description)

	row := utilserver.DBConnection.QueryRow(getQueryDefinition, queryID)
	res := new(string)
	err := row.Scan(res)
	if err != nil {
		err = fmt.Errorf("while scanning SQL record: %s, DB operation: %s", err.Error(), description)
		return "", err
	}

	logrus.Debugf("DB operation successful, result: %s, DB operation: %s", *res, description)

	return *res, nil

}

const getPatientList string = `
SELECT clear_result_set FROM query_tools.explore_query_results
WHERE query_id = (SELECT query_id FROM query_tools.saved_cohorts WHERE user_id = $1 AND cohort_name = $2 AND query_status = 'completed');
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

const getQueryDefinition string = `
SELECT query_definition FROM query_tools.explore_query_results
WHERE query_id = $1
`
