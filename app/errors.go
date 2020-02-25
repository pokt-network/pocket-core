package app

import "errors"

var (
	UninitializedKeybaseError = errors.New(`no keys stored in keybase, create a key pair by using "./main accounts create"`)
	InvalidChainsError        = errors.New("invalid chains.json")
)

func NewInvalidChainsError(err error) error {
	return errors.New(InvalidChainsError.Error() + ": " + err.Error())
}
