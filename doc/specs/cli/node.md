# Node Namespace
Functions for Node management.

- `pocket node stake <fromAddr> <amount> <relayChainIDs> <serviceURI> <chainID> <fee> <legacyCodec=(true | false)>`
> Stakes the Node into the network, making it available for service. Prompts the user for the `<fromAddr>` account passphrase.
> After the 0.6.X upgrade, if the node is already staked, this transaction acts as an *update* transaction.
> A node can updated relayChainIDs, serviceURI, and raise the stake amount with this transaction.
> If the node is currently staked at X and you submit an update with new stake Y. Only Y-X will be subtracted from an account
> If no changes are desired for the parameter, just enter the current param value just as before
> Arguments:
> - `<fromAddr>`: Target Address to stake.
> - `<amount>`: The amount of uPOKT to stake. Must be higher than the current minimum amount of Node Stake parameter.
> - `<relayChainIDs>`: A comma separated list of chain Network Identifiers.
> - `<serviceURI>`: The Service URI Applications will use to communicate with Nodes for Relays.
> - `<chainID>`: The pocket chain identifier.
> - `<fee>`:  An amount of uPOKT for the network.
> - `<legacyCodec>`: Enlble/Disable amino encoding for transaction.
> Example output:
```
Transaction submitted with hash: <Transaction Hash>
```

- `pocket node unstake <fromAddr> <chainID> <fee> <legacyCodec=(true | false)>`
> Unstakes a Node from the network, changing its status to `Unstaking`. Prompts the user for the `<fromAddr>` account passphrase.
>
> Arguments:
> - `<fromAddr>`: Target staked address.
> - `<chainID>`: The pocket chain identifier.
> - `<fee>`:  An amount of uPOKT for the network.
> - `<legacyCodec>`: Enlble/Disable amino encoding for transaction.
> Example output:
```
Transaction submitted with hash: <Transaction Hash>
```

- `pocket node unjail <fromAddr> <chainID> <fee> <legacyCodec=(true | false)>`
> Unjails a Node from the network, allowing it to participate in service and consensus again. Prompts the user for the `<fromAddr>` account passphrase.
>
> Arguments:
> - `<fromAddr>`: Target jailed address.
> - `<chainID>`: The pocket chain identifier.
> - `<fee>`:  An amount of uPOKT for the network.
> - `<legacyCodec>`: Enlble/Disable amino encoding for transaction.
> Example output:
```
Transaction submitted with hash: <Transaction Hash>
```
