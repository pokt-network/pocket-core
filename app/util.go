package app

import (
	"encoding/hex"
	"encoding/json"
	pocket "github.com/pokt-network/pocket-core/x/pocketcore"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/gov"
)

func GenerateChain(ticker, netid, version, client, inter string) (chain string, err error) {
	chain, err = pocket.GenerateChain(ticker, netid, version, client, inter)
	return
}

func GenerateAAT(appPubKey, clientPubKey, passphrase string) (aatjson []byte, err error) {
	aat, err := pocket.GenerateAAT(MustGetKeybase(), appPubKey, clientPubKey, passphrase)
	return json.MarshalIndent(aat, "", "  ")
}

func BuildMultisig(fromAddr, jsonMessage, passphrase string, pk crypto.PublicKeyMultiSig) ([]byte, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	var m sdk.Msg
	if err := Codec().UnmarshalJSON([]byte(jsonMessage), &m); err != nil {
		return nil, err
	}
	return gov.BuildAndSignMulti(Codec(), fa, pk, m, getTMClient(), MustGetKeybase(), passphrase)
}

func SignMultisigNext(fromAddr, txHex, passphrase string) ([]byte, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	bz, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, err
	}
	return gov.SignMulti(Codec(), fa, bz, nil, getTMClient(), MustGetKeybase(), passphrase)
}

func SignMultisigOutOfOrder(fromAddr, txHex, passphrase string, keys []crypto.PublicKey) ([]byte, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	bz, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, err
	}
	return gov.SignMulti(Codec(), fa, bz, keys, getTMClient(), MustGetKeybase(), passphrase)
}
