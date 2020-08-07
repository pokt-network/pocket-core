package keeper

import (
	"encoding/hex"

	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
)

// "AATGeneration" - Generates an application authentication token with an application public key hex string
// a client public key hex string, a passphrase and a keybase. The contract is that the keybase contains the app pub key
// and the passphrase corresponds to the app public key keypair.
func AATGeneration(appPubKey string, clientPubKey string, key crypto.PrivateKey) (pc.AAT, sdk.Error) {
	// create the aat object
	aat := pc.AAT{
		Version:              pc.SupportedTokenVersions[0],
		ApplicationPublicKey: appPubKey,
		ClientPublicKey:      clientPubKey,
		ApplicationSignature: "",
	}
	// marshal aat using json
	sig, err := key.Sign(aat.Hash())
	if err != nil {
		return pc.AAT{}, pc.NewSignatureError(pc.ModuleName, err)
	}
	// stringify the signature into hex
	aat.ApplicationSignature = hex.EncodeToString(sig)
	return aat, nil
}
