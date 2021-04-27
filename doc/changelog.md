## RC-0.6.2
- Implemented a Custom Transaction Indexer
- Return AllClaims if no address is passed to nodeClaims query
- No ABCI query during newTx() function in pocket core module
- Change StdSignature from Base64 to Hex in RPC
- Added config to drop events from account and blocktxs

## RC-0.6.1
- Fixed RelaysToTokensMultiplier bug
- Added utility CLI command to convert evidence from amino to protobuf
- Updated RPC spec to include stdTx
- Return Dispatch for certain failed relay codes to save a hop on client side
- Fix simulate relay to use basic auth
- Updated User Guide to use RC-0.6.0
- Added unsafe delete command to the keybase
- Added optional flag to bypass interactive prompt for commands that require passphrase input

## RC-0.6.0
- Security patch in merkle sum index
- Security patch for BlockHash
- Security patch for Proposer Selection
- Protobuf Encoding implemented
- Pocketcore module event handler fixed
- Max age of evidence enforced in blocks
- Deleted applications from state after unstake

## RC-0.5.2
- Delete local Relay/Challenge Evidence on Code 66 failures
- Log relay errors to nodes (don't just return to clients)
- Added configuration to pre-validate auto transactions
- Sending proofs/claims moved to EndBlock
- Load only Blockmeta for PrevCtx
- Added configurable cache PrevCtx, Validators, and Applications
- Don't broadcast claims/proofs if syncing
- Spread out claims/proofs between non-session blocks
- Added max claim age configuration for proof submission
- Reorganized non-consensus breaking code in Relay/Merkle Verify for efficiency before reads from state
- Configuration to remove ABCILogs
- Fixed (pseudo) memory leak in Tendermints RecvPacketMsg()
- Sessions only store addresses and not entire structs
- Only load bare minimum for relay processing
- Add order to AccountTxs query & blockTxsQuery RPC
- Reduce AccountTxsQuery & blockTxsQuery memory footprint
- Nondeterministic hash fix
- Code 89 Fix
- Evidence Seal Fix
- Fixes header.TotalTxs !=
- Fixes header.NumTxs !=
- Updating TM version and Version Number to BETA-0.5.2.3
- Upgraded AccountTxs and BlockTxs to use ReducedTxSearch
- Implemented Reduced TxSearch in Tendermint
- Moved IAVL from Tendermint to Pocket Core
- Call LazyLoadVersion/Store for all queries and PrevCtx()
- Reduced Tendermint P2P EnsurePeers actions to prevent leak
- Lowered P2P config to far more conservative numbers
- Updated FastSync to default to V1
- Exposed default leveldb options
- Switched to only go-leveldb for leak benchmarking/performance reasons
- Child process to run madvdontneed if not set
- Updated P2P configs
- fixed nil txIndexer bug (Tendermint now sets txindexer and blockstore)
- removed event type and used Tendermint's abci.Event

## RC-0.5.1
- Add terminal completions
- AllowDuplicateIP config default value is now true
- Added Relay Metrics Endpoint
- Ensured pocket-evidence is deleted in all lifecycles
- Claims are now not sent if below minimum proofs
- Improved BlockTxs and AccountTxs Endpoint Output with std.msg
- AAT now printing as json string and not bytes
- Enabled exporting of new genesis file command via CLI for Decentralized Reset
- Roll back blocks / Override App Hash in case of non-deterministic state
- Stop CLI when invalid character creating multi signature accounts
- Fixed "superfluous response.WriteHeader" log message showing on the logs
- Merged POSMint into core

## RC-0.5.0
- Fixed Incorrect Logging and misleading logging for max_signed_blocks
- Remove all receipts to alleviate large state size
- State now doesn't accept claims for under the given param threshold
- AAT now printing as json string and not bytes // not fixed
- Logging fix for unstaked Validator in Handle Validator Signature
- enhanced proof entropy and flush to disk
- Converted `Could not get sessionCtx` in auto send claim tx to an info log

## RC-0.4.3
- Update txIndexer synchronously during commit

## RC-0.4.2
- Update timeout configs for 15 minutes
- pocket-core version checking allows greater version strings and not just the exact version
- Export State doesn't export missed blocks
- Default Fees switched to be .01 POKT or 10,000 uPOKT
- Fee multiplier accounted for in autoTX
- Added memo argument to the send-tx CLI command
- Fixed simulate relay flag
- Fixed jailing lifecycle
- Fixed multisig cli params

## RC-0.4.1
- Update default config
- Fixed Slash() allows zero factor slash
- Fixed ValidateValidatorStaking checks wrong structure for isJailed
- Fixed Check Validator Status before jailing
- Fixed Force Unstaked Validators Never Get Removed
- Fixed Val unstaked continue if not found
- Fixed De-escape json from tx.Response
- Fixed Account-txs and block-txs CLI commands not working
- Fixed App unstake requires 4 arguments but only needs 3
- Fixed DAO Transfer & burn invalid signature


## RC-0.4.0
- Separated Validators By Network ID
- Remove EnsureExists from app & node modules, redundnacy with GetAccounts
- Close open connections on TMClients
- Removed all panics from source code
- Added Bloomfilter for efficient uniqueness checkint of proof objects
- Log on debug level upon failure deleting evidence
- Replace wealdtech/go-merkletree/crypto for golang.org/x/crypto
- Add RelaysToTokens as pamater for the nodes module
- removed return on loop from unstakeAllMaturedValidators
- Remove claim from world state after handled replay attack
- Use sdk.Error for app module keeper coinsFromUnstakedToStaked
- Updated minting scheme to mint on the tx handler level
- Fixed minting issue with unused fees
- Updated burning to burn on tx handler level
- Removed slash.go in apps module
- Ensured validator !isJailed in removeValidatorTokens
- Added max chains params to nodes and apps module
- Updated pocket-core to have 1 msg and 1 signature per tx
- Added RemoteCLIURL flag and config
- Removed round trip from tendermint
- Routed CLI through pocket RPC
- Change ClaimSubmissionWindow and BlocksPerSession call on getPseudorandomIndex to use the session context
- Changed getPseudorandomIndex add a new parameter for sessionCtx
- log mint errors
- Catch error writing json for RPC endpoints
- log mint errors
- Fixed denomination issue in CLI POKT to uPOKT
- Updated PopModel so no empty body `{}` is needed in RPC
- Added evidence type to proof message
- Added Node claims and claim queries *RPC*
- Removed chains prompt and defaults to no chains if none found
- Changed SendClaimTx to use sessionContext for supportedBlockchains
- Refactor ABCI BeginBlock for app module, avoid looping twice on mature applications
- Swapped ExecuteProof and SetReceipt order on handleProofMsg func
- Added prove to queryTx call in *RPC* and a flag in CLI
- Removed Proof Object From Relay Response in *RPC*
- Updated regex to allow only hex character.
- Added empty/nil check to claim struct on handleProofMsg
- Added multiple optimization for efficiency and consistency between validations
- Added basic auth to chains.json
- Added pagination to receipts query in pocket core *RPC*
- Removed .md and json RPC spec
- Removing orphaned jailed nodes after X amount of Blocks jailed.
- Fixed CORS requests on `/v1/client/dispatch`, `/v1/client/relay` and `/v1/client/challenge`.
- Split DAOTx transfer & burn to their own CLI commands to avoid user prone errors
- Added minimum proofs to pocketcore module params
- Split burn and transfer DAO CLI command, avoid error prone errors
- Cache flushes to the database periodically instead of per relay for efficiency
- Added dynamic fees for each message type
- Update RPC Docs with new unified params endpoint

## RC-0.3.0
- Added governance module from pocket-core
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
  - Updated pocket-core version
- Removed default genesis and seeds
- Prevent unstaking time reset if not equals 0 (for vesting)
- Relay Response now has the entire proof when sent back to the client
- Simplified slashing with static burn for validators and apps (non consensus burns)
- Payment for challenge tx
- Added export app command to cli
- Changed Struct used to generate RequestHash to remove empty proof object
- Changed public key field json marshalling
  - update pocket-core version
  - transitive change in Account JSON response for RPC
  - Updated rpc-spec to reflect change
- Updated pocket-core to address duplicated minting logs
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
- Add flag to to simulate relays, default false.

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
- Update pocket-core module to use sdk.Ctx interface
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
