package utilclient

import (
	"bufio"
	"errors"
	"os"
	"strconv"
	"strings"

	"github.com/ldsec/medco-loader/loader/identifiers"
	"github.com/sirupsen/logrus"
)

// ParseQueryString parses the query string given as input
func ParseQueryString(queryString string) (panelsItemKeys [][]int64, panelsIsNot []bool, err error) {
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
				queryItem = strings.TrimSpace(queryItem)
				parsedInt, parsedErr := strconv.ParseInt(queryItem, 10, 64)
				if parsedErr != nil {
					logrus.Debug("Caught exception from strconv %s: ", parsedErr)
				}

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
