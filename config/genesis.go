package config

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/tendermint/tendermint/libs/common"
	tmtypes "github.com/tendermint/tendermint/types"
)

type GenesisState map[string]json.RawMessage

//  expected usage
//  func (app *nameServiceApp) InitChainer(ctx sdk.Ctx, req abci.RequestInitChain) abci.ResponseInitChain {
//	genesisState := GetGensisFromFile(app.cdc, "genesis.go")
//	return app.mm.InitGenesis(ctx, genesisState)
//}
func GenesisStateFromFile(cdc *codec.Codec, genFile string) GenesisState {
	if !common.FileExists(genFile) {
		panic(fmt.Errorf("%s does not exist, run `init` first", genFile))
	}
	genDoc := GenesisFileToGenDoc(genFile)
	return GenesisStateFromGenDoc(cdc, *genDoc)
}

func GenesisFileToGenDoc(genFile string) *tmtypes.GenesisDoc {
	if !common.FileExists(genFile) {
		panic(fmt.Errorf("%s does not exist, run `init` first", genFile))
	}
	genDoc, err := tmtypes.GenesisDocFromFile(genFile)
	if err != nil {
		panic(err)
	}
	return genDoc
}

func GenesisStateFromGenDoc(cdc *codec.Codec, genDoc tmtypes.GenesisDoc) (genesisState map[string]json.RawMessage) {
	if err := cdc.UnmarshalJSON(genDoc.AppState, &genesisState); err != nil {
		panic(err)
	}
	return genesisState
}

// InitConfig common config options for init
type InitConfig struct {
	ChainID   string
	GenTxsDir string
	Name      string
	NodeID    string
	ValPubKey crypto.PublicKey
}

// NewInitConfig creates a new InitConfig object
func NewInitConfig(chainID, genTxsDir, name, nodeID string, valPubKey crypto.PublicKey) InitConfig {
	return InitConfig{
		ChainID:   chainID,
		GenTxsDir: genTxsDir,
		Name:      name,
		NodeID:    nodeID,
		ValPubKey: valPubKey,
	}
}
