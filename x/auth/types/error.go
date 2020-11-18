package types

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
)

// Param module codespace constants
const (
	CodeInvalidMemo         sdk.CodeType = 1
	CodeEmptyPublicKey      sdk.CodeType = 2
	CodeAccNotFound         sdk.CodeType = 3
	CodeInsufficientFee     sdk.CodeType = 4
	CodeSignatureLimit      sdk.CodeType = 5
	CodeDupTx               sdk.CodeType = 6
	CodeInsufficientBalance sdk.CodeType = 7
	CodeTxIndexerNil        sdk.CodeType = 8
)

// ErrUnknownSubspace returns an unknown subspace error.
func ErrInvalidMemo(codespace sdk.CodespaceType, err error) sdk.Error {
	return sdk.NewError(codespace, CodeInvalidMemo, fmt.Sprintf("the transaction memo is invalid: %s", err.Error()))
}

func ErrEmptyPublicKey(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeEmptyPublicKey, "the public key in the transaction is empty and the public key cannot be found in the world state")
}

func ErrAccountNotFound(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeAccNotFound, "the account cannot be found in the world state")
}

func ErrNilTxIndexer(codespace sdk.CodespaceType) sdk.Error {
	return sdk.NewError(codespace, CodeTxIndexerNil, "the transaction indexer in the node is nil")
}

func ErrInsufficientFee(codespace sdk.CodespaceType, expectedFee sdk.Coins, actualFee sdk.Coins) sdk.Error {
	return sdk.NewError(codespace, CodeInsufficientFee, fmt.Sprintf("the fee for the transaction was insufficient. \nExpected: %s\nActual: %s", expectedFee.String(), actualFee.String()))
}

func ErrTooManySignatures(codespace sdk.CodespaceType, sigLimit uint64) sdk.Error {
	return sdk.NewError(codespace, CodeSignatureLimit, fmt.Sprintf("the limit for signatures (%d) is reached", sigLimit))
}

func ErrDuplicateTx(codespace sdk.CodespaceType, txHash string) sdk.Error {
	return sdk.NewError(codespace, CodeDupTx, fmt.Sprintf("the transaction hash is already found, this is a duplicate transaction: %s", txHash))
}

func ErrInsufficientBalance(codespace sdk.CodespaceType, signer sdk.Address, neededFee sdk.Coins) sdk.Error {
	return sdk.NewError(codespace, CodeDupTx, fmt.Sprintf("the signer account : %s, does not have enough coins for the tx. Need %s", signer, neededFee.String()))
}
