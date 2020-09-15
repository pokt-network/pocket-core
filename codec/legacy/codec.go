package legacy

import (
	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/crypto"
)

// Cdc defines a global generic sealed Amino codec to be used throughout sdk. It
// has all Tendermint crypto and evidence types registered.
//
// TODO: Deprecated - remove this global.
var Cdc *codec.LegacyAmino

func init() {
	Cdc = codec.NewLegacyAminoCodec()
	crypto.RegisterAmino(Cdc.Amino)
	codec.RegisterEvidences(Cdc, nil)
	Cdc.Seal()
}
