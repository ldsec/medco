package medcoserver

import (
	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/ldsec/medco-connector/wrappers/i2b2"
	"github.com/ldsec/medco-connector/wrappers/unlynx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

// todo: log query (with associated status)
// todo: put user + query type + unique ID in query id

// ExploreQuery represents an i2b2-MedCo query to be executed
type ExploreQuery struct {
	ID     string
	Query  *models.ExploreQuery
	Result struct {
		EncCount       string
		EncPatientList []string
		Timers         map[string]time.Duration
	}
}

// NewExploreQuery creates a new query object and checks its validity
func NewExploreQuery(queryName string, query *models.ExploreQuery) (new ExploreQuery, err error) {
	new = ExploreQuery{
		ID:    queryName,
		Query: query,
	}
	new.Result.Timers = make(map[string]time.Duration)
	return new, new.isValid()
}

// addTimers adds timers to the query results
func (q *ExploreQuery) addTimers(timerName string, since time.Time, additionalTimers map[string]time.Duration) {
	if timerName != "" {
		q.Result.Timers[timerName] = time.Since(since)
	}

	if additionalTimers != nil {
		for k, v := range additionalTimers {
			q.Result.Timers[k] = v
		}
	}
}

// Execute implements the I2b2 MedCo query logic
func (q *ExploreQuery) Execute(queryType ExploreQueryType) (err error) {
	overallTimer := time.Now()
	var timer time.Time

	// todo: breakdown in i2b2 / count / patient list

	err = q.isValid()
	if err != nil {
		return
	}

	// tag query terms
	timer = time.Now()
	taggedQueryTerms, ddtTimers, err := unlynx.DDTagValues(q.ID, q.getEncQueryTerms())
	if err != nil {
		return
	}
	q.addTimers("medco-connector-DDT", timer, ddtTimers)
	logrus.Info(q.ID, ": tagged ", len(taggedQueryTerms), " elements with unlynx")

	// i2b2 PSM query with tagged items
	timer = time.Now()
	panelsItemKeys, panelsIsNot, err := q.getI2b2PsmQueryTerms(taggedQueryTerms)
	if err != nil {
		return
	}

	patientCount, patientSetID, err := i2b2.ExecutePsmQuery(q.ID, panelsItemKeys, panelsIsNot)
	if err != nil {
		return
	}
	q.addTimers("medco-connector-i2b2-PSM", timer, nil)
	logrus.Info(q.ID, ": got ", patientCount, " in patient set ", patientSetID, " with i2b2")

	// i2b2 PDO query to get the dummy flags
	timer = time.Now()
	patientIDs, patientDummyFlags, err := i2b2.GetPatientSet(patientSetID)
	if err != nil {
		return
	}
	q.addTimers("medco-connector-i2b2-PDO", timer, nil)
	logrus.Info(q.ID, ": got ", len(patientIDs), " patient IDs and ", len(patientDummyFlags), " dummy flags with i2b2")

	// aggregate patient dummy flags
	timer = time.Now()
	aggPatientFlags, err := unlynx.LocallyAggregateValues(patientDummyFlags)
	if err != nil {
		return
	}
	q.addTimers("medco-connector-local-agg", timer, nil)

	// compute and key switch count (returns optionally global aggregate or shuffled results)
	timer = time.Now()
	var encCount string
	var ksCountTimers map[string]time.Duration
	if queryType.CountType == Global {
		logrus.Info(q.ID, ": global aggregate requested")
		encCount, ksCountTimers, err = unlynx.AggregateAndKeySwitchValue(q.ID, aggPatientFlags, q.Query.UserPublicKey)
	} else if queryType.Shuffled {
		logrus.Info(q.ID, ": count per site requested, shuffle enabled")
		encCount, ksCountTimers, err = unlynx.ShuffleAndKeySwitchValue(q.ID, aggPatientFlags, q.Query.UserPublicKey)
	} else {
		logrus.Info(q.ID, ": count per site requested, shuffle disabled")
		encCount, ksCountTimers, err = unlynx.KeySwitchValue(q.ID, aggPatientFlags, q.Query.UserPublicKey)
	}
	if err != nil {
		return
	}
	q.addTimers("medco-connector-unlynx-key-switch-count", timer, ksCountTimers)

	// optionally obfuscate the count
	if queryType.Obfuscated {
		logrus.Info(q.ID, ": (local) obfuscation requested")
		timer = time.Now()
		encCount, err = unlynx.LocallyObfuscateValue(encCount, 5, q.Query.UserPublicKey) // todo: fixed distribution to make dynamic
		if err != nil {
			return
		}
		q.addTimers("medco-connector-local-obfuscation", timer, nil)
	}

	logrus.Info(q.ID, ": processed count")
	q.Result.EncCount = encCount

	// optionally prepare the patient list
	if queryType.PatientList {
		logrus.Info(q.ID, ": patient list requested")

		if len(patientIDs) == 0 {
			logrus.Info(q.ID, ": empty patient list. Skipping masking and key switching")
		} else {
			// mask patient IDs
			timer = time.Now()
			maskedPatientIDs, err := q.maskPatientIDs(patientIDs, patientDummyFlags)
			if err != nil {
				return err
			}

			logrus.Info(q.ID, ": masked ", len(maskedPatientIDs), " patient IDs")
			q.addTimers("medco-connector-local-patient-list-masking", timer, nil)

			// key switch the masked patient IDs
			timer = time.Now()
			ksMaskedPatientIDs, ksPatientListTimers, err := unlynx.KeySwitchValues(q.ID, maskedPatientIDs, q.Query.UserPublicKey)
			if err != nil {
				return err
			}
			q.addTimers("medco-connector-unlynx-key-switch-patient-list", timer, ksPatientListTimers)
			q.Result.EncPatientList = ksMaskedPatientIDs
			logrus.Info(q.ID, ": key switched patient IDs")
		}
	}

	q.addTimers("medco-connector-overall", overallTimer, nil)
	return
}

// maskPatientIDs multiplies homomorphically patient IDs with their dummy flags to mask the dummy patients
func (q *ExploreQuery) maskPatientIDs(patientIDs []string, patientDummyFlags []string) (maskedPatientIDs []string, err error) {

	if len(patientIDs) != len(patientDummyFlags) {
		err = errors.New("patient IDs and dummy flags do not have matching lengths")
		logrus.Error(err)
		return
	}

	for idx, patientID := range patientIDs {

		patientIDInt, err := strconv.ParseInt(patientID, 10, 64)
		if err != nil {
			logrus.Error("error parsing patient ID " + patientID + " as an integer")
			return nil, err
		}

		maskedPatientID, err := unlynx.LocallyMultiplyScalar(patientDummyFlags[idx], patientIDInt)
		if err != nil {
			return nil, err
		}

		maskedPatientIDs = append(maskedPatientIDs, maskedPatientID)
	}

	return
}

func (q *ExploreQuery) getEncQueryTerms() (encQueryTerms []string) {
	for _, panel := range q.Query.Panels {
		for _, item := range panel.Items {
			if *item.Encrypted {
				encQueryTerms = append(encQueryTerms, *item.QueryTerm)
			}
		}
	}
	return
}

func (q *ExploreQuery) getI2b2PsmQueryTerms(taggedQueryTerms map[string]string) (panelsItemKeys [][]string, panelsIsNot []bool, err error) {
	for panelIdx, panel := range q.Query.Panels {
		panelsIsNot = append(panelsIsNot, *panel.Not)

		panelsItemKeys = append(panelsItemKeys, []string{})
		for _, item := range panel.Items {
			var itemKey string
			if *item.Encrypted {

				if tag, ok := taggedQueryTerms[*item.QueryTerm]; ok {
					itemKey = `\\SENSITIVE_TAGGED\medco\tagged\` + tag + `\`
				} else {
					err = errors.New("query error: encrypted term does not have corresponding tag")
					logrus.Error(err)
					return
				}

			} else {
				itemKey = *item.QueryTerm
			}
			panelsItemKeys[panelIdx] = append(panelsItemKeys[panelIdx], itemKey)
		}
	}
	return
}

// isValid checks the validity of the query
func (q *ExploreQuery) isValid() (err error) {
	if len(q.ID) == 0 || q.Query == nil || len(q.Query.UserPublicKey) == 0 {
		err = errors.New("Query " + q.ID + " is invalid")
		logrus.Error(err)
	}
	return
}
