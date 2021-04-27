Pocket Core User Guide
========================
## Introduction
Pocket Core is a golang implementation of the Pocket Network Protocol. 

Go [here](https://forum.pokt.network/) for Pocket Network Protocol documentation

Pocket Core is software for Pocket Network 'node runners'. This software implements full node and validator capabilities. 

#### Contents
- Install
- Quickstart
- Config
- Operation
- Resources
- FAQ

#### Disclaimer
*PNI is not liable for any slashing or economic penalty that may occur*

## Install

### From Source
#### Prerequisite Installations
[go](https://golang.org/doc/install)

[go environment](https://golang.org/doc/gopath_code.html#Workspaces) GOPATH & GOBIN

[git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)
#### Create source code directory
`mkdir -p $GOPATH/src/github.com/pok-network && cd $GOPATH/src/github.com/pok-network`
#### Download the source code
`git clone https://github.com/pokt-network/pocket-core.git && cd pocket-core`
#### Checkout the [latest release](https://github.com/pokt-network/pocket-core/releases)
`git checkout tags/<release tag>`

Example: *git checkout tags/RC-0.6.2*

#### Build

`go build -o <destination directory> <source code directory>/...`

Example: *go build -o $GOPATH/bin/pocket $GOPATH/src/github.com/pok-network/pocket-core*

#### Test installation
```
$ pocket version
> RC-0.6.2
```
### From Homebrew
`brew tap pokt-network/pocket-core && brew install pokt-network/pocket-core/pocket`

#### Test installation
```
$ pocket version
> RC-0.6.2
```
### From Deployment Artifact
See [pokt-network/pocket-core-deployments](https://github.com/pokt-network/pocket-core-deployments)

## Quickstart
#### Prerequisite Knowledge
*This section does not cover the protocol specificiation, rather how to participate in the network as a node runner. For more information on the Pocket Network Protocol, read the wiki [here](https://forum.pokt.network/)*

A **Validator** is an infrastructure provider in Pocket Network

**Staking** a Validator locks up a certain **amount** of balance that can be burned as a security mechanism for bad acting

A **Relay Chain** is blockchain infrastructure **Validators** expose for application access *Ex: Ethereum, Bitcoin, Pocket Network* (identified by 4 hexadecimal characters. *Ex: 0001*)

Apps access **Relay Chains** through the **serviceURI**: the endpoint where **Validators** publicly expose the Pocket API *Ex: https://www.node1.mainnet.pokt.network*
#### Environment Setup
Hardware Requirements: 4 CPU’s (or vCPU’s) | 8 GB RAM | 100GB Disk

Reverse Proxy: For SSL termination and request management

Ports: Expose Pocket RPC (Default :8081) and P2P port (Default: 26656)

SSL Cert: Required for **Validator's serivceURI**

**Open Files Limit**: `ulimit -Sn 16384`

*NOTE*: **Open Files Limit** is very important for the operation of Pocket Core. See **Config** section for ulimit calculation
#### Create an account
An account is needed to participate at any level of the network.
```
pocket accounts create
> Enter Passphrase
> Account generated successfully:
> Address: <address>
```
#### Fund the account
To stake a Validator in Pocket Network, the account must have a balance above the **minimum stake**:

`15,000 POKT` or `15,000,000,000 uPOKT`

Send POKT with the following command:
```
pocket accounts send-tx <fromAddr> <toAddr> <uPOKT amount> mainnet 10000 "" true
```
#### Set the account as Validator
```
pocket accounts set-validator <address>
```
NOTE: Check with *pocket accounts get-validator*
#### Set [Relay Chains](https://forum.pokt.network/t/supportedblockchains/607)
```
pocket util generate-chains
> Enter the chain of the network identifier:
<Relay Chain ID> (Example: 0001)
> Enter the URL of the network identifier:
<Secure URL to Relay Chain>
Would you like to enter another network identifier? (y/n)
n
```
NOTE: *Can test with simulate relay flag and endpoint. See [RPC Specification](rpc_spec.yaml) for details*
#### Sync the blockchain
```
pocket start --seeds=<seeds> --mainnet
```
Example: *pocket start --seeds="64c91701ea98440bc3674fdb9a99311461cdfd6f@node1.mainnet.pokt.network:21656" --mainnet*

[Seeds](https://forum.pokt.network/t/list-of-seeds/647)

NOTE: *Ensure the node is all the way synced before proceeding to the next step*
#### Stake the Validator
Stake the account to participate in the Network as a **Validator**
```
pocket nodes stake <address> <amount> <relay_chains> <serviceURI> mainnet 10000 true
```
Example: *pocket nodes stake 3ee61299d5bbbd2974cddcc194d9b547c7629546 20000000000 ["0001", "0002"] https://pokt.rocks mainnet 10000 true*

**Important:** *Stake 'well over' the minimum stake to avoid force-unstake burning*
## Config
### Data Directory
Pocket Core files are located in a **Data Directory** Default: `$HOME/.pocket/`
### Configuration File
Pocket Core provides a configuration file found in `<datadir>/config/config.json`
### Pocket
  * **"data_dir"**: The data directory of Pocket Core (should be the same directory as Tendermint data dir)
  * **"genesis_file"**: The name of the genesis file
  * **"chains_name"**: The name of the chains file
  * **"session_db_name"**: The name of the SessionDB (where Pocket Core store's Sessions)
  * **"evidence_db_name"**: The name of the EvidenceDB (where Pocket Core store's Relay Evidence)
  * **"tendermint_uri"**: The RPC Port of Tendermint (also defined above in Tendermint/RPC)
  * **"keybase_name"**: The name of the keybase
  * **"rpc_port"**: The port of Pocket Core's RPC
  * **"client_block_sync_allowance"**: The +/- allowance in blocks for of a relay request (security mechanism that can help filter misconfigured clients)
  * **"max_evidence_cache_entries"**: Maximum number of relay evidence stored in cache memory
  * **"max_session_cache_entries"**: Maximum number of sessions stored in cache memory
  * **"json_sort_relay_responses"**: Detect and sort if relay response is in json (can help response comparisons if app client is configured for relay consensus)
  * **"remote_cli_url"**: The URL of the CLI (default is local)
  * **"user_agent"**: Custom user agents defined here during http requests
  * **"validator_cache_size"**: Maximum number of validators stored in cache memory
  * **"application_cache_size"**: Maximum number of applications stored in cache memory
  * **"pocket_prometheus_port"**: Pocket port for Prometheus metrics (5.1 +)
  * **"prometheus_max_open_files"**: Max connections to Pocket prometheus 
  * **"max_claim_age_for_proof_retry"**: Maximum age of a claim where a proof transaction will be sent
  * **"proof_prevalidation"**:  Avoid invalid proof transactions by prevalidating claims (extra compute)
  * **"ctx_cache_size"**: Size of the state cache
  * **"abci_logging"**: Log output for transactions and other ABCI calls
  * **"show_relay_errors"**: Print errors for relays executed by the client
### Tendermint
The official Tendermint explanation of the configuration is found [here](https://docs.tendermint.com/master/tendermint-core/configuration.html)
#### Main
  * **"RootDir"**: The data directory of Tendermint (should be the same directory as Pocket Core's data dir)
  * ProxyApp": Pocket Core is always run "in-process", so this typically isn't applicable. However, this configuration is the path of the the TCP connection exposed by Pocket Core.
  * **Moniker"**: The P2P name that will be shown in `Tendermint Peers
  * **"FastSyncMode"**:  Fast sync allows you to process blocks faster when `catching up` to the latest height. With this mode true, the node checks the merkle tree of validators, and doesn't run the real-time consensus gossip protocol.
  * **"LevelDBOptions"**: goleveldb configuration options
  * **"DBPath"**: Path of Tendermint databases local to data directory ("data")
  * **"LogLevel":** The setting for log output in Pocket Core. These levels can be filtered using a simple [log level language](https://blog.cosmos.network/one-of-the-exciting-new-features-in-0-10-0-release-is-smart-log-level-flag-e2506b4ab756): `<Module>:<Level>` in a comma separated list: `main:info, state:debug, p2p:error, *:`
  * **"LogFormat"**: Colored text ("plain") or JSON format ("json")
  * **"Genesis"**: The path of the genesis file local to the data directory (config/genesis.json)
  * **"PrivValidatorKey"**: The path to the keyfile of your private validator (key Tendermint uses for validator operations) local to the data directory "priv_val_key.json"
  * **"PrivValidatorState"**: The path to the validator state file (file Tendermint uses for validator state operations) local to the data directory "priv_val_state.json"
  * **"PrivValidatorListenAddr"**: TCP or UNIX socket address for Tendermint to listen on for
  * connections from an external PrivValidator process. Pocket Core does not utilize the external validator feature, so likely this can be left blank.
  * **"NodeKey"**: The path to the keyfile of your p2p node (key Tendermint uses for p2p operations) NOTE: In Pocket Core, this should always be the same key as the PrivvalKey file.
  * **"ABCI"**: The type of connection between the proxy app and the Tendermint process (grpc or socket)
  * **"ProfListenAddress"**: the path of the profiling server to listen on.
  * **"FilterPeers"**: Allow the ABCI application to filter peers. Pocket Core currently does not utilize this feature of Tendermint (False)

#### RPC
  * **"RootDir"**: The data directory of Tendermint's RPC (should be the same directory as Pocket Core's data dir)
  * **"ListenAddress"**: Tendermint RPC's listening address ("tcp://127.0.0.1:26657")
  * **"CORSAllowedOrigins"**: list of origins a cross-domain request can be executed from. The default value '[]' disables cors support while '["*"]' to allow any origin.
  * **"CORSAllowedMethods"**: String array of allowed Cross Origin Methods ["POST", "GET"]
  * **"CORSAllowedHeaders"**: String array of allowed Cross Origin Headers["Origin", "Accept",],
  * **"GRPCListenAddress"**: TCP or UNIX socket address for the gRPC server to listen (Pocket Core does not utilize gRPC at this time)
  * **"GRPCMaxOpenConnections"**: Maximum allowed conns to the gRPC server
  * **"Unsafe"**: Activate Tendermint unsafe RPC commands like /dial_seeds and /unsafe_flush_mempool.
  * **"MaxOpenConnections"**: Max connections (including WebSocket) to process. (NOTE: this can greatly affect setting System File Descriptors). If set too low, this can affect Consensus participation at scale, if set too high, this can cause `Too Many Open Files`/Resource Consumption. See `guides` of the documentation to properly set your {ulimit -Sn} and subsequently this option.
  * **"MaxSubscriptionClients"**: Maximum number of unique clientIDs that can /subscribe.
  * **"MaxSubscriptionsPerClient"**: Maximum number of unique queries a given client can /subscribe to.
  * **"TimeoutBroadcastTxCommit"**: How long to wait for a tx to be committed during /broadcast_tx_commit (in ns).
  * **"MaxBodyBytes"**: Maximum size of request body, in bytes
  * **"MaxHeaderBytes"**: Maximum size of request header, in bytes
  * **"TLSCertFile"**: The path to a file containing a certificate that is used to create the HTTPS server. NOTE: this option does not affect Pocket Core RPC in any way.
  * **"TLSKeyFile"**: The path to a file containing corresponding private_key that is used to create the HTTPS server. NOTE: this option does not affect Pocket Core RPC in any way.
       
#### P2P
  * **"RootDir"**: The data directory of Tendermint's P2P config (should be the same directory as Pocket Core's data dir)
  * **"ListenAddress"**: The listening address Tendermint will use for peer connections.
  * **"ExternalAddress"**: Address to advertise to peers for them to dial. NOTE: If empty, will use the same port as the laddr
  * **"Seeds"**: The seed nodes used to connect to the network. Must be a comma-separated list in this format: <ADDRESS>@<P2P listening address> (Ex: 03b74fa3c68356bb40d58ecc10129479b159a145@seed1.mainnet.pokt.network:20656). Click here to see a list of seed nodes on [Mainnet](https://docs.pokt.network/docs/mainnet-dispatcher-and-seed-list) or [Testnet](https://docs.pokt.network/docs/known-dispatcher-list)
  * **"PersistentPeers"**: Comma separated list of nodes to keep persistent connections to. Must be a comma separated list in this format: <ADDRESS>@<P2P listening address> (Ex: 03b74fa3c68356bb40d58ecc10129479b159a145@seed1.mainnet.pokt.network:20656)
  * **"UPNP"**: Enable or disable UPNP forwarding.
  * **"AddrBook"**: The path to the addrbook.json file local to the datadir ("config/addrbook.json")
  * **"AddrBookStrict"**: Set true for strict address routability rules, false for local nets.
  *  "MaxNumInboundPeers": Maximum number of simultaneous peer inbound connections.
  * **"MaxNumOutboundPeers"**: Maximum number of simultaneous peer outbound connections.
  * **"FlushThrottleTimeout"**: Time to wait before flushing messages out on the connection in ns
  * **"MaxPacketMsgPayloadSize"**: Maximum size of a message packet payload, in bytes
  * **"SendRate"**: Rate at which packets can be sent, in bytes/second
  * **"RecvRate"**: Rate at which packets can be received, in bytes/second
  * **"PexReactor"**: Set true to enable the (peer-exchange reactor)[https://docs.tendermint.com/master/spec/reactors/pex/pex.html]
  * "SeedMode": Is this node a seed_node? (in which node constantly crawls the network and looks for peers. If another node asks it for addresses, it responds and disconnects)
  * **"PrivatePeerIDs"**: Comma separated list of peer IDs to keep private (will not be gossiped to other peers)
  * **"AllowDuplicateIP"**: Allow peers with duplicated IP's (according to address book)
  * **"HandshakeTimeout"**: Timeout in ns of peer handshaking
  * **"DialTimeout"**: Timeout in ns of peer dialing
  * **"TestDialFail"**: Testing params. Force dial to fail. Ignore if not testing Tendermint
  *  **TestFuzz"**: Testing params. FUzz connection. Ignore if not testing Tendermint.
  * **"TestFuzzConfig"**: Testing params. Fuzz conn config. Ignore if not testing Tendermint.

#### Mempool
* **"RootDir"**: The data directory of Tendermint's Mempool config (should be the same directory as Pocket Core's data dir)
* **"Recheck"**: Recheck determines if the mempool rechecks all pending transactions after a block was committed. Once a block is committed, the mempool removes all valid transactions that were successfully included in the block.
* **"Broadcast"**: Determines whether this node gossips any valid transactions that arrive in mempool. Default is to gossip anything that passes checktx. If this is disabled, transactions are not gossiped, but instead stored locally and added to the next block this node is the proposer.
* **"WalPath"**: This defines the directory where mempool writes the write-ahead logs. These files can be used to reload unbroadcasted transactions if the node crashes.
* **"Size"**: MaxSize of mempool in Transactions
MaxTxsBytes": Max size of ALL Txs in bytes
* **"CacheSize"**: Max memory cache size of mempool in transactions.
* **"MaxTxBytes"**: Max size of Tx in bytes

#### FastSync (only if main/fastsync_mode=true)
Fast Sync version to use:
* **"Version"**: "v1"
- "v0" - the legacy fast sync implementation
- "v1" (default) - refactor of v0 version for better testability
- "v2" - complete redesign of v0, optimized for testability & readability

#### Consensus 
  * **"RootDir"**: The data directory of Tendermint's Consensus config (should be the same directory as Pocket Core's data dir)
  * **"WalPath"**: Path to Conesusns WAL file relative to datadir. Consensus module writes every message to the WAL (write-ahead log) and will replay all the messages of the last height written to WAL before a crash (if such occurs). [See More](https://docs.tendermint.com/master/spec/consensus/wal.html)
  * **"TimeoutPropose**": The timeout in ns, to receive a proposal block from the designated proposer
  * **"TimeoutProposeDelta"**: The timeout `difference` in ns between the current round and the last round (round is reset every valid proposal block)
  * **"TimeoutPrevote"**: The timeout in ns to get 2/3 prevotes from validators
  * **"TimeoutPrevoteDelta"**: The timeout `difference` in ns between the current round of prevoting and the last round (round is reset every valid proposal block)
  * **"TimeoutPrecommit"**: The timeout in ns to get 2/3 precommits from validators
  * **"TimeoutPrecommitDelta"**: The timeout `difference` in ns between the current round of prevoting and the last round (round is reset every valid proposal block)
  * **"TimeoutCommit"**: The timeout in ns to get 2/3 commits from validators
  * **"SkipTimeoutCommit"**: Make progress as soon as we have all the precommits and don't wait for the designated time. (Pocket Network maintains a steady blocktime by marking this option false)
  * **"CreateEmptyBlocks"**: Create empty blocks if no transactions are submitted/in mempool during the interval.
  * **"CreateEmptyBlocksInterval"**: The timeout that must pass in ns before creating an empty block
  * **"PeerGossipSleepDuration"**: Sleep timer for consensus reactor [More Here](https://docs.tendermint.com/master/spec/reactors/consensus/consensus-reactor.html#gossip-data-routine)
  * **"PeerQueryMaj23SleepDuration"**: Sleep timer for consensus reactor [More Here](https://docs.tendermint.com/master/spec/reactors/consensus/consensus-reactor.html#gossip-data-routine)

#### TxIndex
  * **"Indexer"**: What indexer to use for transactions? (Pocket Core currently must use "KV")
  * **"IndexTags"**: Tags (or events) used to index transactions (Pocket Core depends on this functionality for replay attacks)
  * **"IndexAllTags"**: Would you like to index all tags (events)?

#### Instrumentation
  * **"Prometheus"**: Are you using prometheus to track tendermint metrics?
  * **"PrometheusListenAddr"**: If so, on what port?
  * **"MaxOpenConnections"**: What is the maximum number of simultaneous connections you'd like to allow on prometheus?
  * **"Namespace"**: What namespace would you like to use for prometheus?

### Open Files Calculation
Pocket Core operation requires an elevated Ulimit:
```
({ulimit -Sn} >= {MaxNumInboundPeers} + {MaxNumOutboundPeers} + {GRPCMaxOpenConnections} + {MaxOpenConnections} + {Desired Concurrent Pocket RPC connections} + {100 (Constant number of wal, db and other open files)}
```
### Genesis File
Located:
`$HOME/.pocket/config/genesis.json`

[Testnet Genesis File](https://raw.githubusercontent.com/pokt-network/pocket-network-genesis/master/testnet/genesis.json)

[Mainnet Genesis File](https://raw.githubusercontent.com/pokt-network/pocket-network-genesis/master/mainnet/genesis.json)

Use pocket core flags --mainnet or --testnet to automatically write
### Chains.json
Use the CLI or Manually Edit: `$HOME/.pocket/config/chains.json`

Relay Chain ID's can be found [here](https://forum.pokt.network/t/supportedblockchains/607)
```
[
  {
    "id": "0002",
    "url": "http://eth-geth.com",
    "basic_auth": {
      "username": "",
      "password": ""
    }
  }
]
```
## Operation
Operating a Validator requires (at a minimum) some prerequisite basic knowledge of the Pocket Network. 

*This section will cover the basics of:*
- Slashing and Jailing
- Force Unstake
- Economic Incentives

**Slashing And Jailing**

Jailing and Slashing are high level protocol concepts:

*Jailing* a Validator removes them from both protocol service and consensus. 

*Slashing* a Validator burns a percentage of the 'Staked Tokens'

A Validator is jailed and subsequently slashed for not signing (or incorrectly signing) block proposals. More often than not, this is the reason why Validators are jailed. 

Common reasons for not signing blocks are addressed [here](https://github.com/pokt-network/pocket-core/issues/1092)

If a Validator is jailed for too long it will be forcibly removed by the protocol and all Staked Tokens burned

NOTE: `pocket query params` to see protocol level values like `max_jailed_blocks`

**Force Unstake**

If a Validator falls below the minimum stake (due to slashing) it will be forcibly removed by the protocol and all Staked Tokens burned. This feature of the protcol highlights the importance of staking 'well above' the minimum stake.

**Economic Incentives**

For providing infrastructure access to applications, Validators are rewarded proportional to the work they provide. Pocket Core attempts to send a *Claim* and subsequent *Proof* transaction automatically after the `proof_waiting_period` elapses. If both transactions are successful, Tokens are minted to the address of the Validator.

## Resources
[Pocket Core CLI Spec](cli-interface-spec/)

[Pocket Core RPC Spec](rpc-spec.yaml)

[Software Specific Architecture](software%20specific%20architecture.md)

[Protocol Forum](https://forum.pokt.network/)

[Protocol Wiki](https://forum.pokt.network/c/wiki)

[Community Discord](https://discord.com/invite/KRrqfd3tAK)
## FAQs
[Frequently Asked Questions](https://github.com/pokt-network/pocket-core/issues?q=is%3Aissue+label%3Afaq)
