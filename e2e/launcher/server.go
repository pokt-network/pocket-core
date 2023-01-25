package launcher

import (
	"errors"
	"os/exec"
)

type Stream int

// matching standard file descriptors
const (
	StdIn  = Stream(0)
	StdOut = Stream(1)
	StdErr = Stream(2)
)

type PocketServer interface {
	Start(...string) error
	Kill() error
	RegisterPatternActor(patternActor PatternActor, stream Stream) error
}

func NewPocketServer(executablePath string) PocketServer {
	return &pocketServer{
		executablePath: executablePath,
		pocketInstance: nil,
	}
}

var _ PocketServer = &pocketServer{}

type pocketServer struct {
	executablePath string
	pocketInstance *exec.Cmd
	stdOutPipeline *PatternActorPipeline
	stdErrPipeline *PatternActorPipeline
}

func (ps *pocketServer) Start(args ...string) error {
	if ps.pocketInstance != nil {
		return errors.New("pocket instance already started")
	}

	ps.stdOutPipeline = &PatternActorPipeline{patternActors: []PatternActor{}}
	ps.stdErrPipeline = &PatternActorPipeline{patternActors: []PatternActor{}}

	ps.pocketInstance = exec.Command(ps.executablePath, append([]string{"start"}, args...)...)
	ps.pocketInstance.Stdout = ps.stdOutPipeline
	ps.pocketInstance.Stderr = ps.stdErrPipeline

	return ps.pocketInstance.Start()
}

func (ps *pocketServer) Kill() error {
	return ps.pocketInstance.Process.Kill()
}

func (ps *pocketServer) RegisterPatternActor(patternActor PatternActor, stream Stream) error {
	switch stream {
	case StdOut:
		ps.stdOutPipeline.patternActors = append(ps.stdOutPipeline.patternActors, patternActor)
	case StdErr:
		ps.stdErrPipeline.patternActors = append(ps.stdOutPipeline.patternActors, patternActor)
	default:
		return errors.New("invalid stream")
	}
	return nil
}
