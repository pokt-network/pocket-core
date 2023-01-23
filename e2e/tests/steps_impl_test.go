package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/pokt-network/pocket-core/e2e/launcher"
)

// IMPROVE: Is there a better way to maintain this state rather than global variables?
// IMPROVE: The use of these global vars without checking there values can lead to a panic; the flow of calls is implicit
var (
	pocketExecutablePath                       = ""
	pocketClient         launcher.PocketClient = nil
	cmdResult            *launcher.CommandResult
	pocketNetwork        *launcher.Network
	verbose              bool
)

// Internal helper for `theUserHasAPocketClient`
func verifyPocketExecutable() error {
	pocketExecutablePath = os.Getenv("POCKET_EXE")
	if pocketExecutablePath == "" {
		return errors.New("`POCKET_EXE` not set. Please, set it to the output of `go build -o pocket app/cmd/pocket_core/main.go`")
	} else if _, err := os.Stat(pocketExecutablePath); err != nil {
		return errors.New(fmt.Sprintf("No executable found at ${POCKET_EXE}: %s", pocketExecutablePath))
	}
	return nil
}

// Step: `the user has a pocket client`
func theUserHasAPocketClient() error {
	if err := verifyPocketExecutable(); err != nil {
		return err
	}

	pocketClient = launcher.NewPocketClient(pocketExecutablePath, verbose)
	return nil
}

// Step: `the user runs the command`
func theUserRunsTheCommand(cmd string) (err error) {
	if pocketClient == nil {
		return errors.New("pocket executable not available. Scenario might lack a `the user has a pocket client` or `the user is running the network` step")
	}

	commands := strings.Split(cmd, " ")
	cmdResult, err = pocketClient.RunCommand(commands...)
	return
}

// Step: `the user should be able to see standard output containing`
func theUserShouldBeAbleToSeeStandardOutputContaining(expected string) error {
	if cmdResult == nil {
		return errors.New("no command result to analyze. The scenario might lack a `the user runs the command` or a `the user runs the command against validator` step")
	}

	if !strings.Contains(cmdResult.Stdout, expected) {
		errorMessage := fmt.Sprintf("standard output does not contain expected text. Expected: \"%s\".\n Actual: \n%s\n. Standard error contained:\n%s\n", expected, cmdResult.Stdout, cmdResult.Stderr)
		return errors.New(errorMessage)
	}

	return nil
}

// Step: `the user should be able to see standard error containing`
func theUserShouldBeAbleToSeeStandardErrorContaining(expected string) error {
	if cmdResult == nil {
		return errors.New("no command result to analyze. The scenario might lack a `the user runs the command` or a `the user runs the command against validator` step")
	}

	if !strings.Contains(cmdResult.Stderr, expected) {
		errorMessage := fmt.Sprintf("standard output does not contain expected text. Expected: \"%s\".\n Actual:\n%s\n Standard error contained:\n%s\n", expected, cmdResult.Stdout, cmdResult.Stderr)
		return errors.New(errorMessage)
	}

	return nil
}

// Step: `the pocket client should have exited without error`
func pocketClientShouldHaveExitedWithoutError() (err error) {
	if cmdResult == nil {
		err = errors.New("no command result to analyze. The scenario might lack a `the user runs the command` or a `the user runs the command against validator` step")
	} else if cmdResult.Err != nil {
		err = errors.New(fmt.Sprintf("expected no error, but pocket exited with the following error: %v\n", cmdResult.Err))
	}
	return
}

// Step: `the pocket client should have exited with error`
func pocketClientShouldHaveExitedWithError() (err error) {
	if cmdResult == nil {
		err = errors.New("no command result to analyze. The scenario might lack a `the user runs the command` or a `the user runs the command against validator` step")
	} else if cmdResult.Err == nil {
		err = errors.New("expected an error, but exited without error")
	}
	return
}

// Step: `the user is running the network "([^"]*)`
func theUserIsRunningTheNetwork(netName string) (err error) {
	if err := verifyPocketExecutable(); err != nil {
		return err
	}

	pocketNetwork, err = launcher.LaunchNetwork(netName, pocketExecutablePath)
	if err != nil {
		return
	}

	// IMPROVE: Should the printer pattern actor be registered by default?
	if verbose {
		printerActor := launcher.NewPrinterPatternActor()
		if err = pocketNetwork.Nodes[0].PocketServer.RegisterPatternActor(printerActor, launcher.StdOut); err != nil {
			return
		}
	}

	time.Sleep(time.Second * 2) // DOCUMENT: Why is this necessary?
	return
}

// Step: `the user runs the command "([^"]*)" against validator (-?\d+)`
func theUserRunsTheCommandAgainstValidator(command string, validatorIdx int) (err error) {
	if len(pocketNetwork.Nodes) <= validatorIdx {
		return errors.New(fmt.Sprintf("tried to use validator index %d, only index -1 through %d available", validatorIdx, len(pocketNetwork.Nodes)-1))
	}

	invalidValidator := validatorIdx < 0
	var datadir string
	if !invalidValidator {
		datadir = pocketNetwork.Nodes[validatorIdx].DataDir
	} else {
		datadir = pocketNetwork.Nodes[0].DataDir
	}

	_ = theUserHasAPocketClient()
	commands := strings.Split(command, " ")
	for i := 0; i < len(commands); i++ {
		if commands[i] == "{{address}}" {
			commands[i] = "this_is_an_obviously_invalid_address"
			if !invalidValidator {
				commands[i] = pocketNetwork.Nodes[validatorIdx].Address
			}
		}
	}
	commands = append(commands, "--datadir="+datadir)
	cmdResult, err = pocketClient.RunCommand(commands...)
	return
}

// Step: `the user waits for (\d+) blocks`
func theUserWaitsForBlocks(count int) error {
	blockWriterActor := launcher.NewBlockWaiterPatternActor(count, verbose)
	if err := pocketNetwork.Nodes[0].PocketServer.RegisterPatternActor(blockWriterActor, launcher.StdOut); err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Starting to wait for blocks; %d remaining.\n", count)
	}

	blockWriterActor.Wait()
	return nil
}
