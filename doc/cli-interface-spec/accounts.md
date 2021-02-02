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
> - `<address>`: Target address.
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
> - `<address>`: Target address.
> Example output:
```
KeyPair 0x... deleted successfully.
```

- `pocket accounts update-passphrase <address>`
> Updates the passphrase for the indicated account. Will prompt the user for the current account passphrase and the new account passphrase.
>
> Arguments:
> - `<address>`: Target address.
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

- `pocket accounts get-validator`
> Return the main validator from the priv_val file

- `pocket accounts set-validator <address>`
> Return the main validator from the priv_val file
>
> Arguments:
> - `<address>`: Target address.

- `pocket accounts import-raw <private-key-hex>`
> Imports an account using the provided `<private-key-hex>`. Will prompt the user for a [BIP-0039](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki) password for the imported private key and for a passphrase to encrypt the generated keypair.
>
> Arguments:
> - `<import-raw>`: Target raw private key bytes to be stored on keybase
> Example output:
```
Account imported successfully.
Address: 0x....
```

- `pocket accounts import-armored <aromredJSONFile>`
> Imports an account using the Encrypted ASCII armored `<armoredJSONFile>`. Will prompt the user for a decryption passphrase of the `<armoredJSONFile>` string and for an encryption passphrase to store in the Keybase.
>
> Arguments:
> - `<armoredJSONFile>`: Target file with encrypted encoded prvate key.
> Example output:
```
Account imported successfully.
Address: 0x....
```

- `pocket accounts export [--path <path>] <address>`
> Exports the account with `<address>`, encrypted and ASCII armored. Will prompt the user for the account passphrase and an encryption passphrase for the exported account.
>
> Options:
> - `--path`: Target path to send armored private key 

> Arguments:
> - `<address>`: Target address for  exportation.
> Example output:
```
Exported account: <armored string>
```

- `pocket accounts export-raw <address>`
> Exports the raw private key in hex format. Will prompt the user for the account passphrase. ***NOTE***: THIS METHOD IS NOT RECOMMENDED FOR SECURITY REASONS, USE AT YOUR OWN RISK.*

> Arguments:
> - `<address>`: Target address for exportation.
> Example output:
```
Exported account: 0x...
```

- `pocket accounts send-tx <fromAddr> <toAddr> <amount> <chainID> <fee> <memo> <legacyCodec=(true | false)>`
> Sends `<amount>` POKT `<fromAddr>` to `<toAddr>`. Prompts the user for `<fromAddr>` account passphrase.
>
> Arguments:
> - `<fromAddr>`: Sender address.
> - `<toAddr>`: Recipient address for the transaction.
> - `<amount>`: The amount of POKT to be sent.
> - `<chainID>`: The pocket chain identifier.
> - `<fee>`:  An amount of POKT for the network .
> - `<memo>`: Written message.
> - `<legacyCodec>`: Enlble/Disable amino encoding for transaction.

> Example output:
```
Transaction submitted with hash: <Transaction Hash>
```

- `pocket accounts send-raw-tx <fromAddr> <txBytes>`
> Sends presigned transactin thrugh tendermint node. 
>
> Arguments:
> - `<fromAddr>`: Sender address.
> - <txBytes>: Encoded and signed byte representation of the tx

- `pocket accounts create-multi-public  <hex-pubkeys>`
> Sends `<amount>` POKT `<fromAddr>` to `<toAddr>`. Prompts the user for `<fromAddr>` account passphrase.
>
> Arguments:
> - `<hex-pubkeys>`: ordered comma separated keys. WARNING: changing the order creates a different addressr.

- `pocket accounts build-MS-Tx  <signer-address> <json-message> <hex-pubkeys> <chainID> <fee> <legacyCodec=(true|false)>`
>  Build and sign a multisignature transaction from scratchL result is hex encoded std transaction object.
>
> Arguments:
> = `<signer-address>`: Address  building & signing.
> = `<json-message>`: Message structure for the transaction.
> - `<hex-pubkeys>`: ordered comma separated keys. WARNING: changing the order creates a different addressr.
> - `<chainID>`: The pocket chain identifier.
> - `<fee>`:  An amount of POKT for the network.
> - `<legacyCodec>`: Enlble/Disable amino encoding for transaction.

- `pocket accounts sign-ms-tx  <signer-address> <hex-stdtx> <hex-pubkeys> <chainID> <fee> <legacyCodec=(true|false)>`
>  Sign a multisignature transaction using public keys, and the transaction object out of order. result is hex encoded standard transaction object.
>
> Arguments:
> = `<signer-address>`: Address  building & signing.
> = `<hex-stdtx>`: Prebuilt hexadecimal standard transaction.
> - `<hex-pubkeys>`: ordered comma separated keys. WARNING: changing the order creates a different addressr.
> - `<chainID>`: The pocket chain identifier.
> - `<fee>`:  An amount of POKT for the network.
> - `<legacyCodec>`: Enlble/Disable amino encoding for transaction.

- `pocket accounts sign-ms-next <signer-address> <hex-stdtx> <hex-pubkeys> <chainID> <fee> <legacyCodec=(true|false)>`
>  Sign a multisignature transaction object, result is hex encoded standard transaction object
>  WARNING: signer addres MUST be the next signer (in order of public keys in the multisignature) or the signature will be invalid
>
> Arguments:
> = `<signer-address>`: Address  building & signing.
> = `<hex-stdtx>`: Prebuilt hexadecimal standard transaction.
> - `<hex-pubkeys>`: ordered comma separated keys. WARNING: changing the order creates a different addressr.
> - `<chainID>`: The pocket chain identifier.
> - `<fee>`:  An amount of POKT for the network.
> - `<legacyCodec>`: Enlble/Disable amino encoding for transaction.
