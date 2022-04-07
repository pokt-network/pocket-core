---
description: >- Version 0.0.1. The Pocket Network protocol contemplates the use of Application Auth Tokens to allow
Application Clients to access Service Nodes on behalf of the Application.
---

# Application Authentication Token

## Data Structure Schema

An AAT must contain the following fields:

### version

> type: `string`

A semver string specifying the spec version under which this AAT needs to be interpreted.

### signature

> type: `string`

The application will sign a hash of the `message` property within this token with the specified `appPubKey` and
corresponding private key.

### applicationPublicKey

> type: `string`

The hexadecimal publicKey of the Application

### clientPublicKey

> type: `string`

Required for signature verification, the hexadecimal public of each individual client allowing for granular control of
who can use the AAT

## ECDSA ed25519 Signature Scheme

The protocol wide ed25519 ECDSA will be used for any signatures and verifications that are used within this
specification.

The proper way to sign the token is as follows:

1. JSON Encode AAT with an empty string signature field:
2. SHA3 \(256\) Hash the json bytes
3. Sign with ed25519 ECDSA
4. HexEncode the result bytes into a string

```text
AAT {
    ApplicationSignature: "",
    ApplicationPublicKey: a.ApplicationPublicKey,
    ClientPublicKey:      a.ClientPublicKey,
    Version:              a.Version,
}
```

`AATBytes = JSON.Encode(AAT)`

`Message = SHA3-256(AATBytes)`

`AAT.Signature = ED25519.Sign(Message)`

