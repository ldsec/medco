package servicesmedco

import (
	"fmt"
	"github.com/ldsec/unlynx/lib"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/network"
	"time"
)

func init() {
	network.RegisterMessage(ProtocolConfig{})
	network.RegisterMessage(SurveyDDTRequest{})
}

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
	TaggingTimeExec          = "TaggingTimeExec"
	TaggingTimeCommunication = "TaggingTimeCommunication"
	DDTRequestTime           = "DDTRequestTime"

	KSTimeExec          = "KSTimeExec"
	KSTimeCommunication = "KSTimeCommunication"
	KSRequestTime       = "KSRequestTime"

	ShuffleTimeExec          = "ShuffleTimeExec"
	ShuffleTimeCommunication = "ShuffleTimeCommunication"
	ShuffleRequestTime       = "ShuffleRequestTime"

	AggrTime        = "AggrTime"
	AggrRequestTime = "AggrRequestTime"
)

// ResultDDT will contain final results of the DDT of the query terms.
type ResultDDT struct {
	Result []libunlynx.GroupingKey
	TR     map[string]time.Duration
}

// Result will contain the final results for the other queries
type Result struct {
	Result libunlynx.CipherVector
	TR     TimeResults
}

// SurveyID unique ID for each survey.
type SurveyID string

// ProtocolConfig holds the configuration that will be passed from node to
// node. It replaces the map and calling all nodes.
type ProtocolConfig struct {
	SurveyID SurveyID
	TypeQ    string
	Data     []byte
}

// SurveyDDTRequest is the message used trigger the DDT of the query parameters
type SurveyDDTRequest struct {
	SurveyID SurveyID
	Roster   onet.Roster
	Proofs   bool
	Testing  bool

	Terms libunlynx.CipherVector // query terms

	// message handling
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
	MessageSource *network.ServerIdentity
}

// SurveyAggRequest is the message used trigger the aggregation of the final results (+ key switching)
type SurveyAggRequest struct {
	SurveyID     SurveyID
	Roster       onet.Roster
	Proofs       bool
	ClientPubKey kyber.Point // we need this for the key switching

	AggregateTarget libunlynx.CipherVector // target results to aggregate. the root node adds the results from the other nodes here
	KSTarget        libunlynx.CipherVector // the final aggregated results to be key switched
}

// SurveyKS is the struct that we persist in the service that contains all the data for the Key Switch request phase
type SurveyKS struct {
	SurveyID SurveyID
	Request  SurveyKSRequest
	TR       TimeResults
}

// SurveyShuffle is the struct that we persist in the service that contains all the data for the Shuffle (+KS) request phase
type SurveyShuffle struct {
	SurveyID            SurveyID
	Request             SurveyShuffleRequest
	SurveyChannel       chan int // To wait for all the aggregate results to be received by the root node
	FinalResultsChannel chan int
	TR                  TimeResults
}

// SurveyAgg is the struct that we persist in the service that contains all the data for the Aggregation request phase
type SurveyAgg struct {
	SurveyID SurveyID
	Request  SurveyAggRequest
	TR       TimeResults
}

// SurveyShuffleGenerated is used to ensure that the root server creates the survey before all the other nodes send it their results
type SurveyShuffleGenerated struct {
	SurveyID SurveyID
}

// SurveyAggGenerated is used to ensure that the root server creates the survey before all the other nodes send it their results
type SurveyAggGenerated struct {
	SurveyID SurveyID
}

func (s *Service) deleteSurveyKS(sid SurveyID) (SurveyKS, error) {
	surv, err := s.MapSurveyKS.Remove(string(sid))
	if err != nil {
		return SurveyKS{}, fmt.Errorf("error while deleting surveyID ("+string(sid)+"): %v", err.Error())
	}
	if surv == nil {
		return SurveyKS{}, fmt.Errorf("no entry in map with surveyID (" + string(sid) + ")")
	}
	return surv.(SurveyKS), nil
}

func (s *Service) deleteSurveyShuffle(sid SurveyID) (SurveyShuffle, error) {
	surv, err := s.MapSurveyShuffle.Remove(string(sid))
	if err != nil {
		return SurveyShuffle{}, fmt.Errorf("error while deleting surveyID ("+string(sid)+"): %v", err.Error())
	}
	if surv == nil {
		return SurveyShuffle{}, fmt.Errorf("no entry in map with surveyID (" + string(sid) + ")")
	}
	return surv.(SurveyShuffle), nil
}

func (s *Service) deleteSurveyAgg(sid SurveyID) (SurveyAgg, error) {
	surv, err := s.MapSurveyAgg.Remove(string(sid))
	if err != nil {
		return SurveyAgg{}, fmt.Errorf("error while deleting surveyID ("+string(sid)+"): %v", err.Error())
	}
	if surv == nil {
		return SurveyAgg{}, fmt.Errorf("no entry in map with surveyID (" + string(sid) + ")")
	}
	return surv.(SurveyAgg), nil
}

func (s *Service) getSurveyKS(sid SurveyID) (SurveyKS, error) {
	surv, err := s.MapSurveyKS.Get(string(sid))
	if err != nil {
		return SurveyKS{}, fmt.Errorf("error while getting surveyID ("+string(sid)+"): %v", err.Error())
	}
	if surv == nil {
		return SurveyKS{}, fmt.Errorf("empty map entry while getting surveyID (" + string(sid) + ")")
	}
	return surv.(SurveyKS), nil
}

func (s *Service) getSurveyShuffle(sid SurveyID) (SurveyShuffle, error) {
	surv, err := s.MapSurveyShuffle.Get(string(sid))
	if err != nil {
		return SurveyShuffle{}, fmt.Errorf("error while getting surveyID ("+string(sid)+"): %v", err.Error())
	}
	if surv == nil {
		return SurveyShuffle{}, fmt.Errorf("empty map entry while getting surveyID (" + string(sid) + ")")
	}
	return surv.(SurveyShuffle), nil
}

func (s *Service) getSurveyAgg(sid SurveyID) (SurveyAgg, error) {
	surv, err := s.MapSurveyAgg.Get(string(sid))
	if err != nil {
		return SurveyAgg{}, fmt.Errorf("error while getting surveyID ("+string(sid)+"): %v", err.Error())
	}
	if surv == nil {
		return SurveyAgg{}, fmt.Errorf("empty map entry while getting surveyID (" + string(sid) + ")")
	}
	return surv.(SurveyAgg), nil
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

func unmarshalProtocolConfig(buf []byte) (pc ProtocolConfig, err error) {
	_, pcInt, err := network.Unmarshal(buf, libunlynx.SuiTe)
	if err != nil {
		return
	}
	pc = *pcInt.(*ProtocolConfig)
	return
}

func newProtocolConfig(sid SurveyID, tq string, data interface{}) (
	pc ProtocolConfig, err error) {
	pc.SurveyID = sid
	pc.TypeQ = tq
	pc.Data, err = network.Marshal(data)
	return
}

func (pc ProtocolConfig) getConfig() (gc onet.GenericConfig, err error) {
	gc.Data, err = network.Marshal(&pc)
	return
}

func (pc ProtocolConfig) getTarget() SurveyID {
	if pc.TypeQ == "" {
		return pc.SurveyID
	}
	return SurveyID(fmt.Sprintf("%s/%s", pc.SurveyID, pc.TypeQ))
}
