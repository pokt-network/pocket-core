package app

import (
	"encoding/hex"
	"encoding/json"
	pocket "github.com/pokt-network/pocket-core/x/pocketcore"
	"github.com/pokt-network/posmint/crypto"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/pokt-network/posmint/x/gov"
	"github.com/tendermint/tendermint/types"
	"os"
	"time"
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

func ExportState(filepath string) error {
	if pca == nil {
		return NewNilPocketCoreAppError()
	}
	j, err := pca.ExportAppState(false, nil)
	if err != nil {
		return err
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
	// create a new file
	jsonFile, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}
	// write to the file
	_, err = jsonFile.Write(newDefaultGenesisState())
	if err != nil {
		return err
	}
	// close the file
	err = jsonFile.Close()
	if err != nil {
		return err
	}
	return nil
}
