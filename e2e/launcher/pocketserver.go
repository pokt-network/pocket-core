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
	RegisterPatternActor(patternActor PatternAction, stream Stream) error
}

func NewPocketServer(executablePath string) PocketServer {
	return &pocketServer{
		executablePath: executablePath,
		pocketInstance: nil,
	}
}

type pocketServer struct {
	executablePath string
	pocketInstance *exec.Cmd
	stdOutPipeline *PatternActionPipeline
	stdErrPipeline *PatternActionPipeline
}

func (ps *pocketServer) Start(arguments ...string) error {
	if ps.pocketInstance != nil {
		return errors.New("pocket instance already started")
	}
	invocation := []string{"start"}
	invocation = append(invocation, arguments...)
	ps.pocketInstance = exec.Command(ps.executablePath, invocation...)
	ps.stdOutPipeline = &PatternActionPipeline{patternActors: []PatternAction{}}
	ps.stdErrPipeline = &PatternActionPipeline{patternActors: []PatternAction{}}
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

func (ps *pocketServer) RegisterPatternActor(patternActor PatternAction, stream Stream) error {
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

type PatternAction interface {
	MaybeAct(string)
}

type PatternActionPipeline struct {
	patternActors []PatternAction
}

func (r PatternActionPipeline) Write(p []byte) (n int, err error) {
	for _, cmd := range r.patternActors {
		cmd.MaybeAct(string(p))
	}
	return len(p), nil
}

type PrinterPatternAction struct {
}

func (*PrinterPatternAction) MaybeAct(line string) {
	fmt.Printf("Printer prints: %s\n", line)
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
