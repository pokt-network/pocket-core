package types

import (
	"encoding/hex"
	"encoding/json"
	merkle "github.com/pokt-network/merkle"
	sdk "github.com/pokt-network/posmint/types"
	"sync"
)

// proof of relay per application
type ProofOfRelay struct {
	SessionHeader
	TotalRelays int64
	Proofs      []Proof // slice[index] -> Proofs
	Tree        Tree
}

// generates the corresponding merkle tree for the proof of relay
func (por *ProofOfRelay) GenerateMerkleTree() {
	var data [][]byte
	// create a merkle tree using SHA3_256 Hashing algorithm
	tree := merkle.NewTree(Hasher.New())
	// for each proof in proofs
	for i, proof := range por.Proofs {
		por.Proofs[i].Index = int64(i)
		data = append(data, proof.Hash())
	}
	// add the data to the tree
	tree.AddData(data...)
	// generate the tree
	err := tree.Generate()
	if err != nil {
		panic(err)
	}
	// set the tree in the por structure
	por.Tree = Tree(tree)
}

// extended the merkle.Tree structure to return sdk.Errors
type Tree merkle.Tree

// get the root of the merkle tree
func (t Tree) GetMerkleRoot() ([]byte, sdk.Error) {
	// check for empty tree
	if len(t.Nodes) == 0 {
		return nil, NewEmptyMerkleTreeError(ModuleName)
	}
	// get the merkle root of the tree
	root := merkle.Tree(t).Root()
	// if the root is empty
	if root == nil || len(root) == 0 {
		return nil, NewNodeNotFoundErr(ModuleName)
	}
	return root, nil
}

// get the proof needed to prove index
func (t Tree) GetMerkleProof(index int) (MerkleProof, sdk.Error) {
	// check for empty tree
	if len(t.Nodes) == 0 {
		return nil, NewEmptyMerkleTreeError(ModuleName)
	}
	// get the proof from the tree
	proof := merkle.Tree(t).GetProof(index)
	if proof == nil || len(proof) == 0 {
		return nil, NewNodeNotFoundErr(ModuleName)
	}
	// return the proof object
	return MerkleProof(proof), nil
}

// every proof the node holds
type AllProofs struct {
	M map[string]ProofOfRelay // map[porkey] -> ProofOfRelay
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
			*globalAllProofs = AllProofs{M: make(map[string]ProofOfRelay)}
		}
	})
	return globalAllProofs
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
	Index              int64
	SessionBlockHeight int64
	ServicerPubKey     string
	Blockchain         string
	Token              AAT
	Signature          string
}

func (p Proof) Validate(maxRelays int64, hb HostedBlockchains, verifyPubKey string) sdk.Error {
	// validate the session block height
	if p.SessionBlockHeight < 0 {
		return NewInvalidBlockHeightError(ModuleName)
	}
	// validate blockchain
	if err := HashVerification(p.Blockchain); err != nil {
		return err
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
	Index              int64
	SessionBlockHeight int64
	ServicerPubKey     string
	Blockchain         string
	Signature          string
	Token              string
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
