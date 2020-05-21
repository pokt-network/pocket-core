package types

import (
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/codec"
)

// ModuleCdc is the codec for the module
var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
	codec.RegisterCrypto(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(MsgClaim{}, "pocketcore/claim", nil)
	cdc.RegisterConcrete(MsgProof{}, "pocketcore/Proof", nil)
	cdc.RegisterConcrete(Relay{}, "pocketcore/relay", nil)
	cdc.RegisterConcrete(Session{}, "pocketcore/session", nil)
	cdc.RegisterConcrete(RelayResponse{}, "pocketcore/relay_response", nil)
	cdc.RegisterInterface((*Proof)(nil), nil)
	cdc.RegisterConcrete(RelayProof{}, "pocketcore/relay_proof", nil)
	cdc.RegisterConcrete(ChallengeProofInvalidData{}, "pocketcore/challenge_proof_invalid_data", nil)
	cdc.RegisterConcrete(EvidencePersisted{}, "pocketcore/evidence_persisted", nil)
	cdc.RegisterInterface((*exported.ValidatorI)(nil), nil)
	cdc.RegisterConcrete(nodesTypes.Validator{}, "pos/Validator", nil) // todo does this really need to depend on nodes/types
}
