package types

import (
	"errors"
	sdk "github.com/pokt-network/posmint/types"
	"strconv"
)

const (
	CodeHttpStatusCodeError            = 100
	CodeInvalidTokenError              = 101
	CodePublKeyDecodeError             = 102
	CodeEmptyChainError                = 103
	CodeEmptyBlockIDError              = 104
	CodeEmptyProofsError               = 105
	CodeUnsupportedBlockchainAppError  = 106
	CodeInvalidSessionError            = 107
	CodeInsufficientNodesError         = 108
	CodeEmptyNonNativeChainError       = 109
	CodeInvalidSessionKeyError         = 110
	CodeFilterNodesError               = 111
	CodeXORError                       = 112
	CodeInvalidHashError               = 113
	CodeEmptyBlockHashError            = 114
	CodeEmptyPayloadDataError          = 115
	CodeUnsupportedBlockchainNodeError = 116
	CodeHTTPExecutionError             = 117
	CodeInvalidIncrementCounterError   = 118
	CodeEmptyResponseError             = 119
	CodeResponseSignatureError         = 120
	CodeInvalidNodePubKeyError         = 121
	CodeDuplicateProofError            = 122
	CodeInvalidSignatureSizeError      = 123
	CodeSigDecodeError                 = 124
	CodeMsgDecodeError                 = 125
	CodeInvalidSigError                = 126
	CodePubKeySizeError                = 127
	CodeEmptyKeybaseError              = 128
	CodeSelfNotFoundError              = 129
	CodeAppNotFoundError               = 130
	CodeChainNotHostedError            = 131
	CodeInvalidHostedChainsError       = 132
	CodeNodeNotFoundError              = 133
	CodeInvalidProofsError             = 134
	CodeInvalidChainParamsError        = 135
	CodeNewHexDecodeError              = 136
	CodeChainNotSupportedErr           = 137
	CodePubKeyError                    = 138
	CodeSignatureError                 = 139
	CodeJSONMarshalError               = 140
	CodeInvalidBlockHeightError        = 141
	CodeInvalidAppPubKeyError          = 142
	CodeInvalidHashLengthError         = 143
	CodeInvalidLeafCousinProofsCombo   = 144
	CodeEmptyAddressError              = 145
	CodeClaimNotFoundError             = 146
	CodeInvalidMerkleVerifyError       = 147
	CodeEmptyMerkleTreeError           = 148
	CodeMerkleNodeNotFoundError        = 149
	CodeExpiredProofsSubmissionError   = 150
	CodeAddressError                   = 151
	CodeOverServiceError               = 152
	CodeCousinLeafEquivalentError      = 153
	CodeInvalidRootError               = 154
)

var (
	MissingTokenVersionError         = errors.New("the application authentication token version is missing")
	UnsupportedTokenVersionError     = errors.New("the application authentication token version is not supported")
	MissingApplicationPublicKeyError = errors.New("the applicaiton public key included in the AAT is not valid")
	MissingClientPublicKeyError      = errors.New("the client public key included in the AAT is not valid")
	InvalidTokenSignatureErorr       = errors.New("the application signature on the AAT is not valid")
	InvalidNodePubKeyError           = errors.New("the node public key in the service RelayProof does not match this nodes public key")
	InvalidTokenError                = errors.New("the application authentication token is invalid")
	EmptyProofsError                 = errors.New("the service proofs object is empty")
	DuplicateProofError              = errors.New("the RelayProof at index[increment counter] is not empty")
	InvalidIncrementCounterError     = errors.New("the increment counter included in the relay request is invalid")
	EmptyResponseError               = errors.New("the relay response payload is empty")
	ResponseSignatureError           = errors.New("response signing errored out: ")
	EmptyPayloadDataError            = errors.New("the payload data of the relay request is empty")
	UnsupportedBlockchainAppError    = errors.New("the blockchain in the relay request is not supported for this app")
	UnsupportedBlockchainNodeError   = errors.New("the blockchain in the relay request is not supported on this node")
	HttpStatusCodeError              = errors.New("HTTP status code returned not okay: ")
	InvalidSessionError              = errors.New("this node (self) is not responsible for this session provided by the client")
	EmptyNonNativeChainError         = errors.New("the non-native chain is of length 0")
	EmptyBlockIDError                = errors.New("the block addr is of length 0")
	InsufficientNodesError           = errors.New("there are less than the minimum session nodes found")
	MismatchedByteArraysError        = errors.New("the byte arrays are not of the same length")
	FilterNodesError                 = errors.New("unable to filter nodes: ")
	XORError                         = errors.New("error XORing the keys: ")
	PubKeyDecodeError                = errors.New("error decoding the string into hex bytes")
	InvalidHashError                 = errors.New("the hash is invalid: ")
	HTTPExecutionError               = errors.New("error executing the http request: ")
	InvalidSignatureSizeError        = errors.New("the signature length is invalid")
	MessageDecodeError               = errors.New("the message could not be hex decoded")
	SigDecodeError                   = errors.New("the signature could not be message decoded")
	InvalidSignatureError            = errors.New("the signature could not be verified with the message and pub key")
	PubKeySizeError                  = errors.New("the public key is not the correct size")
	KeybaseError                     = errors.New("the keybase is invalid: ")
	SelfNotFoundError                = errors.New("the self node is not within the world state")
	AppNotFoundError                 = errors.New("the app could not be found in the world state")
	InvalidHostedChainError          = errors.New("invalid hosted chain error")
	ChainNotHostedError              = errors.New("the blockchain requested is not hosted")
	NodeNotFoundErr                  = errors.New("the node is not found in world state")
	InvalidProofsError               = errors.New("the proofs provided are invalid")
	InvalidChainParamsError          = errors.New("the required params for a nonNative blockchain are invalid")
	HexDecodeError                   = errors.New("the hex string could not be decoded: ")
	ChainNotSupportedErr             = errors.New("the chain is not pocket supported")
	PubKeyError                      = errors.New("could not convert hex string to pub key: ")
	SignatureError                   = errors.New("there was a problem signing the message: ")
	JSONMarshalError                 = errors.New("unable to marshal object into json: ")
	InvalidBlockHeightError          = errors.New("the block height passed has been invalid")
	InvalidAppPubKeyError            = errors.New("the app public key is invalid")
	InvalidHashLengthError           = errors.New("the addr length is not valid")
	InvalidLeafCousinProofsCombo     = errors.New("the merkle relayProof combo for the cousin and leaf is invalid")
	EmptyAddressError                = errors.New("the address provided is empty")
	ClaimNotFoundError               = errors.New("the unverified RelayProof was not found for the key given")
	InvalidMerkleVerifyError         = errors.New("claim resulted in an invalid merkle RelayProof")
	EmptyMerkleTreeError             = errors.New("the merkle tree is empty")
	ExpiredProofsSubmissionError     = errors.New("the opportunity of window to submit the RelayProof has closed because the secret has been revealed")
	AddressError                     = errors.New("the address is invalid")
	OverServiceError                 = errors.New("the max number of relays serviced for this node is exceeded")
	CousinLeafEquivalentError        = errors.New("the cousin and leaf cannot be equal")
	InvalidRootError                 = errors.New("the merkle root passed is invalid")
	MerkleNodeNotFoundError          = errors.New("the merkle node cannot be found")
)

func NewOverServiceError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeOverServiceError, OverServiceError.Error())
}

func NewAddressInvalidLengthError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeAddressError, AddressError.Error())
}

func NewExpiredProofsSubmissionError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeExpiredProofsSubmissionError, ExpiredProofsSubmissionError.Error())
}

func NewMerkleNodeNotFoundError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMerkleNodeNotFoundError, MerkleNodeNotFoundError.Error())
}

func NewEmptyMerkleTreeError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyMerkleTreeError, EmptyMerkleTreeError.Error())
}

func NewInvalidMerkleVerifyError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMerkleVerifyError, InvalidMerkleVerifyError.Error())
}

func NewClaimNotFoundError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeClaimNotFoundError, ClaimNotFoundError.Error())
}

func NewEmptyAddressError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyAddressError, EmptyAddressError.Error())
}

func NewCousinLeafEquivalentError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeCousinLeafEquivalentError, CousinLeafEquivalentError.Error())
}

func NewInvalidLeafCousinProofsComboError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidLeafCousinProofsCombo, InvalidLeafCousinProofsCombo.Error())
}

func NewInvalidRootError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidRootError, InvalidRootError.Error())
}

func NewInvalidHashLengthError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidHashLengthError, InvalidHashLengthError.Error())
}
func NewInvalidAppPubKeyError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidAppPubKeyError, InvalidAppPubKeyError.Error())
}

func NewInvalidBlockHeightError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidBlockHeightError, InvalidBlockHeightError.Error())
}

func NewJSONMarshalError(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeJSONMarshalError, JSONMarshalError.Error()+err.Error())
}

func NewSignatureError(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeSignatureError, SignatureError.Error()+err.Error())
}

func NewPubKeyError(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodePubKeyError, PubKeyError.Error()+err.Error())
}

func NewChainNotSupportedErr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeChainNotSupportedErr, ChainNotSupportedErr.Error())
}

func NewHexDecodeError(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeNewHexDecodeError, HexDecodeError.Error()+err.Error())
}

func NewInvalidChainParamsError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidChainParamsError, InvalidChainParamsError.Error())
}

func NewInvalidProofsError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidProofsError, InvalidProofsError.Error())
}

func NewNodeNotFoundErr(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeNodeNotFoundError, NodeNotFoundErr.Error())
}

func NewInvalidHostedChainError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidHostedChainsError, InvalidHostedChainError.Error())
}

func NewErrorChainNotHostedError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeChainNotHostedError, ChainNotHostedError.Error())
}

func NewAppNotFoundError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeAppNotFoundError, AppNotFoundError.Error())
}

func NewSelfNotFoundError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSelfNotFoundError, SelfNotFoundError.Error())
}

func NewKeybaseError(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyKeybaseError, KeybaseError.Error()+err.Error())
}

func NewPubKeySizeError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodePubKeySizeError, PubKeySizeError.Error())
}

func NewInvalidSignatureError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSigError, InvalidSignatureError.Error())
}

func NewMsgDecodeError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeMsgDecodeError, MessageDecodeError.Error())
}

func NewSigDecodeError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeSigDecodeError, SigDecodeError.Error())
}

func NewInvalidSignatureSizeError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSignatureSizeError, InvalidSignatureSizeError.Error())
}

func NewDuplicateProofError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeDuplicateProofError, DuplicateProofError.Error())
}

func NewInvalidNodePubKeyError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidNodePubKeyError, InvalidNodePubKeyError.Error())
}

func NewResponseSignatureError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeResponseSignatureError, ResponseSignatureError.Error())
}

func NewEmptyResponseError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyResponseError, EmptyResponseError.Error())
}

func NewInvalidIncrementCounterError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidIncrementCounterError, InvalidIncrementCounterError.Error())
}

func NewHTTPExecutionError(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeHTTPExecutionError, HTTPExecutionError.Error()+err.Error())
}

func NewUnsupportedBlockchainNodeError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeUnsupportedBlockchainNodeError, UnsupportedBlockchainNodeError.Error())
}

func NewEmptyPayloadDataError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyPayloadDataError, EmptyPayloadDataError.Error())
}

func NewInvalidHashError(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidHashError, InvalidHashError.Error()+err.Error())
}

func NewEmptyHashError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyBlockHashError, InvalidHashError.Error())
}

func NewPubKeyDecodeError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodePublKeyDecodeError, PubKeyDecodeError.Error())
}

func NewXORError(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeXORError, XORError.Error()+err.Error())
}

func NewFilterNodesError(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeFilterNodesError, FilterNodesError.Error()+err.Error())
}

func NewInvalidSessionKeyError(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSessionKeyError, InvalidSessionError.Error()+err.Error())
}

func NewEmptyNonNativeChainError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyNonNativeChainError, EmptyNonNativeChainError.Error())
}

func NewInsufficientNodesError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInsufficientNodesError, InsufficientNodesError.Error())
}

func NewInvalidSessionError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidSessionError, InvalidSessionError.Error())
}

func NewUnsupportedBlockchainAppError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeUnsupportedBlockchainAppError, UnsupportedBlockchainAppError.Error())
}

func NewEmptyProofsError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyProofsError, EmptyProofsError.Error())
}

func NewEmptyBlockIDError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyBlockIDError, EmptyBlockIDError.Error())
}
func NewEmptyChainError(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyChainError, EmptyNonNativeChainError.Error())
}

func NewHTTPStatusCodeError(codespace sdk.CodespaceType, statusCode int) sdk.Error {
	return sdk.NewError(codespace, CodeHttpStatusCodeError, HttpStatusCodeError.Error()+strconv.Itoa(statusCode))
}

func NewInvalidTokenError(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidTokenError, InvalidTokenError.Error()+" : "+err.Error())
}
