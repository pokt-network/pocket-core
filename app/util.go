package app

import (
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth"
	pocketKeeper "github.com/pokt-network/pocket-core/x/pocketcore/keeper"
	"github.com/tendermint/tendermint/types"
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

func ExportState(height int64) (string, error) {
	j, err := PCA.ExportAppState(height, false, nil)
	if err != nil {
		return "", err
	}
	j, _ = Codec().MarshalJSONIndent(types.GenesisDoc{
		ChainID: "<Input New ChainID>",
		ConsensusParams: &types.ConsensusParams{
			Block: types.BlockParams{
				MaxBytes:   4000000,
				MaxGas:     -1,
				TimeIotaMs: 1,
			},
			Evidence: types.EvidenceParams{
				MaxAge: 1000000,
			},
			Validator: types.ValidatorParams{
				PubKeyTypes: []string{"ed25519"},
			},
		},
		Validators: nil,
		AppHash:    nil,
		AppState:   j,
	}, "", "    ")
	return SortJSON(j), err
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

func UnmarshalTxStr(txStr string) auth.StdTx {
	txBytes, err := base64.StdEncoding.DecodeString(txStr)
	if err != nil {
		log.Fatal("error:", err)
	}
	return UnmarshalTx(txBytes)
}

func UnmarshalTx(txBytes []byte) auth.StdTx {
	defaultTxDecoder := auth.DefaultTxDecoder(cdc)
	tx, err := defaultTxDecoder(txBytes)
	if err != nil {
		log.Fatalf("Could not decode transaction: " + err.Error())
	}
	return tx.(auth.StdTx)
}
