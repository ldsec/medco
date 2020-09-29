package unlynx

import (
	"errors"
	"github.com/ldsec/medco/connector/util/server"
	"github.com/ldsec/medco/unlynx/services"
	"github.com/ldsec/unlynx/lib"
	"github.com/sirupsen/logrus"
	"strconv"
	"time"
)

// DDTagValues makes request through unlynx to compute distributed deterministic tags of encrypted values
func DDTagValues(queryName string, values []string) (taggedValues map[string]string, times map[string]time.Duration, err error) {
	unlynxClient, cothorityRoster := newUnlynxClient()

	// deserialize values
	desValues, err := deserializeCipherVector(values)
	if err != nil {
		return
	}

	// execute DDT
	type DDTResults struct {
		Results []libunlynx.GroupingKey
		Times servicesmedco.TimeResults
		Err error
	}
	ddtResultsChan := make(chan DDTResults)

	go func() {
		_, ddtResults, ddtTimes, ddtErr := unlynxClient.SendSurveyDDTRequestTerms(
			cothorityRoster,
			servicesmedco.SurveyID(queryName + "_DDT"),
			desValues,
			false,
			false,
		)
		ddtResultsChan <- DDTResults{ddtResults, ddtTimes, ddtErr}
	}()

	select {
	case ddtResults := <-ddtResultsChan:
		err = ddtResults.Err
		if err != nil {
			logrus.Error("unlynx error executing DDT: ", err)

		} else if len(ddtResults.Results) == 0 || len(ddtResults.Results) != len(values) {
			err = errors.New("unlynx inconsistent DDT results: #results=" + strconv.Itoa(len(ddtResults.Results)) +
				", #terms=" + strconv.Itoa(len(values)))

		} else {
			times = ddtResults.Times.MapTR
			taggedValues = make(map[string]string)
			for i, result := range ddtResults.Results {
				taggedValues[values[i]] = string(result)
			}
		}

	case <-time.After(time.Duration(utilserver.UnlynxTimeoutSeconds) * time.Second):
		err = errors.New("unlynx timeout")
		logrus.Error(err)
	}
	return
}

// KeySwitchValue makes request through unlynx to key switch a single encrypted value (convenience function)
func KeySwitchValue(queryName string, value string, targetPubKey string) (string, map[string]time.Duration, error) {
	results, times, err := KeySwitchValues(queryName, []string{value}, targetPubKey)
	return results[0], times, err
}

// KeySwitchValues makes request through unlynx to key switch encrypted values
func KeySwitchValues(queryName string, values []string, targetPubKey string) (keySwitchedValues []string, times map[string]time.Duration, err error) {
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
	type KSResults struct {
		Results libunlynx.CipherVector
		Times servicesmedco.TimeResults
		Err error
	}
	ksResultsChan := make(chan KSResults)

	go func() {
		_, ksResult, ksTimes, ksErr := unlynxClient.SendSurveyKSRequest(
			cothorityRoster,
			servicesmedco.SurveyID(queryName + "_KS"),
			desTargetKey,
			desValues,
			false,
		)
		ksResultsChan <- KSResults{ksResult, ksTimes, ksErr}
	}()

	select {
	case ksResult := <-ksResultsChan:
		err = ksResult.Err
		if err != nil {
			logrus.Error("unlynx error executing key switching: ", err)

		} else {
			times = ksResult.Times.MapTR
			keySwitchedValues, err = serializeCipherVector(ksResult.Results)
			if err != nil {
				logrus.Error("unlynx error serializing: ", err)
			}
		}

	case <-time.After(time.Duration(utilserver.UnlynxTimeoutSeconds) * time.Second):
		err = errors.New("unlynx timeout")
		logrus.Error(err)
	}
	return
}

// ShuffleAndKeySwitchValue makes request through unlynx to shuffle and key switch one value per node
func ShuffleAndKeySwitchValue(queryName string, value string, targetPubKey string) (shuffledKsValue string, times map[string]time.Duration, err error) {
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
	type ShuffleKSResults struct {
		Results libunlynx.CipherText
		Times servicesmedco.TimeResults
		Err error
	}
	shuffleKsResultsChan := make(chan ShuffleKSResults)

	go func() {
		_, shuffleKsResult, shuffleKsTimes, shuffleKsErr := unlynxClient.SendSurveyShuffleRequest(
			cothorityRoster,
			servicesmedco.SurveyID(queryName + "_SHUFFLE"),
			desTargetKey,
			&desValue,
			false,
		)
		shuffleKsResultsChan <- ShuffleKSResults{shuffleKsResult, shuffleKsTimes, shuffleKsErr}
	}()

	select {
	case shuffleKsResult := <-shuffleKsResultsChan:
		err = shuffleKsResult.Err
		if err != nil {
			logrus.Error("unlynx error executing shuffle and key switching: ", err)

		} else {
			times = shuffleKsResult.Times.MapTR
			shuffledKsValue, err = shuffleKsResult.Results.Serialize()
			if err != nil {
				logrus.Error("unlynx error serializing: ", err)
			}
		}

	case <-time.After(time.Duration(utilserver.UnlynxTimeoutSeconds) * time.Second):
		err = errors.New("unlynx timeout")
		logrus.Error(err)
	}
	return
}

// AggregateAndKeySwitchValue makes request through unlynx to aggregate and key switch one value per node
func AggregateAndKeySwitchValue(queryName string, value string, targetPubKey string) (aggValue string, times map[string]time.Duration, err error) {
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
	type AggKSResults struct {
		Results libunlynx.CipherText
		Times servicesmedco.TimeResults
		Err error
	}
	aggKsResultsChan := make(chan AggKSResults)

	go func() {
		_, aggKsResult, aggKsTimes, aggKsErr := unlynxClient.SendSurveyAggRequest(
			cothorityRoster,
			servicesmedco.SurveyID(queryName + "_AGG"),
			desTargetKey,
			desValue,
			false,
		)
		aggKsResultsChan <- AggKSResults{aggKsResult, aggKsTimes, aggKsErr}
	}()

	select {
	case aggKsResult := <-aggKsResultsChan:
		err = aggKsResult.Err
		if err != nil {
			logrus.Error("unlynx error executing aggregate and key switching: ", err)

		} else {
			times = aggKsResult.Times.MapTR
			aggValue, err = aggKsResult.Results.Serialize()
			if err != nil {
				logrus.Error("unlynx error serializing: ", err)
			}
		}

	case <-time.After(time.Duration(utilserver.UnlynxTimeoutSeconds) * time.Second):
		err = errors.New("unlynx timeout")
		logrus.Error(err)
	}
	return
}
