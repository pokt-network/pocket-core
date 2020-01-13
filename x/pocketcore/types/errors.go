package types

import (
	"errors"
	sdk "github.com/pokt-network/posmint/types"
	"strconv"
)

const ( // todo re-number
	CodeSessionGenerationError           = 1111
	CodeHttpStatusCodeError              = 1112
	CodeInvalidTokenError                = 1114
	CodePublKeyDecodeError               = 1116
	CodeEmptyChainError                  = 1118
	CodeEmptyBlockIDError                = 1119
	CodeAppPubKeyError                   = 1120
	CodeEmptyProofsError                 = 1121
	CodeUnsupportedBlockchainAppError    = 1123
	CodeInvalidSessionError              = 1124
	CodeInsufficientNodesError           = 1127
	CodeEmptyNonNativeChainError         = 1128
	CodeInvalidSessionKeyError           = 1129
	CodeFilterNodesError                 = 1130
	CodeXORError                         = 1131
	CodeInvalidHashError                 = 1132
	CodeEmptyBlockHashError              = 1133
	CodeEmptyBlockchainError             = 1134
	CodeEmptyPayloadDataError            = 1135
	CodeUnsupportedBlockchainNodeError   = 1136
	CodeNotStakedBlockchainError         = 1137
	CodeHTTPExecutionError               = 1138
	CodeInvalidIncrementCounterError     = 1139
	CodeEmptyResponseError               = 1140
	CodeResponseSignatureError           = 1141
	CodeNegativeICCounterError           = 1142
	CodeMaximumIncrementCounterError     = 1143
	CodeInvalidNodePubKeyError           = 1144
	CodeTicketsNotFoundError             = 1145
	CodeDuplicateTicketError             = 1146
	CodeDuplicateProofError              = 1147
	CodeInvalidSignatureSizeError        = 1148
	CodeSigDecodeError                   = 1149
	CodeMsgDecodeError                   = 1150
	CodeInvalidSigError                  = 1151
	CodePubKeySizeError                  = 1152
	CodeEmptyKeybaseError                = 1153
	CodeSelfNotFoundError                = 1154
	CodeAppNotFoundError                 = 1155
	CodeChainNotHostedError              = 1156
	CodeInvalidHostedChainsError         = 1157
	CodeNodeNotFoundError                = 1158
	CodeInvalidProofsError               = 1159
	CodeInconsistentPubKeyError          = 1160
	CodeInvalidChainParamsError          = 1161
	CodeNewHexDecodeError                = 1162
	CodeChainNotSupportedErr             = 1163
	CodePubKeyError                      = 1164
	CodeSignatureError                   = 1165
	CodeInvalidChainError                = 1166
	CodeJSONMarshalError                 = 1167
	CodeInvalidBlockchainHashLengthError = 1168
	CodeEmptySessionKeyError             = 1169
	CodeInvalidBlockHeightError          = 1170
	CodeInvalidAppPubKeyError            = 1171
	CodeInvalidHashLengthError           = 1172
	CodeInvalidLeafCousinProofsCombo     = 1173
	CodeEmptyAddressError                = 1174
	CodeClaimNotFoundError               = 1175
	CodeInvalidMerkleVerifyError         = 1176
	CodeEmptyMerkleTreeError             = 1177
	CodeMerkleNodeNotFoundError          = 1178
	CodeExpiredProofsSubmissionError     = 1179
	CodeAddressError                     = 1180
	CodeOverServiceError                 = 1181
	CodeCousinLeafEquivalentError        = 1182
	CodeInvalidRootError                 = 1183
)

var (
	MissingTokenVersionError         = errors.New("the application authentication token version is missing")
	UnsupportedTokenVersionError     = errors.New("the application authentication token version is not supported")
	MissingApplicationPublicKeyError = errors.New("the applicaiton public key included in the AAT is not valid")
	MissingClientPublicKeyError      = errors.New("the client public key included in the AAT is not valid")
	InvalidTokenSignatureErorr       = errors.New("the application signature on the AAT is not valid")
	NegativeICCounterError           = errors.New("the IC counter is less than 0")
	MaximumIncrementCounterError     = errors.New("the increment counter exceeds the maximum allowed relays")
	InvalidNodePubKeyError           = errors.New("the node public key in the service RelayProof does not match this nodes public key")
	InvalidTokenError                = errors.New("the application authentication token is invalid")
	EmptyProofsError                 = errors.New("the service proofs object is empty")
	DuplicateProofError              = errors.New("the RelayProof at index[increment counter] is not empty")
	InvalidIncrementCounterError     = errors.New("the increment counter included in the relay request is invalid")
	EmptyResponseError               = errors.New("the relay response payload is empty")
	ResponseSignatureError           = errors.New("response signing errored out: ")
	EmptyBlockchainError             = errors.New("the blockchain included in the relay request is empty")
	EmptyPayloadDataError            = errors.New("the payload data of the relay request is empty")
	UnsupportedBlockchainAppError    = errors.New("the blockchain in the relay request is not supported for this app")
	UnsupportedBlockchainNodeError   = errors.New("the blockchain in the relay request is not supported on this node")
	HttpStatusCodeError              = errors.New("HTTP status code returned not okay: ")
	InvalidSessionError              = errors.New("this node (self) is not responsible for this session provided by the client")
	ServiceSessionGenerationError    = errors.New("unable to generate a session for the seed data: ")
	NotStakedBlockchainError         = errors.New("the blockchain is not staked for this application")
	EmptyAppPubKeyError              = errors.New("the public key of the application is of length 0")
	EmptyNonNativeChainError         = errors.New("the non-native chain is of length 0")
	EmptyBlockIDError                = errors.New("the block addr is of length 0")
	InsufficientNodesError           = errors.New("there are less than the minimum session nodes found")
	EmptySessionKeyError             = errors.New("the session key passed is of length 0")
	MismatchedByteArraysError        = errors.New("the byte arrays are not of the same length")
	FilterNodesError                 = errors.New("unable to filter nodes: ")
	XORError                         = errors.New("error XORing the keys: ")
	PubKeyDecodeError                = errors.New("error decoding the string into hex bytes")
	InvalidHashError                 = errors.New("the hash is invalid: ")
	HTTPExecutionError               = errors.New("error executing the http request: ")
	TicketsNotFoundError             = errors.New("the tickets requested could not be found")
	DuplicateTicketError             = errors.New("the ticket is a duplicate")
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
	InconsistentPubKeyError          = errors.New("the public keys in the proofs are inconsistent")
	InvalidChainParamsError          = errors.New("the required params for a nonNative blockchain are invalid")
	HexDecodeError                   = errors.New("the hex string could not be decoded: ")
	ChainNotSupportedErr             = errors.New("the chain is not pocket supported")
	PubKeyError                      = errors.New("could not convert hex string to pub key: ")
	SignatureError                   = errors.New("there was a problem signing the message: ")
	InvalidChainError                = errors.New("the non native chain passed was invalid: ")
	JSONMarshalError                 = errors.New("unable to marshal object into json: ")
	InvalidBlockchainHashLength      = errors.New("the addr length is invalid")
	InvalidBlockHeightError          = errors.New("the block height passed has been invalid")
	InvalidAppPubKeyError            = errors.New("the app public key is invalid")
	InvalidHashLengthError           = errors.New("the addr length is not valid")
	InvalidLeafCousinProofsCombo     = errors.New("the merkle relayProof combo for the cousin and leaf is invalid")
	EmptyAddressError                = errors.New("the address provided is empty")
	ClaimNotFoundError               = errors.New("the unverified RelayProof was not found for the key given")
	InvalidMerkleVerifyError         = errors.New("claim resulted in an invalid merkle RelayProof")
	EmptyMerkleTreeError             = errors.New("the merkle tree is empty")
	NodeNotFoundError                = errors.New("the node of the merkle tree requested is not found")
	ExpiredProofsSubmissionError     = errors.New("the opportunity of window to submit the RelayProof has closed because the secret has been revealed")
	AddressError                     = errors.New("the address is invalid")
	OverServiceError                 = errors.New("the max number of relays serviced for this node is exceeded")
	UninitializedKeybaseError        = errors.New("the keybase is nil")
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
