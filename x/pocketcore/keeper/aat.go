package keeper

import (
	"encoding/hex"
	pc "github.com/pokt-network/pocket-core/x/pocketcore/types"
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/crypto/keys"
	sdk "github.com/pokt-network/posmint/types"
)

func (k Keeper) AATGeneration(appPubKey string, clientPubKey string, passphrase string, keybase keys.Keybase) (pc.AAT, sdk.Error) {
	pk, err := crypto.NewPublicKey(appPubKey)
	if err != nil {
		return pc.AAT{}, pc.NewPubKeyError(pc.ModuleName, err)
	}
	aat := pc.AAT{
		Version:              pc.SUPPORTEDTOKENVERSION,
		ApplicationPublicKey: appPubKey,
		ClientPublicKey:      clientPubKey,
		ApplicationSignature: "",
	}
	res := k.cdc.MustMarshalBinaryBare(aat)
	sig, _, err := keybase.Sign(sdk.AccAddress(pk.Address()), passphrase, res)
	if err != nil {
		return pc.AAT{}, pc.NewSignatureError(pc.ModuleName, err)
	}
	aat.ApplicationSignature = hex.EncodeToString(sig)
	return aat, nil
}
