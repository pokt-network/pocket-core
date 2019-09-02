# Application Auth Token
## Version 0.0.1

### Overview
The Pocket Network protocol contemplates the use of Application Auth Tokens to allow Application Clients to access Service Nodes on behalf of the Application.

This specification will serve to describe the AAT system attributes such as:

- Data Structure Schema
- Encoding/Decoding

### Data Structure Schema
An AAT must contain the following fields:

#### version (required)
> type: `string`
>
> A semver string specifying the spec version under which this ATT needs to be interpreted.

#### message (required)
> type: `ATTMessage`
>
> The message which will be used by the Service Node to configure service to the bearer of this token.
> The schema of the message is represented below by the `AATMessage` data structure.

#### appPubKey (required)
> type: `string`
>
> The hexadecimal representation of the application public key which was used to sign the `message`

#### signature (required)
> type: `string`
>
> The application will sign a hash of the `message` property within this token with the specified `appPubKey` and corresponding private key.

### The `AATMessage` Structure Schema
The AAT data structure contains the information that will be first encoded to the Amino format and then signed using the protocol wide ECDSA. It contains the fields described below:

#### applicationAddress (required)
> type: `string`
>
> The hexadecimal address of the Application

#### clientAddress (optional)
> type: `string`
>
> Optionally for added security the application can sign the hexadecimal address of each individual client allowing for granular control of who can use the ATT

### Encoding/Decoding
An ATT token will always be represented using the [Amino encoding.](https://github.com/tendermint/go-amino)

### ECDSA Signature
The protocol wide ECDSA will be used for any signatures and verifications that are used within this specification.
