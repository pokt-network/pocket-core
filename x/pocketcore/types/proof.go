package types

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
	sdk "github.com/pokt-network/posmint/types"
	"time"
)
// todo broken ** needs to be proofBatches[ProofSummary] = proofs
// one proof per relay
type Proof struct {
	Counter       int       `json:"counter"`       // needed for pseudorandom proof selection
	NodePublicKey string    `json:"nodePublicKey"` // needed to prove the node was the servicer
	Timestamp     time.Time `json:"timestamp"`     // needed to ensure the client signed before the session ended
	Token         AAT       `json:"token"`         // application authentication for client (token)
	Signature     string    `json:"signature"`     // client's signature for the service in hex
}

// many proofs per batch
type Proofs []Proof

// one header per batch
type ProofsHeader struct {
	Chain              string
	SessionBlockHash   string
	SessionBlockHeight int64
	ApplicationPubKey  string
}

type ProofSummary struct { // todo naming conventions
	ProofsHeader
	NodeAddress     sdk.ValAddress
	RelaysCompleted int64
}

// one batch per application
type ProofBatch struct {
	Proofs
}

// all of the batches the node holds
type ProofBatches types.List // list of type ProofBatch

// adds a proof routed to the correct proof batch
func (pbs *ProofBatches) AddProof(proof Proof, sessionBlockIDHex, chain string, sessionBlockHeight int64, maxNumberOfRelays int) error {
	(*types.List)(pbs).Mux.Lock()
	defer (*types.List)(pbs).Mux.Unlock()
	psh := ProofsHeader{
		SessionBlockHash:   sessionBlockIDHex,
		SessionBlockHeight: sessionBlockHeight,
		Chain:              chain,
		ApplicationPubKey:  proof.Token.ApplicationPublicKey,
	}
	if relayBatch, contains := pbs.M[psh]; contains {
		return relayBatch.(ProofBatch).Proofs.AddProof(proof)
	} else {
		return pbs.NewProofBatch(proof, psh, maxNumberOfRelays)
	}
}

func (rbs *ProofBatches) NewProofBatch(proof Proof, proofsHeader ProofsHeader, maxNumberOfRelays int) error {
	rb := ProofBatch{
		ProofsHeader: proofsHeader,
		Proofs:       make([]Proof, maxNumberOfRelays),
	}
	err := rb.AddProof(proof)
	if err != nil {
		return NewBatchCreationErr(ModuleName, err)
	}
	return nil
}

func (rbs *ProofBatches) AddBatch(batch ProofBatch) {
	(*types.List)(rbs).Add(batch.ProofsHeader, batch)
}

func (rbs *ProofBatches) Getbatch(relayBatchHeader ProofsHeader) {
	(*types.List)(rbs).Get(relayBatchHeader)
}

func (rbs *ProofBatches) Removebatch(relayBatchHeader ProofsHeader) {
	(*types.List)(rbs).Remove(relayBatchHeader)
}

func (rbs *ProofBatches) Len() int {
	return (*types.List)(rbs).Count()
}

func (rbs *ProofBatches) Contains(relayBatchHeader ProofsHeader) bool {
	return (*types.List)(rbs).Contains(relayBatchHeader)
}

func (rbs *ProofBatches) Clear() {
	(*types.List)(rbs).Clear()
}

func (p Proof) Hash() ([]byte, error) {
	return crypto.SHA3FromBytes(append([]byte(p.NodePublicKey))), nil // todo this needs to hash everything in this message for sig verification
}

func (p Proof) Validate(maxRelays int64) error {
	// check for negative counter
	if p.Counter < 0 {
		return NegativeICCounterError
	}
	if int64(p.Counter) > maxRelays {
		return MaximumIncrementCounterError
	}
	// decode the client public key todo remove and have sig verification convert to bytes
	cpkBytes, err := hex.DecodeString(p.Token.ClientPublicKey)
	if err != nil {
		return NewClientPubKeyDecodeError(ModuleName, err)
	}
	// validate the service token
	if err := p.Token.Validate(); err != nil {
		return NewInvalidTokenError(ModuleName, err)
	}
	// validate the public key correctness // todo consider abstraction so that this function can also be used for external service proof(s)
	if p.NodePublicKey != "" { // TODO check for current public key *** possibly remove cause may want to reuse for non-self service proofs
		return InvalidNodePubKeyError // the public key is not this nodes, so they would not get paid
	}
	hash, err := p.Hash()
	if err != nil {
		return NewServiceProofHashError(ModuleName, err)
	}
	if !crypto.MockVerifySignature(cpkBytes, hash, []byte(p.Signature)) { // todo change to real sig verification
		return InvalidICSignatureError
	}
	return nil
}

// add proof of work completed (type service proof) to the proof structure
func (ps Proofs) AddProof(p Proof) error {
	if ps == nil || len(ps) == 0 {
		return EmptyProofsError
	}
	// if the increment counter is less than the proof slice
	if len(ps) < p.Counter {
		return InvalidProofSizeError
	}
	// if the proof at index[increment counter] is not empty
	if ps[p.Counter].Signature != "" {
		return DuplicateProofError
	}
	// set proof at index[service proof] = proof of work completed (Service Proof)
	ps[p.Counter] = p
	return nil
}
