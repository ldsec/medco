package survivalserver

import (
	"errors"
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
