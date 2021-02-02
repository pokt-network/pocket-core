# Node Namespace
Functions for Node management.

- `pocket node stake <fromAddr> <amount> <chains> <serviceURI> <chainID> <fee> <legacyCodec=(true | false)>`
> Stakes the Node into the network, making it available for service. Prompts the user for the `<fromAddr>` account passphrase.
>
> Arguments:
> - `<fromAddr>`: Target Address to stake.
> - `<amount>`: The amount of POKT to stake. Must be higher than the current minimum amount of Node Stake parameter.
> - `<chains>`: A comma separated list of chain Network Identifiers.
> - `<serviceURI>`: The Service URI Applications will use to communicate with Nodes for Relays.
> - `<chainID>`: The pocket chain identifier.
> - `<fee>`:  An amount of POKT for the network.
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
> - `<fee>`:  An amount of POKT for the network.
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
> - `<fee>`:  An amount of POKT for the network.
> - `<legacyCodec>`: Enlble/Disable amino encoding for transaction.
> Example output:
```
Transaction submitted with hash: <Transaction Hash>
```