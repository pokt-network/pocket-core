---
description: Functions for Node management.
---

# Nodes Namespace

Terminology:

- `Operator Address`: `Non-Custodial Address` aka `Node Runner` aka `Devops`
- `Output Address`: `Custodial Address` aka `Deposit Owner` aka `Reward Earner`

Basic Rules

1. The `Operator Address` is the **only** valid signer for:
   1. Block Tx
   2. Claim & Proof Tx
2. `Operator` and `Output` Address are valid signers for:
   1. Stake Tx
   2. EditStake Tx
   3. Unstake Tx
   4. Unjail Tx
3. The `Output Address` is where:
   1. Rewards are sent
   2. Unstaked funds are sent (after unstaking)
4. `Operator` and `Output` Address **cannot** be edited once staked
   1. Requires unstaking to change either one

## Stake a Node / Update Stake (Custodial)

```text

pocket nodes stake custodial <fromAddr> <amount> <relayChainIDs> <serviceURI> <networkID> <fee> <isBefore8.0>
```

Stakes a custodial Node into the network, making it available for service. The signer must provide the operator address passphrase when prompted for the `<fromAddr>` account passphrase.

if the node is already staked, this transaction acts as an _update_ transaction. A node can update `<relayChainIDs>`
, `<serviceURI>`, and increase the stake `<amount>` with this transaction. If the node is currently staked at `X` and
you submit an update with new stake `Y`, only `Y-X` will be subtracted from an account. If no changes are desired for
the parameter, just enter the current parameter value \(the same one you entered for your initial stake\).

Arguments:

- `<fromAddr>`: Target Address to stake.
- `<amount>`: The amount of uPOKT to stake. Must be higher than the current value of the `StakeMinimum` parameter,
  found [here](https://docs.pokt.network/learn/protocol-parameters/#stakeminimum).
- `<relayChainIDs>`: A comma separated list of RelayChain Network Identifiers. Find the RelayChain Network
  Identifiers [here](https://docs.pokt.network/supported-blockchains/).
- `<serviceURI>`: The Service URI Applications will use to communicate with Nodes for Relays.
- `<networkID>`: The Pocket chain identifier; "mainnet" or "testnet".
- `<fee>`: An amount of uPOKT for the network.
- `<isBefore8.0>`: true or false depending if non-custodial upgrade is activated.

Example output:

```text
Transaction submitted with hash: <Transaction Hash>
```

## Stake a Node / Update Stake (Non-custodial / 0.8.X)

```text
pocket nodes stake non-custodial <operatorPublicKey> <outputAddress> <amount> <RelayChainIDs> <serviceURI> <networkID> <fee> <isBefore8.0>
```

Stakes a non-custodial node in the network, making it available for service. The signer may be the operator or the output address. The signer must specify the passphrase of either the output or operator address when prompted for the `<fromAddr>` account passphrase. This will determine where the staked funds are taken from. 

if the node is already staked, this transaction acts as an _update_ transaction. A node can update `<relayChainIDs>`
, `<serviceURI>`, and increase the stake `<amount>` with this transaction. If the node is currently staked at `X` and
you submit an update with new stake `Y`, only `Y-X` will be subtracted from an account. If no changes are desired for
the parameter, just enter the current parameter value \(the same one you entered for your initial stake\).

Arguments:

- `<operatorPublicKey>`: operatorAddress is the only valid signer for blocks & relays.
- `<outputAddress>`: outputAddress is where reward and staked funds are directed.
- `<amount>`: The amount of uPOKT to stake. Must be higher than the current value of the `StakeMinimum` parameter,
  found [here](https://docs.pokt.network/learn/protocol-parameters/#stakeminimum).
- `<relayChainIDs>`: A comma separated list of RelayChain Network Identifiers. Find the RelayChain Network
  Identifiers [here](https://docs.pokt.network/supported-blockchains/).
- `<serviceURI>`: The Service URI Applications will use to communicate with Nodes for Relays.
- `<networkID>`: The Pocket chain identifier; "mainnet" or "testnet".
- `<fee>`: An amount of uPOKT for the network.
- `<isBefore8.0>`: true or false depending if non custodial upgrade is activated.

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

- `<operatorAddr>`: Target staked operator address.
- `<fromAddr>`: Signer address.
- `<networkID>`: The Pocket chain identifier; "mainnet" or "testnet".
- `<fee>`: An amount of uPOKT for the network.
- `<isBefore8.0>`: true or false depending if non custodial upgrade is activated.

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

- `<operatorAddr>`: Target jailed operator address.
- `<fromAddr>`: Signer address.
- `<networkID>`: The Pocket chain identifier; "mainnet" or "testnet".
- `<fee>`: An amount of uPOKT for the network.
- `<isBefore8.0>`: true or false depending if non custodial upgrade is activated.

Example output:

```text
Transaction submitted with hash: <Transaction Hash>
```
