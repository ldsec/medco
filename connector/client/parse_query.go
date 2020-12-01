package medcoclient

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/ldsec/medco/connector/restapi/models"

	"github.com/ldsec/medco/loader/identifiers"
	"github.com/sirupsen/logrus"
)

// ParseQueryString parses the query string given as input
func ParseQueryString(queryString string) (panels []*models.Panel, err error) {
	logrus.Info("Client query is: ", queryString)

	panels = make([]*models.Panel, 0)

	for _, queryPanel := range strings.Split(queryString, " AND ") {

		var newPanel models.Panel
		var not bool

		// parse panel negation
		if strings.HasPrefix(queryPanel, "NOT ") {
			not = true
			queryPanel = queryPanel[4:]
		} else {
			not = false
		}
		newPanel.Not = &not

		// parse query items
		for _, queryItem := range strings.Split(queryPanel, " OR ") {

			// queryItem can contain two fields separeted by "^"
			// the first field is in the format parsed by parseQueryItem,
			// the second field contains the number of times the first field must be OR-ed with itself
			// e.g. enc::5^3 --> enc::5 OR enc::5 OR enc::5
			if strings.Contains(queryItem, "^") {
				logrus.Debug("Client query integer repeated item: ", queryItem)

				elements := strings.Split(queryItem, "^")
				if len(elements) != 2 {
					err = errors.New("query item contains more than one ^")
					logrus.Error(err)
					return
				}

				items, err := ParseQueryItem(elements[0])
				if err != nil {
					return nil, err
				}

				multiplier, err := strconv.ParseInt(elements[1], 10, 64)
				if err != nil {
					logrus.Error("invalid multiplier ", elements[1])
					return nil, err
				}

				for i := 0; i < int(multiplier); i++ {
					newPanel.Items = append(newPanel.Items, items...)
				}

			} else {

				items, err := ParseQueryItem(queryItem)
				if err != nil {
					return nil, err
				}

				newPanel.Items = append(newPanel.Items, items...)
			}
		}
		panels = append(panels, &newPanel)
	}
	return
}

// ParseQueryItem parses a string into an item
// queryItem is composed of one field, the content field, separated by "::"
// possible values of the type field are: "enc", "clr", "file"
// when the type field is equal to "enc", the content field contains the concept ID
// when the type field is equal to "clr", the content field contains the concept path
// and, possibly, the modifier field, which in turn contains the modifier key and applied path fields, separated by ":"
// when the type field is equal to "file", the content field contains the path of the file containing the items
func ParseQueryItem(queryItem string) (items []*models.PanelItemsItems0, err error) {

	queryItemFields := strings.Split(queryItem, "::")
	if len(queryItemFields) < 2 {
		return nil, fmt.Errorf("invalid query item format: %s", queryItem)
	}

	switch queryItemFields[0] {
	case "enc":
		if len(queryItemFields) != 2 {
			return nil, fmt.Errorf("invalid enc query item format: %v", queryItemFields)
		}

		_, err = strconv.ParseInt(queryItemFields[1], 10, 64)
		if err != nil {
			logrus.Error("invalid ID ", queryItemFields[1])
			return nil, err
		}

		item := new(models.PanelItemsItems0)
		encrypted := true

		item.Encrypted = &encrypted
		item.QueryTerm = &queryItemFields[1]
		items = append(items, item)
	case "clr":
		if len(queryItemFields) > 3 {
			return nil, fmt.Errorf("invalid clr query item format: %v", queryItemFields)
		}
		item := new(models.PanelItemsItems0)
		encrypted := false

		item.Encrypted = &encrypted
		item.QueryTerm = &queryItemFields[1]

		if len(queryItemFields) == 3 {
			modifierFields := strings.Split(queryItemFields[2], ":")
			if len(modifierFields) != 2 {
				return nil, fmt.Errorf("invalid modifier term format: %v", modifierFields)
			}

			modifier := &models.PanelItemsItems0Modifier{
				AppliedPath: modifierFields[1],
				ModifierKey: modifierFields[0],
			}

			item.Modifier = modifier
		}
		items = append(items, item)
	case "file":
		if len(queryItemFields) != 2 {
			return nil, fmt.Errorf("invalid file query item format: %v", queryItemFields)
		}
		items, err = loadQueryFile(queryItemFields[1])
	default:
		return nil, fmt.Errorf("invalid query item type: %s", queryItemFields[0])
	}

	return
}

// TODO: might fail if alleles of queries are too big, what to do? ignore or fail?
// loadQueryFile load and parse a query file (containing either regular or genomic IDs) into query items
// A regular ID is an ID in the format parsable by parseQueryItem, a genomic ID is a sequence of 4 comma separated values
// Each query item (i.e. regular or genomic ID) must occupy a line in the file.
// As expected, all query items in a file are part of the same panel, and therefore OR-ed
func loadQueryFile(queryFilePath string) (queryItems []*models.PanelItemsItems0, err error) {
	logrus.Debug("Client query: loading file ", queryFilePath)

	queryFile, err := os.Open(queryFilePath)
	if err != nil {
		logrus.Error("error opening query file: ", err)
		return
	}

	fileScanner := bufio.NewScanner(queryFile)
	for fileScanner.Scan() {
		queryTermFields := strings.Split(fileScanner.Text(), ",")
		var queryItem []*models.PanelItemsItems0

		if len(queryTermFields) == 1 { // regular ID

			queryItem, err = ParseQueryItem(queryTermFields[0])
			if err != nil {
				return nil, err
			}

		} else if len(queryTermFields) == 4 { // genomic ID, generate the variant ID

			startPos, err := strconv.ParseInt(queryTermFields[1], 10, 64)
			if err != nil {
				logrus.Error("error parsing start position: ", err)
				return nil, err
			}

			variantID, err := identifiers.GetVariantID(queryTermFields[0], startPos, queryTermFields[2], queryTermFields[3])
			if err != nil {
				logrus.Error("error generating genomic query term: ", err)
				return nil, err
			}

			queryItem, err = ParseQueryItem("enc::" + strconv.FormatInt(variantID, 64))
			if err != nil {
				return nil, err
			}

		} else {
			err = errors.New("dataset with " + fmt.Sprint(len(queryTermFields)) + " fields is not supported")
			logrus.Error(err)
			return
		}

		queryItems = append(queryItems, queryItem...)
	}

	return
}
