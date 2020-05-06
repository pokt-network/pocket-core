## RC-0.4.0
- Separated Validators By Network ID
- Remove EnsureExists from app & node modules, redundnacy with GetAccounts
- Close open connections on TMClients
- log on debug level upon failure deleting evidence
- replace wealdtech/go-merkletree/crypto for golang.org/x/crypto

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
- Prevent unstaking time reset if not equals 0 (for vesting)
- Relay Response now has the entire proof when sent back to the client
- Simplified slashing with static burn for validators and apps (non consensus burns)
- Payment for challenge tx
- Added export app command to cli
- Changed Struct used to generate RequestHash to remove empty proof object
- Changed public key field json marshalling
  - update POSMint version
  - transitive change in Account JSON response for RPC
  - Updated rpc-spec to reflect change
- Updated POSMint to address duplicated minting logs
- Add pagination to Application Queries for RPC & CLI
- Add pagination to Nodes Queries for RPC & CLI
- Converted NonNative Chains to 128 bit encryption (MD5)
- Renamed sessionFrequency to blocksPerSession
- Changed `Proof` field from relayResponse to `proof`
- Changed `Proofs` in genesis struct to `Receipts`
- Changed HostedBlockchain field from `addr` to `id`
- Modified the key generation for receipts and claims by using the header hash
- Added headers back to payload (affects request hash!)
- New Dispatch Formula See Spec
- Added configurability for Pocket and Tendermint in a config.json file
- Patch for Fixed April 17, 2020 consensus failure (0 power consensus failure)
- Updated RPC call for Querying Validators and Apps (See Opts in RPC Spec)
- Updated chains.json to be a slice instead of a mapping
- Chanded defer.body.close() position until after error checking in pocketcore/types/service.go executeHttpRequest
- Added protocol level enforcement of network identifier format
- Added protocol level enforcement of service url `https`
- Updated the import-armored and export commands for better UX keeping backward compatibility with pre RC 0.3.0 .
- Fixed issue querying apps/validators was ignoring the blockchain key
- Added protocol level enforcement of service url `https|http`
- Added `/v1/query/blocktxs` and `/v1/query/acounttxs` to get the list of transactions in a block or sent/received by an account respectively.
- Added `query account-txs` and `query block-txs` to the CLI matching the above mentioned endpoints.
- Fixed `getTMClient()` function to avoid Tendermint opening files every time a Relay/Dispatch.

## RC-0.2.4
- Removed trailing slash when performing a relay when there's no path present in the relay payload.

## RC-0.2.3
- Changed Tendermint's dbbackend to cleveldb after multiple validator crashes on April's 20th load test.

## RC-0.2.2
- Fixed April 17, 2020 consensus failure (0 power consensus failure)
Seed patch to allow more incoming connections
- Changed querySupplyResponse to all string to prevent overflow

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
- Add off chain relay RPC for testing purposes, wont' create proof & does not affect validator
- Add flag to to simulate relays, default false.
