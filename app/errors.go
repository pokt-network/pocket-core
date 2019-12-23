package app

import "errors"

var (
	UninitializedKeybaseError    = errors.New("uninitialized keybase")
	InvalidChainsError           = errors.New("invalid chains.json")
	UninitializedTendermintError = errors.New("uninitialized tendermint node")
)

func NewInvalidChainsError(err error) error {
	return errors.New(InvalidChainsError.Error() + ": " + err.Error())
}
