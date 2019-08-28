package legacy

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
	isZero                  = " is zero"
	incompatibleSizes       = " have different sizes"
	hexString               = " the hex string"
	DefaultHTTPMethod       = "POST"
	UnreachableAt           = " is not reachable at "
	session                 = " the session"
	capacity                = " the session capacity"
	nonce                   = " the relay nonce "
	blkNum                  = " the block number "
	relaySummaryCount       = " the relay summary count "
)

var (
	NoDevIDError                = errors.New(devId + isNotInSeed)
	NoBlockHashError            = errors.New(blockhash + isNotInSeed)
	NoNodeListError             = errors.New(nodelist + isNotInSeed)
	NoReqChainError             = errors.New(requestedchain + isNotInSeed)
	NoCapacityError             = errors.New(capacity + isZero)
	InvalidBlockHashFormatError = errors.New(blockhash + isInvalidFormat)
	InvalidDevIDFormatError     = errors.New(devId + isInvalidFormat)
	InsufficientNodesError      = errors.New(insufficientNodeString)
	IncompleteSessionError      = errors.New(incompleteSessionString)
	MissingBlockchainError      = errors.New(blockchainHash + isMissing)
	MissingPayloadError         = errors.New(payload + isMissing) // TODO is it possible for a payload to be empty and it be expected behavior?
	MissingSignatureError       = errors.New(signature + isMissing)
	MissingDevidError           = errors.New(devId + isMissing)
	MissingPathError            = errors.New(path + isMissing) // not used
	InvalidTokenError           = errors.New(dToken + isInvalid)
	InvalidDevIDError           = errors.New(devId + isInvalid)
	UnsupportedBlockchainError  = errors.New(blockchainHash + isNotSupported)
	MismatchedByteArraysError   = errors.New(byteArrays + incompatibleSizes)
	InvalidPublicKeyError       = errors.New(publickey + isInvalid)
	InvalidPrivateKeyError      = errors.New(privatekey + isInvalid)
	InvalidHexStringError       = errors.New(hexString + isInvalid)
	InvalidSessionError         = errors.New(session + isInvalid)
	ZeroNonceError              = errors.New(nonce + isZero)
	ZeroBlockError              = errors.New(blkNum + isZero)
	ZeroRelaySummaryError       = errors.New(relaySummaryCount + isZero)
)

func UnreachableBlockchainErr(blockchainHash, url string) error {
	return errors.New(blockchainHash + UnreachableAt + url)
}
