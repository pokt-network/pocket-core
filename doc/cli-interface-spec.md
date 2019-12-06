# Pocket Core CLI Interface Specification
## Version 0.0.1

### Overview
This document serves as a specification for the Command Line Interface of the Pocket Core application. There's no protocol verification for these commands, however because they map closely to protocol functions.

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

### Accounts Namespace Functions

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
KeyPair 0x... deleted succesfully.
```

- `pocket accounts update-passphrase <address>`
> Updates the passphrase for the indicated account. Will prompt the user for the current account passphrase and the new account passphrase.
>
> Arguments:
> - `<address>`: The address to be deleted.
> Example output:
```
KeyPair 0x... passphrase updated succesfully.
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
Account generated succesfully.
Address: 0x....
```

- `pocket accounts import <mnemonic>`
> Imports an account using the provided `<mnemonic>`. Will prompt the user for a [BIP-0039](https://github.com/bitcoin/bips/blob/master/bip-0039.mediawiki) password for the imported mnemonic and for a passphrase to encrypt the generated keypair.
>
> Arguments:
> - `<mnemonic>`: The mnemonic of the account to be imported.
> Example output:
```
Account imported succesfully.
Address: 0x....
```

- `pocket accounts import-armored <armor>`
> Imports an account using the Encrypted ASCII armored `<armor>` string. Will prompt the user for a decryption passphrase of the `<armor>` string and for an encryption passphrase to store in the Keybase.
>
> Arguments:
> - `<armor>`: The encrypted encoded private key to be imported.
> Example output:
```
Account imported succesfully.
Address: 0x....
```

- `pocket accounts export <address>`
> Exports the account with `<address>`, encrypted and ASCII armored. Will prompt the user for the account passphrase and for an encryption passphrase for the exported account.
>
> Arguments:
> - `<address>`: The address of the account to be exported.
> Example output:
```
Exported account: <armored string>
```

- `pocket accounts export-raw <address>`
> Exports the raw private key in hex format. Will prompt the user for the account passphrase. *NOTE: THIS METHOD IS NOT RECOMMENDED FOR SECURITY REASONS, USE AT YOUR OWN RISK.*
>
> Arguments:
> - `<address>`: The address of the account to be exported.
> Example output:
```
Exported account: 0x...
```




