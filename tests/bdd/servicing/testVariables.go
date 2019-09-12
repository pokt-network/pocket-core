package servicing

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/service"
)

var (
	hostedBlockchains = map[string]struct{}{"eth":{}}

	validTokenVersion = "0.0.1"

	validAppPubKey = hex.EncodeToString(crypto.NewPrivateKey().PubKey().Bytes())

	validClientPublicKey = hex.EncodeToString(crypto.NewPrivateKey().PubKey().Bytes())

	invalidAppSignature = crypto.Signature{}

	invalidClientSignature = crypto.Signature{}

	validAppSignature = []byte("todosignature") // todo

	validClientSignature = []byte("todosignature") // todo

	validBlockchain = []byte("eth") //todo change to amino

	unsupportedBlockchain = []byte("btc") // todo change to amino

	unsupportedVersion = "0.0.0"

	validICCount = 100

	invalidICCount = -1

	validData = "jsonrpc\":\"2.0\",\"method\":\"net_peerCount\",\"params\":[],\"id\":74\""

	validTokenMessage = types.AATMessage{
		ApplicationPublicKey: validAppPubKey,
		ClientPublicKey:      validClientPublicKey,
	}
	missingAppPublicKeyMessage = types.AATMessage{
		ApplicationPublicKey: "",
		ClientPublicKey:      validClientPublicKey,
	}

	missingClientPubKeyMessage = types.AATMessage{
		ApplicationPublicKey: validAppPubKey,
		ClientPublicKey:      "",
	}

	unsupportedTokenVersion = service.ServiceToken{
		Version:    types.AATVersion(unsupportedVersion),
		AATMessage: validTokenMessage,
		Signature:  validAppSignature,
	}

	missingTokenVersion = service.ServiceToken{
		Version:    "",
		AATMessage: validTokenMessage,
		Signature:  validAppSignature,
	}

	missingApplicationPublicKeyTokenMessage = service.ServiceToken{
		Version:    types.AATVersion(validTokenVersion),
		AATMessage: missingAppPublicKeyMessage,
		Signature:  validAppSignature,
	}

	missingClientPublicKeyTokenMessage = service.ServiceToken{
		Version:    types.AATVersion(validTokenVersion),
		AATMessage: missingClientPubKeyMessage,
		Signature:  validAppSignature,
	}

	invalidTokenSignature = service.ServiceToken{
		Version:    types.AATVersion(validTokenVersion),
		AATMessage: validTokenMessage,
		Signature:  invalidAppSignature,
	}

	validToken = service.ServiceToken{
		Version:    types.AATVersion(validTokenVersion),
		AATMessage: validTokenMessage,
		Signature:  validAppSignature,
	}

	validPayload = service.ServicePayload{
		Data: service.ServiceData(validData),
		HttpServicePayload: service.HttpServicePayload{
			Method: "",
			Path:   "",
		},
	}

	validIncrementCounter = service.IncrementCounter{
		Counter:   validICCount,
		Signature: validClientSignature,
	}

	invalidIncrementCounterCount = service.IncrementCounter{
		Counter:   invalidICCount,
		Signature: validClientSignature,
	}

	invalidIncrementCounterSignature = service.IncrementCounter{
		Counter:   validICCount,
		Signature: invalidClientSignature,
	}

	relayMissingBlockchain = service.Relay{
		Blockchain:       nil,
		Payload:          validPayload,
		ServiceToken:     validToken,
		IncrementCounter: validIncrementCounter,
	}

	relayMissingPayload = service.Relay{
		Blockchain:       validBlockchain,
		Payload:          service.ServicePayload{},
		ServiceToken:     validToken,
		IncrementCounter: validIncrementCounter,
	}

	relayUnsupportedBlockchain = service.Relay{
		Blockchain:       unsupportedBlockchain,
		Payload:          validPayload,
		ServiceToken:     validToken,
		IncrementCounter: validIncrementCounter,
	}

	relayUnsupportedTokenVersion = service.Relay{
		Blockchain:       validBlockchain,
		Payload:          validPayload,
		ServiceToken:     unsupportedTokenVersion,
		IncrementCounter: validIncrementCounter,
	}

	relayMissingTokenVersion = service.Relay{
		Blockchain:       validBlockchain,
		Payload:          validPayload,
		ServiceToken:     missingTokenVersion,
		IncrementCounter: validIncrementCounter,
	}

	relayMissingTokenAppPubKey = service.Relay{
		Blockchain:       validBlockchain,
		Payload:          validPayload,
		ServiceToken:     missingApplicationPublicKeyTokenMessage,
		IncrementCounter: validIncrementCounter,
	}

	relayMissingTokenCliPubKey = service.Relay{
		Blockchain:       validBlockchain,
		Payload:          validPayload,
		ServiceToken:     missingClientPublicKeyTokenMessage,
		IncrementCounter: validIncrementCounter,
	}

	relayInvalidTokenSignature = service.Relay{
		Blockchain:       validBlockchain,
		Payload:          validPayload,
		ServiceToken:     invalidTokenSignature,
		IncrementCounter: validIncrementCounter,
	}

	relayInvalidICCount = service.Relay{
		Blockchain:       validBlockchain,
		Payload:          validPayload,
		ServiceToken:     validToken,
		IncrementCounter: invalidIncrementCounterCount,
	}

	relayInvalidICSignature = service.Relay{
		Blockchain:       validBlockchain,
		Payload:          validPayload,
		ServiceToken:     validToken,
		IncrementCounter: invalidIncrementCounterSignature,
	}

	validRelay = service.Relay{
		Blockchain:       validBlockchain,
		Payload:          validPayload,
		ServiceToken:     validToken,
		IncrementCounter: validIncrementCounter,
	}
)
