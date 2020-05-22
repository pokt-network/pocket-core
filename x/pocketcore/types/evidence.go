package types

import (
	"fmt"
	"github.com/pokt-network/posmint/types"
	"github.com/willf/bloom"
	"strings"
)

// "Evidence" - A proof of work/burn for nodes.
type Evidence struct {
	Bloom         bloom.BloomFilter        `json:"bloom_filter"` // used to check if proof contains
	SessionHeader `json:"evidence_header"` // the session h serves as an identifier for the evidence
	NumOfProofs   int64                    `json:"num_of_proofs"` // the total number of proofs in the evidence
	Proofs        []Proof                  `json:"proofs"`        // a slice of Proof objects (Proof per relay or challenge)
	EvidenceType  EvidenceType             `json:"evidence_type"`
}

// "GenerateMerkleRoot" - Generates the merkle root for an evidence object
func (e *Evidence) GenerateMerkleRoot() (root HashSum) {
	// generate the root object
	root, sortedProofs := GenerateRoot(e.Proofs)
	// sort the proofs
	e.Proofs = sortedProofs
	// set the evidence in cache
	SetEvidence(*e)
	return
}

// "AddProof" - Adds a proof obj to the evidence field
func (e *Evidence) AddProof(p Proof) {
	// add proof to evidence
	e.Proofs = append(e.Proofs, p)
	// increment total proof count
	e.NumOfProofs = e.NumOfProofs + 1
	// add proof to bloom filter
	e.Bloom.Add(p.Hash())
}

// "GenerateMerkleProof" - Generates the merkle Proof for an evidence
func (e *Evidence) GenerateMerkleProof(index int) (proofs MerkleProofs, cousinIndex int) {
	// generate the merkle proof
	proofs, cousinIndex = GenerateProofs(e.Proofs, index)
	// set the evidence in memory
	SetEvidence(*e)
	return
}

// "Evidence" - A proof of work/burn for nodes.
type evidence struct {
	BloomBytes    []byte                   `json:"bloom_bytes"`
	SessionHeader `json:"evidence_header"` // the session h serves as an identifier for the evidence
	NumOfProofs   int64                    `json:"num_of_proofs"` // the total number of proofs in the evidence
	Proofs        []Proof                  `json:"proofs"`        // a slice of Proof objects (Proof per relay or challenge)
	EvidenceType  EvidenceType             `json:"evidence_type"`
}

var _ CacheObject = Evidence{} // satisfies the cache object interface

func (e Evidence) Marshal() ([]byte, error) {
	encodedBloom, err := e.Bloom.GobEncode()
	if err != nil {
		return nil, err
	}
	ep := evidence{
		BloomBytes:    encodedBloom,
		SessionHeader: e.SessionHeader,
		NumOfProofs:   e.NumOfProofs,
		Proofs:        e.Proofs,
		EvidenceType:  e.EvidenceType,
	}
	return ModuleCdc.MarshalBinaryBare(ep)
}

func (e Evidence) Unmarshal(b []byte) (CacheObject, error) {
	ep := evidence{}
	err := ModuleCdc.UnmarshalBinaryBare(b, &ep)
	if err != nil {
		return Evidence{}, fmt.Errorf("could not unmarshal into evidence from cache, moduleCdc unmarshal binary bare: %s", err.Error())
	}
	bloomFilter := bloom.BloomFilter{}
	err = bloomFilter.GobDecode(ep.BloomBytes)
	if err != nil {
		return Evidence{}, fmt.Errorf("could not unmarshal into evidence from cache, bloom bytes gob decode: %s", err.Error())
	}
	return Evidence{
		Bloom:         bloomFilter,
		SessionHeader: e.SessionHeader,
		NumOfProofs:   e.NumOfProofs,
		Proofs:        e.Proofs,
		EvidenceType:  e.EvidenceType,
	}, nil
}

func (e Evidence) Key() ([]byte, error) {
	return KeyForEvidence(e.SessionHeader, e.EvidenceType)
}

// "EvidenceType" type to distinguish the types of evidence (relay/challenge)
type EvidenceType int

const (
	RelayEvidence EvidenceType = iota + 1 // essentially an enum for evidence types
	ChallengeEvidence
)

// "Convert evidence type to bytes
func (et EvidenceType) Byte() (byte, error) {
	switch et {
	case RelayEvidence:
		return 0, nil
	case ChallengeEvidence:
		return 1, nil
	default:
		return 0, fmt.Errorf("unrecognized evidence type")
	}
}

// "Receipt" - Is a structure used to store proof of evidence after verification
type Receipt struct {
	SessionHeader   `json:"header"` // header to identify the session
	ServicerAddress string          `json:"address"`       // the address responsible
	Total           int64           `json:"total"`         // the number of proofs
	EvidenceType    EvidenceType    `json:"evidence_type"` // the type (relay/challenge)
}

func EvidenceTypeFromString(evidenceType string) (et EvidenceType, err types.Error) {
	switch strings.ToLower(evidenceType) {
	case "relay":
		et = RelayEvidence
	case "challenge":
		et = ChallengeEvidence
	default:
		err = types.ErrInternal("type in the receipt query is not recognized: (relay or challenge)")
	}
	return
}
