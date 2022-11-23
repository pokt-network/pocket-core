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
var pocketNetwork launcher.Network

func theUserHasAPocketExecutable() error {
	pocketExecutablePath = os.Getenv("POCKET_EXE")
	// TODO: try to get default pocket client before erroring out
	if pocketExecutablePath == "" {
		return errors.New("environment variable POCKET_EXE not set, can't run any tests without an executable.")
	} else {
		return nil
	}
}

func theUserHasAPocketClient() error {
	err := theUserHasAPocketExecutable()
	if err != nil {
		return err
	}
	pocketClient = launcher.NewPocketClient(pocketExecutablePath)
	return nil
}

func theUserRunsTheCommand(cmd string) (err error) {
	if pocketClient == nil {
		return errors.New("pocket executable not available")
	}
	cmdr, err = pocketClient.RunCommand(cmd)
	return
}

func theUserShouldBeAbleToSeeStandardOutputContaining(needle string) error {
	if cmdr == nil {
		return errors.New("no command result to analyze")
	}
	if !strings.Contains(cmdr.Stdout, needle) {
		errorMessage := fmt.Sprintf("standard output does not contain sought text. Sought %s.\n Stdout contained:\n%s\n Stderr contained:\n%s\n", needle, cmdr.Stdout, cmdr.Stderr)
		return errors.New(errorMessage)
	}
	return nil
}

func pocketClientShouldHaveExitedWithoutErrors() (err error) {
	if cmdr == nil {
		err = errors.New("no command result to analyze")
	} else if cmdr.Err != nil {
		err = errors.New(fmt.Sprintf("expected no errors. Got %v\n", cmdr.Err))
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
	_ = pocketNetwork.Nodes[0].PocketServer.RegisterPatternActor(&launcher.PrinterPatternActor{}, launcher.StdOut)
	time.Sleep(time.Second * 2)
	return nil
}
func theUserRunsTheCommandAgainstValidator(command string, validatorIdx int) error {
	if len(pocketNetwork.Nodes) <= validatorIdx {
		return errors.New(fmt.Sprintf("tried to use validator index %d, only %d available.", validatorIdx, len(pocketNetwork.Nodes)))
	}
	datadir := pocketNetwork.Nodes[validatorIdx].DataDir
	_ = theUserHasAPocketClient()
	commands := strings.Split(command, " ")
	for i := 0; i < len(commands); i++ {
		if commands[i] == "{{address}}" {
			commands[i] = pocketNetwork.Nodes[validatorIdx].Address
		}
	}
	commands = append(commands, "--datadir="+datadir)
	var err error = nil
	cmdr, err = pocketClient.RunCommand(commands...)
	return err
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

	// this step checks the status with which the client exited; should only be called after a `the user runs the command` step
	// example: `And pocket cliekt should have exited without errors`
	// this step is crucial for any command we expect to be run by an automated tool, and good to have elsewhere
	ctx.Step(`^the pocket client should have exited without errors$`, pocketClientShouldHaveExitedWithoutErrors)

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
	ctx.Step(`^the user runs the command "([^"]*)" against validator (\d+)$`, theUserRunsTheCommandAgainstValidator)
}
