package types

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
	"strconv"
)

// Proof per relay
type Proof struct {
	Index              int
	SessionBlockHeight int64
	ServicerPubKey     string
	Token              AAT
	Signature          string
}

type PORHeader struct {
	ApplicationPubKey  string
	Chain              string
	SessionBlockHeight int64
}

// ProofOfRelay per application
type ProofOfRelay struct {
	PORHeader
	TotalRelays int64
	Proofs      []Proof // map[clientPubKey] -> Proofs
}

// structure to map out all proofs
type AllProofs map[string]ProofOfRelay // map[appPubKey+chain+blockheight] -> ProofOfRelay

func (ap AllProofs) AddProof(header PORHeader, p Proof, maxRelays int) error { // todo need mutex
	var por = ProofOfRelay{}
	porKey := KeyForPOR(header.ApplicationPubKey, header.Chain, strconv.Itoa(int(p.SessionBlockHeight)))
	// first check to see if all proofs contain a proof of relay for that specific application and session
	if _, found := ap[porKey]; found {
		por = ap[porKey]
	} else {
		// if not found fill in the header info
		por.SessionBlockHeight = header.SessionBlockHeight
		por.ApplicationPubKey = header.ApplicationPubKey
		por.Chain = header.Chain
		por.Proofs = make([]Proof, maxRelays)
		por.TotalRelays = 0
	}
	// check to see if ticket was already punched
	if pf := por.Proofs[p.Index]; pf.Signature != "" {
		return DuplicateProofError
	}
	// else add the proof to the slice
	por.Proofs[p.Index] = p
	// increment total relay count
	por.TotalRelays = por.TotalRelays + 1
	// update POR
	ap[porKey] = por
	return nil
}

func (ap AllProofs) GetProof(header PORHeader, index int) *Proof {
	por := ap[KeyForPOR(header.ApplicationPubKey, header.Chain, strconv.Itoa(int(header.SessionBlockHeight)))].Proofs
	if por == nil {
		return nil
	}
	return &por[index]
}

func (ap AllProofs) ClearProofs(header PORHeader, maxRelays int) {
	porKey := KeyForPOR(header.ApplicationPubKey, header.Chain, strconv.Itoa(int(header.SessionBlockHeight)))
	ap[porKey] = ProofOfRelay{
		PORHeader:   header,
		Proofs:      make([]Proof, maxRelays),
		TotalRelays: 0,
	}
}

func (p Proof) Validate(maxRelays int64, servicerVerifyPubKey string) error {
	// check for negative counter
	if p.Index < 0 {
		return NegativeICCounterError
	}
	if int64(p.Index) > maxRelays {
		return MaximumIncrementCounterError
	}
	// validate the service token
	if err := p.Token.Validate(); err != nil {
		return NewInvalidTokenError(ModuleName, err)
	}
	// validate the public key correctness
	if p.ServicerPubKey != servicerVerifyPubKey {
		return InvalidNodePubKeyError // the public key is not this nodes, so they would not get paid
	}
	hash, err := p.Hash()
	if err != nil {
		return NewServiceProofHashError(ModuleName, err)
	}
	if !crypto.MockVerifySignature(p.Token.ClientPublicKey, hex.EncodeToString(hash), p.Signature) { // todo real signature verification
		return InvalidICSignatureError
	}
	return nil
}

func (p Proof) Hash() ([]byte, error) {
	return crypto.SHA3FromBytes(append([]byte(p.ServicerPubKey))), nil // !!!!todo this needs to hash everything in this message for sig verification (need byte representation of token)
}

func (ph PORHeader) String() string {
	return strconv.Itoa(int(ph.SessionBlockHeight)) + ph.ApplicationPubKey + ph.Chain // todo standardize
}
