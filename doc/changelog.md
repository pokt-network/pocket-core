## RC-0.3.0
- Added governance module from posmint
  - Multisignature public keys / tx building to cli
  - Governance level protocol upgrade signalling
  - Access control list for params
  - Ability to edit params TX
  - Introduced DAO-Owner
  - Ability to send and receive from DAO
- Added persistence to Sessions and Evidence through
  - LRU cache sessions/evidence
  - Sessions/evidence level-db
- Added Start without passphrase on pocket core
  - sign msgs with private key or keybase
  - removed password from keeper
  - Updated POSMint version
- Removed default genesis and seeds

## RC-0.2.1
- Add version command to CLI
- Disallow double sign on invalid operations, disallow consensus breaks

## RC-0.2.0
- Renamed RelayProof to Proof (in JSON)
- Renamed Invoice (memory) to Evidence
- Renamed StoredInvoice (blockchain persisted) to Receipt
- Renamed ProofWaitingPeriod to ClaimSubmissionWindow
- Changed RPC and from `node-proof` to `node-receipt`
- Update posmint module to use sdk.Ctx interface
- Fix `pseudorandomGenerator` unexported properties would return empty json
- Evidence now holds proof interface to allow for challenge proofs
- Added Relay Request Hash (Hash of payload + meta) to RelayProof object
- Added Block to Dispatch Request
- Added Relay Meta field to relay request
- Added Challenge Functionality
- Added Challenge Request to RPC
- Changed dispatch response to an actual structure and not just a session
- Added block height to dispatch response
- Removed all MustGetPrevCtx and used PrevCtx for panic safety
- Changed receipt structure (added evidence type)
- Change `querySupplyResponse` struct to use `totalStaked`, `totalUnstaked` & `Total` as `*big.Int` due to memory overflow
- Added RPC SPEC doc in yaml and json with swagger support
