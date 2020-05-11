package app

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	pocket "github.com/pokt-network/pocket-core/x/pocketcore"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/gov"
	"github.com/tendermint/tendermint/types"
	"os"
	"time"
)

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

func ExportState() (string, error) {
	j, err := pca.ExportAppState(false, nil)
	if err != nil {
		return "", err
	}
	j, _ = Codec().MarshalJSONIndent(types.GenesisDoc{
		GenesisTime: time.Now(),
		ChainID:     "pocket-test",
		ConsensusParams: &types.ConsensusParams{
			Block: types.BlockParams{
				MaxBytes:   15000,
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
		fmt.Println("could not unmarshal json in SortJSON: " + err.Error())
		os.Exit(1)
	}
	js, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		fmt.Println("could not marshal back to json in SortJSON: " + err.Error())
		os.Exit(1)
	}
	return string(js)
}
