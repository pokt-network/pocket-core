package pocketcore

import (
	"github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	"github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/auth"
	"github.com/pokt-network/posmint/x/auth/util"
)

// "ClaimTx" - A transaction that sends the total number of proofs (claim), the merkle root (for data integrity), and the header (for identification)
func ClaimTx(kp crypto.PrivateKey, cliCtx util.CLIContext, txBuilder auth.TxBuilder, header types.SessionHeader, totalProofs int64, root types.HashSum, evidenceType types.EvidenceType) (*sdk.TxResponse, error) {
	msg := types.MsgClaim{
		SessionHeader:    header,
		TotalProofs:      totalProofs,
		MerkleRoot:       root,
		FromAddress:      sdk.Address(kp.PublicKey().Address()),
		EvidenceType:     evidenceType,
		ExpirationHeight: 0, // leave as zero
	}
	err := msg.ValidateBasic()
	if err != nil {
		return nil, err
	}
	return util.CompleteAndBroadcastTxCLI(txBuilder, cliCtx, []sdk.Msg{msg})
}

// "ProofTx" - A transaction to prove the claim that was previously sent (Merkle Proofs and leaf/cousin)
func ProofTx(cliCtx util.CLIContext, txBuilder auth.TxBuilder, branches [2]types.MerkleProof, leafNode, cousinNode types.Proof) (*sdk.TxResponse, error) {
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

// "GenerateAAT" - Exported call to generate an application authentication token
func GenerateAAT(keybase keys.Keybase, appPubKey, cliPubKey, passphrase string) (types.AAT, error) {
	return keeper.AATGeneration(appPubKey, cliPubKey, passphrase, keybase)
}
