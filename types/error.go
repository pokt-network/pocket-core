package types

import (
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

var (
	InvalidHostedChainError = errors.New("invalid hosted chain error")
)

// NewError - create an error
func NewError(code sdk.CodeType, msg string) sdk.Error {
	return sdk.NewError(POCKETERRORCODESPACE, code, msg)
}
