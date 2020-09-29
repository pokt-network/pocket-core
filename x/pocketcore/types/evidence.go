package types

import (
	"fmt"
	"github.com/pokt-network/pocket-core/types"
	"github.com/willf/bloom"
	"strings"
)

// "Evidence" - A proof of work/burn for nodes.
type Evidence struct {
	Bloom         bloom.BloomFilter `json:"bloom_filter"` // used to check if proof contains
	SessionHeader `json:"evidence_header"`                // the session h serves as an identifier for the evidence
	NumOfProofs   int64        `json:"num_of_proofs"`     // the total number of proofs in the evidence
	Proofs        Proofs       `json:"proofs"`            // a slice of Proof objects (Proof per relay or challenge)
	EvidenceType  EvidenceType `json:"evidence_type"`
}

func (e Evidence) IsSealed() bool {
	globalEvidenceCache.l.Lock()
	defer globalEvidenceCache.l.Unlock()
	_, ok := globalEvidenceSealedMap[e.HashString()]
	return ok
}

func (e Evidence) Seal() CacheObject {
	globalEvidenceSealedMap[e.HashString()] = struct{}{}
	return e
}

// "GenerateMerkleRoot" - Generates the merkle root for an evidence object
func (e *Evidence) GenerateMerkleRoot() (root HashRange) {
	// seal the evidence in cache/db
	ev, ok := SealEvidence(*e)
	if !ok {
		return HashRange{}
	}
	// generate the root object
	root, _ = GenerateRoot(ev.Proofs)
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
func (e *Evidence) GenerateMerkleProof(index int) (proof MerkleProof, leaf Proof) {
	// generate the merkle proof
	proof, leaf = GenerateProofs(e.Proofs, index)
	// set the evidence in memory
	return
}

// "Evidence" - A proof of work/burn for nodes.
type evidence struct {
	BloomBytes    []byte `json:"bloom_bytes"`
	SessionHeader `json:"evidence_header"`            // the session h serves as an identifier for the evidence
	NumOfProofs   int64        `json:"num_of_proofs"` // the total number of proofs in the evidence
	Proofs        []Proof      `json:"proofs"`        // a slice of Proof objects (Proof per relay or challenge)
	EvidenceType  EvidenceType `json:"evidence_type"`
}

var _ CacheObject = Evidence{} // satisfies the cache object interface

func (e Evidence) MarshalObject() ([]byte, error) {
	encodedBloom, err := e.Bloom.GobEncode()
	if err != nil {
		return nil, err
	}
	if ModuleCdc.IsAfterUpgrade() { // TODO phase out if else
		ep := ProtoEvidence{
			BloomBytes:    encodedBloom,
			SessionHeader: &e.SessionHeader,
			NumOfProofs:   e.NumOfProofs,
			Proofs:        e.Proofs.ToProofI(),
			EvidenceType:  e.EvidenceType,
		}
		return ModuleCdc.MarshalBinaryBare(&ep)
	} else {
		ep := evidence{
			BloomBytes:    encodedBloom,
			SessionHeader: e.SessionHeader,
			NumOfProofs:   e.NumOfProofs,
			Proofs:        e.Proofs,
			EvidenceType:  e.EvidenceType,
		}
		return ModuleCdc.MarshalBinaryBare(ep)
	}

}

func (e Evidence) UnmarshalObject(b []byte) (CacheObject, error) {
	if ModuleCdc.IsAfterUpgrade() {  // TODO phase out if else
		ep := ProtoEvidence{}
		err := ModuleCdc.UnmarshalBinaryBare(b, &ep)
		if err != nil {
			return Evidence{}, fmt.Errorf("could not unmarshal into ProtoEvidence from cache, moduleCdc unmarshal binary bare: %s", err.Error())
		}
		bloomFilter := bloom.BloomFilter{}
		err = bloomFilter.GobDecode(ep.BloomBytes)
		if err != nil {
			return Evidence{}, fmt.Errorf("could not unmarshal into ProtoEvidence from cache, bloom bytes gob decode: %s", err.Error())
		}
		evidence := Evidence{
			Bloom:         bloomFilter,
			SessionHeader: *ep.SessionHeader,
			NumOfProofs:   ep.NumOfProofs,
			Proofs:        ep.Proofs.FromProofI(),
			EvidenceType:  ep.EvidenceType}
		return evidence, nil
	} else {
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
		evidence := Evidence{
			Bloom:         bloomFilter,
			SessionHeader: ep.SessionHeader,
			NumOfProofs:   ep.NumOfProofs,
			Proofs:        ep.Proofs,
			EvidenceType:  ep.EvidenceType}
		return evidence, nil
	}
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
