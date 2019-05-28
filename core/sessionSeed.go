package core

import (
	"reflect"
)

type SessionSeed struct {
	DevID          []byte
	BlockHash      []byte
	RequestedChain []byte
	Capacity       int
	NodeList       []Node
}

// "NewSessionSeed" is the constructor of the sessionSeed
func NewSessionSeed(devID []byte, nodePoolFilePath string, requestedBlockchain []byte, blockHash []byte, capacity int) (SessionSeed, error) {
	np, err := FileToNodes(nodePoolFilePath)
	return SessionSeed{DevID: devID, BlockHash: blockHash, RequestedChain: requestedBlockchain, NodeList: np, Capacity: capacity}, err
}

// "ErrorCheck()" checks all of the fields of a seed to ensure that it is considered initially valid
func (s *SessionSeed) ErrorCheck() error {
	if s.DevID == nil || len(s.DevID) == 0 {
		return NoDevIDError
	}
	if s.BlockHash == nil || len(s.BlockHash) == 0 {
		return NoBlockHashError
	}
	if s.Capacity == 0 {
		return NoCapacityError
	}
	if s.NodeList == nil || len(s.NodeList) == 0 {
		return NoNodeListError
	}
	if reflect.DeepEqual(s.NodeList[0], Node{}) {
		return InsufficientNodesError
	}
	if s.RequestedChain == nil || len(s.RequestedChain) == 0 {
		return NoReqChainError
	}
	if len(s.BlockHash) != 32 {
		return InvalidBlockHashFormatError
	}
	if len(s.DevID) != 32 {
		return InvalidDevIDFormatError
	}
	return nil
}
