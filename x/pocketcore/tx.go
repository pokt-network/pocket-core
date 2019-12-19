package pocketcore

import (
	"github.com/pokt-network/merkle"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
)

// transaction that sends the total number of relays (claim), the merkle root (for data integrity), and the header (for identification)
func (am AppModule) ClaimTx(cdc *codec.Codec, cliCtx util.CLIContext, txBuilder auth.TxBuilder, header types.SessionHeader, totalRelays int64, root []byte) (*sdk.TxResponse, error) {
	msg := types.MsgClaim{
		SessionHeader: header,
		TotalRelays:   totalRelays,
		Root:          root,
		FromAddress:   sdk.ValAddress(am.node.PrivValidator().GetPubKey().Address()),
	}
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

// transaction to prove the
func (am AppModule) ProofTx(cdc *codec.Codec, cliCtx util.CLIContext, txBuilder auth.TxBuilder, porBranch merkle.Proof, leafNode types.Proof) (*sdk.TxResponse, error) {
	msg := types.MsgProof{
		Proof:    porBranch,
		LeafNode: leafNode,
	}
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

func (am AppModule) GenerateChain(ticker, netid, version, client, inter string) (string, error) {
	return am.keeper.GenerateChain(ticker, netid, version, client, inter)
}

func (am AppModule) GenerateAAT(appPubKey, cliPubKey, passphrase string) (types.AAT, error) {
	return am.keeper.AATGeneration(appPubKey, cliPubKey, passphrase, am.keybase)
}
