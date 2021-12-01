package types

import (
	"encoding/hex"
	"reflect"
	"testing"

	"github.com/pokt-network/pocket-core/types"
	"github.com/stretchr/testify/assert"
)

func TestMsgClaim_Route(t *testing.T) {
	assert.Equal(t, MsgClaim{}.Route(), RouterKey)
}

func TestMsgClaim_Type(t *testing.T) {
	assert.Equal(t, MsgClaim{}.Type(), MsgClaimName)
}

func TestMsgClaim_GetSigners(t *testing.T) {
	addr := getRandomValidatorAddress()
	signers := MsgClaim{
		SessionHeader: SessionHeader{},
		MerkleRoot:    HashRange{},
		TotalProofs:   0,
		FromAddress:   addr,
	}.GetSigners()
	assert.True(t, reflect.DeepEqual(signers, []types.Address{addr}))
}

func TestMsgClaim_ValidateBasic(t *testing.T) {
	appPubKey := getRandomPubKey().RawString()
	nodeAddress := getRandomValidatorAddress()
	ethereum := hex.EncodeToString([]byte{01})
	rootHash := Hash([]byte("fakeRoot"))
	root := HashRange{
		Hash:  rootHash,
		Range: Range{Upper: 100},
	}
	invalidClaimMessageSH := MsgClaim{
		SessionHeader: SessionHeader{
			ApplicationPubKey:  "",
			Chain:              ethereum,
			SessionBlockHeight: 1,
		},
		MerkleRoot:   root,
		TotalProofs:  100,
		FromAddress:  nodeAddress,
		EvidenceType: RelayEvidence,
	}
	invalidClaimMessageRoot := MsgClaim{
		SessionHeader: SessionHeader{
			ApplicationPubKey:  appPubKey,
			Chain:              ethereum,
			SessionBlockHeight: 1,
		},
		MerkleRoot: HashRange{
			Hash: []byte("bad_root"),
		},
		TotalProofs:  100,
		FromAddress:  nodeAddress,
		EvidenceType: RelayEvidence,
	}
	invalidClaimMessageRelays := MsgClaim{
		SessionHeader: SessionHeader{
			ApplicationPubKey:  appPubKey,
			Chain:              ethereum,
			SessionBlockHeight: 1,
		},
		MerkleRoot:   root,
		TotalProofs:  -1,
		FromAddress:  nodeAddress,
		EvidenceType: RelayEvidence,
	}
	invalidClaimMessageFromAddress := MsgClaim{
		SessionHeader: SessionHeader{
			ApplicationPubKey:  appPubKey,
			Chain:              ethereum,
			SessionBlockHeight: 1,
		},
		MerkleRoot:   root,
		TotalProofs:  -1,
		FromAddress:  types.Address{},
		EvidenceType: RelayEvidence,
	}
	invalidClaimMessageNoEvidence := MsgClaim{
		SessionHeader: SessionHeader{
			ApplicationPubKey:  appPubKey,
			Chain:              ethereum,
			SessionBlockHeight: 1,
		},
		MerkleRoot:  root,
		TotalProofs: 100,
		FromAddress: nodeAddress,
	}
	validClaimMessage := MsgClaim{
		SessionHeader: SessionHeader{
			ApplicationPubKey:  appPubKey,
			Chain:              ethereum,
			SessionBlockHeight: 1,
		},
		MerkleRoot:   root,
		TotalProofs:  100,
		FromAddress:  nodeAddress,
		EvidenceType: RelayEvidence,
	}
	tests := []struct {
		name     string
		msg      MsgClaim
		hasError bool
	}{
		{
			name:     "Invalid Claim Message, session header",
			msg:      invalidClaimMessageSH,
			hasError: true,
		},
		{
			name:     "Invalid Claim Message, root",
			msg:      invalidClaimMessageRoot,
			hasError: true,
		},
		{
			name:     "Invalid Claim Message, relays",
			msg:      invalidClaimMessageRelays,
			hasError: true,
		},
		{
			name:     "Invalid Claim Message, From Address",
			msg:      invalidClaimMessageFromAddress,
			hasError: true,
		},
		{
			name:     "Invalid Claim Message, No Evidence",
			msg:      invalidClaimMessageNoEvidence,
			hasError: true,
		},
		{
			name:     "Valid Claim Message",
			msg:      validClaimMessage,
			hasError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.msg.ValidateBasic() != nil, tt.hasError)
		})
	}
}

func TestMsgClaim_GetSignBytes(t *testing.T) {
	assert.NotPanics(t, func() { MsgClaim{}.GetSignBytes() })
}

func TestMsgProof_Route(t *testing.T) {
	assert.Equal(t, MsgProof{}.Route(), RouterKey)
}

func TestMsgProof_Type(t *testing.T) {
	assert.Equal(t, MsgProof{}.Type(), MsgProofName)
}

func TestMsgProof_GetSigners(t *testing.T) {
	pk := getRandomPubKey()
	addr := types.Address(pk.Address())
	signers := MsgProof{
		MerkleProof: MerkleProof{},
		Leaf: RelayProof{
			Entropy:            0,
			RequestHash:        pk.RawString(), // fake
			SessionBlockHeight: 0,
			ServicerPubKey:     pk.RawString(),
			Blockchain:         "",
			Token:              AAT{},
			Signature:          "",
		},
	}.GetSigners()
	assert.True(t, reflect.DeepEqual(signers, []types.Address{addr}))
}

func TestMsgProof_ValidateBasic(t *testing.T) {
	ethereum := hex.EncodeToString([]byte{01})
	servicerPubKey := getRandomPubKey().RawString()
	clientPrivKey := GetRandomPrivateKey()
	clientPubKey := clientPrivKey.PublicKey().RawString()
	appPrivKey := GetRandomPrivateKey()
	appPubKey := appPrivKey.PublicKey().RawString()
	hash1 := merkleHash([]byte("fake1"))
	hash2 := merkleHash([]byte("fake2"))
	hash3 := merkleHash([]byte("fake3"))
	hash4 := merkleHash([]byte("fake4"))
	validProofMessage := MsgProof{
		MerkleProof: MerkleProof{
			TargetIndex: 0,
			HashRanges: []HashRange{
				{
					Hash:  hash1,
					Range: Range{0, 1},
				},
				{
					Hash:  hash2,
					Range: Range{1, 2},
				},
				{
					Hash:  hash3,
					Range: Range{2, 3},
				},
			},
			Target: HashRange{Hash: hash4, Range: Range{3, 4}},
		},
		Leaf: RelayProof{
			Entropy:            1,
			SessionBlockHeight: 1,
			ServicerPubKey:     servicerPubKey,
			Blockchain:         ethereum,
			RequestHash:        servicerPubKey, // fake
			Token: AAT{
				Version:              "0.0.1",
				ApplicationPublicKey: appPubKey,
				ClientPublicKey:      clientPubKey,
				ApplicationSignature: "",
			},
			Signature: "",
		},
		EvidenceType: RelayEvidence,
	}
	vprLeaf := validProofMessage.Leaf.(RelayProof)
	signature, er := appPrivKey.Sign(vprLeaf.Token.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	vprLeaf.Token.ApplicationSignature = hex.EncodeToString(signature)
	clientSig, er := clientPrivKey.Sign(validProofMessage.Leaf.(RelayProof).Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	vprLeaf.Signature = hex.EncodeToString(clientSig)
	validProofMessage.Leaf = vprLeaf
	// invalid entropy
	invalidProofMsgIndex := validProofMessage
	//vprLeaf = validProofMessage.Leaf.LegacyFromProto().(*RelayProof)
	vprLeaf.Entropy = 0
	invalidProofMsgIndex.Leaf = vprLeaf
	// invalid merkleHash sum
	invalidProofMsgHashes := validProofMessage
	invalidProofMsgHashes.MerkleProof.HashRanges = []HashRange{}
	// invalid session block height
	invalidProofMsgSessionBlkHeight := validProofMessage
	//vprLeaf = validProofMessage.Leaf.LegacyFromProto().(*RelayProof)
	vprLeaf.SessionBlockHeight = -1
	invalidProofMsgSessionBlkHeight.Leaf = vprLeaf
	// invalid token
	invalidProofMsgToken := validProofMessage
	//vprLeaf = validProofMessage.Leaf.LegacyFromProto().(*RelayProof)
	vprLeaf.Token.ApplicationSignature = ""
	invalidProofMsgToken.Leaf = vprLeaf
	// invalid blockchain
	invalidProofMsgBlkchn := validProofMessage
	//vprLeaf = validProofMessage.Leaf.LegacyFromProto().(*RelayProof)
	vprLeaf.Blockchain = ""
	invalidProofMsgBlkchn.Leaf = vprLeaf
	// invalid signature
	invalidProofMsgSignature := validProofMessage
	//vprLeaf = validProofMessage.Leaf.LegacyFromProto().(*RelayProof)
	vprLeaf.Signature = hex.EncodeToString([]byte("foobar"))
	invalidProofMsgSignature.Leaf = vprLeaf
	tests := []struct {
		name     string
		msg      MsgProof
		hasError bool
	}{
		{
			name:     "Invalid Proof Message, signature",
			msg:      invalidProofMsgSignature,
			hasError: true,
		},
		{
			name:     "Invalid Proof Message, session block height",
			msg:      invalidProofMsgSessionBlkHeight,
			hasError: true,
		},
		{
			name:     "Invalid Proof Message, hashsum",
			msg:      invalidProofMsgHashes,
			hasError: true,
		},
		{
			name:     "Invalid Proof Message, leafnode index",
			msg:      invalidProofMsgIndex,
			hasError: true,
		},
		{
			name:     "Invalid Proof Message, token",
			msg:      invalidProofMsgToken,
			hasError: true,
		},
		{
			name:     "Invalid Proof Message, blockchain",
			msg:      invalidProofMsgBlkchn,
			hasError: true,
		},
		{
			name:     "Valid Proof Message",
			msg:      validProofMessage,
			hasError: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.msg.ValidateBasic()
			assert.Equal(t, tt.hasError, err != nil, err)
		})
	}
}

func TestMsgProof_GetSignBytes(t *testing.T) {
	assert.NotPanics(t, func() {
		MsgProof{}.GetSignBytes()
	})
}
