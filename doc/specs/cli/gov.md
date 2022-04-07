---
description: >- Functions for governance (DAO) transactions; only relevant to the DAOowner
(the account that has the permission to perform these transactions on behalf of the DAO).
---

# Gov Namespace

## Change Parameter

```text
pocket gov change_param <fromAddr> <chainID> <paramKey module/param> <paramValue (jsonObj)> <fee> <legacyCodec=(true | false)
```

If authorized by the DAO, submit a tx to change any param from any module. Will prompt the user for the account
passphrase.

Arguments:

* `<fromAddr>`: Sender address.
* `<chainID>`: The Pocket chain identifier; "mainnet" or "testnet".
* `<paramKey>`: Target parameter key to change in format module/param, e.g. `pos/ProposerPercentage`.
* `<paramValue>`: New value for key.
* `<fee>`:  An amount of uPOKT for the network.

Example output:

```text
Transaction submitted with hash: <Transaction Hash>
```

## Upgrade Protocol

```text
pocket gov upgrade <fromAddr> <atHeight> <version> <chainID> <fees> <legacyCodec=(true | false)
```

If authorized by the DAO, upgrade the protocol. Will prompt the user for the account passphrase.

Arguments:

* `<fromAddr>`: Sender address.
* `<atHeight>`: The target height at which the protocol will be upgraded.
* `<version>`: The target version the protocol will be upgraded to.
* `<chainID>`: The Pocket chain identifier; "mainnet" or "testnet".
* `<fee>`:  An amount of uPOKT for the network.

Example output:

```text
Transaction submitted with hash: <Transaction Hash>
```

## Send DAO Funds

```text
pocket gov transfer <amount> <fromAddr> <toAddr> <chainID> <fee>
```

If authorized by the DAO, move funds from the DAO treasury account. Will prompt the user for the account passphrase.

Arguments:

* `<amount>`: The amount of uPOKT to be sent.
* `<toAddr>`: Recipient address for the transaction.
* `<fromAddr>`: Leave blank with " " \(because DAO treasury account is the fromAddr\).
* `<chainID>`: The Pocket chain identifier; "mainnet" or "testnet".
* `<fee>`:  An amount of uPOKT for the network.

Example output:

```text
Transaction submitted with hash: <Transaction Hash>
```

## Burn DAO Funds

```text
pocket gov burn <amount> <fromAddr> <toAddr> <chainID> <fee> <legacyCodec=(true | false)
```

If authorized, burn funds from the DAO treasury account. Will prompt the user for the account passphrase.

Arguments:

* `<amount>`: The amount of uPOKT to be sent.
* `<toAddr>`: Recipient address for the transaction.
* `<fromAddr>`: Leave blank with " " \(because DAO treasury account is the fromAddr\).
* `<chainID>`: The Pocket chain identifier; "mainnet" or "testnet".
* `<fee>`:  An amount of uPOKT for the network.

Example output:

```text
Transaction submitted with hash: <Transaction Hash>
```

