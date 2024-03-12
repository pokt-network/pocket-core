package keeper

import (
	"encoding/hex"

	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
)

// GenerateAAT generates an AAT to be used for relay request authentication.
// - appPubKey is the public key of the application that's staking for on-chain service.
// - clientPubKey (a.k.a gatewayPubKey) is the public key of the Gateway that's facilitating relays on behalf of the app.
// - appPubKey and clientPubKey may or may not be the same.
func AATGeneration(appPubKey, clientPubKey string, appPrivKey crypto.PrivateKey) (pc.AAT, sdk.Error) {
	aat := pc.AAT{
		Version:              pc.SupportedTokenVersions[0],
		ApplicationPublicKey: appPubKey,
		ClientPublicKey:      clientPubKey,
		ApplicationSignature: "",
	}

	// marshal the AAT structure
	aatBytes := aat.Hash()

	// This is where the `ApplicationPrivKey` signs (i.e. delegates trust) to
	// the underlying`ClientPublicKey`.
	sig, err := appPrivKey.Sign(aatBytes)
	if err != nil {
		return pc.AAT{}, pc.NewSignatureError(pc.ModuleName, err)
	}

	aat.ApplicationSignature = hex.EncodeToString(sig)
	return aat, nil
}
