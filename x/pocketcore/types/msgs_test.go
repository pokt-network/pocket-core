package types

import (
	"encoding/binary"
	"encoding/hex"
	"github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"testing"
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
		MerkleRoot:    HashSum{},
		TotalProofs:   0,
		FromAddress:   addr,
	}.GetSigner()
	assert.Equal(t, types.Address(signers), addr)
}

func TestMsgClaim_ValidateBasic(t *testing.T) {
	appPubKey := getRandomPubKey().RawString()
	nodeAddress := getRandomValidatorAddress()
	ethereum := hex.EncodeToString([]byte{01})
	rootHash := Hash([]byte("fakeRoot"))
	rootSum := binary.LittleEndian.Uint64(rootHash)
	root := HashSum{
		Hash: rootHash,
		Sum:  rootSum,
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
		MerkleRoot: HashSum{
			Hash: []byte("bad_root"),
			Sum:  0,
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
		MerkleProofs: [2]MerkleProof{},
		Leaf: RelayProof{
			Entropy:            0,
			RequestHash:        pk.RawString(), // fake
			SessionBlockHeight: 0,
			ServicerPubKey:     pk.RawString(),
			Blockchain:         "",
			Token:              AAT{},
			Signature:          "",
		},
	}.GetSigner()
	assert.Equal(t, signers, addr)
}

func TestMsgProof_ValidateBasic(t *testing.T) {
	ethereum := hex.EncodeToString([]byte{01})
	servicerPubKey := getRandomPubKey().RawString()
	clientPrivKey := GetRandomPrivateKey()
	clientPubKey := clientPrivKey.PublicKey().RawString()
	appPrivKey := GetRandomPrivateKey()
	appPubKey := appPrivKey.PublicKey().RawString()
	hash1 := hash([]byte("fake1"))
	hash2 := hash([]byte("fake2"))
	hash3 := hash([]byte("fake3"))
	hash4 := hash([]byte("fake4"))
	hash5 := hash([]byte("fake5"))
	hash6 := hash([]byte("fake6"))
	validProofMessage := MsgProof{
		MerkleProofs: [2]MerkleProof{
			{
				Index: 0,
				HashSums: []HashSum{
					{
						Hash: hash1,
						Sum:  sumFromHash(hash1),
					},
					{
						Hash: hash2,
						Sum:  sumFromHash(hash2),
					},
					{
						Hash: hash3,
						Sum:  sumFromHash(hash3),
					},
				},
			},
			{
				Index: 2,
				HashSums: []HashSum{
					{
						Hash: hash4,
						Sum:  sumFromHash(hash4),
					},
					{
						Hash: hash5,
						Sum:  sumFromHash(hash5),
					},
					{
						Hash: hash6,
						Sum:  sumFromHash(hash6),
					},
				},
			}},
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
		Cousin: RelayProof{
			Entropy:            2,
			SessionBlockHeight: 1,
			ServicerPubKey:     servicerPubKey,
			RequestHash:        servicerPubKey, // fake
			Blockchain:         ethereum,
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
	vprCousin := validProofMessage.Cousin.(RelayProof)
	signature, er := appPrivKey.Sign(vprLeaf.Token.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	vprLeaf.Token.ApplicationSignature = hex.EncodeToString(signature)
	clientSig, er := clientPrivKey.Sign(validProofMessage.Leaf.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	vprLeaf.Signature = hex.EncodeToString(clientSig)
	signature2, er := appPrivKey.Sign(vprCousin.Token.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	vprCousin.Token.ApplicationSignature = hex.EncodeToString(signature2)
	clientSig2, er := clientPrivKey.Sign(validProofMessage.Cousin.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	vprCousin.Signature = hex.EncodeToString(clientSig2)
	validProofMessage.Leaf = vprLeaf
	validProofMessage.Cousin = vprCousin
	// invalid entropy
	invalidProofMsgIndex := validProofMessage
	vprLeaf = validProofMessage.Leaf.(RelayProof)
	vprLeaf.Entropy = 0
	invalidProofMsgIndex.Leaf = vprLeaf
	// invalid hash sum
	invalidProofMsgHashes := validProofMessage
	invalidProofMsgHashes.MerkleProofs[0].HashSums = []HashSum{}
	// invalid session block height
	invalidProofMsgSessionBlkHeight := validProofMessage
	vprLeaf = validProofMessage.Leaf.(RelayProof)
	vprLeaf.SessionBlockHeight = -1
	invalidProofMsgSessionBlkHeight.Leaf = vprLeaf
	// invalid token
	invalidProofMsgToken := validProofMessage
	vprLeaf = validProofMessage.Leaf.(RelayProof)
	vprLeaf.Token.ApplicationSignature = ""
	invalidProofMsgToken.Leaf = vprLeaf
	// invalid blockchain
	invalidProofMsgBlkchn := validProofMessage
	vprLeaf = validProofMessage.Leaf.(RelayProof)
	vprLeaf.Blockchain = ""
	invalidProofMsgBlkchn.Leaf = vprLeaf
	// invalid signature
	invalidProofMsgSignature := validProofMessage
	vprLeaf = validProofMessage.Leaf.(RelayProof)
	vprLeaf.Signature = hex.EncodeToString(clientSig2)
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
			assert.Equal(t, tt.msg.ValidateBasic() != nil, tt.hasError)
		})
	}
}

func TestMsgProof_GetSignBytes(t *testing.T) {
	assert.NotPanics(t, func() { MsgProof{}.GetSignBytes() })
}
