# Pocket Core RPC Specification
## Version 0.0.1

### Overview
This document serves as a specification for the RPC of the Pocket Core application. There's no protocol verification for these commands, however because they map closely to protocol functions.

### Namespaces
The RPC will contain multiple namespaces listed below:

- Client: Contains all the calls pertinent to pocket core clients (relay and dispatch)
- Query: All queries to the world state are contained in this call.

### RPC Functions Format
Each RPC Function will be in the following format:

- Version: The version for Pocket Core rpc, for example: `/v1`
- Namespace: The namespace of the function: `query`
- Function Name: The name of the actual function to be called: `tx`

ALl RPC calls in Pocket Core will be HTTP `POST`

### RPC-SPEC in OPENAPI / SWAGGER

please see rcp-spec.yaml

this document is for references only and will be deprecated soon.


### Client Namespace
The default namespace contains functions that are pertinent to the execution of the Pocket Node.

- /v1/client/dispatch
> sends a dispatch request to the network and get the nodes that will be servicing your requests for the session.

     request: `pocketTypes.SessionHeader`

     response: `pocketTypes.Session`

- /v1/client/relay
> sends a relay request through pocket network to a external chain

    request: `types.Relay`

    response: `pocketTypes.RelayResponse`

- /v1/client/rawtx
> sends a raw transaction to the pocket blockchain

    request: `sendRawTxParams`

    response: `sdk.TxResponse`

### Query Namespace
The `query` namespace handles all queries to the current world state built on the Pocket node.

- /v1/query/block
> Query a block in the pocket blockchain, height = 0 returns latest height.

    request: `heightParams`

- /v1/query/tx
> Query a transaction in the pocket blockchain by transaction hash.

    request: `hashParams`

- /v1/query/height
> Query the current pocket blockchain height

- /v1/query/balance
> Query account balance for a specified Address and height.

    request: `heightAddrParams`

- /v1/query/account
> Query account for a specified Address and height.

    request: `heightAddrParams`

- /v1/query/nodes
> Query Nodes in the pocket network by height and staking_status, empty ("") staking_status returns a page of nodes

    request: `heightAndValidatorsOpts`

- /v1/query/node
> Query a specific node by address and height

    request: `heightAddrParams`

- /v1/query/nodeparams
> Query the POS parameters for nodes at a specified height

    request: `heightParams`

- /v1/query/nodereceipts
> Query a node receipts for an address at the specified height

    request `heightAddrParams`

- /v1/query/nodereceipt
> Query a specific receipt

    request: `queryNodeReceipts`

    response: `pocketTypes.Receipt`

- /v1/query/apps
> Query Apps in the pocket network by height and staking_status, empty ("") staking_status returns a page of apps

    request: `heightAndApplicationsOpts`

- /v1/query/app
> Query a specific app by address and height

    request: `heightAddrParams`

- v1/query/appparams
> Query the parameters for apps at a specified height

    request `heightParams` at a specified height

- /v1/query/pocketparams
> Query General Parameters

    request `heightParams`

- /v1/query/supportedchains
> Query supported chains by the pocket Network

    request `heightParams`

- /v1/query/supply
> Query POKT supply for NODES, APPS and DAO

    request `heightParams`

    response: `querySupplyResponse`
