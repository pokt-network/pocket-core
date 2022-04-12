---
description: Functions for Node management.
---

# Nodes Namespace

## Stake a Node / Update Stake (Custodial)

```text

pocket nodes stake custodial <fromAddr> <amount> <relayChainIDs> <serviceURI> <networkID> <fee> <isBefore8.0>
```

Stakes the Node into the network, making it available for service. Prompts the user for the `<fromAddr>` account
passphrase.

if the node is already staked, this transaction acts as an _update_ transaction. A node can update `<relayChainIDs>`
, `<serviceURI>`, and increase the stake `<amount>` with this transaction. If the node is currently staked at `X` and
you submit an update with new stake `Y`, only `Y-X` will be subtracted from an account. If no changes are desired for
the parameter, just enter the current parameter value \(the same one you entered for your initial stake\).

Arguments:

* `<fromAddr>`: Target Address to stake.
* `<amount>`: The amount of uPOKT to stake. Must be higher than the current value of the `StakeMinimum`  parameter,
  found [here](https://docs.pokt.network/home/references/protocol-parameters#stakeminimum).
* `<relayChainIDs>`: A comma separated list of RelayChain Network Identifiers. Find the RelayChain Network
  Identifiers [here](https://docs.pokt.network/home/references/supported-blockchains).
* `<serviceURI>`: The Service URI Applications will use to communicate with Nodes for Relays.
* `<networkID>`: The Pocket chain identifier; "mainnet" or "testnet".
* `<fee>`:  An amount of uPOKT for the network.
* `<isBefore8.0>`:  true or false depending if non-custodial upgrade is activated.


Example output:

```text
Transaction submitted with hash: <Transaction Hash>
```

## Stake a Node / Update Stake (Non-custodial / 0.8.X)

```text
pocket nodes stake non-custodial <operatorPublicKey> <outputAddress> <amount> <RelayChainIDs> <serviceURI> <networkID> <fee> <isBefore8.0>
```

Stake a node in the network, the signer may be the operator or the output address. The signer must specify the public
key of the output or operator Prompts the user for the `<fromAddr>` account passphrase.

if the node is already staked, this transaction acts as an _update_ transaction. A node can update `<relayChainIDs>`
, `<serviceURI>`, and increase the stake `<amount>` with this transaction. If the node is currently staked at `X` and
you submit an update with new stake `Y`, only `Y-X` will be subtracted from an account. If no changes are desired for
the parameter, just enter the current parameter value \(the same one you entered for your initial stake\).

Arguments:

* `<operatorPublicKey>`: operatorAddress is the only valid signer for blocks & relays.
* `<outputAddress>`: outputAddress is where reward and staked funds are directed.
* `<amount>`: The amount of uPOKT to stake. Must be higher than the current value of the `StakeMinimum`  parameter,
  found [here](https://docs.pokt.network/home/references/protocol-parameters#stakeminimum).
* `<relayChainIDs>`: A comma separated list of RelayChain Network Identifiers. Find the RelayChain Network
  Identifiers [here](https://docs.pokt.network/home/references/supported-blockchains).
* `<serviceURI>`: The Service URI Applications will use to communicate with Nodes for Relays.
* `<networkID>`: The Pocket chain identifier; "mainnet" or "testnet".
* `<fee>`:  An amount of uPOKT for the network.
* `<isBefore8.0>`:  true or false depending if non custodial upgrade is activated.

Example output:

```text
Transaction submitted with hash: <Transaction Hash>
```

## Unstake a Node

```text
pocket nodes unstake <operatorAddr> <fromAddr> <networkID> <fee> <isBefore8.0>
```

Unstakes a Node from the `<networkID>` network, changing its status to `Unstaking`. Prompts the user for
the `<fromAddr>` account passphrase.

Arguments:

* `<operatorAddr>`: Target staked operator address.
* `<fromAddr>`: Signer address.
* `<networkID>`: The Pocket chain identifier; "mainnet" or "testnet".
* `<fee>`:  An amount of uPOKT for the network.
* `<isBefore8.0>`:  true or false depending if non custodial upgrade is activated.

Example output:

```text
Transaction submitted with hash: <Transaction Hash>
```

## Unjail a Node

```text
pocket nodes unjail <operatorAddr> <fromAddr> <networkID> <fee> <isBefore8.0>
```

Unjails a Node from the `<networkID>` network, allowing it to participate in service and consensus again. Prompts the
user for the `<fromAddr>` account passphrase.

Arguments:

* `<operatorAddr>`: Target jailed operator address.
* `<fromAddr>`: Signer address.
* `<networkID>`: The Pocket chain identifier; "mainnet" or "testnet".
* `<fee>`:  An amount of uPOKT for the network.
* `<isBefore8.0>`:  true or false depending if non custodial upgrade is activated.

Example output:

```text
Transaction submitted with hash: <Transaction Hash>
```

