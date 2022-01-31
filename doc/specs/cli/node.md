---
description: Functions for Node management.
---

# Nodes Namespace

## Stake a Node / Update Stake

```text
pocket nodes stake <operatorAddress||signerAddress> <amount> <RelayChainIDs> <serviceURI> <outputAddress||signerAddress> <networkID> <fee> <isBefore8.0>
```
Stake a node in the network, the signer may be the operator or the output address. The signer must specify the public key of the output or operator
Prompts the user for the `<fromAddr>` account passphrase.

After the 0.6.X upgrade, if the node is already staked, this transaction acts as an _update_ transaction. A node can update `<relayChainIDs>`, `<serviceURI>`, and increase the stake `<amount>` with this transaction. If the node is currently staked at `X` and you submit an update with new stake `Y`, only `Y-X` will be subtracted from an account. If no changes are desired for the parameter, just enter the current parameter value \(the same one you entered for your initial stake\).

Arguments:

* `<operatorAddress||signerAddress>`: operatorAddress is the only valid signer for blocks & relays.
* `<outputAddress||signerAddress>`: outputAddress is where reward and staked funds are directed.
* `<amount>`: The amount of uPOKT to stake. Must be higher than the current value of the `StakeMinimum`  parameter, found [here](https://docs.pokt.network/home/references/protocol-parameters#stakeminimum).
* `<relayChainIDs>`: A comma separated list of RelayChain Network Identifiers. Find the RelayChain Network Identifiers [here](https://docs.pokt.network/home/references/supported-blockchains).
* `<serviceURI>`: The Service URI Applications will use to communicate with Nodes for Relays.
* `<networkID>`: The Pocket chain identifier; "mainnet" or "testnet".
* `<fee>`:  An amount of uPOKT for the network.
* `<isBefore8.0>`:  true or false depending if upgrade is activated..

Example output:

```text
Transaction submitted with hash: <Transaction Hash>
```

## Unstake a Node

```text
pocket nodes unstake <fromAddr> <chainID> <fee>
```

Unstakes a Node from the `<chainID>` network, changing its status to `Unstaking`. Prompts the user for the `<fromAddr>` account passphrase.

Arguments:

* `<fromAddr>`: Target staked address.
* `<chainID>`: The Pocket chain identifier; "mainnet" or "testnet".
* `<fee>`:  An amount of uPOKT for the network.

Example output:

```text
Transaction submitted with hash: <Transaction Hash>
```

## Unjail a Node

```text
pocket nodes unjail <fromAddr> <chainID> <fee>
```

Unjails a Node from the `<chainID>` network, allowing it to participate in service and consensus again. Prompts the user for the `<fromAddr>` account passphrase.

Arguments:

* `<fromAddr>`: Target jailed address.
* `<chainID>`: The Pocket chain identifier; "mainnet" or "testnet".
* `<fee>`:  An amount of uPOKT for the network.

Example output:

```text
Transaction submitted with hash: <Transaction Hash>
```

