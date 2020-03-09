package survivalserver

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	utilserver "github.com/ldsec/medco-connector/util/server"
	"github.com/ldsec/medco-connector/wrappers/i2b2"
	"github.com/ldsec/medco-connector/wrappers/unlynx"
	servicesmedco "github.com/ldsec/medco-unlynx/services"
	libunlynx "github.com/ldsec/unlynx/lib"
	"github.com/sirupsen/logrus"
)

// Query holds the ID of the survival analysis, its parameters and a pointer to its results
type Query struct {
	ID            string
	UserPublicKey string
	PatientSetID  string
	TimeCodes     []string

	//TODO also hide that
	Result *struct {
		Timers    map[string]time.Duration
		EncEvents map[string][2]string
	}

	spin *Spin
}

// NewQuery query constructor
func NewQuery(qID, pubKey, patientSetID string, timeCodes []string) *Query {
	return &Query{ID: qID, UserPublicKey: pubKey, PatientSetID: patientSetID, TimeCodes: timeCodes, spin: NewSpin()}
}

// GetID returns query ID
func (q *Query) GetID() string {
	return q.ID
}

//GetUserPublicKey returns the user public key
func (q *Query) GetUserPublicKey() string {
	return q.UserPublicKey
}

//GetPatientSetID returns the patient set ID
func (q *Query) GetPatientSetID() string {
	return q.PatientSetID
}

// GetTimeCodes returns the time points encryption ID encrypted with the collective authority key
func (q *Query) GetTimeCodes() []string {
	return q.TimeCodes
}

//SetResultMap sets the result map structure to return
func (q *Query) SetResultMap(resultMap map[string][2]string) error {
	q.spin.Lock()
	defer q.spin.Unlock()
	if q.Result == nil {
		q.Result = new(struct {
			Timers    map[string]time.Duration
			EncEvents map[string][2]string
		})
	}
	q.Result.EncEvents = resultMap
	return nil

}

// Execute is the function that translates a medco survival query in multiple calls on psql and unlynx
func (q *Query) Execute(batchNumber int) (err error) {
	var storage sync.Map

	err = directAccessDB.Ping()

	if err != nil {
		logrus.Error("Unable to ping database")
		return
	}

	var timePoints []string
	var patientSet []string

	encTimePoints := q.GetTimeCodes()
	if len(encTimePoints) == 0 {
		logrus.Panic("Unexpected empty list of time points")
	}

	timeCodesMap, _, err := NewTimeCodesMap(q.GetID(), encTimePoints)
	if err != nil {
		return
	}
	timePoints = timeCodesMap.GetTagIDList()
	tagIDToEncTimePoints := timeCodesMap.GetTagIDMap()

	patientSetID := q.GetPatientSetID()
	if patientSetID == "" {
		logrus.Panic("Unexpected null string for patient set ID")
	}
	patientSet, _, err = i2b2.GetPatientSet(patientSetID)

	err = NiceError(err)
	if err != nil {
		return
	}
	if len(patientSet) == 0 {
		//TODO magic numbers
		unlynxChannel := make(chan struct {
			event     *unlynxResult
			censoring *unlynxResult
		}, 1024)

		for _, encTimePoint := range tagIDToEncTimePoints {
			var zeroEvents string
			var zeroCensoring string
			zeroEvents, err = zeroPointEncryption()
			if err != nil {
				return
			}
			zeroCensoring, err = zeroPointEncryption()
			if err != nil {
				return
			}

			err = aggregateAndKeySwitchSend(q.GetID(), encTimePoint, zeroEvents, zeroCensoring, q.GetUserPublicKey(), unlynxChannel)
			if err != nil {
				return
			}

		}
		var receptionBarrier sync.WaitGroup
		receptionBarrier.Add(len(tagIDToEncTimePoints))
		for i := 0; i < len(tagIDToEncTimePoints); i++ {
			err = aggregateAndKeySwitchCollect(&storage, unlynxChannel, &receptionBarrier)
			if err != nil {
				return
			}
		}

		err = q.fillResult(&storage)

		return
	}

	if err != nil {
		return
	}

	exceptionHandler, err := NewExceptionHandler(1)
	err = NiceError(err)

	if err != nil {
		return
	}
	go func() {

		batches, err := NewBatchIterator(timePoints, batchNumber)

		err = NiceError(err)

		if err != nil {
			exceptionHandler.PushError(err)
			return
		}

		for !batches.Done() {
			batch := batches.Next()
			var waitGroup sync.WaitGroup
			//TODO magic numbers
			unlynxChannel := make(chan struct {
				event     *unlynxResult
				censoring *unlynxResult
			}, 1024)

			if naiveVersion {
				waitGroup.Add(len(batch))
				var events []string
				var censorings []string
				for _, timeCode := range batch {

					//error chans ??
					go func(timeCode string) {

						result := Result{}
						//TODO not elegant
						deferred := func(err error) {
							result.Error = err
							//storage.Store(tagIDToEncTimePoints[timeCode], result)
							if err != nil {
								exceptionHandler.PushError(err)
							}
							waitGroup.Done()
						}
						query := buildSingleQuery(patientSet, timeCode)

						rows, err := directAccessDB.Query(query)

						//logrus.Panic("\n\n\n                   :3\n\n")
						if err != nil {

							deferred(err)

							return
						}
						debugCount := 0

						for rows.Next() { //there should be at most one row here
							debugCount++
							var blob string
							err = rows.Scan(&blob)

							err = NiceError(err)
							if err != nil {
								deferred(err)

								return
							}
							event, censoring, err := breakBlob(blob)
							err = NiceError(err)
							if err != nil {
								deferred(err)
								return
							}
							events = append(events, event)
							censorings = append(censorings, censoring)
						}
						err = rows.Close()
						err = NiceError(err)
						if err != nil {
							deferred(err)

							return
						}
						if debugCount == 0 {
							err = errors.New("no row !!!" + query)
						}
						if err != nil {
							deferred(err)
							return
						}

						events, err := unlynx.LocallyAggregateValues(events)
						err = NiceError(err)
						if err != nil {
							deferred(err)

							return
						}
						censoringEvents, err := unlynx.LocallyAggregateValues(censorings)
						if err != nil {
							deferred(err)
							return

						}
						encTimeCode, ok := tagIDToEncTimePoints[timeCode]
						if !ok {
							//TODO check orth
							err = errors.New("the map tag -> encrypted ill-formed")
						}

						err = aggregateAndKeySwitchSend(q.GetID(), encTimeCode, events, censoringEvents, q.GetUserPublicKey(), unlynxChannel)
						if err != nil {
							deferred(err)
							return
						}
						err = NiceError(err)
						deferred(err)
						return
					}(timeCode)

				}
				errOccured := false
				var receptionBarrier sync.WaitGroup
				receptionBarrier.Add(len(batch))
				for range batch {
					err = aggregateAndKeySwitchCollect(&storage, unlynxChannel, &receptionBarrier)
					if err != nil {
						errOccured = true
						exceptionHandler.PushError(err)
						break
					}
				}
				if !errOccured {
					receptionBarrier.Wait()
				}

			} else {
				//does not make a  alot of sense to have  a wait group ad a go routine for only one instance here
				waitGroup.Add(1)

				go func() {
					result := Result{}

					defer waitGroup.Done()

					query := buildBatchQuery(patientSet, batch)

					rows, err := directAccessDB.Query(query)
					//debug
					if err != nil {
						err = errors.New(err.Error() + " " + query)
						return
					}
					err = NiceError(err)

					if err != nil {
						result.Error = err
						//TODO handle this kind of error
						storage.Store("", result)
						exceptionHandler.PushError(err)
						return
					}

					set := NewSet(len(batch))
					for _, timeCode := range batch {
						set.Add(timeCode)
					}

					// ------ aggregates and send aggregates for collective agg and key switching
					//   ---- 1) do this for the result encountered in the loval observation fact
					//   ---- 2) do this for remaining data (they will have an encrypted zero value)
					for rows.Next() {
						recipiens := &struct {
							TimeCode          string
							ConcatenatedBlobs string
						}{}
						err = rows.Scan(&(recipiens.TimeCode), &(recipiens.ConcatenatedBlobs))
						err = NiceError(err)
						set.Remove(recipiens.TimeCode)

						str := strings.Split(recipiens.ConcatenatedBlobs, interBlob)
						eventsOfInterest := make([]string, len(str))
						censoringEvents := make([]string, len(str))
						for idx, blob := range str {
							eventsOfInterest[idx], censoringEvents[idx], err = breakBlob(blob)
							err = NiceError(err)
							if err != nil {
								result.Error = err
								storage.Store(recipiens.TimeCode, result)
								rows.Close()
								exceptionHandler.PushError(err)
								return
							}

						}

						//err chans ?
						go func() { //call unlynx for the first kind of events

							events, err := unlynx.LocallyAggregateValues(eventsOfInterest)
							err = NiceError(err)
							if err != nil {
								exceptionHandler.PushError(err)
								return
							}
							censoringEvents, err := unlynx.LocallyAggregateValues(censoringEvents)
							err = NiceError(err)
							if err != nil {
								exceptionHandler.PushError(err)
								return

							}
							encTimeCode, ok := tagIDToEncTimePoints[recipiens.TimeCode]
							if !ok {
								//TODO orthograph
								err = errors.New("ill formed time codes map")
							}
							err = NiceError(err)
							if err != nil {
								exceptionHandler.PushError(err)
								return
							}

							err = aggregateAndKeySwitchSend(q.GetID(), encTimeCode, events, censoringEvents, q.GetUserPublicKey(), unlynxChannel)
							err = NiceError(err)
							if err != nil {
								exceptionHandler.PushError(err)
								return
							}
						}()
					}
					rows.Close()
					//for those that have not be found in this node

					//TODO this is beyound the unlynx batch barrier mechanism
					set.ForEach(func(key string) {

						go func() {
							zeroEvent, err := zeroPointEncryption()
							err = NiceError(err)
							if err != nil {
								exceptionHandler.PushError(err)

							}
							zeroCensoring, err := zeroPointEncryption()
							err = aggregateAndKeySwitchSend(`queryname`, key, zeroEvent, zeroCensoring, q.GetUserPublicKey(), unlynxChannel)
							err = NiceError(err)
							if err != nil {
								exceptionHandler.PushError(err)
							}

						}()
						//storage.Store(key, result)

					})

					// ------- collect the result of the  aggregate and key switch
					errorOccurred := false
					var receptionBarrier sync.WaitGroup
					receptionBarrier.Add(len(batch))
					for range batch {
						err = aggregateAndKeySwitchCollect(&storage, unlynxChannel, &receptionBarrier)
						err = NiceError(err)
						if err != nil {

							errorOccurred = true
							exceptionHandler.PushError(err)
							break
						}
					}
					if !errorOccurred {
						receptionBarrier.Wait()
					}

				}()

			}
			waitGroup.Wait()

		}
		//err chans !!!
		exceptionHandler.Finished()
		return
	}()
	err = exceptionHandler.WaitEndSignal(3000)
	if err != nil {
		return
	}
	err = q.fillResult(&storage)
	err = NiceError(err)
	if err != nil {
		return
	}
	return
}
func (q *Query) fillResult(storage *sync.Map) (err error) {

	targetMap := make(map[string][2]string)
	counter := 0
	storage.Range(func(timeCode /*string*/, events interface{} /*Result*/) bool {
		if timeCodeString, ok1 := timeCode.(string); ok1 {
			if eventStruct, ok2 := events.(Result); ok2 {

				targetMap[timeCodeString] = [2]string{eventStruct.EventValue, eventStruct.CensoringValue}
				counter++
				return true
			}
			logrus.Panic("Wrong type for map value")

		}
		logrus.Panic("Wrong type for map key")
		return false

	})
	err = q.SetResultMap(targetMap)

	return
}

type unlynxResult struct {
	Key   *string
	Value string
	Error error
}

func buildSingleQuery(patients []string, timeCode string) string {

	patientAcc := stringMapAndAdd(patients)
	res := `SELECT ` + blobColumn + ` FROM ` + schema + `.` + table + ` `
	res += `WHERE ` + timeConceptColumn + ` = '` + timeCode + `' AND ` + patientIDColumn + ` IN (` + patientAcc + `)`
	return res

}
func buildBatchQuery(patients []string, timeCodes []string) string {

	patientAcc := stringMapAndAdd(patients)
	timeAcc := stringMapAndAdd(timeCodes)
	res := `SELECT ` + timeConceptColumn + `, STRING_AGG(` + blobColumn + `,'` + interBlob + `') FROM ` + schema + `.` + table + ` `
	res += `WHERE ` + timeConceptColumn + ` IN (` + timeAcc + `) AND ` + patientIDColumn + ` IN (` + patientAcc + `)`
	res += `GROUP BY (` + timeConceptColumn + `)`
	return res

}

func stringMapAndAdd(inputList []string) string {
	outputList := make([]string, len(inputList))
	for i, str := range inputList {
		outputList[i] = `'` + str + `'`
	}
	return strings.Join(outputList, `,`)

}

// AggKSResults holds the return values of the distributed aggregation and key switch processes
type AggKSResults struct {
	Results libunlynx.CipherText
	Times   servicesmedco.TimeResults
	Err     error
}
type individualRequest struct {
	timeCode   *string
	cipherText *libunlynx.CipherText
	err        error
}

func aggregateAndKeySwitchSend(queryName, timeCode, eventValue, censoringValue, targetPubKey string, aggKsResultsChan chan struct {
	event     *unlynxResult
	censoring *unlynxResult
}) (err error) {
	//TODO verify that one unlynx Client is in deed serial
	/*
		unlynxClient, cothorityRoster := unlynx.NewUnlynxClient()
	*/
	// deserialize value and target public key
	eventValueDeserialized := libunlynx.CipherText{}
	err = eventValueDeserialized.Deserialize(eventValue)
	err = NiceError(err)
	if err != nil {
		return
	}
	censoringValueDeserialized := libunlynx.CipherText{}
	err = censoringValueDeserialized.Deserialize(censoringValue)
	err = NiceError(err)
	if err != nil {
		return
	}

	desTargetKey, err := libunlynx.DeserializePoint(targetPubKey)
	err = NiceError(err)
	if err != nil {
		logrus.Error("unlynx error deserializing target public key: ", err)
		return
	}

	// execute shuffle and key switching request
	type AggKSResults struct {
		Results libunlynx.CipherText
		Times   servicesmedco.TimeResults
		Err     error
	}
	type individualRequest struct {
		timeCode   *string
		cipherText *libunlynx.CipherText
		err        error
	}
	eventChannel := make(chan *individualRequest)
	censoringChannel := make(chan *individualRequest)
	//connectCallback := func() {}

	connectCallback := func() {
		var eventRequest *individualRequest
		var censoringRequest *individualRequest

		eventRequest = <-eventChannel
		censoringRequest = <-censoringChannel

		cipherEvent, err1 := eventRequest.cipherText.Serialize()
		cipherCensoring, err2 := censoringRequest.cipherText.Serialize()
		//both keys are extpected to be the same here
		eventResult := &unlynxResult{eventRequest.timeCode, cipherEvent, err1}
		censoringResult := &unlynxResult{censoringRequest.timeCode, cipherCensoring, err2}

		//TODO define this struct
		aggKsResultsChan <- struct {
			event     *unlynxResult
			censoring *unlynxResult
		}{eventResult, censoringResult}
	}

	individualAggSend := func(desValue libunlynx.CipherText, requestName string, timeCode string, individualUnlynxChannel chan *individualRequest) {

		unlynxClient, cothorityRoster := unlynx.NewUnlynxClient()
		surveyID, aggKsResult, _, aggKsErr := unlynxClient.SendSurveyAggRequest(
			cothorityRoster,
			servicesmedco.SurveyID(requestName+timeCode),
			desTargetKey,
			desValue,
			false,
		)
		logrus.Debug("Received results from survey " + string(*surveyID))
		aggKsErr = NiceError(aggKsErr)

		if aggKsErr != nil {
			return
		}

		result := &individualRequest{
			timeCode:   &timeCode,
			cipherText: &aggKsResult,
			err:        aggKsErr,
		}

		individualUnlynxChannel <- result

	}

	go individualAggSend(eventValueDeserialized, queryName+"_Event_", timeCode, eventChannel)
	logrus.Debug("Sent survey " + queryName + "_Event_" + timeCode)
	go individualAggSend(censoringValueDeserialized, queryName+"_Censoring_Event_", timeCode, censoringChannel)
	logrus.Debug("Sent survey " + queryName + "_Censoring_Event_" + timeCode)
	go connectCallback()

	return
}

func aggregateAndKeySwitchCollect(finalResultMap *sync.Map, aggKsResultsChan chan struct {
	event     *unlynxResult
	censoring *unlynxResult
}, waitGroup *sync.WaitGroup) (err error) {
	defer waitGroup.Done()
	select {
	case aggKsResult := <-aggKsResultsChan:
		if aggKsResult.event.Error != nil || aggKsResult.censoring.Error != nil {
			err = fmt.Errorf("Error during unlynx process on time points  : %s , %s", aggKsResult.event.Error.Error(), aggKsResult.censoring.Error.Error())

			return
		}

		if *aggKsResult.event.Key != *aggKsResult.censoring.Key {
			err = fmt.Errorf("concept key %s and %s are not the same", *aggKsResult.event.Key, *aggKsResult.censoring.Key)

			return
		}
		res := Result{
			EventValue:     aggKsResult.event.Value,
			CensoringValue: aggKsResult.censoring.Value,
			Error:          err,
		}
		//TODO map result ky tio string time point integer indentifier
		logrus.Debug("Received key switched events for time point" + *aggKsResult.event.Key)
		finalResultMap.Store(*aggKsResult.event.Key, res)

		return
	case <-time.After(time.Duration(utilserver.UnlynxTimeoutSeconds) * time.Second):
		err = errors.New("unlynx timeout")
		err = NiceError(err)

		return
	}

}
