package servicing

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/tests/fixtures"
	"github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/service"
	"path/filepath"
)

const (
	GOODENDPOINT = "http://a.fake"
	BADENDPOINT  = "http://b.fake"
	GOODRESULT   = "Mist/v0.9.3/darwin/go1.4.1"
)

var (
	chainsfile, _       = filepath.Abs("../../fixtures/chains.json")
	brokenchainsfile, _ = filepath.Abs("../../fixtures/legacy/brokenChains.json")
	hc                  = types.HostedBlockchains{M: map[interface{}]interface{}{
		hex.EncodeToString(validBlockchain): types.HostedBlockchain{
			Hash: hex.EncodeToString(validBlockchain),
			URL:  GOODENDPOINT,
		},
		hex.EncodeToString(validBlockchain2): types.HostedBlockchain{
			Hash: hex.EncodeToString(validBlockchain2),
			URL:  BADENDPOINT,
		}}}

	hostedBlockchains = service.ServiceBlockchains(hc)

	validTokenVersion = "0.0.1"

	_, appPubKey, _ = crypto.NewKeypair()

	_, cliPubKey, _ = crypto.NewKeypair()

	validAppPubKey = hex.EncodeToString(appPubKey.Bytes())

	validNodePubKey = "" //todo

	invalidNodePubKey = "" //todo

	validClientPublicKey = hex.EncodeToString(cliPubKey.Bytes())

	invalidAppSignature = crypto.Signature{}

	invalidClientSignature = "" // todo

	validAppSignature = []byte("todosignature") // todo

	validClientSignature = "todosignature" // todo

	validBlockchain = fixtures.GenerateNonNativeBlockchainFromTicker("eth")

	validBlockchain2 = fixtures.GenerateNonNativeBlockchainFromTicker("btc")

	unsupportedBlockchain = []byte("aion") // todo change to amino

	unsupportedVersion = "0.0.0"

	validICCount = 100

	invalidICCount = -1

	validData = `"jsonrpc":"2.0","method":"web3_clientVersion","params":[],"id":67`

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

	unsupportedTokenVersion = service.ServiceCertificate{
		Signature: validClientSignature,
		ServiceCertificatePayload: service.ServiceCertificatePayload{
			Counter:       validICCount,
			NodePublicKey: validNodePubKey,
			ServiceToken: service.ServiceToken{
				Version:    types.AATVersion(unsupportedVersion),
				AATMessage: validTokenMessage,
				Signature:  validAppSignature,
			}}}

	missingTokenVersion = service.ServiceCertificate{
		Signature: validClientSignature,
		ServiceCertificatePayload: service.ServiceCertificatePayload{
			Counter:       validICCount,
			NodePublicKey: validNodePubKey,
			ServiceToken: service.ServiceToken{
				Version:    "",
				AATMessage: validTokenMessage,
				Signature:  validAppSignature,
			}}}

	missingApplicationPublicKeyTokenMessage = service.ServiceCertificate{
		Signature: validClientSignature,
		ServiceCertificatePayload: service.ServiceCertificatePayload{
			Counter:       validICCount,
			NodePublicKey: validNodePubKey,
			ServiceToken: service.ServiceToken{
				Version:    types.AATVersion(validTokenVersion),
				AATMessage: missingAppPublicKeyMessage,
				Signature:  validAppSignature,
			}}}

	missingClientPublicKeyTokenMessage = service.ServiceCertificate{
		Signature: validClientSignature,
		ServiceCertificatePayload: service.ServiceCertificatePayload{
			Counter:       validICCount,
			NodePublicKey: validNodePubKey,
			ServiceToken: service.ServiceToken{
				Version:    types.AATVersion(validTokenVersion),
				AATMessage: missingClientPubKeyMessage,
				Signature:  validAppSignature,
			}}}

	invalidTokenSignature = service.ServiceCertificate{
		ServiceCertificatePayload: service.ServiceCertificatePayload{
			Counter:       validICCount,
			NodePublicKey: validNodePubKey,
			ServiceToken: service.ServiceToken{
				Version:    types.AATVersion(validTokenVersion),
				AATMessage: validTokenMessage,
				Signature:  invalidAppSignature,
			}}, Signature: validClientSignature,
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

	validServiceAuthentication = service.ServiceCertificate{
		ServiceCertificatePayload: service.ServiceCertificatePayload{
			Counter:       validICCount,
			NodePublicKey: validNodePubKey,
			ServiceToken:  validToken,
		},
		Signature: validClientSignature,
	}

	invalidServiceAuthenticationCounter = service.ServiceCertificate{
		ServiceCertificatePayload: service.ServiceCertificatePayload{
			Counter:       invalidICCount,
			NodePublicKey: validNodePubKey,
			ServiceToken:  validToken,
		},
		Signature: validClientSignature,
	}

	invalidServiceAuthenticationSignature = service.ServiceCertificate{
		ServiceCertificatePayload: service.ServiceCertificatePayload{
			Counter:       validICCount,
			NodePublicKey: validNodePubKey,
			ServiceToken:  validToken,
		},
		Signature: invalidClientSignature,
	}

	relayMissingBlockchain = service.Relay{
		Blockchain:         nil,
		Payload:            validPayload,
		ServiceCertificate: validServiceAuthentication,
	}

	relayMissingPayload = service.Relay{
		Blockchain:         validBlockchain,
		Payload:            service.ServicePayload{},
		ServiceCertificate: validServiceAuthentication,
	}

	relayUnsupportedBlockchain = service.Relay{
		Blockchain:         unsupportedBlockchain,
		Payload:            validPayload,
		ServiceCertificate: validServiceAuthentication,
	}

	relayUnsupportedTokenVersion = service.Relay{
		Blockchain:         validBlockchain,
		Payload:            validPayload,
		ServiceCertificate: unsupportedTokenVersion,
	}

	relayMissingTokenVersion = service.Relay{
		Blockchain:         validBlockchain,
		Payload:            validPayload,
		ServiceCertificate: missingTokenVersion,
	}

	relayMissingTokenAppPubKey = service.Relay{
		Blockchain:         validBlockchain,
		Payload:            validPayload,
		ServiceCertificate: missingApplicationPublicKeyTokenMessage,
	}

	relayMissingTokenCliPubKey = service.Relay{
		Blockchain:         validBlockchain,
		Payload:            validPayload,
		ServiceCertificate: missingClientPublicKeyTokenMessage,
	}

	relayInvalidTokenSignature = service.Relay{
		Blockchain:         validBlockchain,
		Payload:            validPayload,
		ServiceCertificate: invalidTokenSignature,
	}

	relayInvalidICCount = service.Relay{
		Blockchain:         validBlockchain,
		Payload:            validPayload,
		ServiceCertificate: invalidServiceAuthenticationCounter,
	}

	relayInvalidICSignature = service.Relay{
		Blockchain:         validBlockchain,
		Payload:            validPayload,
		ServiceCertificate: invalidServiceAuthenticationSignature,
	}

	validEthRelay = service.Relay{
		Blockchain:         validBlockchain,
		Payload:            validPayload,
		ServiceCertificate: validServiceAuthentication,
	}

	validBtcRelay = service.Relay{
		Blockchain:         validBlockchain2,
		Payload:            validPayload,
		ServiceCertificate: validServiceAuthentication,
	}
)
