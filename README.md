<div align="center">
  <a href="https://www.pokt.network">
    <img src="pocket-core.png" alt="Pocket Network logo" width="340"/>
  </a>
</div>

# Pocket Core <!-- omit in toc -->

Official golang implementation of the Pocket Network Protocol.

<div>
  <a href="https://godoc.org/github.com/pokt-network/pocket-core"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"/></a>
  <a href="https://goreportcard.com/report/github.com/pokt-network/pocket-core"><img src="https://goreportcard.com/badge/github.com/pokt-network/pocket-core"/></a>
  <a href="https://golang.org"><img  src="https://img.shields.io/badge/golang-v1.21-red.svg"/></a>
  <a href="https://github.com/tools/godep" ><img src="https://img.shields.io/badge/godep-dependency-71a3d9.svg"/></a>
</div>

## Overview

<div>
    <a href="https://discord.gg/pokt"><img src="https://img.shields.io/discord/553741558869131266"></a>
    <a  href="https://github.com/pokt-network/pocket-core/releases"><img src="https://img.shields.io/github/release-pre/pokt-network/pocket-core.svg"/></a>
    <a href="https://circleci.com/gh/pokt-network/pocket-core"><img src="https://circleci.com/gh/pokt-network/pocket-core.svg?style=svg"/></a>
    <a  href="https://github.com/pokt-network/pocket-core/pulse"><img src="https://img.shields.io/github/contributors/pokt-network/pocket-core.svg"/></a>
    <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-blue.svg"/></a>
    <!--<a href="https://github.com/pokt-network/pocket-core/pulse"><img src="https://img.shields.io/github/last-commit/pokt-network/pocket-core.svg"/></a>-->
    <a href="https://github.com/pokt-network/pocket-core/pulls"><img src="https://img.shields.io/github/issues-pr/pokt-network/pocket-core.svg"/></a>
    <a href="https://github.com/pokt-network/pocket-core/releases"><img src="https://img.shields.io/badge/platform-linux%20%7C%20macos-pink.svg"/></a>
    <!--<a href="https://github.com/pokt-network/pocket-core/issues"><img src="https://img.shields.io/github/issues-closed/pokt-network/pocket-core.svg"/></a>-->
</div>

The Pocket Core application will allow anyone to spin up a Pocket Network full node, with options to enable/disable functionality and modules according to each deployment. For more information on Pocket Network, visit [pokt.network](https://pokt.network).

## Table of Contents <!-- omit in toc -->

- [Overview](#overview)
- [Installation](#installation)
- [Usage](#usage)
- [Documentation](#documentation)
- [Portal](#portal)
- [Database Snapshots](#database-snapshots)
- [Pocket Pruner](#pocket-pruner)
- [Accessing TestNet](#accessing-testnet)
- [Charts \& Analytics](#charts--analytics)
- [Running the tests](#running-the-tests)
- [Contributing](#contributing)
- [Seeds (MainNet \& TestNet)](#seeds-mainnet--testnet)
- [Docker Image](#docker-image)
- [Chain Halt Rollback Recovery Guide](#chain-halt-rollback-recovery-guide)
- [Support \& Contact](#support--contact)
  - [GPokT](#gpokt)
- [License](#license)

## Installation

```bash
# Build local binary
git clone git@github.com:pokt-network/pocket && \
cd pocket && \
go build app/cmd/pocket_core/pocket.go

# Assign local binary and add to your `PATH`  if you'd like to use it without direct reference to the binary.
export POKT=$(pwd)/main
```

TIP: You can find alternative ways of installing `pocket` (e.g. homebrew) via the instructions [here](doc/guides/quickstart.md).

## Usage

To run the Pocket Core binary you can use the following flags alongside the `pocket` executable:

```bash
Usage:
  pocket [command]

Available Commands:
  accounts    account management
  apps        application management
  completion  Generate the autocompletion script for the specified shell
  gov         governance management
  help        Help about any command
  nodes       node management
  query       query the blockchain
  reset       Reset pocket-core
  start       starts pocket-core daemon
  stop        Stop pocket-core
  util        utility functions
  version     Get current version

Flags:
      --datadir string            data directory (default is $HOME/.pocket/
  -h, --help                      help for pocket
      --node string               takes a remote endpoint in the form <protocol>://<host>:<port>
      --persistent_peers string   a comma separated list of PeerURLs: '<ID>@<IP>:<PORT>,<ID2>@<IP2>:<PORT>...<IDn>@<IPn>:<PORT>'
      --remoteCLIURL string       takes a remote endpoint in the form of <protocol>://<host> (uses RPC Port)
      --seeds string              a comma separated list of PeerURLs: '<ID>@<IP>:<PORT>,<ID2>@<IP2>:<PORT>...<IDn>@<IPn>:<PORT>'

Use "pocket [command] --help" for more information about a command.
```

For more detailed command information, see the [usage section](doc/specs/cli/).

## Documentation

[Visit our user documentation](https://docs.pokt.network) for tutorials and technical information on the Pocket Network.

## Portal

The Portal to the Pocket Network is provided by [Pocket Network Inc](https://portal.pokt.network/).

## Database Snapshots

Snapshots are provided by [Liquify LTD](https://www.liquify.io/) details on how to access the snapshots can be found in [snapshot.md](doc/guides/snapshot.md)

## Pocket Pruner

An offline pruning tool is provided by [C0D3R](https://c0d3r.org/). The tool is available in [their GitHub repository](https://github.com/msmania/pocket-pruner/).

## Accessing TestNet

TestNet information can be found at [testnet.md](doc/guides/testnet.md) and is maintained by the [nodefleet.org](https://nodefleet.org/) team.

## Charts & Analytics

Key charts & analytics are provided by [POKTScan](https://poktscan.com/) and [C0D3R](https://c0d3r.org).

## Running the tests

To run the Pocket Core unit tests, `go test -short -v -p 1 ./...`

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on contributions and the process of submitting pull requests.

## Seeds (MainNet & TestNet)

Seeds are maintained by [NodeFleet](https://nodefleet.org/).

You can find all the details at [pokt-network/pocket-seeds](https://github.com/pokt-network/pocket-seeds).

## Docker Image

[GitHub Packages](https://github.com/features/packages) is used to maintain
docker images via [this workflow](https://github.com/pokt-network/pocket-core/blob/staging/.github/workflows/build-images.yaml).

The latest images can be found [here](https://github.com/pokt-network/pocket-core/pkgs/container/pocket-v0).

The latest image can be pulled like so:

```bash
docker pull ghcr.io/pokt-network/pocket-v0:latest
```

## Chain Halt Rollback Recovery Guide

See the rollback guide [here](doc/guides/rollback.md) put together by @msmania based on a real world scenario.

## Support & Contact

<div>
  <a  href="https://twitter.com/poktnetwork" ><img src="https://img.shields.io/twitter/url/http/shields.io.svg?style=social"></a>
  <a href="https://t.me/POKTnetwork"><img src="https://img.shields.io/badge/Telegram-blue.svg"></a>
  <a href="https://research.pokt.network"><img src="https://img.shields.io/discourse/https/research.pokt.network/posts.svg"></a>
</div>

### GPokT

You can also use our chatbot, [GPokT](https://gpoktn.streamlit.app), to ask questions about Pocket Network. As of updating this documentation, please note that it may require you to provide your own LLM API token. If the deployed version of GPokT is down, you can deploy your own version by following the instructions [here](https://github.com/pokt-network/gpokt).

## License

This project is licensed under the MIT License; see the [LICENSE.md](LICENSE.md) file for details
