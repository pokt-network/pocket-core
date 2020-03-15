package types

// Proof of relay per application
type Evidence struct {
	SessionHeader `json:"evidence_header"` // the session h serves as an identifier for the evidence
	NumOfProofs   int64                    `json:"num_of_proofs"` // the total number of proofs in the evidence
	Proofs        []Proof                  `json:"proofs"`        // a slice of Proof objects (Proof per relay or challenge)
}

// generate the merkle root of an evidence
func (e *Evidence) GenerateMerkleRoot() (root HashSum) {
	root, sortedProofs := GenerateRoot(e.Proofs)
	e.Proofs = sortedProofs
	SetEvidence(*e, e.Proofs[0].EvidenceType())
	return
}

func (e *Evidence) AddProof(p Proof) {
	// add proof to evidence
	e.Proofs = append(e.Proofs, p)
	// increment total proof count
	e.NumOfProofs = e.NumOfProofs + 1
}

// generate the merkle Proof for an evidence
func (e *Evidence) GenerateMerkleProof(index int) (proofs MerkleProofs, cousinIndex int) {
	proofs, cousinIndex = GenerateProofs(e.Proofs, index)
	SetEvidence(*e, e.Proofs[0].EvidenceType())
	return
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

// structure used to store the proof of work
type Receipt struct {
	SessionHeader   `json:"header"`
	ServicerAddress string       `json:"address"`
	Total           int64        `json:"total"`
	EvidenceType    EvidenceType `json:"evidence_type"`
}
