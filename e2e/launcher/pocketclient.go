package launcher

import (
	"log"
	"os/exec"
)

type CommandResult struct {
	Stdout string
	Stderr string
	Err    error
}

type writerByteArray []byte

func (so *writerByteArray) Write(p []byte) (n int, err error) {
	*so = append(*so, p...)
	return len(p), nil
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

	so := &writerByteArray{}
	se := &writerByteArray{}

	cmd.Stdout = so
	cmd.Stderr = se
	err := cmd.Run()

	return &CommandResult{
		Stdout: string(*so),
		Stderr: string(*se),
		Err:    err,
	}, nil
}
