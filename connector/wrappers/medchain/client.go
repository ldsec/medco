package medchain

import (
	"github.com/ldsec/medchain/contracts"
	"go.dedis.ch/cothority/v3/byzcoin"
	"go.dedis.ch/cothority/v3/byzcoin/bcadmin/lib"
	"go.dedis.ch/cothority/v3/darc"
	"go.dedis.ch/protobuf"
	"golang.org/x/xerrors"
)

// QueryStatus denotes the status of a query.
type QueryStatus string

const (
	// QuerySuccessStatus denotes a "success" status
	QuerySuccessStatus QueryStatus = contracts.QuerySuccessStatus
	// QueryFailedStatus denotes a "failed" status
	QueryFailedStatus QueryStatus = contracts.QueryFailedStatus
)

// Client is used to interact with medchain.
type Client struct {
	bcClient *byzcoin.Client
	signer   darc.Signer
}

// NewClient creates a new Client by parsing a Byzcoin configuration file
// at `path` and uses `signer` for signing the transactions
func NewClient(path string, signer darc.Signer) (*Client, error) {
	_, bc, err := lib.LoadConfig(path)
	if err != nil {
		return nil, xerrors.Errorf("failed to initialize byzcoin client: %v", err)
	}

	return &Client{bcClient: bc, signer: signer}, nil
}

// GetAuthorization returns the Authorization a user has for a given query in a project.
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
	}

	ctx, err := c.bcClient.CreateTransaction(instr)
	if err != nil {
		return nil, byzcoin.InstanceID{}, xerrors.Errorf("failed to create transaction: %v", err)
	}

	err = c.bcClient.SignTransaction(ctx, c.signer)
	if err != nil {
		return nil, byzcoin.InstanceID{}, xerrors.Errorf("failed to sign transaction :%v", err)
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

// UpdateStatus updates the status for a given query.
func (c *Client) UpdateStatus(queryID byzcoin.InstanceID, status QueryStatus) error {
	instruction := byzcoin.Instruction{
		InstanceID: queryID,
		Invoke: &byzcoin.Invoke{
			Command:    contracts.QueryUpdateAction,
			ContractID: contracts.QueryContractID,
			Args: []byzcoin.Argument{{
				Name:  contracts.QueryStatusKey,
				Value: []byte(status),
			}},
		},
	}

	ctx, err := c.bcClient.CreateTransaction(instruction)
	if err != nil {
		return xerrors.Errorf("failed to create transaction: %v", err)
	}

	err = c.bcClient.SignTransaction(ctx, c.signer)
	if err != nil {
		return xerrors.Errorf("failed to sign transaction :%v", err)
	}

	_, err = c.bcClient.AddTransactionAndWait(ctx, 10)
	if err != nil {
		return xerrors.Errorf("failed to add transaction: %v", err)
	}

	return nil
}
