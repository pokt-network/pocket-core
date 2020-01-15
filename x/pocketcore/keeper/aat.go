package keeper

import (
	"encoding/hex"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
)

func AATGeneration(appPubKey string, clientPubKey string, passphrase string, keybase keys.Keybase) (pc.AAT, sdk.Error) {
	// get the public key from string
	pk, err := crypto.NewPublicKey(appPubKey)
	if err != nil {
		return pc.AAT{}, pc.NewPubKeyError(pc.ModuleName, err)
	}
	// create the aat object
	aat := pc.AAT{
		Version:              pc.SUPPORTEDTOKENVERSION,
		ApplicationPublicKey: appPubKey,
		ClientPublicKey:      clientPubKey,
		ApplicationSignature: "",
	}
	// marshal aat using json
	sig, _, err := (keybase).Sign(sdk.Address(pk.Address()), passphrase, aat.Hash())
	if err != nil {
		return pc.AAT{}, pc.NewSignatureError(pc.ModuleName, err)
	}
	// stringify the signature
	aat.ApplicationSignature = hex.EncodeToString(sig)
	return aat, nil
}
