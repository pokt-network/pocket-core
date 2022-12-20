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

func NewPocketClient(executablePath string, verbose bool) PocketClient {
	return &pocketClient{
		executablePath: executablePath,
		verbose:        verbose,
	}
}

type pocketClient struct {
	executablePath string
	verbose        bool
}

func (pc *pocketClient) RunCommand(commandAndArgs ...string) (*CommandResult, error) {
	if pc.verbose {
		log.Printf("Running Command: %v\n", commandAndArgs)
	}
	cmd := exec.Command(pc.executablePath, commandAndArgs...)

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
