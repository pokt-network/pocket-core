package core

import "errors"

const (
	blockhash               = "blockhash"
	nodelist                = "nodelist"
	requestedchain          = "requested chain"
	isNotInSeed             = " is empty or nil"
	isInvalidFormat         = " is not in the correct format"
	insufficientNodeString  = "not enough nodes to fulfill the session"
	incompleteSessionString = "invalid session, missing information needed for key generation"
	blockchainHash          = " the blockchain"
	payload                 = " the data payload"
	signature               = " the client signature"
	devId                   = " the developer id"
	publickey               = " the public key"
	privatekey              = " the private key"
	dToken                  = " the developer token"
	path                    = " the http path"
	byteArrays              = " byte arrays"
	isMissing               = " is empty or nil"
	isInvalid               = " is not valid"
	isNotSupported          = " is not supported"
	incompatibleSizes       = " have different sizes"
	hexString               = " the hex string"
	DefaultHTTPMethod       = "POST"
	UnreachableAt           = " is not reachable at "
	session                 = " the session"
)

var (
	NoDevID                    = errors.New(devId + isNotInSeed)
	NoBlockHash                = errors.New(blockhash + isNotInSeed)
	NoNodeList                 = errors.New(nodelist + isNotInSeed)
	NoReqChain                 = errors.New(requestedchain + isNotInSeed)
	InvalidBlockHashFormat     = errors.New(blockhash + isInvalidFormat)
	InvalidDevIDFormat         = errors.New(devId + isInvalidFormat)
	InsufficientNodes          = errors.New(insufficientNodeString)
	IncompleteSession          = errors.New(incompleteSessionString)
	MissingBlockchainError     = errors.New(blockchainHash + isMissing)
	MissingPayloadError        = errors.New(payload + isMissing) // TODO is it possible for a payload to be empty and it be expected behavior?
	MissingSignatureError      = errors.New(signature + isMissing)
	MissingDevidError          = errors.New(devId + isMissing)
	MissingPathError           = errors.New(path + isMissing) // not used
	InvalidTokenError          = errors.New(dToken + isInvalid)
	InvalidDevIDError          = errors.New(devId + isInvalid)
	UnsupportedBlockchainError = errors.New(blockchainHash + isNotSupported)
	MismatchedByteArraysError  = errors.New(byteArrays + incompatibleSizes)
	InvalidPublicKeyError      = errors.New(publickey + isInvalid)
	InvalidPrivateKeyError     = errors.New(privatekey + isInvalid)
	InvalidHexStringError      = errors.New(hexString + isInvalid)
	InvalidSessionError        = errors.New(session + isInvalid)
)

func UnreachableBlockchainErr(blockchainHash, url string) error {
	return errors.New(blockchainHash + UnreachableAt + url)
}
