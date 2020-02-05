<div align="center">
  <a href="https://www.pokt.network">
    <img src="https://pokt.network/wp-content/uploads/2019/11/pocket-logo.png" alt="Pocket Network logo" width="340"/>
  </a>
</div>

# Pocket Core

Official golang implementation of the Pocket Network Protocol.
<div>
  <a href="https://godoc.org/github.com/pokt-network/pocket-core"><img src="https://img.shields.io/badge/godoc-reference-blue.svg"/></a>
  <a href="https://goreportcard.com/report/github.com/pokt-network/pocket-core"><img src="https://goreportcard.com/badge/github.com/pokt-network/pocket-core"/></a>
  <a href="https://golang.org"><img  src="https://img.shields.io/badge/golang-v1.11-red.svg"/></a>
  <a href="https://github.com/tools/godep" ><img src="https://img.shields.io/badge/godep-dependency-71a3d9.svg"/></a>
</div>

## Overview
<div>
    <a  href="https://github.com/pokt-network/pocket-core/releases"><img src="https://img.shields.io/github/release-pre/pokt-network/pocket-core.svg"/></a>
    <a href="https://circleci.com/gh/pokt-network/pocket-core"><img src="https://circleci.com/gh/pokt-network/pocket-core.svg?style=svg"/></a>
    <a  href="https://github.com/pokt-network/pocket-core/pulse"><img src="https://img.shields.io/github/contributors/pokt-network/pocket-core.svg"/></a>
    <a href="https://opensource.org/licenses/MIT"><img src="https://img.shields.io/badge/License-MIT-blue.svg"/></a>
    <!--<a href="https://github.com/pokt-network/pocket-core/pulse"><img src="https://img.shields.io/github/last-commit/pokt-network/pocket-core.svg"/></a>-->
    <a href="https://github.com/pokt-network/pocket-core/pulls"><img src="https://img.shields.io/github/issues-pr/pokt-network/pocket-core.svg"/></a>
    <a href="https://github.com/pokt-network/pocket-core/releases"><img src="https://img.shields.io/badge/platform-linux%20%7C%20windows%20%7C%20macos-pink.svg"/></a>
    <!--<a href="https://github.com/pokt-network/pocket-core/issues"><img src="https://img.shields.io/github/issues-closed/pokt-network/pocket-core.svg"/></a>-->
</div>

The Pocket Core application will allow anyone to spin up a Pocket Network full node, with options to enable/disable functionality and modules according to each deployment. For more information on the Pocket Network Protocol you can visit [pokt.network](https://pokt.network).

## Getting Started

### Example usage

To run the Pocket Core binary you can use the following flags alongside the `pocket-core` executable:
````
  -cfile string
    	specifies the filepath for chains.json (default "chains.json")
  -clientrpc
    	whether or not to start the rpc server (default true)
  -clientrpcport string
    	specified port to run client rpc (default "8080")
  -datadirectory string
    	setup a custom location for the datadirectory
        	(default: `%APPDATA%\Pocket` for Windows, `~/.pocket` for Linux, `~/Library/Pocket` for Mac)
  -pfile string
    	specifies the filepath for peers.json (default "peers.json")
  -relayrpc
    	whether or not to start the rpc server (default true)
  -relayrpcport string
    	specified port to run relay rpc (default "8081")
````

### Installation

Clone the repository and run `go build pokt-network/pocket-core/cmd/pocket_core/main.go`

## Documentation

[Visit our developer portal](https://pocket-network.readme.io) for tutorials and technical documentation for the Pocket Network.

## Running the tests

To run the Pocket Core unit tests, use the go testing tools and the `go test ./...` command within the tests directory

## Contributing

Please read [CONTRIBUTING.md](https://github.com/pokt-network/pocket-core/blob/master/README.md) for details on contributions and the process of submitting pull requests.

## Support & Contact

<div>
  <a  href="https://twitter.com/poktnetwork" ><img src="https://img.shields.io/twitter/url/http/shields.io.svg?style=social"></a>
  <a href="https://t.me/POKTnetwork"><img src="https://img.shields.io/badge/Telegram-blue.svg"></a>
  <a href="https://www.facebook.com/POKTnetwork" ><img src="https://img.shields.io/badge/Facebook-red.svg"></a>
  <a href="https://research.pokt.network"><img src="https://img.shields.io/discourse/https/research.pokt.network/posts.svg"></a>
</div>

## License

This project is licensed under the MIT License; see the [LICENSE.md](LICENSE.md) file for details
