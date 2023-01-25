package cli

import "github.com/cucumber/godog"

func InitializeScenario(ctx *godog.ScenarioContext) {
	// Ensures the presence of the pocket executable and sets things up to invoke commands.
	// Only necessary when not using the network.
	// Example: `Given the user has a pocket client`
	ctx.Step(`^the user has a pocket client$`, theUserHasAPocketClient)

	// Runs a command with the pocket executable.
	// Example: `Given the user runs the command "query --help"`
	ctx.Step(`^the user runs the command "([^"]*)"$`, theUserRunsTheCommand)

	// Analyzes the stdOut of a command
	// Should only be called after a `the user runs the command` step.
	// Example: `Then the user should be able to see standard output containing "Usage:"`
	ctx.Step(`^the user should be able to see standard output containing "([^"]*)"$`, theUserShouldBeAbleToSeeStandardOutputContaining)

	// Analyzes the stdErr of a command
	// Should only be called after a `the user runs the command` step.
	// Example: `Then the user should be able to see standard output containing "Usage:"`
	ctx.Step(`^the user should be able to see standard error containing "([^"]*)"$`, theUserShouldBeAbleToSeeStandardErrorContaining)

	// Checks the status with which the client exited
	// Should only be called after a `the user runs the command` step
	// Required for any command we expect to be run by an automated tool
	// Example: `And pocket client should have exited without error`
	ctx.Step(`^the pocket client should have exited without error$`, pocketClientShouldHaveExitedWithoutError)

	// Checks the status with which the client exited
	// Should only be called after a `the user runs the command` step
	// Required for any command we expect to be run by an automated tool
	// Example: `And pocket client should have exited with error`
	ctx.Step(`^the pocket client should have exited with error$`, pocketClientShouldHaveExitedWithError)

	// Spins up a particular network so commands can be run against its validators.
	// Example: `Given the user is running the network "single_node_network"` spins up the network described in `e2e/launcher/network_configs/single_node_network`
	ctx.Step(`^the user is running the network "([^"]*)"$`, theUserIsRunningTheNetwork)

	// Allows running a client command against a specific validator in the network.
	// Should only be called after a `the user is running the network` step.
	// NOTE: validators in the network are 0 indexed.
	// 		If you need to use the validator's address, represent it with {{}}, as in: "query account {{address}}"
	// 		If you need to have an invalid address, use validator -1
	// Example: "When the user runs the command "query accounts" against validator 0"
	ctx.Step(`^the user runs the command "([^"]*)" against validator (-?\d+)$`, theUserRunsTheCommandAgainstValidator)

	// Waits for the network to process a certain number of blocks before continuing.
	// Should only be called after a `the user is running the network` step.
	// Example: "Then the user waits for 1 block"
	ctx.Step(`^the user waits for (\d+) blocks?$`, theUserWaitsForBlocks)

	// Allows to see server output for test diagnosis and development purposes
	ctx.Step(`^verbose server$`, func() { verbose = true })
}
