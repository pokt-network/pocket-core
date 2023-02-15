Feature: Accounts Namespace

  The scenarios herein presented exercise all the actions exposed via the CLI for the accounts namespace.
  The accounts namespace of the CLI allows for the execution of tasks related to accounts and signatures:
  accounts creation, deletion, construction of multi-signature accounts, import and export of accounts, among others.

  Several of these scenarios use either a flag ("--pwd") to specify the password or they read the password with the readline
  facility without echoing to screen. Because communicating with the process means we'd use a pipe, reading the password
  from it would fail. Thus, we always use the flag instead.

  Scenario Outline: Help Needed for Accounts Namespace
	Given the user has a pocket client
	When the user runs the command <command>
	Then the user should be able to see standard output containing "Usage:"
	And the pocket client should have exited without error
	Examples:
	  | command           |
	  | "accounts"        |
	  | "accounts --help" |

  Scenario: Account Creation
	Given the user has a pocket client
	When the user runs the command "accounts create --pwd=test"
	Then the user should be able to see standard output containing "Account generated successfully"
	And the pocket client should have exited without error

  Scenario Outline: Using the Wrong Invocation for Accounts Namespace Commands
	Given the user has a pocket client
	When the user runs the command <mistaken_command>
	Then the user should be able to see standard error containing "Error: "
	And the user should be able to see standard error containing "Usage:"
	And the user should be able to see standard error containing <command_template_start>
	And the pocket client should have exited with error
	Examples:
	  | mistaken_command                                         | command_template_start       |
	  | "accounts delete"                                        | "accounts delete"            |
	  | "accounts export"                                        | "accounts export"            |
	  | "accounts export-raw"                                    | "accounts export-raw"        |
	  | "accounts import-armored"                                | "accounts import-armored"    |
	  | "accounts import-raw"                                    | "accounts import-raw"        |
	  | "accounts send-raw-tx"                                   | "accounts send-raw-tx"       |
	  | "accounts send-tx"                                       | "accounts send-tx"           |
	  | "accounts show"                                          | "accounts show"              |
	  | "accounts sign"                                          | "accounts sign"              |
	  | "accounts sign 7ab712998671b09e1a266ce6901000acb657833b" | "accounts sign"              |
	  | "accounts sign message"                                  | "accounts sign"              |
	  | "accounts update-passphrase"                             | "accounts update-passphrase" |

  Scenario Outline: Running Commands With Correct Syntax and Bad Data
	Given the user has a pocket client
	When the user runs the command <bad_command>
	Then the user should be able to see standard output containing <message_contents>
	And the pocket client should have exited without error
	Examples:
	  | bad_command                                                                   | message_contents               |
	  | "accounts export thisisdefinitelyawrongaccount --pwd-decrypt=badpassword"     | "Address Error"                |
	  | "accounts export-raw thisisdefinitelyawrongaccount --pwd-decrypt=badpassword" | "Address Error"                |
	  | "accounts sign address testing"                                               | "invalid byte"                 |
	  | "accounts create-multi-public thisisdefinitelyawrongaccount"                  | "error in public key creation" |

  Scenario: Account Sign With Non-Existing Account
	Given the user has a pocket client
	When the user runs the command "accounts sign 7ab712998671b09e1a266ce6901000acb657833b 'example_message'"
	Then the user should be able to see standard output containing "invalid byte"
	And the pocket client should have exited without error

  Scenario: Listing Existing Accounts
	Given the user has a pocket client
	When the user runs the command "accounts create --pwd=test"
	And the user runs the command "accounts list"
	Then the user should be able to see standard output containing "(0)"
	And the pocket client should have exited without error

  Scenario: create multi signature account
	Given the user has a pocket client
	When the user runs the command "accounts create-multi-public 883cc39e7f73259b4d5cb601a3251911373e6c10221e5f3b81c321caf5d16403,6047cca57f58f55bbe1e0c829c09d513986344d0081e832af18ad99517fc5c99"
	Then the user should be able to see standard output containing "Sucessfully generated Multisig Public Key"
	And the pocket client should have exited without error
