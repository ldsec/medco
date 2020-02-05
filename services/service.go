package servicesmedco

import (
	"encoding/base64"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/btcsuite/goleveldb/leveldb/errors"
	"github.com/fanliao/go-concurrentMap"
	"github.com/ldsec/medco-unlynx/protocols"
	"github.com/ldsec/unlynx/lib"
	"github.com/ldsec/unlynx/protocols"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/util/random"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
	"golang.org/x/xerrors"
	"os"
	"strings"
	"sync"
	"time"
)

func init() {
	_, err := onet.RegisterNewService(Name, NewService)
	log.ErrFatal(err)

	// Register SurveyShuffleRequest for propagation-protocol
	network.RegisterMessage(&SurveyShuffleRequest{})
}

var propagateShuffleFromChildren = "PropShuffleFromChildren"
var propagateShuffleToChildren = "PropShuffleToChildren"

// Service defines a service in unlynx
type Service struct {
	*onet.ServiceProcessor

	shuffleGetData protocols.PropagationFunc
	shufflePutData protocols.PropagationFunc

	MapSurveyKS      *concurrent.ConcurrentMap
	MapSurveyShuffle *concurrent.ConcurrentMap
	MapSurveyAgg     *concurrent.ConcurrentMap
	Mutex            *sync.Mutex
}

// NewService constructor which registers the needed messages.
func NewService(c *onet.Context) (onet.Service, error) {
	newUnLynxInstance := &Service{
		ServiceProcessor: onet.NewServiceProcessor(c),
		MapSurveyKS:      concurrent.NewConcurrentMap(),
		MapSurveyShuffle: concurrent.NewConcurrentMap(),
		MapSurveyAgg:     concurrent.NewConcurrentMap(),
		Mutex:            &sync.Mutex{},
	}
	var err error
	newUnLynxInstance.shuffleGetData, err =
		protocols.NewPropagationFunc(newUnLynxInstance, propagateShuffleFromChildren, -1)
	if err != nil {
		return nil, fmt.Errorf("couldn't create propagation function: %+v", err)
	}
	newUnLynxInstance.shufflePutData, err =
		protocols.NewPropagationFunc(newUnLynxInstance, propagateShuffleToChildren, -1)
	if err != nil {
		return nil, fmt.Errorf("couldn't create propagation function: %+v", err)
	}

	if cerr := newUnLynxInstance.RegisterHandlers(
		newUnLynxInstance.HandleSurveyDDTRequestTerms,
		newUnLynxInstance.HandleSurveyKSRequest,
		newUnLynxInstance.HandleSurveyShuffleRequest,
		newUnLynxInstance.HandleSurveyAggRequest); cerr != nil {
		log.Error("Wrong Handler.", cerr)
		return nil, cerr
	}

	return newUnLynxInstance, nil
}

// HandleSurveyDDTRequestTerms handles the reception of the query terms to be deterministically tagged
func (s *Service) HandleSurveyDDTRequestTerms(sdq *SurveyDDTRequest) (network.Message, error) {
	// sanitize params
	if err := emptySurveyID(sdq.SurveyID); err != nil {
		return nil, xerrors.Errorf("%+v", err)
	}
	if err := emptyRoster(sdq.Roster); err != nil {
		return nil, xerrors.Errorf("%+v", err)
	}

	// if this server is the one receiving the request from the client
	log.Lvl2(s.ServerIdentity().String(), " received a SurveyDDTRequestTerms:", sdq.SurveyID)

	if len(sdq.Terms) == 0 {
		return nil, xerrors.Errorf(s.ServerIdentity().String() + " for survey" + string(sdq.SurveyID) + "has no data to det tag")
	}

	// initialize timers
	mapTR := make(map[string]time.Duration)

	request := SurveyDDTRequest{
		SurveyID:      sdq.SurveyID,
		Proofs:        sdq.Proofs,
		Testing:       sdq.Testing,
		Terms:         sdq.Terms,
		MessageSource: s.ServerIdentity(),
	}

	deterministicTaggingResult, execTime, communicationTime,
		err := s.TaggingPhase(&request, &sdq.Roster)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	// convert the result to of the tagging for something close to the response of i2b2 (array of tagged terms)
	listTaggedTerms := make([]libunlynx.GroupingKey, 0)
	for _, el := range deterministicTaggingResult {
		listTaggedTerms = append(listTaggedTerms, libunlynx.GroupingKey(el.String()))
	}

	mapTR[TaggingTimeExec] = execTime
	mapTR[TaggingTimeCommunication] = communicationTime
	return &ResultDDT{Result: listTaggedTerms, TR: mapTR}, nil
}

// HandleSurveyKSRequest handles the reception of the aggregate local result to be key switched
func (s *Service) HandleSurveyKSRequest(skr *SurveyKSRequest) (network.Message, error) {
	// sanitize params
	if err := emptySurveyID(skr.SurveyID); err != nil {
		return nil, xerrors.Errorf("%+v", err)
	}
	if err := emptyRoster(skr.Roster); err != nil {
		return nil, xerrors.Errorf("%+v", err)
	}
	if skr.ClientPubKey == nil {
		return nil, xerrors.Errorf("no target public key")
	}
	if skr.KSTarget == nil && len(skr.KSTarget) == 0 {
		return nil, xerrors.Errorf(s.ServerIdentity().String() + " for survey" + string(skr.SurveyID) + "has no data to key switch")
	}

	log.Lvl2(s.ServerIdentity().String(), " received a SurveyKSRequest:", skr.SurveyID)

	mapTR := make(map[string]time.Duration)
	err := s.putSurveyKS(skr.SurveyID, SurveyKS{
		SurveyID: skr.SurveyID,
		Request:  *skr,
		TR:       TimeResults{MapTR: mapTR},
	})
	if err != nil {
		return nil, xerrors.Errorf("%+v", err)
	}

	// key switch the results
	keySwitchingResult, execTime, communicationTime, err := s.KeySwitchingPhase(skr.SurveyID, KSRequestName, &skr.Roster)
	if err != nil {
		return nil, xerrors.Errorf("key switching error: %+v", err)
	}

	surveyKS, err := s.getSurveyKS(skr.SurveyID)
	if err != nil {
		return nil, xerrors.Errorf("%+v", err)
	}
	surveyKS.TR.MapTR[KSTimeExec] = execTime
	surveyKS.TR.MapTR[KSTimeCommunication] = communicationTime

	// remove query from map
	_, err = s.deleteSurveyKS(skr.SurveyID)
	if err != nil {
		return nil, xerrors.Errorf("%+v", err)
	}

	return &Result{Result: keySwitchingResult, TR: surveyKS.TR}, nil
}

// HandleSurveyShuffleRequest handles the reception of the aggregate local result to be shared/shuffled/switched
func (s *Service) HandleSurveyShuffleRequest(ssr *SurveyShuffleRequest) (network.Message, error) {
	// sanitize params
	if err := emptySurveyID(ssr.SurveyID); err != nil {
		return nil, xerrors.Errorf("%+v", err)
	}
	if err := emptyRoster(ssr.Roster); err != nil {
		return nil, xerrors.Errorf("%+v", err)
	}
	if ssr.ClientPubKey == nil {
		return nil, xerrors.Errorf("no target public key")
	}

	var root bool
	if s.ServerIdentity().String() == ssr.Roster.List[0].String() {
		root = true
	} else {
		root = false
	}

	log.Lvl2(s.ServerIdentity().String(), " received a SurveyShuffleRequest:", ssr.SurveyID, "(root =", root, ")")

	if root {
		//Message sent to root node:
		//1. collect encrypted data from children
		//2. run shuffling protocol
		//3. distributed shuffled data to children
		//4. start key-switching
		//5. return data to client

		if ssr.ShuffleTarget == nil || len(ssr.ShuffleTarget) == 0 {
			return nil, xerrors.Errorf(s.ServerIdentity().String() + " for survey" + string(ssr.SurveyID) + "has no data to shuffle")
		}

		childrenMsgs, err := s.shuffleGetData(&ssr.Roster,
			&ProtocolConfig{SurveyID: ssr.SurveyID}, 10*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("couldn't get children data: %+v", err)
		}

		mapTR := make(map[string]time.Duration)
		surveyShuffle := SurveyShuffle{
			SurveyID:      ssr.SurveyID,
			Request:       *ssr,
			SurveyChannel: make(chan int, 100),
			TR:            TimeResults{MapTR: mapTR},
		}
		for _, msg := range childrenMsgs {
			req, ok := msg.(*SurveyShuffleRequest)
			if !ok {
				return nil, xerrors.New("couldn't convert msg to Request")
			}
			surveyShuffle.Request.ShuffleTarget = append(surveyShuffle.
				Request.ShuffleTarget, req.ShuffleTarget...)
		}

		err = s.putSurveyShuffle(ssr.SurveyID, surveyShuffle)
		if err != nil {
			return nil, xerrors.Errorf("%+v", err)
		}

		// shuffle the results
		shufflingResult, execTime, communicationTime, err := s.ShufflingPhase(ssr.SurveyID, &ssr.Roster)
		if err != nil {
			return nil, xerrors.Errorf("shuffling error: %+v", err)
		}

		shufflingFinalResult := make(libunlynx.CipherVector, 0)
		for _, el := range shufflingResult {
			shufflingFinalResult = append(shufflingFinalResult, el[0])
		}

		surveyShuffle.Request.KSTarget = shufflingFinalResult
		surveyShuffle.TR.MapTR[ShuffleTimeExec] = execTime
		surveyShuffle.TR.MapTR[ShuffleTimeCommunication] = communicationTime

		err = s.putSurveyShuffle(ssr.SurveyID, surveyShuffle)
		if err != nil {
			return nil, xerrors.Errorf("%+v", err)
		}

		// send the shuffled results to all the other nodes
		ssr.KSTarget = shufflingFinalResult
		ssr.MessageSource = s.ServerIdentity()

		// let's delete what we don't need (less communication time)
		ssr.ShuffleTarget = nil

		// signal the other nodes that they need to prepare to execute a key switching
		// basically after shuffling the results the root server needs to send them back
		// to the remaining nodes for key switching
		_, err = s.shufflePutData(&ssr.Roster, ssr, 10*time.Minute)
		if err != nil {
			return nil, fmt.Errorf("couldn't send data to children: %+v", err)
		}

		// key switch the results
		keySwitchingResult, execTime, communicationTime, err := s.KeySwitchingPhase(ssr.SurveyID, ShuffleRequestName, &ssr.Roster)
		if err != nil {
			return nil, xerrors.Errorf("key switching error: %+v", err)
		}

		// get server index
		index, _ := ssr.Roster.Search(s.ServerIdentity().ID)
		if index < 0 {
			return nil, xerrors.New("couldn't find this node in the roster")
		}

		surveyShuffle.TR.MapTR[KSTimeExec] = execTime
		surveyShuffle.TR.MapTR[KSTimeCommunication] = communicationTime

		// remove query from map
		_, err = s.deleteSurveyShuffle(ssr.SurveyID)
		if err != nil {
			return nil, xerrors.Errorf("%+v", err)
		}

		return &Result{Result: libunlynx.CipherVector{keySwitchingResult[index]}, TR: surveyShuffle.TR}, nil

	}
	//if message sent to children node:
	//1. Send encrypted data to root node
	//2. participate in shuffling
	//3. receive shuffled data from root node
	//4. start key-switching
	//5. return data to client

	mapTR := make(map[string]time.Duration)
	surveyShuffle := SurveyShuffle{
		SurveyID:            ssr.SurveyID,
		Request:             *ssr,
		SurveyChannel:       make(chan int, 100),
		FinalResultsChannel: make(chan int, 100),
		TR:                  TimeResults{MapTR: mapTR},
	}
	err := s.putSurveyShuffle(ssr.SurveyID, surveyShuffle)
	if err != nil {
		return nil, xerrors.Errorf("%+v", err)
	}

	ssr.MessageSource = s.ServerIdentity()

	// wait for root to be ready to send the local aggregate result
	<-surveyShuffle.SurveyChannel

	// update the local survey with the shuffled results
	surveyShuffle, err = s.getSurveyShuffle(ssr.SurveyID)
	if err != nil {
		return nil, xerrors.Errorf("%+v", err)
	}

	// key switch the results
	keySwitchingResult, execTime, communicationTime, err := s.KeySwitchingPhase(ssr.SurveyID, ShuffleRequestName, &ssr.Roster)
	if err != nil {
		return nil, xerrors.Errorf("key switching error: %+v", err)
	}

	surveyShuffle.TR.MapTR[KSTimeExec] = execTime
	surveyShuffle.TR.MapTR[KSTimeCommunication] = communicationTime

	// get server index
	index, _ := ssr.Roster.Search(s.ServerIdentity().ID)
	if index < 0 {
		return nil, xerrors.New("couldn't find this node in the roster")
	}

	// remove query from map
	_, err = s.deleteSurveyShuffle(ssr.SurveyID)
	if err != nil {
		return nil, xerrors.Errorf("%+v", err)
	}

	return &Result{Result: libunlynx.CipherVector{keySwitchingResult[index]},
		TR: surveyShuffle.TR}, nil
}

// HandleSurveyAggRequest handles the reception of the aggregate local result to be shared/shuffled/switched
func (s *Service) HandleSurveyAggRequest(sar *SurveyAggRequest) (network.Message, error) {
	// sanitize params
	if err := emptyRoster(sar.Roster); err != nil {
		return nil, xerrors.Errorf("%+v", err)
	}
	if err := emptySurveyID(sar.SurveyID); err != nil {
		return nil, xerrors.Errorf("%+v", err)
	}
	if sar.AggregateTarget.K == nil || sar.AggregateTarget.C == nil {
		return nil, xerrors.Errorf(s.ServerIdentity().String() + " for survey" + string(sar.SurveyID) + "has no data to aggregate")
	}
	if sar.ClientPubKey == nil {
		return nil, xerrors.Errorf("no target public key")
	}

	log.Lvl2(s.ServerIdentity().String(), " received a SurveyAggRequest:", sar.SurveyID)

	mapTR := make(map[string]time.Duration)
	surveyAgg := SurveyAgg{
		SurveyID: sar.SurveyID,
		Request:  *sar,
		TR:       TimeResults{MapTR: mapTR},
	}
	err := s.putSurveyAgg(sar.SurveyID, surveyAgg)

	if err != nil {
		return nil, xerrors.Errorf("%+v", err)
	}

	// collectively aggregate the results
	aggregationResult, aggrTime, err := s.CollectiveAggregationPhase(sar.SurveyID, &sar.Roster)
	if err != nil {
		return nil, xerrors.Errorf("aggregation error: %+v", err)
	}

	surveyAgg.Request.KSTarget = aggregationResult
	surveyAgg.TR.MapTR[AggrTime] = aggrTime

	err = s.putSurveyAgg(sar.SurveyID, surveyAgg)
	if err != nil {
		return nil, xerrors.Errorf("%+v", err)
	}

	// key switch the results
	keySwitchingResult, execTime, communicationTime, err := s.KeySwitchingPhase(sar.SurveyID, AggRequestName, &sar.Roster)
	if err != nil {
		return nil, xerrors.Errorf("key switching error: %+v", err)
	}
	surveyAgg.TR.MapTR[KSTimeExec] = execTime
	surveyAgg.TR.MapTR[KSTimeCommunication] = communicationTime

	// remove query from map
	_, err = s.deleteSurveyAgg(sar.SurveyID)
	if err != nil {
		return nil, xerrors.Errorf("%+v", err)
	}

	return &Result{Result: keySwitchingResult, TR: surveyAgg.TR}, nil
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
		return false, nil, nil, errors.New("Could not identify the request:" + typeQ)
	}
	return proofs, data, cPubKey, nil
}

// NewProtocol creates a protocol instance executed by all nodes
func (s *Service) NewProtocol(tn *onet.TreeNodeInstance,
	conf *onet.GenericConfig) (onet.ProtocolInstance, error) {
	var pi onet.ProtocolInstance
	var err error
	var target SurveyID
	var protoConf ProtocolConfig

	if conf != nil && conf.Data != nil {
		protoConf, err = unmarshalProtocolConfig(conf.Data)
		if err != nil {
			return nil, err
		}
		target = protoConf.getTarget()
	}

	switch tn.ProtocolName() {
	case protocolsunlynx.DeterministicTaggingProtocolName:
		_, sti, err := network.Unmarshal(protoConf.Data, libunlynx.SuiTe)
		if err != nil {
			log.Fatal(err)
			return nil, fmt.Errorf("couldn't unmarshal: %+v", err)
		}
		surveyRequest := sti.(*SurveyDDTRequest)

		pi, err = protocolsunlynx.NewDeterministicTaggingProtocol(tn)
		if err != nil {
			return nil, err
		}
		hashCreation := pi.(*protocolsunlynx.DeterministicTaggingProtocol)

		var serverIDMap *network.ServerIdentity

		if tn.IsRoot() {
			dataToDDT := surveyRequest.Terms
			hashCreation.TargetOfSwitch = &dataToDDT

			surveyRequest.Terms = libunlynx.CipherVector{}

			pc, err := newProtocolConfig(surveyRequest.SurveyID, "",
				surveyRequest)
			if err != nil {
				return nil, fmt.Errorf("couldn't update protocolConfig: %+v",
					err)
			}
			newConfig, err := pc.getConfig()
			if err != nil {
				return nil, fmt.Errorf("couldn't set config again: %+v", err)
			}
			conf = &newConfig

			serverIDMap = s.ServerIdentity()
		} else {
			serverIDMap = surveyRequest.MessageSource
		}

		s.Mutex.Lock()
		var aux kyber.Scalar
		if surveyRequest.Testing {
			path := DDTSecretsPath + "_" + s.ServerIdentity().Address.Host() + ":" + s.ServerIdentity().Address.Port() + ".toml"
			aux, err = CheckDDTSecrets(path, serverIDMap.Address, nil)
			if err != nil || aux == nil {
				log.Fatal(err)
				return nil, errors.New("Error while reading the DDT secrets from file")
			}
		} else {
			aux, err = CheckDDTSecrets(os.Getenv("UNLYNX_DDT_SECRETS_FILE_PATH"), serverIDMap.Address, nil)
			if err != nil || aux == nil {
				log.Fatal(err)
				return nil, errors.New("Error while reading the DDT secrets from file")
			}
		}
		hashCreation.SurveySecretKey = &aux
		hashCreation.Proofs = surveyRequest.Proofs
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
		if err != nil {
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
		var surveyAgg SurveyAgg
		maxLoop := 10 * 60
		for i := 1; i <= maxLoop; i++ {
			surveyAgg, err = s.getSurveyAgg(target)
			if err != nil {
				log.Lvl3(s.ServerIdentity(), "Waiting for data to arrive for survey", target)
				if i == maxLoop {
					return nil, xerrors.New(
						"didn't get data within 10 minutes - aborting")
				}
				time.Sleep(1000 * time.Millisecond)
			} else {
				break
			}
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

	case propagateShuffleFromChildren:
		pi, err = protocols.NewPropagationProtocol(tn)
		if err != nil {
			return nil, xerrors.Errorf("couldn't create protocol: %+v", err)
		}
		prop := pi.(*protocols.Propagate)
		surveyIDChan := make(chan SurveyID)
		prop.RegisterOnDataToChildren(func(msg network.Message) error {
			pc, ok := msg.(*ProtocolConfig)
			if !ok {
				return errors.New("didn't get ProtocolConfig")
			}
			surveyIDChan <- pc.SurveyID
			return nil
		})
		surveyShuffleChan := make(chan SurveyShuffle)
		prop.RegisterOnDataToRoot(func() network.Message {
			ss := <-surveyShuffleChan
			return &ss.Request
		})
		go func() {
			surveyID := <-surveyIDChan
			for {
				surveyShuffle, err := s.getSurveyShuffle(surveyID)
				if err != nil {
					time.Sleep(100 * time.Millisecond)
				} else {
					surveyShuffleChan <- surveyShuffle
					break
				}
			}
		}()

	case propagateShuffleToChildren:
		pi, err = protocols.NewPropagationProtocol(tn)
		if err != nil {
			return nil, xerrors.Errorf("couldn't create new protocol: %+v",
				err)
		}
		prop := pi.(*protocols.Propagate)
		prop.RegisterOnDataToChildren(func(msg network.Message) error {
			ssr, ok := msg.(*SurveyShuffleRequest)
			if !ok {
				return xerrors.New("didn't receive SurveyShuffleRequest" +
					" message")
			}
			surveyShuffle, err := s.getSurveyShuffle(ssr.SurveyID)
			if err != nil {
				return xerrors.Errorf("couldn't get survey: %+v", err)
			}
			surveyShuffle.Request.KSTarget = ssr.KSTarget
			err = s.putSurveyShuffle(ssr.SurveyID, surveyShuffle)
			if err != nil {
				return xerrors.Errorf(
					"couldn't store new surveyShuffle: %+v", err)
			}
			surveyShuffle.SurveyChannel <- 1
			return nil
		})

	default:
		return nil, errors.New("Service attempts to start an unknown protocol: " + tn.ProtocolName() + ".")
	}

	if err := tn.SetConfig(conf); err != nil {
		return nil, xerrors.Errorf("couldn't set config: %+v", err)
	}

	return pi, nil
}

// StartProtocol starts a specific protocol (Shuffling, KeySwitching, etc.)
func (s *Service) StartProtocol(name, typeQ string, pc ProtocolConfig,
	roster *onet.Roster) (onet.ProtocolInstance, error) {
	tree := roster.GenerateNaryTreeWithRoot(2, s.ServerIdentity())
	tn := s.NewTreeNodeInstance(tree, tree.Root, name)

	if name == protocolsunlynx.KeySwitchingProtocolName {
		pc.TypeQ = typeQ
	}

	conf, err := pc.getConfig()
	if err != nil {
		return nil, fmt.Errorf("couldn't get config: %+v", err)
	}
	pi, err := s.NewProtocol(tn, &conf)
	if err != nil || pi == nil {
		return nil, fmt.Errorf("couldn't start new protocol: %+v", err)
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
func (s *Service) TaggingPhase(targetSurvey *SurveyDDTRequest,
	roster *onet.Roster) ([]libunlynx.DeterministCipherText, time.Duration, time.Duration, error) {
	start := time.Now()
	pc, err := newProtocolConfig(targetSurvey.SurveyID, "", targetSurvey)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("couldn't get protoConfig: %+v", err)
	}
	pi, err := s.StartProtocol(protocolsunlynx.
		DeterministicTaggingProtocolName, "", pc, roster)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("couldn't start protocol: %+v", err)
	}
	deterministicTaggingResult := <-pi.(*protocolsunlynx.DeterministicTaggingProtocol).FeedbackChannel

	execTime := pi.(*protocolsunlynx.DeterministicTaggingProtocol).ExecTime
	return deterministicTaggingResult, execTime, time.Since(start) - execTime, nil
}

// CollectiveAggregationPhase performs a collective aggregation between the participating nodes
func (s *Service) CollectiveAggregationPhase(targetSurvey SurveyID, roster *onet.Roster) (libunlynx.CipherText, time.Duration, error) {
	start := time.Now()
	pi, err := s.StartProtocol(protocolsunlynx.CollectiveAggregationProtocolName, "",
		ProtocolConfig{targetSurvey, "", nil}, roster)
	if err != nil {
		return libunlynx.CipherText{}, 0, err
	}
	aggregationResult := <-pi.(*protocolsunlynx.CollectiveAggregationProtocol).FeedbackChannel

	// in the resulting map there is only one element
	var finalResult libunlynx.CipherText
	for _, v := range aggregationResult.GroupedData {
		finalResult = v.AggregatingAttributes[0]
		break
	}
	return finalResult, time.Since(start), nil
}

// ShufflingPhase performs the shuffling aggregated results from each of the nodes
func (s *Service) ShufflingPhase(targetSurvey SurveyID, roster *onet.Roster) ([]libunlynx.CipherVector, time.Duration, time.Duration, error) {
	start := time.Now()
	pi, err := s.StartProtocol(protocolsunlynx.ShufflingProtocolName, "",
		ProtocolConfig{targetSurvey, "", nil}, roster)
	if err != nil {
		return nil, 0, 0, err
	}
	shufflingResult := <-pi.(*protocolsunlynx.ShufflingProtocol).FeedbackChannel

	execTime := pi.(*protocolsunlynx.ShufflingProtocol).ExecTime
	return shufflingResult, execTime, time.Since(start) - execTime, nil
}

// KeySwitchingPhase performs the switch to the querier key on the currently aggregated data.
func (s *Service) KeySwitchingPhase(targetSurvey SurveyID, typeQ string, roster *onet.Roster) (libunlynx.CipherVector, time.Duration, time.Duration, error) {
	start := time.Now()
	pi, err := s.StartProtocol(protocolsunlynx.KeySwitchingProtocolName, typeQ,
		ProtocolConfig{targetSurvey, "", nil}, roster)
	if err != nil {
		return nil, 0, 0, err
	}
	keySwitchedAggregatedResponses := <-pi.(*protocolsunlynx.KeySwitchingProtocol).FeedbackChannel

	execTime := pi.(*protocolsunlynx.KeySwitchingProtocol).ExecTime
	return keySwitchedAggregatedResponses, execTime, time.Since(start) - execTime, nil
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

func createTOMLSecrets(path string, id network.Address, secret kyber.Scalar) (kyber.Scalar, error) {
	var fileHandle *os.File
	var err error
	defer fileHandle.Close()

	fileHandle, err = os.Create(path)

	encoder := toml.NewEncoder(fileHandle)

	// generate random secret if not provided
	if secret == nil {
		secret = libunlynx.SuiTe.Scalar().Pick(random.New())
	}
	b, err := secret.MarshalBinary()
	if err != nil {
		return nil, err
	}

	aux := make([]secretDDT, 0)
	aux = append(aux, secretDDT{ServerID: id.String(), Secret: base64.URLEncoding.EncodeToString(b)})
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
func CheckDDTSecrets(path string, id network.Address, secret kyber.Scalar) (kyber.Scalar, error) {
	var err error

	if _, err = os.Stat(path); os.IsNotExist(err) {
		return createTOMLSecrets(path, id, secret)
	}

	contents := privateTOML{}
	if _, err := toml.DecodeFile(path, &contents); err != nil {
		return nil, err
	}

	for _, el := range contents.Secrets {
		if el.ServerID == id.String() {
			secret := libunlynx.SuiTe.Scalar()

			b, err := base64.URLEncoding.DecodeString(el.Secret)
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

	// no secret for this 'source' server; set to the provided one or generate a random one
	if secret == nil {
		secret = libunlynx.SuiTe.Scalar().Pick(random.New())
	}

	b, err := secret.MarshalBinary()

	if err != nil {
		return nil, err
	}

	contents.Secrets = append(contents.Secrets, secretDDT{ServerID: id.String(), Secret: base64.URLEncoding.EncodeToString(b)})

	err = addTOMLSecret(path, contents)
	if err != nil {
		return nil, err
	}

	return secret, nil
}

func emptySurveyID(id SurveyID) error {
	if id == "" {
		return errors.New("survey id is empty")
	}
	return nil
}

func emptyRoster(roster onet.Roster) error {
	if len(roster.List) == 0 {
		return errors.New("roster is empty")
	}
	return nil
}
