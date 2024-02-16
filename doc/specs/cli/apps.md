---
description: Functions for Application management.
---

# Apps Namespace

## Stake an App / Update Stake

```text
pocket apps stake <fromAddr> <amount> <relayChainIDs> <chainID> <fee>
```

Stakes the Application into the network, making it available to receive service. Prompts the user for the `<fromAddr>`
account passphrase.

After the 0.6.X upgrade, if the app is already staked, this transaction acts as an _update_ transaction. An app can
update `<relayChainIDs>` and increase the stake `<amount>`. If the app is currently staked at `X` and you submit an
update with new stake `Y`, only `Y-X` will be subtracted from an account. If no changes are desired for the parameter,
just enter the current parameter value \(the same one you entered for your initial stake\).

Arguments:

- `<fromAddr>`: Target Address to stake.
- `<amount>`: The amount of uPOKT to stake. Must be higher than the current value of the `ApplicationStakeMinimum`
  parameter, found [here](https://docs.pokt.network/learn/protocol-parameters/#applicationstakeminimum).
- `<relayChainIDs>`: A comma separated list of RelayChain Network Identifiers. Find the RelayChain Network
  Identifiers [here](https://docs.pokt.network/supported-blockchains/).
- `<chainID>`: The Pocket chain identifier; "mainnet" or "testnet".
- `<fee>`: An amount of uPOKT for the network.

Example output:

```text
Transaction submitted with hash: <Transaction Hash>
```

## Unstake an App

```text
pocket apps unstake <fromAddr> <chainID> <fee>
```

Unstakes an Application from the `<chainID>` network, changing its status to `Unstaking`. Prompts the user for
the `<fromAddr>` account passphrase.

Arguments:

- `<fromAddr>`: The address of the sender.
- `<chainID>`: The Pocket chain identifier; "mainnet" or "testnet".
- `<fee>`: An amount of uPOKT for the network.

Example output:

```text
Transaction submitted with hash: <Transaction Hash>
```

## Create an Application Authentication Token \(AAT\)

```text
pocket apps create-aat <appAddr> <clientPubKey>
```

This CLI is hard-coded to generate an AAT with spec version 0.0.1.

The output is intended to be embedded into the Gateway for Relay servicing.

Upon AAT generation, the user will be prompted for the <appAddr> account passphrase.

Read the Application Authentication Token documentation here:

{% page-ref page="../application-auth-token.md" %}

Arguments:

- `<appAddr>`: Application Address
  The address of the `Application` account to use to produce this AAT.
  This is the account that has to be staked on-chain to be able to use the Pocket Network.
- `<clientPubKey>`: Gateway Public Key
  The public key of the gateway that will be signing and sending Relays sent to the Pocket Network.

Example output:

```javascript
{
    "version" : "0.0.1",
    "applicationPublicKey": "0x...",
    "clientPublicKey": "0x...",
    "signature": "0x..."
}
```
