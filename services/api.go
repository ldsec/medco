package servicesmedco

import (
	"github.com/ldsec/unlynx/lib"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/util/key"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
	"time"
)

// API represents a client with the server to which he is connected and its public/private key pair.
type API struct {
	*onet.Client
	ClientID   string
	entryPoint *network.ServerIdentity
	public     kyber.Point
	private    kyber.Scalar
}

// NewMedCoClient constructor of a client.
func NewMedCoClient(entryPoint *network.ServerIdentity, clientID string) *API {
	keys := key.NewKeyPair(libunlynx.SuiTe)
	newClient := &API{
		Client:     onet.NewClient(libunlynx.SuiTe, Name),
		ClientID:   clientID,
		entryPoint: entryPoint,
		public:     keys.Public,
		private:    keys.Private,
	}
	return newClient
}

// Send Queries
//______________________________________________________________________________________________________________________

// SendSurveyDDTRequestTerms sends the encrypted query terms and DDT tags those terms (the array of terms is ordered).
func (c *API) SendSurveyDDTRequestTerms(entities *onet.Roster, surveyID SurveyID, terms libunlynx.CipherVector, proofs bool, testing bool) (*SurveyID, []libunlynx.GroupingKey, TimeResults, error) {
	start := time.Now()
	log.Lvl2("Client", c.ClientID, "is creating a DDT survey with ID:", surveyID)

	sdq := SurveyDDTRequest{
		SurveyID: surveyID,
		Roster:   *entities,
		Proofs:   proofs,
		Testing:  testing,

		// query parameters to DDT
		Terms: terms,

		IntraMessage: false,
	}

	resp := ResultDDT{}
	err := c.SendProtobuf(c.entryPoint, &sdq, &resp)
	if err != nil {
		return nil, nil, TimeResults{}, err
	}
	resp.TR.MapTR[DDTRequestTime] = time.Since(start)
	return &surveyID, resp.Result, resp.TR, nil
}

// SendSurveyKSRequest performs key switching in a list of values
func (c *API) SendSurveyKSRequest(entities *onet.Roster, surveyID SurveyID, cPK kyber.Point, values libunlynx.CipherVector, proofs bool) (*SurveyID, libunlynx.CipherVector, TimeResults, error) {
	start := time.Now()
	log.Lvl2("Client", c.ClientID, "is creating a KS survey with ID:", surveyID)

	skr := SurveyKSRequest{
		SurveyID:     surveyID,
		Roster:       *entities,
		Proofs:       proofs,
		ClientPubKey: cPK,
		KSTarget:     values,
	}

	resp := Result{}
	err := c.SendProtobuf(c.entryPoint, &skr, &resp)
	if err != nil {
		return nil, nil, TimeResults{}, err
	}
	resp.TR.MapTR[KSRequestTime] = time.Since(start)
	return &surveyID, resp.Result, resp.TR, nil
}

// SendSurveyShuffleRequest performs shuffling + key switching on a list of values
func (c *API) SendSurveyShuffleRequest(entities *onet.Roster, surveyID SurveyID, cPK kyber.Point, value *libunlynx.CipherText, proofs bool) (*SurveyID, libunlynx.CipherText, TimeResults, error) {
	start := time.Now()
	log.Lvl2("Client", c.ClientID, "is creating a Shuffle survey with ID:", surveyID)

	target := make(libunlynx.CipherVector, 0)
	if value != nil {
		target = append(target, *value)
	}
	ssr := SurveyShuffleRequest{
		SurveyID:      surveyID,
		Roster:        *entities,
		Proofs:        proofs,
		ClientPubKey:  cPK,
		ShuffleTarget: target,
	}

	resp := Result{}
	err := c.SendProtobuf(c.entryPoint, &ssr, &resp)
	if err != nil {
		return nil, libunlynx.CipherText{}, TimeResults{}, err
	}
	resp.TR.MapTR[ShuffleRequestTime] = time.Since(start)
	return &surveyID, resp.Result[0], resp.TR, nil
}

// SendSurveyAggRequest sends the encrypted aggregate local results at each node and aggregates these values (result is the same for all nodes)
func (c *API) SendSurveyAggRequest(entities *onet.Roster, surveyID SurveyID, cPK kyber.Point, value libunlynx.CipherText, proofs bool) (*SurveyID, libunlynx.CipherText, TimeResults, error) {
	start := time.Now()
	log.Lvl2("Client", c.ClientID, "is creating a Agg survey with ID:", surveyID)

	sar := SurveyAggRequest{
		SurveyID:        surveyID,
		Roster:          *entities,
		Proofs:          proofs,
		ClientPubKey:    cPK,
		AggregateTarget: value,
	}

	resp := Result{}
	err := c.SendProtobuf(c.entryPoint, &sar, &resp)
	if err != nil {

		return nil, libunlynx.CipherText{}, TimeResults{}, err
	}
	resp.TR.MapTR[AggrRequestTime] = time.Since(start)
	return &surveyID, resp.Result[0], resp.TR, nil

}
