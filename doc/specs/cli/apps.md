# Apps Namespace

Functions for Application management.

* `pocket app stake <fromAddr> <amount> <relayChainIDs> <chainID> <fee> <legacyCodec=(true | false)>`

  > Stakes the Application into the network, making it available to receive service. Prompts the user for the `<fromAddr>` account passphrase. After the 0.6.X upgrade, if the app is already staked, this transaction acts as an _update_ transaction. A app can updated relayChainIDs, and raise the stake/max\_relays amount with this transaction. If the app is currently staked at X and you submit an update with new stake Y. Only Y-X will be subtracted from an account If no changes are desired for the parameter, just enter the current param value just as before Arguments:
  >
  > * `<fromAddr>`: Target Address to stake.
  > * `<amount>`: The amount of uPOKT to stake. Must be higher than the current minimum amount of Node Stake parameter.
  > * `<relayChainIDs>`: A comma separated list of chain Network Identifiers.
  > * `<chainID>`: The pocket chain identifier.
  > * `<fee>`:  An amount of uPOKT for the network.
  > * `<legacyCodec>`: Enlble/Disable amino encoding for transaction.
  >
  >   Example output:
  >
  >   ```text
  >   Transaction submitted with hash: <Transaction Hash>
  >   ```

* `pocket app unstake <fromAddr> <chainID> <fee> <legacyCodec=(true | false)>`

  > Unstakes an Application from the network, changing its status to `Unstaking`. Prompts the user for the `<fromAddr>` account passphrase.
  >
  > Arguments:
  >
  > * `<fromAddr>`: The address of the sender.
  > * `<chainID>`: The pocket chain identifier
  >
  >   Example output:
  >
  >   ```text
  >   Transaction submitted with hash: <Transaction Hash>
  >   ```

* `pocket app create-aat <appAddr> <clientPubKey>`

  > Creates a signed application authentication token \(version `0.0.1` of the AAT spec\), that can be embedded into application software for Relay servicing. Will prompt the user for the `<appAddr>` account passphrase. Read the Application Authentication Token documentation [here](https://github.com/pokt-network/pocket-core/tree/18c6dc91425e4e058b5edfbb6d5fe7b7ddb4ed3c/doc/specs/cli/application-auth-token.md). _**NOTE**_: USE THIS METHOD AT YOUR OWN RISK. READ THE APPLICATION SECURITY GUIDELINES TO UNDERSTAND WHAT'S THE RECOMMENDED AAT CONFIGURATION FOR YOUR APPLICATION:
  >
  > Arguments:
  >
  > * `<appAddr>`: The address of the Application account to use to produce this AAT.
  > * `<clientPubKey>`: The account public key of the client that will be signing and sending Relays sent to the Pocket Network.
  >
  >   Example output:
  >
  >   ```javascript
  >   {
  >   "version": "0.0.1",
  >   "applicationPublicKey": "0x...",
  >   "clientPublicKey": "0x...",
  >   "signature": "0x..."
  >   }
  >   ```

