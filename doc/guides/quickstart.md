# Quickstart

{% hint style="danger" %}
_Disclaimer: Pocket Network Inc. \(PNI\) is not liable for any slashing or economic penalty that may occur._
{% endhint %}

## Install

### From Source

#### Prerequisite Installations

* [go](https://golang.org/doc/install)
* [go environment](https://golang.org/doc/gopath_code.html#Workspaces) GOPATH & GOBIN
* [git](https://git-scm.com/book/en/v2/Getting-Started-Installing-Git)

#### Create source code directory

```text
mkdir -p $GOPATH/src/github.com/pokt-network && cd $GOPATH/src/github.com/pokt-network
```

#### Download the source code

```text
git clone https://github.com/pokt-network/pocket-core.git && cd pocket-core
```

#### Checkout the [latest release](https://github.com/pokt-network/pocket-core/releases)

```text
git checkout tags/<release tag>
```

Example:

```text
git checkout tags/RC-0.8.0
```

#### Build

```text
go build -o <destination directory> <source code directory>/...
```

Example:

```text
go build -o $GOPATH/bin/pocket $GOPATH/src/github.com/pokt-network/pocket-core/app/cmd/pocket_core/main.go
```

#### Test installation

{% tabs %} {% tab title="Command" %}

```text
pocket version
```

{% endtab %}

{% tab title="Response" %}

```
> RC 0.8.0
```

{% endtab %} {% endtabs %}

### From Homebrew

```text
brew tap pokt-network/pocket-core && brew install pokt-network/pocket-core/pocket
```

#### Test installation

{% tabs %} {% tab title="Command" %}

```text
pocket version
```

{% endtab %}

{% tab title="Response" %}

```
> RC-0.8.0
```

{% endtab %} {% endtabs %}

### From Deployment Artifact

See [pokt-network/pocket-core-deployments](https://github.com/pokt-network/pocket-core-deployments)

## Quickstart

### Prerequisite Knowledge

{% hint style="info" %} This section does not cover the protocol specification, rather how to participate in the network
as a node runner. For more information on the Pocket Network protocol, read the
wiki [here](https://docs.pokt.network/main-concepts/protocol). {% endhint %}

* A **Validator** is an infrastructure provider in Pocket Network
* **Staking** a Validator locks up a certain **amount** of balance that can be burned as a security mechanism for bad
  acting
* A **Relay Chain** is blockchain infrastructure **Validators** expose for application access _Ex: Ethereum, Bitcoin,
  Pocket Network_ \(identified by 4 hexadecimal characters. _Ex: 0001_\)
* Apps access **Relay Chains** through the **serviceURI**: the endpoint where **Validators** publicly expose the Pocket
  API _Ex:_ [https://www.node1.mainnet.pokt.network](https://www.node1.mainnet.pokt.network)

### Environment Setup

* **Hardware Requirements:** 4 CPU’s \(or vCPU’s\) \| 8 GB RAM \| 100GB Disk
* **Reverse Proxy:** For SSL termination and request management
* **Ports:** Expose Pocket RPC \(Default :8081\) and P2P port \(Default: 26656\)
* **SSL Cert:** Required for **Validator's serviceURI**
* **Open Files Limit:** `ulimit -Sn 16384`

{% hint style="warning" %}
**Open Files Limit** is very important for the operation of Pocket Core. See [**
Config**](quickstart.md#open-files-calculation) section for ulimit calculation. {% endhint %}

### Create an account

An account is needed to participate at any level of the network.

{% tabs %} {% tab title="Command" %}

```text
pocket accounts create
```

{% endtab %}

{% tab title="Response" %}

```
> Enter Passphrase
> Account generated successfully:
> Address: <address>
```

{% endtab %} {% endtabs %}

### Fund the account

To stake a Validator in Pocket Network, the account must have a balance above the **minimum stake**:

`15,000 POKT` or `15,000,000,000 uPOKT`

{% hint style="danger" %} If your stake falls below `15,000 POKT` your node will be force-unstake burned. We recommend
having a buffer above the 15,000 minimum \(e.g. 15,100-16,000\), so that minor slashing doesn't result in loss of the
entire stake. {% endhint %}

Send POKT with the following command:

```text
pocket accounts send-tx <fromAddr> <toAddr> <uPOKT amount> mainnet 10000 "" true
```

### Set the account as Validator

```text
pocket accounts set-validator <address>
```

{% hint style="info" %} Check with `pocket accounts get-validator`
{% endhint %}

### Set [Relay Chains](https://docs.pokt.network/home/references/supported-blockchains)

{% tabs %} {% tab title="Command" %}

```text
pocket util generate-chains
```

{% endtab %}

{% tab title="Response" %}

```
> Enter the chain of the network identifier:
<Relay Chain ID> (Example: 0001)
> Enter the URL of the network identifier:
<Secure URL to Relay Chain>
Would you like to enter another network identifier? (y/n)
n
```

{% endtab %} {% endtabs %}

{% hint style="info" %} Can test with simulate relay flag and endpoint. See [RPC Specification](../specs/rpc-spec.md)
for details. {% endhint %}

### Sync the blockchain

```text
pocket start --seeds=<seeds> --mainnet
```

Example:

```text
pocket start --seeds="64c91701ea98440bc3674fdb9a99311461cdfd6f@node1.mainnet.pokt.network:21656" --mainnet
```

[Seeds](https://docs.pokt.network/references/seeds)

{% hint style="warning" %} Ensure the node is all the way synced before proceeding to the next step. {% endhint %}

### Stake the Validator

Stake the account to participate in the Network as a **Validator**

```text
pocket nodes stake <address> <amount> <relay_chains> <serviceURI>:<port> mainnet 10000
```

Example:

```text
pocket nodes stake 3ee61299d5bbbd2974cddcc194d9b547c7629546 20000000000 ["0001","0002"] https://pokt.rocks:443 mainnet 10000
```

{% hint style="danger" %} Remember, stake 'well over' the minimum stake to avoid force-unstake burning. {% endhint %}

## Config

### Data Directory

Pocket Core files are located in a **Data Directory** Default: `$HOME/.pocket/`

### Configuration File

Pocket Core provides a configuration file found in `<datadir>/config/config.json`

### Pocket

* **"data\_dir"**: The data directory of Pocket Core \(should be the same directory as Tendermint data dir\)
* **"genesis\_file"**: The name of the genesis file
* **"chains\_name"**: The name of the chains file
* **"evidence\_db\_name"**: The name of the EvidenceDB \(where Pocket Core store's Relay Evidence\)
* **"tendermint\_uri"**: The RPC Port of Tendermint \(also defined above in Tendermint/RPC\)
* **"keybase\_name"**: The name of the keybase
* **"rpc\_port"**: The port of Pocket Core's RPC
* **"client\_block\_sync\_allowance"**: The +/- allowance in blocks for of a relay request \(security mechanism that can
  help filter misconfigured clients\)
* **"max\_evidence\_cache\_entries"**: Maximum number of relay evidence stored in cache memory
* **"max\_session\_cache\_entries"**: Maximum number of sessions stored in cache memory
* **"json\_sort\_relay\_responses"**: Detect and sort if relay response is in json \(can help response comparisons if
  app client is configured for relay consensus\)
* **"remote\_cli\_url"**: The URL of the CLI \(default is local\)
* **"user\_agent"**: Custom user agents defined here during http requests
* **"validator\_cache\_size"**: Maximum number of validators stored in cache memory
* **"application\_cache\_size"**: Maximum number of applications stored in cache memory
* **"pocket\_prometheus\_port"**: Pocket port for Prometheus metrics \(5.1 +\)
* **"prometheus\_max\_open\_files"**: Max connections to Pocket prometheus
* **"max\_claim\_age\_for\_proof\_retry"**: Maximum age of a claim where a proof transaction will be sent
* **"proof\_prevalidation"**:  Avoid invalid proof transactions by prevalidating claims \(extra compute\)
* **"ctx\_cache\_size"**: Size of the state cache
* **"abci\_logging"**: Log output for transactions and other ABCI calls
* **"show\_relay\_errors"**: Print errors for relays executed by the client

  **Tendermint**

  The official Tendermint explanation of the configuration is
  found [here](https://docs.tendermint.com/master/tendermint-core/configuration.html)

  **Main**

* **"RootDir"**: The data directory of Tendermint \(should be the same directory as Pocket Core's data dir\)
* ProxyApp": Pocket Core is always run "in-process", so this typically isn't applicable. However, this configuration is
  the path of the the TCP connection exposed by Pocket Core.
* **Moniker"**: The P2P name that will be shown in \`Tendermint Peers
* **"FastSyncMode"**:  Fast sync allows you to process blocks faster when `catching up` to the latest height. With this
  mode true, the node checks the merkle tree of validators, and doesn't run the real-time consensus gossip protocol.
* **"LevelDBOptions"**: goleveldb configuration options
* **"DBPath"**: Path of Tendermint databases local to data directory \("data"\)
* **"LogLevel":** The setting for log output in Pocket Core. These levels can be filtered using a
  simple [log level language](https://blog.cosmos.network/one-of-the-exciting-new-features-in-0-10-0-release-is-smart-log-level-flag-e2506b4ab756): `<Module>:<Level>`
  in a comma separated list: `main:info, state:debug, p2p:error, *:`
* **"LogFormat"**: Colored text \("plain"\) or JSON format \("json"\)
* **"Genesis"**: The path of the genesis file local to the data directory \(config/genesis.json\)
* **"PrivValidatorKey"**: The path to the keyfile of your private validator \(key Tendermint uses for validator
  operations\) local to the data directory "priv\_val\_key.json"
* **"PrivValidatorState"**: The path to the validator state file \(file Tendermint uses for validator state operations\)
  local to the data directory "priv\_val\_state.json"
* **"PrivValidatorListenAddr"**: TCP or UNIX socket address for Tendermint to listen on for
* connections from an external PrivValidator process. Pocket Core does not utilize the external validator feature, so
  likely this can be left blank.
* **"NodeKey"**: The path to the keyfile of your p2p node \(key Tendermint uses for p2p operations\) NOTE: In Pocket
  Core, this should always be the same key as the PrivvalKey file.
* **"ABCI"**: The type of connection between the proxy app and the Tendermint process \(grpc or socket\)
* **"ProfListenAddress"**: the path of the profiling server to listen on.
* **"FilterPeers"**: Allow the ABCI application to filter peers. Pocket Core currently does not utilize this feature of
  Tendermint \(False\)

#### RPC

* **"RootDir"**: The data directory of Tendermint's RPC \(should be the same directory as Pocket Core's data dir\)
* **"ListenAddress"**: Tendermint RPC's listening address \("tcp://127.0.0.1:26657"\)
* **"CORSAllowedOrigins"**: list of origins a cross-domain request can be executed from. The default value '\[\]'
  disables cors support while '\["\*"\]' to allow any origin.
* **"CORSAllowedMethods"**: String array of allowed Cross Origin Methods \["POST", "GET"\]
* **"CORSAllowedHeaders"**: String array of allowed Cross Origin Headers\["Origin", "Accept",\],
* **"GRPCListenAddress"**: TCP or UNIX socket address for the gRPC server to listen \(Pocket Core does not utilize gRPC
  at this time\)
* **"GRPCMaxOpenConnections"**: Maximum allowed conns to the gRPC server
* **"Unsafe"**: Activate Tendermint unsafe RPC commands like /dial\_seeds and /unsafe\_flush\_mempool.
* **"MaxOpenConnections"**: Max connections \(including WebSocket\) to process. \(NOTE: this can greatly affect setting
  System File Descriptors\). If set too low, this can affect Consensus participation at scale, if set too high, this can
  cause `Too Many Open Files`/Resource Consumption. See `guides` of the documentation to properly set your {ulimit -Sn}
  and subsequently this option.
* **"MaxSubscriptionClients"**: Maximum number of unique clientIDs that can /subscribe.
* **"MaxSubscriptionsPerClient"**: Maximum number of unique queries a given client can /subscribe to.
* **"TimeoutBroadcastTxCommit"**: How long to wait for a tx to be committed during /broadcast\_tx\_commit \(in ns\).
* **"MaxBodyBytes"**: Maximum size of request body, in bytes
* **"MaxHeaderBytes"**: Maximum size of request header, in bytes
* **"TLSCertFile"**: The path to a file containing a certificate that is used to create the HTTPS server. NOTE: this
  option does not affect Pocket Core RPC in any way.
* **"TLSKeyFile"**: The path to a file containing corresponding private\_key that is used to create the HTTPS server.
  NOTE: this option does not affect Pocket Core RPC in any way.

#### P2P

* **"RootDir"**: The data directory of Tendermint's P2P config \(should be the same directory as Pocket Core's data
  dir\)
* **"ListenAddress"**: The listening address Tendermint will use for peer connections.
* **"ExternalAddress"**: Address to advertise to peers for them to dial. NOTE: If empty, will use the same port as the
  laddr
* **"Seeds"**: The seed nodes used to connect to the network. Must be a comma-separated list in this format:
  &lt;ADDRESS&gt;@ \(Ex: 03b74fa3c68356bb40d58ecc10129479b159a145@seed1.mainnet.pokt.network:20656\).
  Click [here](https://docs.pokt.network/home/resources/references/seeds) to see a list of seed nodes on Mainnet or
  Testnet
* **"PersistentPeers"**: Comma separated list of nodes to keep persistent connections to. Must be a comma separated list
  in this format: &lt;ADDRESS&gt;@ \(Ex: 03b74fa3c68356bb40d58ecc10129479b159a145@seed1.mainnet.pokt.network:20656\)
* **"UPNP"**: Enable or disable UPNP forwarding.
* **"AddrBook"**: The path to the addrbook.json file local to the datadir \("config/addrbook.json"\)
* **"AddrBookStrict"**: Set true for strict address routability rules, false for local nets.
* "MaxNumInboundPeers": Maximum number of simultaneous peer inbound connections.
* **"MaxNumOutboundPeers"**: Maximum number of simultaneous peer outbound connections.
* **"FlushThrottleTimeout"**: Time to wait before flushing messages out on the connection in ns
* **"MaxPacketMsgPayloadSize"**: Maximum size of a message packet payload, in bytes
* **"SendRate"**: Rate at which packets can be sent, in bytes/second
* **"RecvRate"**: Rate at which packets can be received, in bytes/second
* **"PexReactor"**: Set true to enable the \(peer-exchange
  reactor\)\[[https://docs.tendermint.com/master/spec/reactors/pex/pex.html](https://docs.tendermint.com/master/spec/reactors/pex/pex.html)\]
* "SeedMode": Is this node a seed\_node? \(in which node constantly crawls the network and looks for peers. If another
  node asks it for addresses, it responds and disconnects\)
* **"PrivatePeerIDs"**: Comma separated list of peer IDs to keep private \(will not be gossiped to other peers\)
* **"AllowDuplicateIP"**: Allow peers with duplicated IP's \(according to address book\)
* **"HandshakeTimeout"**: Timeout in ns of peer handshaking
* **"DialTimeout"**: Timeout in ns of peer dialing
* **"TestDialFail"**: Testing params. Force dial to fail. Ignore if not testing Tendermint
* **TestFuzz"**: Testing params. FUzz connection. Ignore if not testing Tendermint.
* **"TestFuzzConfig"**: Testing params. Fuzz conn config. Ignore if not testing Tendermint.

#### Mempool

* **"RootDir"**: The data directory of Tendermint's Mempool config \(should be the same directory as Pocket Core's data
  dir\)
* **"Recheck"**: Recheck determines if the mempool rechecks all pending transactions after a block was committed. Once a
  block is committed, the mempool removes all valid transactions that were successfully included in the block.
* **"Broadcast"**: Determines whether this node gossips any valid transactions that arrive in mempool. Default is to
  gossip anything that passes checktx. If this is disabled, transactions are not gossiped, but instead stored locally
  and added to the next block this node is the proposer.
* **"WalPath"**: This defines the directory where mempool writes the write-ahead logs. These files can be used to reload
  unbroadcasted transactions if the node crashes.
* **"Size"**: MaxSize of mempool in Transactions

  MaxTxsBytes": Max size of ALL Txs in bytes

* **"CacheSize"**: Max memory cache size of mempool in transactions.
* **"MaxTxBytes"**: Max size of Tx in bytes

#### FastSync \(only if main/fastsync\_mode=true\)

Fast Sync version to use:

* **"Version"**: "v1"
* "v0" - the legacy fast sync implementation
* "v1" \(default\) - refactor of v0 version for better testability
* "v2" - complete redesign of v0, optimized for testability & readability

#### Consensus

* **"RootDir"**: The data directory of Tendermint's Consensus config \(should be the same directory as Pocket Core's
  data dir\)
* **"WalPath"**: Path to Conesusns WAL file relative to datadir. Consensus module writes every message to the WAL
  \(write-ahead log\) and will replay all the messages of the last height written to WAL before a crash \(if such
  occurs\). [See More](https://docs.tendermint.com/master/spec/consensus/wal.html)
* **"TimeoutPropose**": The timeout in ns, to receive a proposal block from the designated proposer
* **"TimeoutProposeDelta"**: The timeout `difference` in ns between the current round and the last round \(round is
  reset every valid proposal block\)
* **"TimeoutPrevote"**: The timeout in ns to get 2/3 prevotes from validators
* **"TimeoutPrevoteDelta"**: The timeout `difference` in ns between the current round of prevoting and the last round
  \(round is reset every valid proposal block\)
* **"TimeoutPrecommit"**: The timeout in ns to get 2/3 precommits from validators
* **"TimeoutPrecommitDelta"**: The timeout `difference` in ns between the current round of prevoting and the last round
  \(round is reset every valid proposal block\)
* **"TimeoutCommit"**: The timeout in ns to get 2/3 commits from validators
* **"SkipTimeoutCommit"**: Make progress as soon as we have all the precommits and don't wait for the designated time.
  \(Pocket Network maintains a steady blocktime by marking this option false\)
* **"CreateEmptyBlocks"**: Create empty blocks if no transactions are submitted/in mempool during the interval.
* **"CreateEmptyBlocksInterval"**: The timeout that must pass in ns before creating an empty block
* **"PeerGossipSleepDuration"**: Sleep timer for consensus
  reactor [More Here](https://docs.tendermint.com/master/spec/reactors/consensus/consensus-reactor.html#gossip-data-routine)
* **"PeerQueryMaj23SleepDuration"**: Sleep timer for consensus
  reactor [More Here](https://docs.tendermint.com/master/spec/reactors/consensus/consensus-reactor.html#gossip-data-routine)

#### TxIndex

* **"Indexer"**: What indexer to use for transactions? \(Pocket Core currently must use "KV"\)
* **"IndexTags"**: Tags \(or events\) used to index transactions \(Pocket Core depends on this functionality for replay
  attacks\)
* **"IndexAllTags"**: Would you like to index all tags \(events\)?

#### Instrumentation

* **"Prometheus"**: Are you using prometheus to track tendermint metrics?
* **"PrometheusListenAddr"**: If so, on what port?
* **"MaxOpenConnections"**: What is the maximum number of simultaneous connections you'd like to allow on prometheus?
* **"Namespace"**: What namespace would you like to use for prometheus?

### Open Files Calculation

Pocket Core operation requires an elevated Ulimit:

```text
({ulimit -Sn} >= {MaxNumInboundPeers} + {MaxNumOutboundPeers} + {GRPCMaxOpenConnections} + {MaxOpenConnections} + {Desired Concurrent Pocket RPC connections} + {100 (Constant number of wal, db and other open files)}
```

### Genesis File

Located: `$HOME/.pocket/config/genesis.json`

* [Testnet Genesis File](https://raw.githubusercontent.com/pokt-network/pocket-network-genesis/master/testnet/genesis.json)
* [Mainnet Genesis File](https://raw.githubusercontent.com/pokt-network/pocket-network-genesis/master/mainnet/genesis.json)

Use pocket core flags --mainnet or --testnet to automatically write

### Chains.json

Use the CLI or Manually Edit: `$HOME/.pocket/config/chains.json`

{% hint style="info" %} Relay Chain ID's can be
found [here](https://docs.pokt.network/home/references/supported-blockchains). {% endhint %}

```text
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

Operating a Validator requires \(at a minimum\) some prerequisite basic knowledge of the Pocket Network.

This section will cover the basics of:

* Slashing and Jailing
* Force Unstake
* Economic Incentives

### **Slashing And Jailing**

Jailing and Slashing are high level protocol concepts:

* _Jailing_ a Validator removes them from both protocol service and consensus.
* _Slashing_ a Validator burns a percentage of the 'Staked Tokens'

A Validator is jailed and subsequently slashed for not signing \(or incorrectly signing\) block proposals. More often
than not, this is the reason why Validators are jailed.

Common reasons for not signing blocks are addressed [here](https://github.com/pokt-network/pocket-core/issues/1092).

If a Validator is jailed for too long it will be forcibly removed by the protocol and all Staked Tokens burned.

{% hint style="info" %}
`pocket query params` to see protocol level values like `max_jailed_blocks`
{% endhint %}

### **Force Unstake**

If a Validator falls below the minimum stake \(due to slashing\) it will be forcibly removed by the protocol and all
Staked Tokens burned. This feature of the protocol highlights the importance of staking 'well above' the minimum stake.

{% hint style="danger" %} If your stake falls below `15,000 POKT` your node will be force-unstake burned. We recommend
having a buffer above the 15,000 minimum \(e.g. 15,100-16,000\), so that minor slashing doesn't result in loss of the
entire stake. {% endhint %}

### **Economic Incentives**

For providing infrastructure access to applications, Validators are rewarded proportional to the work they provide.
Pocket Core attempts to send a _Claim_ and subsequent _Proof_ transaction automatically after the `proof_waiting_period`
elapses. If both transactions are successful, Tokens are minted to the address of the Validator.

