package launcher

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

type PocketServer interface {
	Start(...string) error
	Kill() error
	RegisterPatternActor(patternActor PatternActor, stream Stream) error
}

func NewPocketServer(executableLocation string) PocketServer {
	return &pocketServer{
		executableLocation: executableLocation,
		pocketInstance:     nil,
	}
}

type pocketServer struct {
	executableLocation string
	pocketInstance     *exec.Cmd
	stdOutPipeline     *PatternActionPipeline
	stdErrPipeline     *PatternActionPipeline
}

func (ps *pocketServer) Start(arguments ...string) error {
	if ps.pocketInstance != nil {
		return errors.New("pocket instance already started")
	}
	invocation := []string{"start"}
	invocation = append(invocation, arguments...)
	ps.pocketInstance = exec.Command(ps.executableLocation, invocation...)
	ps.stdOutPipeline = &PatternActionPipeline{patternActors: []PatternActor{}}
	ps.stdErrPipeline = &PatternActionPipeline{patternActors: []PatternActor{}}
	ps.pocketInstance.Stdout = ps.stdOutPipeline
	ps.pocketInstance.Stderr = ps.stdErrPipeline
	err := ps.pocketInstance.Start()
	return err
}

func (ps *pocketServer) Kill() error {
	return ps.pocketInstance.Process.Kill()
}

type Stream int

// matching standard file descriptors
const (
	StdIn  = Stream(0)
	StdOut = Stream(1)
	StdErr = Stream(2)
)

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

type PatternActor interface {
	MaybeAct(string)
}

type PatternActionPipeline struct {
	patternActors []PatternActor
}

func (r PatternActionPipeline) Write(p []byte) (n int, err error) {
	for _, cmd := range r.patternActors {
		cmd.MaybeAct(string(p))
	}
	return len(p), nil
}

type PrinterPatternActor struct {
}

func (*PrinterPatternActor) MaybeAct(line string) {
	fmt.Println("Printer prints: " + line)
}

func NewBlockWaiter(blocks int, verbose bool) *BlockWaiter {
	wg := new(sync.WaitGroup)
	wg.Add(1)
	return &BlockWaiter{
		Wg:              wg,
		RemainingBlocks: blocks,
		done:            false,
		verbose:         verbose,
	}
}

type BlockWaiter struct {
	Wg              *sync.WaitGroup
	RemainingBlocks int
	done            bool
	verbose         bool
}

func (b *BlockWaiter) MaybeAct(line string) {
	if b.Wg != nil && !b.done {
		if strings.Contains(line, "Executed block") {
			b.RemainingBlocks--
			if b.verbose {
				fmt.Printf("Block elapsed; remaining: %d\n", b.RemainingBlocks)
			}
		}
		if b.RemainingBlocks <= 0 {
			b.Wg.Done()
			b.done = true
		}
	}
}
