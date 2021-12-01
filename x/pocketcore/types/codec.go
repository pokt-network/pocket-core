package types

import (
	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/codec/types"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
)

// module wide codec
var ModuleCdc *codec.Codec

func init() {
	ModuleCdc = codec.NewCodec(types.NewInterfaceRegistry())
	RegisterCodec(ModuleCdc)
	crypto.RegisterAmino(ModuleCdc.AminoCodec().Amino)
}

// RegisterCodec registers concrete types on the codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterStructure(MsgClaim{}, "pocketcore/claim")
	cdc.RegisterStructure(MsgProtoProof{}, "pocketcore/protoProof")
	cdc.RegisterStructure(MsgProof{}, "pocketcore/proof")
	cdc.RegisterStructure(Relay{}, "pocketcore/relay")
	cdc.RegisterStructure(Session{}, "pocketcore/session")
	cdc.RegisterStructure(RelayResponse{}, "pocketcore/relay_response")
	cdc.RegisterStructure(RelayProof{}, "pocketcore/relay_proof")
	cdc.RegisterStructure(ChallengeProofInvalidData{}, "pocketcore/challenge_proof_invalid_data")
	cdc.RegisterStructure(ProofI_RelayProof{}, "pocketcore/proto_relay_proofI")
	cdc.RegisterStructure(ProofI_ChallengeProof{}, "pocketcore/proto_challenge_proofI")
	cdc.RegisterStructure(ProtoEvidence{}, "pocketcore/evidence_persisted")
	cdc.RegisterStructure(nodesTypes.Validator{}, "pos/8.0Validator")    // todo does this really need to depend on nodes/types
	cdc.RegisterStructure(nodesTypes.LegacyValidator{}, "pos/Validator") // todo does this really need to depend on nodes/types
	cdc.RegisterInterface("x.pocketcore.Proof", (*Proof)(nil), &RelayProof{}, &ChallengeProofInvalidData{})
	cdc.RegisterInterface("types.isProofI_Proof", (*isProofI_Proof)(nil))
	cdc.RegisterImplementation((*sdk.ProtoMsg)(nil), &MsgClaim{}, &MsgProof{})
	cdc.RegisterImplementation((*sdk.Msg)(nil), &MsgClaim{}, &MsgProof{})
	ModuleCdc = cdc
}
