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
// A query string is a list of panels concatenated by " AND ".
// Each panel is a list of query items, in the format parsed by parseQueryItem, concatenated by " OR ".
// Each query item can be OR-ed n times with itself (and so lengthening the panel's query items list) using the syntax query_item^n.
// The first element of a panel can be a "NOT". In this case the panel is negated.
// The last element of a panel can be the panel timing, whose value either "any", "samevisit", or "sameinstancenum".
// If omitted, the panel timing is defaulted to "any".
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

		// parse panel timing
		queryPanelFields := strings.Split(queryPanel, " ")
		panelTiming := strings.ToLower(queryPanelFields[len(queryPanelFields)-1])
		switch panelTiming {
		case string(models.TimingSamevisit):
			newPanel.PanelTiming = models.TimingSamevisit
			queryPanel = queryPanel[:len(queryPanel)-len(panelTiming)-1]
		case string(models.TimingSameinstancenum):
			newPanel.PanelTiming = models.TimingSameinstancenum
			queryPanel = queryPanel[:len(queryPanel)-len(panelTiming)-1]
		case string(models.TimingAny):
			newPanel.PanelTiming = models.TimingAny
			queryPanel = queryPanel[:len(queryPanel)-len(panelTiming)-1]
		default:
			newPanel.PanelTiming = models.TimingAny
		}

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

				conceptItems, cohortItems, err := ParseQueryItem(elements[0])
				if err != nil {
					return nil, err
				}

				multiplier, err := strconv.ParseInt(elements[1], 10, 64)
				if err != nil {
					logrus.Error("invalid multiplier ", elements[1])
					return nil, err
				}

				for i := 0; i < int(multiplier); i++ {
					newPanel.ConceptItems = append(newPanel.ConceptItems, conceptItems...)
					newPanel.CohortItems = append(newPanel.CohortItems, cohortItems...)
				}

			} else {

				conceptItems, cohortItems, err := ParseQueryItem(queryItem)
				if err != nil {
					return nil, err
				}

				newPanel.ConceptItems = append(newPanel.ConceptItems, conceptItems...)
				newPanel.CohortItems = append(newPanel.CohortItems, cohortItems...)
			}
		}
		panels = append(panels, &newPanel)
	}
	return
}

// ParseQueryItem parses a string into an item.
// queryItem is composed of two mandatory fields, the type field and the content field,
// and an optional field, the constraint field, separated by "::".
//		type::content[::constraint]
// Possible values of the type field are: "chr", "enc", "clr", "file".
// 1. When the type field is equal to "chr", the content field contains the cohort name. The constraint field is not present in this case.
// 2. When the type field is equal to "enc", the content field contains the concept ID. The constraint field is not present in this case.
// 3. When the type field is equal to "clr", the content field contains the concept field (containing the concept path)
// 		and, possibly, the modifier field, which in turn contains the modifier key and applied path fields, separated by ":".
// 		The optional constraint field can be present, containing the operator, type and value fields separated by ":".
//		The constraint field applies to the concept or, if the modifier field is present, to the modifier.
// 4. When the type field is equal to "file", the content field contains the path of the file containing the items. The constraint field is not present in this case.
func ParseQueryItem(queryItem string) (conceptItems []*models.PanelConceptItemsItems0, cohortItems []string, err error) {

	queryItemFields := strings.Split(queryItem, "::")
	if len(queryItemFields) != 2 && len(queryItemFields) != 3 {
		return nil, nil, fmt.Errorf("invalid query item format: %s", queryItem)
	}

	switch queryItemFields[0] {
	case "chr":
		if len(queryItemFields) != 2 {
			return nil, nil, fmt.Errorf("invalid chr query item format: %v", queryItemFields)
		}

		cohortItems = append(cohortItems, queryItemFields[1])
	case "enc":
		if len(queryItemFields) != 2 {
			return nil, nil, fmt.Errorf("invalid enc query item format: %v", queryItemFields)
		}

		_, err = strconv.ParseInt(queryItemFields[1], 10, 64)
		if err != nil {
			logrus.Error("invalid ID ", queryItemFields[1])
			return nil, nil, err
		}

		item := new(models.PanelConceptItemsItems0)
		encrypted := true
		item.Encrypted = &encrypted
		item.QueryTerm = &queryItemFields[1]

		conceptItems = append(conceptItems, item)
	case "clr":
		contentFieldFields := strings.Split(queryItemFields[1], ":")
		if len(contentFieldFields) != 1 && len(contentFieldFields) != 3 {
			return nil, nil, fmt.Errorf("invalid content field format: %v", contentFieldFields)
		}

		item := new(models.PanelConceptItemsItems0)
		encrypted := false

		item.Encrypted = &encrypted
		item.QueryTerm = &contentFieldFields[0]

		if len(contentFieldFields) == 3 { // there is a modifier field
			modifier := &models.PanelConceptItemsItems0Modifier{
				AppliedPath: &contentFieldFields[2],
				ModifierKey: &contentFieldFields[1],
			}

			item.Modifier = modifier
		}

		if len(queryItemFields) == 3 { // there is a constrain field
			constrainFieldFields := strings.Split(queryItemFields[2], ":")
			if len(constrainFieldFields) != 3 {
				return nil, nil, fmt.Errorf("invalid constrain field format: %v", constrainFieldFields)
			}

			item.Operator = constrainFieldFields[0]
			item.Type = constrainFieldFields[1]
			item.Value = constrainFieldFields[2]
		}

		conceptItems = append(conceptItems, item)
	case "file":
		if len(queryItemFields) != 2 {
			return nil, nil, fmt.Errorf("invalid file query item format: %v", queryItemFields)
		}
		conceptItems, cohortItems, err = loadQueryFile(queryItemFields[1])
	default:
		return nil, nil, fmt.Errorf("invalid query item type: %s", queryItemFields[0])
	}

	return
}

// TODO: might fail if alleles of queries are too big, what to do? ignore or fail?
// loadQueryFile loads and parses a query file (containing either regular or genomic IDs) into query items.
// A regular ID is an ID in the format parsable by parseQueryItem, a genomic ID is a sequence of 4 comma separated values.
// Each query item (i.e. regular or genomic ID) must occupy a line in the file.
// As expected, all query items in a file are part of the same panel, and therefore OR-ed.
func loadQueryFile(queryFilePath string) (conceptItems []*models.PanelConceptItemsItems0, cohortItems []string, err error) {
	logrus.Debug("Client query: loading file ", queryFilePath)

	queryFile, err := os.Open(queryFilePath)
	if err != nil {
		logrus.Error("error opening query file: ", err)
		return
	}

	fileScanner := bufio.NewScanner(queryFile)
	for fileScanner.Scan() {
		queryTermFields := strings.Split(fileScanner.Text(), ",")
		var conceptItem []*models.PanelConceptItemsItems0
		var cohortItem []string

		if len(queryTermFields) == 1 { // regular ID

			conceptItem, cohortItem, err = ParseQueryItem(queryTermFields[0])
			if err != nil {
				return nil, nil, err
			}

		} else if len(queryTermFields) == 4 { // genomic ID, generate the variant ID

			startPos, err := strconv.ParseInt(queryTermFields[1], 10, 64)
			if err != nil {
				logrus.Error("error parsing start position: ", err)
				return nil, nil, err
			}

			variantID, err := identifiers.GetVariantID(queryTermFields[0], startPos, queryTermFields[2], queryTermFields[3])
			if err != nil {
				logrus.Error("error generating genomic query term: ", err)
				return nil, nil, err
			}

			conceptItem, cohortItem, err = ParseQueryItem("enc::" + strconv.FormatInt(variantID, 64))
			if err != nil {
				return nil, nil, err
			}

		} else {
			err = errors.New("dataset with " + fmt.Sprint(len(queryTermFields)) + " fields is not supported")
			logrus.Error(err)
			return
		}

		conceptItems = append(conceptItems, conceptItem...)
		cohortItems = append(cohortItems, cohortItem...)
	}

	return
}

// ParseSequences parses a string to a list of temporal sequence information.
// Multiple sequence information groups must be separated with columns ":".
// The different attributes inside a group must be separated with commas ",".
func ParseSequences(sequenceString string) (sequences []*models.TimingSequenceInfo, err error) {

	if sequenceString == "" {
		return
	}

	var seq *models.TimingSequenceInfo

	for _, sequenceString := range strings.Split(sequenceString, ":") {
		seq, err = parseSequence(sequenceString)
		if err != nil {
			err = fmt.Errorf("while parsing temporal sequence information: %s", err.Error())
			return
		}
		sequences = append(sequences, seq)
	}
	return

}

func parseSequence(sequenceString string) (sequence *models.TimingSequenceInfo, err error) {
	sequenceInfoStrings := strings.Split(sequenceString, ",")

	// the 5 items are:
	// 1. the operator (before, before or same time, same time)
	// 2. which occurence should be considered for the left operand (first, any, last)
	// 3. what date should be considered fot the left operand (startdate, enddate)
	// 4. which occurence should be considered for the right operand (first, any, last)
	// 5. what date should be considered fot the right operand (startdate, enddate)
	if len(sequenceInfoStrings) != 5 {
		err = fmt.Errorf("sequence info is expected to be composed of 5 elements separated by commas: sequence info string \"%s\"", sequenceString)
		return
	}

	sequence = &models.TimingSequenceInfo{}

	sequence.When = new(string)

	switch sequenceInfoStrings[0] {
	case "before":
		*sequence.When = models.TimingSequenceInfoWhenLESS
	case "beforeorsametime":
		*sequence.When = models.TimingSequenceInfoWhenLESSEQUAL
	case "sametime":
		*sequence.When = models.TimingSequenceInfoWhenEQUAL
	default:
		err = fmt.Errorf(`the first element of the info string is expected to be "before", "beforeorsametime" or "sametime": got "%s"`, sequenceInfoStrings[0])
		return
	}

	sequence.WhichObservationFirst = new(string)

	switch sequenceInfoStrings[1] {
	case "first":
		*sequence.WhichObservationFirst = models.TimingSequenceInfoWhichObservationFirstFIRST
	case "any":
		*sequence.WhichObservationFirst = models.TimingSequenceInfoWhichObservationFirstANY
	case "last":
		*sequence.WhichObservationFirst = models.TimingSequenceInfoWhichObservationFirstLAST
	default:
		err = fmt.Errorf(`the second element of the info string is expected to be "first", "any" or "last": got "%s"`, sequenceInfoStrings[1])
		return
	}

	sequence.WhichDateFirst = new(string)

	switch sequenceInfoStrings[2] {
	case "startdate":
		*sequence.WhichDateFirst = models.TimingSequenceInfoWhichDateFirstSTARTDATE
	case "enddate":
		*sequence.WhichDateFirst = models.TimingSequenceInfoWhichDateFirstENDDATE
	default:
		err = fmt.Errorf(`the third element of the info string is expected to be "startdate" or "enddate": got "%s"`, sequenceInfoStrings[2])
		return
	}

	sequence.WhichObservationSecond = new(string)

	switch sequenceInfoStrings[3] {
	case "first":
		*sequence.WhichObservationSecond = models.TimingSequenceInfoWhichObservationSecondFIRST
	case "any":
		*sequence.WhichObservationSecond = models.TimingSequenceInfoWhichObservationSecondANY
	case "last":
		*sequence.WhichObservationSecond = models.TimingSequenceInfoWhichObservationSecondLAST
	default:
		err = fmt.Errorf(`the fourth element of the info string is expected to be "first", "any" or "last": got "%s"`, sequenceInfoStrings[3])
		return
	}

	sequence.WhichDateSecond = new(string)

	switch sequenceInfoStrings[4] {
	case "startdate":
		*sequence.WhichDateSecond = models.TimingSequenceInfoWhichDateSecondSTARTDATE
	case "enddate":
		*sequence.WhichDateSecond = models.TimingSequenceInfoWhichDateSecondENDDATE
	default:
		err = fmt.Errorf(`the fifth element of the info string is expected to be "startdate" or "enddate": got "%s"`, sequenceInfoStrings[4])
		return
	}
	return
}
