package legacy

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/google/flatbuffers/go"
	"github.com/pokt-network/pocket-core/logs"
	"io/ioutil"
	"log"
	"os"
	"sort"
	"strings"
)

type NodeWorldState struct {
	Enode  string       `json:"enode"`
	Stake  int          `json:"stake"`
	Active bool         `json:"status"`
	Karma  int8         `json:"karma"`
	Chains []Blockchain `json:"chains"`
}

const (
	NODECOUNT = 5
)

func (nws NodeWorldState) EnodeSplit() (gid string, ip string, port string, discport string) {
	var url []string
	e := nws.Enode
	enodeSplit := strings.Split(e, "@")
	if strings.Contains(enodeSplit[1], "?") {
		contact := strings.Split(enodeSplit[1], "?")
		url = strings.Split(contact[0], ":")
		discport = strings.Split(contact[1], "=")[1]
	}
	url = strings.Split(enodeSplit[1], ":")
	ip = url[0]
	port = url[1]
	hash := strings.TrimPrefix(enodeSplit[0], "enode://")
	gid = hash
	return
}

type Node struct {
	GID    string `json:"gid"`
	IP     string `json:"ip"`
	Port   string `json:"port"`
	Karma  int8   `json:"karma"`
	Chains Set    `json:"chains"`
	XOR    []byte `json:"xor"`
}

type NodePool []Node

// "GetNodes" filters by blockchash, and returns the closest nodes to the key
func (n NodePool) GetNodes(s Session) (SessionNodes, error) {
	n.Filter(hex.EncodeToString(s.Chain))
	n.XOR(s)
	n.Sort()
	return n.GetClosestNodes(s)
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

// "xor" performs bitwise operation on two byte arrays and returns the destination
func xor(dst, a []byte, b []byte) error {
	if len(a) != len(b) {
		return MismatchedByteArraysError
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

// "Filter" removes any nodes from the slice that do not contain the blockchainHash
// TODO:
// This is slow O(n). If the world state is stored as a slice or as a map, this could affect
// implementation. For now this solution is acceptable, but optimizations should be made.
func (n *NodePool) Filter(blockchainHash string) {
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
	if len(*n) < 5 {
		return nil, InsufficientNodesError
	}
	sn := SessionNodes((*n)[0:5])
	return sn, nil
}

// "FileToNodes" converts the world state noodPool.json file into a slice of session.Node
func FileToNodes(nodePoolFilePath string) ([]Node, error) {
	nws := FileToNWSSlice(nodePoolFilePath)
	return nwsToNodes(nws)
}

// "FileToNWSSlice" converts a file to a slice of NodeWorldState SessionNodes
func FileToNWSSlice(nodePoolFilePath string) []NodeWorldState {
	jsonFile, err := os.Open(nodePoolFilePath)
	if err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
		fmt.Println(err.Error())
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
		fmt.Println(err.Error())
	}
	var nodes []NodeWorldState
	err = json.Unmarshal(byteValue, &nodes)
	if err != nil {
		log.Fatal(err.Error())
	}
	return nodes
}

// "nwsToNodes" converts the NodeWorldState slice of nodes from the json file
// into a []Node which is used in our session seed
func nwsToNodes(nws []NodeWorldState) ([]Node, error) {
	var nodeList []Node
	for _, node := range nws {
		if !node.Active {
			continue
		}
		n, err := nwsToNode(node)
		if err != nil {
			return nodeList, err
		}
		nodeList = append(nodeList, n)
	}
	return nodeList, nil
}

// "nwsToNode" is a helper function to NWSToNode which takes a NodeWorldState Node
// and converts it to a session.Node
func nwsToNode(nws NodeWorldState) (Node, error) {
	chains := NewSet()
	gid, ip, port, _ := nws.EnodeSplit()
	for _, c := range nws.Chains {
		marshalChain, err := MarshalBlockchain(flatbuffers.NewBuilder(0), c)
		if err != nil {
			return Node{}, err
		}
		chains.Add(hex.EncodeToString(SHA3FromBytes(marshalChain)))
	}
	return Node{GID: gid, IP: ip, Port: port, Chains: *chains, Karma: nws.Karma}, nil
}
