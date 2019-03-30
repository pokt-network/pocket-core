package session

import "errors"

const (
	IsNotInSeed     = " is empty or nil "
	IsInvalidFormat = " is not in the correct format"
)

var (
	NoDevID                = errors.New("devid" + IsNotInSeed)
	NoBlockHash            = errors.New("blockhash" + IsNotInSeed)
	NoNodeList             = errors.New("nodelist" + IsNotInSeed)
	NoReqChain             = errors.New("requestedchain" + IsNotInSeed)
	InvalidBlockHashFormat = errors.New("blockhash" + IsInvalidFormat)
	InvalidDevIDFormat     = errors.New("devid" + IsInvalidFormat)
	InsufficientNodes      = errors.New(" not enough nodes to fulfill session")
	IncompleteSession      = errors.New("invalid session, missing information needed for key generation")
	NoSessionKey = errors.New("no session key is found")
)
