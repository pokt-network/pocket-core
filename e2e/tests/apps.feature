Feature: Pocket App Commands

  Scenario Outline: Help Needed for Apps Namespace
	Given the user has a pocket client
	When the user runs the command <cmd>
	Then the user should be able to see standard output containing "Usage:"
	And the pocket client should have exited without error

	Examples:
	  | cmd           |
	  | "apps"        |
	  | "apps --help" |

  Scenario Outline: Using the Wrong Invocation for Apps Namespace Commands
	Given the user has a pocket client
	When the user runs the command <mistaken_command>
	Then the user should be able to see standard error containing "Error: "
	And the user should be able to see standard error containing "Usage:"
	And the user should be able to see standard error containing <command_template_start>
	And the pocket client should have exited with error

	Examples:
	  | mistaken_command  | command_template_start |
	  | "apps create-aat" | "apps create-aat"      |
	  | "apps stake"      | "apps stake"           |
	  | "apps unstake"    | "apps unstake"         |



