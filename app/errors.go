package app

import (
	"errors"
)

var (
	UninitializedKeybaseError = errors.New(`no Keys stored in keybase, create a key pair by using "./main accounts create"`)
	InvalidChainsError        = errors.New("the chains.json file input is invalid")
)

func NewInvalidChainsError(err error) error {
	return errors.New(InvalidChainsError.Error() + ": " + err.Error())
}
