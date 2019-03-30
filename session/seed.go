package session

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/common"
	"github.com/pokt-network/pocket-core/logs"
	"github.com/pokt-network/pocket-core/types"
	"io/ioutil"
	"log"
	"os"
)

type Seed struct {
	DevID          []byte
	BlockHash      []byte
	RequestedChain []byte
	NodeList       []Node
}

// "NewSeed" is the constructor of the sessionSeed
func NewSeed(devID []byte, nodePoolFilePath string, requestedBlockchain []byte, blockHash []byte) Seed {
	return Seed{DevID: devID, BlockHash: blockHash, RequestedChain: requestedBlockchain, NodeList: FileToNodes(nodePoolFilePath)}
}

// "FileToNodes" converts the world state noodPool.json file into a slice of session.Node
func FileToNodes(nodePoolFilePath string) []Node {
	nws := FileToNWSSlice(nodePoolFilePath)
	return NWSToNodes(nws)
}

// "FileToNWSSlice" converts a file to a slice of NodeWorldState Nodes
func FileToNWSSlice(nodePoolFilePath string) []common.NodeWorldState {
	jsonFile, err := os.Open(nodePoolFilePath)
	if err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
	}
	defer jsonFile.Close()
	byteValue, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		logs.NewLog(err.Error(), logs.ErrorLevel, logs.JSONLogFormat)
	}
	var nodes []common.NodeWorldState
	err = json.Unmarshal(byteValue, &nodes)
	if err != nil {
		log.Fatal(err.Error())
	}
	return nodes
}

// "NWSToNodes" converts the NodeWorldState slice of nodes from the json file
// into a []Node which is used in our session seed
func NWSToNodes(nws []common.NodeWorldState) []Node {
	var nodeList []Node
	for _, node := range nws {
		if !node.Active {
			continue
		}
		nodeList = append(nodeList, nwsToNode(node))
	}
	return nodeList
}

// "nwsToNode" is a helper function to NWSToNode which takes a NodeWorldState Node
// and converts it to a session.Node
func nwsToNode(nws common.NodeWorldState) Node {
	chains := types.NewSet()
	var role role
	gid, ip, port, _ := nws.EnodeSplit()
	for _, c := range nws.Chains {
		chains.Add(common.SHA256FromString(fmt.Sprintf("%v", c)))
	}
	switch nws.IsVal {
	case true:
		role = VALIDATE
	case false:
		role = SERVICE
	}
	return Node{GID: gid, IP: ip, Port: port, Chains: *chains, Role: role}
}

func (s *Seed) ErrorCheck() error {
	if s.DevID == nil || len(s.DevID) == 0 {
		return NoDevID
	}
	if s.BlockHash == nil || len(s.BlockHash) == 0 {
		return NoBlockHash
	}
	if s.NodeList == nil || len(s.NodeList) == 0 {
		return NoNodeList
	}
	if s.RequestedChain == nil || len(s.RequestedChain) == 0 {
		return NoReqChain
	}
	if len(s.BlockHash) != 32 {
		return InvalidBlockHashFormat
	}
	if len(s.DevID) != 32 {
		return InvalidDevIDFormat
	}
	return nil
}
