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

// GetQueryTermsDDT makes request through unlynx to compute distributed deterministic tags
func GetQueryTermsDDT(queryName string, encQueryTerms []string) (taggedQueryTerms map[string]string, err error) {
	unlynxClient, cothorityRoster := newUnlynxClient()

	// todo: wrap ddt in go routine, have channel for result + error, select with result / error / timeout
	// todo: get time measurements
	// 	tr.DDTRequestTimeCommunication = totalTime - tr.DDTRequestTimeExec
	//	tr.DDTParsingTime = parsingTime
	//	tr.DDTRequestTimeExec += tr.DDTParsingTime
	//  totalTime = around request

	// deserialize query terms
	deserializedEncQueryTerms, err := deserializeCipherVector(encQueryTerms)
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
			deserializedEncQueryTerms,
			false,
			false,
		)
		if ddtErr != nil {
			ddtErrChan <- ddtErr
		} else if len(ddtResults) == 0 || len(ddtResults) != len(encQueryTerms) {
			ddtErrChan <- errors.New("unlynx inconsistent DDT results: #results=" + strconv.Itoa(len(ddtResults)) + ", #terms=" + strconv.Itoa(len(encQueryTerms)))
		} else {
			ddtResultsChan <- ddtResults
		}
	}()

	select {
	case ddtResults := <-ddtResultsChan:
		taggedQueryTerms = make(map[string]string)
		for i, result := range ddtResults {
			taggedQueryTerms[encQueryTerms[i]] = string(result)
		}

	case err = <-ddtErrChan:
		logrus.Error("unlynx error executing DDT: ", err)

	case <-time.After(time.Duration(util.UnlynxTimeoutSeconds) * time.Second):
		err = errors.New("unlynx timeout")
		logrus.Error(err)
	}
	return
}

// AggregateAndKeySwitchDummyFlags makes request through unlynx to aggregate and key switch encrypted values
func AggregateAndKeySwitchDummyFlags(queryName string, dummyFlags []string, clientPubKey string) (agg string, err error) {
	unlynxClient, cothorityRoster := newUnlynxClient()


	// todo: get time measurements
	// 	tr.AggRequestTimeCommunication = totalTime - tr.DDTRequestTimeExec
	//	tr.LocalAggregationTime = aggregationTime
	//	tr.AggParsingTime = parsingTime
	//	tr.AggRequestTimeExec += tr.AggParsingTime + tr.LocalAggregationTime

	// deserialize dummy flags and client public key
	deserializedDummyFlags, err := deserializeCipherVector(dummyFlags)
	if err != nil {
		return
	}
	deserializedClientPubKey, err := libunlynx.DeserializePoint(clientPubKey)
	if err != nil {
		logrus.Error("unlynx error deserializing client public key: ", err)
		return
	}

	// local aggregation
	aggregate :=  &deserializedDummyFlags[0]
	for i := 1; i < len(deserializedDummyFlags); i++ {
		aggregate.Add(*aggregate, deserializedDummyFlags[i])
	}

	// execute aggregate request
	aggResultChan := make(chan libunlynx.CipherText)
	aggErrChan := make(chan error)
	go func() {
		_, aggResult, _, aggErr := unlynxClient.SendSurveyAggRequest(
			cothorityRoster,
			servicesmedco.SurveyID(queryName + "_AGG"),
			deserializedClientPubKey,
			*aggregate,
			false,
		)
		if aggErr != nil {
			aggErrChan <- aggErr
		} else {
			aggResultChan <- aggResult
		}
	}()

	select {
	case aggResult := <-aggResultChan:
		agg = aggResult.Serialize()

	case err = <-aggErrChan:
		logrus.Error("unlynx error executing aggregation: ", err)

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
