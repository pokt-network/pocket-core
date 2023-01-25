package launcher

import (
	"log"
	"strings"
	"sync"
)

// IMPROVE: Need to document how how build/add new pattern actors

// All PatternActors must implement this interface
type PatternActor interface {
	MaybeAct(string) // The action that the pattern actor may execute depending on the output being passed in
}

// A series of PatternActors
type PatternActorPipeline struct {
	patternActors []PatternActor
}

func (r PatternActorPipeline) Write(p []byte) (n int, err error) {
	for _, cmd := range r.patternActors {
		cmd.MaybeAct(string(p))
	}
	return len(p), nil
}

// PrinterPatternActor - A pattern actor that prints all output to the screen

type PrinterPatternActor interface {
	PatternActor
}

func NewPrinterPatternActor() PrinterPatternActor {
	return &printerPatternActor{}
}

var _ PrinterPatternActor = &printerPatternActor{}

type printerPatternActor struct {
}

func (*printerPatternActor) MaybeAct(line string) {
	log.Printf("Printer prints: %s\n", line)
}

// BlockWaiterPatternActor - A pattern actor that prints all output to the screen

type BlockWaiterPatternActor interface {
	PatternActor
	Wait() // Blocks until the actor has finished its work
}

func NewBlockWaiterPatternActor(blocks int, verbose bool) BlockWaiterPatternActor {
	wg := new(sync.WaitGroup)
	wg.Add(1) // IMPROVE: Consider making this counter equal to `blocks`
	return &blockWaiter{
		wg:              wg,
		remainingBlocks: blocks,
		done:            false,
		verbose:         verbose,
	}
}

var _ BlockWaiterPatternActor = &blockWaiter{}

type blockWaiter struct {
	wg              *sync.WaitGroup
	remainingBlocks int
	done            bool
	verbose         bool
}

func (b *blockWaiter) MaybeAct(line string) {
	if b.wg != nil && !b.done {
		if strings.Contains(line, "Executed block") {
			b.remainingBlocks--
			if b.verbose {
				log.Printf("Block elapsed; remaining: %d\n", b.remainingBlocks)
			}
		}
		if b.remainingBlocks <= 0 {
			b.wg.Done()
			b.done = true
		}
	}
}

func (b *blockWaiter) Wait() {
	b.wg.Wait()
}
