package types

import (
	"errors"
	sdk "github.com/pokt-network/posmint/types"
)

var (
	InvalidHostedChainError = errors.New("invalid hosted chain error")
	ChainNotHostedError     = errors.New("the blockchain requested is not hosted")
)

const (
	CodeChainNotHosted = 1
)

// NewError - create an error
func NewError(code sdk.CodeType, msg string) sdk.Error {
	return sdk.NewError(POCKETERRORCODESPACE, code, msg)
}

func NewErrorChainNotHostedError() sdk.Error {
	return NewError(CodeChainNotHosted, ChainNotHostedError.Error())
}
