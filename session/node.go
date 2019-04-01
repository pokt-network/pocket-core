package session

import (
	"encoding/hex"
	"errors"
	"github.com/pokt-network/pocket-core/types"
	"log"
	"sort"
)

type role int

const (
	NODECOUNT          = 5
	MAXVALIDATORS      = (NODECOUNT / 2) + 1
	MAXSERVICERS       = NODECOUNT / 2
	VALIDATE      role = iota + 1
	SERVICE
	DELEGATEDMINTER
)

type Node struct {
	GID    string    `json:"gid"`
	IP     string    `json:"ip"`
	Port   string    `json:"port"`
	Role   role      `json:"role"`
	Chains types.Set `json:"chains"`
	XOR    []byte    `json:"xor"`
}

type NodePool []Node

// "GetSessionNodes" filters by blockchash, and returns the closest nodes to the key
func (n NodePool) GetSessionNodes(s Session) (SessionNodes, error) {
	n.Filter(hex.EncodeToString(s.Chain))
	n.XOR(s)
	n.Sort()
	return n.GetClosestNodes(s)
}

// "Init" creates a nodePool from a seed
func (n NodePool) Init(seed Seed) {
	// n = seed.NodeList
}

// "Len" returns the length of the node pool -> needed for sort.Sort() interface
func (n NodePool) Len() int { return len(n) }

// "Swap" swaps two elements in the node pool -> needed for sort.Sort() interface
func (n NodePool) Swap(i, j int) { n[i], n[j] = n[j], n[i] }

// "Less" returns if node i is less than node j by XOR value
// it assumes big endian encoding
func (n NodePool) Less(i, j int) bool {
	if len(n[i].XOR) < len(n[j].XOR) {
		return false
	}
	for a := range n[i].XOR {
		if n[i].XOR[a] < n[j].XOR[a] {
			return true
		}
		if n[i].XOR[a] < n[i].XOR[a] {
			return false
		}
	}
	return false
}

// "Sort" sorts the nodes based on XOR proximity of GID and Session Key
func (n NodePool) Sort() {
	sort.Sort(sort.Reverse(n))
}

// "XOR" performs xor operation on each node's gid and session key
func (n *NodePool) XOR(s Session) {
	for a := range *n {
		gid, err := hex.DecodeString((*n)[a].GID)
		if err != nil {
			log.Fatal(err.Error())
		}
		xor((*n)[a].XOR, gid, s.Key)
	}
}

// "xor" performs bitwise operation on two byte arrays and returns the destination
func xor(dst, a []byte, b []byte) error {
	if len(a) != len(b) {
		return errors.New("mismatched size error")
	}
	n := len(a)
	if len(dst) < n {
		n = len(dst)
	}
	for i := 0; i < n; i++ {
		dst[i] = a[i] ^ b[i]
	}
	return nil
}

// "Filter" removes any nodes from the slice that do not contain the blockchainHash
func (n *NodePool) Filter(blockchainHash string) {
	// TODO possible optimizations from node slice to map
	tmp := (*n)[:0]
	for _, node := range *n {
		if node.Chains.Contains(blockchainHash) {
			tmp = append(tmp, node)
		}
	}
	*n = tmp
}

// "GetClosestNodes" returns the 'proper' closest nodes to the session key
func (n *NodePool) GetClosestNodes(s Session) (SessionNodes, error) {
	var sessionNodes SessionNodes
	for _, node := range *n {
		if node.Role == VALIDATE {
			if len(sessionNodes.ValidatorNodes) == 0 {
				sessionNodes.DelegatedMinter = node
				node.Role = DELEGATEDMINTER
			}
			if len(sessionNodes.ValidatorNodes) != MAXVALIDATORS {
				sessionNodes.ValidatorNodes = append(sessionNodes.ValidatorNodes, node)
				continue
			}
		}
		if len(sessionNodes.ServiceNodes) != MAXSERVICERS {
			node.Role = SERVICE
			sessionNodes.ServiceNodes = append(sessionNodes.ServiceNodes, node)
			continue
		}
		return sessionNodes, nil
	}
	return sessionNodes, InsufficientNodes
}
