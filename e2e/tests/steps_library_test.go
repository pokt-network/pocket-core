package cli

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/cucumber/godog"
	"github.com/pokt-network/pocket-core/e2e/launcher"
)

var pocketExecutablePath = ""
var pocketClient launcher.PocketClient = nil
var cmdr *launcher.CommandResult
var pocketNetwork *launcher.Network
var verbose bool

func theUserHasAPocketExecutable() error {
	pocketExecutablePath = os.Getenv("POCKET_EXE")
	if pocketExecutablePath == "" {
		return errors.New("environment variable POCKET_EXE not set. Please, set it to the path of the tested executable")
	} else if _, err := os.Stat(pocketExecutablePath); err != nil {
		return errors.New(fmt.Sprintf("pocket executable doesn't seem to be at location specified by POCKET_EXE (%s)", pocketExecutablePath))
	} else {
		return nil
	}
}

func theUserHasAPocketClient() error {
	err := theUserHasAPocketExecutable()
	if err != nil {
		return err
	}

	pocketClient = launcher.NewPocketClient(pocketExecutablePath, verbose)
	return nil
}

func theUserRunsTheCommand(cmd string) (err error) {
	if pocketClient == nil {
		return errors.New("pocket executable not available. The scenario might lack a `the user has a pocket client` or `the user is running the network` step")
	}

	commands := strings.Split(cmd, " ")
	cmdr, err = pocketClient.RunCommand(commands...)
	return
}

func theUserShouldBeAbleToSeeStandardOutputContaining(needle string) error {
	if cmdr == nil {
		return errors.New("no command result to analyze. The scenario might lack a `the user runs the command` or a `the user runs the command against validator` step")
	}

	if !strings.Contains(cmdr.Stdout, needle) {
		errorMessage := fmt.Sprintf("standard output does not contain sought text. Text sought: \"%s\".\n Standard output contained:\n%s\n Standard error contained:\n%s\n", needle, cmdr.Stdout, cmdr.Stderr)
		return errors.New(errorMessage)
	}

	return nil
}

func theUserShouldBeAbleToSeeStandardErrorContaining(needle string) error {
	if cmdr == nil {
		return errors.New("no command result to analyze. The scenario might lack a `the user runs the command` or a `the user runs the command against validator` step")
	}

	if !strings.Contains(cmdr.Stderr, needle) {
		errorMessage := fmt.Sprintf("standard output does not contain sought text. Text sought: \"%s\".\n Standard output contained:\n%s\n Standard error contained:\n%s\n", needle, cmdr.Stdout, cmdr.Stderr)
		return errors.New(errorMessage)
	}

	return nil
}

func pocketClientShouldHaveExitedWithoutError() (err error) {
	if cmdr == nil {
		err = errors.New("no command result to analyze. The scenario might lack a `the user runs the command` or a `the user runs the command against validator` step")
	} else if cmdr.Err != nil {
		err = errors.New(fmt.Sprintf("expected no error, but pocket exited with the following error: %v\n", cmdr.Err))
	}
	return
}

func pocketClientShouldHaveExitedWithError() (err error) {
	if cmdr == nil {
		err = errors.New("no command result to analyze. The scenario might lack a `the user runs the command` or a `the user runs the command against validator` step")
	} else if cmdr.Err == nil {
		err = errors.New("expected an error, but pocket exited without error")
	}
	return
}

func theUserIsRunningTheNetwork(netName string) error {
	err := theUserHasAPocketExecutable()
	if err != nil {
		return err
	}

	net, launchErr := launcher.LaunchNetwork(netName, pocketExecutablePath)
	if launchErr != nil {
		return launchErr
	}

	pocketNetwork = net
	if verbose {
		err = pocketNetwork.Nodes[0].PocketServer.RegisterPatternActor(&launcher.PrinterPatternAction{}, launcher.StdOut)
		if err != nil {
			return err
		}
	}

	time.Sleep(time.Second * 2)
	return nil
}
func theUserRunsTheCommandAgainstValidator(command string, validatorIdx int) error {
	if len(pocketNetwork.Nodes) <= validatorIdx {
		return errors.New(fmt.Sprintf("tried to use validator index %d, only index -1 through %d available", validatorIdx, len(pocketNetwork.Nodes)-1))
	}

	invalidValidator := validatorIdx < 0
	datadir := pocketNetwork.Nodes[0].DataDir
	if !invalidValidator {
		datadir = pocketNetwork.Nodes[validatorIdx].DataDir
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
	var err error = nil
	cmdr, err = pocketClient.RunCommand(commands...)
	return err
}

func theUserWaitsForBlocks(count int) error {
	bw := launcher.NewBlockWaiter(count, verbose)
	err := pocketNetwork.Nodes[0].PocketServer.RegisterPatternActor(bw, launcher.StdOut)
	if err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Starting to wait for blocks; %d remaining.\n", count)
	}
	bw.Wg.Wait()
	return nil
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	// this step ensures the presence of the pocket executable and sets things up to invoke commands
	// it is only necessary when not using the network.
	// example: `Given the user has a pocket client`
	ctx.Step(`^the user has a pocket client$`, theUserHasAPocketClient)

	// this step runs a command with the pocket executable.
	// example usage: `Given the user runs the command "query --help"`
	ctx.Step(`^the user runs the command "([^"]*)"$`, theUserRunsTheCommand)

	// this step analyzes the standard output of a command; it should only be called after a `the user runs the command` step.
	// example: `Then the user should be able to see standard output containing "Usage:"`
	ctx.Step(`^the user should be able to see standard output containing "([^"]*)"$`, theUserShouldBeAbleToSeeStandardOutputContaining)

	// this step analyzes the standard error of a command; it should only be called after a `the user runs the command` step.
	// example: `Then the user should be able to see standard output containing "Usage:"`
	ctx.Step(`^the user should be able to see standard error containing "([^"]*)"$`, theUserShouldBeAbleToSeeStandardErrorContaining)

	// this step checks the status with which the client exited; should only be called after a `the user runs the command` step
	// example: `And pocket client should have exited without error`
	// this step is crucial for any command we expect to be run by an automated tool
	ctx.Step(`^the pocket client should have exited without error$`, pocketClientShouldHaveExitedWithoutError)

	// this step checks the status with which the client exited; should only be called after a `the user runs the command` step
	// example: `And pocket client should have exited with error`
	// this step is crucial for any command we expect to be run by an automated tool
	ctx.Step(`^the pocket client should have exited with error$`, pocketClientShouldHaveExitedWithError)

	// this step spins up a particular network so commands can be run against its validators.
	// example: `Given the user is running the network "single_node_network"` spins up the network described in
	// `e2e/launcher/network_configs/single_node_network`
	ctx.Step(`^the user is running the network "([^"]*)"$`, theUserIsRunningTheNetwork)

	// this step allows running a client command against a specific validator in the network.
	// should only be called after a `the user is running the network` step.
	// example: "When the user runs the command "query accounts" against validator 0"
	// validators in the network are 0 indexed.
	// if you need to use the validator's address, represent it with {{}}, as in: "query account {{address}}"
	// if you need to have an invalid address, use validator -1
	ctx.Step(`^the user runs the command "([^"]*)" against validator (-?\d+)$`, theUserRunsTheCommandAgainstValidator)

	// this step allows to wait for the network to process a certain number of blocks before the tests go on.
	// should only be called after a `the user is running the network` step.
	// example: "Then the user waits for 1 block"
	ctx.Step(`^the user waits for (\d+) blocks?$`, theUserWaitsForBlocks)

	// this steps allows to see server output for test diagnosis and development purposes
	ctx.Step(`^verbose server$`, func() { verbose = true })
}
