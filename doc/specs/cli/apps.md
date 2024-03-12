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
  parameter.
- `<relayChainIDs>`: A comma separated list of RelayChain Network Identifiers. Find the RelayChain Network
  Identifiers [here](https://docs.pokt.network/reference/supported-chains).
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

Creates a signed Application Authentication Token.
Creates a signed AAT (= Application Authentication Token) where the version is
hardcoded as "0.0.1" that is the only version supported by the protocol.

This command prompts you to input the `<appAddr>` account passphrase.
When you send a relay request with AAT, `<appAddr>` needs to be a staked
application.

Please read [application-auth-token.md](../application-auth-token.md)
for additional details.

Arguments:

- `<appAddr>`:
  The address of an `Application` account to use to produce this AAT.
  The account has to be staked on-chain to be able to use the Pocket Network.
- `<clientPubKey>`:
  The public key of a client that will be signing and sending Relays to the Pocket Network.

Example output:

```javascript
{
    "version" : "0.0.1",
    "applicationPublicKey": "0x...",
    "clientPublicKey": "0x...",
    "signature": "0x..."
}
```
