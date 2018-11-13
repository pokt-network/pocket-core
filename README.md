# Pocket Core
Official implementation of the Pocket Network Protocol.

# Overview
The Pocket Core application will allow anyone to spin up a Pocket Network Full Node, with options to enable/disable functionality and modules according to each deployment. For more information on the Pocket Network Protocol you can visit [https://pokt.network](https://pokt.network).

# How to run it
To run the Pocket Core binary you can use the following flags alongside the `pocket-core` executable:

- `--datadir <absolute path>` to configure the data directory where the blockchain information will be stored. The default value is: `/path/to/data/dir`.
- `--clientrpc` to enable the Client RPC endpoints. The default value is `false`.
- `--clientrpcport <port number>` the port on which the Client RPC endpoints will run. The default value is `8545`.
- `--relayrpc` to enable the Relay RPC endpoints. The default value is `false`.
- `--relayrpcport <port number>` the port on which the Relay RPC endpoints will run. The default value is `8546`.

# How to contribute
Pocket Core is an open source project, and as such we welcome any contribution from anyone on the internet. Please read our [Developer Setup Guide](https://github.com/pokt-network/pocket-core/wiki/Developer-Setup-Guide) for a guide on how get started.

Please fork, code and submit a Pull Request for the Pocket Core Team to review and merge. We ask that you please follow the guidelines below in order to submit your contributions for review:

## High impact or architectural changes
Reach us out in our [Slack](https://www.pokt.network/slack-pokt) and start a discussion with the Pocket Core Team regarding your change before you start working, communication is key for open source projects and asynchronous contributions.

## Coding style
- Code must adhere to the official Go formatting guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt)).
- (Optional) Use [Editor Config](https://editorconfig.org) to help your Text Editor keep the same formatting used throughout the project.
- Code must be documented adhering to the official Go commentary guidelines.
- Pull requests need to be based on and opened against the `staging` branch.

# License
The pocket-core code is licensed under the MIT License, also included in the repository in the LICENSE file.
