Feature: Root Namespace

  Scenario: User Needs Help
	Given the user has a pocket client
	When the user runs the command "help"
	Then the user should be able to see standard output containing "Available Commands"
	And the pocket client should have exited without error
