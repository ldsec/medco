package survivalserver

import (
	"errors"
	"fmt"
	"sync"
	"time"

	utilserver "github.com/ldsec/medco-connector/util/server"
	"github.com/ldsec/medco-connector/wrappers/unlynx"
	servicesmedco "github.com/ldsec/medco-unlynx/services"
	libunlynx "github.com/ldsec/unlynx/lib"
	"github.com/sirupsen/logrus"
)

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
