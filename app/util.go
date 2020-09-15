package app

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth"
	"github.com/pokt-network/pocket-core/x/auth/types"
	pocketKeeper "github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	"log"
)

func GenerateAAT(appPubKey, clientPubKey string, key crypto.PrivateKey) (aatjson []byte, err error) {
	aat, er := pocketKeeper.AATGeneration(appPubKey, clientPubKey, key)
	if er != nil {
		return nil, er
	}
	return json.MarshalIndent(aat, "", "  ")
}

func BuildMultisig(fromAddr, jsonMessage, passphrase, chainID string, pk crypto.PublicKeyMultiSig, fees int64) ([]byte, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	var m sdk.Msg
	if err := Codec().UnmarshalJSON([]byte(jsonMessage), &m); err != nil {
		return nil, err
	}
	kb, err := GetKeybase()
	if err != nil {
		return nil, err
	}
	txBuilder := auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc),
		auth.DefaultTxDecoder(cdc),
		chainID,
		"", nil).WithKeybase(kb)
	return txBuilder.BuildAndSignMultisigTransaction(fa, pk, m, passphrase, fees)
}

func SignMultisigNext(fromAddr, txHex, passphrase, chainID string) ([]byte, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	bz, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, err
	}
	kb, err := GetKeybase()
	if err != nil {
		return nil, err
	}
	txBuilder := auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc),
		auth.DefaultTxDecoder(cdc),
		chainID,
		"", nil).WithKeybase(kb)
	return txBuilder.SignMultisigTransaction(fa, nil, passphrase, bz)
}

func SignMultisigOutOfOrder(fromAddr, txHex, passphrase, chainID string, keys []crypto.PublicKey) ([]byte, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	bz, err := hex.DecodeString(txHex)
	if err != nil {
		return nil, err
	}
	kb, err := GetKeybase()
	if err != nil {
		return nil, err
	}
	txBuilder := auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc),
		auth.DefaultTxDecoder(cdc),
		chainID,
		"", nil).WithKeybase(kb)
	return txBuilder.SignMultisigTransaction(fa, keys, passphrase, bz)
}

func SortJSON(toSortJSON []byte) string {
	var c interface{}
	err := json.Unmarshal(toSortJSON, &c)
	if err != nil {
		log.Fatal("could not unmarshal json in SortJSON: " + err.Error())
	}
	js, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		log.Fatalf("could not marshal back to json in SortJSON: " + err.Error())
	}
	return string(js)
}

func UnmarshalTxStr(txStr string) types.StdTxI {
	txBytes, err := base64.StdEncoding.DecodeString(txStr)
	if err != nil {
		log.Fatal("error:", err)
	}
	return UnmarshalTx(txBytes)
}

func UnmarshalTx(txBytes []byte) types.StdTxI {
	defaultTxDecoder := auth.DefaultTxDecoder(cdc)
	tx, err := defaultTxDecoder(txBytes)
	if err != nil {
		log.Fatalf("Could not decode transaction: " + err.Error())
	}
	if cdc.IsAfterUpgrade() {
		return tx.(auth.StdTx)
	} else {
		return tx.(auth.LegacyStdTx)

	}
}
