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
		TotalRelays:   0,
		FromAddress:   addr,
	}.GetSigners()
	assert.Len(t, signers, 1)
	assert.Equal(t, types.Address(signers[0]), addr)
}

func TestMsgClaim_ValidateBasic(t *testing.T) {
	appPubKey := getRandomPubKey().RawString()
	nodeAddress := getRandomValidatorAddress()
	ethereum, err := NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "geth",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
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
		MerkleRoot:  root,
		TotalRelays: 100,
		FromAddress: nodeAddress,
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
		TotalRelays: 100,
		FromAddress: nodeAddress,
	}
	invalidClaimMessageRelays := MsgClaim{
		SessionHeader: SessionHeader{
			ApplicationPubKey:  appPubKey,
			Chain:              ethereum,
			SessionBlockHeight: 1,
		},
		MerkleRoot:  root,
		TotalRelays: -1,
		FromAddress: nodeAddress,
	}
	invalidClaimMessageFromAddress := MsgClaim{
		SessionHeader: SessionHeader{
			ApplicationPubKey:  appPubKey,
			Chain:              ethereum,
			SessionBlockHeight: 1,
		},
		MerkleRoot:  root,
		TotalRelays: -1,
		FromAddress: types.Address{},
	}
	validClaimMessage := MsgClaim{
		SessionHeader: SessionHeader{
			ApplicationPubKey:  appPubKey,
			Chain:              ethereum,
			SessionBlockHeight: 1,
		},
		MerkleRoot:  root,
		TotalRelays: 100,
		FromAddress: nodeAddress,
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
			SessionBlockHeight: 0,
			ServicerPubKey:     pk.RawString(),
			Blockchain:         "",
			Token:              AAT{},
			Signature:          "",
		},
	}.GetSigners()
	assert.Len(t, signers, 1)
	assert.Equal(t, signers[0], addr)
}

func TestMsgProof_ValidateBasic(t *testing.T) {
	ethereum, err := NonNativeChain{
		Ticker:  "eth",
		Netid:   "4",
		Version: "v1.9.9",
		Client:  "geth",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	servicerPubKey := getRandomPubKey().RawString()
	clientPrivKey := getRandomPrivateKey()
	clientPubKey := clientPrivKey.PublicKey().RawString()
	appPrivKey := getRandomPrivateKey()
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
			Blockchain:         ethereum,
			Token: AAT{
				Version:              "0.0.1",
				ApplicationPublicKey: appPubKey,
				ClientPublicKey:      clientPubKey,
				ApplicationSignature: "",
			},
			Signature: "",
		},
	}
	signature, er := appPrivKey.Sign(validProofMessage.Leaf.Token.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validProofMessage.Leaf.Token.ApplicationSignature = hex.EncodeToString(signature)
	clientSig, er := clientPrivKey.Sign(validProofMessage.Leaf.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validProofMessage.Leaf.Signature = hex.EncodeToString(clientSig)
	signature2, er := appPrivKey.Sign(validProofMessage.Cousin.Token.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validProofMessage.Cousin.Token.ApplicationSignature = hex.EncodeToString(signature2)
	clientSig2, er := clientPrivKey.Sign(validProofMessage.Cousin.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validProofMessage.Cousin.Signature = hex.EncodeToString(clientSig2)
	// invalid entropy
	invalidProofMsgIndex := validProofMessage
	invalidProofMsgIndex.Leaf.Entropy = 0
	// invalid hash sum
	invalidProofMsgHashes := validProofMessage
	invalidProofMsgHashes.MerkleProofs[0].HashSums = []HashSum{}
	// invalid session block height
	invalidProofMsgSessionBlkHeight := validProofMessage
	invalidProofMsgSessionBlkHeight.Leaf.SessionBlockHeight = -1
	// invalid token
	invalidProofMsgToken := validProofMessage
	invalidProofMsgToken.Leaf.Token.ApplicationSignature = ""
	// invalid blockchain
	invalidProofMsgBlkchn := validProofMessage
	invalidProofMsgBlkchn.Leaf.Blockchain = ""
	// invalid signature
	invalidProofMsgSignature := validProofMessage
	invalidProofMsgSignature.Leaf.Signature = hex.EncodeToString(clientSig2)
	tests := []struct {
		name     string
		msg      MsgProof
		hasError bool
	}{
		{
			name:     "Invalid RelayProof Message, signature",
			msg:      invalidProofMsgSignature,
			hasError: true,
		},
		{
			name:     "Invalid RelayProof Message, session block height",
			msg:      invalidProofMsgSessionBlkHeight,
			hasError: true,
		},
		{
			name:     "Invalid RelayProof Message, hashsum",
			msg:      invalidProofMsgHashes,
			hasError: true,
		},
		{
			name:     "Invalid RelayProof Message, leafnode index",
			msg:      invalidProofMsgIndex,
			hasError: true,
		},
		{
			name:     "Invalid RelayProof Message, token",
			msg:      invalidProofMsgToken,
			hasError: true,
		},
		{
			name:     "Invalid RelayProof Message, blockchain",
			msg:      invalidProofMsgBlkchn,
			hasError: true,
		},
		{
			name:     "Valid RelayProof Message",
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
