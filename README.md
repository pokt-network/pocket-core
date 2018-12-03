# Pocket Core
Official implementation of the Pocket Network Protocol.

# Overview
The Pocket Core application will allow anyone to spin up a Pocket Network full node, with options to enable/disable functionality and modules according to each deployment. For more information on the Pocket Network Protocol you can visit [https://pokt.network](https://pokt.network).

# How to run it
To run the Pocket Core binary you can use the following flags alongside the `pocket-core` executable:

  `-bitcoin`
    	whether or not bitcoin is hosted
    	
 ` -btcrpcport string`
    	specified port to run bitcoin rpc (default "8333")
    	
 ` -clientrpc`
    	whether or not to start the rpc server
    	
 ` -clientrpcport string`
    	specified port to run client rpc (default "8080")
    	
 ` -datadir string`
    	setup the data director for the DB and keystore 
    	`%APPDATA%\Pocket` for Windows, 
    	`~/.pocket` for Linux, 
    	`~/Library/Pocket` for Mac
    	
 ` -ethereum`
    	whether or not ethereum is hosted
    	
 ` -ethrpcport string`
    	specified port to run ethereum rpc (default "8545")
    	
  `-manpeers`
    	specifies if peers are manually added
    	
  `-peerFile string`
    	specifies the filepath for peers.json (default "<DATADIR>/peers.json")
    	
  `-relayrpc`
    	whether or not to start the rpc server
    	
  `-relayrpcport string`
    	specified port to run relay rpc (default "8081")
# How to test
To run the Pocket Core unit tests, use the go testing tools and the `go test ./...` command within the tests directory

# How to contribute
Pocket Core is an open source project, and as such we welcome any contribution from anyone on the internet. Please read our [Developer Setup Guide](https://github.com/pokt-network/pocket-core/wiki/Developer-Setup-Guide) on how get started.

Please fork, code and submit a Pull Request for the Pocket Core Team to review and merge. We ask that you please follow the guidelines below in order to submit your contributions for review:

## High impact or architectural changes
Reach out to us on [Slack](https://www.pokt.network/slack-pokt) and start a discussion with the Pocket Core Team regarding your change before you start working. Communication is key for open source projects and asynchronous contributions.

## Coding style
- Code must adhere to the official Go formatting guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt)).
- (Optional) Use [Editor Config](https://editorconfig.org) to help your Text Editor keep the same formatting used throughout the project.
- Code must be documented adhering to the official Go commentary guidelines.
- Pull requests need to be based on and opened against the `staging` branch.

# How to build
`go build pokt-network/pocket-core/cmd/pocket_core/main.go`

### Where to find logs
`<datadir>/logs`

# License
The pocket-core code is licensed under the MIT License, also included in the repository in the LICENSE file.
