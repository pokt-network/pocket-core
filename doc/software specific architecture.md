software specific architecture
========================
### Preface/Disclaimer
Pocket Core is the Official Golang implementation of the Pocket Network Protocol. 

*This document DOES NOT cover the architecture of the Pocket Network Protocol. Rather it only covers the software specific implementation details of Pocket Core and for the sake of clarity, orientation, and modularity, explicitly avoids covering any details that might duplicate the contents of the Pocket Network Protocol Specification.*

## Contents

1.  Cosmos SDK And Tendermint
2. Evidence Collection
3. Automatic Randomized Proof and Claim submissions

### Tendermint and Cosmos SDK

Pocket Core is technically an ABCI application built on top of a modified version of [Tendermint](https://tendermint.com/): a state machine replication engine built natively in Golang. Tendermint abstracts a BFT consensus algorithm that replicates Pocket Core's state machine logic among all nodes, creating the decentralized network. 

The decision to use Tendermint was made out of both necessity and convenience. Pocket Core *needs* a state machine replication module and Tendermint *conveniently* modularized just that. Though multiple **protocol level** modifications were made to Tendermint, this document will only cover the **implementation level** modifications: 
1. Pocket Core restricts the use of GoLevelDB only
2. Pocket Core exposes LevelDB options in configuration file
3. Pocket Core patches some minor P2P issues
4. Pocket Core patches some minor resource issues 

NOTE: *This document does not cover points 3 and 4, however it is covered in the [commit history](https://github.com/pokt-network/tendermint/commits/)*

Pocket Core restricts the use of GolevelDB for a few basic reasons.   
  
1) GoLevelDB [has proven](https://github.com/pokt-network/pocket-core/releases/tag/RC-0.5.2.9) to be the most efficient and reliable database for Tendermint
2) GolevelDB is not *always* cross compatible with ClevelDB and certainly not compatible with other non-LevelDB implementations
3) The effort of maintaining multiple *experimental* (as Tendermint refers to them) database implementations is quite high, and is currently not a priority for PNI. 
4) GolevelDB is able to be monitored and benchmarked with Golang Library tools like [PPROF](https://github.com/google/pprof) and a single database implementation is much easier to standardize resource consumption.

Pocket Core exposes LevelDB options in the Tendermint configuration file

1) LevelDB configurations are provided by the library to allow for users to perform application level modifications to the storage layer of the blockchain application.
2) Tendermint [often](https://github.com/cosmos/cosmos-sdk/issues/1394#issuecomment-402819672) suffers from 'Open File Limit' issues. GolevelDB provides a configuration for this

Though Pocket Core does rely on Tendermint for its replication / consensus module, Pocket  Core is often confused as a modified Cosmos SDK: *An SDK for building application specific blockchains.*

Transparently, the original idea was to use the Cosmos SDK with little modification to suit the Proof-Of-Stake needs of the Pocket Network Protocol. However, the PNI Engineering team quickly realized the Cosmos SDK was really not an SDK in the traditional sense, rather a somewhat loose forkable framework. 1K+ commits later, Pocket Core shares very little code with the Cosmos SDK other than its loose Module structure and the [KV store implementation](https://github.com/pokt-network/pocket-core/tree/staging/store) (both of which are planned to be overhauled in the future). 

### Relay Evidence Collection

Terminology: *Relay Evidence - Proof of a Relay Completed*

Pocket Network implementations must collect Relay Evidence in order to participate in the network's economic incentive mechanisms. 

Pocket Core's Relay Evidence collection is broken down into two major pieces:
1) A standard unordered golang slice (dynamic array) to store the evidence 
2) An adjacent [bloom filter](https://www.eecs.harvard.edu/~michaelm/NEWWORK/postscripts/BloomFilterSurvey.pdf) for time/space efficient membership checks

**Evidence Store:** 

For speed and quality of service reasons, Pocket Core chooses to store the evidence undordered (as it comes in from the client). The evidence is then ordered using a quicksort algorithm at time of the claim/proof submission.

Pocket Core also dynamically grows the storage array (slice) with the current assumption that most Applications will not use close to their maximum relay capacity.

*Alternatives:* 

Pocket Core could store the evidence ordered with the tradeoff of quality of service and speed to the client. This comes with a compute benefit at the time of the Merkle Tree generation (proof/claim submisison)

Pocket Core could fix the array to the maximum size for each application, because the relay capacity is known at the time of array creation. This comes with a tradeoff in higher memory consumption but for a decrease in compute.

**Bloom Filter**

Duplicate Relay Evidence can easlity result in a *code 66* invalid Merkle Proof Protocol error, or worse, trigger an economic Replay Attack defense mechanism. 

In order to prevent duplicate relay evidence, the protocol implementation should check for membership of the Relay Evidence before accepting anymore incoming requests.

Pocket Core uses a bit array Bloom Filter implementation to provide an efficient structure for quick membership checks. Due to the properties of a Bloom Filter, Pocket Core will never accept duplicate evidence and is able to check for membership with O(1) time

*Alternatives:* 

Use a standard Golang Map to store the Evidence and just use that structure to check for membership. Also O(1) time. The tradeoff here is more memory because the membership is stored as a hash vs bits.

### Proof and Claim Submissions

Inorder to participate in the Pocket Network economic incentive mechanisms, proof and claim messages must be submitted to claim and prove the work (relays) done.

The protocol, however, does not specifiy how or when to execute the proof/claim transactions, and leaves it up to the implementation.

Pocket Core's Proof and Claim submission algorithm is summarized in two unique points:
1. Automatic Submission
2. Address Randomaztion

**Automatic Submission**

Contrary to most blockchain client's, Pocket Core attempts to automatically submit claim and proof transactions if there is Relay Evidence stored in the Evidence database and the validator is staked for work. Since the private key of any Tendermint Validator is already available and exposed for block signing, Pocket Core uses the PK to sign the Claim/Proof transactions.

*Alternative*

Pocket Core can ask the users to manually submit the claim and proof transaction or offload the work to a third party application. The tradeoff is both conveinence and robustness of the application

**Address Randomization**

In order to reduce network bandwidth and ensure consistent successful delivery of the claim and proof transactions to the finality layer, Pocket Core randomizes the transaction submission by address. The offset calculation is simple:
> (Block Height *+* Byte 1 of ValAddress) *modulo* Blocks Per Session

This algorithm ensures an efficient offset caluclation and an equal distribution of automatic transactions.

### TXIndexer
Pocket Core uses a custom transaction indexer for optimal usage of resources. Tendermint's default TxIndexer paginates
and sorts transactions during time of query, resulting in a large consumption of resources. Tendermint's default TxIndexer
also relies on events to index transactions, and due to a bug in RC-0.5.X pocket core module events are all concatenated
together. This makes indexing a nightmare.

**Custom Indexer**
Since LevelDB comparators order lexicongraphically, the implementation uses ELEN to encode numbers to ensure alphanumerical
ordering at insertion time. https://www.zanopha.com/docs/elen.pdf
Since the keys are sorted alphanumerically from the start, we don't have to:
    - Load all results to memory 
    - Paginate and sort transactions after
This indexer inserts in sorted order so it can paginate and return based on the db iterator resulting in a significant 
reduction in resource consumption

The custom pocket core transaction indexer also reduces the scope of the Search() functionality to optimize strictly for 
the following use cases:
    - BlockTxs (Get transactions at a certain height)
    - AccountTxs (Get transactions for a certain account (sent and received))
    
The custom pocket core transaction indexer also injects the message_type into the struct to provide an easier method of 
parsing the transactions. `json:"message_type"`




