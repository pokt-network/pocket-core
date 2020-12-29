# Pocket Core CLI Interface Specification
## Version RC-0.3.0

### Overview
This document serves as a specification for the Command Line Interface of the Pocket Core application. There's no protocol verification for these commands because they map closely to protocol functions.

### Namespaces
The CLI will contain multiple namespaces listed below:

- Default Namespace: These functions will be called when the namespace is blank
- Accounts: Contains all the calls pertinent to accounts and their local storage.
- Nodes: Contains all the functions for Node upkeep.
- Apps: Contains all the functions for app upkeep.
- Query: All queries to the world state are contained in this call.

### CLI Functions Format
Each CLI Function will be in the following format:

- Binary Name: The name of the binary for Pocket Core, for example: `pocket`
- Namespace: The namespace of the function, or blank for the default namespace: `accounts`
- Function Name: The name of the actual function to be called: `create`
- (Optional): Space separated function arguments, e.g.: `pocket accounts create <passphrase>`

### Default Namespace
The default namespace contains functions that are pertinent to the execution of the Pocket Node.

- `pocket start <datadir>`
> Starts the Pocket Node, picks up the config from the assigned `<datadir>`.
>
> Arguments:
> - `<datadir>`: The data directory where the configuration files for this node are specified.

- `pocket reset`
> Reset the Pocket node.
> Deletes the following files / folders:
> - .pocket/data
> - priv_val_key
> - priv_val_state
> - node_keys
>

### Accounts Namespace Functions
The `accounts` namespace handles all account related interactions, from creating and deleting accounts, to importing and exporting accounts.

- `pocket accounts list`
> Lists all the account addresses stored in the keybase.
> Example output:
```
0xb3746D30F2A579a2efe7F2F6E8E06277a78054C1
0xab514F27e98DE7E3ecE3789b511dA955C3F09Bbc
```

- `pocket accounts show <address>`
> Lists an account address and public key.
>
> Arguments:
> - `<address>`: The address to be fetched.
>
> Example output:
```
Address: 0x.....
Public Key: 0x....
```

- `pocket accounts delete <address>`
> Deletes an account from the keybase. Will prompt the user for the account passphrase
>
> Arguments:
> - `<address>`: The address to be deleted.
> Example output:
```
KeyPair 0x... deleted successfully.
```

- `pocket accounts update-passphrase <address>`
> Updates the passphrase for the indicated account. Will prompt the user for the current account passphrase and the new account passphrase.
>
> Arguments:
> - `<address>`: The address to be deleted.
> Example output:
```
KeyPair 0x... passphrase updated successfully.
```

- `pocket accounts sign <address> <msg>`
> Signs the specified `<msg>` using the specified `<address>` account credentials. Will prompt the user for the account passphrase.
>
> Arguments:
> - `<address>`: The address to be deleted.
> - `<msg>`: The message to be signed in hex string format.
> Example output:
```
Original Message: 0x...
Signature: 0x...
```

- `pocket accounts create`
> Creates and persists a new account in the Keybase. Will prompt the user for a [BIP-0039](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki) password for the generated mnemonic and for a passphrase to encrypt the generated keypair.
>
> Example output:
```
Account generated successfully.
Address: 0x....
```

- `pocket accounts import <mnemonic>`
> Imports an account using the provided `<mnemonic>`. Will prompt the user for a [BIP-0039](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki) password for the imported mnemonic and for a passphrase to encrypt the generated keypair.
>
> Arguments:
> - `<mnemonic>`: The mnemonic of the account to be imported.
> Example output:
```
Account imported successfully.
Address: 0x....
```

- `pocket accounts import-armored <armor>`
> Imports an account using the Encrypted ASCII armored `<armor>` string. Will prompt the user for a decryption passphrase of the `<armor>` string and for an encryption passphrase to store in the Keybase.
>
> Arguments:
> - `<armor>`: The encrypted encoded private key to be imported.
> Example output:
```
Account imported successfully.
Address: 0x....
```

- `pocket accounts export <address>`
> Exports the account with `<address>`, encrypted and ASCII armored. Will prompt the user for the account passphrase and an encryption passphrase for the exported account.
>
> Arguments:
> - `<address>`: The address of the account to be exported.
> Example output:
```
Exported account: <armored string>
```

- `pocket accounts export-raw <address>`
> Exports the raw private key in hex format. Will prompt the user for the account passphrase. ***NOTE***: THIS METHOD IS NOT RECOMMENDED FOR SECURITY REASONS, USE AT YOUR RISK.*
>
> Arguments:
> - `<address>`: The address of the account to be exported.
> Example output:
```
Exported account: 0x...
```

- `pocket accounts send-tx <fromAddr> <toAddr> <amount>`
> Sends `<amount>` POKT `<fromAddr>` to `<toAddr>`. Prompts the user for `<fromAddr>` account passphrase.
>
> Arguments:
> - `<fromAddr>`: The address of the sender.
> - `<toAddr>`: The address of the receiver.
> - `<amount>`: The amount of POKT to be sent.
> - `<chainID>`: The pocket chain identifier
> Example output:

> Flags
> - `prove`: Get a proof of the transaction
```
Transaction submitted with hash: <Transaction Hash>
```

### Node Namespace
Functions for Node management.

- `pocket node stake <fromAddr> <amount> <chains> <serviceURI>`
> Stakes the Node into the network, making it available for service. Prompts the user for the `<fromAddr>` account passphrase.
>
> Arguments:
> - `<fromAddr>`: The address of the sender.
> - `<amount>`: The amount of POKT to stake. Must be higher than the current minimum amount of Node Stake parameter.
> - `<chains>`: A comma separated list of chain Network Identifiers.
> - `<serviceURI>`: The Service URI Applications will use to communicate with Nodes for Relays.
> - `<chainID>`: The pocket chain identifier
> Example output:
```
Transaction submitted with hash: <Transaction Hash>
```

- `pocket node unstake <fromAddr>`
> Unstakes a Node from the network, changing its status to `Unstaking`. Prompts the user for the `<fromAddr>` account passphrase.
>
> Arguments:
> - `<fromAddr>`: The address of the sender.
> - `<chainID>`: The pocket chain identifier
> Example output:
```
Transaction submitted with hash: <Transaction Hash>
```

- `pocket node unjail <fromAddr>`
> Unjails a Node from the network, allowing it to participate in service and consensus again. Prompts the user for the `<fromAddr>` account passphrase.
>
> Arguments:
> - `<fromAddr>`: The address of the sender.
> - `<chainID>`: The pocket chain identifier
> Example output:
```
Transaction submitted with hash: <Transaction Hash>
```

### Pocket App Namespace
Functions for Application management.

- `pocket app stake <fromAddr> <amount> <chains>`
> Stakes the Application into the network, making it available to receive service. Prompts the user for the `<fromAddr>` account passphrase.
>
> Arguments:
> - `<fromAddr>`: The address of the sender.
> - `<amount>`: The amount of POKT to stake. Must be higher than the current minimum amount of Application Stake parameter.
> - `<chains>`: A comma separated list of chain Network Identifiers.
> - `<chainID>`: The pocket chain identifier
> Example output:
```
Transaction submitted with hash: <Transaction Hash>
```

- `pocket app unstake <fromAddr>`
> Unstakes an Application from the network, changing its status to `Unstaking`. Prompts the user for the `<fromAddr>` account passphrase.
>
> Arguments:
> - `<fromAddr>`: The address of the sender.
> - `<chainID>`: The pocket chain identifier
> Example output:
```
Transaction submitted with hash: <Transaction Hash>
```

- `pocket app create-aat <appAddr> <clientPubKey>`
> Creates a signed application authentication token (version `0.0.1` of the AAT spec), that can be embedded into application software for Relay servicing. Will prompt the user for the `<appAddr>` account passphrase. Read the Application Authentication Token documentation [here](application-auth-token.md). ***NOTE***: USE THIS METHOD AT YOUR OWN RISK. READ THE APPLICATION SECURITY GUIDELINES TO UNDERSTAND WHAT'S THE RECOMMENDED AAT CONFIGURATION FOR YOUR APPLICATION:
>
> Arguments:
> - `<appAddr>`: The address of the Application account to use to produce this AAT.
> - `<clientPubKey>`: The account public key of the client that will be signing and sending Relays sent to the Pocket Network.
> Example output:
```json
{
    "version": "0.0.1",
    "applicationPublicKey": "0x...",
    "clientPublicKey": "0x...",
    "signature": "0x..."
}
```

### Pocket Util Namespace
Generic utility functions for diverse use cases.

- `pocket util generate-chains`
> Generate the chains.json file for network identifiers.
>
> Example output:
```
Enter the ID of the network identifier:
>0001
Enter the URL of the network identifier:
>https://ethnode.test.com:8085
Would you like to enter another network identifier? (y/n)
>n
chains.json contains:

0001 @ https://ethnode.test.com:8085
If incorrect: please remove the chains.json with the delete-chains command

```

- `pocket util delete-chains`
> Delete the chains.json file for network identifiers.
>
> Example output:
```
successfully deleted chains.json
```

- `pocket util decode-tx <tx>`
> Decodes a given transaction encoded in Amino base64 bytes
>
> Arguments:
> - `<tx>`: the transaction amino encoded bytes.
> Example output:
```
% pocket util decode-tx qgLbCxcNCp0Bq4P6fApLCkA3ZWFjZWFjZTYwNzY1YzhiYjU0NDAzOGUxNGRjOGMyNjQ1NWRmODJmNTVmOGVkZDc1M2EwNDU5ZmY4MzYxZmViEgQwMDIxGP1wEi8KIEd86o3r3PIS6aK3CW+8L3E9JZMEHFdM1kMmy7XmuSQ/EgsQ8YrPpKGm95f/ARieAiIUjDp8K56yjpfHbsHBoLReW9EfapcoARIOCgV1cG9rdBIFMTAwMDAaaQolnVRHdCAO6zUJvs6taFLJzycYSzl2lPHXTYkxOnru2wG+T5y3PxJAckq7juFqII9kg/QPK2JmnLYNUthqZXNbEEQ5Zb/Jk/yqA2kwKUKS9yAZMPX8anDHj5Ycrtkw+LWnyha7aKFFBCiFvpiZ3YOT2JQB
Type:           claim
Msg:            {{7eaceace60765c8bb544038e14dc8c26455df82f55f8edd753a0459ff8361feb 0021 14461} {[71 124 234 141 235 220 242 18 233 162 183 9 111 188 47 113 61 37 147 4 28 87 76 214 67 38 203 181 230 185 36 63] {0 18388159010740356465}} 286 8C3A7C2B9EB28E97C76EC1C1A0B45E5BD11F6A97 1 0}
Fee:            10000upokt
Entropy:        -7732596869214888187
Memo:
Signer          8c3a7c2b9eb28e97c76ec1c1a0b45e5bd11f6a97
Sig:            0eeb3509becead6852c9cf27184b397694f1d74d89313a7aeedb01be4f9cb73f

```

- `pocket export-genesis-for-reset <height> <newChainID>`
> In the event of a network reset, this will export a genesis file based on the previous state
>
> Arguments:
> - `<height>`: the height to export.
> - `<newChainID>`: the chainID to use for exporting.
	> Example output:
```json
{
	"app_hash": "",
	"app_state": {
		"application": {
			"applications": [],
			"exported": true,
			"params": {
				"app_stake_minimum": "1000000",
				"base_relays_per_pokt": "167",
				"max_applications": "9223372036854775807",
				"maximum_chains": "15",
				"participation_rate_on": false,
				"stability_adjustment": "0",
				"unstaking_time": "1814000000000000"
			}
		},
		"auth": {
			"accounts": [
				{}
			]
		}
	}
}
...
```
- `pocket unsafe-rollback <height>`
> Rollbacks the blockchain, the state, and app to a previous height
>
> Arguments:
> - `<height>`: the height you want to rollback to.
>
> Flags
> - `blocks`: rollback block store and state
>

- `pocket completion [bash|zsh|fish|powershell]`
> Generate completion script for the specified shell
>
> Arguments:
> - `<shell>`: the shell you currently use. Supported options: **bash / zsh / fish / powershell**
>

- `pocket print-configs`
> Prints Default config.json to console.
>
> Example output:
```json
{
	"tendermint_config": {
		"RootDir": "/Users/admin/.pocket",
		"ProxyApp": "tcp://127.0.0.1:26658",
		"Moniker": "ultima.local",
		"FastSyncMode": true,
		"DBBackend": "goleveldb",
		"LevelDBOptions": {
			"block_cache_capacity": 83886,
			"block_cache_evict_removed": false,
			"block_size": 4096,
			"disable_buffer_pool": true,
			"open_files_cache_capacity": -1,
			"write_buffer": 838860
		},
		"DBPath": "data",
		"LogLevel": "*:info, *:error",
		"LogFormat": "plain",
		"Genesis": "config/genesis.json",
		"PrivValidatorKey": "priv_val_key.json",
		"PrivValidatorState": "priv_val_state.json",
		"PrivValidatorListenAddr": "",
		"NodeKey": "node_key.json",
		"ABCI": "socket",
		"ProfListenAddress": "",
		"FilterPeers": false,
		"RPC": {
			"RootDir": "/Users/admin/.pocket",
			"ListenAddress": "tcp://127.0.0.1:26657",
			"CORSAllowedOrigins": [],
			"CORSAllowedMethods": [
				"HEAD",
				"GET",
				"POST"
			],
			"CORSAllowedHeaders": [
				"Origin",
				"Accept",
				"Content-Type",
				"X-Requested-With",
				"X-Server-Time"
			],
			"GRPCListenAddress": "",
			"GRPCMaxOpenConnections": 2500,
			"Unsafe": false,
			"MaxOpenConnections": 2500,
			"MaxSubscriptionClients": 100,
			"MaxSubscriptionsPerClient": 5,
			"TimeoutBroadcastTxCommit": 10000000000,
			"MaxBodyBytes": 1000000,
			"MaxHeaderBytes": 1048576,
			"TLSCertFile": "",
			"TLSKeyFile": ""
		},
		"P2P": {
			"RootDir": "/Users/admin/.pocket",
			"ListenAddress": "tcp://0.0.0.0:26656",
			"ExternalAddress": "",
			"Seeds": "",
			"PersistentPeers": "",
			"UPNP": false,
			"AddrBook": "config/addrbook.json",
			"AddrBookStrict": false,
			"MaxNumInboundPeers": 10,
			"MaxNumOutboundPeers": 10,
			"UnconditionalPeerIDs": "",
			"PersistentPeersMaxDialPeriod": 0,
			"FlushThrottleTimeout": 100000000,
			"MaxPacketMsgPayloadSize": 1024,
			"SendRate": 5120000,
			"RecvRate": 5120000,
			"PexReactor": true,
			"SeedMode": false,
			"PrivatePeerIDs": "",
			"AllowDuplicateIP": true,
			"HandshakeTimeout": 20000000000,
			"DialTimeout": 3000000000,
			"TestDialFail": false,
			"TestFuzz": false,
			"TestFuzzConfig": {
				"Mode": 0,
				"MaxDelay": 3000000000,
				"ProbDropRW": 0.2,
				"ProbDropConn": 0,
				"ProbSleep": 0
			}
		},
		"Mempool": {
			"RootDir": "/Users/admin/.pocket",
			"Recheck": true,
			"Broadcast": true,
			"WalPath": "",
			"Size": 9000,
			"MaxTxsBytes": 1073741824,
			"CacheSize": 9000,
			"MaxTxBytes": 1048576
		},
		"FastSync": {
			"Version": "v1"
		},
		"Consensus": {
			"RootDir": "/Users/admin/.pocket",
			"WalPath": "data/cs.wal/wal",
			"TimeoutPropose": 60000000000,
			"TimeoutProposeDelta": 10000000000,
			"TimeoutPrevote": 60000000000,
			"TimeoutPrevoteDelta": 10000000000,
			"TimeoutPrecommit": 60000000000,
			"TimeoutPrecommitDelta": 10000000000,
			"TimeoutCommit": 900000000000,
			"SkipTimeoutCommit": false,
			"CreateEmptyBlocks": true,
			"CreateEmptyBlocksInterval": 900000000000,
			"PeerGossipSleepDuration": 100000000000,
			"PeerQueryMaj23SleepDuration": 200000000000
		},
		"TxIndex": {
			"Indexer": "kv",
			"IndexKeys": "tx.hash,tx.height,message.sender,transfer.recipient",
			"IndexAllKeys": false
		},
		"Instrumentation": {
			"Prometheus": false,
			"PrometheusListenAddr": ":26660",
			"MaxOpenConnections": 3,
			"Namespace": "tendermint"
		}
	},
	"pocket_config": {
		"data_dir": "/Users/admin/.pocket",
		"genesis_file": "genesis.json",
		"chains_name": "chains.json",
		"session_db_name": "session",
		"evidence_db_name": "pocket_evidence",
		"tendermint_uri": "tcp://localhost:26657",
		"keybase_name": "pocket-keybase",
		"rpc_port": "8081",
		"client_block_sync_allowance": 10,
		"max_evidence_cache_entries": 500,
		"max_session_cache_entries": 500,
		"json_sort_relay_responses": true,
		"remote_cli_url": "http://localhost:8081",
		"user_agent": "",
		"validator_cache_size": 100,
		"application_cache_size": 100,
		"rpc_timeout": 3000,
		"pocket_prometheus_port": "8083",
		"prometheus_max_open_files": 3,
		"max_claim_age_for_proof_retry": 32,
		"proof_prevalidation": false,
		"ctx_cache_size": 20,
		"abci_logging": false,
		"show_relay_errors": true
	}
}

```

- `pocket util update-configs`
> Update the config file with new defaults params for **consensus / leveldbopts / p2p / cache / mempool / fastsync** .
>
> Creates a backup file named _config.json.bk_ under config/
>
> Example output:
```
successfully Updated Config file
```

### Pocket Query Namespace
Queries the current world state built on the Pocket node.

- `pocket query block <height>`
> Returns the block at the specified height.
>
> Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query block-height`
> Returns the current block height known by this node.
>
> Example output:
```
Block Height: <current block height>
```

- `pocket query tx <hash>`
> Returns a result transaction object
>> Arguments:
 > - `<hash>`: The hash of the transaction to query.

- `pocket query balance <accAddr> <height>`
> Returns the balance of the specified `<accAddr>` at the specified `<height>`.
>
> Arguments:
> - `<accAddr>`: The address of the account to query.
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.
> Example output:
```
Account balance: <balance of the account>
```

- `pocket query nodes --staking-status=<stakingStatus> --page=<page> --limit=<limit> <height>`
> Returns a page containing a list of nodes known at the specified `<height>`.
>
> Options:
> - `--staking-status`: Filters the node list with a staking status. Supported statuses are: `STAKED`, `UNSTAKED` and `UNSTAKING`.
> - `--page`: The current page you want to query.
> - `--limit`: The maximum amount of nodes per page.
>
> Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query node <nodeAddr> <height>`
> Returns the node at the specified `<height>`.
>
> Arguments:
> - `<nodeAddr>`: The node address to be queried.
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query node-params <height>`
> Returns the list of node params specified in the `<height>`.
>
> Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query signing-info <nodeAddr> <height>`
> Returns the signing info of the node with `<nodeAddr>` at `<height>`.
>
> Arguments:
> - `<nodeAddr>`: The node address to be queried.
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query supply <height>`
> Returns the total amount of POKT staked/unstaked by nodes, apps, DAO, and totals at the specified `<height>`.
>
> Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query apps --staking-status=<stakingStatus> --page=<page> --limit=<limit> <height>`
> Returns a page containing a  list of applications known at the specified `<height>`.
>
> Options:
> - `--staking-status`: Filters the node list with a staking status. Supported statuses are: `STAKED`, `UNSTAKED` and `UNSTAKING`.
> - `--page`: The current page you want to query.
> - `--limit`: The maximum amount of nodes per page.
>
> Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query app <appAddr> <height>`
> Returns the application at the specified `<height>`.
>
> Arguments:
> - `<appAddr>`: The application address to be queried.
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query app-params <height>`
> Returns the list of node params specified in the `<height>`.
>
> Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query node-receipts <nodeAddr> <height>`
> Returns the list of all receipts for work done by `<nodeAddr>`.
>
> Arguments:
> - `<nodeAddr>`: The node address to be queried.
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

>- `pocket query node-claims <nodeAddr> <height>`
 > Returns the list of all pending claims submitted by `<nodeAddr>`.
 >
 > Arguments:
 > - `<nodeAddr>`: The node address to be queried.
 > - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query node-proof <nodeAddr> <appPubKey> <networkId> <sessionHeight> <height>`
> Returns the receipt of work completed specific to the arguments.
>
> Arguments:
> - `<nodeAddr>`: The address of the node that submitted the proof.
> - `<appPubKey>`: The public key of the application the Node serviced.
> - `<networkId>`: The Network Identifier of the blockchain that was serviced.
> - `<sessionHeight>`: The session block for which the proof was submitted.
> - `<receiptType>`: "relay" or "challenge"
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query node-claim <nodeAddr> <appPubKey> <networkId> <sessionHeight> <height>`
> Returns the claim specific to the arguments.
>
> Arguments:
> - `<nodeAddr>`: The address of the node that submitted the proof.
> - `<appPubKey>`: The public key of the application the Node serviced.
> - `<networkId>`: The Network Identifier of the blockchain that was serviced.
> - `<sessionHeight>`: The session block for which the proof was submitted.
> - `<receiptType>`: "relay" or "challenge"
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query supported-networks <height>`
> Returns the list Network Identifiers supported by the network at the specified `<height>`.
>
> Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query pocket-params <height>`
> Returns the list of Pocket Network params specified in the `<height>`.
>
> Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.
