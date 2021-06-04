package medchain

import (
	"sync/atomic"

	"github.com/ldsec/medchain/contracts"
	"go.dedis.ch/cothority/v3/byzcoin"
	"go.dedis.ch/cothority/v3/byzcoin/bcadmin/lib"
	"go.dedis.ch/cothority/v3/darc"
	"go.dedis.ch/protobuf"
	"golang.org/x/xerrors"
)

type QueryStatus string

const (
	QuerySuccessStatus QueryStatus = contracts.QuerySuccessStatus
	QueryFailedStatus  QueryStatus = contracts.QueryFailedStatus
)

type Client struct {
	bcClient *byzcoin.Client
	signer   darc.Signer
	counter  uint64
}

func NewClient(path string) (*Client, error) {
	_, bc, err := lib.LoadConfig(path)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize byzcoin client: %v", err)
	}

	return &Client{bcClient: bc}, nil
}

func (c *Client) GetAuthorization(projectInstID byzcoin.InstanceID, userID, queryID, queryDef string) (*contracts.Authorization, byzcoin.InstanceID, error) {
	instr := byzcoin.Instruction{
		InstanceID: projectInstID,
		Spawn: &byzcoin.Spawn{
			ContractID: contracts.QueryContractID,
			Args: []byzcoin.Argument{{
				Name:  contracts.QueryUserIDKey,
				Value: []byte(userID),
			}, {
				Name:  contracts.QueryQueryIDKey,
				Value: []byte(queryID),
			}, {
				Name:  contracts.QueryQueryDefinitionKey,
				Value: []byte(queryDef),
			}},
		},
		SignerCounter: []uint64{atomic.AddUint64(&c.counter, 1)},
	}

	ctx, err := c.bcClient.CreateTransaction(instr)
	if err != nil {
		return nil, byzcoin.InstanceID{}, xerrors.Errorf("failed to create transaction: %v", err)
	}

	err = ctx.FillSignersAndSignWith(c.signer)
	if err != nil {
		return nil, byzcoin.InstanceID{}, xerrors.Errorf("failed to sign transaction: %v", err)
	}

	_, err = c.bcClient.AddTransactionAndWait(ctx, 10)
	if err != nil {
		return nil, byzcoin.InstanceID{}, xerrors.Errorf("failed to add transaction: %v", err)
	}

	queryInstID := ctx.Instructions[0].DeriveID("")

	resp, err := c.bcClient.GetProofFromLatest(queryInstID[:])
	if err != nil {
		return nil, byzcoin.InstanceID{}, xerrors.Errorf("failed to get proof: %v", err)
	}

	_, val, _, _, _ := resp.Proof.KeyValue()
	query := contracts.QueryContract{}
	err = protobuf.Decode(val, &query)
	if err != nil {
		return nil, byzcoin.InstanceID{}, xerrors.Errorf("failed to decode query contract: %v", err)
	}

	if query.Status == contracts.QueryRejectedStatus {
		return nil, byzcoin.InstanceID{}, xerrors.Errorf("failed to authorize query: status %s", contracts.QueryRejectedStatus)
	}

	projectResponse, err := c.bcClient.GetProofFromLatest(projectInstID[:])
	if err != nil {
		return nil, byzcoin.InstanceID{}, xerrors.Errorf("failed to get project proof: %v", err)
	}

	_, val, _, _, _ = projectResponse.Proof.KeyValue()
	project := contracts.ProjectContract{}
	err = protobuf.Decode(val, &project)
	if err != nil {
		return nil, byzcoin.InstanceID{}, xerrors.Errorf("failed to decode project contract: %v", err)
	}

	return project.Authorizations.Find(userID), queryInstID, nil
}

func (c *Client) UpdateStatus(queryID byzcoin.InstanceID, status QueryStatus) error {
	instruction := byzcoin.Instruction{
		// that's the key part, where we provide the instanceID of the project
		// instance we just spawned. This instance will spawn the query.
		InstanceID: queryID,
		Invoke: &byzcoin.Invoke{
			Command:    contracts.QueryUpdateAction,
			ContractID: contracts.QueryContractID,
			Args: []byzcoin.Argument{{
				Name:  contracts.QueryStatusKey,
				Value: []byte(status),
			}},
		},
		SignerCounter: []uint64{atomic.AddUint64(&c.counter, 1)},
	}

	ctx, err := c.bcClient.CreateTransaction(instruction)
	if err != nil {
		return xerrors.Errorf("failed to create transaction: %v", err)
	}

	err = ctx.FillSignersAndSignWith(c.signer)
	if err != nil {
		return xerrors.Errorf("failed to sign transaction: %v", err)
	}

	_, err = c.bcClient.AddTransactionAndWait(ctx, 10)
	if err != nil {
		return xerrors.Errorf("failed to add transaction: %v", err)
	}

	return nil
}
