---
description: Queries the current world state built on the Pocket node.
---

# Query Namespace

## Network

### Block Height

```text
pocket query height
```

Returns the current block height known by this node.

Example output:

```text
Block Height: <current block height>
```

### Total POKT Supply

```text
pocket query supply [<height>]
```

Returns the total amount of POKT staked/unstaked by nodes, apps, DAO, and totals at the specified `<height>`.

Optional Arguments:

* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

Example Output:

```text
http://localhost:8082/v1/query/supply
{
    "app_staked": "1542099237724",
    "dao": "50000096646730",
    "node_staked": "5001315185588861",
    "total": "147579001422608721340",
    "total_staked": "5052857381473315",
    "total_unstaked": "147573948565227248025"
}
```

_Note: the `dao` value is currently included in the `total_staked` value._

### Supported Blockchains

```text
pocket query supported-networks [<height>]
```

Returns the list of RelayChain Network Identifiers supported by the network at the specified `<height>`, meaning they've
been whitelisted to generate revenue for nodes.

Optional Arguments:

* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

Example Output:

```text
pocket query supported-networks
[
    "0001",
    "0021"
]
```

### Block Details

```text
pocket query block <height>
```

Returns the block at the specified height.

Arguments:

* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

### Block Transactions

```text
pocket query block-txs <height> [<page> <per_page> <prove> <order>]
```

Retrieves the transactions in the block `<height>` .

Arguments:

* `<address>`: The specified height of the block to be queried, defaults to `0` which brings the latest block known to
  this node

Optional arguments:

* `<page>`: the page of the transaction list that you want to focus on.
* `<per_page>`: how many transactions you want to see per page of the transaction list.
* `<prove>`: the Tendermint merkle proof that the transaction exists. This can be **true** or **false**.
* `<order>`: Sort of the results. Default is desc.

## Parameters

### All Parameters

```text
pocket query params [<height>]
```

Returns all the parameters at the specified `<height>`.

Optional Arguments:

* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

### Specific Parameter

```text
pocket query param <key> [<height>]
```

Get a parameter with the given `<key>` at the specified `<height>`.

Arguments:

* `<key>`: key identifier of the param.

Optional Arguments:

* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

### Pocket Core Parameters

```text
pocket query pocket-params [<height>]
```

Returns the list of parameters in the Pocket Core module at the specified `<height>`.

Optional Arguments:

* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

```text
{
    "claim_expiration": "120",
    "minimum_number_of_proofs": "10",
    "proof_waiting_period": "3",
    "replay_attack_burn_multiplier": "3",
    "session_node_count": "5",
    "supported_blockchains": [
        "0001",
        "0021"
    ]
}
```

### Node Parameters

```text
pocket query node-params [<height>]
```

Returns the list of parameters in the PoS \(Node\) module at the specified `<height>`.

Optional Arguments:

* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

```text
{
    "dao_allocation": "10",
    "downtime_jail_duration": "3600000000000",
    "max_evidence_age": "120000000000",
    "max_jailed_blocks": "37960",
    "max_validators": "5000",
    "maximum_chains": "15",
    "min_signed_per_window": "6.000000000000000000",
    "proposer_allocation": "1",
    "relays_to_tokens_multiplier": "0",
    "session_block_frequency": "4",
    "signed_blocks_window": "10",
    "slash_fraction_double_sign": "0.050000000000000000",
    "slash_fraction_downtime": "0.000001000000000000",
    "stake_denom": "upokt",
    "stake_minimum": "15000000000",
    "unstaking_time": "3600000000000"
}
```

### App Parameters

```text
pocket query app-params [<height>]
```

Returns the list of parameters in the Application module at the specified `<height>`.

Optional Arguments:

* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

```text
http://localhost:8082/v1/query/appparams
{
    "app_stake_minimum": "1000000",
    "base_relays_per_pokt": "167",
    "max_applications": "9223372036854775807",
    "maximum_chains": "15",
    "participation_rate_on": false,
    "stability_adjustment": "0",
    "unstaking_time": "3600000000000"
}
```

## Accounts

### Account Details

```text
pocket query account <address> [<height>]
```

Returns the account structure for a specific `<address>`.

Arguments:

* `<address>`: Target address.

Optional Arguments:

* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

### Account Transactions

```text
pocket query account-txs <address> [<page> <per_page> <prove> <received=(true | false)> <order=(asc | desc)>]
```

Retrieves the transactions sent by the address.

Arguments:

* `<address>`: Target address.

Optional arguments:

* `<page>`: The current page you want to query. Default to first page.
* `<per_page>`: The maximum amount elements per page. Default is 30 elements per page.
* `<prove>`: Shows proof. Default is false.
* `<received>`: Check if target address is recipient. Default is false.
* `<order>`: Sort of the results. Default is desc.

### Transaction

```text
pocket query tx <hash>
```

Returns a result transaction object for a specified transaction `<hash>`.

Arguments:

* `<hash>`: The hash of the transaction to query.

### POKT Balance of Account

```text
pocket query balance <address> [<height>]
```

Returns the balance of the specified `<address>` at the specified `<height>`.

Arguments:

* `<address>`: Target address.

Optional Arguments:

* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

Example output:

```text
Account balance: <balance of the account>
```

## Nodes

### List of All Nodes at Height

```text
pocket query nodes [--staking-status=(staked | unstaking)] [--jailed-status=(jailed | unjailed)][page=<page>] [--limit=<limit>] <height>
```

Returns a page containing a list of nodes known at the specified `<height>`.

Options:

* `--staking-status`: Filters the node list with a staking status. Supported statuses are: `staked` and `unstaking`.
* `--jailed-status`: Filters the node list with jailed/unjailed validators. Supported statuses are: `jailed`
  and `unjailed`.
* `--page`: The current page you want to query.
* `--limit`: The maximum amount of nodes per page.

Arguments:

* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

### Node Details

```text
pocket query node <address> <height>
```

Returns the node at the specified `<height>`.

Arguments:

* `<address>`: Target address.
* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

### Node Signing Info

```text
pocket query signing-info <address> [<height>]
```

Returns the signing info of the node `<address>` at `<height>`.

Arguments:

* `<address>`: Target address.
* `<height>`: The specified height of the block to be queried, defaults to `0` which brings the latest block known to
  this node.

### List of Relay Proofs Submitted by Node

```text
pocket query node-claims [<address>] [<height>]
```

Returns the list of all pending claims submitted by `<address>` at specified `<height>`.

Optional Arguments:

* `<address>`: Target address. Defaults to `0` which brings all the claims for the specified height.
* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

### Relay Proof Details

```text
pocket query node-claim <address> <appPubKey> <claimType=(relay | challenge)> <chainID> <sessionHeight> [<height>]
```

Returns the claim specific to the arguments.

Arguments:

* `<address>`: The address of the node that submitted the proof.
* `<appPubKey>`: The public key of the application the Node serviced.
* `<chainID>`: The Network Identifier of the blockchain that was serviced.
* `<sessionHeight>`: The session block for which the proof was submitted.
* `<receiptType>`: An enum string that can be "relay" or "challenge".

Optional Arguments:

* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

## Apps

### List of All Apps at Height

```text
pocket query apps [--staking-status=(staked | unstaking)] [--page=<page>] [--limit=<limit>] [<height>]
```

Returns a page containing a list of applications known at the specified `<height>`.

Options:

* `--staking-status`: Filters the app list with a staking status. Supported statuses are: `staked` and `unstaking`.
* `--page`: The current page you want to query.
* `--limit`: The maximum amount of apps per page.

Optional Arguments:

* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

### App Details

```text
pocket query app <address> [<height>]
```

Returns the application at the specified `<height>`.

Arguments:

* `<address>`:Target address.

Optional Arguments:

* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

## Gov

### daoOwner

```text
pocket query daoOwner [<height>]
```

Retrieves the owner of the DAO, the account which has the permission to submit governance transactions on behalf of the
DAO.

Optional Arguments:

* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

### ACL

```text
pocket query acl [<height>]
```

Returns the access control list of governance param \(which account can change the param\).

Optional Arguments:

* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.

### Latest Upgrade

```text
pocket query upgrade [<height>]
```

Retrieves the latest protocol upgrade executed by governance using the `pocket gov upgrade` command.

Optional Arguments:

* `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to
  this node.
