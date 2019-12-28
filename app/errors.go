package app

import "errors"

var (
	UninitializedKeybaseError = errors.New("uninitialized keybase")
	InvalidChainsError        = errors.New("invalid chains.json")
)

func NewInvalidChainsError(err error) error {
	return errors.New(InvalidChainsError.Error() + ": " + err.Error())
}
