package node

import (
	"encoding/json"
	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/sirupsen/logrus"
	"strings"
)

type testI2b2QueryParameters struct {
	query       string
	queryResult []string
}

var i2b2QueryParameters = []testI2b2QueryParameters{
	{
		query: "SELECT c_fullname FROM medco_ont.e2etest ORDER BY c_fullname",
		queryResult: []string{
			"\\DeathStatus-status\\",
			"\\DeathStatus-status\\death",
			"\\DeathStatus-status\\unknown",
			"\\e2etest\\",
			"\\e2etest\\1\\",
			"\\e2etest\\2\\",
			"\\e2etest\\3\\",
			"\\FophDiagnosis-code\\ICD10\\",
			"\\FophDiagnosis-code\\ICD10\\Conditions on the perinatal period(760-779)\\",
			"\\FophDiagnosis-code\\ICD10\\Conditions on the perinatal period(760-779)\\Maternally caused (760-763)\\",
			"\\FophDiagnosis-code\\ICD10\\Conditions on the perinatal period(760-779)\\Maternally caused (760-763)\\(762) Fetus or newborn affected b~\\(762-3) Placental transfusion syn~\\",
			"\\I2B2\\",
			"\\I2B2\\Demographics\\",
			"\\I2B2\\Demographics\\Gender\\",
			"\\I2B2\\Demographics\\Gender\\Female\\",
			"\\I2B2\\Demographics\\Gender\\Male\\",
			"\\modifiers\\",
			"\\modifiers\\1\\",
			"\\modifiers\\2\\",
			"\\modifiers\\3\\",
			"\\SPHNv2020.1\\",
			"\\SPHNv2020.1\\DeathStatus\\",
			"\\SPHNv2020.1\\FophDiagnosis\\",
		},
	},
}

func statusI2b2() (testPassed bool) {

	for _, testParams := range i2b2QueryParameters {
		if !testI2b2Query(testParams) {
			log := "test failed: "
			text, err := json.Marshal(testParams)
			if err == nil {
				log += string(text)
			} else {
				log += err.Error()
			}
			logrus.Warn(log)
			return false
		}
	}

	return true

}

func testI2b2Query(testParams testI2b2QueryParameters) (testPassed bool) {

	err := utilserver.I2b2DBConnection.Ping()
	if err != nil {
		logrus.Error("Impossible to connect to i2b2 DB: " + err.Error())
		return false
	}

	var queryResult []string
	var row string

	rows, err := utilserver.I2b2DBConnection.Query(testParams.query)

	if err != nil {
		logrus.Error("i2b2 query execution error: " + err.Error())
		return false
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&row)
		if err != nil {
			logrus.Error("i2b2 query result reading error: " + err.Error())
			return false
		}
		queryResult = append(queryResult, row)
	}

	if !areEqual(queryResult, testParams.queryResult) {
		logrus.Error("Wrong query result: " + strings.Join(queryResult, ","))
		return false
	}

	return true

}

func areEqual(slice1, slice2 []string) bool {

	if len(slice1) != len(slice2) {
		return false
	}

	for i, element := range slice1 {
		if element != slice2[i] {
			return false
		}
	}

	return true
}
