package dispatch

import "errors"

var (
	SessionGenerationError           = errors.New("dispatch session generation errored out: ")
	InsufficientAlivePeersError      = errors.New("there is an insufficient amount of alive peers")
	InsufficientTendermintPeersError = errors.New("tenermint did not return any peers")
)

func NewSessionGenerationError(err error) error {
	return errors.New(SessionGenerationError.Error() + err.Error())
}
