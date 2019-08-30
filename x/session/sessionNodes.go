package session

import (
	"encoding/hex"
	"github.com/pokt-network/pocket-core/crypto"
	types "github.com/pokt-network/pocket-core/types"
	"sort"
)

// A simple slice abstraction of type `Node`
// All of the active nodes on the Pocket Network
type AllActiveNodes []types.Node

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

func NewSessionNodes(allActiveNodes AllActiveNodes, nonNativeChain types.AminoBuffer, sessionKey SessionKey) (SessionNodes, error) {
	// session nodes are just a wrapper around node slice
	var result SessionNodes
	// node distance is just a node with a computational field attached to it
	var xorResult []NodeDistance
	// ensure there is atleast the minimum amount of nodes
	if len(allActiveNodes) < SESSIONNODECOUNT {
		return nil, InsufficientNodesError
	}
	// filter `allActiveNodes` by the HASH(nonNativeChain)
	result, err := filter(allActiveNodes, crypto.Hash(nonNativeChain))
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
	result = sorty(xorResult)

	// return the top 5 nodes
	return result[0:5], nil
}

func filter(allActiveNodes AllActiveNodes, nonNativeChainHash []byte) (SessionNodes, error) {
	var result SessionNodes
	for _, node := range allActiveNodes {
		if !node.SupportedChains.Contains(hex.EncodeToString(nonNativeChainHash)) {
			continue
		}
		result = append(result, node)
	}
	if len(result)==0 {
		return nil, InsufficientNodesError
	}
	return result, nil
}

func xor(sessionNodes SessionNodes, sessionkey SessionKey) (NodeDistances, error){
	var keyLength = len(sessionkey)
	result := make([]NodeDistance, keyLength)
	// for every node, find the distance between it's pubkey and the sesskey
	for index, node := range sessionNodes{
		pubKey := node.PubKey.Bytes()
		if len(pubKey) != keyLength {
			return nil, MismatchedByteArraysError
		}
		result[index].Node = node
		for i := 0; i < keyLength; i++ {
			result[index].distance[i] = pubKey[i] ^ sessionkey[i]
		}
	}
	return result, nil
}

func sorty(sessionNodes NodeDistances) SessionNodes {
	sort.Sort(sort.Reverse(sessionNodes))
}

// "Len" returns the length of the node pool -> needed for sort.Sort() interface
func (n NodeDistances) Len() int { return len(n) }

// "Swap" swaps two elements in the node pool -> needed for sort.Sort() interface
func (n NodeDistances) Swap(i, j int) { n[i], n[j] = n[j], n[i] }

// "Less" returns if node i is less than node j by XOR value
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
