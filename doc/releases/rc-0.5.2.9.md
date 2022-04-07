# RC-0.5.2.9

## RC-0.5.2.9

After nine Beta releases, two Month's worth of continuous internal and external testing, and investigation and QA,
Pocket Network's Engineering team feels the resource problems of RC-0.5.0 are _fixed_ \(see below for known issues\)
with the upcoming RC-0.5.2. Official upgrade guide [here](https://docs.pokt.network/docs/network-upgrade-guide)

### Important Release Notes

1\) **Delete Session.DB before upgrading from RC-0.5.1**

* `rm -rf <datadir>/session.db`

2\) **Run this release with the following environment
variable: `export GODEBUG="madvdontneed=1"`** [Link to Golang Issue](https://github.com/golang/go/issues/42330) 3\) **
Use the default config for all options \(except unique configurations like moniker, external addr, etc\). You have two
options:**

* Remove`/config/config.json` file, execute a CLI command, and update the custom configurations
* Run `pocket util update-configs` command \(creates a new config file and backs up old config file\)

**GoLevelDB is the only supported database from RC-0.5.2 onward**

* If previously using CLevelDB, users might experience incompatibility issues due to known incompatibilities between the
  two
* PNI temporarily will provide a backup datadir to download to avoid syncing from scratch:

  [13K .zip](https://storage.googleapis.com/blockchains-data/backup_5.2.zip)

  [13K .tar.gz](https://storage.googleapis.com/blockchains-data/backup_5.2.tar.gz)

* After uncompressing theses files, place the contents in the `<datadir>/data` folder

### Context And Original Issues

After a series related issues of Pocket Core's RC-0.5.0 were opened \(\#1115 \#1094 \#1116 \#1117 ++\) in October 2020,
PNI opened a formal investigation into the related resource consumption issues of RC-0.5.0 \(and subsequently the more
stable RC-0.5.1\). The main metric of concern with RC-0.5.0 Resources is 'Memory' \(virtual, real, RSS, you name it\),
with a very tangible 'Memory Leak'. 'Relay Stability', though a primary concern for any release, is a secondary concern
for RC-0.5.2 as RC-0.5.1 seemed to solve the immediate, emergency level `Code 66` errors that plagued blocks 6K-7.5K.
Speed is a tertiary concern with RC-0.5.0, taking 10+ hours
to [sync](https://github.com/pokt-network/pocket-core/issues/1089#issuecomment-708439110) to Mainnet Block 7000.

### Tooling

To debug the issues above, several tools were utilized to determine the root causes of all.

Listed in no particular order:

* [x] [Grafana](https://grafana.com/) \(Observibility/Visibility of resources and consensus issues\)
* [x] [Google's PProf](https://github.com/google/pprof) \(CPU and Memory visibility and profile snapshot differences\)
* [x] [GCVIS](https://github.com/davecheney/gcvis) \(Golang garbage collector monitoring\)
* [x] [Docker/Docker-Compose](https://www.docker.com/) \(Clean room simulations\)
* [x] [GCP](https://cloud.google.com/) \(Load testing\)
* [x] [Golang Runtime Pkg](https://golang.org/pkg/runtime/) \(Memstats Testing\)
* [x] [Golang Debug Pkg](https://golang.org/pkg/runtime/debug/) \(FreeOsMemory Testing\)
* [x] [GoLand+Debugger](https://www.jetbrains.com/go/) \(IDE and Debugger\)

  **Debugging and Changelog**

  Immediately, PNI's team recognized
  many [optimizations](https://github.com/pokt-network/pocket-core/commits/staging?after=d637db6bc5d397812fa3b8e68d9ba661f89fc0cc+69&branch=staging)
  to be made within Pocket Core's own source code. This includes the following:

  \`\`\`

* Delete local Relay/Challenge Evidence on Code 66 failures
* Log relay errors to nodes \(don't just return to clients\)
* Added configuration to pre-validate auto transactions
* Sending proofs/claims moved to EndBlock
* Load only Blockmeta for PrevCtx
* Added configurable cache PrevCtx, Validators, and Applications
* Don't broadcast claims/proofs if syncing
* Spread out claims/proofs between non-session blocks
* Added max claim age configuration for proof submission
* Reorganized non-consensus breaking code in Relay/Merkle Verify for efficiency before reads from state
* Configuration to remove ABCILogs
* Fixed \(pseudo\) memory leak in Tendermints RecvPacketMsg\(\)
* Sessions only store addresses and not entire structs
* Only load bare minimum for relay processing
* Add order to AccountTxs query & blockTxsQuery RPC
* Reduce AccountTxsQuery & blockTxsQuery memory footprint

  ```text
  The [results](https://github.com/pokt-network/pocket-core/issues/1089#issuecomment-713197007) were quite significant in both speed and initial resource usage. Subsequently, the following BETA releases targeted bug fixes and small improvements that were a result of the drastic breaking changes from the original Beta.
  ```

* Nondeterministic hash fix
* Code 89 Fix
* Evidence Seal Fix
* Fixes header.TotalTxs !=
* Fixes header.NumTxs !=
* Updating TM version and Version Number to BETA-0.5.2.3
* Upgraded AccountTxs and BlockTxs to use ReducedTxSearch
* Implemented Reduced TxSearch in Tendermint

  ```text
  Will all of this, the speed and 'Relay Stability' concerns seem to be solved. However, the 'Memory Leak' was not fixed. Transparently, the team was surprised and unsure on how to proceed with tackling the issue. One thing that was clear, more visibility was needed to solve the issue. With the addition some much needed tooling (see above), the hunt was on for the leak culprit. Here's a taste of the testing the team did to hunt down this issue:
  ```

* 72 hour simulations in Docker
* Clean Room Relay Stress Tests in GCP
* Mainnet `Validator`and `Full Node` Simulations
* Snapshot comparisons between different versions
* Upgrade Path \(0.5.1-0.5.2\) simulations
* And Much Much More XD

  ```text
  With the help of some close partners and community members, memory offenders were checked off the list:
  ```

* Moved IAVL from Tendermint to Pocket Core
* Call LazyLoadVersion/Store for all queries and PrevCtx\(\)
* Reduced Tendermint P2P EnsurePeers actions to prevent leak
* Lowered P2P config to far more conservative numbers
* Updated FastSync to default to V1
* Exposed default leveldb options
* Switched to only go-leveldb for leak benchmarking/performance reasons
* Child process to run madvdontneed if not set
* Updated P2P configs
* fixed nil txIndexer bug \(Tendermint now sets txindexer and blockstore\)
* removed event type and used Tendermint's abci.Event

  \`\`\`

  Finally, in Beta-0.5.2.8, memory seemed to be at a constant rate.

  **Evidence**

  **IAVL ISSUES**

  ![Screen Shot 2020-11-23 at 4 14 21 PM](https://user-images.githubusercontent.com/18633790/103254979-92ddf680-4955-11eb-830f-52a2e6a61715.png)

#### Memory Bump during a block

![Screen Shot 2020-11-27 at 5 15 32 PM](https://user-images.githubusercontent.com/18633790/103254989-996c6e00-4955-11eb-81cd-3e8f48ebdc7d.png)

#### IAVL NODE CLONE

![Screen Shot 2020-12-04 at 9 42 12 AM](https://user-images.githubusercontent.com/18633790/103254996-9d988b80-4955-11eb-9d55-411368e7a58f.png)

#### Append Events

![Screen Shot 2020-12-11 at 5 08 55 PM](https://user-images.githubusercontent.com/18633790/103255001-a1c4a900-4955-11eb-8478-c0a1fa6c2ca7.png)

#### Tendermint True Bit Indicies

![Screen Shot 2020-11-28 at 3 12 50 PM](https://user-images.githubusercontent.com/18633790/103255044-bc971d80-4955-11eb-92f4-bf8b976ab152.png)

#### Multiple GCVIS heap stability at Beta 5.2.7

![Screen Shot 2020-12-07 at 3 23 43 PM](https://user-images.githubusercontent.com/18633790/103255048-c1f46800-4955-11eb-9054-efa197900d0f.png) ![Screen Shot 2020-12-11 at 3 53 26 PM](https://user-images.githubusercontent.com/18633790/103255051-c6b91c00-4955-11eb-8014-2b31e3be08fc.png) ![Screen Shot 2020-12-11 at 6 04 52 PM](https://user-images.githubusercontent.com/18633790/103255064-d5073800-4955-11eb-8341-0ef7110d28c1.png)

#### Evidence of cache growth from mempool

![Screen Shot 2020-12-11 at 6 13 15 PM](https://user-images.githubusercontent.com/18633790/103255069-d9335580-4955-11eb-9545-0232ff6278bf.png)

### External Reports from community members

![image](https://user-images.githubusercontent.com/18633790/103255381-206e1600-4957-11eb-881c-e5aae7b25404.png) ![image](https://user-images.githubusercontent.com/18633790/103255393-2a901480-4957-11eb-9050-949dce4c0a42.png) ![image](https://user-images.githubusercontent.com/18633790/103255405-3380e600-4957-11eb-9bb2-473228bc62f1.png)

### Disclaimer

Though, the memory seems to be both significantly decreased and stabilized, the team is still not convinced the memory
growth issue is fully fixed \(though not supported with evidence currently\). The team expects to dive deeper and
provide even more visibility into Tendermint and Pocket Core in future releases.

