package directaccess

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ldsec/medco-connector/wrappers/i2b2"

	survivalserver "github.com/ldsec/medco-connector/survival/server"
	utilserver "github.com/ldsec/medco-connector/util/server"
	"github.com/ldsec/medco-connector/wrappers/unlynx"
	servicesmedco "github.com/ldsec/medco-unlynx/services"
	libunlynx "github.com/ldsec/unlynx/lib"
	"github.com/sirupsen/logrus"
)

const (
	schema            string = "i2b2demodata_i2b2"
	table             string = "observation_fact"
	blobColumn        string = "observation_blob"
	timeConceptColumn string = "concept_cd"
	patientIDColumn   string = "patient_num"
	naiveVersion      bool   = false
	interBlob         string = "," //blobconcept does not contain any comma
)

var DirectAccessDB *sql.DB
var DBName string

//TODO enable once connection debugged

func init() {

	host := os.Getenv("DIRECT_ACCESS_DB_HOST")
	port, err := strconv.ParseInt(os.Getenv("DIRECT_ACCESS_DB_PORT"), 10, 64)
	if err != nil || port < 0 || port > 65535 {
		logrus.Warn("Invalid port, defaulted")
		port = 5432
	}
	name := os.Getenv("DIRECT_ACCESS_DB_NAME")
	DBName = name
	loginUser := os.Getenv("DIRECT_ACCESS_DB_USER")
	loginPw := os.Getenv("DIRECT_ACCESS_DB_PW")

	DirectAccessDB, err = utilserver.InitializeConnectionToDB(host, int(port), name, loginUser, loginPw)
	if err != nil {
		logrus.Error("Unable to connect database for direct access to I2B2")
		return
	}
	logrus.Info("Connected I2B2 DB for direct access")
	return
}

//type Map survivalserver.Map

//type PatientSet []PatientID

// QueryTimePoints is the function that translates a medco survival query in multiple calls on psql and unlynx
type SurvivalQuery interface {
	Execute() error
	GetTimeCodes() []string
	GetPatientSetID() string
	GetID() string
	GetUserPublicKey() string
	SetResultMap(map[string][2]string) error
}

func QueryTimePoints(q SurvivalQuery, batchNumber int) (err error) {
	var storage sync.Map

	err = DirectAccessDB.Ping()

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
	//unlynx.DDTagValues()
	TagIDRetrievalMethod := DirectAccessTags(getTagIDs)

	timeCodesMap, _, err := survivalserver.NewTimeCodesMap(q.GetID(), encTimePoints, &TagIDRetrievalMethod)
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

	err = survivalserver.NiceError(err)
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
			zeroEvents, err = survivalserver.ZeroPointEncryption()
			if err != nil {
				return
			}
			zeroCensoring, err = survivalserver.ZeroPointEncryption()
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

		err = fillResult(q, &storage)

		return
	}

	if err != nil {
		return
	}

	exceptionHandler, err := survivalserver.NewExceptionHandler(1)
	err = survivalserver.NiceError(err)

	if err != nil {
		return
	}
	go func() {

		batches, err := survivalserver.NewBatchItertor(timePoints, batchNumber)

		err = survivalserver.NiceError(err)

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

						result := survivalserver.Result{}
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

						rows, err := DirectAccessDB.Query(query)

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

							err = survivalserver.NiceError(err)
							if err != nil {
								deferred(err)

								return
							}
							event, censoring, err := survivalserver.BreakBlob(blob)
							err = survivalserver.NiceError(err)
							if err != nil {
								deferred(err)
								return
							}
							events = append(events, event)
							censorings = append(censorings, censoring)
						}
						err = rows.Close()
						err = survivalserver.NiceError(err)
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
						err = survivalserver.NiceError(err)
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
						err = survivalserver.NiceError(err)
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
					result := survivalserver.Result{}

					defer waitGroup.Done()

					query := buildBatchQuery(patientSet, batch)

					rows, err := DirectAccessDB.Query(query)
					//debug
					if err != nil {
						err = errors.New(err.Error() + " " + query)
						return
					}
					err = survivalserver.NiceError(err)

					if err != nil {
						result.Error = err
						//TODO handle this kind of error
						storage.Store("", result)
						exceptionHandler.PushError(err)
						return
					}

					set := survivalserver.NewSet(len(batch))
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
						err = survivalserver.NiceError(err)
						set.Remove(recipiens.TimeCode)

						str := strings.Split(recipiens.ConcatenatedBlobs, interBlob)
						eventsOfInterest := make([]string, len(str))
						censoringEvents := make([]string, len(str))
						for idx, blob := range str {
							eventsOfInterest[idx], censoringEvents[idx], err = survivalserver.BreakBlob(blob)
							err = survivalserver.NiceError(err)
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
							err = survivalserver.NiceError(err)
							if err != nil {
								exceptionHandler.PushError(err)
								return
							}
							censoringEvents, err := unlynx.LocallyAggregateValues(censoringEvents)
							err = survivalserver.NiceError(err)
							if err != nil {
								exceptionHandler.PushError(err)
								return

							}
							encTimeCode, ok := tagIDToEncTimePoints[recipiens.TimeCode]
							if !ok {
								//TODO orthograph
								err = errors.New("ill formed time codes map")
							}
							err = survivalserver.NiceError(err)
							if err != nil {
								exceptionHandler.PushError(err)
								return
							}

							err = aggregateAndKeySwitchSend(q.GetID(), encTimeCode, events, censoringEvents, q.GetUserPublicKey(), unlynxChannel)
							err = survivalserver.NiceError(err)
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
							zeroEvent, err := survivalserver.ZeroPointEncryption()
							err = survivalserver.NiceError(err)
							if err != nil {
								exceptionHandler.PushError(err)

							}
							zeroCensoring, err := survivalserver.ZeroPointEncryption()
							err = aggregateAndKeySwitchSend(`queryname`, key, zeroEvent, zeroCensoring, q.GetUserPublicKey(), unlynxChannel)
							err = survivalserver.NiceError(err)
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
						err = survivalserver.NiceError(err)
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
	err = fillResult(q, &storage)
	err = survivalserver.NiceError(err)
	if err != nil {
		return
	}
	return
}
func fillResult(q SurvivalQuery, storage *sync.Map) (err error) {

	targetMap := make(map[string][2]string)
	counter := 0
	storage.Range(func(timeCode /*string*/, events interface{} /*survivalserver.Result*/) bool {
		if timeCodeString, ok1 := timeCode.(string); ok1 {
			if eventStruct, ok2 := events.(survivalserver.Result); ok2 {

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
	err = survivalserver.NiceError(err)
	if err != nil {
		return
	}
	censoringValueDeserialized := libunlynx.CipherText{}
	err = censoringValueDeserialized.Deserialize(censoringValue)
	err = survivalserver.NiceError(err)
	if err != nil {
		return
	}

	desTargetKey, err := libunlynx.DeserializePoint(targetPubKey)
	err = survivalserver.NiceError(err)
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
		aggKsErr = survivalserver.NiceError(aggKsErr)

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
		res := survivalserver.Result{
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
		err = survivalserver.NiceError(err)

		return
	}

}
