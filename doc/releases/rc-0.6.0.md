# RC-0.6.0

Another release, another major milestone.

Presenting RC-0.6.0: Pocket Network's first Consensus Rule Change upgrade.

The development of this release began months before the development of RC-0.5.2.9 and just concluded.

The upcoming Pocket Core release \(0.6.0\) offers a higher level of security \(2 mission critical patches in the merkle
tree\) plus provides a higher level of network stability through the removal/patching of events, in addition to a change
in the encoding algorithm \(Amino to Google's Protobuf\).

## Upgrade

### 1. Shutdown Pocket Core

{% hint style="info" %} If Validator, 4 blocks until jailed {% endhint %}

### 2. Ensure golang version 1.16 or &gt; [golang upgrade](https://golang.org/doc/install)

### 3. Build from Source, Homebrew or Docker

#### From Source

To build the latest binary from source, follow these steps:

Navigate into your pocket-core directory: Example: cd ~/go/src/github.com/pokt-network/pocket-core

Enter: pocket version You should see: `RC-0.5.2.10` \(or older\)

To grab the latest packages and tags we are going to swap branches to the latest tag using:

```text
git pull
git checkout tags/RC-0.6.0
```

Once you checked out the latest tag and branch, we are going to rebuild the binary by entering
in: `go build -o $GOPATH/bin/pocket ./app/cmd/pocket_core/main.go`

After it builds, make sure you are on the latest release version by entering in: `pocket version`

Output will be `RC-0.6.0`

#### Homebrew

If you built your binary using Homebrew, follow these steps to upgrade your binary:

In a terminal window, we are going to pull the latest tap by entering: `$ brew upgrade pokt-network/pocket-core/pocket`

After it builds, make sure you are on the latest version by entering in: `pocket version`

Output will be `RC-0.6.0`

#### Docker

For individuals using Docker, all you will need to do to get the new container image is run:

`docker pull poktnetwork/pocket-core:RC-0.6.0`

or

`docker pull poktnetwork/pocket:RC-0.6.0`

Depending on which of the 2 Docker images you want to use.

### 4. **Upgrade your config.json**

Use the default config for all options \(except unique configurations like moniker, external addr, etc\).

```text
You have two options:

- Remove`/config/config.json` file, execute a CLI command, and update the custom configurations
- Run `pocket util update-configs` command (creates a new config file and backs up old config file)
    >In order to use the most performant values, you will need to upgrade your config.json file in within your datadir.
    The following example assumes your datadir to be ~/.pocket, but feel free to swap out your actual datadir with the location in your system.
    Run the pocket util update-configs command to backup your old config and generate a new default one.
    Manually go over your datadir/config/config.json.bk file and update your new datadir/config/config.json with any pertinent values such as moniker, external addr, etc.
```

### 5. **Delete Session.DB before upgrading**

* `rm -rf <datadir>/session.db`

### 6\) NOTE: **Step 6 IS ONLY NEEDED IF RUNNING VERSION &lt; RC-0.5.2.9 OR DB CORRUPTED** **GoLevelDB is the only
supported database from RC-0.5.2 onward**

* If previously using CLevelDB, users might experience incompatibility issues due to known incompatibilities between the
  two
* PNI temporarily will provide a backup datadir to download to avoid syncing from scratch:

  [22K .zip](https://storage.googleapis.com/blockchains-data/backup_data_22150.zip)

* After uncompressing theses files, place the contents in the `<datadir>/data` folder

### Special Notes

#### **Consensus rule change**

Consensus rules are the fabric of the protocol that requires 66% &gt; agreement of Validators in order to reach quorum
on the blockchain data.

This edit to the existing Pocket Core Software defines a \(dormant\) new set of consensus rules that can be activated
with a Validator approved governance transaction \(See Changelog Below\)

This release of Pocket Core supports legacy \(RC-0.5\) consensus rules as well as the new \(RC-0.6\) consensus rules.

The software will not change to the new consensus rules until activated by a 66% majority Validator support of a DAO
initiated transaction that specifies the height at which the 'rules change'.

#### **Transaction Param \(Legacy-Codec\)**

As described below, the 6.0 upgrade contains a codec upgrade.

Submitting transactions with this release before the upgrade height \(will add height here once the DAO votes\) will
require the 'legacy codec' argument to be TRUE

Submitting transactions with this release after the upgrade height will require the 'legacy codec' argument to be FALSE
\(DEFAULT\)

## FAQ

### Can I upgrade before of the upgrade height?

Yes, as soon as the release is published. Watch the repo here to be notified when the release is published.

We’ll also make an announcement about the release in the \#announcements channel of our discord.

### **What happens if I do not upgrade in proposed time?**

The Pocket Core process will not be able to continue \(shutdown automatically and cannot be restarted\)

### **What do I do If I am using a third-party service provider to run my nodes?**

If you are using a third-party service provider it will be up to them to upgrade your nodes, but we do recommend that
you contact them for their upgrade plan.

We are also coordinating with the major third-party providers directly in order to ensure a smooth upgrade.

### **What happens if I do not upgrade in time?**

Your node will not be able to continue. Full nodes will shut down and be unable to restart without an upgrade.
Validators will be slashed and jailed and be unable to restart without an upgrade. If this happens to you, the only way
to get your node back online will be to follow the directions outlined in the upgrade guides.

Note: if you manage to recover within 6 blocks of missing the upgrade time, you won’t actually be slashed. But why take
that risk?

### **If I do not upgrade in time, will I lose my POKT?**

As explained above, the upgrade force shuts down any node running older versions after the upgrade height.

For Validators, the shutdown can result in standard offline slashing.

## Proposal

[LINK](https://forum.pokt.network/t/pip-4-consensus-rule-change-0-6-0/834)

### **Motivation**

There are two major security issues in the merkle tree proof/claim implementation as well as an exploitable prediction
attack due to a misimplementation at the block hash generation. The current encoding scheme is both 'custom' and
unsupported across most all programming languages which hinders ecosystem growth and future development. Lastly, PUP-4
is somewhat addressed in this release.

### Specification

* Convert all consensus level amino encoding \(including but not limited to the internal storage codecs\) to protobuf
  encoding while maintaining as many legacy structures as possible
* Introduce Previous Block Validator Voting structure into the block hash used for session and proposer selection
  algorithms.
* Use the index of the leafs of the plasma core merkle tree as part of the parent hash to lock in the values using the
  Claim merkle root
* Ensure consensus level events are not concatenated in the pocket core module by initializing in the transaction
  handler
* Change ABCIValidatorUpdate to ABCIValidatorZeroUpdate for separation of service and validation

### Rationale

The bug fixes in the merkle tree result in an increased level of network stability. Applications and node runners will
experience an even higher degree of reliability through the new found network security.

Through the addition of Protobuf encoding, client-side tooling such as SDK development and improvements just got a lot
easier, which will make expanding our potential app user bases easier and an all around better development experience
while using Pocket Network.

In addition a bug identified in event-handling has now been fixed, which creates smaller block sizes and should enable
faster txs and overall better service.

We have successfully separated servicing and validation, which allows us to have more nodes overall and more scalability
- no longer capped to 5000\(technically\). That said, PUP4 will likely still try to limit nodes to 5000 due to the lack
of jailing available to servicers which may lead to service degradation.

Protobuf encoding will also lower transaction latency a tad because of less resource demand on nodes.

### Viability

An extensive number of tests, functional, integration, unit, load, and simulation were completed leading up to this
upgrade. The reports will be included in the release notes but we can't release until the DAO has finalized an upgrade
height.

### Implementation

The implementation of 0.6.0 is near complete. A few pending tests, the agreement of an Upgrade height, and the approval
of this proposal, will result in a complete implementation.

### Audit

There was no external audit, refer to Viability.

## Changelog

* Security patch in merkle sum index \(CONSENSUS RULE CHANGE\)

  > Hashes the index of the leafs in the merkle sum index tree

* Security patch for BlockHash \(CONSENSUS RULE CHANGE\)

  > Uses the current block hash for session generation and not the lastBlockHash

* Security patch for Proposer Selection \(CONSENSUS RULE CHANGE\)

  > Uses the current evidence in the proposer selection and not the lastBlockHash

* Protobuf Encoding implemented \(CONSENSUS RULE CHANGE\)

  > Switching away from amino completely \(except for keybase encodings for legacy compatibility\)

* Pocketcore module event handler fixed \(CONSENSUS RULE CHANGE\)

  > Events are no longer concatenated

* Max age of evidence enforced in blocks \(CONSENSUS RULE CHANGE\)

  > Ensure double signs from blocks exceeding the param age are no longer slashable events

* Deleted applications from state after unstake \(CONSENSUS RULE CHANGE\)

  > Remove unstaked apps from the state to save blockchain size

## **Tooling**

To debug the issues above, several tools were utilized to determine the root causes of all.

Listed in no particular order:

* [x] [Grafana](https://grafana.com/) \(Observibility/Visibility of resources and consensus issues\)
* [x] [Google's PProf](https://github.com/google/pprof) \(CPU and Memory visibility and profile snapshot differences\)
* [x] [Docker/Docker-Compose](https://www.docker.com/) \(Clean room simulations\)
* [x] [GCP](https://cloud.google.com/) \(Load testing\)
* [x] [GoLand+Debugger](https://www.jetbrains.com/go/) \(IDE and Debugger\)

## Reports

### Load Tests

[https://github.com/pokt-network/organization/tree/main/milestones/1-08-01-2021/6.0/projects/func\_test](https://github.com/pokt-network/organization/tree/main/milestones/1-08-01-2021/6.0/projects/func_test)

### Resource Benchmarks

[https://github.com/pokt-network/organization/tree/main/milestones/1-08-01-2021/6.0/projects/benchmark\_test](https://github.com/pokt-network/organization/tree/main/milestones/1-08-01-2021/6.0/projects/benchmark_test)

### Functional Tests

[https://github.com/pokt-network/organization/blob/main/milestones/1-08-01-2021/6.0/projects/func\_test/Functional Test Execution - RC-0.6.0.xlsx](https://github.com/pokt-network/organization/blob/main/milestones/1-08-01-2021/6.0/projects/func_test/Functional%20Test%20Execution%20-%20RC-0.6.0.xlsx)

### Unit/Behavior Tests

[https://app.circleci.com/pipelines/github/pokt-network/pocket-core](https://app.circleci.com/pipelines/github/pokt-network/pocket-core)

