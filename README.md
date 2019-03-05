<div align="center">
  <a href="https://www.pokt.network">
    <img src="https://pokt.network/wp-content/uploads/2018/12/Logo-488x228-px.png" alt="drawing" width="340"/>
  </a>
</div>
<h1 align="center">Pocket Core</h1>
<h6 align="center">Official golang implementation of the Pocket Network Protocol.</h6>
<div align="center">
  <a  href="https://godoc.org/github.com/pokt-network/pocket-core">
    <img src="https://img.shields.io/badge/godoc-reference-blue.svg"/>
  </a>
  <a  href="https://goreportcard.com/report/github.com/pokt-network/pocket-core">
    <img src="https://goreportcard.com/badge/github.com/pokt-network/pocket-core"/>
    </a>
  <a href="https://golang.org">
  <img  src="https://img.shields.io/badge/golang-v1.11-red.svg"/>
    </a>
  <a  href="https://github.com/tools/godep" >
    <img src="https://img.shields.io/badge/godep-dependency-71a3d9.svg"/>
  </a>
</div>

<h1 align="center"> Overview</h1>
  <div align="center">
    <a  href="https://github.com/pokt-network/pocket-core/releases">
      <img src="https://img.shields.io/github/release-pre/pokt-network/pocket-core.svg"/>
    </a>
    <a href="https://circleci.com/gh/pokt-network/pocket-core/tree/staging">
      <img src="https://circleci.com/gh/pokt-network/pocket-core/tree/staging.svg?style=svg"/>
    </a>
    <a  href="https://github.com/pokt-network/pocket-core/pulse">
      <img src="https://img.shields.io/github/contributors/pokt-network/pocket-core.svg"/>
    </a>
    <a href="https://opensource.org/licenses/MIT">
      <img src="https://img.shields.io/badge/License-MIT-blue.svg"/>
    </a>
    <br >
    <a href="https://github.com/pokt-network/pocket-core/pulse">
      <img src="https://img.shields.io/github/last-commit/pokt-network/pocket-core.svg"/>
    </a>
    <a href="https://github.com/pokt-network/pocket-core/pulls">
      <img src="https://img.shields.io/github/issues-pr/pokt-network/pocket-core.svg"/>
    </a>
    <a href="https://github.com/pokt-network/pocket-core/releases">
      <img src="https://img.shields.io/badge/platform-linux%20%7C%20windows%20%7C%20macos-pink.svg"/>
    </a>
    <a href="https://github.com/pokt-network/pocket-core/issues">
      <img src="https://img.shields.io/github/issues-closed/pokt-network/pocket-core.svg"/>
    </a>
</div>

The Pocket Core application will allow anyone to spin up a Pocket Network full node, with options to enable/disable functionality and modules according to each deployment. For more information on the Pocket Network Protocol you can visit <a href="https://pokt.network">pokt.network</a>.

<h1 align="center">How to run it</h1>

If you don't have the `pocket-core` binaries, scroll down to the "How to build" section for instructions on building from source.

The Pocket Core binary, `pocket-core`, accepts many arguments.

Here, we break them all down by category:

<h2> Arguments for service nodes</h2>

```
  -cfile string
    	specifies the filepath for chains.json 
	(default "[datadir]/chains.json")
  -clientrpc
    	whether or not to start the rpc server 
	(default true)
  -clientrpcport string
    	specified port to run client rpc 
	(default "8080")
  -datadirectory string
    	setup a custom location for the datadirectory 
	(default: `%APPDATA%\Pocket` for Windows, `~/.pocket` for Linux, `~/Library/Pocket` for Mac)
  -dfile string
    	specifies the filepath for developer_whitelist.json 
	(default "[datadir]/developer_whitelist.json")
  -discport string
    	specifies the client port of the centralized dispatcher 
	(default "8080")
  -disip string
    	specifies the address of the centralized dispatcher 
	(default "localhost")
  -disrport string
    	specifies the relay port of the centralized dispatcher 
	(default "8081")
  -gid string
    	set the selfNode.GID for pocket core mvp 
	(default "GID1")
  -pfile string
    	specifies the filepath for peers.json 
	(default "[datadir]/peers.json")
  -relayrpc
    	whether or not to start the rpc server 
	(default true)
  -relayrpcport string
    	specified port to run relay rpc 
	(default "8081")
  -sfile string
    	specifies the filepath for service_whitelist.json 
	(default "[datadir]/service_whitelist.json")
```

<h2>Arguments for dispatcher nodes</h2>

```
  -dbend string
    	specifies the database endpoint for the centralized dispatcher 
	(default "dynamodb.us-east-2.amazonaws.com")

  -dbtable string
    	specifies the database tablename for the centralized dispatcher 
	(default "dispatchPeers")
		
  -dismode int
    	specifies the mode by which the dispatcher is operating 
	(0) Normal, (1) Migrate, (2) Deprecated
	
  -dispatch
      	specifies if this node is operating as a dispatcher
```

<h1 align="center">How to build</h1>
If your environment is not set up, visit our <a href="https://github.com/pokt-network/pocket-core/wiki/Developer-Setup-Guide">Developer Setup Guide</a> to make sure you have everything you need to get the project up and running.

After your environment is set up, run: `go build pokt-network/pocket-core/cmd/pocket_core/main.go`

<h1 align="center">How to test</h1>

To run the Pocket Core unit tests, use the go testing tools and the `go test ./...` command within the tests directory

<h1 align="center">How to contribute</h1>
Pocket Core is an open source project, and as such we welcome any contribution from anyone on the internet. Please read our <a href="https://github.com/pokt-network/pocket-core/wiki/Developer-Setup-Guide">Developer Setup Guide</a> on how get started.

Please fork, code and submit a Pull Request for the Pocket Core Team to review and merge. We ask that you please follow the guidelines below in order to submit your contributions for review:

<h3 align="center">High impact or architectural changes</h3>
Reach out to us on <a href="https://t.me/POKTnetwork">Telegram</a> and start a discussion with the Pocket Core Team regarding your change before you start working. Communication is key for open source projects and asynchronous contributions.

For an active research forum, checkout and post on <a href="https://research.pokt.network">our forum</a>.

<h3 align="center">Coding style</h3>
<ul>
	<li>Code must adhere to the official Go formatting guidelines (i.e. uses <a href="https://golang.org/cmd/gofmt">gofmt</a>).</li>

  <li>(Optional) Use <a href="https://editorconfig.org">Editor Config</a> to help your Text Editor keep the same formatting used throughout the project.</li>

  <li>Code must be documented adhering to the official Go commentary guidelines.</li>

  <li>Pull requests need to be based on and opened against the `staging` branch.</.i>
</ul>

<h1 align="center">Contact</h1>
<div align="center">
  <a  href="https://twitter.com/poktnetwork" >
    <img src="https://img.shields.io/twitter/url/http/shields.io.svg?style=social">
  </a>
  <a href="https://t.me/POKTnetwork">
    <img src="https://img.shields.io/badge/Telegram-blue.svg">
  </a>
  <a href="https://www.facebook.com/POKTnetwork" >
  <img src="https://img.shields.io/badge/Facebook-red.svg">
  </a>
  <a href="https://research.pokt.network">
  <img src="https://img.shields.io/discourse/https/research.pokt.network/posts.svg">
  </a>
</div>
