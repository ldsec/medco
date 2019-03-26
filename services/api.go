package servicesmedco

import (
	"github.com/lca1/unlynx/lib"
	"github.com/satori/go.uuid"
	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/util/key"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/onet/v3/log"
	"go.dedis.ch/onet/v3/network"
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
	log.Lvl1("Client", c.ClientID, "is creating a DDT survey with ID:", surveyID)

	rndUUID := uuid.NewV4()
	sdq := SurveyDDTRequest{
		SurveyID: SurveyID(rndUUID.String()),
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
		return nil, resp.Result, TimeResults{}, err
	}
	return &surveyID, resp.Result, resp.TR, nil
}

// SendSurveyAggRequest sends the encrypted aggregate local results at each node and expects a shuffling and a key switching of these data.
func (c *API) SendSurveyAggRequest(entities *onet.Roster, surveyID SurveyID, cPK kyber.Point, aggregate libunlynx.CipherText, proofs bool) (*SurveyID, libunlynx.CipherText, TimeResults, error) {
	log.Lvl1("Client", c.ClientID, "is creating a Agg survey with ID:", surveyID)

	listAggregate := make([]libunlynx.CipherText, 0)
	listAggregate = append(listAggregate, aggregate)

	sar := SurveyAggRequest{
		SurveyID:     surveyID,
		Roster:       *entities,
		Proofs:       proofs,
		ClientPubKey: cPK,

		Aggregate:         listAggregate,
		AggregateShuffled: make(libunlynx.CipherVector, 0),

		IntraMessage: false,
	}

	resp := ResultAgg{}
	err := c.SendProtobuf(c.entryPoint, &sar, &resp)
	if err != nil {

		return nil, resp.Result, TimeResults{}, err
	}
	return &surveyID, resp.Result, resp.TR, nil
}
