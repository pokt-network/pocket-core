package session

import (
	"bytes"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
	"sort"
)

// The computational distance between two byte arrays
// this is judged by XORing the two
type ComputationalDistance []byte

// A simple slice abstraction of type `Node`
// These nodes are linked to the session
type SessionNodes []types.Node

// A node linked to it's computational distance
type NodeDistance struct {
	types.Node
	distance ComputationalDistance
}

type NodeDistances []NodeDistance

func (sn SessionNodes) Validate() error {
	if len(sn) < SESSIONNODECOUNT {
		return InsufficientNodesError
	}
	return nil
}

func (sn SessionNodes) Contains(n types.Node) bool { // todo use a map instead of a slice to save time
	err := sn.Validate()
	if err != nil {
		return false
	}
	// todo o(n) is too slow, see above
	for _, node := range sn {
		npkb, err := node.PubKey.Bytes()
		if err != nil {
			return false
		}
		npkb2, err := node.PubKey.Bytes()
		if bytes.Equal(npkb, npkb2) {
			return true
		}
	}
	return false
}

func NewSessionNodes(nonNativeChain SessionBlockchain, sessionKey SessionKey, allActiveNodes types.Nodes) (SessionNodes, error) {
	// validate params
	if err := nonNativeChain.Validate(); err != nil {
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
	if len(allActiveNodes) < SESSIONNODECOUNT {
		return nil, InsufficientNodesError
	}
	// filter `allActiveNodes` by the HASH(nonNativeChain)
	result, err := filter(allActiveNodes, nonNativeChain)
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
	return result[0:SESSIONNODECOUNT], nil
}

// filter the nodes by non native chain
func filter(allActiveNodes types.Nodes, nonNativeChainHash SessionBlockchain) (SessionNodes, error) {
	var result SessionNodes
	for _, node := range allActiveNodes {
		if !node.SupportedChains.Contains(string(nonNativeChainHash)) {
			continue
		}
		result = append(result, node)
	}
	if err := result.Validate(); err != nil {
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
		pubKey, err := node.PubKey.Bytes()
		pubKeyHash := crypto.SHA3FromBytes(pubKey) // currently hashing public key but could easily just take the first n bytes to compare
		if err != nil {
			return nil, err
		}
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
