package service

import (
	"errors"
	"strconv"
)

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
	UnsupportedPayloadTypeError      = errors.New("the payload type is not supported")
	EmptyResponseError               = errors.New("the relay response payload is empty")
	ResponseSignatureError           = errors.New("response signing errored out: ")
	EmptyHostedChainsError           = errors.New("the hosted chains object is of length 0")
	HttpStatusCodeError              = errors.New("HTTP status code returned not okay: ")
	InvalidNodePubKeyError           = errors.New("the node public key in the service certificate does not match this nodes public key")
	ServiceCertificateHashError      = errors.New("the service certificate object was unable to be hashed: ")
	InvalidSessionError              = errors.New("this node (self) is not responsible for this session provided by the client")
	ServiceSessionGenerationError    = errors.New("unable to generate a session for the seed data: ")
	BlockHashHexDecodeError          = errors.New("the block hash was unable to be decoded into hex format")
	ServiceCertificateError          = errors.New("the service is unauthorized: ")
	EmptyEvidenceError               = errors.New("the evidence object type([]ServiceCertificate) is nil or empty")
	InvalidEvidenceSizeError         = errors.New("the size of the evidence container is less than the counter")
	DuplicateEvidenceError           = errors.New("DuplicateEvidenceError: the evidence is already recorded for that increment counter")
	RelayBatchCreationError          = errors.New("there was a problem creating a new relay batch: ")
)

func NewRelayBatchCreationError(err error) error {
	return errors.New(RelayBatchCreationError.Error() + err.Error())
}

func NewServiceCertificateError(err error) error {
	return errors.New(ServiceCertificateError.Error() + err.Error())
}
func NewBlockHashHexDecodeError(err error) error {
	return errors.New(BlockHashHexDecodeError.Error() + err.Error())
}

func NewSignatureError(err error) error {
	return errors.New(ResponseSignatureError.Error() + err.Error())
}

func NewServiceSessionGenerationError(err error) error {
	return errors.New(ServiceSessionGenerationError.Error() + err.Error())
}

func NewHTTPStatusCodeError(statusCode int) error {
	return errors.New(HttpStatusCodeError.Error() + strconv.Itoa(statusCode))
}

func NewInvalidTokenError(err error) error {
	return errors.New(InvalidTokenError.Error() + " : " + err.Error())
}

func NewServiceCertificateHashError(err error) error {
	return errors.New(ServiceCertificateHashError.Error() + err.Error())
}

func NewClientPubKeyDecodeError(err error) error {
	return errors.New(ClientPubKeyDecodeError.Error() + " : " + err.Error())
}
