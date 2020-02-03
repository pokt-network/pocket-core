# Pocket Core CLI Interface Specification
## Version 0.0.1

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
> Example output:
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
> Example output:
```
Transaction submitted with hash: <Transaction Hash>
```

- `pocket node unstake <fromAddr>`
> Unstakes a Node from the network, changing its status to `Unstaking`. Prompts the user for the `<fromAddr>` account passphrase.
>
> Arguments:
> - `<fromAddr>`: The address of the sender.
> Example output:
```
Transaction submitted with hash: <Transaction Hash>
```

- `pocket node unjail <fromAddr>`
> Unjails a Node from the network, allowing it to participate in service and consensus again. Prompts the user for the `<fromAddr>` account passphrase.
>
> Arguments:
> - `<fromAddr>`: The address of the sender.
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
> Example output:
```
Transaction submitted with hash: <Transaction Hash>
```

- `pocket app unstake <fromAddr>`
> Unstakes an Application from the network, changing its status to `Unstaking`. Prompts the user for the `<fromAddr>` account passphrase.
>
> Arguments:
> - `<fromAddr>`: The address of the sender.
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

- `pocket util generate-chain <ticker> <netid> <client> <version> <interface>`
> Creates a Network Identifier hash, used as a parameter for both Node and App stake.
>
> Arguments:
> - `<ticker>`: The ticker of the blockchain that will be accessed, e.g. `ETH`, `BTC`, `AION`.
> - `<netid>`: The network identifier for the blockchain that will be access, `ETH` e.g. `1`, `4`.
> - `<client>`: The node client for the specified blockchain that will be accessed, `ETH` e.g.: `geth`, `parity`.
> - `<version>`: The version of the aforementioned node client, e.g.: `1.9.2`.
> - `<interface>`: The interface used to talk to the aforementioned node client, e.g.: `wss://`, `https://`.
> Example output:
```
Network Identifier: 0x...
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

- `pocket query node-status`
> Returns the current node status.

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

- `pocket query all-nodes --staking-status=<stakingStatus> <height>`
> Returns the list of all nodes known at the specified `<height>`.
>
> Options:
> - `--staking-status`: Filters the node list with a staking status. Supported statuses are: `STAKED`, `UNSTAKED` and `UNSTAKING`.
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

- `pocket query all-apps --staking-status=<stakingStatus> <height>`
> Returns the list of all applications known at the specified `<height>`.
>
> Options:
> - `--staking-status`: Filters the node list with a staking status. Supported statuses are: `STAKED`, `UNSTAKED` and `UNSTAKING`.
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

- `pocket query node-proofs <nodeAddr> <height>`
> Returns the list of all Relay Batch proofs submitted by `<nodeAddr>`.
>
> Arguments:
> - `<nodeAddr>`: The node address to be queried.
> - `<height>`: The specified height of the block to be queried. Defaults to `0` which brings the latest block known to this node.

- `pocket query node-proof <nodeAddr> <appPubKey> <networkId> <sessionHeight> <height>`
> Returns the Relay Batch proof specific to the arguments.
>
> Arguments:
> - `<nodeAddr>`: The address of the node that submitted the proof.
> - `<appPubKey>`: The public key of the application the Node serviced.
> - `<networkId>`: The Network Identifier of the blockchain that was serviced.
> - `<sessionHeight>`: The session block for which the proof was submitted.
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

