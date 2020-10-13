package node

import (
	"encoding/json"
	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/sirupsen/logrus"
	"strings"
)

type testGenomicAnnotationsQueryParameters struct {
	query       string
	queryResult []string
}

var genomicAnnotationsQueryParameters = []testGenomicAnnotationsQueryParameters{
	{
		query: "SELECT variant_id FROM genomic_annotations.e2etest ORDER BY variant_id",
		queryResult: []string{
			"vID1",
			"vID2",
			"vID3",
			"vID4",
			"vID5",
		},
	},
}

func statusGenomicAnnotations() (testPassed bool) {

	for _, testParams := range genomicAnnotationsQueryParameters {
		if !testGenomicAnnotationsQuery(testParams) {
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

func testGenomicAnnotationsQuery(testParams testGenomicAnnotationsQueryParameters) (testPassed bool) {

	err := utilserver.MedcoDBConnection.Ping()
	if err != nil {
		logrus.Error("Impossible to connect to MedCo DB: " + err.Error())
		return false
	}

	var queryResult []string
	var row string

	rows, err := utilserver.MedcoDBConnection.Query(testParams.query)

	if err != nil {
		logrus.Error("genomic annotation query execution error: " + err.Error())
		return false
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&row)
		if err != nil {
			logrus.Error("genomic annotation query result reading error: " + err.Error())
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
