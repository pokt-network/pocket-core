package types

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/stretchr/testify/assert"
)

func TestRelayProof_ValidateLocal(t *testing.T) {
	appPrivateKey := GetRandomPrivateKey()
	clientPrivateKey := GetRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	sPK := getRandomPubKey()
	servicerPubKey := sPK.RawString()
	verifyAddr := sPK.Address()
	clientPubKey := clientPrivateKey.PublicKey().RawString()
	ethereum := hex.EncodeToString([]byte{01})
	hbs := HostedBlockchains{
		M: map[string]HostedBlockchain{ethereum: {ID: ethereum, URL: "https://www.google.com:443"}},
	}
	bitcoin := hex.EncodeToString([]byte{02})
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
	// invalid Proof Request ID
	invalidProofRequestHash := validProof
	invalidProofRequestHash.RequestHash = servicerPubKey
	// invalid Proof no client signature
	invalidProofClientSignature := validProof
	invalidProofClientSignature.Signature = hex.EncodeToString(appSignature) // wrong signature
	tests := []struct {
		name             string
		proof            Proof
		maxRelays        int64
		numOfChains      int
		sessionNodeCount int
		verifyPubKey     string
		hb               *HostedBlockchains
		hasError         bool
	}{
		{
			name:             "Invalid Proof: session block",
			proof:            invalidProofSessionBlock,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               &hbs,
			hasError:         true,
		},
		{
			name:             "Invalid Proof: servicer pub key",
			proof:            invalidProofServicerPubKey,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               &hbs,
			hasError:         true,
		},
		{
			name:             "Invalid Proof: blockchain",
			proof:            invalidProofBlockchain,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               &hbs,
			hasError:         true,
		},
		{
			name:             "Invalid Proof: invalid AAT",
			proof:            invalidProofInvalidAAT,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               &hbs,
			hasError:         true,
		},
		{
			name:             "Invalid Proof: client signature",
			proof:            invalidProofClientSignature,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               &hbs,
			hasError:         true,
		},
		{
			name:             "Invalid Proof: invalid request merkleHash from payload",
			proof:            invalidProofRequestHash,
			maxRelays:        5,
			numOfChains:      2,
			sessionNodeCount: 0,
			verifyPubKey:     servicerPubKey,
			hb:               &hbs,
			hasError:         true,
		},
		{
			name:             "Valid Proof",
			proof:            validProof,
			maxRelays:        100,
			numOfChains:      2,
			sessionNodeCount: 5,
			verifyPubKey:     servicerPubKey,
			hb:               &hbs,
			hasError:         false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.proof.(RelayProof).ValidateLocal([]string{getTestSupportedBlockchain()}, tt.sessionNodeCount, 1, sdk.Address(verifyAddr)) != nil, tt.hasError)
		})
	}
}

func TestRelayProof_Bytes(t *testing.T) {
	appPubKey := getRandomPubKey().RawString()
	servicerPubKey := getRandomPubKey().RawString()
	clientPubKey := getRandomPubKey().RawString()
	ethereum := hex.EncodeToString([]byte{01})
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

func TestRelayProof_HashAndHashString(t *testing.T) {
	appPubKey := getRandomPubKey().RawString()
	servicerPubKey := getRandomPubKey().RawString()
	clientPubKey := getRandomPubKey().RawString()
	ethereum := hex.EncodeToString([]byte{01})
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

func TestRelayProof_ValidateBasic(t *testing.T) {
	appPrivateKey := GetRandomPrivateKey()
	clientPrivateKey := GetRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	servicerPubKey := getRandomPubKey().RawString()
	clientPubKey := clientPrivateKey.PublicKey().RawString()
	ethereum := hex.EncodeToString([]byte{01})
	validProof := RelayProof{
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
	// invalid session block height
	invalidSessionBlock := validProof
	invalidSessionBlock.SessionBlockHeight = -1
	// invalid public key format
	invalidPubkeyFormat := validProof
	invalidPubkeyFormat.ServicerPubKey = "abc"
	// invalid blockchain merkleHash
	invalidBCHash := validProof
	invalidBCHash.Blockchain = "abc"
	// invalid request merkleHash
	invalidReqHash := validProof
	invalidReqHash.RequestHash = "abc"
	// invalid Entropy
	invalidEntropy := validProof
	invalidEntropy.Entropy = -1
	// invalid token
	invalidToken := validProof
	invalidToken.Token.ClientPublicKey = "abc"
	// invalid signature
	invalidSig := validProof
	invalidSig.Signature = "abc"
	tests := []struct {
		name     string
		proof    Proof
		hasError bool
	}{
		{
			name:     "valid proof",
			proof:    validProof,
			hasError: false,
		},
		{
			name:     "invalid proof, invalidSessionBlockHeight",
			proof:    invalidSessionBlock,
			hasError: true,
		},
		{
			name:     "invalid proof, invalidPubkeyFormat",
			proof:    invalidPubkeyFormat,
			hasError: true,
		},
		{
			name:     "invalid proof, invalid Blockchain merkleHash",
			proof:    invalidBCHash,
			hasError: true,
		},
		{
			name:     "invalid proof, invalid request merkleHash",
			proof:    invalidReqHash,
			hasError: true,
		},
		{
			name:     "invalid proof, invalid entropy",
			proof:    invalidEntropy,
			hasError: true,
		},
		{
			name:     "invalid proof, invalid token",
			proof:    invalidToken,
			hasError: true,
		},
		{
			name:     "invalid proof, invalid signature",
			proof:    invalidSig,
			hasError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.proof.(RelayProof).ValidateBasic(); (err != nil) != tt.hasError {
				t.Fatalf(err.Error())
			}
		})
	}
}

func TestRelayProof_SessionHeader(t *testing.T) {
	appPrivateKey := GetRandomPrivateKey()
	clientPrivateKey := GetRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	servicerPubKey := getRandomPubKey().RawString()
	clientPubKey := clientPrivateKey.PublicKey().RawString()
	ethereum := hex.EncodeToString([]byte{01})
	validProof := RelayProof{
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
	appSignature, er := appPrivateKey.Sign(validProof.Token.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validProof.Token.ApplicationSignature = hex.EncodeToString(appSignature)
	sh := SessionHeader{
		ApplicationPubKey:  validProof.Token.ApplicationPublicKey,
		Chain:              validProof.Blockchain,
		SessionBlockHeight: validProof.SessionBlockHeight,
	}
	assert.Equal(t, validProof.SessionHeader(), sh)
}

func TestChallengeProofInvalidData_ValidateBasic(t *testing.T) {
	validChallengeProofIVD, _, _, _, _, _, _ := NewValidChallengeProof(t)
	// invalid empty reporter
	invalidEmptyRep := validChallengeProofIVD
	invalidEmptyRep.ReporterAddress = nil
	// invalid signature
	invalidSignature := validChallengeProofIVD
	invalidSignature.MinorityResponse.Signature = ";"
	// mismatched request hashes
	invalidRequestHashes := validChallengeProofIVD
	invalidRequestHashes.MinorityResponse.Proof.RequestHash = "xyz"
	tests := []struct {
		name     string
		proof    ChallengeProofInvalidData
		hasError bool
	}{
		{
			name:     "valid proof",
			proof:    validChallengeProofIVD,
			hasError: false,
		},
		{
			name:     "invalid proof, empty reporter",
			proof:    invalidEmptyRep,
			hasError: true,
		},
		{
			name:     "invalid proof, invalid signature",
			proof:    invalidEmptyRep,
			hasError: true,
		},
		{
			name:     "invalid proof, mismatched request hashes",
			proof:    invalidRequestHashes,
			hasError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.proof.ValidateBasic(); (err != nil) != tt.hasError {
				t.Fatalf(err.Error())
			}
		})
	}
}

func TestChallengeProofInvalidData_ValidateLocal(t *testing.T) {
	validChallengeProofIVD, servicer1PK, servicer2PK, servicer3PK, appPK, _, reporterPK := NewValidChallengeProof(t)
	safeCopyChallenge := func(c ChallengeProofInvalidData) ChallengeProofInvalidData {
		p := ChallengeProofInvalidData{
			MajorityResponses: make([]RelayResponse, 2),
			MinorityResponse:  validChallengeProofIVD.MinorityResponse,
			ReporterAddress:   validChallengeProofIVD.ReporterAddress,
		}
		_ = copy(p.MajorityResponses, c.MajorityResponses)
		return p
	}
	ser1PubKey := servicer1PK.PublicKey()
	ser2PubKey := servicer2PK.PublicKey()
	ser3PubKey := servicer3PK.PublicKey()
	appPubKey := appPK.PublicKey()
	reporterPubKey := reporterPK.PublicKey()
	// invalid challenge Proof duplicate
	invalidProofDup := safeCopyChallenge(validChallengeProofIVD)
	invalidProofDup.MajorityResponses[1] = invalidProofDup.MajorityResponses[0]
	// invalid proof no majority
	invalidProofNoMajority := safeCopyChallenge(validChallengeProofIVD)

	majResp := invalidProofNoMajority.MajorityResponses[0]
	majResp.Response = "foo.bar"
	sig, err := servicer1PK.Sign(majResp.Hash())
	if err != nil {
		t.Fatalf(err.Error())
	}
	majResp.Signature = hex.EncodeToString(sig)
	invalidProofNoMajority.MajorityResponses[0] = majResp

	// invalid proof all majority
	invalidProofAllMajority := safeCopyChallenge(validChallengeProofIVD)

	minResp := invalidProofAllMajority.MinorityResponse
	minResp.Response = invalidProofAllMajority.MajorityResponses[0].Response
	sig, err = servicer3PK.Sign(minResp.Hash())
	if err != nil {
		t.Fatalf(err.Error())
	}
	minResp.Signature = hex.EncodeToString(sig)
	invalidProofAllMajority.MinorityResponse = minResp
	ethereum := hex.EncodeToString([]byte{01})
	sessionNodes := SessionNodes{sdk.Address(ser1PubKey.Address()), sdk.Address(ser2PubKey.Address()), sdk.Address(ser3PubKey.Address()), sdk.Address(reporterPubKey.Address()), sdk.Address(appPubKey.Address())}
	tests := []struct {
		name                 string
		proof                ChallengeProofInvalidData
		maxRelays            sdk.BigInt
		supportedBlockchains []string
		sessionNodes         SessionNodes
		reporterAddress      sdk.Address
		hasError             bool
	}{
		{
			name:                 "valid proof",
			proof:                validChallengeProofIVD,
			maxRelays:            sdk.NewInt(100000),
			supportedBlockchains: []string{ethereum},
			sessionNodes:         sessionNodes,
			reporterAddress:      sdk.Address(reporterPubKey.Address()),
			hasError:             false,
		},
		{
			name:                 "invalidProof, reporter (self) not in session",
			proof:                validChallengeProofIVD,
			maxRelays:            sdk.NewInt(10000),
			supportedBlockchains: []string{ethereum},
			sessionNodes:         sessionNodes,
			reporterAddress:      sdk.Address([]byte("fake")),
			hasError:             true,
		},
		{
			name:                 "invalidProof, duplicate",
			proof:                invalidProofDup,
			maxRelays:            sdk.NewInt(10000),
			supportedBlockchains: []string{ethereum},
			sessionNodes:         sessionNodes,
			reporterAddress:      sdk.Address(reporterPubKey.Address()),
			hasError:             true,
		},
		{
			name:                 "invalidProof, no majority",
			proof:                invalidProofNoMajority,
			maxRelays:            sdk.NewInt(10000),
			supportedBlockchains: []string{ethereum},
			sessionNodes:         sessionNodes,
			reporterAddress:      sdk.Address(reporterPubKey.Address()),
			hasError:             true,
		},
		{
			name:                 "invalidProof, all majority",
			proof:                invalidProofAllMajority,
			maxRelays:            sdk.NewInt(10000),
			supportedBlockchains: []string{ethereum},
			sessionNodes:         sessionNodes,
			reporterAddress:      sdk.Address(reporterPubKey.Address()),
			hasError:             true,
		},
		{
			name:                 "invalidProof, proof overflow",
			proof:                validChallengeProofIVD,
			maxRelays:            sdk.ZeroInt(),
			supportedBlockchains: []string{ethereum},
			sessionNodes:         sessionNodes,
			reporterAddress:      sdk.Address(reporterPubKey.Address()),
			hasError:             true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := SessionHeader{
				ApplicationPubKey:  tt.proof.MinorityResponse.Proof.Token.ApplicationPublicKey,
				Chain:              tt.proof.MinorityResponse.Proof.Blockchain,
				SessionBlockHeight: tt.proof.MinorityResponse.Proof.SessionBlockHeight,
			}
			if err := tt.proof.ValidateLocal(h, tt.maxRelays, tt.supportedBlockchains, 5, tt.sessionNodes, tt.reporterAddress); (err != nil) != tt.hasError {
				fmt.Println(tt.name)
				fmt.Println(err)
				t.Fatalf(err.Error())
			}
		})
	}
}

func NewValidChallengeProof(t *testing.T) (challenge ChallengeProofInvalidData, ser1 crypto.PrivateKey, ser2 crypto.PrivateKey, ser3 crypto.PrivateKey, app crypto.PrivateKey, cli crypto.PrivateKey, repor crypto.PrivateKey) {
	appPrivateKey := GetRandomPrivateKey()
	servicerPrivKey1 := GetRandomPrivateKey()
	servicerPrivKey2 := GetRandomPrivateKey()
	servicerPrivKey3 := GetRandomPrivateKey()
	clientPrivateKey := GetRandomPrivateKey()
	appPubKey := appPrivateKey.PublicKey().RawString()
	servicerPubKey := servicerPrivKey1.PublicKey().RawString()
	servicerPubKey2 := servicerPrivKey2.PublicKey().RawString()
	servicerPubKey3 := servicerPrivKey3.PublicKey().RawString()
	reporterPrivKey := GetRandomPrivateKey()
	reporterPubKey := reporterPrivKey.PublicKey()
	reporterAddr := reporterPubKey.Address()
	clientPubKey := clientPrivateKey.PublicKey().RawString()
	ethereum := hex.EncodeToString([]byte{01})
	validProof := RelayProof{
		Entropy:            0,
		SessionBlockHeight: 1,
		ServicerPubKey:     servicerPubKey,
		RequestHash:        clientPubKey, // fake
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
	// valid proof 2
	validProof2 := RelayProof{
		Entropy:            0,
		SessionBlockHeight: 1,
		ServicerPubKey:     servicerPubKey2,
		RequestHash:        clientPubKey, // fake
		Blockchain:         ethereum,
		Token: AAT{
			Version:              "0.0.1",
			ApplicationPublicKey: appPubKey,
			ClientPublicKey:      clientPubKey,
			ApplicationSignature: "",
		},
		Signature: "",
	}
	appSignature, er = appPrivateKey.Sign(validProof2.Token.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validProof2.Token.ApplicationSignature = hex.EncodeToString(appSignature)
	clientSignature, er = clientPrivateKey.Sign(validProof2.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validProof2.Signature = hex.EncodeToString(clientSignature)
	// valid proof 3
	validProof3 := RelayProof{
		Entropy:            0,
		SessionBlockHeight: 1,
		ServicerPubKey:     servicerPubKey3,
		RequestHash:        clientPubKey, // fake
		Blockchain:         ethereum,
		Token: AAT{
			Version:              "0.0.1",
			ApplicationPublicKey: appPubKey,
			ClientPublicKey:      clientPubKey,
			ApplicationSignature: "",
		},
		Signature: "",
	}
	appSignature, er = appPrivateKey.Sign(validProof3.Token.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validProof3.Token.ApplicationSignature = hex.EncodeToString(appSignature)
	clientSignature, er = clientPrivateKey.Sign(validProof3.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	validProof3.Signature = hex.EncodeToString(clientSignature)
	// create responses
	majorityResponsePayload := `{"id":67,"jsonrpc":"2.0","result":"Mist/v0.9.3/darwin/go1.4.1"}`
	minorityResponsePayload := `{"id":67,"jsonrpc":"2.0","result":"Mist/v0.9.3/darwin/go1.4.2"}`
	// majority response 1
	majResp1 := RelayResponse{
		Signature: "",
		Response:  majorityResponsePayload,
		Proof:     validProof,
	}
	sig, er := servicerPrivKey1.Sign(majResp1.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	majResp1.Signature = hex.EncodeToString(sig)
	// majority response 2
	majResp2 := RelayResponse{
		Signature: "",
		Response:  majorityResponsePayload,
		Proof:     validProof2,
	}
	sig, er = servicerPrivKey2.Sign(majResp2.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	majResp2.Signature = hex.EncodeToString(sig)
	// minority response
	minResp := RelayResponse{
		Signature: "",
		Response:  minorityResponsePayload,
		Proof:     validProof3,
	}
	sig, er = servicerPrivKey3.Sign(minResp.Hash())
	if er != nil {
		t.Fatalf(er.Error())
	}
	minResp.Signature = hex.EncodeToString(sig)
	// create valid challenge proof
	return ChallengeProofInvalidData{
		MajorityResponses: []RelayResponse{
			majResp1,
			majResp2,
		},
		MinorityResponse: minResp,
		ReporterAddress:  sdk.Address(reporterAddr),
	}, servicerPrivKey1, servicerPrivKey2, servicerPrivKey3, appPrivateKey, clientPrivateKey, reporterPrivKey
}

func TestChallengeProofInvalidData_SessionHeader(t *testing.T) {
	c, _, _, _, _, _, _ := NewValidChallengeProof(t)
	assert.Equal(t, c.SessionHeader(), SessionHeader{
		ApplicationPubKey:  c.MinorityResponse.Proof.Token.ApplicationPublicKey,
		Chain:              c.MinorityResponse.Proof.Blockchain,
		SessionBlockHeight: c.MinorityResponse.Proof.SessionBlockHeight,
	})
}
