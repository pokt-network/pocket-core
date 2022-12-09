# Claims, Proofs, and Merkle Sum Index Trees

## Abstract

V0 of Pocket Network utilizes a custom reward protocol that enables a succinct and probabilistic proof of work completed. Specifically, this procedure allows staked Service Nodes to be retroactively rewarded POKT in exchange for a proof of Relays executed during a Session via a "commit-reveal" scheme.

The document below assumes an understanding of various Pocket Network-related concepts, including but not limited to Relays, Sessions, Servicers, etc. Please visit [docs.pokt.network](https://docs.pokt.network) prior to reading this document if you are not familiar with this terminology.

## High Level Process

The reward protocol is executed in a **3 step** process:

\* *Node creates a Merkle Tree locally* \* 

**1) Generate Merkle Sum Index (MSI) Tree using Relay Proofs**

\* *A commit-reveal cycle begins by interacting with the state of the protocol* \* 

**2) Claim Transaction (Node claims the volume of Relays completed)**

**3) Prove Transaction (Node proves Claim using MSI Tree)**

```mermaid
graph TD
    subgraph "Reward Protocol High Level Flow"
        A(Session Ends) --> |Aggregate Relay Proofs|B
        B[Generate MSI Tree] --> |Use MSI Root|C
        C[Send Claim Message] --> |Use MSI Proof|D
        D[Send Proof Message] --> E(Rewarded)
    end
```

## Generate Merkle Sum Index Tree

At a high level, the generation of the MSI tree is simply aggregating the Relay Proofs from the applicable Session and using each Proof as a leaf to build the Merkle Tree.

By building a tree out of the Relay Proofs, the protocol is able to verify a single leaf of the tree with a succinct < O(log(n)) > Merkle Proof.

![merkle_tree](https://user-images.githubusercontent.com/18633790/204642445-08b558d5-9a17-4ec3-b073-8ccb865f25e7.png)


However, it is essential to note the nuances of the MSI tree and why this specific flavour of Merkle Tree is chosen.

The MSI Tree is based on [Plasma Core's Merkle Sum Tree](https://plasma-core.readthedocs.io/en/latest/specs/sum-tree.html).

Specifically, the property that is most interesting from Plasma Core's MST is:

> [...] it allows us to calculate a specific range which a branch describes, and be 100% confident that no other valid branches exist which overlap that range.

This non-overlap property of branches in MSTs is foundational to preventing Replay Attacks of completed Relay Proofs in Pocket Network.

*NOTE:  The only difference between a Merkle Sum Tree and a Merkle Sum Index Tree is that the index of the leaf is utilized in the `ParentHash` function to add an extra layer of security around the swapping unique relays within a Merkle Tree.*

A replay attack is simply duplicating completed Relay Proofs inside a Merkle Tree to falsify additional Relays completed.

![replay](https://user-images.githubusercontent.com/18633790/204642475-24dc64e9-8d01-4837-b0a7-bd975295c003.png)


By ordering the leafs lexicographically by their hash and using the `Hash(left sibling)+1` as the `lower range` and the hash of the leaf as the `upper range`, Pocket Network is able to build a non-overlapping leaf - Merkle Tree.

For example, assume that these are ordered lexicographically: `[hash(A), hash(B), hash(C), hash(D)]`

- Range of A: `[0 + 1, hash(A)]`
- Range of B: `[hash(A) + 1, hash(B)]`
- Range of C: `[hash(B) + 1, hash(C)]`
- Range of D: `[hash(C) + 1, hash(D)]`


It is **very important** for the reader to note two related parts to Pocket Network's Servicing Protocol:
- *Each Relay Proof is signed by the Application Client using the `SignBytes()` which is used to generate the `Hash(leafSignBytes)`.*
- *Service Nodes reject all `identical Relays` with a `Code: 37` Duplicate Proof Error*

Due to the signature scheme of Relay Servicing, a Service Node is unable to both falsify a Relay Proof in a Merkle Tree and have that Leaf be verifiable via the signature - completing the Replay Attack protection.

The Merkle Sum Index Tree enables:
- Enforceably Unique Relay Proof Leafs
- Succinct Proof of a single Relay Proof Leaf
- Verifiable integrity of the indices of the Leafs using a Merkle Root

## Claim Transaction

After an active Session finishes and before the closing of the `Claim Submission Window`, a **Claim Transaction** is sent by the Service Nodes to be stored on-chain to **claim** work completed.

```go
type MsgClaim interface {  
  SessionHeader()     // link the work to a specific prior session
  MerkleRoot()        // lock in the integrity of the Leafs indices
  TotalProofs()       // lock in the range ( 0-TotalProofs ) 
  FromAddress()       // sender of the claim transaction (i.e. the Service Node)
  EvidenceType()      // relay (or challenge which is outside the scope of this document)
  ExpirationHeight()  // the expiration height of the claim, prior to which the proof must be submitted (dictated by on-chain gov params)
}
```

```mermaid
sequenceDiagram
    participant ClaimTx as ClaimTx
    participant WorldState as World State
    participant MSITree as MSI Tree
    participant RelayProofs as Relay Proofs
    
    ClaimTx->>+WorldState: IsSessionEnded()
    WorldState-->>-ClaimTx: SessionHeader
    ClaimTx->>+MSITree: GenerateMerkleRoot()
    MSITree->>+RelayProofs: GenerateTree(sessionHeader)
    RelayProofs-->>-MSITree: MSI Tree
    MSITree-->>-ClaimTx: Root & TotalProofs
    ClaimTx->>WorldState: SubmitClaimTx()
```

The Claim Transaction enables the protocol to:
- Lock in the Merkle Root from the MSI Tree (a single commitment for all the relays the Service Node performed in a session)
- Establish a **range** or max indices to select from (ranging from 0 to the number of relays performed by the Service Node during a session)

## Proof Transaction

The Claim Transaction sets up the protocol to be able to challenge the Service Node for a succinct proof of the work completed.

The network state now contains the Claim's Merkle Root and the maximum index (i.e. max relay number) it may select from.

Since the Merkle Root preserves the integrity of the Relay Proof indices and the MSI Tree ensures unique Relay Proofs as leafs, the protocol is able to use a single randomly selected leaf from the tree and use that as proof of the entire set as long as there exists a MSI Branch from the selected leaf to the given Merkle Root.

```mermaid
sequenceDiagram
    ProofTx->>+WorldState: IsClaimSubmissionWindowEnded()
    WorldState->>+PseudorandomSelection: Select(Blockheader, SessionHeader, TotalProofs)
    PseudorandomSelection-->>-WorldState: SelectedIndex
    WorldState-->>-ProofTx: SelectedIndex
    ProofTx->>+MSITree: GenerateMerkleProof(SelectedIndex)
    MSITree->>+RelayProofs: GetLeaf(SelectedIndex)
    RelayProofs-->>-MSITree: Leaf
    MSITree-->>-ProofTx: Leaf & MerkleProof
    ProofTx->>+WorldState: SubmitProofTx()
    WorldState->>WorldState: ValidateProofTx()
    WorldState-->>-ProofTx: IfValid() -> Mint POKT
```

### Pseudorandom Selection Of Proof Index

Similar to the Session Protocol, the Pseudorandom Selection algorithm uses the BlockHash of `SessionEndHeight`+`ClaimSubmissionWindow` block as seed data for the Pseudorandom Selection.

```go
headerHash := sessionHeader.HashString()  
pseudoGenerator := pseudorandomGenerator{blockHash, headerHash}
r, err := json.Marshal(pseudoGenerator)  
...
PseudorandomSelection(totalRelays, Hash(r))
```

To increase security and to tailor unique index selections for each Node, the seed data is also calculated from the SessionHeader.

```go
func PseudorandomSelection(max sdk.BigInt, hash []byte) (index sdk.BigInt) {  
   // first 8 bytes of the seedHash
  intHash := sdk.NewIntFromBigInt(new(big.Int).SetBytes(hash[:8]))  
   // modulo the selection against the maximum of the range
  return intHash.Mod(max)  
}
```

*NOTE: The algorithm must be **pseudo**random because it must also be deterministic for the network to have Consensus on the agreed upon index.*

### Proof of the Claim

Once the Pseudorandom index is revealed to the network, the final step is for the Service Node to submit a `Proof Transaction` to the network to satisfy:
> the protocol is able to use a single randomly selected leaf from the tree and use that as proof of the entire set **as long as there exists a MSI Branch from the selected leaf to the given Merkle Root**

Using the previously generated MSI Tree, the Service Node is able to generate a MSI Branch from the Pseudorandomly Selected Leaf to the Merkle Root.

Including the Leaf and the MSI Branch, the Service Node is able to complete the lifecycle by submitting a **Proof Transaction** to the network.

```go
type MsgProof interface {  
   MerkleBranch()  // the MSI Branch
   Leaf()          // the pseudrandomly selected leaf
   EvidenceType()  // relay (or challenge which is outside the scope of this document))
}
```

### Reward

Like any Transaction, the proof is validated by all who process blocks.

In order to determine if a `ProofTx` is valid, several Validations must occur.

```go
func ValidateProofTx(proof MsgProof) error {  
  // find corresponding claim from world state
  ...
  // validate session exists by comparing the session header against the world state
  ...
  // generate pseudorandom index locally (using deterministic seeds) and validate against leaf provided
  ...
  // validate the merkle proof (common cryptographic approach) 
  ...
  // validate the leaf
}
```
If a Proof Transaction is `Code:0` (valid) then POKT reward is minted to the Service Node based on the current Governance Parameters.

## Additional References & Resources

**Documentation:**

* Pocket docs: https://docs.pokt.network/home/v0/protocol/servicing#claim-proof-lifecycle
* Twitter Thread: https://twitter.com/o_rourke/status/1263847357122326530
* Plasma core Merkle Sum Tree: https://plasma-core.readthedocs.io/en/latest/specs/sum-tree.html

**Source code references:**

* Merkle: https://github.com/pokt-network/pocket-core/blob/staging/x/pocketcore/types/merkle.go
* Claim: https://github.com/pokt-network/pocket-core/blob/staging/x/pocketcore/keeper/claim.go
* Proof: https://github.com/pokt-network/pocket-core/blob/staging/x/pocketcore/keeper/proof.go

