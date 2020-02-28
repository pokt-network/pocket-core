package types

import (
	sdk "github.com/pokt-network/posmint/types"
	"sync"
)

var (
	globalAllEvidences *Evidences // holds every Proof of the node
	allEvidencesOnce   sync.Once // ensure only made once
)

// Proof of relay per application
type Evidence struct {
	SessionHeader `json:"evidence_header"`       // the session evidenceHeader serves as an identifier for the evidence
	TotalRelays   int64   `json:"total_relays"` // the total number of relays completed
	Proofs        []Proof `json:"proofs"`       // a slice of Proof objects (Proof per relay)
}

// generate the merkle root of an evidence
func (i *Evidence) GenerateMerkleRoot() (root HashSum) {
	root, sortedProofs := GenerateRoot(i.Proofs)
	i.Proofs = sortedProofs
	return
}

// generate the merkle Proof for an evidence
func (i *Evidence) GenerateMerkleProof(index int) (proofs MerkleProofs, cousinIndex int) {
	return GenerateProofs(i.Proofs, index)
}

// every `evidence` the node holds in memory
type Evidences struct {
	M map[string]Evidence `json:"evidences"` // map[evidenceKey] -> Evidence
	l sync.Mutex         // a lock in the case of concurrent calls
}

// get all evidences the node holds
func GetAllEvidences() *Evidences {
	// only do once
	allEvidencesOnce.Do(func() {
		// if the all proofs object is nil
		if globalAllEvidences == nil {
			// initialize
			globalAllEvidences = &Evidences{M: make(map[string]Evidence)}
		}
	})
	return globalAllEvidences
}

func (i Evidences) GetEvidence(evidenceHeader SessionHeader) (evidence Evidence, found bool) {
	key := evidenceHeader.HashString()
	// lock the shared data
	i.l.Lock()
	defer i.l.Unlock()
	evidence, found = i.M[key]
	return
}

func (i Evidences) IsUniqueProof(evidenceHeader SessionHeader, p Proof) bool {
	key := evidenceHeader.HashString()
	// lock the shared data
	i.l.Lock()
	defer i.l.Unlock()
	if _, found := i.M[key]; found {
		// if Proof already stored in allProofs
		evidence := i.M[key]
		// iterate over evidences to see if unique // todo efficiency (store hashes in map)
		for _, proof := range evidence.Proofs {
			if proof.HashStringWithSignature() == p.HashStringWithSignature() {
				return false
			}
		}
	}
	return true
}

// add the Proof to the Evidences object
func (i Evidences) AddToEvidence(evidenceHeader SessionHeader, p Proof) sdk.Error {
	var evidence Evidence
	// generate the key for this specific Proof
	key := evidenceHeader.HashString()
	// lock the shared data
	i.l.Lock()
	defer i.l.Unlock()
	if _, found := i.M[key]; found {
		// if Proof already stored in allProofs
		evidence = i.M[key]
	} else {
		// if Proof is not already stored, initialize all
		evidence.SessionHeader = evidenceHeader
		evidence.Proofs = make([]Proof, 0)
		evidence.TotalRelays = 0
	}
	// add Proof to the proofs object
	evidence.Proofs = append(evidence.Proofs, p)
	// increment total relay count
	evidence.TotalRelays = evidence.TotalRelays + 1
	// update POR
	i.M[key] = evidence
	return nil
}

func (i Evidences) GetTotalRelays(evidenceHeader SessionHeader) int64 {
	// lock the shared data
	i.l.Lock()
	defer i.l.Unlock()
	// return the proofs object, corresponding to the evidenceHeader
	return i.M[evidenceHeader.HashString()].TotalRelays
}

// retrieve the single Proof from the all proofs object
func (i Evidences) GetProof(evidenceHeader SessionHeader, index int) Proof {
	// lock the shared data
	i.l.Lock()
	defer i.l.Unlock()
	// return the proofs object, corresponding to the evidenceHeader
	evidence := i.M[evidenceHeader.HashString()].Proofs
	// do a nil check before indexing
	if evidence == nil {
		return Proof{}
	}
	// return the Proof at specific index
	return evidence[index]
}

// retrieve the proofs from the all proofs object
func (i Evidences) GetProofs(evidenceHeader SessionHeader) []Proof {
	// lock the shared data
	i.l.Lock()
	defer i.l.Unlock()
	// return the proofs object, corresponding to the evidenceHeader
	return i.M[evidenceHeader.HashString()].Proofs
}

// delete evidence
func (i Evidences) DeleteEvidence(evidenceHeader SessionHeader) {
	// lock the shared data
	i.l.Lock()
	defer i.l.Unlock()
	// delete the value corresponding to the evidenceHeader
	delete(i.M, evidenceHeader.HashString())
}

// structure used to store the Proof after verification
type StoredEvidence struct {
	SessionHeader   `json:"header"`
	ServicerAddress string `json:"address"`
	TotalRelays     int64  `json:"relays"`
}
