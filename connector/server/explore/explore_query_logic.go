package medcoserver

import (
	"fmt"
	"strconv"
	"time"

	medcomodels "github.com/ldsec/medco/connector/models"

	querytoolsserver "github.com/ldsec/medco/connector/server/querytools"

	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/ldsec/medco/connector/wrappers/i2b2"
	"github.com/ldsec/medco/connector/wrappers/unlynx"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// todo: log query (with associated status)
// todo: put user + query type + unique ID in query id

// ExploreQuery represents an i2b2-MedCo query to be executed
type ExploreQuery struct {
	ID     string
	Query  *models.ExploreQuery
	User   *models.User
	Result struct {
		EncCount       string
		EncPatientList []string
		Timers         medcomodels.Timers
		QueryID        int
		PatientSetID   int
	}
}

// NewExploreQuery creates a new query object and checks its validity
func NewExploreQuery(queryName string, query *models.ExploreQuery, user *models.User) (new ExploreQuery, err error) {
	new = ExploreQuery{
		ID:    queryName,
		Query: query,
		User:  user,
	}
	new.Result.Timers = medcomodels.NewTimers()
	return new, new.isValid()
}

// Execute implements the I2b2 MedCo query logic
func (q *ExploreQuery) Execute(queryType ExploreQueryType) (err error) {
	overallTimer := time.Now()
	var timer time.Time

	// create medco connector result instance
	queryDefinition, err := q.Query.MarshalBinary()
	if err != nil {
		err = fmt.Errorf("while marshalling query: %s", err.Error())
		return
	}
	logrus.Info("Creating Explore Result instance")
	timer = time.Now()
	queryID, err := querytoolsserver.InsertExploreResultInstance(q.User.ID, q.ID, string(queryDefinition))
	q.Result.Timers.AddTimers("medco-connector-create-result-instance", timer, nil)
	logrus.Infof("Created Explore Result Instance : %d", queryID)
	if err != nil {
		err = fmt.Errorf("while inserting explore result instance: %s", err.Error())
		return
	}
	defer func() {
		if err != nil {
			logrus.Info("Updating Explore Result instance with error status")
			qtError := querytoolsserver.UpdateErrorExploreResultInstance(queryID)
			if qtError != nil {
				err = fmt.Errorf("while inserting a status error in result instance table: %s", qtError.Error())
			} else {
				logrus.Info("Updating Explore Result instance with error status")
			}
		}
	}()

	// todo: breakdown in i2b2 / count / patient list

	err = q.isValid()
	if err != nil {
		err = fmt.Errorf("while checking validity: %s", err.Error())
		return
	}

	// tag query terms
	taggedQueryTerms := make(map[string]string)
	if encQueryTerms := q.getEncQueryTerms(); len(encQueryTerms) > 0 {
		var ddtTimers map[string]time.Duration
		timer = time.Now()
		taggedQueryTerms, ddtTimers, err = unlynx.DDTagValues(q.ID, encQueryTerms)
		if err != nil {
			return
		}
		q.Result.Timers.AddTimers("medco-connector-DDT", timer, ddtTimers)
		logrus.Info(q.ID, ": tagged ", len(taggedQueryTerms), " elements with unlynx")
	}

	// i2b2 PSM query with tagged items
	timer = time.Now()
	err = q.convertI2b2PsmQueryPanels(taggedQueryTerms)
	if err != nil {
		err = fmt.Errorf("while converting I2B2 panels: %s", err.Error())
		return
	}

	patientCount, patientSetID, err := i2b2.ExecutePsmQuery(q.ID, q.Query.Panels, q.Query.QueryTimingSequence, q.Query.QueryTiming)
	if err != nil {
		err = fmt.Errorf("during I2B2 PSM query exection: %s", err.Error())
		return
	}

	q.Result.Timers.AddTimers("medco-connector-i2b2-PSM", timer, nil)
	logrus.Info(q.ID, ": got ", patientCount, " in patient set ", patientSetID, " with i2b2")

	// i2b2 PDO query to get the dummy flags
	timer = time.Now()
	patientIDs, patientDummyFlags, err := i2b2.GetPatientSet(patientSetID, true)
	if err != nil {
		err = fmt.Errorf("during I2B2 patient set query exection: %s", err.Error())
		return
	}
	q.Result.Timers.AddTimers("medco-connector-i2b2-PDO", timer, nil)
	logrus.Info(q.ID, ": got ", len(patientIDs), " patient IDs and ", len(patientDummyFlags), " dummy flags with i2b2")

	// aggregate patient dummy flags
	timer = time.Now()
	aggPatientFlags, err := unlynx.LocallyAggregateValues(patientDummyFlags)
	if err != nil {
		err = fmt.Errorf("during local aggregation %s", err.Error())
		return
	}
	q.Result.Timers.AddTimers("medco-connector-local-agg", timer, nil)

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
		err = fmt.Errorf("during key switch/shuffle operation: %s", err.Error())
		return
	}
	q.Result.Timers.AddTimers("medco-connector-unlynx-key-switch-count", timer, ksCountTimers)

	// optionally obfuscate the count
	if queryType.Obfuscated {
		logrus.Info(q.ID, ": (local) obfuscation requested")
		timer = time.Now()
		encCount, err = unlynx.LocallyObfuscateValue(encCount, 5, string(q.Query.UserPublicKey)) // todo: fixed distribution to make dynamic
		if err != nil {
			err = fmt.Errorf("during key obfuscation operation: %s", err.Error())
			return
		}
		q.Result.Timers.AddTimers("medco-connector-local-obfuscation", timer, nil)
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
				err = fmt.Errorf("while producing patient masks: %s", err.Error())
				return err
			}

			logrus.Info(q.ID, ": masked ", len(maskedPatientIDs), " patient IDs")
			q.Result.Timers.AddTimers("medco-connector-local-patient-list-masking", timer, nil)

			// key switch the masked patient IDs
			timer = time.Now()
			ksMaskedPatientIDs, ksPatientListTimers, err := unlynx.KeySwitchValues(q.ID, maskedPatientIDs, string(q.Query.UserPublicKey))
			if err != nil {
				err = fmt.Errorf("while key-switching patient masks: %s", err.Error())
				return err
			}
			q.Result.Timers.AddTimers("medco-connector-unlynx-key-switch-patient-list", timer, ksPatientListTimers)
			q.Result.EncPatientList = ksMaskedPatientIDs
			logrus.Info(q.ID, ": key switched patient IDs")
		}
		q.Result.QueryID = queryID

		logrus.Info(q.ID, ": patient set ID requested")

		q.Result.PatientSetID, err = strconv.Atoi(patientSetID)
		if err != nil {
			return fmt.Errorf("while parsing patient set id: %v", err)
		}
	}
	//update medco connector result instance
	timer = time.Now()
	err = updateResultInstanceTable(queryID, patientCount, patientIDs, patientSetID)
	q.Result.Timers.AddTimers("medco-connector-update-result-instance", timer, nil)
	if err != nil {
		return err
	}

	q.Result.Timers.AddTimers("medco-connector-overall", overallTimer, nil)
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
		for _, item := range panel.ConceptItems {
			if *item.Encrypted {
				encQueryTerms = append(encQueryTerms, *item.QueryTerm)
			}
		}
	}
	return
}

func (q *ExploreQuery) convertI2b2PsmQueryPanels(taggedQueryTerms map[string]string) (err error) {

	for _, panel := range q.Query.Panels {
		for _, item := range panel.ConceptItems {
			if *item.Encrypted {
				if tag, ok := taggedQueryTerms[*item.QueryTerm]; ok {
					itemKey := `/SENSITIVE_TAGGED/medco/tagged/` + tag + `\`
					item.QueryTerm = &itemKey
				} else {
					err = errors.New("query error: encrypted term does not have corresponding tag")
					logrus.Error(err)
					return
				}
			}
		}

		// change the cohort name with the patient set ID
		if len(panel.CohortItems) > 0 {
			cohorts, err := querytoolsserver.GetSavedCohorts(q.User.ID, 100)
			if err != nil {
				return fmt.Errorf("while retrieving cohorts: %v", err)
			}
			if len(cohorts) == 0 {
				return fmt.Errorf("no cohorts for user: %v", q.User.ID)
			}
			for i, item := range panel.CohortItems {
				foundCohort := false
				for _, cohort := range cohorts {
					if cohort.CohortName == item {
						psID, err := querytoolsserver.GetPatientSetID(cohort.QueryID)
						if err != nil {
							return fmt.Errorf("while retrieving patient set ID: %v", cohort.QueryID)
						}
						panel.CohortItems[i] = "patient_set_coll_id:" + strconv.Itoa(psID)
						foundCohort = true
						break
					}
				}
				if !foundCohort {
					return fmt.Errorf("no cohort with name: %v found for user: %v", item, q.User.ID)
				}
			}
		}
	}

	return
}

// isValid checks the validity of the query
func (q *ExploreQuery) isValid() (err error) {
	if len(q.ID) == 0 || q.Query == nil || len(q.Query.UserPublicKey) == 0 {
		err = errors.New("query " + q.ID + " is invalid")
		logrus.Error(err)
	}

	nOfSequenceItems := len(q.Query.QueryTimingSequence)
	nOfConceptPanels := 0
	isASequenceQuery := nOfSequenceItems > 0

	containsOnlyConcepts := true
	//check number of sequential query
	for _, panel := range q.Query.Panels {
		containsOnlyConcepts = !(len(panel.CohortItems) > 0)
		if len(panel.ConceptItems) > 0 {
			nOfConceptPanels++
		}
	}
	if !containsOnlyConcepts && isASequenceQuery {
		err = fmt.Errorf("query %s is invalid: a sequential query can only have concept items in its selection panels, cohort item(s) was/were found", q.ID)
		logrus.Error(err)
		return
	}
	if (isASequenceQuery) && (nOfSequenceItems+1 != nOfConceptPanels) {
		err = fmt.Errorf("query %s is invalid: in a sequential query, the number of sequence information items + 1 must be equal to the number of panels: found %d sequence info items against %d panels", q.ID, nOfSequenceItems, nOfConceptPanels)
		logrus.Error(err)
	}
	return
}

func updateResultInstanceTable(queryID int, patientCount string, patientIDs []string, patientSetID string) (err error) {

	pCount, err := strconv.Atoi(patientCount)
	if err != nil {
		err = fmt.Errorf("while parsing integer from patient count string \"%s\": %s", patientCount, err.Error())
		return
	}
	pIDs := make([]int, len(patientIDs))
	for i, patientID := range patientIDs {
		pIDs[i], err = strconv.Atoi(patientID)
		if err != nil {
			err = fmt.Errorf("while parsing integer from patient ID string \"%s\": %s", patientID, err.Error())
			return
		}
	}

	patientSetIDNum, err := strconv.Atoi(patientSetID)
	if err != nil {
		err = fmt.Errorf("while parsing integer from patient set ID string \"%s\": %s", patientSetID, err.Error())
		return err
	}

	logrus.Info("Updating Explore Result instance")
	querytoolsserver.UpdateExploreResultInstance(queryID, pCount, pIDs, nil, &patientSetIDNum)
	if err != nil {
		err = fmt.Errorf("while updating result instance table: %s", err.Error())
		return
	}
	logrus.Info("Updated Explore Result instance")
	return
}
