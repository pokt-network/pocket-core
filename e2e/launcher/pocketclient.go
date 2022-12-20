package launcher

import (
	"log"
	"os/exec"
	"strings"
)

type CommandResult struct {
	Stdout string
	Stderr string
	Err    error
}

type PocketClient interface {
	RunCommand(...string) (*CommandResult, error)
}

func NewPocketClient(executableLocation string, verbose bool) PocketClient {
	return &pocketClient{
		executableLocation: executableLocation,
		verbose:            verbose,
	}
}

type pocketClient struct {
	executableLocation string
	verbose            bool
}

func (pc *pocketClient) RunCommand(commandAndArgs ...string) (*CommandResult, error) {
	if pc.verbose {
		log.Printf("Running Command: %v\n", commandAndArgs)
	}
	cmd := exec.Command(pc.executableLocation, commandAndArgs...)

	so := &strings.Builder{}
	se := &strings.Builder{}

	cmd.Stdout = so
	cmd.Stderr = se
	err := cmd.Run()

	return &CommandResult{
		Stdout: so.String(),
		Stderr: se.String(),
		Err:    err,
	}, nil
}
