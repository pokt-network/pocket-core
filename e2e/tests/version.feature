Feature: Pocket Version

  Scenario: User Checks Pocket Version
	Given the user has a pocket client
	When the user runs the command "version"
	Then the user should be able to see standard output containing "Version:"
	And the pocket client should have exited without error
