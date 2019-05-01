package servicesmedco

import (
	"encoding/base64"
	"github.com/BurntSushi/toml"
	"github.com/btcsuite/goleveldb/leveldb/errors"
	"github.com/fanliao/go-concurrentMap"
	"github.com/lca1/unlynx/lib"
	"github.com/lca1/unlynx/lib/tools"
	"github.com/lca1/unlynx/protocols"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/util/random"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
	"os"
	"sync"
	"time"
)

// Name is the registered name for the medco service.
const Name = "medco"

// DDTSecretsPath filename
const DDTSecretsPath = "secrets"

// TimeResults includes all variables that will store the durations (to collect the execution/communication time)
type TimeResults struct {
	DDTParsingTime              time.Duration // Total parsing time (i2b2 -> unlynx client)
	DDTRequestTimeExec          time.Duration // Total DDT (of the request) execution time
	DDTRequestTimeCommunication time.Duration // Total DDT (of the request) communication time

	AggParsingTime              time.Duration // Total parsing time (i2b2 -> unlynx client)
	AggRequestTimeExec          time.Duration // Total Agg (of the request) execution time
	AggRequestTimeCommunication time.Duration // Total Agg (of the request) communication time
	LocalAggregationTime        time.Duration // Total local aggregation time
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

// SurveyAggRequest is the message used trigger the aggregation of the final results (well it's mostly shuffling and key switching)
type SurveyAggRequest struct {
	SurveyID     SurveyID
	Roster       onet.Roster
	Proofs       bool
	ClientPubKey kyber.Point // we need this for the key switching

	Aggregate          libunlynx.CipherVector // aggregated final result. It is an array because we the root node adds the results from the other nodes here
	AggregateShuffled  libunlynx.CipherVector // aggregated final results after they are shuffled
	AggregateKSwitched libunlynx.CipherVector // the final results after the key switching

	// message handling
	IntraMessage  bool
	MessageSource *network.ServerIdentity
}

// SurveyTag is the struct that we persist in the service that contains all the data for the DDT protocol
type SurveyTag struct {
	SurveyID      SurveyID
	Request       SurveyDDTRequest
	SurveyChannel chan int    // To wait for the survey to be created before the DDT protocol
	TR            TimeResults // contains all the time measurements
}

// SurveyAgg is the struct that we persist in the service that contains all the data for the Aggregation request phase
type SurveyAgg struct {
	SurveyID            SurveyID
	Request             SurveyAggRequest
	SurveyChannel       chan int    // To wait for all the aggregate results to be received by the root node
	FinalResultsChannel chan int    // To wait for the final key switched results
	TR                  TimeResults // contains all the time measurements
}

// SurveyTagGenerated is used to ensure that all servers get the survey before starting the DDT protocol
type SurveyTagGenerated struct {
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

func (s *Service) putSurveyAgg(sid SurveyID, surv SurveyAgg) error {
	_, err := s.MapSurveyAgg.Put(string(sid), surv)
	return err
}

/*func castToSurveyTag(object interface{}, err error) SurveyTag {
	if err != nil {
		log.Error("Error reading SurveyTag map")
	}
	return object.(SurveyTag)
}

func castToSurveyAgg(object interface{}, err error) SurveyAgg {
	if err != nil {
		log.Error("Error reading SurveyAgg map")
	}
	return object.(SurveyAgg)
}*/

// MsgTypes defines the Message Type ID for all the service's intra-messages.
type MsgTypes struct {
	msgSurveyDDTRequestTerms network.MessageTypeID
	msgSurveyTagGenerated    network.MessageTypeID
	msgSurveyAggRequest      network.MessageTypeID
	msgSurveyAggGenerated    network.MessageTypeID
}

var msgTypes = MsgTypes{}

func init() {
	_, err := onet.RegisterNewService(Name, NewService)
	log.ErrFatal(err)

	// messages for DDT Request
	msgTypes.msgSurveyDDTRequestTerms = network.RegisterMessage(&SurveyDDTRequest{})
	msgTypes.msgSurveyTagGenerated = network.RegisterMessage(&SurveyTagGenerated{})
	network.RegisterMessage(&ResultDDT{})

	// messages for Agg Request
	msgTypes.msgSurveyAggRequest = network.RegisterMessage(&SurveyAggRequest{})
	msgTypes.msgSurveyAggGenerated = network.RegisterMessage(&SurveyAggGenerated{})
	network.RegisterMessage(&ResultAgg{})
}

// ResultDDT will contain final results of the DDT of the query terms.
type ResultDDT struct {
	Result []libunlynx.GroupingKey
	TR     TimeResults // contains all the time measurements
}

// ResultAgg will contain final aggregate result to sent to the client.
type ResultAgg struct {
	Result libunlynx.CipherText
	TR     TimeResults // contains all the time measurements
}

// Service defines a service in unlynx
type Service struct {
	*onet.ServiceProcessor

	MapSurveyTag *concurrent.ConcurrentMap
	MapSurveyAgg *concurrent.ConcurrentMap
	Mutex        *sync.Mutex
}

// NewService constructor which registers the needed messages.
func NewService(c *onet.Context) (onet.Service, error) {

	newUnLynxInstance := &Service{
		ServiceProcessor: onet.NewServiceProcessor(c),
		MapSurveyTag:     concurrent.NewConcurrentMap(),
		MapSurveyAgg:     concurrent.NewConcurrentMap(),
		Mutex:            &sync.Mutex{},
	}

	if cerr := newUnLynxInstance.RegisterHandler(newUnLynxInstance.HandleSurveyDDTRequestTerms); cerr != nil {
		log.Error("Wrong Handler.", cerr)
	}
	if cerr := newUnLynxInstance.RegisterHandler(newUnLynxInstance.HandleSurveyAggRequest); cerr != nil {
		log.Error("Wrong Handler.", cerr)
	}

	c.RegisterProcessor(newUnLynxInstance, msgTypes.msgSurveyDDTRequestTerms)
	c.RegisterProcessor(newUnLynxInstance, msgTypes.msgSurveyTagGenerated)

	c.RegisterProcessor(newUnLynxInstance, msgTypes.msgSurveyAggRequest)
	c.RegisterProcessor(newUnLynxInstance, msgTypes.msgSurveyAggGenerated)

	return newUnLynxInstance, nil
}

// Process implements the processor interface and is used to recognize messages broadcasted between servers
func (s *Service) Process(msg *network.Envelope) {
	if msg.MsgType.Equal(msgTypes.msgSurveyDDTRequestTerms) {
		tmp := (msg.Msg).(*SurveyDDTRequest)
		_, err := s.HandleSurveyDDTRequestTerms(tmp)
		if err != nil {
			log.Error(err)
		}
	} else if msg.MsgType.Equal(msgTypes.msgSurveyTagGenerated) {
		tmp := (msg.Msg).(*SurveyTagGenerated)
		_, err := s.HandleSurveyTagGenerated(tmp)
		if err != nil {
			log.Error(err)
		}
	} else if msg.MsgType.Equal(msgTypes.msgSurveyAggRequest) {
		tmp := (msg.Msg).(*SurveyAggRequest)
		_, err := s.HandleSurveyAggRequest(tmp)
		if err != nil {
			log.Error(err)
		}
	} else if msg.MsgType.Equal(msgTypes.msgSurveyAggGenerated) {
		tmp := (msg.Msg).(*SurveyAggGenerated)
		_, err := s.HandleSurveyAggGenerated(tmp)
		if err != nil {
			log.Error(err)
		}
	} else {
		log.Error("Cannot identify the intra-message")
	}
}

// Request Handlers
//______________________________________________________________________________________________________________________

// HandleSurveyTagGenerated handles triggers the SurveyDDTChannel
func (s *Service) HandleSurveyTagGenerated(recq *SurveyTagGenerated) (network.Message, error) {
	surveyTag, err := s.getSurveyTag(recq.SurveyID)
	if err != nil {
		return nil, err
	}
	surveyTag.SurveyChannel <- 1
	return nil, nil
}

// HandleSurveyDDTRequestTerms handles the reception of the query terms to be deterministically tagged
func (s *Service) HandleSurveyDDTRequestTerms(sdq *SurveyDDTRequest) (network.Message, error) {

	// if this server is the one receiving the request from the client
	if !sdq.IntraMessage {
		log.Lvl1(s.ServerIdentity().String(), " received a SurveyDDTRequestTerms:", sdq.SurveyID)

		if len(sdq.Terms) == 0 {
			log.Lvl1(s.ServerIdentity(), " for survey", sdq.SurveyID, "has no data to det tag")
			return &ResultDDT{}, nil
		}

		// initialize timers
		err := s.putSurveyTag(sdq.SurveyID,
			SurveyTag{
				SurveyID:      sdq.SurveyID,
				Request:       *sdq,
				SurveyChannel: make(chan int, 100),
				TR:            TimeResults{DDTRequestTimeExec: 0, DDTRequestTimeCommunication: 0},
			})
		if err != nil {
			return nil, err
		}

		// signal the other nodes that they need to prepare to execute a DDT (no need to send the terms
		// we only need the message source so that they know which node requested the DDT and fetch the secret accordingly)
		err = libunlynxtools.SendISMOthers(s.ServiceProcessor, &sdq.Roster,
			&SurveyDDTRequest{
				SurveyID:      sdq.SurveyID,
				Roster:        sdq.Roster,
				IntraMessage:  true,
				MessageSource: s.ServerIdentity(),
				Proofs:        sdq.Proofs,
				Testing:       sdq.Testing,
			})
		if err != nil {
			return nil, err
		}

		surveyTag, err := s.getSurveyTag(sdq.SurveyID)
		if err != nil {
			return nil, err
		}

		// waits for all other nodes to receive the survey
		counter := len(sdq.Roster.List) - 1
		for counter > 0 {
			counter = counter - <-surveyTag.SurveyChannel
		}

		deterministicTaggingResult, err := s.TaggingPhase(sdq.SurveyID, &sdq.Roster)
		if err != nil {
			return nil, err
		}

		start := time.Now()

		// convert the result to of the tagging for something close to the XML response of i2b2 (array of tagged terms)
		listTaggedTerms := make([]libunlynx.GroupingKey, 0)

		for _, el := range deterministicTaggingResult {
			listTaggedTerms = append(listTaggedTerms, libunlynx.GroupingKey(el.String()))
		}

		surveyTag, err = s.getSurveyTag(sdq.SurveyID)
		if err != nil {
			return nil, err
		}
		surveyTag.TR.DDTRequestTimeExec += time.Since(start)

		tr := surveyTag.TR
		return &ResultDDT{Result: listTaggedTerms, TR: tr}, nil
	}

	log.Lvl1(s.ServerIdentity().String(), " is notified of survey:", sdq.SurveyID)

	err := s.putSurveyTag(sdq.SurveyID, SurveyTag{
		SurveyID: sdq.SurveyID,
		Request:  *sdq,
	})

	// sends a signal to unlock waiting channel
	err = s.SendRaw(sdq.MessageSource, &SurveyTagGenerated{SurveyID: sdq.SurveyID})
	if err != nil {
		log.Error("sending error ", err)
	}

	return nil, nil
}

// HandleSurveyAggGenerated handles triggers the SurveyDDTChannel
func (s *Service) HandleSurveyAggGenerated(recq *SurveyAggGenerated) (network.Message, error) {
	var el interface{}
	el = nil
	for el == nil {
		el, _ = s.MapSurveyAgg.Get((string)(recq.SurveyID))
		if el != nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	surveyAgg, err := s.getSurveyAgg(recq.SurveyID)
	if err != nil {
		return nil, err
	}
	surveyAgg.SurveyChannel <- 1
	return nil, nil
}

// HandleSurveyAggRequest handles the reception of the aggregate local result to be shared/shuffled/switched
func (s *Service) HandleSurveyAggRequest(sar *SurveyAggRequest) (network.Message, error) {
	var root bool
	if s.ServerIdentity().String() == sar.Roster.List[0].String() {
		root = true
	} else {
		root = false
	}

	log.Lvl1(s.ServerIdentity().String(), " received a SurveyAggRequest:", sar.SurveyID, "(root =", root, "- intra =", sar.IntraMessage, ")")

	// (root = true - intra = false )
	if !sar.IntraMessage && root {
		// initialize timers

		err := s.putSurveyAgg(sar.SurveyID, SurveyAgg{
			SurveyID:      sar.SurveyID,
			Request:       *sar,
			SurveyChannel: make(chan int, 100),
			TR:            TimeResults{AggRequestTimeExec: 0, AggRequestTimeCommunication: 0}})

		// send signal to unlock the other nodes
		err = libunlynxtools.SendISMOthers(s.ServiceProcessor, &sar.Roster, &SurveyAggGenerated{SurveyID: sar.SurveyID})
		if err != nil {
			log.Error("broadcasting error ", err)
		}

		surveyAgg, err := s.getSurveyAgg(sar.SurveyID)
		if err != nil {
			return nil, err
		}

		// wait until you've got all the aggregate results from the other nodes
		counter := len(sar.Roster.List) - 1
		for counter > 0 {
			counter = counter - <-surveyAgg.SurveyChannel
		}

		surveyAgg, err = s.getSurveyAgg(sar.SurveyID)
		if err != nil {
			return nil, err
		}
		if len(surveyAgg.Request.Aggregate) == 0 {
			log.Lvl1(s.ServerIdentity(), " no data to shuffle")
		} else {
			// shuffle the results
			shufflingResult, err := s.ShufflingPhase(sar.SurveyID, &sar.Roster)

			if err != nil {
				log.Error("shuffling error", err)
				return nil, err
			}

			shufflingFinalResult := make(libunlynx.CipherVector, 0)
			for _, el := range shufflingResult {
				shufflingFinalResult = append(shufflingFinalResult, el[0])
			}

			surveyAgg.Request.AggregateShuffled = shufflingFinalResult

			err = s.putSurveyAgg(sar.SurveyID, surveyAgg)
			if err != nil {
				return nil, err
			}

			// send the shuffled results to all the other nodes
			sar.AggregateShuffled = shufflingFinalResult
			sar.IntraMessage = true
			sar.MessageSource = s.ServerIdentity()

			// let's delete what we don't need (less communication time)
			sar.Aggregate = nil

			// signal the other nodes that they need to prepare to execute a key switching
			// basically after shuffling the results the root server needs to send them back
			// to the remaining nodes for key switching
			err = libunlynxtools.SendISMOthers(s.ServiceProcessor, &sar.Roster, sar)
			if err != nil {
				log.Error("broadcasting error ", err)
			}

			// key switch the results
			keySwitchingResult, err := s.KeySwitchingPhase(sar.SurveyID, &sar.Roster)

			if err != nil {
				log.Error("key switching error", err)
				return nil, err
			}

			// get server index
			index := 0
			for i, r := range sar.Roster.List {
				if r.String() == s.ServerIdentity().String() {
					index = i
					break
				}
			}

			surveyAgg, err = s.getSurveyAgg(sar.SurveyID)
			if err != nil {
				return nil, err
			}
			return &ResultAgg{Result: keySwitchingResult[index], TR: surveyAgg.TR}, nil
		}
		//(root = false - intra = false )
	} else if !sar.IntraMessage && !root { // if message sent by client and not to root
		err := s.putSurveyAgg(sar.SurveyID, SurveyAgg{
			SurveyID:            sar.SurveyID,
			Request:             *sar,
			SurveyChannel:       make(chan int, 100),
			FinalResultsChannel: make(chan int, 100),
			TR:                  TimeResults{AggRequestTimeExec: 0, AggRequestTimeCommunication: 0},
		})

		sar.IntraMessage = true
		sar.MessageSource = s.ServerIdentity()

		surveyAgg, err := s.getSurveyAgg(sar.SurveyID)
		if err != nil {
			return nil, err
		}
		<-surveyAgg.SurveyChannel

		// send your local aggregate result to the root server (index 0)
		err = s.SendRaw(sar.Roster.List[0], sar)
		if err != nil {
			log.Error(s.ServerIdentity().String()+"could not send its aggregate value", err)
		}

		surveyAgg, err = s.getSurveyAgg(sar.SurveyID)
		if err != nil {
			return nil, err
		}
		//waits for the final results to be ready
		<-surveyAgg.FinalResultsChannel

		// get server index
		index := 0
		for i, r := range sar.Roster.List {
			if r.String() == s.ServerIdentity().String() {
				index = i
				break
			}
		}

		surveyAgg, err = s.getSurveyAgg(sar.SurveyID)
		if err != nil {
			return nil, err
		}

		return &ResultAgg{Result: surveyAgg.Request.AggregateKSwitched[index], TR: surveyAgg.TR}, nil

		// (root = true - intra = true )
	} else if sar.IntraMessage && root { // if message sent by another node and root
		s.Mutex.Lock()
		surveyAgg, err := s.getSurveyAgg(sar.SurveyID)
		if err != nil {
			return nil, err
		}
		surveyAgg.Request.Aggregate = append(surveyAgg.Request.Aggregate, sar.Aggregate...)
		err = s.putSurveyAgg(sar.SurveyID, surveyAgg)
		if err != nil {
			return nil, err
		}
		s.Mutex.Unlock()

		// get the request from the other non-root nodes
		surveyAgg, err = s.getSurveyAgg(sar.SurveyID)
		if err != nil {
			return nil, err
		}
		surveyAgg.SurveyChannel <- 1
		// (root = false - intra = true )
	} else { // if message sent by another node and not root
		// update the local survey with the shuffled results
		s.Mutex.Lock()
		surveyAgg, err := s.getSurveyAgg(sar.SurveyID)
		if err != nil {
			return nil, err
		}
		surveyAgg.Request.AggregateShuffled = sar.AggregateShuffled
		err = s.putSurveyAgg(sar.SurveyID, surveyAgg)
		if err != nil {
			return nil, err
		}
		s.Mutex.Unlock()

		// key switch the results
		keySwitchingResult, err := s.KeySwitchingPhase(sar.SurveyID, &sar.Roster)
		if err != nil {
			log.Error("key switching error", err)
			return nil, err
		}

		s.Mutex.Lock()
		surveyAgg, err = s.getSurveyAgg(sar.SurveyID)
		if err != nil {
			return nil, err
		}
		surveyAgg.Request.AggregateKSwitched = keySwitchingResult
		err = s.putSurveyAgg(sar.SurveyID, surveyAgg)
		if err != nil {
			return nil, err
		}
		s.Mutex.Unlock()

		surveyAgg, err = s.getSurveyAgg(sar.SurveyID)
		if err != nil {
			return nil, err
		}
		surveyAgg.FinalResultsChannel <- 1
	}

	return nil, nil
}

// Protocol Handlers
//______________________________________________________________________________________________________________________

// NewProtocol creates a protocol instance executed by all nodes
func (s *Service) NewProtocol(tn *onet.TreeNodeInstance, conf *onet.GenericConfig) (onet.ProtocolInstance, error) {
	if err := tn.SetConfig(conf); err != nil {
		return nil, err
	}

	var pi onet.ProtocolInstance
	target := SurveyID(string(conf.Data))

	switch tn.ProtocolName() {
	case protocolsunlynx.DeterministicTaggingProtocolName:
		surveyTag, err := s.getSurveyTag(target)
		if err != nil {
			return nil, err
		}

		pi, err = protocolsunlynx.NewDeterministicTaggingProtocol(tn)
		if err != nil {
			return nil, err
		}
		hashCreation := pi.(*protocolsunlynx.DeterministicTaggingProtocol)

		var serverIDMap *network.ServerIdentity

		if tn.IsRoot() {
			dataToDDT := surveyTag.Request.Terms
			hashCreation.TargetOfSwitch = &dataToDDT

			serverIDMap = s.ServerIdentity()
		} else {
			serverIDMap = surveyTag.Request.MessageSource
		}

		s.Mutex.Lock()

		var aux kyber.Scalar
		if surveyTag.Request.Testing {
			aux, err = CheckDDTSecrets(DDTSecretsPath+"_"+s.ServerIdentity().Address.Host()+":"+s.ServerIdentity().Address.Port()+".toml", serverIDMap.Address)
			if err != nil || aux == nil {
				return nil, errors.New("Error while reading the DDT secrets from file")
			}
		} else {
			aux, err = CheckDDTSecrets(os.Getenv("UNLYNX_DDT_SECRETS_FILE_PATH"), serverIDMap.Address)
			if err != nil || aux == nil {
				return nil, errors.New("Error while reading the DDT secrets from file")
			}
		}

		hashCreation.SurveySecretKey = &aux
		hashCreation.Proofs = surveyTag.Request.Proofs
		s.Mutex.Unlock()
	case protocolsunlynx.ShufflingProtocolName:
		surveyAgg, err := s.getSurveyAgg(target)
		if err != nil {
			return nil, err
		}

		pi, err := protocolsunlynx.NewShufflingProtocol(tn)
		if err != nil {
			return nil, err
		}

		shuffle := pi.(*protocolsunlynx.ShufflingProtocol)

		shuffle.Proofs = surveyAgg.Request.Proofs
		shuffle.Precomputed = nil

		if tn.IsRoot() {
			dataToShuffle := protocolsunlynx.AdaptCipherTextArray(surveyAgg.Request.Aggregate)
			shuffle.ShuffleTarget = &dataToShuffle
		}
		return pi, nil
	case protocolsunlynx.KeySwitchingProtocolName:
		surveyAgg, err := s.getSurveyAgg(target)
		if err != nil {
			return nil, err
		}

		pi, err = protocolsunlynx.NewKeySwitchingProtocol(tn)
		if err != nil {
			return nil, err
		}

		keySwitch := pi.(*protocolsunlynx.KeySwitchingProtocol)
		keySwitch.Proofs = surveyAgg.Request.Proofs

		if tn.IsRoot() {
			dataToSwitch := surveyAgg.Request.AggregateShuffled
			keySwitch.TargetOfSwitch = &dataToSwitch
			tmp := surveyAgg.Request.ClientPubKey
			keySwitch.TargetPublicKey = &tmp
		}
	default:
		return nil, errors.New("Service attempts to start an unknown protocol: " + tn.ProtocolName() + ".")
	}

	return pi, nil
}

// StartProtocol starts a specific protocol (Shuffling, KeySwitching, etc.)
func (s *Service) StartProtocol(name string, targetSurvey SurveyID, roster *onet.Roster) (onet.ProtocolInstance, error) {
	tree := roster.GenerateNaryTreeWithRoot(2, s.ServerIdentity())
	tn := s.NewTreeNodeInstance(tree, tree.Root, name)

	conf := onet.GenericConfig{Data: []byte(string(targetSurvey))}
	pi, err := s.NewProtocol(tn, &conf)

	err = s.RegisterProtocolInstance(pi)
	if err != nil {
		return nil, err
	}

	go func(pname string) {
		if tmpErr := pi.Dispatch(); tmpErr != nil {
			log.Error("Error running Dispatch ->" + name + " :" + err.Error())
		}
	}(name)
	go func(pname string) {
		if tmpErr := pi.Start(); tmpErr != nil {
			log.Error("Error running Start ->" + name + " :" + err.Error())
		}
	}(name)

	return pi, err
}

// Service Phases
//______________________________________________________________________________________________________________________

// TaggingPhase performs the private grouping on the currently collected data.
func (s *Service) TaggingPhase(targetSurvey SurveyID, roster *onet.Roster) ([]libunlynx.DeterministCipherText, error) {
	start := time.Now()
	pi, err := s.StartProtocol(protocolsunlynx.DeterministicTaggingProtocolName, targetSurvey, roster)
	if err != nil {
		return nil, err
	}

	surveyTag, err := s.getSurveyTag(targetSurvey)
	if err != nil {
		return nil, err
	}
	surveyTag.TR.DDTRequestTimeExec += time.Since(start)
	err = s.putSurveyTag(surveyTag.SurveyID, surveyTag)
	if err != nil {
		return nil, err
	}

	deterministicTaggingResult := <-pi.(*protocolsunlynx.DeterministicTaggingProtocol).FeedbackChannel

	surveyTag, err = s.getSurveyTag(targetSurvey)
	if err != nil {
		return nil, err
	}
	surveyTag.TR.DDTRequestTimeExec += pi.(*protocolsunlynx.DeterministicTaggingProtocol).ExecTime
	surveyTag.TR.DDTRequestTimeCommunication = time.Since(start) - surveyTag.TR.DDTRequestTimeExec
	err = s.putSurveyTag(surveyTag.SurveyID, surveyTag)
	if err != nil {
		return nil, err
	}

	return deterministicTaggingResult, nil
}

// ShufflingPhase performs the shuffling aggregated results from each of the nodes
func (s *Service) ShufflingPhase(targetSurvey SurveyID, roster *onet.Roster) ([]libunlynx.CipherVector, error) {
	start := time.Now()
	pi, err := s.StartProtocol(protocolsunlynx.ShufflingProtocolName, targetSurvey, roster)
	if err != nil {
		return nil, err
	}
	shufflingTimeExec := time.Since(start)

	shufflingResult := <-pi.(*protocolsunlynx.ShufflingProtocol).FeedbackChannel

	shufflingTimeExec += pi.(*protocolsunlynx.ShufflingProtocol).ExecTime
	shufflingTimeCommun := time.Since(start) - shufflingTimeExec

	if shufflingTimeCommun < 0 {
		shufflingTimeCommun = 0
	}

	surveyAgg, err := s.getSurveyAgg(targetSurvey)
	if err != nil {
		return nil, err
	}
	surveyAgg.TR.AggRequestTimeExec += shufflingTimeExec
	surveyAgg.TR.AggRequestTimeCommunication += shufflingTimeCommun
	err = s.putSurveyAgg(surveyAgg.SurveyID, surveyAgg)
	if err != nil {
		return nil, err
	}

	return shufflingResult, nil
}

// KeySwitchingPhase performs the switch to the querier key on the currently aggregated data.
func (s *Service) KeySwitchingPhase(targetSurvey SurveyID, roster *onet.Roster) (libunlynx.CipherVector, error) {
	start := time.Now()
	pi, err := s.StartProtocol(protocolsunlynx.KeySwitchingProtocolName, targetSurvey, roster)
	if err != nil {
		return nil, err
	}
	keySwitchedAggregatedResponses := <-pi.(*protocolsunlynx.KeySwitchingProtocol).FeedbackChannel

	// *(nbr of servers) because this protocol happens sequentially
	keySTimeExec := pi.(*protocolsunlynx.KeySwitchingProtocol).ExecTime * time.Duration(len(roster.List))
	keySTimeCommun := time.Since(start) - keySTimeExec

	if keySTimeCommun < 0 {
		keySTimeCommun = 0
	}

	surveyAgg, err := s.getSurveyAgg(targetSurvey)
	if err != nil {
		return nil, err
	}
	surveyAgg.TR.AggRequestTimeExec += keySTimeExec
	surveyAgg.TR.AggRequestTimeCommunication += keySTimeCommun
	err = s.putSurveyAgg(surveyAgg.SurveyID, surveyAgg)
	if err != nil {
		return nil, err
	}

	return keySwitchedAggregatedResponses, nil
}

// Support functions
//______________________________________________________________________________________________________________________

type secretDDT struct {
	ServerID string
	Secret   string
}

type privateTOML struct {
	Public      string
	Private     string
	Address     string
	Description string
	Secrets     []secretDDT
}

func createTOMLSecrets(path string, id network.Address) (kyber.Scalar, error) {
	var fileHandle *os.File
	var err error
	defer fileHandle.Close()

	fileHandle, err = os.Create(path)

	encoder := toml.NewEncoder(fileHandle)

	secret := libunlynx.SuiTe.Scalar().Pick(random.New())
	b, err := secret.MarshalBinary()

	if err != nil {
		return nil, err
	}

	aux := make([]secretDDT, 0)
	aux = append(aux, secretDDT{ServerID: id.String(), Secret: base64.StdEncoding.EncodeToString(b)})
	endR := privateTOML{Public: "", Private: "", Address: "", Description: "", Secrets: aux}

	err = encoder.Encode(&endR)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func addTOMLSecret(path string, content privateTOML) error {
	var fileHandle *os.File
	defer fileHandle.Close()

	fileHandle, err := os.Create(path)

	encoder := toml.NewEncoder(fileHandle)

	err = encoder.Encode(&content)
	if err != nil {
		return err
	}

	return nil
}

// CheckDDTSecrets checks for the existence of the DDT secrets on the private_*.toml (we need to ensure that we use the same secrets always)
func CheckDDTSecrets(path string, id network.Address) (kyber.Scalar, error) {
	var err error

	if _, err = os.Stat(path); os.IsNotExist(err) {
		return createTOMLSecrets(path, id)
	}

	contents := privateTOML{}
	if _, err := toml.DecodeFile(path, &contents); err != nil {
		return nil, err
	}

	for _, el := range contents.Secrets {
		if el.ServerID == id.String() {
			secret := libunlynx.SuiTe.Scalar()

			b, err := base64.StdEncoding.DecodeString(el.Secret)
			if err != nil {
				return nil, err
			}

			err = secret.UnmarshalBinary(b)
			if err != nil {
				return nil, err
			}

			return secret, nil
		}
	}

	// no secret for this 'source' server
	secret := libunlynx.SuiTe.Scalar().Pick(random.New())
	b, err := secret.MarshalBinary()

	if err != nil {
		return nil, err
	}

	contents.Secrets = append(contents.Secrets, secretDDT{ServerID: id.String(), Secret: base64.StdEncoding.EncodeToString(b)})

	err = addTOMLSecret(path, contents)
	if err != nil {
		return nil, err
	}

	return secret, nil
}
