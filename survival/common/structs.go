package common

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ldsec/medco-connector/restapi/models"
	"github.com/ldsec/medco-connector/wrappers/i2b2"
	"github.com/ldsec/medco-connector/wrappers/unlynx"
	"github.com/sirupsen/logrus"
	"go.dedis.ch/onet/v3/log"
)

type Result struct {
	EventValue     string //or kyber.point
	CensoringValue string //or kyber.point
	Delay          time.Duration
	Error          error
}

const eventFlagSeparator = ` `

//type Map sync.Map

type PatientID string

func (str PatientID) String() string {
	return string(str)
}

var ResultMap sync.Map

type BatchIterator struct {
	timeCodes         []string
	length            int
	batchNumber       int
	batchSize         float64
	currentBatchIndex int
	//currentStateLower int
	//currentStateUpper int
	endReached bool
}

func NewBatchItertor(timePoints []string, batchNumber int) (batches *BatchIterator, err error) {
	length := len(timePoints)
	if length == 0 {
		err = errors.New("Input array must contain at least 1 time code")
		return
	}
	if batchNumber < 1 {
		err = errors.New("Number of batch should be at least 1")
		return
	}

	if batchNumber > length {
		logrus.Info(fmt.Sprintf("Batch number %d higher than lenght of time points array %d. Changing batch number to %d to avoid empty batch", batchNumber, length, length))
		batchNumber = length
	}

	batches = &BatchIterator{
		timeCodes:   timePoints,
		length:      length,
		batchNumber: batchNumber,
		batchSize:   float64(length) / float64(batchNumber),
	}
	return
}

func (batches *BatchIterator) Next() (res []string) {
	resLower := int(math.Floor(float64(batches.currentBatchIndex) * batches.batchSize))
	resUpper := int(math.Floor(float64(batches.currentBatchIndex+1) * batches.batchSize))
	res = batches.timeCodes[resLower:resUpper]
	if batches.currentBatchIndex < batches.batchNumber-1 {
		batches.currentBatchIndex++
	} else if batches.endReached == false {
		batches.endReached = true
	}

	return res
}

func (batches *BatchIterator) Done() bool {
	return batches.endReached
}

func BreakBlob(blobValue string) (eventOfInterest, censoringEvent string, err error) {
	res := strings.Split(blobValue, eventFlagSeparator)
	//TODO magic
	if len(res) != 2 {
		err = fmt.Errorf(`Blob value %s is ill-formed, it should be  to base64 encoded sequence separated by "%s" (without quotes)`, blobValue, eventFlagSeparator)
		return
	}
	eventOfInterest = res[0]
	censoringEvent = res[1]
	return
}

const minusOne32 = int32(-1)

//this is an ad-hoc barrier mechanism wrapping a waitgroup
type Barrier struct {
	condition int32
	value     int32
	waitGroup sync.WaitGroup
}

func NewBarrier(condition int) (barrier *Barrier, err error) {
	if condition < 1 {
		err = fmt.Errorf("The condition number must be at least 1, here %d", condition)
		return
	}
	if condition > math.MaxInt32 {
		err = fmt.Errorf("%d exceeds max int32", condition)
		return
	}

	barrier = &Barrier{
		condition: int32(condition),
	}
	return

}

func (barrier *Barrier) Add(delta int32) {
	barrier.waitGroup.Add(int(delta))
	atomic.AddInt32(&(barrier.value), delta)
	//check for overflow ??

}

func (barrier *Barrier) Done() {
	barrier.waitGroup.Done()
	atomic.AddInt32(&(barrier.value), minusOne32)
	//check for overflow or negative value ?
}

func (barrier *Barrier) ConditionalWait() {
	var conditionExceeded bool
	var old int32
	for consistent := false; !consistent; {
		old = atomic.LoadInt32(&(barrier.value))
		conditionExceeded = old >= barrier.condition
		//looks like a CAS without swap, maybe a mutex instead is more appropriate ??
		consistent = atomic.LoadInt32(&(barrier.value)) == old
	}

	if conditionExceeded {
		barrier.waitGroup.Wait()

		for consistent := false; !consistent; {
			old = atomic.LoadInt32(&(barrier.value))
			//this looks more appropriate here, but still a mutex is simpleer to use ?
			consistent = atomic.CompareAndSwapInt32(&(barrier.value), old, int32(0))
		}
	} //don't wait otherwise !!
}

func (barrier *Barrier) AbsoluteWait() {
	barrier.waitGroup.Wait()
	for consistent := false; !consistent; {
		old := atomic.LoadInt32(&(barrier.value))
		//this looks more appropriate here, but still a mutex is simpleer to use ?
		consistent = atomic.CompareAndSwapInt32(&(barrier.value), old, int32(0))
	}
}

type Set struct {
	data map[string]struct{}
}

func NewSet(size int) *Set {
	return &Set{data: make(map[string]struct{}, size)}
}

func (set *Set) Add(key string) {
	set.data[key] = struct{}{}
}

func (set *Set) Remove(key string) {
	_, ok := set.data[key]
	if ok {
		delete(set.data, key)
	}
}

func (set *Set) ForEach(instruction func(string)) {
	for key := range set.data {
		instruction(key)
	}

}

const timeConceptRootPath = `/SurvivalAnalysis/`

//for the entire survival query, maybe wrap this around a query anin a structure that is created for each new survival query..

var errorChannel = make(chan error)
var endOfProcessChannel = make(chan bool)

func PushError(err error) {
	select {
	case errorChannel <- err:
		log.Lvl2("pushed an error")
	default:
		log.Lvl2("error channels already full")
	}
}
func Finished() {
	endOfProcessChannel <- true
}

func WaitEndSignal(timeoutInSeconds int) (err error) {
	select {
	case err = <-errorChannel:
		return
	case <-endOfProcessChannel:
		return
	case <-time.After(time.Duration(timeoutInSeconds) * time.Second):
		err = errors.New("Survival query Timeout")
		return
	}

}

//TODO another function already exists in unlynx wrapper
func ZeroPointEncryption() (res string, err error) {

	events, err := unlynx.EncryptWithCothorityKey(int64(0))
	if err != nil {
		return
	}
	censoringEvents, err := unlynx.EncryptWithCothorityKey(int64(0))
	if err != nil {
		return
	}
	res = events + eventFlagSeparator + censoringEvents
	return
}

func UnlynxRequestName(queryName, timecode string) string {
	return queryName + timecode
}

//this redundance is a quick solutiio to avoid cyclic imports with medcoservers
type ExploreQuery struct {
	ID     string
	Query  *models.ExploreQuery
	Result struct {
		EncCount       string
		EncPatientList []string
		Timers         map[string]time.Duration
		EncEvents      map[string][2]string
	}
	ExecuteCallback func(*ExploreQuery, []string, []string, int) error
}

var ExecCallback func(*ExploreQuery, []string, []string, int) error

func (q *ExploreQuery) Execute(patientIDs []string) error {
	timeCodes, err := GetTimeCodes()
	if err != nil {
		return err
	}

	return q.ExecuteCallback(q, patientIDs, timeCodes, 1)
}

func GetTimeCodes() (timeCode []string, err error) {
	results, err := i2b2.GetOntologyChildren(timeConceptRootPath)
	if err != nil {
		return
	}
	for _, result := range results {
		timeCode = append(timeCode, strconv.FormatInt(*(result.MedcoEncryption.ID), 10))
	}
	return

}
