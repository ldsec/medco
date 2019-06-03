package servicesmedco

import (
	"github.com/btcsuite/goleveldb/leveldb/errors"
	"github.com/lca1/unlynx/lib"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/network"
	"time"
)

// Name is the registered name for the medco service.
const Name = "medco"

// DDTSecretsPath filename
const DDTSecretsPath = "secrets"

// Name of query/request types (important to distinguish which map to use during key switching)

// KSRequestName the name of this type of query
const KSRequestName = "KSRequestName"

// ShuffleRequestName the name of this type of query
const ShuffleRequestName = "ShuffleRequestName"

// AggRequestName the name of this type of query
const AggRequestName = "AggRequestName"

// TimeResults includes all variables that will store the durations (to collect the execution/communication time)
type TimeResults struct {
	MapTR map[string]time.Duration
}

// timer constant names
const (
	TaggingTimeExec = "TaggingTimeExec"
	TaggingTimeCommunication = "TaggingTimeCommunication"
	DDTRequestTime = "DDTRequestTime"

	KSTimeExec = "KSTimeExec"
	KSTimeCommunication = "KSTimeCommunication"
	KSRequestTime = "KSRequestTime"

	ShuffleTimeExec = "ShuffleTimeExec"
	ShuffleTimeCommunication = "ShuffleTimeCommunication"
	ShuffleRequestTime = "ShuffleRequestTime"

	AggrTime = "AggrTime"
	AggrRequestTime = "AggrRequestTime"
)

// ResultDDT will contain final results of the DDT of the query terms.
type ResultDDT struct {
	Result []libunlynx.GroupingKey
	TR     TimeResults
}

// Result will contain the final results for the other queries
type Result struct {
	Result libunlynx.CipherVector
	TR     TimeResults
}

// SurveyID unique ID for each survey.
type SurveyID string

// SurveyDDTRequest is the message used trigger the DDT of the query parameters
type SurveyDDTRequest struct {
	SurveyID SurveyID
	Roster   onet.Roster
	Proofs   bool
	Testing  bool

	Terms libunlynx.CipherVector // query terms

	// message handling
	IntraMessage  bool
	MessageSource *network.ServerIdentity
}

// SurveyKSRequest is the message used trigger the key switching of the results
type SurveyKSRequest struct {
	SurveyID     SurveyID
	Roster       onet.Roster
	Proofs       bool
	ClientPubKey kyber.Point // we need this for the key switching

	KSTarget libunlynx.CipherVector // target values to key switch
}

// SurveyShuffleRequest is the message used trigger the shuffling and key switching of the final results
type SurveyShuffleRequest struct {
	SurveyID     SurveyID
	Roster       onet.Roster
	Proofs       bool
	ClientPubKey kyber.Point // we need this for the key switching

	ShuffleTarget libunlynx.CipherVector // target results to shuffle. the root node adds the results from the other nodes here
	KSTarget      libunlynx.CipherVector // the final results to be key switched

	// message handling
	IntraMessage  bool
	MessageSource *network.ServerIdentity
}

// SurveyAggRequest is the message used trigger the aggregation of the final results (+ key switching)
type SurveyAggRequest struct {
	SurveyID     SurveyID
	Roster       onet.Roster
	Proofs       bool
	ClientPubKey kyber.Point // we need this for the key switching

	AggregateTarget libunlynx.CipherText // target results to aggregate. the root node adds the results from the other nodes here
	KSTarget        libunlynx.CipherText // the final aggregated result to be key switched
}

// SurveyTag is the struct that we persist in the service that contains all the data for the DDT protocol
type SurveyTag struct {
	SurveyID      SurveyID
	Request       SurveyDDTRequest
	SurveyChannel chan int // To wait for the survey to be created before the DDT protocol
	TR 			  TimeResults
}

// SurveyKS is the struct that we persist in the service that contains all the data for the Key Switch request phase
type SurveyKS struct {
	SurveyID SurveyID
	Request  SurveyKSRequest
	TR 	     TimeResults
}

// SurveyShuffle is the struct that we persist in the service that contains all the data for the Shuffle (+KS) request phase
type SurveyShuffle struct {
	SurveyID            SurveyID
	Request             SurveyShuffleRequest
	SurveyChannel       chan int // To wait for all the aggregate results to be received by the root node
	FinalResultsChannel chan int
	TR 	     			TimeResults
}

// SurveyAgg is the struct that we persist in the service that contains all the data for the Aggregation request phase
type SurveyAgg struct {
	SurveyID      SurveyID
	Request       SurveyAggRequest
	SurveyChannel chan int // To wait for all the aggregate results to be received by the root node
	TR 			  TimeResults
}

// SurveyTagGenerated is used to ensure that all servers get the survey before starting the DDT protocol
type SurveyTagGenerated struct {
	SurveyID SurveyID
}

// SurveyShuffleGenerated is used to ensure that the root server creates the survey before all the other nodes send it their results
type SurveyShuffleGenerated struct {
	SurveyID SurveyID
}

// SurveyAggGenerated is used to ensure that the root server creates the survey before all the other nodes send it their results
type SurveyAggGenerated struct {
	SurveyID SurveyID
}

func (s *Service) getSurveyTag(sid SurveyID) (SurveyTag, error) {
	surv, err := s.MapSurveyTag.Get(string(sid))
	if err != nil {
		return SurveyTag{}, errors.New("Error" + err.Error() + "while getting surveyID: " + string(sid))
	}
	if surv == nil {
		return SurveyTag{}, errors.New("Empty map entry while getting surveyID: " + string(sid))
	}
	return surv.(SurveyTag), nil
}

func (s *Service) getSurveyKS(sid SurveyID) (SurveyKS, error) {
	surv, err := s.MapSurveyKS.Get(string(sid))
	if err != nil {
		return SurveyKS{}, errors.New("Error" + err.Error() + "while getting surveyID: " + string(sid))
	}
	if surv == nil {
		return SurveyKS{}, errors.New("Empty map entry while getting surveyID: " + string(sid))
	}
	return surv.(SurveyKS), nil
}

func (s *Service) getSurveyShuffle(sid SurveyID) (SurveyShuffle, error) {
	surv, err := s.MapSurveyShuffle.Get(string(sid))
	if err != nil {
		return SurveyShuffle{}, errors.New("Error" + err.Error() + "while getting surveyID: " + string(sid))
	}
	if surv == nil {
		return SurveyShuffle{}, errors.New("Empty map entry while getting surveyID: " + string(sid))
	}
	return surv.(SurveyShuffle), nil
}

func (s *Service) getSurveyAgg(sid SurveyID) (SurveyAgg, error) {
	surv, err := s.MapSurveyAgg.Get(string(sid))
	if err != nil {
		return SurveyAgg{}, errors.New("Error" + err.Error() + "while getting surveyID" + string(sid))
	}
	if surv == nil {
		return SurveyAgg{}, errors.New("Empty map entry while getting surveyID" + string(sid))
	}
	return surv.(SurveyAgg), nil
}

func (s *Service) putSurveyTag(sid SurveyID, surv SurveyTag) error {
	_, err := s.MapSurveyTag.Put(string(sid), surv)
	return err
}

func (s *Service) putSurveyKS(sid SurveyID, surv SurveyKS) error {
	_, err := s.MapSurveyKS.Put(string(sid), surv)
	return err
}

func (s *Service) putSurveyShuffle(sid SurveyID, surv SurveyShuffle) error {
	_, err := s.MapSurveyShuffle.Put(string(sid), surv)
	return err
}

func (s *Service) putSurveyAgg(sid SurveyID, surv SurveyAgg) error {
	_, err := s.MapSurveyAgg.Put(string(sid), surv)
	return err
}
