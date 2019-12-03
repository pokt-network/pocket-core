# Application Auth Token
## Version 0.0.1

### Overview
The Pocket Network protocol contemplates the use of Application Auth Tokens to allow Application Clients to access Service Nodes on behalf of the Application.

This specification will serve to describe the AAT system attributes such as:

- Data Structure Schema
- Encoding/Decoding

### Data Structure Schema
An AAT must contain the following fields:

#### version
> type: `string`
>
> A semver string specifying the spec version under which this ATT needs to be interpreted.

#### signature
> type: `string`
>
> The application will sign a hash of the `message` property within this token with the specified `appPubKey` and corresponding private key.

#### applicationPublicKey
> type: `string`
>
> The hexadecimal publicKey of the Application

#### clientPublicKey
> type: `string`
>
> Required for signature verification, the hexadecimal public of each individual client allowing for granular control of who can use the ATT

### Encoding/Decoding
An ATT token will always be represented using the [Amino encoding.](https://github.com/tendermint/go-amino)

### ECDSA ed25519 Signature
The protocol wide ed25519 ECDSA will be used for any signatures and verifications that are used within this specification.

The proper way to sign the token is as follows:

1. JSON Encode AAT with an empty string signature field:
````


````
