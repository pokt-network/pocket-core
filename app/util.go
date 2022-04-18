package app

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"reflect"

	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth"
	"github.com/pokt-network/pocket-core/x/auth/types"
	pocketKeeper "github.com/pokt-network/pocket-core/x/pocketcore/keeper"
)

func GenerateAAT(appPubKey, clientPubKey string, key crypto.PrivateKey) (aatjson []byte, err error) {
	aat, er := pocketKeeper.AATGeneration(appPubKey, clientPubKey, key)
	if er != nil {
		return nil, er
	}
	return json.MarshalIndent(aat, "", "  ")
}

func BuildMultisig(fromAddr, jsonMessage, passphrase, chainID string, pk crypto.PublicKeyMultiSig, fees int64, legacyCodec bool) ([]byte, error) {
	fa, err := sdk.AddressFromHex(fromAddr)
	if err != nil {
		return nil, err
	}
	var m sdk.Msg
	if err := Codec().UnmarshalJSON([]byte(jsonMessage), &m); err != nil {
		return nil, err
	}
	// use reflection to convert to proto msg
	val := reflect.ValueOf(m)
	vp := reflect.New(val.Type())
	vp.Elem().Set(val)
	protoMsg := vp.Interface().(sdk.ProtoMsg)
	kb, err := GetKeybase()
	if err != nil {
		return nil, err
	}
	txBuilder := auth.NewTxBuilder(
		auth.DefaultTxEncoder(cdc),
		auth.DefaultTxDecoder(cdc),
		chainID,
		"", nil).WithKeybase(kb)
	return txBuilder.BuildAndSignMultisigTransaction(fa, pk, protoMsg, passphrase, fees, legacyCodec)
}

func SignMultisigNext(fromAddr, txHex, passphrase, chainID string, legacyCodec bool) ([]byte, error) {
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
	return txBuilder.SignMultisigTransaction(fa, nil, passphrase, bz, legacyCodec)
}

func SignMultisigOutOfOrder(fromAddr, txHex, passphrase, chainID string, keys []crypto.PublicKey, legacyCodec bool) ([]byte, error) {
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
	return txBuilder.SignMultisigTransaction(fa, keys, passphrase, bz, legacyCodec)
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

func UnmarshalTxStr(txStr string, height int64) (types.StdTx, error) {
	txBytes, err := base64.StdEncoding.DecodeString(txStr)
	if err != nil {
		return types.StdTx{}, err
	}
	return UnmarshalTx(txBytes, height)
}

func UnmarshalTx(txBytes []byte, height int64) (types.StdTx, error) {
	defaultTxDecoder := auth.DefaultTxDecoder(cdc)
	tx, err := defaultTxDecoder(txBytes, height)
	if err != nil {
		return types.StdTx{}, fmt.Errorf("Could not decode transaction: " + err.Error())
	}
	return tx.(auth.StdTx), nil
}
