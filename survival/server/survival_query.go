package survivalserver

import (
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ldsec/medco-connector/wrappers/i2b2"
	"github.com/ldsec/medco-connector/wrappers/unlynx"

	"github.com/sirupsen/logrus"
)

// Query holds the ID of the survival analysis, its parameters and a pointer to its results
type Query struct {
	ID            string
	UserPublicKey string
	PatientSetID  string
	TimeCodes     []EncryptedEncID

	//TODO also hide that
	Result *struct {
		Timers    map[string]time.Duration
		EncEvents map[string][2]string
	}
	timers map[string]time.Duration

	spin *Spin
}

// NewQuery query constructor
func NewQuery(qID, pubKey, patientSetID string, timeCodes []string) *Query {
	encryptedEncIDs := make([]EncryptedEncID, len(timeCodes))
	for idx, timeCode := range timeCodes {
		encryptedEncIDs[idx] = EncryptedEncID(timeCode)
	}
	return &Query{ID: qID, UserPublicKey: pubKey, PatientSetID: patientSetID, TimeCodes: encryptedEncIDs, timers: make(map[string]time.Duration), spin: NewSpin()}
}

func (q *Query) addTimer(label string, since time.Time) (err error) {
	if _, exists := q.timers[label]; exists {
		err = errors.New("label " + label + " already exists in timers")
		return
	}
	q.timers[label] = time.Since(since)
	return
}
func (q *Query) addTimers(timers map[string]time.Duration) (err error) {
	for label, duration := range timers {
		if _, exists := q.timers[label]; exists {
			err = errors.New("label " + label + " already exists in timers")
			return
		}
		q.timers[label] = duration
	}
	return
}

// GetTimers returns the execution time information
func (q *Query) GetTimers() map[string]time.Duration {
	return q.timers
}

// PrintTimers print the execution time information if the log level is Debug
func (q *Query) PrintTimers() {
	for label, duration := range q.timers {
		logrus.Debug(label + duration.String())
	}
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

	var timePoints []TagID
	var patientSet []string

	encTimePoints := q.TimeCodes
	if len(encTimePoints) == 0 {
		logrus.Panic("Unexpected empty list of time points")
	}

	timeCodesMap, _, err := NewTimeCodesMap(q.ID, encTimePoints)
	if err != nil {
		return
	}
	timePoints = timeCodesMap.tagIDs
	tagIDToEncTimePoints := timeCodesMap.tagIDsToEncTimeCodes

	patientSetID := q.PatientSetID
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

		for _, encryptedEncIDs := range tagIDToEncTimePoints {
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

			err = aggregateAndKeySwitchSend(q.ID, encryptedEncIDs, zeroEvents, zeroCensoring, q.UserPublicKey, unlynxChannel)
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
		start := time.Now()
		defer func(err error) {
			err = q.addTimer("Total time for a non-empty set", start)
		}(err)

		batches, err := NewBatchIterator(timePoints, batchNumber)

		err = NiceError(err)

		if err != nil {
			exceptionHandler.PushError(err)
			return
		}

		for batchCounter := 0; !batches.Done(); batchCounter++ {
			batch := batches.Next()
			//TODO magic numbers
			unlynxChannel := make(chan struct {
				event     *unlynxResult
				censoring *unlynxResult
			}, 1024)

			batchIdx := batchCounter

			//TODO not elegant !!!!!!!
			result := Result{}

			query := buildBatchQuery(patientSet, batch)

			SQLQueryTime := time.Now()
			rows, err := directAccessDB.Query(query)
			if err != nil {
				err = errors.New(err.Error() + " " + query)
				return
			}
			err = q.addTimer("SQLQueryTime of batch "+strconv.Itoa(batchIdx)+" ", SQLQueryTime)
			if err != nil {
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
			unlynxMachineryTime := time.Now()
			for rows.Next() {
				recipiens := &struct {
					TimeCode          string
					ConcatenatedBlobs string
				}{}
				err = rows.Scan(&(recipiens.TimeCode), &(recipiens.ConcatenatedBlobs))
				err = NiceError(err)
				set.Remove(TagID(recipiens.TimeCode))

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

					err = NiceError(err)
					if err != nil {
						exceptionHandler.PushError(err)
						return
					}

					if EID, ok := timeCodesMap.tagIDsToEncTimeCodes[TagID(recipiens.TimeCode)]; ok {
						err = aggregateAndKeySwitchSend(q.ID, EID, events, censoringEvents, q.UserPublicKey, unlynxChannel)
					} else {
						err = errors.New("missing TAGID in the time codes maps")
					}

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
			set.ForEach(func(key TagID) {

				go func() {
					zeroEvent, err := zeroPointEncryption()
					err = NiceError(err)
					if err != nil {
						exceptionHandler.PushError(err)
					}
					zeroCensoring, err := zeroPointEncryption()
					if EID, ok := timeCodesMap.tagIDsToEncTimeCodes[key]; ok {
						err = aggregateAndKeySwitchSend(`queryname`, EID, zeroEvent, zeroCensoring, q.UserPublicKey, unlynxChannel)
					} else {
						err = errors.New("missing key in the time codes maps")
					}
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

			err = q.addTimer("unlynx machinery time for batch "+strconv.Itoa(batchIdx)+" ", unlynxMachineryTime)
			if err != nil {
				return
			}

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
		if timeCodeString, ok1 := timeCode.(EncryptedEncID); ok1 {
			if eventStruct, ok2 := events.(Result); ok2 {

				targetMap[string(timeCodeString)] = [2]string{eventStruct.EventValue, eventStruct.CensoringValue}
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
	Key   *EncryptedEncID
	Value string
	Error error
}

func buildSingleQuery(patients []string, timeCode string) string {

	patientAcc := stringMapAndAdd(patients)
	res := `SELECT ` + blobColumn + ` FROM ` + schema + `.` + table + ` `
	res += `WHERE ` + timeConceptColumn + ` = '` + timeCode + `' AND ` + patientIDColumn + ` IN (` + patientAcc + `)`
	return res

}
func buildBatchQuery(patients []string, timeCodes []TagID) string {

	patientAcc := stringMapAndAdd(patients)
	timeCodesString := make([]string, len(timeCodes))
	for idx, timeCode := range timeCodes {
		timeCodesString[idx] = string(timeCode)
	}
	timeAcc := stringMapAndAdd(timeCodesString)
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
