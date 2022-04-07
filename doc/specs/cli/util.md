---
description: Useful operations.
---

# Util Namespace

## Create chains.json

```text
pocket util generate-chains
```

Generate the chains.json file for RelayChain Network Identifiers. Find the RelayChain Network
Identifiers [here](https://docs.pokt.network/home/references/supported-blockchains).

Example output:

```text
Enter the ID of the network identifier:
0001
Enter the URL of the network identifier:
https://ethnode.test.com:8085
Would you like to enter another network identifier? (y/n)
n
chains.json contains:

0001 @ https://ethnode.test.com:8085
If incorrect: please remove the chains.json with the delete-chains command
```

## Delete chains.json

```text
pocket util delete-chains
```

Delete the chains.json file for RelayChain Network Identifiers. Find the RelayChain Network
Identifiers [here](https://docs.pokt.network/home/references/supported-blockchains).

Example Output:

```text
Successfully deleted chains.json.
```

## Rollback the Chain

```text
pocket util unsafe-rollback <height>
```

Rolls back the blockchain, the state, and app to a previous height

Arguments:

* `<height>`: the height you want to rollback to.

Flags

* `blocks`: rollback block store and state

## Decode Transaction

```text
pocket util decode-tx <tx> <legacyCodec=(true | false)>
```

Decodes a given transaction encoded in Amino/Proto base64 bytes

Arguments:

* `<tx>`: The transaction amino encoded bytes.
* `<legacyCodec>`: Enable/Disable amino encoding for transaction.

Example output:

```text
% pocket util decode-tx qgLbCxcNCp0Bq4P6fApLCkA3ZWFjZWFjZTYwNzY1YzhiYjU0NDAzOGUxNGRjOGMyNjQ1NWRmODJmNTVmOGVkZDc1M2EwNDU5ZmY4MzYxZmViEgQwMDIxGP1wEi8KIEd86o3r3PIS6aK3CW+8L3E9JZMEHFdM1kMmy7XmuSQ/EgsQ8YrPpKGm95f/ARieAiIUjDp8K56yjpfHbsHBoLReW9EfapcoARIOCgV1cG9rdBIFMTAwMDAaaQolnVRHdCAO6zUJvs6taFLJzycYSzl2lPHXTYkxOnru2wG+T5y3PxJAckq7juFqII9kg/QPK2JmnLYNUthqZXNbEEQ5Zb/Jk/yqA2kwKUKS9yAZMPX8anDHj5Ycrtkw+LWnyha7aKFFBCiFvpiZ3YOT2JQB
Type:           claim
Msg:             286 8C3A7C2B9EB28E97C76EC1C1A0B45E5BD11F6A97 1 0}
Fee:            10000upokt
Entropy:        -7732596869214888187
Memo:
Signer          8c3a7c2b9eb28e97c76ec1c1a0b45e5bd11f6a97
Sig:            0eeb3509becead6852c9cf27184b397694f1d74d89313a7aeedb01be4f9cb73f
```

## Generate Completion Script

```text
pocket util completion (shell=bash | zsh | fish | powershell)>
```

Generate completion script for the specified shell

Arguments:

* `<shell>`: the shell you currently use. Supported options: **bash / zsh / fish / powershell**

## Convert Evidence to Proto from Amino

```text
pocket util convert-pocket-evidence-db
```

Convert pocket-evidence-db to Proto from Amino.

Example Output:

```
Successfully converted pocket evidence db
```

## Update config.json With New Param Defaults

```text
pocket util update-configs
```

Update the config file with new params defaults for consensus / leveldbopts / p2p / cache / mempool / fastsync.

Example Output:

```text
Successfuly updated config file.
```

## Print Default config.json

```text
pocket util print-configs
```

Prints Default config.json to console.

Example Output:

```text
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

## Export Genesis for Reset

```text
pocket util export-genesis-for-reset <height> <newChainID>
```

In the event of a network reset, this will export a genesis file based on the previous state.

Arguments:

* `<height>`: the height from which to generate genesis.json
* `<newChainID>`: the chainID to use for exporting.

Example Output:

```text
{
    "app_hash": "",/
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
```

