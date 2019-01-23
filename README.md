<div align="center">
  <a href="https://www.pokt.network">
    <img src="https://pokt.network/wp-content/uploads/2018/12/Logo-488x228-px.png" alt="drawing" width="340"/>
  </a>
</div>
<h1 align="center">Pocket Core</h1>
<h6 align="center">Official golang implementation of the Pocket Network Protocol.</h6>
<div align="center">
  <img href="https://godoc.org/github.com/pokt-network/pocket-core" src="https://img.shields.io/badge/godoc-reference-blue.svg"/>
  <img href="https://goreportcard.com/report/github.com/pokt-network/pocket-core" src="https://goreportcard.com/badge/github.com/pokt-network/pocket-core"/>
  <img href="https://golang.org" src="https://img.shields.io/badge/golang-v1.11-red.svg"/>
  <img href="https://github.com/tools/godep" src="https://img.shields.io/badge/godep-dependency-71a3d9.svg"/>
</div>

<h1 align="center"> Overview</h1>
<div align="center">
    <img href="https://github.com/pokt-network/pocket-core/releases" src="https://img.shields.io/github/release-pre/pokt-network/pocket-core.svg"/>
  <img href="https://github.com/pokt-network/pocket-core/pulse" src="https://img.shields.io/github/languages/code-size/pokt-network/pocket-core.svg"/>
  <img href="https://github.com/pokt-network/pocket-core/pulse" src="https://img.shields.io/github/contributors/pokt-network/pocket-core.svg"/>
  <img href="https://opensource.org/licenses/MIT" src="https://img.shields.io/badge/License-MIT-blue.svg"/>
    <br >
  <img href="https://github.com/pokt-network/pocket-core/pulse" src="https://img.shields.io/github/last-commit/pokt-network/pocket-core.svg"/>
  <img href="https://github.com/pokt-network/pocket-core/issues" src="https://img.shields.io/github/issues-pr/pokt-network/pocket-core.svg"/>
  <img alt="issues" href="https://github.com/pokt-network/pocket-core/issues?q=is%3Aissue+is%3Aclosed" src="https://img.shields.io/github/issues-closed/pokt-network/pocket-core.svg"/>
  <img href="https://github.com/pokt-network/pocket-core/commits/staging" src="https://img.shields.io/github/commit-activity/w/pokt-network/pocket-core.svg"/>
</div>

The Pocket Core application will allow anyone to spin up a Pocket Network full node, with options to enable/disable functionality and modules according to each deployment. For more information on the Pocket Network Protocol you can visit [pokt.network](https://pokt.network).

<h1 align="center">How to run it</h1>

To run the Pocket Core binary you can use the following flags alongside the `pocket-core` executable:

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

<h1 align="center">How to test</h1>
To run the Pocket Core unit tests, use the go testing tools and the `go test ./...` command within the tests directory

<h1 align="center">How to contribute</h1>
Pocket Core is an open source project, and as such we welcome any contribution from anyone on the internet. Please read our [Developer Setup Guide](https://github.com/pokt-network/pocket-core/wiki/Developer-Setup-Guide) on how get started.

Please fork, code and submit a Pull Request for the Pocket Core Team to review and merge. We ask that you please follow the guidelines below in order to submit your contributions for review:

<h3 align="center">High impact or architectural changes</h3>
Reach out to us on [Telegram](https://t.me/POKTnetwork) and start a discussion with the Pocket Core Team regarding your change before you start working. Communication is key for open source projects and asynchronous contributions.

For an active research forum, checkout and post on [our forum](https://research.pokt.network).

<h3 align="center">Coding style</h3>
<ul>
  <li>Code must adhere to the official Go formatting guidelines (i.e. uses [gofmt](https://golang.org/cmd/gofmt)).</li>

  <li>(Optional) Use [Editor Config](https://editorconfig.org) to help your Text Editor keep the same formatting used throughout the project.</li>

  <li>Code must be documented adhering to the official Go commentary guidelines.</li>

  <li>Pull requests need to be based on and opened against the `staging` branch.</.i>
</ul>
<h1 align="center"> How to build </h1>
run: `go build pokt-network/pocket-core/cmd/pocket_core/main.go`

<h1 align="center">Contact</h1>
<div align="center">
  <img href="https://twitter.com/poktnetwork" src="https://img.shields.io/twitter/url/http/shields.io.svg?style=social">
  <img href="https://t.me/POKTnetwork" src="https://img.shields.io/badge/Telegram-blue.svg">
  <img href="https://www.facebook.com/POKTnetwork" src="https://img.shields.io/badge/Facebook-red.svg">
  <img href="https://research.pokt.network" src="https://img.shields.io/discourse/https/research.pokt.network/posts.svg">
</div>
