package types

import (
	"errors"
	"strconv"
)

var (
	MissingTokenVersionError         = errors.New("the application authentication token version is missing")
	UnsupportedTokenVersionError     = errors.New("the application authentication token version is not supported")
	MissingApplicationPublicKeyError = errors.New("the applicaiton public key included in the AAT is not valid")
	MissingClientPublicKeyError      = errors.New("the client public key included in the AAT is not valid")
	InvalidTokenSignatureErorr       = errors.New("the application signature on the AAT is not valid")
	NegativeICCounterError           = errors.New("the IC counter is less than 0")
	InvalidICSignatureError          = errors.New("the client signature is not valid on the increment counter")
	InvalidICError                   = errors.New("the increment counter proof provided does not match the needed proof")
	MaximumIncrementCounterError     = errors.New("the increment counter exceeds the maximum allowed relays")
	InvalidNodePubKeyError           = errors.New("the node public key in the service proof does not match this nodes public key")
	ClientPubKeyDecodeError          = errors.New("unable to hex.Decode( clientPublicKey )")
	InvalidTokenError                = errors.New("the application authentication token is invalid")
	ServiceProofHashError            = errors.New("the service proof object was unable to be hashed: ")
	EmptyProofsError                 = errors.New("the service proofs object is empty")
	InvalidProofSizeError            = errors.New("the size of the proofs object is bigger than the max number of relays")
	DuplicateProofError              = errors.New("the proof at index[increment counter] is not empty")
	BatchCreationErr                 = errors.New("there was a problem creating the proof batch: ")
	InvalidIncrementCounterError     = errors.New("the increment counter included in the relay request is invalid")
	EmptyResponseError               = errors.New("the relay response payload is empty")
	ResponseSignatureError           = errors.New("response signing errored out: ")
	EmptyBlockchainError             = errors.New("the blockchain included in the relay request is empty")
	EmptyPayloadDataError            = errors.New("the payload data of the relay request is empty")
	UnsupportedBlockchainError       = errors.New("the blockchain in the relay request is not supported on this node")
	UnsupportedPayloadTypeError      = errors.New("the payload type is not supported")
	HttpStatusCodeError              = errors.New("HTTP status code returned not okay: ")
	InvalidSessionError              = errors.New("this node (self) is not responsible for this session provided by the client")
	ServiceSessionGenerationError    = errors.New("unable to generate a session for the seed data: ")
	BlockHashHexDecodeError          = errors.New("the block hash was unable to be decoded into hex format")
	ServiceProofError                = errors.New("the service is unauthorized: ")
	NotStakedBlockchainError         = errors.New("the blockchain is not staked for this application")
	NotEveryICProvidedError          = errors.New("not every requested proof was provided by the node")
	EmptyAppPubKeyError              = errors.New("the public key of the application is of length 0")
	EmptyNonNativeChainError         = errors.New("the non-native chain is of length 0")
	EmptyBlockIDError                = errors.New("the block hash is of length 0")
	InsufficientNodesError           = errors.New("there are less than the minimum session nodes found")
	EmptySessionKeyError             = errors.New("the session key passed is of length 0")
	MismatchedByteArraysError        = errors.New("the byte arrays are not of the same length")
)

func NewServiceProofError(err error) error {
	return errors.New(ServiceProofError.Error() + err.Error())
}
func NewBlockHashHexDecodeError(err error) error {
	return errors.New(BlockHashHexDecodeError.Error() + err.Error())
}

func NewServiceSessionGenerationError(err error) error {
	return errors.New(ServiceSessionGenerationError.Error() + err.Error())
}

func NewHTTPStatusCodeError(statusCode int) error {
	return errors.New(HttpStatusCodeError.Error() + strconv.Itoa(statusCode))
}

func NewBatchCreationErr(err error) error {
	return errors.New(BatchCreationErr.Error() + err.Error())
}

func NewInvalidTokenError(err error) error {
	return errors.New(InvalidTokenError.Error() + " : " + err.Error())
}

func NewServiceProofHashError(err error) error {
	return errors.New(ServiceProofHashError.Error() + err.Error())
}

func NewClientPubKeyDecodeError(err error) error {
	return errors.New(ClientPubKeyDecodeError.Error() + " : " + err.Error())
}

func NewSignatureError(err error) error {
	return errors.New(ResponseSignatureError.Error() + err.Error())
}
