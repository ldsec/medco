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
	description := fmt.Sprintf("InsertExploreResultInstance (user ID %s, query name %s, query definition %s), procedure: %s", userID, queryName, queryDefinition, "query_tools.insert_explore_result_instance")
	logrus.Debugf("running: %s", description)
	row := utilserver.DBConnection.QueryRow("SELECT query_tools.insert_explore_result_instance($1, $2, $3);", userID, queryName, queryDefinition)
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
	var res *sql.Row
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
	procedure := "query_tools.update_explore_query_instance"
	if i2b2EncryptedPatientSetID == nil {
		description = fmt.Sprintf("UpdateExploreResultInstanceOnlyClear (query ID: %d, clear patient set size %d, set definition: %s,clear patient set %d): procedure: %s", queryID, clearResultSetSize, setDefinition, *i2b2NonEncryptedPatientSetID, procedure)
		logrus.Debugf("running: %s", description)
		res = utilserver.DBConnection.QueryRow("SELECT query_tools.update_explore_query_instance($1, $2, $3, $4, $5);", queryID, clearResultSetSize, setDefinition, nil, *i2b2NonEncryptedPatientSetID)
	} else if i2b2NonEncryptedPatientSetID == nil {
		description = fmt.Sprintf("UpdateExploreResultInstanceOnlyEncrypted (query ID: %d, clear patient set size %d, set definition: %s, encrypted patient set: %d): procedure: %s", queryID, clearResultSetSize, setDefinition, *i2b2EncryptedPatientSetID, procedure)
		logrus.Debugf("running: %s", description)
		res = utilserver.DBConnection.QueryRow("SELECT query_tools.update_explore_query_instance($1, $2, $3, $4, $5);", queryID, clearResultSetSize, setDefinition, *i2b2EncryptedPatientSetID, nil)
	} else {
		description = fmt.Sprintf(
			"UpdateUpdateExploreResultInstanceBoth (query ID: %d, clear patient set size %d, set definition: %s, clear patient set: %d, encrypted patient set: %d), procedure: %s",
			queryID,
			clearResultSetSize,
			setDefinition,
			*i2b2NonEncryptedPatientSetID,
			*i2b2EncryptedPatientSetID,
			procedure,
		)
		logrus.Debugf("running: %s", description)
		res = utilserver.DBConnection.QueryRow("SELECT query_tools.update_explore_query_instance($1, $2, $3, $4, $5);", queryID, clearResultSetSize, setDefinition, *i2b2EncryptedPatientSetID, *i2b2NonEncryptedPatientSetID)
	}
	modifiedQueryID := new(int)
	err = res.Scan(modifiedQueryID)

	if err == sql.ErrNoRows {
		err = fmt.Errorf("nothing was updated, DB operation: %s", description)
		return err
	}

	if err != nil {
		err = fmt.Errorf("while executing procedure: %s, DB operation: %s", err.Error(), description)
		return err
	}

	if *modifiedQueryID != queryID {
		err = fmt.Errorf(
			"ID of query instance to modify is not equal to the modified one: (input query ID: %d, output query ID: %d), DB operation: %s",
			queryID,
			*modifiedQueryID,
			description,
		)
		return err
	}

	logrus.Debugf("DB operation successful, ID of affected query instance: %d, operation: %s", *modifiedQueryID, description)

	return nil

}

// UpdateErrorExploreResultInstance updates the instance corresponding to the given queryID. Its status is changed to 'error'.
// UpdateErrorExploreResultInstance should be called whenever any I2B2 project throws an error while executing a query.
func UpdateErrorExploreResultInstance(queryID int) error {
	description := fmt.Sprintf("UpdateErrorExploreResultInstance (query ID: %d), procedure: %s", queryID, "query_tools.update_error_explore_query_instance")
	logrus.Debugf("running: %s", description)
	res := utilserver.DBConnection.QueryRow("SELECT query_tools.update_error_explore_query_instance($1);", queryID)
	modifiedQueryID := new(int)

	err := res.Scan(modifiedQueryID)
	if err == sql.ErrNoRows {
		err = fmt.Errorf("no rows were affected: %s, DB operation: %s", err.Error(), description)
		return err
	}

	if err != nil {
		err = fmt.Errorf("during DB operation execution: %s, DB operation: %s", err.Error(), description)
		return err
	}
	if err != nil {
		err = fmt.Errorf("while checking SQL row affected: %s, DB operation: %s", err.Error(), description)
		return err
	}

	if *modifiedQueryID != queryID {
		err = fmt.Errorf(
			"ID of query instance to modify is not equal to the modified one: (input query ID: %d, output query ID: %d), DB operation: %s",
			queryID,
			*modifiedQueryID,
			description,
		)
		return err
	}

	logrus.Debugf("DB operation successful, ID of affected query instance: %d, operation: %s", *modifiedQueryID, description)

	return err
}

// CheckQueryID checks whether the user really has a query before inserting a new cohort defined by that query's id
func CheckQueryID(userID string, queryID int) (bool, error) {
	description := fmt.Sprintf("CheckQueryID (user ID: %s, query ID: %d), procedure: %s", userID, queryID, "query_tools.check_query_id")
	logrus.Debugf("running: %s", description)
	row := utilserver.DBConnection.QueryRow("SELECT query_tools.check_query_id($1, $2);", userID, queryID)
	res := new(bool)
	err := row.Scan(res)
	if err != nil {
		err = fmt.Errorf("while scanning SQL record: %s, DB operation: %s", err.Error(), description)
		return false, err
	}

	logrus.Debugf("DB operation successful, result: %t, DB operation: %s", *res, description)

	return *res, err

}

// GetQueryDefinition is called when the query is created. A new row is inserted in explore query results table with status 'running'.
func GetQueryDefinition(queryID int) (string, error) {

	description := fmt.Sprintf("GetQueryDefinition (ID %d), procedure: %s", queryID, "query_tools.get_query_definition")
	logrus.Debugf("running: %s", description)

	row := utilserver.DBConnection.QueryRow("SELECT query_tools.get_query_definition($1);", queryID)
	res := new(string)
	err := row.Scan(res)
	if err != nil {
		err = fmt.Errorf("while scanning SQL record: %s, DB operation: %s", err.Error(), description)
		return "", err
	}

	logrus.Debugf("DB operation successful, result: %s, DB operation: %s", *res, description)

	return *res, nil

}
