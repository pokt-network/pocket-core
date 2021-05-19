---
description: >-
  This specification outlines all of the CLI commands to perform any function in
  Pocket Core. There's no protocol verification for these commands because they
  map closely to protocol functions.
---

# CLI

## Accessing Command List in Terminal

To quickly remind yourself of a command while in the terminal, simply enter `pocket help` or `pocket <namespace> help` to generate a list of all of the available commands and their associated flags.

## CLI Overview

### Namespaces

The CLI contains multiple namespaces listed below:

* \*\*\*\*[**Default**](default.md)**:** called when the namespace is blank
* [**Accounts**](accounts.md)**:** calls pertinent to accounts and their local storage
* [**Apps**](apps.md)**:** functions for app upkeep
* [**Nodes**](https://github.com/pokt-network/pocket-core/tree/7fe089d32a648f138a51d44ebf42ef38b57e196e/doc/specs/cli/nodes.md)**:** functions for node upkeep
* \*\*\*\*[**Query**](query.md)**:** queries to the world state
* [**Util**](util.md)**:** useful operations
* [**Gov**](gov.md)**:** functions for governance \(DAO\) transactions, only relevant to the DAOowner \(the account that has the permission to perform these transactions on behalf of the DAO\)

### CLI Functions Format

Each CLI function is constructed in the following format:

* **Binary Name:** The name of the binary for Pocket Core, for example: `pocket`
* **Global Options:** any number of global options, for example: `pocket --datadir /tmp/.pocket`
* **Namespace:** The namespace of the function, or blank for the default namespace: `accounts`
* **Function Name:** The name of the actual function to be called: `create`
* **Function Options:** Options that modify behaviour of the function `pocket query nodes --staking_status unstaking`
* **Arguments/Flags \(Optional\):** Space separated function arguments, e.g.: `pocket query nodes --staking_status unstaking <height>`

