package launcher

import (
	"fmt"
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

func NewPocketClient(executableLocation string) PocketClient {
	return &pocketClient{
		executableLocation: executableLocation,
	}
}

type pocketClient struct {
	executableLocation string
}

func (pc *pocketClient) RunCommand(commandAndArgs ...string) (*CommandResult, error) {
	fmt.Printf("Running Command: %v\n", commandAndArgs)
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
