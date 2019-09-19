package service

import "errors"

var (
	NegativeICCounterError           = errors.New("the IC counter is less than 0")
	ClientPubKeyDecodeError          = errors.New("unable to hex.Decode( clientPublicKey )")
	InvalidICSignatureError          = errors.New("the client signature is not valid on the increment counter")
	InvalidTokenSignatureErorr       = errors.New("the application signature on the AAT is not valid")
	EmptyBlockchainError             = errors.New("the blockchain included in the relay request is empty")
	EmptyPayloadDataError            = errors.New("the payload data of the relay request is empty")
	InvalidTokenError                = errors.New("the application authentication token is invalid")
	InvalidIncrementCounterError     = errors.New("the increment counter included in the relay request is invalid")
	UnsupportedBlockchainError       = errors.New("the blockchain in the relay request is not supported on this node")
	MissingTokenVersionError         = errors.New("the application authentication token version is missing")
	UnsupportedTokenVersionError     = errors.New("the application authentication token version is not supported")
	MissingApplicationPublicKeyError = errors.New("the applicaiton public key included in the AAT is not valid")
	MissingClientPublicKeyError      = errors.New("the client public key included in the AAT is not valid")
)
