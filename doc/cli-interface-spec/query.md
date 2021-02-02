# Pocket Query Namespace
Queries the current world state built on the Pocket node.

- `pocket query block <height>`
> Returns the block at the specified height.
>
> Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query height`
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

- `pocket query balance <address> [<height>]`
> Returns the balance of the specified `<accAddr>` at the specified `<height>`.
>
> Arguments:
> - `<address>`: Target address.
> Optional Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.
> Example output:
```
Account balance: <balance of the account>
```

- `pocket query nodes [--staking-status=(staked | unstaking)] [--jailed-status=(jailed | unjailed)][page=<page>] [--limit=<limit>] <height>`
> Returns a page containing a list of nodes known at the specified `<height>`.
>
> Options:
> - `--staking-status`: Filters the node list with a staking status. Supported statuses are: `staked` and `unstaking`.
> - `--jailed-status`: Filters the node list with jailed/unjailed validators. Supported statuses are: `jailed` and `unjailed`.
> - `--page`: The current page you want to query.
> - `--limit`: The maximum amount of nodes per page.
>
> Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query node <address> <height>`
> Returns the node at the specified `<height>`.
>
> Arguments:
> - `<address>`: Target address.
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query node-params [<height>]`
> Returns the list of node params specified in the `<height>`.
>
> Optional Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query supply [<height>]`
> Returns the total amount of POKT staked/unstaked by nodes, apps, DAO, and totals at the specified `<height>`.
>
> Optional Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query apps [--staking-status=(staked | unstaking)] [--page=<page>] [--limit=<limit>] [<height>]`
> Returns a page containing a  list of applications known at the specified `<height>`.
>
> Options:
> - `--staking-status`: Filters the app list with a staking status. Supported statuses are: `staked` and `unstaking`.
> - `--page`: The current page you want to query.
> - `--limit`: The maximum amount of nodes per page.
>
> Optional Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query app <address> [<height>]`
> Returns the application at the specified `<height>`.
>
> Arguments:
> - `<address>`:Target address. 
 Optional Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query app-params [<height>]`
> Returns the list of node params specified in the `<height>`.
>
> Optional Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query node-claims <address> [<height>]`
 > Returns the list of all pending claims submitted by `<address>`.
 >
 > Arguments:
 > - `<address>`: Target address.
> Optional Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query node-claim <address> <appPubKey> <claimType=(relay | challenge)> <networkId> <sessionHeight> [<height>]`
> Returns the claim specific to the arguments.
>
> Arguments:
> - `<address>`: The address of the node that submitted the proof.
> - `<appPubKey>`: The public key of the application the Node serviced.
> - `<networkId>`: The Network Identifier of the blockchain that was serviced.
> - `<sessionHeight>`: The session block for which the proof was submitted.
> - `<receiptType>`: An enum string that can be "relay" or "challenge".
> Optional Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query supported-networks [<height>]`
> Returns the list Network Identifiers supported by the network at the specified `<height>`.
>
> Optional Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query pocket-params [<height>]`
> Returns the list of Pocket Network params specified in the `<height>`.
>
> Optional Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query daoOwner [<height>]`
> Retrieves the owner of the DAO

> Optional Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query acl [<height>]`
> Returns the access control list of governance param (which account can change the param).

> Optional Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query params [<height>]`
> Returns all the parameters at the specified <height>

> Optional Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query param <key> [<height>]`
> Get a parameter with the given key.
>
> Arguments:
> - `<key>`: key identifier of the param.
> Optional Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query upgrade [<height>]`
> Returns the latest protocol upgrade by governance
>
> Optional Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query account <address> [<height>]`
> Returns the account structure for a specific address
>
> Arguments:
> - `<address>`: Target address.
> Optional Arguments:
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query account-txs <address> <page> <per_page> <prove> <received=(true | false)> <order=(asc | desc)>`
> Returns the latest protocol upgrade by governance
>
> Arguments:
> - `<address>`: Target address.
> Optional arguments:
> - `<page>`: The current page you want to query. Default to first page
> - `<per_page>`: The maximum amount elements per page. Default is 30 elements per page
> - `<prove>`: Shows proof. Default is false.
> - `<received>`: Check if target address is recipient. Default is false.
> - `<order>`: Sort of the results. Default is desc. Default is desc.

- `pocket query block-txs <address> [<page> <per_page> <prove> <received=(true | false)> <order=(asc | desc)>]`
> Returns the latest protocol upgrade by governance
>
> Arguments:
> - `<address>`: Target address.
> Optional arguments:
> - `<page>`: The current page you want to query. Default to first page
> - `<per_page>`: The maximum amount elements per page. Default is 30 elements per page
> - `<prove>`: Shows proof. Default is false.
> - `<received>`: Check if target address is recipient. Default is false.
> - `<order>`: Sort of the results. Default is desc. Default is desc.