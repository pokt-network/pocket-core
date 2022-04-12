# 0.8.X Changes & Features Description:

## Nodes Namespace CLI changes

For 0.8.X, in preparation for the upcoming consensus rule change that will enable non-custodial staking, we have some
changes on our CLI commands.

#### Regular Stake command has changed (Custodial)

from

```text
pocket nodes stake <fromAddr> <amount> <relayChainIDs> <serviceURI> <networkID> <fee>
```

to

```text
pocket nodes stake custodial <fromAddr> <amount> <relayChainIDs> <serviceURI> <networkID> <fee> <isBefore8.0>
```

The updated command expects 1 new parameter:

* `<isBefore8.0>`:  true or false depending if non-custodial upgrade is activated.

Before the upgrade is activated be sure to use `true` for `<isBefore8.0>` or the transaction won't go through

#### New non-custodial stake command (0.8.X)

```text
pocket nodes stake non-custodial <operatorPublicKey> <outputAddress> <amount> <RelayChainIDs> <serviceURI> <networkID> <fee> <isBefore8.0>
```

The new stake command expects 3 new parameters:

* `<operatorPublicKey>`: operatorAddress is the only valid signer for blocks & relays. Should match the running node.
* `<outputAddress>`: outputAddress is where reward and staked funds are directed.
* `<isBefore8.0>`:  true or false depending if non custodial upgrade is activated.

Replacing the `<fromAddr>` we have the `<operatorPublicKey>`, notice the change from using an address to using the
public key of the node. Also, we have the `<outputAddress>`, where rewards and funds will be delivered after the update
is activated. Keep in mind that even if you are using the command before the upgrade is activated you are required to
complete all the parameter.

Before the upgrade is activated be sure to use `true` for `<isBefore8.0>` or the transaction won't go through

**NOTE**: Even if the command is available, Non-custodial won't work until the consensus rule change is activated , we
recommend using `pocket nodes stake custodial [...] true` before the upgrade is activated.

## Chains.json hot reload + DIY Bonus Endpoint

### Hot reloading

For 0.8.X we added a feature to allow the reloading of the chains.json without restarting your node.

To enable the feature be sure to either run `pocket util update-configs` to have the new parameter added to your
config.json automatically (must edit manually afterwards to turn on as feature ships disabled)
or manually add the new key `"chains_hot_reload": true` to the `"pocket_config"` inside your config.json

After starting your node with this value set to `true` pocket will reload the information from the chains.json every
minute.

### Endpoint

And for those of you that need more control or automations, we have also enabled a protected endpoint that allows you to
check or update your hosted chains.

To use the endpoint:

* First make sure to set `"chains_hot_reload": false` before beginning, hot reloading and the update endpoint can't be
  used at the same time (enabling hot reload disables the update endpoint).
* Locate your `auth.json` file in the config directory, this file contains your `authtoken` you need its value to send
  request to protected endpoints (../v1/private/..).
    ```json
    {
    "Value": "<TOKEN>",
    "Issued": "2022-04-08T07:35:35.858373-04:00"
    }
    ```
* Copy the token value and replace the placeholder `<TOKEN>` with its value on this example call
  ```text
    curl --location --request POST 'http://localhost:8081/v1/private/updatechains?authtoken=<TOKEN>' \
    --header 'Content-Type: application/json' \
    --data-raw '[
    {
    "id": "0001",
    "url": "http://localhost:8081",
    "basic_auth": {
    "username": "",
    "password": ""
    }
    }
    ]'
    ```
* Check If your request is successful using this other call
  ```
  curl --location --request POST 'http://localhost:8081/v1/private/chains?authtoken=<TOKEN>'
  ```
  you should receive a response like this :
  ```json
  {
    "0001": {
        "basic_auth": {
            "password": "",
            "username": ""
        },
        "id": "0001",
        "url": "http://localhost:8081"
    }
  }
    ```
* Remember to manually update your chains.json with the desired changes before or after using this method as any changes
  done using the `updatechains` endpoint will be overwritten at restart by loading the chains.json
* Also, if you are writing any automations remember the `authtoken` is recreated on every restart.
