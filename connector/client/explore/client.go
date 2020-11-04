package exploreclient

import (
	"fmt"
	"os"
	"sort"
	"strconv"
	"time"

	medcoclient "github.com/ldsec/medco/connector/client"
	"github.com/ldsec/medco/connector/wrappers/unlynx"

	"github.com/ldsec/medco/connector/restapi/models"
	utilclient "github.com/ldsec/medco/connector/util/client"
	"github.com/sirupsen/logrus"
)

// ExecuteClientQuery executes and displays the result of the MedCo client query
func ExecuteClientQuery(token, username, password, queryType, queryString, resultOutputFilePath string, disableTLSCheck bool) (err error) {

	// get token
	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)

	// parse query type
	queryTypeParsed := models.ExploreQueryType(queryType)
	err = queryTypeParsed.Validate(nil)
	if err != nil {
		logrus.Error("invalid query type")
		return
	}

	// parse query string
	panelsItemKeys, panelsIsNot, err := medcoclient.ParseQueryString(queryString)
	if err != nil {
		return
	}

	// encrypt item key
	encPanelsItemKeys, err := unlynx.EncryptMatrix(panelsItemKeys)
	if err != nil {
		return
	}

	// execute query
	clientQuery, err := NewExploreQuery(accessToken, queryTypeParsed, encPanelsItemKeys, panelsIsNot, disableTLSCheck)
	if err != nil {
		return
	}

	nodesResult, err := clientQuery.Execute()
	if err != nil {
		return
	}

	// output results
	err = printResultsCSV(nodesResult, resultOutputFilePath)
	return
}

// printResultsCSV prints on a specified output in a CSV format the results, each node being one line
func printResultsCSV(nodesResult map[int]*ExploreQueryResult, CSVFileURL string) (err error) {

	CSVFile, err := utilclient.NewCSV(CSVFileURL)
	if err != nil {
		logrus.Warn("erro opening csv file result: ", err)
		return
	}

	csvHeaders := []string{"node_name", "count", "patient_list", "patient_set_id"}
	csvNodesResults := make([][]string, 0)

	// CSV values: results
	for nodeIdx, queryResult := range nodesResult {
		csvNodesResults = append(csvNodesResults, []string{
			strconv.Itoa(nodeIdx),
			strconv.FormatInt(queryResult.Count, 10),
			fmt.Sprint(queryResult.PatientList),
			strconv.FormatInt(queryResult.PatientSetID, 10),
		})
	}

	// CSV headers: timers
	for _, queryResult := range nodesResult {

		// sort the timers by name for deterministic output
		timerNames := make([]string, 0)
		for timerName := range queryResult.Times {
			timerNames = append(timerNames, timerName)
		}
		sort.Strings(timerNames)

		// add to headers
		for _, timerName := range timerNames {
			csvHeaders = append(csvHeaders, timerName)
		}
		break
	}

	// CSV values: timers: iter over results, then iter over timer names from csvHeaders
	for csvNodeResultsIdx, csvNodeResults := range csvNodesResults {
		nodeIdx, err := strconv.Atoi(csvNodeResults[0])
		if err != nil {
			logrus.Error("error parsing node number: ", err)
			return err
		}

		for timerNameIdx := 3; timerNameIdx < len(csvHeaders); timerNameIdx++ {
			timerName := csvHeaders[timerNameIdx]
			timerValue := nodesResult[nodeIdx].Times[timerName]

			csvNodesResults[csvNodeResultsIdx] = append(csvNodesResults[csvNodeResultsIdx],
				strconv.FormatInt(int64(timerValue/time.Millisecond), 10))
		}
	}

	// write to output

	err = CSVFile.Write(csvHeaders)
	if err != nil {
		logrus.Warn("error printing results: ", err)
	}
	err = CSVFile.WriteAll(csvNodesResults)
	if err != nil {
		logrus.Warn("error printing results: ", err)
	}
	err = CSVFile.Flush()
	if err != nil {
		logrus.Warn("error flushing csv file results: ", err)
	}

	return
}

// ExecuteClientSearchConcept executes and displays the result of the MedCo sconcept search
func ExecuteClientSearchConcept(token, username, password, conceptPath, resultOutputFilePath string, disableTLSCheck bool) (err error) {

	// get token
	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		return err
	}

	// execute search
	clientSearchConcept, err := NewExploreSearchConcept(accessToken, conceptPath, disableTLSCheck)
	if err != nil {
		return err
	}

	result, err := clientSearchConcept.Execute()
	if err != nil {
		return
	}

	output := "PATH" + "\t" + "TYPE" + "\n"
	for _, child := range result.Payload.Results {
		output += child.Path + "\t" + child.Type + "\n"
	}

	if resultOutputFilePath == "" {
		fmt.Println(output)
	} else {
		outputFile, err := os.Create(resultOutputFilePath)
		if err != nil {
			logrus.Error("error opening file: ", err)
		}
		outputFile.WriteString(output)
		outputFile.Close()
	}

	return
}

// ExecuteClientSearchModifier executes and displays the result of the MedCo modifier search
func ExecuteClientSearchModifier(token, username, password, modifierPath, appliedPath, appliedConcept, resultOutputFilePath string, disableTLSCheck bool) (err error) {

	// get token
	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		return err
	}

	// execute search
	clientSearchConcept, err := NewExploreSearchModifier(accessToken, modifierPath, appliedPath, appliedConcept, disableTLSCheck)
	if err != nil {
		return err
	}

	result, err := clientSearchConcept.Execute()
	if err != nil {
		return
	}

	output := "PATH" + "\t" + "TYPE" + "\n"
	for _, child := range result.Payload.Results {
		output += child.Path + "\t" + child.Type + "\n"
	}

	if resultOutputFilePath == "" {
		fmt.Println(output)
	} else {
		outputFile, err := os.Create(resultOutputFilePath)
		if err != nil {
			logrus.Error("error opening file: ", err)
		}
		outputFile.WriteString(output)
		outputFile.Close()
	}

	return
}

// ExecuteClientGenomicAnnotationsGetValues displays the genomic annotations values matching the "annotation" parameter
func ExecuteClientGenomicAnnotationsGetValues(token, username, password, annotation, value string, limit int64, disableTLSCheck bool) (err error) {

	// get token
	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		return err
	}

	// execute query
	clientGenomicAnnotationsGetValues, err := NewGenomicAnnotationsGetValues(accessToken, annotation, value, &limit, disableTLSCheck)
	if err != nil {
		return
	}

	result, err := clientGenomicAnnotationsGetValues.Execute()
	if err != nil {
		return
	}

	for _, annotation := range result {
		fmt.Printf("%s\n", annotation)
	}

	return

}

// ExecuteClientGenomicAnnotationsGetVariants displays the variant ids corresponding to the annotation and value parameters
func ExecuteClientGenomicAnnotationsGetVariants(token, username, password, annotation, value string, zygosity string, encrypted bool, disableTLSCheck bool) (err error) {

	// get token
	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		return err
	}

	// execute query
	clientGenomicAnnotationsGetVariants, err := NewGenomicAnnotationsGetVariants(accessToken, annotation, value, zygosity, &encrypted, disableTLSCheck)
	if err != nil {
		return
	}

	result, err := clientGenomicAnnotationsGetVariants.Execute()
	if err != nil {
		return
	}

	for _, variant := range result {
		fmt.Printf("%s\n", variant)
	}

	return

}
