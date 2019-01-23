<div id="header" style="text-align:center">
  <img src="https://pokt.network/wp-content/uploads/2018/12/Logo-488x228-px.png" alt="drawing" width="340"/>
  <h1>Pocket Core</h1>
  <h6>Official golang implementation of the Pocket Network Protocol.</h6>
  <img alt="undefined" href="https://godoc.org/github.com/pokt-network/pocket-core" src="https://img.shields.io/badge/godoc-reference-blue.svg">
  <img alt="undefined" href="https://goreport.com" src="https://goreportcard.com/badge/github.com/pokt-network/pocket-core">
  <img alt="undefined" src="https://img.shields.io/badge/golang-v1.11-red.svg">
  <img alt="undefined" src="https://img.shields.io/badge/godep-dependency-71a3d9.svg">

  <h1> Overview</h1>
  <img alt="undefined" src="https://img.shields.io/github/release-pre/pokt-network/pocket-core.svg">
  <img alt="undefined" src="https://img.shields.io/github/languages/code-size/pokt-network/pocket-core.svg">
  <img alt="undefined" src="https://img.shields.io/github/contributors/pokt-network/pocket-core.svg">
  <img alt="undefined" src="https://img.shields.io/badge/License-MIT-blue.svg">

  <img alt="undefined" src="https://img.shields.io/github/last-commit/pokt-network/pocket-core.svg">
  <img alt="undefined" src="https://img.shields.io/github/issues-pr/pokt-network/pocket-core.svg">
  <img alt="issues" src="https://img.shields.io/github/issues-closed/pokt-network/pocket-core.svg">
  <img alt="undefined" src="https://img.shields.io/github/commit-activity/w/pokt-network/pocket-core.svg">

The Pocket Core application will allow anyone to spin up a Pocket Network full node, with options to enable/disable functionality and modules according to each deployment. For more information on the Pocket Network Protocol you can visit [pokt.network](https://pokt.network).

<h1>How to run it</h1>

To run the Pocket Core binary you can use the following flags alongside the `pocket-core` executable:
</div>

    -clientrpc
      	whether or not to start the rpc server
    -clientrpcport string
      	specified port to run client rpc (default "8080")
    -datadir string
      	setup the data director for the DB and keystore 
      	(default: `%APPDATA%\Pocket` for Windows, `~/.pocket` for Linux, `~/Library/Pocket` for Mac)
    -hostedchains string
      	specifies the filepath for hosted chains (default "datadir/chains.json")
    -peerFile string
      	specifies the filepath for peers.json (default "datadir/peers.json")
    -relayrpc
      	whether or not to start the rpc server
    -relayrpcport string
      	specified port to run relay rpc (default "8081")

<div style="text-align:center">
  <h1>How to test</h1>
  To run the Pocket Core unit tests, use the go testing tools and the `go test ./...` command within the tests directory

  # How to contribute
  Pocket Core is an open source project, and as such we welcome any contribution from anyone on the internet. Please read our [Developer Setup Guide](https://github.com/pokt-network/pocket-core/wiki/Developer-Setup-Guide) on how get started.

  Please fork, code and submit a Pull Request for the Pocket Core Team to review and merge. We ask that you please follow the guidelines below in order to submit your contributions for review:

  ## High impact or architectural changes
  Reach out to us on [Telegram](https://t.me/POKTnetwork) and start a discussion with the Pocket Core Team regarding your change before you start working. Communication is key for open source projects and asynchronous contributions.

  For an active research forum, checkout and post on [our forum](https://research.pokt.network).

  ## Coding style
  - Code must adhere to the official Go formatting guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt)).
  - (Optional) Use [Editor Config](https://editorconfig.org) to help your Text Editor keep the same formatting used throughout the project.
  - Code must be documented adhering to the official Go commentary guidelines.
  - Pull requests need to be based on and opened against the `staging` branch.

  # How to build
  `go build pokt-network/pocket-core/cmd/pocket_core/main.go`

  # Contact
  <img alt="undefined" href="https://twitter.com/poktnetwork" src="https://img.shields.io/twitter/url/http/shields.io.svg?style=social">
  <img alt="undefined" href="https://t.me/POKTnetwork" src="https://img.shields.io/badge/Telegram-blue.svg">
  <img alt="undefined" href="https://www.facebook.com/POKTnetwork" src="https://img.shields.io/badge/Facebook-red.svg">
  <img alt="undefined" src="https://img.shields.io/discourse/https/research.pokt.network/posts.svg">
</div>
