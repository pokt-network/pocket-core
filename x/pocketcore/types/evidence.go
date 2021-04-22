package types

import (
	"fmt"
	"strings"

	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/types"
	"github.com/willf/bloom"
)

// "Evidence" - A proof of work/burn for nodes.
type Evidence struct {
	Bloom         bloom.BloomFilter        `json:"bloom_filter"` // used to check if proof contains
	SessionHeader `json:"evidence_header"` // the session h serves as an identifier for the evidence
	NumOfProofs   int64                    `json:"num_of_proofs"` // the total number of proofs in the evidence
	Proofs        Proofs                   `json:"proofs"`        // a slice of Proof objects (Proof per relay or challenge)
	EvidenceType  EvidenceType             `json:"evidence_type"`
}

func (e Evidence) IsSealed() bool {
	globalEvidenceCache.l.Lock()
	defer globalEvidenceCache.l.Unlock()
	_, ok := globalEvidenceSealedMap.Load(e.HashString())
	return ok
}

func (e Evidence) Seal() CacheObject {
	globalEvidenceSealedMap.Store(e.HashString(), struct{}{})
	return e
}

// "GenerateMerkleRoot" - Generates the merkle root for an GOBEvidence object
func (e *Evidence) GenerateMerkleRoot(height int64) (root HashRange) {
	// seal the evidence in cache/db
	ev, ok := SealEvidence(*e)
	if !ok {
		return HashRange{}
	}
	// generate the root object
	root, _ = GenerateRoot(height, ev.Proofs)
	return
}

// "AddProof" - Adds a proof obj to the GOBEvidence field
func (e *Evidence) AddProof(p Proof) {
	// add proof to GOBEvidence
	e.Proofs = append(e.Proofs, p)
	// increment total proof count
	e.NumOfProofs = e.NumOfProofs + 1
	// add proof to bloom filter
	e.Bloom.Add(p.Hash())
}

// "GenerateMerkleProof" - Generates the merkle Proof for an GOBEvidence
func (e *Evidence) GenerateMerkleProof(height int64, index int) (proof MerkleProof, leaf Proof) {
	// generate the merkle proof
	proof, leaf = GenerateProofs(height, e.Proofs, index)
	// set the evidence in memory
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

func (e Evidence) LegacyAminoMarshal() ([]byte, error) {
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
	return ModuleCdc.MarshalBinaryBare(ep, 0)
}

func (e Evidence) LegacyAminoUnmarshal(b []byte) (CacheObject, error) {
	ep := evidence{}
	err := ModuleCdc.UnmarshalBinaryBare(b, &ep, 0)
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
		EvidenceType:  ep.EvidenceType,
	}
	return evidence, nil
}

var (
	_ CacheObject          = Evidence{} // satisfies the cache object interface
	_ codec.ProtoMarshaler = &Evidence{}
)

func (e *Evidence) Reset() {
	*e = Evidence{}
}

func (e *Evidence) String() string {
	return fmt.Sprintf("SessionHeader: %v\nNumOfProofs: %v\nProofs: %v\nEvidenceType: %vBloomFilter: %v\n",
		e.SessionHeader, e.NumOfProofs, e.Proofs, e.EvidenceType, e.Bloom)
}

func (e *Evidence) ProtoMessage() {}

func (e *Evidence) Marshal() ([]byte, error) {
	pe, err := e.ToProto()
	if err != nil {
		return nil, err
	}
	return pe.Marshal()
}

func (e *Evidence) MarshalTo(data []byte) (n int, err error) {
	pe, err := e.ToProto()
	if err != nil {
		return 0, err
	}
	return pe.MarshalTo(data)
}

func (e *Evidence) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	pe, err := e.ToProto()
	if err != nil {
		return 0, err
	}
	return pe.MarshalToSizedBuffer(dAtA)
}

func (e *Evidence) Size() int {
	pe, err := e.ToProto()
	if err != nil {
		return 0
	}
	return pe.Size()
}

func (e *Evidence) Unmarshal(data []byte) error {
	pe := ProtoEvidence{}
	err := pe.Unmarshal(data)
	if err != nil {
		return err
	}
	*e, err = pe.FromProto()
	return err
}

func (e *Evidence) ToProto() (*ProtoEvidence, error) {
	encodedBloom, err := e.Bloom.GobEncode()
	if err != nil {
		return nil, err
	}
	return &ProtoEvidence{
		BloomBytes:    encodedBloom,
		SessionHeader: &e.SessionHeader,
		NumOfProofs:   e.NumOfProofs,
		Proofs:        e.Proofs.ToProofI(),
		EvidenceType:  e.EvidenceType,
	}, nil
}

func (pe *ProtoEvidence) FromProto() (Evidence, error) {
	bloomFilter := bloom.BloomFilter{}
	err := bloomFilter.GobDecode(pe.BloomBytes)
	if err != nil {
		return Evidence{}, fmt.Errorf("could not unmarshal into ProtoEvidence from cache, bloom bytes gob decode: %s", err.Error())
	}
	return Evidence{
		Bloom:         bloomFilter,
		SessionHeader: *pe.SessionHeader,
		NumOfProofs:   pe.NumOfProofs,
		Proofs:        pe.Proofs.FromProofI(),
		EvidenceType:  pe.EvidenceType}, nil
}

func (e Evidence) MarshalObject() ([]byte, error) {
	pe, err := e.ToProto()
	if err != nil {
		return nil, err
	}
	return ModuleCdc.ProtoMarshalBinaryBare(pe)
}

func (e Evidence) UnmarshalObject(b []byte) (CacheObject, error) {
	pe := ProtoEvidence{}
	err := ModuleCdc.ProtoUnmarshalBinaryBare(b, &pe)
	if err != nil {
		return Evidence{}, fmt.Errorf("could not unmarshal into ProtoEvidence from cache, moduleCdc unmarshal binary bare: %s", err.Error())
	}
	return pe.FromProto()
}

func (e Evidence) Key() ([]byte, error) {
	return KeyForEvidence(e.SessionHeader, e.EvidenceType)
}

// "EvidenceType" type to distinguish the types of GOBEvidence (relay/challenge)
type EvidenceType int

const (
	RelayEvidence EvidenceType = iota + 1 // essentially an enum for GOBEvidence types
	ChallengeEvidence
)

// "Convert GOBEvidence type to bytes
func (et EvidenceType) Byte() (byte, error) {
	switch et {
	case RelayEvidence:
		return 0, nil
	case ChallengeEvidence:
		return 1, nil
	default:
		return 0, fmt.Errorf("unrecognized GOBEvidence type")
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
