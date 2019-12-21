package types

import (
	"encoding/hex"
	"encoding/json"
	"github.com/pokt-network/merkle"
	sdk "github.com/pokt-network/posmint/types"
	"math"
	"sync"
)

// proof of relay per application
type ProofOfRelay struct {
	SessionHeader `json:"header"`
	TotalRelays   int64   `json:"total_relays"`
	Proofs        []Proof `json:"proofs"` // slice[index] -> ProofsMap
}

// every proof the node holds
type AllProofs struct {
	M map[string]ProofOfRelay `json:"proofs"` // map[porkey] -> ProofOfRelay
	l sync.Mutex
}

var (
	globalAllProofs *AllProofs // holds every proof of the node
	apOnce          sync.Once  // ensure only made once
)

// get all proofs
func GetAllProofs() *AllProofs {
	apOnce.Do(func() {
		if globalAllProofs == nil {
			ap := AllProofs{M: make(map[string]ProofOfRelay)}
			globalAllProofs = &ap
		}
	})
	return globalAllProofs
}

func (por ProofOfRelay) GenerateMerkleRoot() (root []byte) {
	var data [][]byte
	// todo this can be way more efficient (store the hashes)
	for i, proof := range por.Proofs {
		proof.Index = int64(i)
		por.Proofs[i] = proof
		data = append(data, proof.Bytes())
	}
	return merkle.GenerateRoot(data)
}

func (por ProofOfRelay) GenerateProof(index int) merkle.Proof {
	var data [][]byte
	// todo this can be way more efficient (store the hashes)
	for i, proof := range por.Proofs {
		proof.Index = int64(i)
		por.Proofs[i] = proof
		data = append(data, proof.Bytes())
	}
	return merkle.GenerateProof(data, index)
}

func (por ProofOfRelay) VerifyProof(root, leaf []byte, proof merkle.Proof) bool {
	return merkle.VerifyProof(root, leaf, proof)
}

// add the proof to the AllProofs object
func (ap AllProofs) AddProof(header SessionHeader, p Proof, maxRelays int64) sdk.Error {
	var por = ProofOfRelay{}
	// generate the key for this specific proof
	key := header.HashString()
	// lock the shared data
	ap.l.Lock()
	defer ap.l.Unlock()
	if _, found := ap.M[key]; found {
		// if proof already stored in allProofs
		por = ap.M[key]
	} else {
		// if proof is not already stored, initialize all
		por.SessionHeader = header
		por.Proofs = make([]Proof, maxRelays)
		por.TotalRelays = 0
	}
	// add proof to the proofs object
	por.Proofs = append(por.Proofs, p)
	// increment total relay count
	por.TotalRelays = por.TotalRelays + 1
	// update POR
	ap.M[key] = por
	return nil
}

func (ap AllProofs) GetTotalRelays(header SessionHeader) int64 {
	// lock the shared data
	ap.l.Lock()
	defer ap.l.Unlock()
	// return the proofs object, corresponding to the header
	return ap.M[header.HashString()].TotalRelays
}

// retrieve the single proof from the all proofs object
func (ap AllProofs) GetProof(header SessionHeader, index int) *Proof {
	// lock the shared data
	ap.l.Lock()
	defer ap.l.Unlock()
	// return the proofs object, corresponding to the header
	por := ap.M[header.HashString()].Proofs
	// do a nil check before indexing
	if por == nil {
		return nil
	}
	// return the proof at specific index
	return &por[index]
}

// retrieve the proofs from the all proofs object
func (ap AllProofs) GetProofs(header SessionHeader) []Proof {
	// lock the shared data
	ap.l.Lock()
	defer ap.l.Unlock()
	// return the proofs object, corresponding to the header
	return ap.M[header.HashString()].Proofs
}

// delete proofs from the all proofs object
func (ap AllProofs) DeleteProofs(header SessionHeader) {
	// lock the shared data
	ap.l.Lock()
	defer ap.l.Unlock()
	// delete the value corresponding to the header
	delete(ap.M, header.HashString())
}

// proof per relay
type Proof struct {
	Index              int64  `json:"index"`
	SessionBlockHeight int64  `json:"session_block_height"`
	ServicerPubKey     string `json:"servicer_pub_key"`
	Blockchain         string `json:"blockchain"`
	Token              AAT    `json:"aat"`
	Signature          string `json:"signature"`
}

func (p Proof) Validate(maxRelays int64, numberOfChains, sessionNodeCount int, hb HostedBlockchains, verifyPubKey string) sdk.Error {
	// validate the session block height
	if p.SessionBlockHeight < 0 {
		return NewInvalidBlockHeightError(ModuleName)
	}
	// validate blockchain
	if err := HashVerification(p.Blockchain); err != nil {
		return err
	}
	// validate not over service
	totalRelays := GetAllProofs().GetTotalRelays(SessionHeader{
		ApplicationPubKey:  p.Token.ApplicationPublicKey,
		Chain:              p.Blockchain,
		SessionBlockHeight: p.SessionBlockHeight,
	})
	if totalRelays >= int64(math.Ceil(float64(maxRelays)/float64(numberOfChains))/(float64(sessionNodeCount))) {
		return NewOverServiceError(ModuleName)
	}
	// validate the public key correctness
	if p.ServicerPubKey != verifyPubKey {
		return NewInvalidNodePubKeyError(ModuleName) // the public key is not this nodes, so they would not get paid
	}
	// ensure the blockchain is supported
	if !hb.ContainsFromString(p.Blockchain) {
		return NewUnsupportedBlockchainNodeError(ModuleName)
	}
	// validate the proof public key format
	if err := PubKeyVerification(p.ServicerPubKey); err != nil {
		return NewInvalidNodePubKeyError(ModuleName)
	}
	// validate the verify public key format
	if err := PubKeyVerification(verifyPubKey); err != nil {
		return NewInvalidNodePubKeyError(ModuleName)
	}
	// validate the service token
	if err := p.Token.Validate(); err != nil {
		return NewInvalidTokenError(ModuleName, err)
	}
	return SignatureVerification(p.Token.ClientPublicKey, p.HashString(), p.Signature)
}

// structure used to json marshal the proof
type proof struct {
	Index              int64  `json:"index"`
	SessionBlockHeight int64  `json:"session_block_height"`
	ServicerPubKey     string `json:"servicer_pub_key"`
	Blockchain         string `json:"blockchain"`
	Signature          string `json:"signature"`
	Token              string `json:"token"`
}

// hash the proof bytes
func (p Proof) Hash() []byte {
	res := p.Bytes()
	return Hash(res)
}

// hex encode the proof hash
func (p Proof) HashString() string {
	return hex.EncodeToString(p.Hash())
}

// convert the proof to bytes
func (p Proof) Bytes() []byte {
	res, err := json.Marshal(proof{
		Index:              p.Index,
		ServicerPubKey:     p.ServicerPubKey,
		Blockchain:         p.Blockchain,
		SessionBlockHeight: p.SessionBlockHeight,
		Signature:          "", // omit the signature
		Token:              p.Token.HashString(),
	})
	if err != nil {
		panic(err)
	}
	return res
}

type MerkleProof merkle.Proof

// structure used to store the proof after verification
type StoredProof struct {
	SessionHeader
	ServicerAddress string
	TotalRelays     int64
}
