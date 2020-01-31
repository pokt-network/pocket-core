package types

import (
	"encoding/hex"
	appsType "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	nodesTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"reflect"
	"testing"
	"time"
)

func TestRelay_Validate(t *testing.T) {
	clientPrivateKey := getRandomPrivateKey()
	clientPubKey := clientPrivateKey.PublicKey().RawString()
	appPrivateKey := getRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	npk := getRandomPubKey()
	nodePubKey := npk.RawString()
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
	validRelay := Relay{
		Payload: Payload{
			Data:    "{\"jsonrpc\":\"2.0\",\"method\":\"web3_clientVersion\",\"params\":[],\"id\":67}",
			Method:  "",
			Path:    "",
			Headers: nil,
		},
		Proof: RelayProof{
			Entropy:            1,
			SessionBlockHeight: 1,
			ServicerPubKey:     nodePubKey,
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
	appSig, er := appPrivateKey.Sign(validRelay.Proof.Token.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validRelay.Proof.Token.ApplicationSignature = hex.EncodeToString(appSig)
	clientSig, er := clientPrivateKey.Sign(validRelay.Proof.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validRelay.Proof.Signature = hex.EncodeToString(clientSig)
	// invalid payload empty data and method
	invalidPayloadEmpty := validRelay
	invalidPayloadEmpty.Payload.Data = ""
	selfNode := nodesTypes.Validator{
		Address:                 sdk.Address(npk.Address()),
		PublicKey:               npk,
		Jailed:                  false,
		Status:                  sdk.Staked,
		Chains:                  []string{ethereum, bitcoin},
		ServiceURL:              "www.google.com",
		StakedTokens:            sdk.NewInt(100000),
		UnstakingCompletionTime: time.Time{},
	}
	var allNodes []exported.ValidatorI
	for i := 0; i < 4; i++ {
		pubKey := getRandomPubKey()
		allNodes = append(allNodes, nodesTypes.Validator{
			Address:                 sdk.Address(pubKey.Address()),
			PublicKey:               pubKey,
			Jailed:                  false,
			Status:                  sdk.Staked,
			Chains:                  []string{ethereum, bitcoin},
			ServiceURL:              "www.google.com",
			StakedTokens:            sdk.NewInt(100000),
			UnstakingCompletionTime: time.Time{},
		})
	}
	var noEthereumNodes []exported.ValidatorI
	for i := 0; i < 4; i++ {
		pubKey := getRandomPubKey()
		noEthereumNodes = append(noEthereumNodes, nodesTypes.Validator{
			Address:                 sdk.Address(pubKey.Address()),
			PublicKey:               pubKey,
			Jailed:                  false,
			Status:                  sdk.Staked,
			Chains:                  []string{bitcoin},
			ServiceURL:              "www.google.com",
			StakedTokens:            sdk.NewInt(100000),
			UnstakingCompletionTime: time.Time{},
		})
	}
	allNodes = append(allNodes, selfNode)
	noEthereumNodes = append(noEthereumNodes, selfNode)
	hb := HostedBlockchains{
		M: map[string]HostedBlockchain{ethereum: {
			Hash: ethereum,
			URL:  "www.google.com",
		}},
	}
	hbNotSupported := HostedBlockchains{
		M: map[string]HostedBlockchain{bitcoin: {
			Hash: bitcoin,
			URL:  "www.google.com",
		}},
	}
	pubKey := getRandomPubKey()
	app := appsType.Application{
		Address:                 sdk.Address(pubKey.Address()),
		PublicKey:               pubKey,
		Jailed:                  false,
		Status:                  sdk.Staked,
		Chains:                  []string{ethereum},
		StakedTokens:            sdk.NewInt(1000),
		MaxRelays:               sdk.NewInt(1000),
		UnstakingCompletionTime: time.Time{},
	}
	tests := []struct {
		name     string
		relay    Relay
		node     nodesTypes.Validator
		app      appsType.Application
		allNodes []exported.ValidatorI
		hb       HostedBlockchains
		hasError bool
	}{
		{
			name:     "valid relay",
			relay:    validRelay,
			node:     selfNode,
			app:      app,
			allNodes: allNodes,
			hb:       hb,
			hasError: false,
		},
		{
			name:     "invalid relay: payload empty",
			relay:    invalidPayloadEmpty,
			node:     selfNode,
			app:      app,
			allNodes: allNodes,
			hb:       hb,
			hasError: true,
		},
		{
			name:     "invalid relay: unsupported blockchain",
			relay:    validRelay,
			node:     selfNode,
			app:      app,
			allNodes: allNodes,
			hb:       hbNotSupported,
			hasError: true,
		},
		{
			name:     "invalid relay: not enough service nodes",
			relay:    validRelay,
			node:     selfNode,
			app:      app,
			allNodes: noEthereumNodes,
			hb:       hb,
			hasError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.relay.Validate(newContext(t, false), tt.node,
				tt.hb, 1, 5, tt.allNodes, tt.app) != nil, tt.hasError)
		})
	}
}

func TestRelay_Execute(t *testing.T) {
	clientPrivateKey := getRandomPrivateKey()
	clientPubKey := clientPrivateKey.PublicKey().RawString()
	appPrivateKey := getRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	npk := getRandomPubKey()
	nodePubKey := npk.RawString()
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
	validRelay := Relay{
		Payload: Payload{
			Data:    "foo",
			Method:  "POST",
			Path:    "",
			Headers: nil,
		},
		Proof: RelayProof{
			Entropy:            1,
			SessionBlockHeight: 1,
			ServicerPubKey:     nodePubKey,
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
	defer gock.Off() // Flush pending mocks after test execution

	gock.New("https://server.com").
		Post("/relay/").
		Reply(200).
		BodyString("bar")

	hb := HostedBlockchains{
		M: map[string]HostedBlockchain{ethereum: {
			Hash: ethereum,
			URL:  "https://server.com/relay/",
		}},
	}
	response, err := validRelay.Execute(hb)
	assert.True(t, err == nil)
	assert.Equal(t, response, "bar")
}

func TestRelay_HandleProof(t *testing.T) {
	clientPrivateKey := getRandomPrivateKey()
	clientPubKey := clientPrivateKey.PublicKey().RawString()
	appPrivateKey := getRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	npk := getRandomPubKey()
	nodePubKey := npk.RawString()
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
	validRelay := Relay{
		Payload: Payload{
			Data:    "foo",
			Method:  "POST",
			Path:    "",
			Headers: nil,
		},
		Proof: RelayProof{
			Entropy:            1,
			SessionBlockHeight: 1,
			ServicerPubKey:     nodePubKey,
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
	err = validRelay.HandleProof(newContext(t, false), 1)
	assert.Nil(t, err)
	res := GetAllInvoices().GetProof(SessionHeader{
		ApplicationPubKey:  appPubKey,
		Chain:              ethereum,
		SessionBlockHeight: 1,
	}, 0)
	assert.True(t, reflect.DeepEqual(validRelay.Proof, res))
}

func TestRelayResponse_BytesAndHash(t *testing.T) {
	nodePrivKey := getRandomPrivateKey()
	nodePubKey := nodePrivKey.PublicKey().RawString()
	appPrivKey := getRandomPrivateKey()
	appPublicKey := appPrivKey.PublicKey().RawString()
	cliPrivKey := getRandomPrivateKey()
	cliPublicKey := cliPrivKey.PublicKey().RawString()
	relayResp := RelayResponse{
		Signature: "",
		Response:  "foo",
		Proof: RelayProof{
			Entropy:            230942034,
			SessionBlockHeight: 1,
			ServicerPubKey:     nodePubKey,
			Blockchain:         hex.EncodeToString(hash([]byte("foo"))),
			Token: AAT{
				Version:              "0.0.1",
				ApplicationPublicKey: appPublicKey,
				ClientPublicKey:      cliPublicKey,
				ApplicationSignature: "",
			},
			Signature: "",
		},
	}
	appSig, err := appPrivKey.Sign(relayResp.Proof.Token.Hash())
	if err != nil {
		t.Fatalf(err.Error())
	}
	relayResp.Proof.Token.ApplicationSignature = hex.EncodeToString(appSig)
	assert.NotNil(t, relayResp.Hash())
	assert.Equal(t, hex.EncodeToString(relayResp.Hash()), relayResp.HashString())
	storedHashString := relayResp.HashString()
	nodeSig, err := nodePrivKey.Sign(relayResp.Hash())
	if err != nil {
		t.Fatalf(err.Error())
	}
	relayResp.Signature = hex.EncodeToString(nodeSig)
	assert.Equal(t, storedHashString, relayResp.HashString())
}
