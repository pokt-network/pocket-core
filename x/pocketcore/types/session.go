package types

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
	nodeexported "github.com/pokt-network/pocket-core/x/nodes/exported"
	"sort"
)

type Session struct {
	SessionKey     SessionKey                `json:"sessionkey"`
	AppPubKey      string                    `json:"appPubKey"`
	NonNativeChain string                    `json:"chain"`
	BlockHash      string                    `json:"blockHash"`
	BlockHeight    int64                     `json:"blockHeight"`
	Nodes          []nodeexported.ValidatorI `json:"nodes"`
}

// Create a new session from seed data
func NewSession(appPubKey string, nonNativeChain string, blockID string, blockHeight int64, allActiveNodes []nodeexported.ValidatorI, sessionNodesCount int) (*Session, error) { // todo possibly convert block id to block hash
	// first generate session key
	sessionKey, err := NewSessionKey(appPubKey, nonNativeChain, blockID)
	if err != nil {
		return nil, err
	}
	// then generate the service nodes for that session
	sessionNodes, err := NewSessionNodes(nonNativeChain, sessionKey, allActiveNodes, sessionNodesCount)
	if err != nil {
		return nil, err
	}
	// then populate the structure and return
	return &Session{SessionKey: sessionKey, AppPubKey: appPubKey, BlockHeight: blockHeight, NonNativeChain: nonNativeChain, BlockHash: blockID, Nodes: sessionNodes}, nil
}

// A simple slice abstraction of type `Node`
// These nodes are linked to the session
type SessionNodes []nodeexported.ValidatorI

// A node linked to it's computational distance
type NodeDistance struct {
	Node     nodeexported.ValidatorI
	distance []byte
}

type NodeDistances []NodeDistance

func (sn SessionNodes) Validate(sessionNodesCount int) error {
	if len(sn) < sessionNodesCount {
		return InsufficientNodesError
	}
	return nil
}

func (sn SessionNodes) Contains(nodeVerify nodeexported.ValidatorI, sessionNodesCount int) bool { // todo use a map instead of a slice to save time
	if nodeVerify == nil {
		return false
	}
	err := sn.Validate(sessionNodesCount)
	if err != nil {
		return false
	}
	// todo o(n) is too slow, see above
	for _, node := range sn {
		if node.GetConsPubKey().Equals(nodeVerify.GetConsPubKey()) {
			return true
		}
	}
	return false
}

func NewSessionNodes(nonNativeChain string, sessionKey SessionKey, allActiveNodes []nodeexported.ValidatorI, sessionNodesCount int) (SessionNodes, error) {
	// validate params
	if len(nonNativeChain) == 0 {
		return nil, EmptyNonNativeChainError
	}
	// validate params
	if err := sessionKey.Validate(); err != nil {
		return nil, err
	}
	// session nodes are just a wrapper around node slice
	var result SessionNodes
	// node distance is just a node with a computational field attached to it
	var xorResult []NodeDistance
	// ensure there is atleast the minimum amount of nodes
	if len(allActiveNodes) < sessionNodesCount {
		return nil, InsufficientNodesError
	}
	// filter `allActiveNodes` by the HASH(nonNativeChain)
	result, err := filter(allActiveNodes, nonNativeChain, sessionNodesCount)
	if err != nil {
		return nil, err
	}
	// xor each node's public key and session key
	// return NodeDistance array to be ordered
	xorResult, err = xor(result, sessionKey)
	if err != nil {
		return nil, err
	}
	// sort the nodes based off of distance
	result = revSort(xorResult)

	// return the top 5 nodes
	return result[0:sessionNodesCount], nil
}

// filter the nodes by non native chain
func filter(allActiveNodes []nodeexported.ValidatorI, nonNativeChainHash string, sessionNodesCount int) (SessionNodes, error) {
	var result SessionNodes
	for _, node := range allActiveNodes {
		if _, contains := node.GetChains()[nonNativeChainHash]; !contains {
			continue
		}
		result = append(result, node)
	}
	if err := result.Validate(sessionNodesCount); err != nil {
		return nil, err
	}
	return result, nil
}

// xor the sessionNodes.publicKey against the sessionKey to find the computationally
// closest nodes
func xor(sessionNodes SessionNodes, sessionkey SessionKey) (NodeDistances, error) {
	var keyLength = len(sessionkey)
	result := make([]NodeDistance, len(sessionNodes))
	// for every node, find the distance between it's pubkey and the sesskey
	for index, node := range sessionNodes {
		pubKey := node.GetConsPubKey()
		pubKeyHash := crypto.SHA3FromBytes(pubKey.Bytes()) // currently hashing public key but could easily just take the first n bytes to compare
		if len(pubKeyHash) != keyLength {
			return nil, MismatchedByteArraysError
		}
		result[index].Node = node
		result[index].distance = make([]byte, keyLength)
		for i := 0; i < keyLength; i++ {
			result[index].distance[i] = pubKeyHash[i] ^ sessionkey[i]
		}
	}
	return result, nil
}

// sort the nodes by shortest computational distance
func revSort(sessionNodes NodeDistances) SessionNodes {
	result := make(SessionNodes, len(sessionNodes))
	sort.Sort(sort.Reverse(sessionNodes))
	for _, node := range sessionNodes {
		result = append(result, node.Node)
	}
	return result
}

// returns the length of the node pool -> needed for sort.Sort() interface
func (n NodeDistances) Len() int { return len(n) }

// swaps two elements in the node pool -> needed for sort.Sort() interface
func (n NodeDistances) Swap(i, j int) { n[i], n[j] = n[j], n[i] }

// returns if node i is less than node j by XOR value
// it assumes big endian encoding
func (n NodeDistances) Less(i, j int) bool {
	// compare size of byte arrays
	if len(n[i].distance) < len(n[j].distance) {
		return false
	}
	// bitwise comparison
	for a := range n[i].distance {
		if n[i].distance[a] < n[j].distance[a] {
			return true
		}
		if n[i].distance[a] < n[i].distance[a] {
			return false
		}
	}
	return false
}

type SessionKey []byte

// Generates the session key = SessionHashingAlgo(devid+chain+blockhash)
func NewSessionKey(appPublicKey string, nonNativeChain string, blockHash string) (SessionKey, error) {
	// validate session application
	if len(appPublicKey) == 0 {
		return nil, EmptyAppPubKeyError
	}
	// get the public key from the appPublicKey structure
	appPubKey, err := hex.DecodeString(appPublicKey)
	if err != nil {
		return nil, err
	}
	if len(nonNativeChain) == 0 {
		return nil, EmptyNonNativeChainError
	}
	if len(blockHash) == 0 {
		return nil, EmptyBlockIDError
	}
	nnBytes, err := hex.DecodeString(nonNativeChain)
	if err != nil {
		return nil, err
	}
	// append them all together
	// in the order of appPubKey - > nonnativeChain -> blockHash
	// TODO consider using amino buffer to find the session key
	seed := append(appPubKey, nnBytes...)
	seed = append(seed, blockHash...)

	// return the hash of the result
	return crypto.SHA3FromBytes(seed), nil
}

func (sk SessionKey) Validate() error {
	// todo more validation
	if len(sk) == 0 {
		return EmptySessionKeyError
	}
	return nil
}
