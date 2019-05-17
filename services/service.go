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
	"strings"
	"sync"
	"time"
)

// MsgTypes defines the Message Type ID for all the service's intra-messages.
type MsgTypes struct {
	msgSurveyDDTRequestTerms  network.MessageTypeID
	msgSurveyTagGenerated     network.MessageTypeID
	msgSurveyKSRequest        network.MessageTypeID
	msgSurveyShuffleRequest   network.MessageTypeID
	msgSurveyShuffleGenerated network.MessageTypeID
	msgSurveyAggRequest       network.MessageTypeID
	msgSurveyAggGenerated     network.MessageTypeID
}

var msgTypes = MsgTypes{}

func init() {
	_, err := onet.RegisterNewService(Name, NewService)
	log.ErrFatal(err)

	// messages for DDT Request
	msgTypes.msgSurveyDDTRequestTerms = network.RegisterMessage(&SurveyDDTRequest{})
	msgTypes.msgSurveyTagGenerated = network.RegisterMessage(&SurveyTagGenerated{})
	network.RegisterMessage(&ResultDDT{})

	// messages for the other requests
	msgTypes.msgSurveyKSRequest = network.RegisterMessage(&SurveyKSRequest{})
	msgTypes.msgSurveyShuffleRequest = network.RegisterMessage(&SurveyShuffleRequest{})
	msgTypes.msgSurveyShuffleGenerated = network.RegisterMessage(&SurveyShuffleGenerated{})
	msgTypes.msgSurveyAggRequest = network.RegisterMessage(&SurveyAggRequest{})
	msgTypes.msgSurveyAggGenerated = network.RegisterMessage(&SurveyAggGenerated{})
	network.RegisterMessage(&Result{})
}

// Service defines a service in unlynx
type Service struct {
	*onet.ServiceProcessor

	MapSurveyTag     *concurrent.ConcurrentMap
	MapSurveyKS      *concurrent.ConcurrentMap
	MapSurveyShuffle *concurrent.ConcurrentMap
	MapSurveyAgg     *concurrent.ConcurrentMap
	Mutex            *sync.Mutex
}

// NewService constructor which registers the needed messages.
func NewService(c *onet.Context) (onet.Service, error) {
	newUnLynxInstance := &Service{
		ServiceProcessor: onet.NewServiceProcessor(c),
		MapSurveyTag:     concurrent.NewConcurrentMap(),
		MapSurveyKS:      concurrent.NewConcurrentMap(),
		MapSurveyShuffle: concurrent.NewConcurrentMap(),
		MapSurveyAgg:     concurrent.NewConcurrentMap(),
		Mutex:            &sync.Mutex{},
	}

	if cerr := newUnLynxInstance.RegisterHandler(newUnLynxInstance.HandleSurveyDDTRequestTerms); cerr != nil {
		log.Error("Wrong Handler.", cerr)
	}
	if cerr := newUnLynxInstance.RegisterHandler(newUnLynxInstance.HandleSurveyKSRequest); cerr != nil {
		log.Error("Wrong Handler.", cerr)
	}
	if cerr := newUnLynxInstance.RegisterHandler(newUnLynxInstance.HandleSurveyShuffleRequest); cerr != nil {
		log.Error("Wrong Handler.", cerr)
	}
	if cerr := newUnLynxInstance.RegisterHandler(newUnLynxInstance.HandleSurveyAggRequest); cerr != nil {
		log.Error("Wrong Handler.", cerr)
	}

	c.RegisterProcessor(newUnLynxInstance, msgTypes.msgSurveyDDTRequestTerms)
	c.RegisterProcessor(newUnLynxInstance, msgTypes.msgSurveyTagGenerated)

	c.RegisterProcessor(newUnLynxInstance, msgTypes.msgSurveyKSRequest)

	c.RegisterProcessor(newUnLynxInstance, msgTypes.msgSurveyShuffleRequest)
	c.RegisterProcessor(newUnLynxInstance, msgTypes.msgSurveyShuffleGenerated)

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
	} else if msg.MsgType.Equal(msgTypes.msgSurveyKSRequest) {
		tmp := (msg.Msg).(*SurveyKSRequest)
		_, err := s.HandleSurveyKSRequest(tmp)
		if err != nil {
			log.Error(err)
		}
	} else if msg.MsgType.Equal(msgTypes.msgSurveyShuffleRequest) {
		tmp := (msg.Msg).(*SurveyShuffleRequest)
		_, err := s.HandleSurveyShuffleRequest(tmp)
		if err != nil {
			log.Error(err)
		}
	} else if msg.MsgType.Equal(msgTypes.msgSurveyShuffleGenerated) {
		tmp := (msg.Msg).(*SurveyShuffleGenerated)
		_, err := s.HandleSurveyShuffleGenerated(tmp)
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

		// convert the result to of the tagging for something close to the XML response of i2b2 (array of tagged terms)
		listTaggedTerms := make([]libunlynx.GroupingKey, 0)

		for _, el := range deterministicTaggingResult {
			listTaggedTerms = append(listTaggedTerms, libunlynx.GroupingKey(el.String()))
		}

		surveyTag, err = s.getSurveyTag(sdq.SurveyID)
		if err != nil {
			return nil, err
		}
		return &ResultDDT{Result: listTaggedTerms}, nil
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

// HandleSurveyKSRequest handles the reception of the aggregate local result to be key /switched
func (s *Service) HandleSurveyKSRequest(skr *SurveyKSRequest) (network.Message, error) {
	log.Lvl1(s.ServerIdentity().String(), " received a SurveyKSRequest:", skr.SurveyID)

	err := s.putSurveyKS(skr.SurveyID, SurveyKS{
		SurveyID: skr.SurveyID,
		Request:  *skr})

	// key switch the results
	keySwitchingResult, err := s.KeySwitchingPhase(skr.SurveyID, KSRequestName, &skr.Roster)
	if err != nil {
		log.Error("key switching error", err)
		return nil, err
	}

	return &Result{Result: keySwitchingResult}, nil
}

// HandleSurveyShuffleRequest handles the reception of the aggregate local result to be shared/shuffled/switched
func (s *Service) HandleSurveyShuffleRequest(ssr *SurveyShuffleRequest) (network.Message, error) {
	var root bool
	if s.ServerIdentity().String() == ssr.Roster.List[0].String() {
		root = true
	} else {
		root = false
	}

	log.Lvl1(s.ServerIdentity().String(), " received a SurveyShuffleRequest:", ssr.SurveyID, "(root =", root, "- intra =", ssr.IntraMessage, ")")

	// (root = true - intra = false )
	if !ssr.IntraMessage && root {

		err := s.putSurveyShuffle(ssr.SurveyID, SurveyShuffle{
			SurveyID:      ssr.SurveyID,
			Request:       *ssr,
			SurveyChannel: make(chan int, 100)})

		// send signal to unlock the other nodes
		err = libunlynxtools.SendISMOthers(s.ServiceProcessor, &ssr.Roster, &SurveyShuffleGenerated{SurveyID: ssr.SurveyID})
		if err != nil {
			log.Error("broadcasting error ", err)
		}

		surveyShuffle, err := s.getSurveyShuffle(ssr.SurveyID)
		if err != nil {
			return nil, err
		}

		// wait until you've got all the aggregate results from the other nodes
		counter := len(ssr.Roster.List) - 1
		for counter > 0 {
			counter = counter - <-surveyShuffle.SurveyChannel
		}

		surveyShuffle, err = s.getSurveyShuffle(ssr.SurveyID)
		if err != nil {
			return nil, err
		}

		if len(surveyShuffle.Request.ShuffleTarget) <= 1 {
			return nil, errors.New("no data to shuffle")
		}

		// shuffle the results
		shufflingResult, err := s.ShufflingPhase(ssr.SurveyID, &ssr.Roster)

		if err != nil {
			log.Error("shuffling error", err)
			return nil, err
		}

		shufflingFinalResult := make(libunlynx.CipherVector, 0)
		for _, el := range shufflingResult {
			shufflingFinalResult = append(shufflingFinalResult, el[0])
		}

		surveyShuffle.Request.KSTarget = shufflingFinalResult

		err = s.putSurveyShuffle(ssr.SurveyID, surveyShuffle)
		if err != nil {
			return nil, err
		}

		// send the shuffled results to all the other nodes
		ssr.KSTarget = shufflingFinalResult
		ssr.IntraMessage = true
		ssr.MessageSource = s.ServerIdentity()

		// let's delete what we don't need (less communication time)
		ssr.ShuffleTarget = nil

		// signal the other nodes that they need to prepare to execute a key switching
		// basically after shuffling the results the root server needs to send them back
		// to the remaining nodes for key switching
		err = libunlynxtools.SendISMOthers(s.ServiceProcessor, &ssr.Roster, ssr)
		if err != nil {
			log.Error("broadcasting error ", err)
		}

		// key switch the results
		keySwitchingResult, err := s.KeySwitchingPhase(ssr.SurveyID, ShuffleRequestName, &ssr.Roster)
		if err != nil {
			log.Error("key switching error", err)
			return nil, err
		}

		// get server index
		index := 0
		for i, r := range ssr.Roster.List {
			if r.String() == s.ServerIdentity().String() {
				index = i
				break
			}
		}

		surveyShuffle, err = s.getSurveyShuffle(ssr.SurveyID)
		if err != nil {
			return nil, err
		}

		return &Result{Result: libunlynx.CipherVector{keySwitchingResult[index]}}, nil

		//(root = false - intra = false )
	} else if !ssr.IntraMessage && !root { // if message sent by client and not a root node
		err := s.putSurveyShuffle(ssr.SurveyID, SurveyShuffle{
			SurveyID:            ssr.SurveyID,
			Request:             *ssr,
			SurveyChannel:       make(chan int, 100),
			FinalResultsChannel: make(chan int, 100),
		})

		ssr.IntraMessage = true
		ssr.MessageSource = s.ServerIdentity()

		surveyShuffle, err := s.getSurveyShuffle(ssr.SurveyID)
		if err != nil {
			return nil, err
		}
		<-surveyShuffle.SurveyChannel

		// send your local aggregate result to the root server (index 0)
		err = s.SendRaw(ssr.Roster.List[0], ssr)
		if err != nil {
			log.Error(s.ServerIdentity().String()+"could not send its aggregate value", err)
		}

		surveyShuffle, err = s.getSurveyShuffle(ssr.SurveyID)
		if err != nil {
			return nil, err
		}
		//waits for the final results to be ready
		<-surveyShuffle.FinalResultsChannel

		// get server index
		index := 0
		for i, r := range ssr.Roster.List {
			if r.String() == s.ServerIdentity().String() {
				index = i
				break
			}
		}

		surveyShuffle, err = s.getSurveyShuffle(ssr.SurveyID)
		if err != nil {
			return nil, err
		}

		return &Result{Result: libunlynx.CipherVector{surveyShuffle.Request.KSTarget[index]}}, nil

		// (root = true - intra = true )
	} else if ssr.IntraMessage && root { // if message sent by another node and is root
		s.Mutex.Lock()
		surveyShuffle, err := s.getSurveyShuffle(ssr.SurveyID)
		if err != nil {
			return nil, err
		}
		surveyShuffle.Request.ShuffleTarget = append(surveyShuffle.Request.ShuffleTarget, ssr.ShuffleTarget...)
		err = s.putSurveyShuffle(ssr.SurveyID, surveyShuffle)
		if err != nil {
			return nil, err
		}
		s.Mutex.Unlock()

		// get the request from the other non-root nodes
		surveyShuffle, err = s.getSurveyShuffle(ssr.SurveyID)
		if err != nil {
			return nil, err
		}
		surveyShuffle.SurveyChannel <- 1
		// (root = false - intra = true )
	} else { // if message sent by another node and not root
		// update the local survey with the shuffled results
		s.Mutex.Lock()
		surveyShuffle, err := s.getSurveyShuffle(ssr.SurveyID)
		if err != nil {
			return nil, err
		}
		surveyShuffle.Request.KSTarget = ssr.KSTarget
		err = s.putSurveyShuffle(ssr.SurveyID, surveyShuffle)
		if err != nil {
			return nil, err
		}
		s.Mutex.Unlock()

		// key switch the results
		keySwitchingResult, err := s.KeySwitchingPhase(ssr.SurveyID, ShuffleRequestName, &ssr.Roster)
		if err != nil {
			log.Error("key switching error", err)
			return nil, err
		}

		s.Mutex.Lock()
		surveyShuffle, err = s.getSurveyShuffle(ssr.SurveyID)
		if err != nil {
			return nil, err
		}
		surveyShuffle.Request.KSTarget = keySwitchingResult
		err = s.putSurveyShuffle(ssr.SurveyID, surveyShuffle)
		if err != nil {
			return nil, err
		}
		s.Mutex.Unlock()
		surveyShuffle.FinalResultsChannel <- 1
	}
	return nil, nil
}

// HandleSurveyShuffleGenerated handles triggers the SurveyChannel
func (s *Service) HandleSurveyShuffleGenerated(recq *SurveyShuffleGenerated) (network.Message, error) {
	var el interface{}
	el = nil
	for el == nil {
		el, _ = s.MapSurveyShuffle.Get((string)(recq.SurveyID))
		if el != nil {
			break
		}
		time.Sleep(time.Millisecond * 100)
	}

	surveyShuffle, err := s.getSurveyShuffle(recq.SurveyID)
	if err != nil {
		return nil, err
	}
	surveyShuffle.SurveyChannel <- 1
	return nil, nil
}

// HandleSurveyAggRequest handles the reception of the aggregate local result to be shared/shuffled/switched
func (s *Service) HandleSurveyAggRequest(sar *SurveyAggRequest) (network.Message, error) {
	log.Lvl1(s.ServerIdentity().String(), " received a SurveyAggRequest:", sar.SurveyID)

	err := s.putSurveyAgg(sar.SurveyID, SurveyAgg{
		SurveyID:      sar.SurveyID,
		Request:       *sar,
		SurveyChannel: make(chan int, 100)})

	if err != nil {
		return nil, err
	}

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

	// collectively aggregate the results
	aggregationResult, err := s.CollectiveAggregationPhase(sar.SurveyID, &sar.Roster)
	if err != nil {
		log.Error("aggregation error", err)
		return nil, err
	}

	surveyAgg.Request.KSTarget = aggregationResult

	err = s.putSurveyAgg(sar.SurveyID, surveyAgg)
	if err != nil {
		return nil, err
	}

	// key switch the results
	keySwitchingResult, err := s.KeySwitchingPhase(sar.SurveyID, AggRequestName, &sar.Roster)

	if err != nil {
		log.Error("key switching error", err)
		return nil, err
	}

	return &Result{Result: keySwitchingResult}, nil
}

// HandleSurveyAggGenerated handles triggers the SurveyChannel
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

// Protocol Handlers
//______________________________________________________________________________________________________________________

// whatRequest fetches the data from the correct map based on the configuration string ('target')
func (s *Service) whatRequest(target string) (bool, libunlynx.CipherVector, kyber.Point, error) {
	var proofs bool
	var data libunlynx.CipherVector
	var cPubKey kyber.Point

	tokens := strings.Split(target, "/")
	sID := SurveyID(tokens[0])
	typeQ := tokens[1]

	switch typeQ {
	case KSRequestName:
		surveyKS, err := s.getSurveyKS(sID)
		if err != nil {
			return false, nil, nil, err
		}
		proofs = surveyKS.Request.Proofs
		data = surveyKS.Request.KSTarget
		cPubKey = surveyKS.Request.ClientPubKey

	case ShuffleRequestName:
		surveyShuffle, err := s.getSurveyShuffle(sID)
		if err != nil {
			return false, nil, nil, err
		}
		proofs = surveyShuffle.Request.Proofs
		data = surveyShuffle.Request.KSTarget
		cPubKey = surveyShuffle.Request.ClientPubKey

	case AggRequestName:
		surveyAgg, err := s.getSurveyAgg(sID)
		if err != nil {
			return false, nil, nil, err
		}
		proofs = surveyAgg.Request.Proofs
		data = libunlynx.CipherVector{surveyAgg.Request.KSTarget}
		cPubKey = surveyAgg.Request.ClientPubKey

	default:
		return false, nil, nil, errors.New("Could not identify request:" + typeQ)
	}
	return proofs, data, cPubKey, nil
}

// NewProtocol creates a protocol instance executed by all nodes
func (s *Service) NewProtocol(tn *onet.TreeNodeInstance, conf *onet.GenericConfig) (onet.ProtocolInstance, error) {
	if err := tn.SetConfig(conf); err != nil {
		return nil, err
	}

	var pi onet.ProtocolInstance
	var err error
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
		surveyShuffle, err := s.getSurveyShuffle(target)
		if err != nil {
			return nil, err
		}

		pi, err = protocolsunlynx.NewShufflingProtocol(tn)
		if err != nil {
			return nil, err
		}

		shuffle := pi.(*protocolsunlynx.ShufflingProtocol)

		shuffle.Proofs = surveyShuffle.Request.Proofs
		shuffle.Precomputed = nil

		if tn.IsRoot() {
			dataToShuffle := protocolsunlynx.AdaptCipherTextArray(surveyShuffle.Request.ShuffleTarget)
			shuffle.ShuffleTarget = &dataToShuffle
		}
	case protocolsunlynx.KeySwitchingProtocolName:
		pi, err = protocolsunlynx.NewKeySwitchingProtocol(tn)
		if err != nil || pi == nil {
			return nil, err
		}

		keySwitch := pi.(*protocolsunlynx.KeySwitchingProtocol)

		if tn.IsRoot() {
			//define which map to retrieve the values to key switch
			proofs, data, cPubKey, err := s.whatRequest(string(target))
			if err != nil {
				return nil, err
			}

			keySwitch.Proofs = proofs
			dataToSwitch := data
			keySwitch.TargetOfSwitch = &dataToSwitch
			tmp := cPubKey
			keySwitch.TargetPublicKey = &tmp
		}
	case protocolsunlynx.CollectiveAggregationProtocolName:
		surveyAgg, err := s.getSurveyAgg(target)
		if err != nil {
			return nil, err
		}

		pi, err = protocolsunlynx.NewCollectiveAggregationProtocol(tn)
		if err != nil {
			return nil, err
		}

		aggr := pi.(*protocolsunlynx.CollectiveAggregationProtocol)
		aggr.Proofs = surveyAgg.Request.Proofs

		data := make([]libunlynx.CipherText, 0)
		data = append(data, surveyAgg.Request.AggregateTarget)
		aggr.SimpleData = &data
	default:
		return nil, errors.New("Service attempts to start an unknown protocol: " + tn.ProtocolName() + ".")
	}

	return pi, nil
}

// StartProtocol starts a specific protocol (Shuffling, KeySwitching, etc.)
func (s *Service) StartProtocol(name, typeQ string, targetSurvey SurveyID, roster *onet.Roster) (onet.ProtocolInstance, error) {
	tree := roster.GenerateNaryTreeWithRoot(2, s.ServerIdentity())
	tn := s.NewTreeNodeInstance(tree, tree.Root, name)

	var confData string
	if name == protocolsunlynx.KeySwitchingProtocolName {
		confData = string(targetSurvey) + "/" + typeQ
	} else {
		confData = string(targetSurvey)
	}

	conf := onet.GenericConfig{Data: []byte(confData)}
	pi, err := s.NewProtocol(tn, &conf)
	if err != nil || pi == nil {
		return nil, err
	}

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
	pi, err := s.StartProtocol(protocolsunlynx.DeterministicTaggingProtocolName, "", targetSurvey, roster)
	if err != nil {
		return nil, err
	}
	deterministicTaggingResult := <-pi.(*protocolsunlynx.DeterministicTaggingProtocol).FeedbackChannel
	return deterministicTaggingResult, nil
}

// CollectiveAggregationPhase performs a collective aggregation between the participating nodes
func (s *Service) CollectiveAggregationPhase(targetSurvey SurveyID, roster *onet.Roster) (libunlynx.CipherText, error) {
	pi, err := s.StartProtocol(protocolsunlynx.CollectiveAggregationProtocolName, "", targetSurvey, roster)
	if err != nil {
		return libunlynx.CipherText{}, err
	}
	aggregationResult := <-pi.(*protocolsunlynx.CollectiveAggregationProtocol).FeedbackChannel

	// in the resulting map there is only one element
	var finalResult libunlynx.CipherText
	for _, v := range aggregationResult.GroupedData {
		finalResult = v.AggregatingAttributes[0]
		break
	}

	return finalResult, nil
}

// ShufflingPhase performs the shuffling aggregated results from each of the nodes
func (s *Service) ShufflingPhase(targetSurvey SurveyID, roster *onet.Roster) ([]libunlynx.CipherVector, error) {
	pi, err := s.StartProtocol(protocolsunlynx.ShufflingProtocolName, "", targetSurvey, roster)
	if err != nil {
		return nil, err
	}
	shufflingResult := <-pi.(*protocolsunlynx.ShufflingProtocol).FeedbackChannel
	return shufflingResult, nil
}

// KeySwitchingPhase performs the switch to the querier key on the currently aggregated data.
func (s *Service) KeySwitchingPhase(targetSurvey SurveyID, typeQ string, roster *onet.Roster) (libunlynx.CipherVector, error) {
	pi, err := s.StartProtocol(protocolsunlynx.KeySwitchingProtocolName, typeQ, targetSurvey, roster)
	if err != nil {
		return nil, err
	}
	keySwitchedAggregatedResponses := <-pi.(*protocolsunlynx.KeySwitchingProtocol).FeedbackChannel
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
