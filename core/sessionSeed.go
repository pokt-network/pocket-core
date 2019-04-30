package core

import (
	"reflect"
)

type SessionSeed struct {
	DevID          []byte
	BlockHash      []byte
	RequestedChain []byte
	NodeList       []Node
}

// "NewSessionSeed" is the constructor of the sessionSeed
func NewSessionSeed(devID []byte, nodePoolFilePath string, requestedBlockchain []byte, blockHash []byte) (SessionSeed, error) {
	np, err := FileToNodes(nodePoolFilePath)
	return SessionSeed{DevID: devID, BlockHash: blockHash, RequestedChain: requestedBlockchain, NodeList: np}, err
}

// "ErrorCheck()" checks all of the fields of a seed to ensure that it is considered initially valid
func (s *SessionSeed) ErrorCheck() error {
	if s.DevID == nil || len(s.DevID) == 0 {
		return NoDevID
	}
	if s.BlockHash == nil || len(s.BlockHash) == 0 {
		return NoBlockHash
	}
	if s.NodeList == nil || len(s.NodeList) == 0 {
		return NoNodeList
	}
	if reflect.DeepEqual(s.NodeList[0], Node{}) {
		return InsufficientNodes
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
