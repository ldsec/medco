package directaccess

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	survivalserver "github.com/ldsec/medco-connector/survival/server"
	utilserver "github.com/ldsec/medco-connector/util/server"
	"github.com/ldsec/medco-connector/wrappers/unlynx"
	servicesmedco "github.com/ldsec/medco-unlynx/services"
	libunlynx "github.com/ldsec/unlynx/lib"
	"github.com/sirupsen/logrus"
)

const (
	schema            = "i2b2demodata_i2b2"
	table             = "observation_fact"
	blobColumn        = "observation_blob"
	timeConceptColumn = "concept_cd"
	patientIDColumn   = "patient_num"
	naiveVersion      = true
	interBlob         = "," //blobconcept does not contain any comma
)

//type Map survivalserver.Map

//type PatientSet []PatientID

// QueryTimePoints is the function that translates a medco survival query in multiple calls on psql and unlynx
type SurvivalQuery interface {
	Execute() error
	GetTimeCodes() ([]string, error)
	GetPatients() ([]string, error)
	GetID() string
	GetUserPublicKey() string
	SetResultMap(map[string][2]string) error
}

func QueryTimePoints(q SurvivalQuery, batchNumber int) (err error) {
	var timePoints []string
	var patientSet []string

	timePoints, err = q.GetTimeCodes()
	if err != nil {
		return
	}
	patientSet, err = q.GetPatients()
	if err != nil {
		return
	}

	exceptionHandler, err := survivalserver.NewExceptionHandler(1)
	if err != nil {
		return
	}
	go func() {
		err = utilserver.DBConnection.Ping()
		if err != nil {
			logrus.Error("Impossible to connect to database")
			exceptionHandler.PushError(err)
			return
		}
		batches, err := survivalserver.NewBatchItertor(timePoints, batchNumber)
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
				waitGroup.Add(batchNumber)
				var events []string
				var censorings []string
				for _, time := range batch {
					//error chans ??
					go func(time string) {

						result := survivalserver.Result{}
						defer waitGroup.Done()
						query := buildSingleQuery(patientSet, time)
						rows, err := utilserver.DBConnection.Query(query)
						if err != nil {
							result.Error = err
							survivalserver.ResultMap.Store(time, result)
							exceptionHandler.PushError(err)
							return
						}

						for rows.Next() { //there should be at most one row here
							var blob string
							err = rows.Scan(&blob)
							if err != nil {
								exceptionHandler.PushError(err)
								return
							}
							event, censoring, err := survivalserver.BreakBlob(blob)
							if err != nil {
								exceptionHandler.PushError(err)
								return
							}
							events = append(events, event)
							censorings = append(censorings, censoring)
						}
						err = rows.Close()
						if err != nil {
							exceptionHandler.PushError(err)
							return
						}

						events, err := unlynx.LocallyAggregateValues(events)
						if err != nil {
							exceptionHandler.PushError(err)
							return
						}
						censoringEvents, err := unlynx.LocallyAggregateValues(censorings)
						if err != nil {
							exceptionHandler.PushError(err)

						}

						err = aggregateAndKeySwitchSend(q.GetID(), time, events, censoringEvents, q.GetUserPublicKey(), unlynxChannel)

					}(time)

				}
				//TODO repeat
				errOccured := false
				var receptionBarrier sync.WaitGroup
				receptionBarrier.Add(len(batch))
				for range batch {
					err = aggregateAndKeySwitchCollect(&survivalserver.ResultMap, unlynxChannel, &receptionBarrier)
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

					rows, err := utilserver.DBConnection.Query(query)

					if err != nil {
						result.Error = err
						//TODO handle this kind of error
						survivalserver.ResultMap.Store("", result)
						exceptionHandler.PushError(err)
						return
					}
					//var unlynxBarrier sync.WaitGroup
					//for the moment do all in one run (dangerous ?)
					//unlynxBarrier.Add(rows)

					unlynxBarrier, err := survivalserver.NewBarrier(100)
					if err != nil {
						result.Error = err
						//TODO handle this kind of error
						survivalserver.ResultMap.Store("", result)
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
						err = rows.Scan(recipiens)
						set.Remove(recipiens.TimeCode)

						str := strings.Split(recipiens.ConcatenatedBlobs, interBlob)
						eventsOfInterest := make([]string, len(str))
						censoringEvents := make([]string, len(str))
						for idx, blob := range str {
							eventsOfInterest[idx], censoringEvents[idx], err = survivalserver.BreakBlob(blob)
							if err != nil {
								result.Error = err
								survivalserver.ResultMap.Store(recipiens.TimeCode, result)
								rows.Close()
								exceptionHandler.PushError(err)
								return
							}

						}

						unlynxBarrier.Add(1)
						//err chans ?
						go func() { //call unlynx for the first kind of events

							events, err := unlynx.LocallyAggregateValues(eventsOfInterest)
							if err != nil {
								exceptionHandler.PushError(err)
								return
							}
							censoringEvents, err := unlynx.LocallyAggregateValues(censoringEvents)
							if err != nil {
								exceptionHandler.PushError(err)
								return

							}

							err = aggregateAndKeySwitchSend(`queryname`, recipiens.TimeCode, events, censoringEvents, q.GetUserPublicKey(), unlynxChannel)

							unlynxBarrier.Done()
						}()

						unlynxBarrier.ConditionalWait()
					}
					rows.Close()
					unlynxBarrier.AbsoluteWait()
					//for those that have not be found in this node

					//TODO this is beyound the unlynx batch barrier mechanism
					set.ForEach(func(key string) {

						go func() {
							zeroEvent, err := survivalserver.ZeroPointEncryption()
							if err != nil {
								exceptionHandler.PushError(err)

							}
							zeroCensoring, err := survivalserver.ZeroPointEncryption()
							err = aggregateAndKeySwitchSend(`queryname`, key, zeroEvent, zeroCensoring, q.GetUserPublicKey(), unlynxChannel)
							if err != nil {
								exceptionHandler.PushError(err)
							}

						}()
						//survivalserver.ResultMap.Store(key, result)

					})

					// ------- collect the result of the  aggregate and key switch
					errorOccurred := false
					var receptionBarrier sync.WaitGroup
					receptionBarrier.Add(len(batch))
					for range batch {
						err = aggregateAndKeySwitchCollect(&survivalserver.ResultMap, unlynxChannel, &receptionBarrier)
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
	targetMap := make(map[string][2]string)
	survivalserver.ResultMap.Range(func(timeCode /*string*/, events interface{} /*survivalserver.Result*/) bool {
		timeCodeString := timeCode.(string)
		eventStruct := events.(survivalserver.Result)
		targetMap[timeCodeString] = [2]string{eventStruct.EventValue, eventStruct.CensoringValue}
		return true
	})
	err = q.SetResultMap(targetMap)
	if err != nil {
		return
	}
	return
}

type unlynxResult struct {
	Key   string
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

func aggregateAndKeySwitchSend(queryName, key, eventValue, censoringValue, targetPubKey string, aggKsResultsChan chan struct {
	event     *unlynxResult
	censoring *unlynxResult
}) (err error) {
	unlynxClient, cothorityRoster := unlynx.NewUnlynxClient()

	// deserialize value and target public key
	eventValueDeserialized := libunlynx.CipherText{}
	err = eventValueDeserialized.Deserialize(eventValue)
	if err != nil {
		return
	}
	censoringValueDeserialized := libunlynx.CipherText{}
	err = censoringValueDeserialized.Deserialize(censoringValue)
	if err != nil {
		return
	}

	desTargetKey, err := libunlynx.DeserializePoint(targetPubKey)
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
		key        string
		cipherText *libunlynx.CipherText
		err        error
	}
	eventChannel := make(chan *individualRequest, 1)
	censoringChannel := make(chan *individualRequest, 1)
	//connectCallback := func() {}

	connectCallback := func() {
		var eventRequest *individualRequest
		var censoringRequest *individualRequest
		select {
		//if event is ready first, does not handle timeout here
		case eventRequest = <-eventChannel:
			censoringRequest = <-censoringChannel
		case censoringRequest = <-censoringChannel:
			eventRequest = <-eventChannel
		}
		cipherEvent, err1 := eventRequest.cipherText.Serialize()
		cipherCensoring, err2 := censoringRequest.cipherText.Serialize()
		//both keys are extpected to be the same here
		eventResult := &unlynxResult{eventRequest.key, cipherEvent, err1}
		censoringResult := &unlynxResult{censoringRequest.key, cipherCensoring, err2}

		select {
		//TODO define this struct
		case aggKsResultsChan <- struct {
			event     *unlynxResult
			censoring *unlynxResult
		}{eventResult, censoringResult}:
		//success in merging chans
		case <-time.After(time.Duration(300) * time.Second): //TODO magic
			logrus.Error("Buffered channel full for too long")
		}
		//TODO define this struct
		aggKsResultsChan <- struct {
			event     *unlynxResult
			censoring *unlynxResult
		}{eventResult, censoringResult}
	}

	individualAggSend := func(desValue libunlynx.CipherText, individualUnlynxChannel chan *individualRequest, requestID string) {
		_, aggKsResult, _, aggKsErr := unlynxClient.SendSurveyAggRequest(
			cothorityRoster,
			servicesmedco.SurveyID(queryName+requestID+"_AGG"),
			desTargetKey,
			desValue,
			false,
		)

		if aggKsErr != nil {
			return
		}
		//TODO define this nested type instead of repeating anonymous struct
		result := &individualRequest{
			key:        key,
			cipherText: &aggKsResult,
			err:        aggKsErr,
		}
		select {
		case individualUnlynxChannel <- result:
			//success
			return
		case <-time.After(time.Duration(300) * time.Second): //TODO magic
			err = errors.New("Unexpected error on local channel")
		}
	}

	go individualAggSend(eventValueDeserialized, eventChannel, key+"event")
	go individualAggSend(censoringValueDeserialized, censoringChannel, key+"censoring")
	go connectCallback()

	return
}

func aggregateAndKeySwitchCollect(finalResultMap *sync.Map, aggKsResultsChan chan struct {
	event     *unlynxResult
	censoring *unlynxResult
}, waitGroup *sync.WaitGroup) (err error) {
	select {
	case aggKsResult := <-aggKsResultsChan:
		if aggKsResult.event.Error != nil || aggKsResult.censoring.Error != nil {
			err = fmt.Errorf("Error during unlynx process on time points  : %s , %s", aggKsResult.event.Error.Error(), aggKsResult.censoring.Error.Error())
			waitGroup.Done()
			return
		}
		//times = aggKsResult.Times.MapTR
		if aggKsResult.event.Key != aggKsResult.censoring.Key {
			err = errors.New("time concept keys are  expected to be the same")
			waitGroup.Done()
			return
		}
		res := survivalserver.Result{
			EventValue:     aggKsResult.event.Value,
			CensoringValue: aggKsResult.censoring.Value,
			Error:          err,
		}
		finalResultMap.Store(aggKsResult.event.Key, res)
		waitGroup.Done()
		return
	case <-time.After(time.Duration(utilserver.UnlynxTimeoutSeconds) * time.Second):
		err = errors.New("unlynx timeout")
		waitGroup.Done()
		return
	}

}
