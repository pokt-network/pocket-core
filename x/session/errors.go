package session

import (
	"errors"
	"strconv"
)

var (
	EmptyDevPubKeyError       = errors.New("the public key of the developer is of length 0")
	EmptyNonNativeChainError  = errors.New("the non-native chain is of length 0")
	EmptyBlockIDError         = errors.New("the block hash is of length 0")
	InsufficientNodesError    = errors.New("there are less than the minimum of " + strconv.FormatUint(uint64(SESSIONNODECOUNT), 10) + " nodes found")
	MismatchedByteArraysError = errors.New("the byte arrays are not of the same length")
)
