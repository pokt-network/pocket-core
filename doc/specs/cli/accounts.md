---
description: >- The accounts namespace handles all account related interactions, from creating and deleting accounts, to
importing and exporting accounts.
---

# Accounts Namespace

## Show All Accounts

```text
pocket accounts list
```

Lists all the account addresses currently stored in the keybase.

Example output:

```text
(0) 53d809964195172f2970219dfcb0007f33150623
(1) 59f08710afbad0e20352340780fdbf4e47622a7c
```

## Show Details of an Account

```text
pocket accounts show <address>
```

Lists an account address and public key.

Arguments:

* `<address>`: The address to be fetched.

Example output:

```text
Address: 0x.....
Public Key: 0x....
```

## Create an Account

```text
pocket accounts create
```

Creates and persists a new account in the Keybase. Will prompt the user for
a [BIP-0039](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki) password for the generated mnemonic and for
a passphrase to encrypt the generated keypair. _**Make sure to keep a note of this passphrase in a secure place.**_

Example output:

```text
Account generated successfully.
Address: 0x....
```

## Import an Account

```text
pocket accounts import-raw <private-key-hex>
```

Imports an account using the provided `<private-key-hex>`. Will prompt the user for
a [BIP-0039](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki) password for the imported private key and
for a passphrase to encrypt the generated keypair. _**Make sure to keep a note of this passphrase in a secure place.**_

Arguments:

* `<private-key-hex>`: Target raw private key bytes to be stored on keybase.

Example output:

```text
Account imported successfully.
Address: 0x....
```

## Encrypted Import

```text
pocket accounts import-armored <armoredJSONFile>
```

Imports an account using the Encrypted ASCII armored `<armoredJSONFile>`. Will prompt the user for a decryption
passphrase of the `<armoredJSONFile>` string and for an encryption passphrase to store in the Keybase. _**Make sure to
keep a note of this passphrase in a secure place.**_

Arguments:

* `<armoredJSONFile>`: Target file with encrypted encoded private key.

Example output:

```text
Account imported successfully.
Address: 0x....
```

## Export an Account

```text
pocket accounts export [--path <path>] <address>
```

Exports the account with `<address>`, encrypted and ASCII armored. Will prompt the user for the account passphrase and
an encryption passphrase for the exported account. _**Make sure to keep a note of this passphrase in a secure place.**_

Options:

* `--path`: Target path to send armored private key.

Arguments:

* `<address>`: The address of the account to be exported.

Example output:

```text
Exported account: <armored string>
```

## Raw Export

```text
pocket accounts export-raw <address>
```

Exports the raw private key in hex format. Will prompt the user for the account passphrase. _**NOTE: THIS METHOD IS NOT
RECOMMENDED FOR SECURITY REASONS, USE AT YOUR OWN RISK.**_

Arguments:

* `<address>`: The address of the account to be exported.

Example output:

```text
Exported account: 0x...
```

## Delete an Account

```text
pocket accounts delete <address>
```

Deletes an account from the Keybase. Will prompt the user for the account passphrase.

Arguments:

* `<address>`: The address to be deleted.

Example output:

```text
KeyPair 0x... deleted successfully.
```

### Unsafe Delete

```text
pocket accounts unsafe-delete <address>
```

Deletes an account from the keybase without a passphrase.

Arguments:

* `<address>`: The address to be deleted.

Example output:

```text
KeyPair 0x... deleted successfully.
```

## Show the Main Validator

```text
pocket accounts get-validator
```

Returns the main validator from the `priv_val` file.

## Set the Main Validator

```text
pocket accounts set-validator <address>
```

Sets a new main validator in the `priv_val` file.

Arguments:

* `<address>`: Target address.

## Update an Account's Passphrase

```text
pocket accounts update-passphrase <address>
```

Updates the passphrase for the indicated account. Will prompt the user for the current account passphrase and the new
account passphrase.

Arguments:

* `<address>`: Target address.

Example output:

```text
KeyPair 0x... passphrase updated successfully.
```

## Sign a Message

```text
pocket accounts sign <address> <msg>
```

Signs the specified `<msg>` using the specified `<address>` account credentials. Will prompt the user for the account
passphrase.

Arguments:

* `<address>`: Target address.
* `<msg>`: The message to be signed in hex string format.

Example output:

```text
Original Message: 0x...
Signature: 0x...
```

## Send Transaction

```text
pocket accounts send-tx <fromAddr> <toAddr> <amount> <chainID> <fee> <memo>
```

Sends `<amount>` uPOKT from `<fromAddr>` to `<toAddr>`. Prompts the user for `<fromAddr>` account passphrase.

Arguments:

* `<fromAddr>`: Sender address.
* `<toAddr>`: Recipient address for the transaction.
* `<amount>`: The amount of uPOKT to be sent.
* `<chainID>`: The Pocket chain identifier; "mainnet" or "testnet".
* `<fee>`:  An amount of uPOKT for the network.
* `<memo>`: Written message.

Example output:

```text
Transaction submitted with hash: <Transaction Hash>
```

## Send Raw Transaction

```text
pocket accounts send-raw-tx <fromAddr> <txBytes>
```

Sends presigned transaction through Tendermint node.

Arguments:

* `<fromAddr>`: Sender address.
* `<txBytes>`: Encoded and signed byte representation of the tx.

## Create a Multi-sig Account

```text
pocket accounts create-multi-public <hex-pubkeys>
```

Multi-signature accounts enable multiple individual accounts to share an account and create transactions that require
signatures from all accounts.

Important notes:

* Pocket Core does not save the multi-sig account in your keybase.
* You will need to remember the order in which the public keys have been assigned to the multi-sig account.

Arguments:

* `<hex-pubkeys>`: ordered comma separated keys. _**WARNING: changing the order creates a different address.**_

## Build a Multi-sig Transaction

```text
pocket accounts build-MS-Tx <signer-address> <json-message> <hex-pubkeys> <chainID> <fee>
```

Build and sign a multisignature transaction from scratch. Result is hex encoded std transaction object.

Arguments:

* `<signer-address>`: Address building & signing.
* `<json-message>`: Message structure for the transaction.
* `<hex-pubkeys>`: Ordered comma separated keys. _**WARNING: must be in the same order as when you created the multi-sig
  account.**_
* `<chainID>`: The Pocket chain identifier; "mainnet" or "testnet".
* `<fee>`:  An amount of uPOKT for the network.

## Sign a Multi-sig Transaction

```text
pocket accounts sign-ms-tx <signer-address> <hex-stdtx> <hex-pubkeys> <chainID> <fee>
```

Sign a multisignature transaction using public keys, and the transaction object out of order. Result is hex encoded
standard transaction object.

Arguments:

* `<signer-address>`: Address building & signing.
* `<hex-stdtx>`: Prebuilt hexadecimal standard transaction.
* `<hex-pubkeys>`: Ordered comma separated keys. _**WARNING: must be in the same order as when you created the multi-sig
  account.**_
* `<chainID>`: The Pocket chain identifier; "mainnet" or "testnet".
* `<fee>`:  An amount of uPOKT for the network.

## Sign a Multi-sig Transaction as the Next Signer

```text
pocket accounts sign-ms-next <signer-address> <hex-stdtx> <hex-pubkeys> <chainID> <fee>
```

Sign a multisignature transaction object, result is hex encoded standard transaction object. _**WARNING: signer address
MUST be the next signer \(in order of public keys in the multisignature\) or the signature will be invalid.**_

Arguments:

* `<signer-address>`: Address building & signing.
* `<hex-stdtx>`: Prebuilt hexadecimal standard transaction.
* `<hex-pubkeys>`: Ordered comma separated keys. _**WARNING: must be in the same order as when you created the multi-sig
  account.**_
* `<chainID>`: The Pocket chain identifier; "mainnet" or "testnet".
* `<fee>`:  An amount of uPOKT for the network.
