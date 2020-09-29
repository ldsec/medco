package medcoclient

import (
	"bufio"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/ldsec/medco-connector/restapi/models"
	utilclient "github.com/ldsec/medco-connector/util/client"
	"github.com/ldsec/medco-connector/wrappers/unlynx"
	"github.com/ldsec/medco-loader/loader/identifiers"
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// ExecuteClientQuery executes and displays the result of the MedCo client query
func ExecuteClientQuery(token, username, password, queryType, queryString, resultOutputFilePath string, disableTLSCheck bool) (err error) {

	// get token
	accessToken, err := getAccessToken(token, username, password, disableTLSCheck)
	if err != nil {
		return err
	}

	// parse query type
	queryTypeParsed := models.ExploreQueryType(queryType)
	err = queryTypeParsed.Validate(nil)
	if err != nil {
		logrus.Error("invalid query type")
		return
	}

	// parse query string
	panelsItemKeys, panelsIsNot, err := parseQueryString(queryString)
	if err != nil {
		return
	}

	// encrypt the item keys
	encPanelsItemKeys := make([][]string, 0)
	for _, panel := range panelsItemKeys {
		encItemKeys := make([]string, 0)
		for _, itemKey := range panel {
			encrypted, err := unlynx.EncryptWithCothorityKey(itemKey)
			if err != nil {
				return err
			}
			encItemKeys = append(encItemKeys, encrypted)
		}
		encPanelsItemKeys = append(encPanelsItemKeys, encItemKeys)
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
	var output io.Writer
	if resultOutputFilePath == "" {
		output = os.Stdout
	} else {
		output, err = os.Create(resultOutputFilePath)
		if err != nil {
			logrus.Error("error opening file: ", err)
		}
		err = os.Chmod(resultOutputFilePath, 0777)
		if err != nil {
			logrus.Error("error setting permissions on file: ", err)
		}
	}
	err = printResultsCSV(nodesResult, output)
	return
}

// printResultsCSV prints on a specified output in a CSV format the results, each node being one line
func printResultsCSV(nodesResult map[int]*ExploreQueryResult, output io.Writer) (err error) {

	csvHeaders := []string{"node_name", "count", "patient_list"}
	csvNodesResults := make([][]string, 0)

	// CSV values: results
	for nodeIdx, queryResult := range nodesResult {
		csvNodesResults = append(csvNodesResults, []string{
			strconv.Itoa(nodeIdx),
			strconv.FormatInt(queryResult.Count, 10),
			fmt.Sprint(queryResult.PatientList),
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
	csvWriter := csv.NewWriter(output)
	err = csvWriter.Write(csvHeaders)
	if err != nil {
		logrus.Warn("error printing results: ", err)
	}
	err = csvWriter.WriteAll(csvNodesResults)
	if err != nil {
		logrus.Warn("error printing results: ", err)
	}
	return
}

// parseQueryString parses the query string given as input
func parseQueryString(queryString string) (panelsItemKeys [][]int64, panelsIsNot []bool, err error) {
	logrus.Info("Client query is: ", queryString)

	panelsItemKeys = make([][]int64, 0)
	panelsIsNot = make([]bool, 0)

	for _, queryPanel := range strings.Split(queryString, " AND ") {

		// parse panel negation
		if strings.HasPrefix(queryPanel, "NOT ") {
			panelsIsNot = append(panelsIsNot, true)
			queryPanel = queryPanel[4:]
		} else {
			panelsIsNot = append(panelsIsNot, false)
		}

		// parse query items
		itemKeys := make([]int64, 0)
		for _, queryItem := range strings.Split(queryPanel, " OR ") {
			// 3 cases: simple integer, integer to be repeated, query file

			// case 1: integer to be repeated
			if strings.Contains(queryItem, "^") {
				logrus.Debug("Client query integer repeated item: ", queryItem)

				elements := strings.Split(queryItem, "^")
				if len(elements) != 2 {
					err = errors.New("query item contains more than one ^")
					logrus.Error(err)
					return
				}

				queryInt, queryIntErr := strconv.ParseInt(elements[0], 10, 64)
				if queryIntErr != nil {
					logrus.Error("could not parse query integer: ", queryIntErr)
					return nil, nil, queryIntErr
				}

				intMultiplier, intMultiplierErr := strconv.ParseInt(elements[1], 10, 64)
				if intMultiplierErr != nil {
					logrus.Error("could not parse query integer multiplier: ", intMultiplierErr)
					return nil, nil, intMultiplierErr
				}

				for i := 0; i < int(intMultiplier); i++ {
					itemKeys = append(itemKeys, queryInt)
				}

			} else {
				parsedInt, parsedErr := strconv.ParseInt(queryItem, 10, 64)

				// case 2: simple integer
				if parsedErr == nil {
					logrus.Debug("Client query integer item: ", queryItem)

					// if a parsable integer: use as is
					itemKeys = append(itemKeys, parsedInt)

					// case 3: query file
				} else {
					logrus.Debug("Client query file item: ", queryItem)

					// else assume it is a file
					itemKeysFromFile, err := loadQueryFile(queryItem)
					if err != nil {
						return nil, nil, err
					}
					itemKeys = append(itemKeys, itemKeysFromFile...)
				}
			}
		}
		panelsItemKeys = append(panelsItemKeys, itemKeys)
	}
	return
}

// todo: might fail if alleles of queries are too big, what to do? ignore or fail?
// loadQueryFile load and parse a query file (either simple integer or genomic) into integers
func loadQueryFile(queryFilePath string) (queryTerms []int64, err error) {
	logrus.Debug("Client query: loading file ", queryFilePath)

	queryFile, err := os.Open(queryFilePath)
	if err != nil {
		logrus.Error("error opening query file: ", err)
		return
	}

	fileScanner := bufio.NewScanner(queryFile)
	for fileScanner.Scan() {
		queryTermFields := strings.Split(fileScanner.Text(), ",")
		var queryTerm int64

		if len(queryTermFields) == 1 {

			// simple integer identifier
			queryTerm, err = strconv.ParseInt(queryTermFields[0], 10, 64)
			if err != nil {
				logrus.Error("error parsing query term: ", err)
				return
			}

		} else if len(queryTermFields) == 4 {

			// genomic identifier, generate the variant ID
			startPos, err := strconv.ParseInt(queryTermFields[1], 10, 64)
			if err != nil {
				logrus.Error("error parsing start position: ", err)
				return nil, err
			}

			queryTerm, err = identifiers.GetVariantID(queryTermFields[0], startPos, queryTermFields[2], queryTermFields[3])
			if err != nil {
				logrus.Error("error generating genomic query term: ", err)
				return nil, err
			}

		} else {
			err = errors.New("dataset with " + string(len(queryTermFields)) + " fields is not supported")
			logrus.Error(err)
			return
		}

		queryTerms = append(queryTerms, queryTerm)
	}

	return
}

// ExecuteClientSearchConcept executes and displays the result of the MedCo sconcept search
func ExecuteClientSearchConcept(token, username, password, conceptPath, resultOutputFilePath string, disableTLSCheck bool) (err error) {

	// get token
	accessToken, err := getAccessToken(token, username, password, disableTLSCheck)
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
	accessToken, err := getAccessToken(token, username, password, disableTLSCheck)
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
	accessToken, err := getAccessToken(token, username, password, disableTLSCheck)
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
	accessToken, err := getAccessToken(token, username, password, disableTLSCheck)
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

func getAccessToken(token, username, password string, disableTLSCheck bool) (accessToken string, err error) {
	if len(token) > 0 {
		accessToken = token
	} else {
		logrus.Debug("No token provided, requesting token for user ", username, ", disable TLS check: ", disableTLSCheck)
		accessToken, err = utilclient.RetrieveAccessToken(username, password, disableTLSCheck)
	}
	return
}
