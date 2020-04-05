package i2b2

import (
	"encoding/json"
	utilserver "github.com/ldsec/medco-connector/util/server"
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
			"\\e2etest\\",
			"\\e2etest\\1\\",
			"\\e2etest\\2\\",
			"\\e2etest\\3\\",
		},
	},
}

func TestI2b2() (testPassed bool) {

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

func testDBConnection() (testPassed bool) {
	var err error
	utilserver.I2b2DBConnection, err = utilserver.InitializeConnectionToDB(utilserver.I2b2DBHost, utilserver.I2b2DBPort, utilserver.I2b2DBName, utilserver.I2b2DBLoginUser, utilserver.I2b2DBLoginPassword)
	if err != nil {
		return false
	}

	err = utilserver.I2b2DBConnection.Ping()
	if err != nil {
		logrus.Error("Impossible to connect to i2b2 DB: " + err.Error())
		return false
	}

	return true
}

func testI2b2Query(testParams testI2b2QueryParameters) (testPassed bool) {

	if !testDBConnection() {
		return false
	}

	var testNames []string
	var testName string
	var err error

	rows, err := utilserver.I2b2DBConnection.Query("SELECT c_fullname FROM medco_ont.e2etest ORDER BY c_fullname")

	if err != nil {
		logrus.Error("i2b2 query execution error: " + err.Error())
		return false
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&testName)
		if err != nil {
			logrus.Error("i2b2 query result reading error: " + err.Error())
			return false
		}
		testNames = append(testNames, testName)
	}

	if !areEqual(testNames, testParams.queryResult) {
		logrus.Error("Wrong query result: " + strings.Join(testNames, ","))
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
