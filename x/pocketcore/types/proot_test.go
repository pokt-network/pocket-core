package types

import (
	"encoding/hex"
	"encoding/json"
	"github.com/pokt-network/posmint/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"sync"
	"testing"
)

func TestProof_Validate(t *testing.T) {
	appPrivateKey := getRandomPrivateKey()
	clientPrivateKey := getRandomPrivateKey()
	appPubKey := crypto.PublicKey(appPrivateKey.PubKey().(ed25519.PubKeyEd25519)).String()
	servicerPubKey := crypto.PublicKey(getRandomPubKey()).String()
	clientPubKey := crypto.PublicKey(clientPrivateKey.PubKey().(ed25519.PubKeyEd25519)).String()
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
	validProof := RelayProof{
		Entropy:            0,
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
	// invalid RelayProof servicer public key
	invalidProofServicerPubKey := validProof
	invalidProofServicerPubKey.ServicerPubKey = ""
	// invalid RelayProof wrong verify pub key
	wrongVerifyPubKey := crypto.PublicKey(getRandomPubKey()).String()
	invalidProofServicerPubKeyVerify := validProof
	invalidProofServicerPubKeyVerify.ServicerPubKey = wrongVerifyPubKey
	// invalid RelayProof blockchain
	invalidProofBlockchain := validProof
	invalidProofBlockchain.Blockchain = ""
	// invalid RelayProof nothosted blockchain
	invalidProofNotHostedBlockchain := validProof
	invalidProofNotHostedBlockchain.Blockchain = bitcoin
	// invalid RelayProof AAT
	invalidProofInvalidAAT := validProof
	invalidProofInvalidAAT.Token.ApplicationSignature = hex.EncodeToString(clientSignature) // wrong signature
	// invalid RelayProof no client signature
	invalidProofClientSignature := validProof
	invalidProofClientSignature.Signature = hex.EncodeToString(appSignature) // wrong signature
	// over max relays
	overMaxRelays := int64(0)
	// over number of chains
	overNumberOfChains := 2
	tests := []struct {
		name             string
		proof            RelayProof
		maxRelays        int64
		numOfChains      int
		sessionNodeCount int
		verifyPubKey     string
		hb               HostedBlockchains
		hasError         bool
	}{
		{
			name:             "Invalid RelayProof: session block",
			proof:            invalidProofSessionBlock,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Invalid RelayProof: servicer pub key",
			proof:            invalidProofServicerPubKey,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Invalid RelayProof: verify pub key",
			proof:            invalidProofServicerPubKeyVerify,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Invalid RelayProof: blockchain",
			proof:            invalidProofBlockchain,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Invalid RelayProof: not hosted chain",
			proof:            invalidProofNotHostedBlockchain,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Invalid RelayProof: invalid AAT",
			proof:            invalidProofInvalidAAT,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Invalid RelayProof: client signature",
			proof:            invalidProofClientSignature,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Invalid RelayProof: over max relays",
			proof:            validProof,
			maxRelays:        overMaxRelays,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Invalid RelayProof: over number of chains",
			proof:            validProof,
			maxRelays:        1,
			numOfChains:      overNumberOfChains,
			sessionNodeCount: 0,
			verifyPubKey:     servicerPubKey,
			hb:               hbs,
			hasError:         true,
		},
		{
			name:             "Valid RelayProof",
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
			assert.Equal(t, tt.proof.Validate(tt.maxRelays, tt.numOfChains, tt.sessionNodeCount, tt.hb, tt.verifyPubKey) != nil, tt.hasError)
		})
	}
}

func TestProof_Bytes(t *testing.T) {
	appPubKey := crypto.PublicKey(getRandomPubKey()).String()
	servicerPubKey := crypto.PublicKey(getRandomPubKey()).String()
	clientPubKey := crypto.PublicKey(getRandomPubKey()).String()
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
	assert.Equal(t, pro.Bytes(), proof2.Bytes())
	assert.NotEqual(t, pro.BytesWithSignature(), proof2.BytesWithSignature())
	var p relayProof
	assert.Nil(t, json.Unmarshal(pro.Bytes(), &p))
}

func TestProof_HashAndHashString(t *testing.T) {
	appPubKey := crypto.PublicKey(getRandomPubKey()).String()
	servicerPubKey := crypto.PublicKey(getRandomPubKey()).String()
	clientPubKey := crypto.PublicKey(getRandomPubKey()).String()
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
