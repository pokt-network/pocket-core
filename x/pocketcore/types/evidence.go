package types

import (
	sdk "github.com/pokt-network/posmint/types"
	"sync"
)

var (
	globalEvidenceMap *EvidenceMap // holds all evidence from relays and challenges
	evidenceMapOnce   sync.Once    // ensure only made once
)

// Proof of relay per application
type Evidence struct {
	SessionHeader `json:"evidence_header"`      // the session h serves as an identifier for the evidence
	TotalRelays   int64   `json:"total_relays"` // the total number of relays completed
	Proofs        []Proof `json:"proofs"`       // a slice of Proof objects (Proof per relay)
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
	l sync.Mutex                                // a lock in the case of concurrent calls
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

func (e EvidenceMap) GetEvidence(h SessionHeader) (evidence Evidence, found bool) {
	key := h.HashString()
	// lock the shared data
	e.l.Lock()
	defer e.l.Unlock()
	evidence, found = e.M[key]
	return
}

// delete evidence
func (e EvidenceMap) DeleteEvidence(h SessionHeader) {
	// lock the shared data
	e.l.Lock()
	defer e.l.Unlock()
	// delete the value corresponding to the h
	delete(e.M, h.HashString())
}

// add the Proof to the EvidenceMap object
func (e EvidenceMap) AddToEvidence(h SessionHeader, p Proof) sdk.Error {
	var evidence Evidence
	// generate the key for this specific Proof
	key := h.HashString()
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
		evidence.TotalRelays = 0
	}
	// add Proof to the proofs object
	evidence.Proofs = append(evidence.Proofs, p)
	// increment total relay count
	evidence.TotalRelays = evidence.TotalRelays + 1
	// update POR
	e.M[key] = evidence
	return nil
}

func (e EvidenceMap) IsUniqueProof(h SessionHeader, p Proof) bool {
	key := h.HashString()
	// lock the shared data
	e.l.Lock()
	defer e.l.Unlock()
	if _, found := e.M[key]; found {
		// if Proof already stored in allProofs
		evidence := e.M[key]
		// iterate over evidence to see if unique // todo efficiency (store hashes in map)
		for _, proof := range evidence.Proofs {
			if proof.HashStringWithSignature() == p.HashStringWithSignature() {
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
	return e.M[h.HashString()].TotalRelays
}

// retrieve the single Proof from the all proofs object
func (e EvidenceMap) GetProof(h SessionHeader, index int) Proof {
	// lock the shared data
	e.l.Lock()
	defer e.l.Unlock()
	// return the proofs object, corresponding to the h
	evidence := e.M[h.HashString()].Proofs
	// do a nil check before indexing
	if evidence == nil {
		return nil
	}
	// return the Proof at specific index
	return evidence[index]
}

// retrieve the proofs from the all proofs object
func (e EvidenceMap) GetProofs(h SessionHeader) []Proof {
	// lock the shared data
	e.l.Lock()
	defer e.l.Unlock()
	// return the proofs object, corresponding to the h
	return e.M[h.HashString()].Proofs
}

// structure used to store the proof of work
type Receipt struct {
	SessionHeader   `json:"header"`
	ServicerAddress string `json:"address"`
	TotalRelays     int64  `json:"relays"`
}
