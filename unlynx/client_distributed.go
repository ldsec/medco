package unlynx

import (
	"errors"
	"github.com/lca1/medco-connector/util"
	"github.com/lca1/medco-unlynx/services"
	"github.com/lca1/unlynx/lib"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

// todo: get time measurements in all functions
// 	tr.DDTRequestTimeCommunication = totalTime - tr.DDTRequestTimeExec
//	tr.DDTParsingTime = parsingTime
//	tr.DDTRequestTimeExec += tr.DDTParsingTime
//  totalTime = around request

// DDTagValues makes request through unlynx to compute distributed deterministic tags of encrypted values
func DDTagValues(queryName string, values []string) (taggedValues map[string]string, err error) {
	unlynxClient, cothorityRoster := newUnlynxClient()

	// deserialize values
	desValues, err := deserializeCipherVector(values)
	if err != nil {
		return
	}

	// execute DDT
	ddtResultsChan := make(chan []libunlynx.GroupingKey)
	ddtErrChan := make(chan error)
	go func() {
		_, ddtResults, ddtErr := unlynxClient.SendSurveyDDTRequestTerms(
			cothorityRoster,
			servicesmedco.SurveyID(queryName + "_DDT"),
			desValues,
			false,
			false,
		)
		if ddtErr != nil {
			ddtErrChan <- ddtErr
		} else if len(ddtResults) == 0 || len(ddtResults) != len(values) {
			ddtErrChan <- errors.New("unlynx inconsistent DDT results: #results=" + strconv.Itoa(len(ddtResults)) + ", #terms=" + strconv.Itoa(len(values)))
		} else {
			ddtResultsChan <- ddtResults
		}
	}()

	select {
	case ddtResults := <-ddtResultsChan:
		taggedValues = make(map[string]string)
		for i, result := range ddtResults {
			taggedValues[values[i]] = string(result)
		}

	case err = <-ddtErrChan:
		logrus.Error("unlynx error executing DDT: ", err)

	case <-time.After(time.Duration(util.UnlynxTimeoutSeconds) * time.Second):
		err = errors.New("unlynx timeout")
		logrus.Error(err)
	}
	return
}

// KeySwitchValues makes request through unlynx to key switch a single encrypted value (convenience function)
func KeySwitchValue(queryName string, value string, targetPubKey string) (string, error) {
	results, err := KeySwitchValues(queryName, []string{value}, targetPubKey)
	return results[0], err
}

// KeySwitchValues makes request through unlynx to key switch encrypted values
func KeySwitchValues(queryName string, values []string, targetPubKey string) (keySwitchedValues []string, err error) {
	unlynxClient, cothorityRoster := newUnlynxClient()

	// deserialize values and target public key
	desValues, err := deserializeCipherVector(values)
	if err != nil {
		return
	}

	desTargetKey, err := libunlynx.DeserializePoint(targetPubKey)
	if err != nil {
		logrus.Error("unlynx error deserializing target public key: ", err)
		return
	}

	// execute key switching request
	ksResultsChan := make(chan libunlynx.CipherVector)
	ksErrChan := make(chan error)
	go func() {
		_, ksResult, ksErr := unlynxClient.SendSurveyKSRequest(
			cothorityRoster,
			servicesmedco.SurveyID(queryName + "_KS"),
			desTargetKey,
			desValues,
			false,
		)
		if ksErr != nil {
			ksErrChan <- ksErr
		} else {
			ksResultsChan <- ksResult
		}
	}()

	select {
	case ksResult := <-ksResultsChan:
		keySwitchedValues, err = serializeCipherVector(ksResult)
		if err != nil {
			logrus.Error("unlynx error serializing: ", err)
		}

	case err = <-ksErrChan:
		logrus.Error("unlynx error executing key switching: ", err)

	case <-time.After(time.Duration(util.UnlynxTimeoutSeconds) * time.Second):
		err = errors.New("unlynx timeout")
		logrus.Error(err)
	}
	return
}

// ShuffleAndKeySwitchValue makes request through unlynx to shuffle and key switch one value per node
func ShuffleAndKeySwitchValue(queryName string, value string, targetPubKey string) (shuffledKsValue string, err error) {
	unlynxClient, cothorityRoster := newUnlynxClient()

	// deserialize value and target public key
	desValue := libunlynx.CipherText{}
	err = desValue.Deserialize(value)
	if err != nil {
		return
	}

	desTargetKey, err := libunlynx.DeserializePoint(targetPubKey)
	if err != nil {
		logrus.Error("unlynx error deserializing target public key: ", err)
		return
	}

	// execute shuffle and key switching request
	shuffleKsResultsChan := make(chan libunlynx.CipherText)
	shuffleKsErrChan := make(chan error)
	go func() {
		_, shuffleKsResult, shuffleKsErr := unlynxClient.SendSurveyShuffleRequest(
			cothorityRoster,
			servicesmedco.SurveyID(queryName + "_SHUFFLE"),
			desTargetKey,
			desValue,
			false,
		)
		if shuffleKsErr != nil {
			shuffleKsErrChan <- shuffleKsErr
		} else {
			shuffleKsResultsChan <- shuffleKsResult
		}
	}()

	select {
	case shuffleKsResult := <-shuffleKsResultsChan:
		shuffledKsValue, err = shuffleKsResult.Serialize()
		if err != nil {
			logrus.Error("unlynx error serializing: ", err)
		}

	case err = <-shuffleKsErrChan:
		logrus.Error("unlynx error executing shuffle and key switching: ", err)

	case <-time.After(time.Duration(util.UnlynxTimeoutSeconds) * time.Second):
		err = errors.New("unlynx timeout")
		logrus.Error(err)
	}
	return
}

// AggregateAndKeySwitchValue makes request through unlynx to aggregate and key switch one value per node
func AggregateAndKeySwitchValue(queryName string, value string, targetPubKey string) (aggValue string, err error) {
	unlynxClient, cothorityRoster := newUnlynxClient()

	// deserialize value and target public key
	desValue := libunlynx.CipherText{}
	err = desValue.Deserialize(value)
	if err != nil {
		return
	}

	desTargetKey, err := libunlynx.DeserializePoint(targetPubKey)
	if err != nil {
		logrus.Error("unlynx error deserializing target public key: ", err)
		return
	}

	// execute shuffle and key switching request
	aggKsResultsChan := make(chan libunlynx.CipherText)
	aggKsErrChan := make(chan error)
	go func() {
		_, aggKsResult, aggKsErr := unlynxClient.SendSurveyAggRequest(
			cothorityRoster,
			servicesmedco.SurveyID(queryName + "_AGG"),
			desTargetKey,
			desValue,
			false,
		)
		if aggKsErr != nil {
			aggKsErrChan <- aggKsErr
		} else {
			aggKsResultsChan <- aggKsResult
		}
	}()

	select {
	case aggKsResult := <-aggKsResultsChan:
		aggValue, err = aggKsResult.Serialize()
		if err != nil {
			logrus.Error("unlynx error serializing: ", err)
		}

	case err = <-aggKsErrChan:
		logrus.Error("unlynx error executing aggregate and key switching: ", err)

	case <-time.After(time.Duration(util.UnlynxTimeoutSeconds) * time.Second):
		err = errors.New("unlynx timeout")
		logrus.Error(err)
	}
	return
}
