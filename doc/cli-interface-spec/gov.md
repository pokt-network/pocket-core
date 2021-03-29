### Governance Namespace Functions
The gov namespace handles all governance related interactions, from DAOTransfer, change parameters; to performing protocol Upgrades.


- `pocket gov transfer <amount> <fromAddr> <toAddr> <chainID> <fee> <legacyCodec=(true | false)>`
> if authorized, move funds from the DAO.
>
> Arguments:
> - `<amount>`: The amount of uPOKT to be sent.
> - `<toAddr>`: Recipient address for the transaction.
> - `<fromAddr>`: Sender address.
> - `<chainID>`: The pocket chain identifier.
> - `<fee>`:  An amount of uPOKT for the network .
> - `<legacyCodec>`: Enlble/Disable amino encoding for transaction.
> Example output:
```
Transaction submitted with hash: <Transaction Hash>
```

- `pocket gov burn <amount> <fromAddr> <toAddr> <chainID> <fee> <legacyCodec=(true | false)`
> if authorized, burn funds from the DAO.
>
> Arguments:
> - `<amount>`: The amount of uPOKT to be sent.
> - `<toAddr>`: Recipient address for the transaction.
> - `<fromAddr>`: Sender address.
> - `<chainID>`: The pocket chain identifier.
> - `<fee>`:  An amount of uPOKT for the network .
> - `<legacyCodec>`: Enlble/Disable amino encoding for transaction.
> Example output:
```
Transaction submitted with hash: <Transaction Hash>
```

- `pocket gov change_param  <fromAddr> <chainID> <paramKey module/param> <paramValue (jsonObj)> <fee>  <legacyCodec=(true | false)`
> if authorized,submit a tx to change any param from any module. Will pormt the user foor the <fromAddr> account passphrase
>
> Arguments:
> - `<fromAddr>`: Sender address.
> - `<chainID>`: The pocket chain identifier.
> - `<paramKey>`: Target parameter key to change in format module/param. i.e: pos/ProposerPercentage
> - `<paramValue>: New value for key.`
> - `<fee>`:  An amount of uPOKT for the network .
> - `<legacyCodec>`: Enlble/Disable amino encoding for transaction.
> Example output:
```
Transaction submitted with hash: <Transaction Hash>
```

- `pocket gov upgrade  <fromAddr> <atHeight> <version> <chainID> <fees> <legacyCodec=(true | false)`
> if authorized, upgrade the protocol. Will prompt the user for the <fromAddr> account passphrase.
>
> Arguments:
> - `<fromAddr>`: Sender address.
> - `<atHeight>`: Target height for upgrade
> - `<version>`: Target version for protocl update.
> - `<chainID>`: The pocket chain identifier.
> - `<fee>`:  An amount of uPOKT for the network .
> - `<legacyCodec>`: Enlble/Disable amino encoding for transaction.
> Example output:
```
Transaction submitted with hash: <Transaction Hash>
```