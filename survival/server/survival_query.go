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
	GroupIDs      []string
	TimeCodes     []EncryptedEncID
	Result        struct {
		Timers    map[string]time.Duration
		EncEvents EventGroups
	}
}

// NewQuery query constructor
func NewQuery(qID, pubKey, patientSetID string, groupIDs []string, timeCodes []string) *Query {
	logrus.Debugf("group IDs %v:", groupIDs)
	encryptedEncIDs := make([]EncryptedEncID, len(timeCodes))
	for idx, timeCode := range timeCodes {
		encryptedEncIDs[idx] = EncryptedEncID(timeCode)
	}
	res := &Query{ID: qID, UserPublicKey: pubKey,
		PatientSetID: patientSetID,
		GroupIDs:     groupIDs,
		TimeCodes:    encryptedEncIDs}
	res.Result.EncEvents = EventGroups{}
	if len(groupIDs) == 0 {
		res.Result.EncEvents = append(res.Result.EncEvents, &EventGroup{GroupID: patientSetID, TimePointResults: make([]*TimePointResult, 0)})
	} else {
		for _, group := range groupIDs {
			res.Result.EncEvents = append(res.Result.EncEvents, &EventGroup{GroupID: group, TimePointResults: make([]*TimePointResult, 0)})
		}
	}
	res.Result.Timers = make(map[string]time.Duration)
	return res
}

func (q *Query) addTimer(label string, since time.Time) (err error) {
	if _, exists := q.Result.Timers[label]; exists {
		err = errors.New("label " + label + " already exists in timers")
		return
	}
	q.Result.Timers[label] = time.Since(since)
	return
}
func (q *Query) addTimers(timers map[string]time.Duration) (err error) {
	for label, duration := range timers {
		if _, exists := q.Result.Timers[label]; exists {
			err = errors.New("label " + label + " already exists in timers")
			return
		}
		q.Result.Timers[label] = duration
	}
	return
}

// PrintTimers print the execution time information if the log level is Debug
func (q *Query) PrintTimers() {
	for label, duration := range q.Result.Timers {
		logrus.Debug(label + duration.String())
	}
}

// Execute is the function that translates a medco survival query in calls on psql and unlynx
func (q *Query) Execute() (err error) {
	err = directAccessDB.Ping()

	//todo patientNumType
	patientSet := make([][]PatientNum, 0)

	errChan := make(chan error, 1)
	resChans := make(map[string]chan []PatientNum)
	var barrier sync.WaitGroup
	timeToGetPatientSets := time.Now()

	if len(q.GroupIDs) == 0 {
		barrier.Add(1)
		patientSet = append(patientSet, make([]PatientNum, 0))
		resChans[q.PatientSetID] = make(chan []PatientNum, 1)
		go getPatientAsync(q.PatientSetID, resChans[q.PatientSetID], errChan, &barrier)
	} else {
		barrier.Add(len(q.GroupIDs))

		for _, group := range q.GroupIDs {
			patientSet = append(patientSet, make([]PatientNum, 0))
			resChans[group] = make(chan []PatientNum, 1)
			//TODO maybe not necessary
			localGroup := group

			go getPatientAsync(localGroup, resChans[group], errChan, &barrier)

		}
	}
	barrierPtr := &barrier
	//eventGroupsChan:= make(chan *EventGroups)

	timeCodeMapChan, timerChan, errorChan := NewTimeCodesMapWithCallback(q.ID, q.TimeCodes, func(timeCodesMap *TimeCodesMap, innerErrChan chan error, innerTimerChan chan map[string]time.Duration) {
		var patientGroups [4][]PatientNum
		var times map[string]time.Duration

		barrierPtr.Wait()

		q.addTimer("time to get patient sets", timeToGetPatientSets)
		select {
		case err := <-errChan:
			logrus.Errorf("error int get patients")
			innerErrChan <- err
			return
		default:

		}
		if len(q.GroupIDs) == 0 {
			//TODO select case
			patientGroups[0] = <-resChans[q.PatientSetID]
			logrus.Debugf("not a pug %v:", patientGroups[0])
			patientGroups[1] = []PatientNum{}
			patientGroups[2] = []PatientNum{}
			patientGroups[3] = []PatientNum{}

			q.GroupIDs = append(q.GroupIDs, q.PatientSetID)
		} else {
			pos := 0

			for i, groupID := range q.GroupIDs {

				if patientGroupChan, ok := resChans[groupID]; ok {
					logrus.Debugf("golden retriever %v:", groupID)
					select {
					case patientGroups[i] = <-patientGroupChan:
						logrus.Debugf("pug %v:", patientGroups[i])
					case <-time.After(20 * time.Second):
						logrus.Panic("Unexpected delay")
					}
				} else {
					logrus.Debug("element not in maps")
				}
				pos++

			}
			for i := pos + 1; i < 4; i++ {
				patientGroups[i] = []PatientNum{}
			}

		}
		if len(timeCodesMap.tagIDs) == 0 {
			innerErrChan <- errors.New("tag id list is empty for time codes")
			return
		}
		logrus.Debugf("doggo %v:", patientGroups)

		queryString := buildGroupedQuery(patientGroups, timeCodesMap.tagIDs)
		logrus.Debugf("survival query string :%s", queryString)

		sqlTime := time.Now()
		rows, err := directAccessDB.Query(queryString)
		q.addTimer("time to execute sql query", sqlTime)

		if err != nil {

			innerErrChan <- err
			return
		}

		for rows.Next() {
			var recTimePointTag TagID
			var recGroupNum int
			var recEventsOfInterestConcat string
			var recCensoringEventsConcat string
			rows.Scan(&recTimePointTag, &recGroupNum, &recEventsOfInterestConcat, &recCensoringEventsConcat)

			encTimeCode := timeCodesMap.tagIDsToEncTimeCodes[recTimePointTag]
			groupID := q.GroupIDs[recGroupNum-1]
			eventOfInterestLocalAgg, err := unlynx.LocallyAggregateValues(strings.Split(recEventsOfInterestConcat, ","))
			if err != nil {
				innerErrChan <- err
				return
			}
			censoringEventLocalAgg, err := unlynx.LocallyAggregateValues(strings.Split(recCensoringEventsConcat, ","))

			//not optimal but there are only four groups

			for _, group := range q.Result.EncEvents {
				if group.GroupID == groupID {
					group.TimePointResults = append(group.TimePointResults, &TimePointResult{
						TimePoint: string(encTimeCode),
						Result: Result{
							EventValueAgg:     eventOfInterestLocalAgg,
							CensoringValueAgg: censoringEventLocalAgg,
						},
					})

				}
			}

		}

		//this could be with a call back or a new go routine, but this is the last step of the query
		wholeAKStime := time.Now()
		q.Result.EncEvents, times, err = AKSgroups(q.ID, q.Result.EncEvents, q.UserPublicKey)

		if err != nil {
			innerErrChan <- err
			return
		}
		q.addTimer("time for the whole AKS", wholeAKStime)
		q.addTimers(times)
		return
	})

	select {
	case err = <-errorChan:
		return
	case <-timeCodeMapChan:
		select {
		case q.Result.Timers = <-timerChan:
		default:
		}
		return
	}
}

type unlynxResult struct {
	Key   *EncryptedEncID
	Value string
	Error error
}

func patientNumToString(pnums []PatientNum) []string {
	res := make([]string, len(pnums))
	for i, pnum := range pnums {
		res[i] = string(pnum)
	}
	return res
}

func buildGroupedQuery(patients [4][]PatientNum, timeCodes []TagID) string {

	var groups = make([]string, 4)
	groups[0] = "'{" + strings.Join(patientNumToString(patients[0]), ",") + "}'::int[]"
	groups[1] = "'{" + strings.Join(patientNumToString(patients[1]), ",") + "}'::int[]"
	groups[2] = "'{" + strings.Join(patientNumToString(patients[2]), ",") + "}'::int[]"
	groups[3] = "'{" + strings.Join(patientNumToString(patients[3]), ",") + "}'::int[]"
	groupString := strings.Join(groups, ", ")

	var timeCodesString []string
	for _, timeCode := range timeCodes {
		timeCodesString = append(timeCodesString, string(timeCode))
	}
	timeAcc := stringMapAndAdd(timeCodesString)

	res := "SELECT * FROM (SELECT " + timeConceptColumn + ", grouping_filter(" + patientNumColumn + ", " + groupString + "), STRING_AGG(" + eventOfInterestColumn + ",','), STRING_AGG(" + censoringEventColumn + ",',')\n"
	res += "FROM survivaldemodata_i2b2.survival_fact\n"
	res += "WHERE " + timeConceptColumn + " IN (" + timeAcc + ") \n"

	res += "GROUP BY " + timeConceptColumn + ", grouping_filter) AS whole_set\n"
	res += "WHERE grouping_filter  != -1\n"
	return res

}

func stringMapAndAdd(inputList []string) string {
	outputList := make([]string, len(inputList))
	for i, str := range inputList {
		outputList[i] = `'` + str + `'`
	}
	return strings.Join(outputList, `,`)

}

func getPatientAsync(patientSetId string, resChannel chan []PatientNum, errChan chan error, barrier *sync.WaitGroup) {
	patients, _, err := i2b2.GetPatientSet(patientSetId)
	defer func() {
		barrier.Done()
	}()
	if err != nil {
		select {
		case errChan <- err:
			//push error
		default:
			//error is not pushed to the unavailable channel as a previous error occurs

		}
		return
	}
	tmpRes := make([]PatientNum, len(patients))
	for i, patient := range patients {
		tmpRes[i] = PatientNum(patient)
	}
	resChannel <- tmpRes

	return

}
