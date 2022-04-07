---
description: >- Pocket Core is the Official Golang implementation of the Pocket Network protocol.
---

# Pocket Core Architecture

{% hint style="info" %} This document DOES NOT cover the architecture of the Pocket Network protocol. Rather it only
covers the software-specific implementation details of Pocket Core and for the sake of clarity, orientation, and
modularity, explicitly avoids covering any details that might duplicate the contents of the Pocket Network protocol
specification, which can be found [here](https://docs.pokt.network/home/main-concepts/protocol). {% endhint %}

## Tendermint and Cosmos SDK

Pocket Core is technically an ABCI application built on top of a modified version
of [Tendermint](https://tendermint.com/): a state machine replication engine built natively in Golang. Tendermint
abstracts a BFT consensus algorithm that replicates Pocket Core's state machine logic among all nodes, creating the
decentralized network.

The decision to use Tendermint was made out of both necessity and convenience. Pocket Core _needs_ a state machine
replication module and Tendermint _conveniently_ modularized just that. Though multiple **protocol level** modifications
were made to Tendermint, this document will only cover the **implementation level** modifications:

1. Pocket Core restricts the use of GoLevelDB only
2. Pocket Core exposes LevelDB options in configuration file
3. Pocket Core patches some minor P2P issues
4. Pocket Core patches some minor resource issues

{% hint style="info" %} This document does not cover points 3 and 4, however it is covered in
the [commit history](https://github.com/pokt-network/tendermint/commits/). {% endhint %}

Pocket Core restricts the use of GolevelDB for a few basic reasons.

1. GoLevelDB [has proven](https://github.com/pokt-network/pocket-core/releases/tag/RC-0.5.2.9) to be the most efficient
   and reliable database for Tendermint
2. GolevelDB is not _always_ cross compatible with ClevelDB and certainly not compatible with other non-LevelDB
   implementations
3. The effort of maintaining multiple _experimental_ \(as Tendermint refers to them\) database implementations is quite
   high, and is currently not a priority for PNI.
4. GolevelDB is able to be monitored and benchmarked with Golang Library tools
   like [PPROF](https://github.com/google/pprof) and a single database implementation is much easier to standardize
   resource consumption.

Pocket Core exposes LevelDB options in the Tendermint configuration file

1. LevelDB configurations are provided by the library to allow for users to perform application level modifications to
   the storage layer of the blockchain application.
2. Tendermint [often](https://github.com/cosmos/cosmos-sdk/issues/1394#issuecomment-402819672) suffers from 'Open File
   Limit' issues. GolevelDB provides a configuration for this

Though Pocket Core does rely on Tendermint for its replication / consensus module, Pocket Core is often confused as a
modified Cosmos SDK: _An SDK for building application specific blockchains._

Transparently, the original idea was to use the Cosmos SDK with little modification to suit the Proof-Of-Stake needs of
the Pocket Network Protocol. However, the PNI Engineering team quickly realized the Cosmos SDK was really not an SDK in
the traditional sense, rather a somewhat loose forkable framework. 1K+ commits later, Pocket Core shares very little
code with the Cosmos SDK other than its loose Module structure and
the [KV store implementation](https://github.com/pokt-network/pocket-core/tree/staging/store) \(both of which are
planned to be overhauled in the future\).

## Relay Evidence Collection

Terminology: _Relay Evidence - Proof of a Relay Completed_

Pocket Network implementations must collect Relay Evidence in order to participate in the network's economic incentive
mechanisms.

Pocket Core's Relay Evidence collection is broken down into two major pieces:

1. A standard unordered golang slice \(dynamic array\) to store the evidence
2. An adjacent [bloom filter](https://www.eecs.harvard.edu/~michaelm/NEWWORK/postscripts/BloomFilterSurvey.pdf) for
   time/space efficient membership checks

**Evidence Store:**

For speed and quality of service reasons, Pocket Core chooses to store the evidence unordered \(as it comes in from the
client\). The evidence is then ordered using a quicksort algorithm at time of the claim/proof submission.

Pocket Core also dynamically grows the storage array \(slice\) with the current assumption that most Applications will
not use close to their maximum relay capacity.

_Alternatives:_

* Pocket Core could store the evidence ordered with the tradeoff of quality of service and speed to the client. This
  comes with a compute benefit at the time of the Merkle Tree generation \(proof/claim submisison\)
* Pocket Core could fix the array to the maximum size for each application, because the relay capacity is known at the
  time of array creation. This comes with a tradeoff in higher memory consumption but for a decrease in compute.

**Bloom Filter**

Duplicate Relay Evidence can easily result in a _code 66_ invalid Merkle Proof Protocol error, or worse, trigger an
economic Replay Attack defense mechanism.

In order to prevent duplicate relay evidence, the protocol implementation should check for membership of the Relay
Evidence before accepting anymore incoming requests.

Pocket Core uses a bit array Bloom Filter implementation to provide an efficient structure for quick membership checks.
Due to the properties of a Bloom Filter, Pocket Core will never accept duplicate evidence and is able to check for
membership with O\(1\) time

_Alternatives:_

* Use a standard Golang Map to store the Evidence and just use that structure to check for membership. Also O\(1\) time.
  The tradeoff here is more memory because the membership is stored as a hash vs bits.

## Automatic Randomized Proof and Claim Submissions

In order to participate in the Pocket Network economic incentive mechanisms, proof and claim messages must be submitted
to claim and prove the work \(relays\) done.

The protocol, however, does not specify how or when to execute the proof/claim transactions, and leaves it up to the
implementation.

Pocket Core's Proof and Claim submission algorithm is summarized in two unique points:

1. Automatic Submission
2. Address Randomization

**Automatic Submission**

Contrary to most blockchain clients, Pocket Core attempts to automatically submit claim and proof transactions if there
is Relay Evidence stored in the Evidence database and the validator is staked for work. Since the private key of any
Tendermint Validator is already available and exposed for block signing, Pocket Core uses the PK to sign the Claim/Proof
transactions.

_Alternative_

* Pocket Core can ask the users to manually submit the claim and proof transaction or offload the work to a third party
  application. The tradeoff is both convenience and robustness of the application

**Address Randomization**

In order to reduce network bandwidth and ensure consistent successful delivery of the claim and proof transactions to
the finality layer, Pocket Core randomizes the transaction submission by address. The offset calculation is simple:

> \(Block Height _+_ Byte 1 of ValAddress\) _modulo_ Blocks Per Session

This algorithm ensures an efficient offset calculation and an equal distribution of automatic transactions.

## TXIndexer

Pocket Core uses a custom transaction indexer for optimal usage of resources. Tendermint's default TxIndexer paginates
and sorts transactions during time of query, resulting in a large consumption of resources. Tendermint's default
TxIndexer also relies on events to index transactions, and due to a bug in RC-0.5.X pocket core module events are all
concatenated together. This makes indexing a nightmare.

### **Custom Indexer**

Since LevelDB comparators order lexicongraphically, the implementation uses ELEN to encode numbers to ensure
alphanumerical ordering at insertion
time. [https://www.zanopha.com/docs/elen.pdf](https://www.zanopha.com/docs/elen.pdf) Since the keys are sorted
alphanumerically from the start, we don't have to:

* Load all results to memory
* Paginate and sort transactions after

This indexer inserts in sorted order so it can paginate and return based on the db iterator resulting in a significant
reduction in resource consumption.

The custom pocket core transaction indexer also reduces the scope of the Search\(\) functionality to optimize strictly
for the following use cases:

* BlockTxs \(Get transactions at a certain height\)
* AccountTxs \(Get transactions for a certain account \(sent and received\)\)

The custom pocket core transaction indexer also injects the message\_type into the struct to provide an easier method of
parsing the transactions. `json:"message_type"`

## Edit Stake

Rules:

1. Must be after 0.6.X upgrade height
2. You should ONLY be able to execute this transaction while status = staked
3. You can also execute this transaction while in jail
4. Only can modify certain parts of the structure \(see below\)

```text
type Application struct {
              Address                 sdk.Address      // SHOULD NOT CHANGE
              PublicKey               crypto.PublicKey // SHOULD NOT CHANGE
              Jailed                  bool             // SHOULD NOT CHANGE
              Status                  sdk.StakeStatus  // SHOULD NOT CHANGE
              Chains                  []string         // CAN CHANGE
              StakedTokens            sdk.BigInt       // CAN GO UP ONLY
              MaxRelays               sdk.BigInt       // CAN GO UP ONLY
              UnstakingCompletionTime time.Time        // SHOULD NOT CHANGE
          }


type Validator struct {
              Address                 sdk.Address      // SHOULD NOT CHANGE
              PublicKey               crypto.PublicKey // SHOULD NOT CHANGE
              Jailed                  bool             // SHOULD NOT CHANGE
              Status                  sdk.StakeStatus  // SHOULD NOT CHANGE
              Chains                  []string         // CAN CHANGE
              ServiceURL              string           // CAN CHANGE
              StakedTokens            sdk.BigInt       // CAN GO UP ONLY
              UnstakingCompletionTime time.Time        // SHOULD NOT CHANGE
          }
```

## Max Validators & Separation of Service and Consensus

The `Max_Validators` DAO param assigns a ceiling threshold on the number of _Tendermint Validators_. However, it does
not limit the number of _Servicer Nodes_ in the network. Max\_Validators caps _Tendermint Validators_ to help with the
current P2P bottlenecks that exist today. Less Validators = Less Consensus P2P traffic. Max\_Validators is not a static
floor. For instance, if the Max\_Validator threshold is exceeded the Validator with the lowest amount of _Stake_ out of
all the the Validators is removed as a Tendermint Validator and simply exists as a _Servicer_. These changes are likely
to happen in between blocks as new Validators stake in the network.

Summary:

* Less Validators = Less Consensus P2P traffic
* The number of Service Nodes is limited
* A Validator can join the Tendermint Validators set by staking more than the lowest staked Tendermint Validator which
  in turn changes the lowest state validator into just being a Service Node.

## Max\_Applications

Max\_Applications is a parameter assigning a ceiling threshold on the number of Applications able to stake in Pocket
Network.

The rules for Max\_Applications are simple:

* Applications cannot stake \(no matter what amount\) if this parameter is enabled and MaxApplications threshold is
  reached
* If an Application unstakes a slot will open up

