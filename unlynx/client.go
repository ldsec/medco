package unlynx

import (
	"errors"
	"github.com/lca1/medco-connector/util"
	"github.com/lca1/medco-unlynx/services"
	"github.com/lca1/unlynx/lib"
	"github.com/sirupsen/logrus"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/app"
	"os"
	"strconv"
	"time"
)

// DDTagValues makes request through unlynx to compute distributed deterministic tags of encrypted values
func DDTagValues(queryName string, values []string) (taggedValues map[string]string, err error) {
	unlynxClient, cothorityRoster := newUnlynxClient()

	// todo: get time measurements
	// 	tr.DDTRequestTimeCommunication = totalTime - tr.DDTRequestTimeExec
	//	tr.DDTParsingTime = parsingTime
	//	tr.DDTRequestTimeExec += tr.DDTParsingTime
	//  totalTime = around request

	// deserialize values
	desValues, err := deserializeCipherVector(values)
	if err != nil {
		return
	}

	// execute DDT
	ddtResultsChan := make(chan []libunlynx.GroupingKey)
	ddtErrChan := make(chan error)
	go func() {
		_, ddtResults, _, ddtErr := unlynxClient.SendSurveyDDTRequestTerms(
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

// AggregateValues adds together several encrypted values homomorphically
func AggregateValues(values []string) (agg string, err error) {

	// deserialize values
	deserialized, err := deserializeCipherVector(values)
	if err != nil {
		return
	}

	// local aggregation
	aggregate :=  &deserialized[0]
	for i := 1; i < len(deserialized); i++ {
		aggregate.Add(*aggregate, deserialized[i])
	}

	return aggregate.Serialize(), nil
}

// KeySwitchValue makes request through unlynx to key switch encrypted values
func KeySwitchValue(queryName string, value string, targetPubKey string) (keySwitchedValue string, err error) {
	unlynxClient, cothorityRoster := newUnlynxClient()

	// todo: get time measurements
	// 	tr.AggRequestTimeCommunication = totalTime - tr.DDTRequestTimeExec
	//	tr.LocalAggregationTime = aggregationTime
	//	tr.AggParsingTime = parsingTime
	//	tr.AggRequestTimeExec += tr.AggParsingTime + tr.LocalAggregationTime

	// deserialize value and target public key
	desValue := libunlynx.CipherText{}
	err = desValue.Deserialize(value)
	if err != nil {
		logrus.Error("unlynx error deserializing cipher text: ", err)
		return
	}

	desTargetKey, err := libunlynx.DeserializePoint(targetPubKey)
	if err != nil {
		logrus.Error("unlynx error deserializing target public key: ", err)
		return
	}

	// execute key switching request
	ksResultChan := make(chan libunlynx.CipherText)
	ksErrChan := make(chan error)
	go func() {
		_, ksResult, _, ksErr := unlynxClient.SendSurveyAggRequest(
			cothorityRoster,
			servicesmedco.SurveyID(queryName + "_KS"),
			desTargetKey,
			desValue,
			false,
		)
		if ksErr != nil {
			ksErrChan <- ksErr
		} else {
			ksResultChan <- ksResult
		}
	}()

	select {
	case ksResult := <-ksResultChan:
		keySwitchedValue = ksResult.Serialize()

	case err = <-ksErrChan:
		logrus.Error("unlynx error executing key switching: ", err)

	case <-time.After(time.Duration(util.UnlynxTimeoutSeconds) * time.Second):
		err = errors.New("unlynx timeout")
		logrus.Error(err)
	}
	return
}

// deserializeCipherVector deserializes string-encoded cipher texts into a vector
func deserializeCipherVector(cipherTexts []string) (cipherVector libunlynx.CipherVector, err error) {
	for _, cipherText := range cipherTexts {
		deserialized := libunlynx.CipherText{}
		err = deserialized.Deserialize(cipherText)
		if err != nil {
			logrus.Error("unlynx error deserializing cipher text: ", err)
			return
		}

		cipherVector = append(cipherVector, deserialized)
	}
	return
}

// newUnlynxClient creates a new client to communicate with unlynx
func newUnlynxClient() (unlynxClient *servicesmedco.API, cothorityRoster *onet.Roster) {

	// initialize medco client
	groupFile, err := os.Open(util.UnlynxGroupFilePath)
	if err != nil {
		logrus.Panic("unlynx error opening group file: ", err)
	}

	group, err := app.ReadGroupDescToml(groupFile)
	if err != nil || len(group.Roster.List) <= 0 {
		logrus.Panic("unlynx error parsing group file: ", err)
	}

	cothorityRoster = group.Roster
	unlynxClient = servicesmedco.NewMedCoClient(
		cothorityRoster.List[util.UnlynxGroupFileIdx],
		strconv.Itoa(util.UnlynxGroupFileIdx),
	)

	return
}
