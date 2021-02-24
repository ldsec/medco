package timepoints

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

const SQLDateFormat = "2006-01-02"

func StartEvent(patientList []int64, conceptCodes, modifierCodes []string, earliest bool) (map[int64]time.Time, error) {

	setStrings := make([]string, len(patientList))

	for i, patient := range patientList {
		setStrings[i] = strconv.FormatInt(patient, 10)
	}
	setDefinition := "{" + strings.Join(setStrings, ",") + "}"
	conceptDefinition := "{" + strings.Join(conceptCodes, ",") + "}"
	modifierDefinition := "{" + strings.Join(modifierCodes, ",") + "}"

	description := fmt.Sprintf("get start event (patient list: %s, start concept codes: %s, start modifier codes: %s, begins with earliest occurence: %t): procedure: %s",
		setDefinition, conceptDefinition, modifierDefinition, earliest, "i2b2demodata_i2b2.start_event")

	logrus.Debugf("Retrieving the start event dates for the patients: %s", description)
	row, err := utilserver.I2B2DBConnection.Query("SELECT i2b2demodata_i2b2.start_event($1,$2,$3,$4)", setDefinition, conceptDefinition, modifierDefinition, earliest)
	if err != nil {
		err = errors.Errorf("while calling database for retrieving start event dates: %s; DB operation: %s", err.Error(), description)
		return nil, err
	}
	var record = new(string)
	for row.Next() {
		err = row.Scan(record)
		if err != nil {
			err = errors.Errorf("while reading database record stream for retrieving start event dates: %s; DB operation: %s", err.Error(), description)
			return nil, err
		}

		recordEntries := strings.Split(strings.Trim(*record, "()"), ",")
		patientID, err := strconv.ParseInt(recordEntries[0], 10, 64)
		if err != nil {
			err = errors.Errorf("while parsing patient number \"%s\": %s; DB operation: %s", recordEntries[0], err.Error(), description)
			return nil, err
		}
		startDate, err := time.Parse(SQLDateFormat, recordEntries[1])
		if err != nil {
			err = errors.Errorf("while parsing patient number \"%s\": %s; DB operation: %s", recordEntries[0], err.Error(), description)
			return nil, err
		}
	}

}

func EndEvents(patientWithStartEventList map[int64]time.Time, conceptCodes, modifierCodes []string) (map[int64][]time.Time, error)

func CensoringEvents(patientWithoutEndEvent map[int64]struct{}) (map[int64]time.Time, error)
