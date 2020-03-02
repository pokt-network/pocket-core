package types

import (
	"encoding/hex"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestProof_Validate(t *testing.T) {
	appPrivateKey := getRandomPrivateKey()
	clientPrivateKey := getRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	servicerPubKey := getRandomPubKey().RawString()
	clientPubKey := clientPrivateKey.PublicKey().RawString()
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
	hbs := HostedBlockchains{
		M: map[string]HostedBlockchain{ethereum: {Hash: ethereum, URL: "https://www.google.com"}},
		l: sync.Mutex{},
		o: sync.Once{},
	}
	bitcoin, err := NonNativeChain{
		Ticker:  "btc",
		Netid:   "1",
		Version: "0.19.0.1",
		Client:  "",
		Inter:   "",
	}.HashString()
	if err != nil {
		t.Fatalf(err.Error())
	}
	payload := Payload{Data: "fake"}
	validProof := RelayProof{
		Entropy:            0,
		SessionBlockHeight: 1,
		ServicerPubKey:     servicerPubKey,
		RequestHash:        payload.HashString(), // fake
		Blockchain:         ethereum,
		Token: AAT{
			Version:              "0.0.1",
			ApplicationPublicKey: appPubKey,
			ClientPublicKey:      clientPubKey,
			ApplicationSignature: "",
		},
		Signature: "",
	}
	appSignature, er := appPrivateKey.Sign(validProof.Token.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validProof.Token.ApplicationSignature = hex.EncodeToString(appSignature)
	clientSignature, er := clientPrivateKey.Sign(validProof.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validProof.Signature = hex.EncodeToString(clientSignature)
	// invalidProof sessionBlockHeight
	invalidProofSessionBlock := validProof
	invalidProofSessionBlock.SessionBlockHeight = -1
	// invalid Proof servicer public key
	invalidProofServicerPubKey := validProof
	invalidProofServicerPubKey.ServicerPubKey = ""
	// invalid Proof wrong verify pub key
	wrongVerifyPubKey := getRandomPubKey().RawString()
	invalidProofServicerPubKeyVerify := validProof
	invalidProofServicerPubKeyVerify.ServicerPubKey = wrongVerifyPubKey
	// invalid Proof blockchain
	invalidProofBlockchain := validProof
	invalidProofBlockchain.Blockchain = ""
	// invalid Proof nothosted blockchain
	invalidProofNotHostedBlockchain := validProof
	invalidProofNotHostedBlockchain.Blockchain = bitcoin
	// invalid Proof AAT
	invalidProofInvalidAAT := validProof
	invalidProofInvalidAAT.Token.ApplicationSignature = hex.EncodeToString(clientSignature) // wrong signature
	// invalid Proof Request Hash
	invalidProofRequestHash := validProof
	invalidProofRequestHash.RequestHash = servicerPubKey
	// invalid Proof no client signature
	invalidProofClientSignature := validProof
	invalidProofClientSignature.Signature = hex.EncodeToString(appSignature) // wrong signature
	// over max relays
	overMaxRelays := int64(0)
	// over number of chains
	overNumberOfChains := 2
	tests := []struct {
		name             string
		proof            Proof
		maxRelays        int64
		numOfChains      int
		sessionNodeCount int
		verifyPubKey     string
		hb               HostedBlockchains
		hasError         bool
	}{
		{
			name:             "Invalid Proof: session block",
			proof:            invalidProofSessionBlock,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Invalid Proof: servicer pub key",
			proof:            invalidProofServicerPubKey,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Invalid Proof: verify pub key",
			proof:            invalidProofServicerPubKeyVerify,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Invalid Proof: blockchain",
			proof:            invalidProofBlockchain,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Invalid Proof: not hosted chain",
			proof:            invalidProofNotHostedBlockchain,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Invalid Proof: invalid AAT",
			proof:            invalidProofInvalidAAT,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Invalid Proof: client signature",
			proof:            invalidProofClientSignature,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Invalid Proof: over max relays",
			proof:            validProof,
			maxRelays:        overMaxRelays,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Invalid Proof: over number of chains",
			proof:            validProof,
			maxRelays:        1,
			numOfChains:      overNumberOfChains,
			sessionNodeCount: 0,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Invalid Proof: invalid request hash from payload",
			proof:            invalidProofRequestHash,
			maxRelays:        5,
			numOfChains:      2,
			sessionNodeCount: 0,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Valid Proof",
			proof:            validProof,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.proof.Validate(tt.maxRelays, tt.numOfChains, tt.sessionNodeCount, 1, tt.hb, payload.HashString(), tt.verifyPubKey) != nil, tt.hasError)
		})
	}
}

func TestProof_Bytes(t *testing.T) {
	appPubKey := getRandomPubKey().RawString()
	servicerPubKey := getRandomPubKey().RawString()
	clientPubKey := getRandomPubKey().RawString()
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
	pro := RelayProof{
		Entropy:            0,
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
	}
	proof2 := pro
	proof2.Signature = hex.EncodeToString([]byte("fake Signature"))
	assert.Equal(t, pro.Hash(), proof2.Hash())
	assert.NotEqual(t, pro.HashWithSignature(), proof2.HashWithSignature())
	var p RelayProof
	assert.Nil(t, json.Unmarshal(pro.Bytes(), &p))
}

func TestProof_HashAndHashString(t *testing.T) {
	appPubKey := getRandomPubKey().RawString()
	servicerPubKey := getRandomPubKey().RawString()
	clientPubKey := getRandomPubKey().RawString()
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
	pro := RelayProof{
		Entropy:            0,
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
	}
	assert.Nil(t, HashVerification(hex.EncodeToString(pro.Hash())))
	assert.Nil(t, HashVerification(pro.HashString()))
	assert.Nil(t, HashVerification(hex.EncodeToString(pro.HashWithSignature())))
	assert.Nil(t, HashVerification(pro.HashStringWithSignature()))
}
