package pocketcore

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
)

// transaction that sends the total number of relays (claim), the merkle root (for data integrity), and the header (for identification)
func ClaimTx(keybase keys.Keybase, cliCtx util.CLIContext, txBuilder auth.TxBuilder, header types.SessionHeader, totalRelays int64, root types.HashSum) (*sdk.TxResponse, error) {
	kp, err := keybase.GetCoinbase()
	if err != nil {
		return nil, err
	}
	msg := types.MsgClaim{
		SessionHeader: header,
		TotalRelays:   totalRelays,
		MerkleRoot:    root,
		FromAddress:   kp.GetAddress(),
	}
	err = msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

// transaction to prove the
func ProofTx(cliCtx util.CLIContext, txBuilder auth.TxBuilder, branches [2]types.MerkleProof, leafNode, cousinNode types.RelayProof) (*sdk.TxResponse, error) {
	msg := types.MsgProof{
		MerkleProofs: branches,
		Leaf:         leafNode,
		Cousin:       cousinNode,
	}
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func GenerateChain(ticker, netid, version, client, inter string) (string, error) {
	return keeper.GenerateChain(ticker, netid, version, client, inter)
}

func GenerateAAT(keybase keys.Keybase, appPubKey, cliPubKey, passphrase string) (types.AAT, error) {
	return keeper.AATGeneration(appPubKey, cliPubKey, passphrase, keybase)
}
