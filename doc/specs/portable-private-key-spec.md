# Portable Private Key
## Version 0.0.1

### Overview
The Pocket Network protocol Portable Private Key.

This specification will serve to describe the Portable Private Key (PPK)
enabling multiple use cases like the creation of wallets.

### Design Basics

Portable Private Key (PPK) design borrows ideas from GPG ASCII Armored previously used by tendermint keys
and mix them with the JSON format that is lightweight and easy readable by humans.

Ensuring a safe and portable store for private key.

### PPK Example
````json
{
 "kdf":"scrypt",
 "salt":"8AA85775977952115075E68278C070A6",
 "secparam":"12",
 "hint":"",
 "ciphertext":"3oGW8vJfpEtW57XQ4AB+wdHfcPGdJb266eD8RMoJ3EAb2bgnUSyxV4oHYtnXoqEQY6kxb9+hB1tvA5TMacYCRZOEDA4Ml0fevUvh2oRTwVE="
}
````

### Data Structure Schema
PPK Structure as a JSON object holds metadata that enables the system to safely decrypt a password protected private key.
PPK must contain the following fields:

#### kdf
> type: `string`
>
> The KDF used for password-based symmetric encryption.
>
> Currently pocket is using scrypt as a sole kdf

#### salt
> type: `string`
>
> The salt used for symmetric encryption/decryption.

#### secparam
> type: `string`
>
> The secParam used at encryption time

#### hint
> type: `string`
>
> An optional hint that was input at the moment of exporting the private key.

#### ciphertext
> type: `string`
>
> Your Pocket private key encrypted and armored using ASCII;

### Aditional Elements

#### Scrypt params
>
>	N    = 32768
> 	r    = 8
> 	p    = 1
> 	keylen = 32
>

#### Symmetric cipher
>
> AES-256-GCM
>
#### Symetric Cipher Params

> nonce = first 12 bytes from decryption key

### How to export a private key using PPK format

1) get the key from the keybase or generate a new one.
2) using a random `salt` and a desired `password` we generate a `key` using Scrypt with the specified params
3) store the used `salt` encoded in hex(base16) used as **salt** on the PPK
5) use the `key` to generate a `nonce` or `iv` consisting of the first `secparam` bytes of the `key`
6) Encrypt the raw bytes from the private key using AES-256-GCM with the `key` and the `nonce`
7) Using base64 encoding, "Armor" the encrypted bytes and store the string as the **ciphertext**
8) store an optional `hint` as a reminder for the `password` used to encrypt
9) create the JSON struct using the values stored, **kdf** value should be "scrypt" and **secparam** currently is `12`
10) store it on a file


### How to import a private key using PPK format

1) Unmarshall or read the PPK JSON file to be able to use the values stored.
2) Validate the PPK passes these validations
    - **kdf** value equals to `"scrypt"`
    - **salt** value is not `""` (empty)
    - **salt** value can be decoded from `hex(base16)`
    - **ciphertext** can be decoded from `base64`
3) Using base64 decoding, "Unarmor" the **ciphertext** armored string and store it as the `encryptedBytes`
4) using the decoded **salt** value, the **secparam** as `cost` and the encryption `password` we generate a `key` using scrypt
6) use the `key` to generate a `nonce` consisting of the first **secparam** bytes of the `key`
7) Decrypt the `encryptedBytes` from 3)  using AES-256-GCM with the `key` and the `nonce`
8) if the password used was the correct one you should have the decrypted private key bytes if not, the password was wrong and should retry from step 4)
