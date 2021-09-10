package exploreclient

import (
	"encoding/xml"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	medcoclient "github.com/ldsec/medco/connector/client"
	"github.com/ldsec/medco/connector/restapi/models"
	utilclient "github.com/ldsec/medco/connector/util/client"
	"github.com/ldsec/medco/connector/wrappers/unlynx"
	"github.com/sirupsen/logrus"
)

// ExecuteClientQuery executes and displays the result of the MedCo client query.
// endpoint on the server: /node/explore/query
func ExecuteClientQuery(token, username, password, queryString, queryTiming, querySequences, resultOutputFilePath string, disableTLSCheck bool) (err error) {

	// get token
	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		return err
	}

	// parse query string
	panels, err := medcoclient.ParseQueryString(queryString)
	if err != nil {
		return
	}

	// parse query sequences
	sequences, err := medcoclient.ParseSequences(querySequences)
	if err != nil {
		return
	}

	// encrypt item keys
	for _, panel := range panels {
		for _, item := range panel.ConceptItems {
			if *item.Encrypted {
				queryTermInt, err := strconv.ParseInt(*item.QueryTerm, 10, 64)
				if err != nil {
					return err
				}
				encrypted, err := unlynx.EncryptWithCothorityKey(queryTermInt)
				if err != nil {
					return err
				}
				item.QueryTerm = &encrypted
			}
		}
	}

	// execute query
	clientQuery, err := NewExploreQuery(accessToken, panels, models.Timing(strings.ToLower(queryTiming)), sequences, disableTLSCheck)
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

// printResultsCSV prints on a specified output in a CSV format the results, each node being one line.
func printResultsCSV(nodesResult map[int]*ExploreQueryResult, CSVFileURL string) (err error) {

	CSVFile, err := utilclient.NewCSV(CSVFileURL)
	if err != nil {
		logrus.Warn("erro opening csv file result: ", err)
		return
	}

	csvHeaders := []string{"node_name", "count", "patient_list", "query_id", "patient_set_id"}
	csvNodesResults := make([][]string, 0)

	// CSV values: results
	for nodeIdx, queryResult := range nodesResult {
		csvNodesResults = append(csvNodesResults, []string{
			strconv.Itoa(nodeIdx),
			strconv.FormatInt(queryResult.Count, 10),
			fmt.Sprint(queryResult.PatientList),
			strconv.FormatInt(queryResult.QueryID, 10),
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

// ExecuteClientSearch executes and displays the result of the MedCo search.
// endpoint on the server: /node/explore/search
func ExecuteClientSearch(token, username, password, searchString, resultOutputFilePath string, disableTLSCheck bool) (err error) {

	// get token
	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		return err
	}

	// execute search
	clientSearch, err := NewExploreSearch(accessToken, searchString, disableTLSCheck)
	if err != nil {
		return err
	}

	result, err := clientSearch.Execute()
	if err != nil {
		return
	}

	output := ""

	for _, element := range result.Payload.Results {
		tmp, err := xml.MarshalIndent(element, "  ", "    ")
		if err != nil {
			return err
		}
		output += fmt.Sprintln(string(tmp))
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

// ExecuteClientSearchConceptChildren executes and displays the result of the MedCo concept children search.
// endpoint on the server: /node/explore/search/concept
func ExecuteClientSearchConceptChildren(token, username, password, conceptPath, resultOutputFilePath string, disableTLSCheck bool) (err error) {

	// get token
	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		return err
	}

	// execute search
	clientSearchConceptChildren, err := NewExploreSearchConcept(accessToken, conceptPath, models.ExploreSearchConceptOperationChildren, disableTLSCheck)
	if err != nil {
		return err
	}

	result, err := clientSearchConceptChildren.Execute()
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

// ExecuteClientSearchModifierChildren executes and displays the result of the MedCo modifier children search.
// endpoint on the server: /node/explore/search/modifier
func ExecuteClientSearchModifierChildren(token, username, password, modifierPath, appliedPath, appliedConcept, resultOutputFilePath string, disableTLSCheck bool) (err error) {

	// get token
	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		return err
	}

	// execute search
	clientSearchModifierChildren, err := NewExploreSearchModifier(accessToken, modifierPath, appliedPath, appliedConcept, models.ExploreSearchModifierOperationChildren, disableTLSCheck)
	if err != nil {
		return err
	}

	result, err := clientSearchModifierChildren.Execute()
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

// ExecuteClientSearchConceptInfo executes and displays the result of the MedCo concept info search.
// endpoint on the server: /node/explore/search/concept
func ExecuteClientSearchConceptInfo(token, username, password, conceptPath, resultOutputFilePath string, disableTLSCheck bool) (err error) {

	// get token
	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		return err
	}

	// execute search
	clientSearchConceptInfo, err := NewExploreSearchConcept(accessToken, conceptPath, models.ExploreSearchConceptOperationInfo, disableTLSCheck)
	if err != nil {
		return err
	}

	result, err := clientSearchConceptInfo.Execute()
	if err != nil {
		return
	}

	output := ""
	for _, concept := range result.Payload.Results {
		tmp, err := xml.MarshalIndent(concept, "  ", "    ")
		if err != nil {
			return err
		}
		output += fmt.Sprintln(string(tmp))
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

// ExecuteClientSearchModifierInfo executes and displays the result of the MedCo modifier info search.
// endpoint on the server: /node/explore/search/modifier
func ExecuteClientSearchModifierInfo(token, username, password, modifierPath, appliedPath, resultOutputFilePath string, disableTLSCheck bool) (err error) {

	// get token
	accessToken, err := utilclient.RetrieveOrGetNewAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		return err
	}

	// execute search
	clientSearchModifierInfo, err := NewExploreSearchModifier(accessToken, modifierPath, appliedPath, "/", models.ExploreSearchModifierOperationInfo, disableTLSCheck)
	if err != nil {
		return err
	}

	result, err := clientSearchModifierInfo.Execute()
	if err != nil {
		return
	}

	output := ""
	for _, modifier := range result.Payload.Results {
		tmp, err := xml.MarshalIndent(modifier, "  ", "    ")
		if err != nil {
			return err
		}
		output += fmt.Sprintln(string(tmp))
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

// ExecuteClientGenomicAnnotationsGetValues displays the genomic annotations values matching the "annotation" parameter.
// endpoint on the server: /genomic-annotations/{annotation}
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

// ExecuteClientGenomicAnnotationsGetVariants displays the variant ids corresponding to the annotation and value parameters.
// endpoint on the server: /genomic-annotations/{annotation}/{value}
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
