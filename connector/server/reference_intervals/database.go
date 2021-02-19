package referenceintervalserver

import (
	"fmt"
	"strconv"

	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/sirupsen/logrus"
)

//QueryResult contains the information about a row that comes out from the query executed in RetrievePatients.
type QueryResult struct {
	NumericValue  float64
	Unit          string
	PatientNumber int64
}

// RetrieveObservationsForConcept returns the numerical values that correspond to the concept passed as argument for the specified cohort.
func RetrieveObservationsForConcept(code string, patients []int64, minObservation float64) (queryResults []QueryResult, err error) {
	logrus.Info("About to execute explore stats SQL query ", sqlModifier, ", concept:", code, ", patients: ", patients)
	return retrieveObservations(sqlConcept, code, patients, minObservation)
}

//RetrieveObservationsForModifier returns the numerical values that correspond to the modifier passed as argument for the specified cohort.
func RetrieveObservationsForModifier(code string, patients []int64, minObservation float64) (queryResults []QueryResult, err error) {
	logrus.Info("About to execute explore stats SQL query ", sqlModifier, ", modifier:", code, ", patients: ", patients)
	return retrieveObservations(sqlModifier, code, patients, minObservation)
}

//retrieveObservations returns the numerical values that correspond to the concept or modifier whose code is passed as argument for the specified cohort.
func retrieveObservations(sqlQuery, code string, patients []int64, minObservation float64) (queryResults []QueryResult, err error) {
	strPatientList := utilserver.ConvertIntListToString(patients)

	rows, err := utilserver.I2B2DBConnection.Query(sqlQuery, code, minObservation, strPatientList)
	if err != nil {
		err = fmt.Errorf("while execution SQL query: %s", err.Error())
		return
	}

	queryResults = make([]QueryResult, 0)

	for rows.Next() {
		numericValue := new(string)
		patientNb := new(string)
		unit := new(string)
		scanErr := rows.Scan(numericValue, patientNb, unit)
		if scanErr != nil {
			err = scanErr
			err = fmt.Errorf("while scanning SQL record: %s", err.Error())
			return
		}

		var queryResult QueryResult

		queryResult.Unit = *unit

		queryResult.NumericValue, err = strconv.ParseFloat(*numericValue, 64)
		if err != nil {
			err = fmt.Errorf("error while converting numerical value %s for the (concept or modifier) with code (%s)", *numericValue, code)
			return
		}

		queryResult.PatientNumber, err = strconv.ParseInt(*patientNb, 10, 64)
		if err != nil {
			err = fmt.Errorf("error while parsing the patient identifier %s for the (concept or modifier) with code (%s)", *patientNb, code)
			return
		}

		queryResults = append(queryResults, queryResult)
	}

	return

}

/*
	* This query will return the numerical values from all observations where
	* the patient_num is contained within the list passed as argument (the list is in principle a list of patient from a specific cohort).

	TODO In the same way I gathered the schema and table in which the ontology is contained, gather the schema in which observations are contained.
	For the moment I hardcode the table and schema.

	We only keep rows where nval_num is exactly equal to a specific values hence the required value of TVAL_CHAR.
	We could keep values which are GE or LE or L or G the problem is that we would need open brackets for intervals.
	VALTYPE_CD = 'N' because we only care about numerical values.
*/
const sqlStart string = `
SELECT nval_num, patient_num, units_cd FROM i2b2demodata_i2b2.observation_fact
	WHERE `

const sqlModifier string = sqlStart + ` modifier_cd = $1 ` + sqlEnd
const sqlConcept string = sqlStart + ` concept_cd = $1 ` + sqlEnd

const sqlEnd = ` AND valtype_cd = 'N' AND tval_char = 'E' AND nval_num is not null AND units_cd is not null AND units_cd != '@'
AND nval_num >= $2
AND patient_num = ANY($3::integer[]) `
