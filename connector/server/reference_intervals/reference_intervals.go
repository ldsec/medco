package referenceintervalserver

import (
	"encoding/base64"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"time"

	medcomodels "github.com/ldsec/medco/connector/models"
	"github.com/ldsec/medco/connector/restapi/models"
	"github.com/ldsec/medco/connector/restapi/server/operations/explore_statistics"
	medcoserver "github.com/ldsec/medco/connector/server/explore"
	querytoolsserver "github.com/ldsec/medco/connector/server/querytools"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	utilserver "github.com/ldsec/medco/connector/util/server"
	"github.com/ldsec/medco/connector/wrappers/i2b2"
	"github.com/ldsec/medco/connector/wrappers/unlynx"
)

// Interval is a structure containing the lower bound and higher bound of an interval.
// The lower bound is inclusive and the higher bound is exclusive: [lower bound, higher bound[
type Interval struct {
	LowerBound  float64
	HigherBound float64
	EncCount    string // contains the count of subject in this interval
}

// StatsResult represents the information necessary to build a histogram about the observations of a specific analyte
type StatsResult struct {
	Intervals []*Interval
	Unit      string
	//concept or modifier name
	AnalyteName string
	//Timers for the construction of an individual histogram
	Timers medcomodels.Timers
}

// Query holds the ID of the survival analysis, its parameters and a pointer to its results
type Query struct {
	UserID        string
	UserPublicKey string
	//The name of the query given by the front end. It is used to distinguish this query from others.
	QueryName string

	QueryType medcoserver.ExploreQueryType

	// The clear list of patients from the cohort created from the inclusion and exclusion criterias.
	PatientsIDs []int64
	// Whether or not the explore query instance was created in the DB
	InstantiatedRecord bool
	//this variable is true if the
	isPanelEmpty bool
	// I2B2 panels defining the population upon which the analysis takes place. This is basically a constraint on the properties of the population.
	Panels []*models.Panel
	// query timing of the explore query defined by the panels
	QueryTiming models.Timing
	Concepts    []string
	// The bucket size for each analyte.  In the future there will be one such value for each analyte
	BucketSize float64
	// The global minimal observation for each analyte. In the future there will be one such value for each analyte
	MinObservation float64
	Modifiers      []*explore_statistics.ExploreStatisticsParamsBodyModifiersItems0 //TODO export this class out of the survival package make it a common thing
	Response       struct {
		medcoserver.PatientSetResult
		//Timers for what happens outside of the construction of the histograms
		GlobalTimers medcomodels.Timers
		Results      []*StatsResult
	}
}

// NewQuery query constructor for explore statistics query
func NewQuery(
	UserID string,
	params explore_statistics.ExploreStatisticsBody,
) (q *Query, err error) {

	if params.BucketSize <= 0 {
		err := fmt.Errorf("the size of each interval specified in the parameters must be strictly greater than 0")
		return nil, err
	}

	res := &Query{
		InstantiatedRecord: false,
		UserID:             UserID,
		UserPublicKey:      params.UserPublicKey,
		Panels:             params.CohortDefinition.Panels,
		QueryTiming:        params.CohortDefinition.QueryTiming,
		QueryName:          params.ID,
		Concepts:           params.Concepts,
		BucketSize:         params.BucketSize,
		MinObservation:     params.MinObservation,
		Modifiers:          params.Modifiers,
		isPanelEmpty:       params.CohortDefinition.IsPanelEmpty,
	}

	res.Response.GlobalTimers = make(map[string]time.Duration)

	return res, nil
}

func outlierRemoval(observations []QueryResult) (outputObs []QueryResult, err error) {
	// |Z| = | (x - x bar) / S | >= 6 (where S is std deviation)
	outputObs = observations
	return
}

type CohortInformation struct {
	PatientIDs   []int64 // a list of patient IDs
	IsEmptyPanel bool    // True iff the client set no constraint in the panel definition. The population selected would in this case consist of all patients.
}

// Execute runs the explore statistics query.
// The histogram is created by cutting the space of observations in equal sized interval. The number of interval is defined by the parameters of the query.
// The observations are then fetched from the database and classified in each interval. Then the count of observations is determined per interval.
// Those counts are then aggregated between all nodes. And the result is returned to the user.
func (q *Query) Execute(principal *models.User) (err error) {

	timer := time.Now()

	conceptsCodesAndNames, modifiersCodesAndNames, timers, patientsInfos, err := q.prepareArguments(principal)

	defer func() {
		// If an error occured and the current statistics query resulted in the creation of an explore query this deferred function will set an error status
		// on the DB record for the explore query underlying this explore statistics query.
		if err != nil && q.InstantiatedRecord {
			logrus.Info("Updating the Explore Result instance with error status that is underlying the explore statistics query")
			qtError := querytoolsserver.UpdateErrorExploreResultInstance(q.Response.QueryID)
			if qtError != nil {
				err = fmt.Errorf("while inserting a status error in result instance table: %s", qtError.Error())
			} else {
				logrus.Info("Updating Explore Result instance with error status")
			}
		}
	}()

	if err != nil {
		modStr := ""
		for _, mod := range q.Modifiers {
			modStr += *mod.ModifierKey + ", "
		}

		err = fmt.Errorf("while retrieving concept codes and patient indices: for concepts %s ... modifiers %s ... error: %s ", q.Concepts, modStr, err.Error())
		return err
	}
	q.Response.GlobalTimers.AddTimers("", timer, timers)

	//TODO define the analysis depending on a bucket width and not a number of intervals

	// A wait group that will allow us to wait for all goroutines created in this function to finish before returning from the Execute function.
	waitGroup := &sync.WaitGroup{}

	nbCodes := len(conceptsCodesAndNames) + len(modifiersCodesAndNames)
	waitGroup.Add(nbCodes)

	//each analyte observations are processed and stored in a channel within statsChannels
	statsChannels := make([]chan *StatsResult, nbCodes)
	errChan := make(chan error)
	signal := make(chan struct{})

	//this function is an abstraction for the processing of retreving observations and then processing them for a concept or modifier (depends on RetrieveObservations)
	processMedicalConcept := func(index int, codeNamePair codeAndName, RetrieveObservations func(string, CohortInformation, float64) (results []QueryResult, err error)) {
		conceptTimer := time.Now()

		defer waitGroup.Done()

		cohortInfo := CohortInformation{
			PatientIDs:   q.PatientsIDs,
			IsEmptyPanel: q.isPanelEmpty,
		}
		conceptObservations, err := RetrieveObservations(codeNamePair.Code, cohortInfo, q.MinObservation)
		if err != nil {
			errChan <- err
			return
		}

		cleanObservations, err := outlierRemoval(conceptObservations)

		if err != nil {
			return
		}

		// TODO define bucket width depending on the observations of the concept
		statsResults, err := q.processObservations(q.BucketSize, cleanObservations, conceptTimer, codeNamePair.Code)

		if err != nil {
			errChan <- err
		}

		statsResults.AnalyteName = codeNamePair.Name

		statsChannels[index] <- statsResults

	}

	for i, concept := range conceptsCodesAndNames {
		statsChannels[i] = make(chan *StatsResult, 1)
		go processMedicalConcept(i, concept, RetrieveObservationsForConcept)
	}

	nbConcepts := len(conceptsCodesAndNames)
	for j, modifier := range modifiersCodesAndNames {

		index := nbConcepts + j
		statsChannels[index] = make(chan *StatsResult, 1)
		go processMedicalConcept(index, modifier, RetrieveObservationsForModifier)
	}

	go func() {
		waitGroup.Wait()
		signal <- struct{}{}
	}()

	select {
	case err = <-errChan:
		return
	case <-signal:
		break
	}

	// We fetch the histogram information for each analyte within each channel that contains such information and append this information to the HTTP response.
	for _, statResultChannel := range statsChannels {
		q.Response.Results = append(q.Response.Results, <-statResultChannel)
	}

	cohortInt := make([]int, 0)
	for i := 0; i < len(q.PatientsIDs); i++ {
		cohortInt = append(cohortInt, int(q.PatientsIDs[i]))
	}
	querytoolsserver.UpdateExploreResultInstance(patientsInfos.QueryID, len(cohortInt), cohortInt, nil, nil)

	err = q.locallyAggregatePatientCount(patientsInfos)
	if err != nil {
		return
	}
	// Process the answers given by the database about the cohort and fill the query
	q.Response.ProcessPatientsList(q.QueryType, q.QueryName, patientsInfos, q.UserPublicKey, &q.Response.GlobalTimers)

	return

}

func (q *Query) locallyAggregatePatientCount(patientsInfos medcoserver.LocalPatientsInfos) (err error) {

	// aggregate patient dummy flags
	timer := time.Now()
	aggPatientFlags, err := unlynx.LocallyAggregateValues(patientsInfos.PatientDummyFlags)
	if err != nil {
		err = fmt.Errorf("during local aggregation %s", err.Error())
		return
	}

	q.Response.GlobalTimers.AddTimers("medco-connector-local-agg", timer, nil)

	// compute and key switch count (returns optionally global aggregate or shuffled results)
	timer = time.Now()
	var encCount string
	var ksCountTimers map[string]time.Duration
	if q.QueryType.PatientList && !q.QueryType.Obfuscated {
		logrus.Info(q.QueryName, ": count per site requested, shuffle disabled")
		encCount, ksCountTimers, err = unlynx.KeySwitchValue(q.QueryName, aggPatientFlags, q.UserPublicKey)
	}
	if err != nil {
		err = fmt.Errorf("during key switch/shuffle operation: %s", err.Error())
		return
	}
	q.Response.GlobalTimers.AddTimers("medco-connector-unlynx-key-switch-count", timer, ksCountTimers)
	logrus.Info(q.QueryName, ": processed count")
	q.Response.EncCount = encCount
	return
}

func (q *Query) locallyProcessObservations(bucketSize float64, queryResults []QueryResult,
	timer time.Time, code string, encryptionMethod func(int64) (string, error)) (encCounts []string, statsResults *StatsResult, err error) {

	//get the minimum and maximum value of the concepts
	var maxResult QueryResult = queryResults[0]
	for _, r := range queryResults {
		if r.NumericValue > maxResult.NumericValue {
			maxResult = r
		}
	}

	logrus.Info(" max value ", maxResult.NumericValue)

	// defining the number of intervals depending on the maximum and minimum observations and the bucket size
	nbIntervals := int(math.Ceil((maxResult.NumericValue - q.MinObservation) / bucketSize))

	logrus.Info("Query results contains ", len(queryResults), " records")

	if len(queryResults) < 2 {
		err = fmt.Errorf("not enough concepts present in order to define buckets")
		return
	}

	statsResults = &StatsResult{
		Intervals: make([]*Interval, nbIntervals),
		Unit:      queryResults[0].Unit,
		Timers:    make(map[string]time.Duration),
		// TODO Later on we will probably have to perform unit conversion. One possibility would be to fetch the metadataXML of a concept to see what
		// are the conversion rules for that concept. For now we take the hypothesis that everything is under the same unit
		// c.f. https://community.i2b2.org/wiki/display/DevForum/Metadata+XML+for+Medication+Modifiers

		//another option is to convert all observations for a same concept to the same unit during the ETL phase.
	}

	//from the minimum and maximum value of the select concept we determine the boundaries of the different buckets.
	step := bucketSize
	logrus.Info("Steps equals ", step)

	current := q.MinObservation
	logrus.Debugf("processObservations: Number of interval = %d", nbIntervals)

	for i := 0; i < nbIntervals; i++ {
		statsResults.Intervals[i] = &Interval{}
		interval := statsResults.Intervals[i]
		logrus.Info("Setting interval bounds. Before ", interval.LowerBound, interval.HigherBound, " After ", current, current+step, step)
		interval.LowerBound = current //TODO trim the zeroes when sending that in json format
		interval.HigherBound = current + step

		current += step
	}

	/* In the following lines of code we group the query results in different buckets depending on their numerical values. We count the number of concepts
	 * that belong to the differents intervals and cypher this value with the cothority key.
	 */

	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(nbIntervals)

	channels := make([]chan struct {
		encCount *string
		medcomodels.Timers
	}, nbIntervals)

	errChan := make(chan error)
	signal := make(chan struct{})

	//TODO In the future we will have to save the queries in the medco database in order to reproduce their results.
	for i, interval := range statsResults.Intervals {
		logrus.Debugf("Starting the processing of the interval with index %d", i)
		if interval.LowerBound >= interval.HigherBound {
			err := fmt.Errorf("the lower bound of the interval #%d is greater than the higher bound: %f >= %f", i, interval.LowerBound, interval.HigherBound)
			errChan <- err
			break
		}

		channels[i] = make(chan struct {
			encCount *string
			medcomodels.Timers
		}, 1)

		logrus.Debugf("processObservations: Assigned struct to channel with index %d", i)

		go func(i int, interval *Interval) {
			defer waitGroup.Done()
			timers := medcomodels.NewTimers()

			count := 0

			logrus.Debugf("About to count the number of observations that fit in interval %d", i)
			//counting the number of numerical values that belong to the [lowerbound, higherbound[ interval.
			for _, queryResult := range queryResults {
				isLastInterval := maxResult.NumericValue == interval.HigherBound
				smallerThanHigherBound :=
					(isLastInterval && queryResult.NumericValue <= interval.HigherBound) ||
						(!isLastInterval && queryResult.NumericValue < interval.HigherBound)

				if queryResult.NumericValue >= interval.LowerBound && smallerThanHigherBound {
					count++
				}
			}

			timer = time.Now()

			logrus.Debugf("Count for bucket [ %f , %f] is %d", interval.LowerBound, interval.HigherBound, count)

			//encrypt the interval count of observations with the collective authority key of medco nodes.
			encCount, err := encryptionMethod(int64(count))

			logrus.Debugf("Count %d,  encrypted %s ", count, encCount)
			timers.AddTimers(fmt.Sprintf("medco-connector-encrypt-interval-count-group%d", i), timer, nil)
			if err != nil {
				err = fmt.Errorf("while encrypting the count of an interval of the future reference interval: %s", err.Error())
				errChan <- err
				return
			}

			logrus.Debugf("Sending information encrypted count information to channel %d", i)
			channels[i] <- struct {
				encCount *string
				medcomodels.Timers
			}{&encCount, timers}

			logrus.Debugf("Done sending information encrypted count information to channel %d", i)
		}(i, interval)

	}
	go func() {
		waitGroup.Wait()
		signal <- struct{}{}
	}()

	select {
	case err = <-errChan:
		return
	case <-signal:
		break
	}

	logrus.Debugf("Before receiving each interval count for the concept/modifier with code %s", code)
	encCounts = make([]string, 0, len(channels))
	for i, channel := range channels {
		chanResult := <-channel

		logrus.Debugf("Receiving the encrypted count in the channel with index %d, %s", i, *chanResult.encCount)

		encCounts = append(encCounts, *chanResult.encCount)
		statsResults.Timers.AddTimers("", timer, chanResult.Timers)
	}

	return
}

//@param code is the code of the modifier or concept whose observations are being processed
func (q *Query) processObservations(bucketSize float64, queryResults []QueryResult, timer time.Time, code string) (statsResults *StatsResult, err error) {
	var encCounts []string
	encCounts, statsResults, err = q.locallyProcessObservations(bucketSize, queryResults, timer, code, unlynx.EncryptWithCothorityKey)
	if err != nil {
		return
	}

	// aggregate and key switch locally encrypted counts of each bucket
	timer = time.Now()
	var aggregationTimers medcomodels.Timers
	var aggValues []string

	qName := q.QueryName + "_AGG_AND_KEYSWITCH_" + code
	logrus.Info("Launching the encrypted aggregation: explore stats", len(encCounts), " with name ", qName)

	aggValues, aggregationTimers, err = unlynx.AggregateAndKeySwitchValues(qName, encCounts, q.UserPublicKey)

	if err != nil {
		err = fmt.Errorf("during aggregation and keyswitch: %s", err.Error())
		logrus.Errorf("%s", err.Error())
		return
	}

	logrus.Debugf("After explore stats aggregate and key switch values")
	//assign the encrypted count to the matching interval
	for i, interval := range statsResults.Intervals {
		interval.EncCount = aggValues[i]
	}

	statsResults.Timers.AddTimers("medco-connector-aggregate-and-key-switch", timer, aggregationTimers)
	return
}

// Validate checks members of a Query instance for early error detection.
// Heading and trailing spaces are silently trimmed.
// If any other wrong member can be defaulted, a warning message is printed, otherwise an error is returned.
func (q *Query) Validate() error {

	q.QueryName = strings.TrimSpace(q.QueryName)
	if q.QueryName == "" {
		return fmt.Errorf("empty query name")
	}

	for i, concept := range q.Concepts {
		q.Concepts[i] = strings.TrimSpace(concept)
		if concept == "" {
			return fmt.Errorf("emtpy concept path, queryID: %s", q.QueryName)
		}

	}

	for _, modifier := range q.Modifiers {
		if modifier == nil {
			continue
		}

		*modifier.ModifierKey = strings.TrimSpace(*modifier.ModifierKey)
		if *modifier.ModifierKey == "" {
			return fmt.Errorf("empty modifier key, queryID: %s", q.QueryName)
		}
		*modifier.AppliedPath = strings.TrimSpace(*modifier.AppliedPath)
		if *modifier.AppliedPath == "" {
			return fmt.Errorf(
				"empty modifier applied path, queryID: %s,  modifier key: %s",
				q.QueryName,
				*modifier.ModifierKey,
			)
		}

	}

	q.UserID = strings.TrimSpace(q.UserID)
	if q.UserID == "" {
		return fmt.Errorf("empty user name, queryID: %s", q.QueryName)
	}

	q.UserPublicKey = strings.TrimSpace(q.UserPublicKey)
	if q.UserPublicKey == "" {
		return fmt.Errorf("empty user public keyqueryID: %s", q.QueryName)
	}
	_, err := base64.URLEncoding.DecodeString(q.UserPublicKey)
	if err != nil {
		return fmt.Errorf("user public key is not valid against the alternate RFC4648 base64 for URL: %s; queryID: %s", err.Error(), q.QueryName)
	}
	return nil
}

type codeAndName struct {
	Code string
	Name string
}

// prepareArguments retrieves concept codes and patients that will be used as the arguments of direct SQL call
func (q *Query) prepareArguments(principal *models.User) (
	conceptsInfos,
	modifiersInfos []codeAndName,
	timers medcomodels.Timers,
	patientsInfos medcoserver.LocalPatientsInfos,
	err error,
) {
	timers = make(map[string]time.Duration)
	// --- cohort patient list
	timer := time.Now()

	logrus.Info("Fetching patients for the explore statistics using an i2b2 query")

	// Creating and executing an explore query in order to fetch the local patients that match the constraints specified by the panels. This explore request will return
	// the patients that define the population upon which the explore statistics analysis will base itself.
	exploreQueryParam := &models.ExploreQuery{
		Panels:        q.Panels,
		QueryTiming:   q.QueryTiming,
		UserPublicKey: q.UserPublicKey,
	}
	exploreQuery, err := medcoserver.NewExploreQuery(q.QueryName+"_explore", exploreQueryParam, principal)

	if err != nil {
		err = fmt.Errorf("error when creating a new explore query in the prepareArgument method of the explore statistics. %s", err)
		return
	}

	if !q.isPanelEmpty {
		patientsInfos, err = exploreQuery.FetchLocalPatients(timer) //before we used querytoolsserver.GetPatientList
	}

	q.InstantiatedRecord = true

	if err != nil {
		logrus.Error("error while getting patient list")
		return
	}

	if len(patientsInfos.PatientIDs) == 0 && !q.isPanelEmpty {
		err = fmt.Errorf("zero patients in the cohort for non empty constraint")
		return
	}

	q.PatientsIDs = make([]int64, 0)

	for _, patientID := range patientsInfos.PatientIDs {
		var patientIDInt int64
		patientIDInt, err = strconv.ParseInt(patientID, 10, 64)
		if err != nil {
			err = fmt.Errorf("while parsing integer from patient ID string \"%s\": %s", patientID, err.Error())
			return
		}

		q.PatientsIDs = append(q.PatientsIDs, patientIDInt)
	}

	timers.AddTimers("medco-connector-get-patient-list", timer, nil)
	logrus.Debug("got patients for the explore statistics cohort: ", patientsInfos.PatientIDs)

	// --- get concept and modifier codes from the ontology
	logrus.Info("get concept and modifier codes")
	err = utilserver.I2B2DBConnection.Ping()
	if err != nil {
		err = fmt.Errorf("while connecting to clear project database: %s", err.Error())
		return
	}

	waitGroup := &sync.WaitGroup{}
	logrus.Info("wait group length = ", len(q.Concepts)+len(q.Modifiers))
	waitGroup.Add(len(q.Concepts) + len(q.Modifiers))

	signal := make(chan struct{})
	// Each channel in this array will contain the code and name of either a modifier or a concept.
	conceptsChannels := make([]chan codeAndName, len(q.Concepts))
	modifiersChannels := make([]chan codeAndName, len(q.Modifiers))
	errChan := make(chan error)

	if q.Concepts != nil && len(q.Concepts) > 0 {
		for i, concept := range q.Concepts {
			conceptsChannels[i] = make(chan codeAndName, 1)
			go func(concept string, index int) {
				defer waitGroup.Done()

				var conceptCode, conceptName string

				//fetch the code and name of the concept using an i2b2 query.
				conceptCode, conceptName, err = getCodeAndName(concept)
				if err != nil {
					errChan <- fmt.Errorf("while retrieving concept code: %s", err.Error())
					return
				}
				conceptsChannels[index] <- codeAndName{Code: conceptCode, Name: conceptName}
				logrus.Info("Got concept code for concept ", conceptName, " code = ", conceptCode)

			}(concept, i)
		}
	}

	if q.Modifiers != nil && len(q.Modifiers) > 0 {
		for i, modifier := range q.Modifiers {

			if modifier == nil {
				continue
			}
			modifiersChannels[i] = make(chan codeAndName, 1)
			go func(index int, modifierKey, appliedPath string) {
				defer waitGroup.Done()

				var modifierCode, modifierName string
				//fetch the code and name of the modifier using an i2b2 query.
				modifierCode, modifierName, err = getModifierCodeAndName(modifierKey, appliedPath)
				logrus.Debugf("get modifier code and name has returned %s, %s", modifierCode, modifierName)

				if err != nil {
					errChan <- fmt.Errorf("while retrieving modifier code: %s", err.Error())
					return
				}
				logrus.Debugf("About to pass modifier code and name to the channel of index %d, %s, %s", index, modifierCode, modifierName)
				modifiersChannels[index] <- codeAndName{Code: modifierCode, Name: modifierName}
				logrus.Info("Got modifier code, modifier name ", modifierName, " code = ", modifierCode)

			}(i, *modifier.ModifierKey, *modifier.AppliedPath)

		}
	}

	go func() {
		waitGroup.Wait()
		signal <- struct{}{}
	}()

	select {
	case err = <-errChan:
		return
	case <-signal:
		break
	}

	for _, channel := range conceptsChannels {
		conceptsInfos = append(conceptsInfos, <-channel)
	}

	for _, channel := range modifiersChannels {
		modifiersInfos = append(modifiersInfos, <-channel)
	}

	logrus.Info("got concept and modifier codes")
	return
}

// getCodeAndName takes the full path of a I2B2 concept and returns its code and name
func getCodeAndName(path string) (string, string, error) {
	logrus.Debugf("get code concept path %s", path)
	res, err := i2b2.GetOntologyConceptInfo(path)
	if err != nil {
		return "", "", err
	}
	if len(res) != 1 {
		return "", "", errors.Errorf("Result length of GetOntologyConceptInfo is expected to be 1. Got: %d", len(res))
	}

	if res[0].Code == "" {
		return "", "", errors.New("Code is empty")
	}

	if res[0].Name == "" {
		return "", "", errors.New("Concept name is empty")
	}

	logrus.Debugf("got concept code %s", res[0].Code)

	return res[0].Code, res[0].DisplayName, nil

}

// getModifierPath takes the full path of a I2B2 modifier and its applied path and returns its code and name
func getModifierCodeAndName(path string, appliedPath string) (string, string, error) {
	logrus.Debugf("get modifier code modifier path %s applied path %s", path, appliedPath)
	res, err := i2b2.GetOntologyModifierInfo(path, appliedPath)
	if err != nil {
		return "", "", err
	}

	if len(res) != 1 {
		return "", "", errors.Errorf("Result length of GetOntologyTermInfo is expected to be 1. Got: %d. "+
			"Is applied path %s available for modifier key %s ?", len(res), appliedPath, path)
	}
	if res[0].Code == "" {
		return "", "", errors.New("Code is empty")
	}

	if res[0].DisplayName == "" {
		return "", "", errors.New("Modifier name is empty")
	}

	logrus.Debugf("got modifier code %s", res[0].Code)

	return res[0].Code, res[0].DisplayName, nil
}
