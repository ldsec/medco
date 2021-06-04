package medchain

import (
	"testing"
	"time"

	"github.com/ldsec/medchain/contracts"
	"github.com/stretchr/testify/require"
	"go.dedis.ch/cothority/v3"
	"go.dedis.ch/cothority/v3/byzcoin"
	"go.dedis.ch/cothority/v3/darc"
	"go.dedis.ch/onet/v3"
	"go.dedis.ch/protobuf"
)

func TestClient_GetAuthorization_Rejected(t *testing.T) {
	local := onet.NewTCPTest(cothority.Suite)
	defer local.CloseAll()

	signer := darc.NewSignerEd25519(nil, nil)
	_, roster, _ := local.GenTree(3, true)

	genesisMsg, err := byzcoin.DefaultGenesisMsg(byzcoin.CurrentVersion, roster,
		[]string{"spawn:project"}, signer.Identity())
	require.NoError(t, err)
	gDarc := &genesisMsg.GenesisDarc

	genesisMsg.BlockInterval = time.Second

	cl, _, err := byzcoin.NewLedger(genesisMsg, false)
	require.NoError(t, err)

	projectName := "name"

	ctx, err := addProject(t, projectName, "d", gDarc, signer, cl)
	require.NoError(t, err)

	client := Client{bcClient: cl, signer: signer}

	projectInstID := ctx.Instructions[0].DeriveID("")

	auth, _, err := client.GetAuthorization(projectInstID, "test", "testQueryID", "testQueryDef")
	require.Nil(t, auth)
	require.EqualError(t, err, "failed to authorize query: status "+contracts.QueryRejectedStatus)

	local.WaitDone(genesisMsg.BlockInterval)
}

func TestClient_GetAuthorization_Pending(t *testing.T) {
	local := onet.NewTCPTest(cothority.Suite)
	defer local.CloseAll()

	signer := darc.NewSignerEd25519(nil, nil)
	_, roster, _ := local.GenTree(3, true)

	genesisMsg, err := byzcoin.DefaultGenesisMsg(byzcoin.CurrentVersion, roster,
		[]string{"spawn:project", "invoke:project.add"}, signer.Identity())
	require.NoError(t, err)
	gDarc := &genesisMsg.GenesisDarc

	genesisMsg.BlockInterval = time.Second

	cl, _, err := byzcoin.NewLedger(genesisMsg, false)
	require.NoError(t, err)

	projectName := "name"
	userID := "testUserID"
	queryTerm := "testQueryTerm"

	ctx, err := addProject(t, projectName, "d", gDarc, signer, cl)
	require.NoError(t, err)

	projectInstID := ctx.Instructions[0].DeriveID("")

	ctx, err = cl.CreateTransaction(byzcoin.Instruction{
		InstanceID: projectInstID,
		Invoke: &byzcoin.Invoke{
			ContractID: contracts.ProjectContractID,
			Command:    "add",
			Args: byzcoin.Arguments{{
				Name:  contracts.ProjectUserIDKey,
				Value: []byte(userID),
			}, {
				Name:  contracts.ProjectQueryTermKey,
				Value: []byte(queryTerm),
			}},
		},
		SignerCounter: []uint64{2},
	})
	require.NoError(t, err)

	err = ctx.FillSignersAndSignWith(signer)
	require.NoError(t, err)

	_, err = cl.AddTransactionAndWait(ctx, 10)
	require.NoError(t, err)

	client := Client{bcClient: cl, signer: signer}

	auth, _, err := client.GetAuthorization(projectInstID, userID, "testQueryID", queryTerm)
	require.NoError(t, err)

	require.NotNil(t, auth)
	require.Equal(t, 1, len(auth.QueryTerms))
	require.Equal(t, queryTerm, auth.QueryTerms[0])
	require.Equal(t, userID, auth.UserID)

	local.WaitDone(genesisMsg.BlockInterval)
}

func TestClient_UpdateStatus_Correct(t *testing.T) {
	local := onet.NewTCPTest(cothority.Suite)
	defer local.CloseAll()

	signer := darc.NewSignerEd25519(nil, nil)
	_, roster, _ := local.GenTree(3, true)

	genesisMsg, err := byzcoin.DefaultGenesisMsg(byzcoin.CurrentVersion, roster,
		[]string{"spawn:project", "invoke:project.add"}, signer.Identity())
	require.NoError(t, err)
	gDarc := &genesisMsg.GenesisDarc

	genesisMsg.BlockInterval = time.Second

	cl, _, err := byzcoin.NewLedger(genesisMsg, false)
	require.NoError(t, err)

	projectName := "name"
	userID := "testUserID"
	queryTerm := "testQueryTerm"

	ctx, err := addProject(t, projectName, "d", gDarc, signer, cl)
	require.NoError(t, err)

	projectInstID := ctx.Instructions[0].DeriveID("")

	ctx, err = cl.CreateTransaction(byzcoin.Instruction{
		InstanceID: projectInstID,
		Invoke: &byzcoin.Invoke{
			ContractID: contracts.ProjectContractID,
			Command:    "add",
			Args: byzcoin.Arguments{{
				Name:  contracts.ProjectUserIDKey,
				Value: []byte(userID),
			}, {
				Name:  contracts.ProjectQueryTermKey,
				Value: []byte(queryTerm),
			}},
		},
		SignerCounter: []uint64{2},
	})
	require.NoError(t, err)

	err = ctx.FillSignersAndSignWith(signer)
	require.NoError(t, err)

	_, err = cl.AddTransactionAndWait(ctx, 10)
	require.NoError(t, err)

	client := Client{bcClient: cl, signer: signer}

	auth, queryInstID, err := client.GetAuthorization(projectInstID, userID, "testQueryID", queryTerm)
	require.NoError(t, err)
	require.NotNil(t, auth)

	err = client.UpdateStatus(queryInstID, QuerySuccessStatus)
	require.NoError(t, err)

	resp, err := cl.GetProofFromLatest(queryInstID.Slice())
	require.NoError(t, err)

	_, val, _, _, _ := resp.Proof.KeyValue()
	query := contracts.QueryContract{}
	err = protobuf.Decode(val, &query)

	require.Equal(t, QuerySuccessStatus, QueryStatus(query.Status))

	local.WaitDone(genesisMsg.BlockInterval)
}

func TestClient_UpdateStatus_Invalid(t *testing.T) {
	local := onet.NewTCPTest(cothority.Suite)
	defer local.CloseAll()

	signer := darc.NewSignerEd25519(nil, nil)
	_, roster, _ := local.GenTree(3, true)

	genesisMsg, err := byzcoin.DefaultGenesisMsg(byzcoin.CurrentVersion, roster,
		[]string{"spawn:project", "invoke:project.add"}, signer.Identity())
	require.NoError(t, err)
	gDarc := &genesisMsg.GenesisDarc

	genesisMsg.BlockInterval = time.Second

	cl, _, err := byzcoin.NewLedger(genesisMsg, false)
	require.NoError(t, err)

	projectName := "name"
	userID := "testUserID"
	queryTerm := "testQueryTerm"

	ctx, err := addProject(t, projectName, "d", gDarc, signer, cl)
	require.NoError(t, err)

	projectInstID := ctx.Instructions[0].DeriveID("")

	ctx, err = cl.CreateTransaction(byzcoin.Instruction{
		InstanceID: projectInstID,
		Invoke: &byzcoin.Invoke{
			ContractID: contracts.ProjectContractID,
			Command:    "add",
			Args: byzcoin.Arguments{{
				Name:  contracts.ProjectUserIDKey,
				Value: []byte(userID),
			}, {
				Name:  contracts.ProjectQueryTermKey,
				Value: []byte(queryTerm),
			}},
		},
		SignerCounter: []uint64{2},
	})
	require.NoError(t, err)

	err = ctx.FillSignersAndSignWith(signer)
	require.NoError(t, err)

	_, err = cl.AddTransactionAndWait(ctx, 10)
	require.NoError(t, err)

	client := Client{bcClient: cl, signer: signer}

	auth, queryInstID, err := client.GetAuthorization(projectInstID, userID, "testQueryID", queryTerm)
	require.NoError(t, err)
	require.NotNil(t, auth)

	err = client.UpdateStatus(queryInstID, "invalid")
	require.Error(t, err)

	local.WaitDone(genesisMsg.BlockInterval)
}

func addProject(t *testing.T, name, description string,
	gDarc *darc.Darc, signer darc.Signer, cl *byzcoin.Client) (byzcoin.ClientTransaction, error) {

	instruction := byzcoin.Instruction{
		InstanceID: byzcoin.NewInstanceID(gDarc.GetBaseID()),
		Spawn: &byzcoin.Spawn{
			ContractID: contracts.ProjectContractID,
			Args: []byzcoin.Argument{{
				Name:  contracts.ProjectDescriptionKey,
				Value: []byte(description),
			}, {
				Name:  contracts.ProjectNameKey,
				Value: []byte(name),
			}},
		},
		SignerCounter: []uint64{1},
	}

	ctx, err := cl.CreateTransaction(instruction)
	require.NoError(t, err)
	require.NoError(t, ctx.FillSignersAndSignWith(signer))

	_, err = cl.AddTransactionAndWait(ctx, 10)
	return ctx, err
}
