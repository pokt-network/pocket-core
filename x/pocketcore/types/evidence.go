package types

import (
	"encoding/hex"
	sdk "github.com/pokt-network/posmint/types"
	"sync"
)

var (
	globalEvidenceMap *EvidenceMap // holds all evidence from relays and challenges
	evidenceMapOnce   sync.Once    // ensure only made once
)

// Proof of relay per application
type Evidence struct {
	SessionHeader `json:"evidence_header"` // the session h serves as an identifier for the evidence
	NumOfProofs   int64                    `json:"num_of_¬proofs"` // the total number of proofs in the evidence
	Proofs        []Proof                  `json:"proofs"`         // a slice of Proof objects (Proof per relay or challenge)
}

// generate the merkle root of an evidence
func (e *Evidence) GenerateMerkleRoot() (root HashSum) {
	root, sortedProofs := GenerateRoot(e.Proofs)
	e.Proofs = sortedProofs
	return
}

// generate the merkle Proof for an evidence
func (e *Evidence) GenerateMerkleProof(index int) (proofs MerkleProofs, cousinIndex int) {
	return GenerateProofs(e.Proofs, index)
}

// every `evidence` the node holds in memory
type EvidenceMap struct {
	M map[string]Evidence `json:"evidence_map"` // map[evidenceKey] -> Evidence
	l sync.Mutex          // a lock in the case of concurrent calls
}

// get all evidence the node holds
func GetEvidenceMap() *EvidenceMap {
	// only do once
	evidenceMapOnce.Do(func() {
		// if the all proofs object is nil
		if globalEvidenceMap == nil {
			// initialize
			globalEvidenceMap = &EvidenceMap{M: make(map[string]Evidence)}
		}
	})
	return globalEvidenceMap
}

func (e EvidenceMap) GetEvidence(h SessionHeader, evidenceType EvidenceType) (evidence Evidence, found bool) {
	key := e.KeyForEvidence(h, evidenceType)
	// lock the shared data
	e.l.Lock()
	defer e.l.Unlock()
	evidence, found = e.M[key]
	return
}

// delete evidence
func (e EvidenceMap) DeleteEvidence(h SessionHeader, evidenceType EvidenceType) {
	// lock the shared data
	e.l.Lock()
	defer e.l.Unlock()
	key := e.KeyForEvidence(h, evidenceType)
	// delete the value corresponding to the h
	delete(e.M, key)
}

// add the Proof to the EvidenceMap object
func (e EvidenceMap) AddToEvidence(h SessionHeader, p Proof) sdk.Error {
	var evidence Evidence
	key := e.KeyForEvidenceByProof(h, p)
	// lock the shared data
	e.l.Lock()
	defer e.l.Unlock()
	if _, found := e.M[key]; found {
		// if Proof already stored in allProofs
		evidence = e.M[key]
	} else {
		// if Proof is not already stored, initialize all
		evidence.SessionHeader = h
		evidence.Proofs = make([]Proof, 0)
		evidence.NumOfProofs = 0
	}
	// add Proof to the proofs object
	evidence.Proofs = append(evidence.Proofs, p)
	// increment total relay count
	evidence.NumOfProofs = evidence.NumOfProofs + 1
	// update POR
	e.M[key] = evidence
	return nil
}

func (e EvidenceMap) IsUniqueProof(h SessionHeader, p Proof) bool {
	key := e.KeyForEvidenceByProof(h, p)
	// lock the shared data
	e.l.Lock()
	defer e.l.Unlock()
	if _, found := e.M[key]; found {
		// if Proof already stored in allProofs
		evidence := e.M[key]
		// iterate over evidence to see if unique // todo efficiency (store hashes in map)
		for _, proof := range evidence.Proofs {
			if proof.HashString() == p.HashString() {
				return false
			}
		}
	}
	return true
}

func (e EvidenceMap) GetTotalRelays(h SessionHeader) int64 {
	// lock the shared data
	e.l.Lock()
	defer e.l.Unlock()
	// return the proofs object, corresponding to the h
	return e.M[e.KeyForEvidence(h, RelayEvidence)].NumOfProofs
}

// retrieve the single Proof from the all proofs object
func (e EvidenceMap) GetProof(h SessionHeader, evidenceType EvidenceType, index int) Proof {
	// lock the shared data
	e.l.Lock()
	defer e.l.Unlock()
	key := e.KeyForEvidence(h, evidenceType)
	// return the proofs object, corresponding to the h
	evidence := e.M[key].Proofs
	// do a nil check before indexing
	if evidence == nil {
		return nil
	}
	// return the Proof at specific index
	return evidence[index]
}

// retrieve the proofs from the all proofs object
func (e EvidenceMap) GetProofs(h SessionHeader, evidenceType EvidenceType) []Proof {
	// lock the shared data
	e.l.Lock()
	defer e.l.Unlock()
	key := e.KeyForEvidence(h, evidenceType)
	// return the proofs object, corresponding to the h
	return e.M[key].Proofs
}

func (e EvidenceMap) Clear() {
	globalEvidenceMap = &EvidenceMap{M: make(map[string]Evidence)}
}

// type to distinguish the types of evidence

type EvidenceType int

const (
	RelayEvidence EvidenceType = iota + 1
	ChallengeEvidence
)

func (et EvidenceType) Byte() byte {
	switch et {
	case RelayEvidence:
		return 0
	case ChallengeEvidence:
		return 1
	default:
		panic("unrecognized evidence type")
	}
}

func (e EvidenceMap) KeyForEvidence(h SessionHeader, evidenceType EvidenceType) string {
	return hex.EncodeToString(append(h.Hash(), evidenceType.Byte()))
}
func (e EvidenceMap) KeyForEvidenceByProof(h SessionHeader, p Proof) string {
	var evidenceType EvidenceType
	switch p.(type) {
	case RelayProof:
		evidenceType = RelayEvidence
	case ChallengeProofInvalidData:
		evidenceType = ChallengeEvidence
	}
	// generate the key for this specific Proof
	return e.KeyForEvidence(h, evidenceType)
}

func EvidenceTypeFromProof(p Proof) EvidenceType {
	switch p.(type) {
	case RelayProof:
		return RelayEvidence
	case ChallengeProofInvalidData:
		return ChallengeEvidence
	}
	panic("unsupported evidence type")
}

// structure used to store the proof of work
type Receipt struct {
	SessionHeader   `json:"header"`
	ServicerAddress string `json:"address"`
	TotalRelays     int64  `json:"relays"`
}
