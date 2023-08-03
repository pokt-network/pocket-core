# Pocket TestNet <!-- omit in toc -->

## Table of Contents <!-- omit in toc -->

- [Validator List w/ Pocket RPC](#validator-list-w-pocket-rpc)
  - [Example Queries](#example-queries)
    - [Querying State](#querying-state)
    - [Querying Binary Version](#querying-binary-version)
- [Validators w/ Tendermint RPC](#validators-w-tendermint-rpc)
  - [Credentials](#credentials)
  - [Querying Net Info](#querying-net-info)
  - [Querying Status](#querying-status)
- [TestNet Seeds](#testnet-seeds)
- [Pocket TestNet metrics dashboard](#pocket-testnet-metrics-dashboard)
  - [Brought to you by NodeFleet](#brought-to-you-by-nodefleet)
- [Helper Examples](#helper-examples)
  - [View All Validator Tendermint Versions](#view-all-validator-tendermint-versions)
  - [View all Validator Binary Versions](#view-all-validator-binary-versions)
  - [View All Validator Heights](#view-all-validator-heights)

## Validator List w/ Pocket RPC

```bash
https://node1.testnet.pokt.network/
https://node2.testnet.pokt.network/
https://node3.testnet.pokt.network/
https://node4.testnet.pokt.network/
https://node5.testnet.pokt.network/
https://node6.testnet.pokt.network/
```

### Example Queries

#### Querying State

```bash
curl -X POST https://node1.testnet.pokt.network/v1/query/state | tee query_state.json | jq
```

Output:

```bash
# {
#   "app_hash": "",
#   "app_state": {
#     "application": {
#       "applications": [
#         {
#           "address": "065013157ffb401642d0418b408474b361ee0836",
#           "chains": [
#             "004A",
#             "004B",
# ...
```

#### Querying Binary Version

```bash
curl https://node1.testnet.pokt.network/v1
```

Output:

```bash
#"BETA-0.10.2"%
```

## Validators w/ Tendermint RPC

```bash
https://node1.tendermint.testnet.pokt.network/
https://node2.tendermint.testnet.pokt.network/
https://node3.tendermint.testnet.pokt.network/
https://node4.tendermint.testnet.pokt.network/
https://node5.tendermint.testnet.pokt.network/
https://node6.tendermint.testnet.pokt.network/
```

### Credentials

These endpoints require authentication and can be accessed at `https://testnet:${NODE_FLEET_PASSWORD}@$node`

You can request TestNet credentials in the [Pocket Node-Chat Discord channel](https://discord.com/channels/553741558869131266/564836328202567725).

### Querying Net Info

```bash
curl https://node1.tendermint.testnet.pokt.network/net_info
```

### Querying Status

```bash
curl https://node1.tendermint.testnet.pokt.network/status
```

## TestNet Seeds

The following seeds can be used to sync with TestNet. Copy-paste the following list of seeds into the `config.json` file on the `seeds` variable:

```bash
d90094952a3a67a99243cca645cdd5bd55fe8d27@seed1.testnet.pokt.network:26668, 2a5258dcdbaa5ca6fd882451f5a725587427a793@seed2.testnet.pokt.network:26669, a37baa84a53f2aab1243986c1cd4eff1591e50d0@seed3.testnet.pokt.network:26668, fb18401cf435bd24a2e8bf75ea7041afcf122acf@seed4.testnet.pokt.network:26669
```

## Pocket TestNet metrics dashboard

TestNet metrics can be viewed at the following links:

- [Loadbalancer metrics](https://monitoring.nodefleet.net/d/O23g2BeWk/testnet-loadbalancer-metrics?orgId=4&var-service=testnet1@file&var-entrypoint=All&from=now-3h&to=now&refresh=5m): Network traffic dashboard
- [Tendermint metrics](https://monitoring.nodefleet.net/d/UJyurCTWz/testnet-validators-tendermint-metrics): Consensus, Blocks, Transactions dashboard information and so on
- [Node exporter metrics](https://monitoring.nodefleet.net/d/Gm5yJc94z/testnet-validators-telegraf-metrics): Instance metrics
- [Loki dashboard](https://monitoring.nodefleet.net/d/_j0yAcrVz/testnet-validators-loki): Testnet Logs and explorer search

You can request TestNet credentials in the [Pocket Node-Chat Discord channel](https://discord.com/channels/553741558869131266/564836328202567725).

### Brought to you by NodeFleet

The support for Testnet infrastructure is given by[nodefleet.org](https://nodefleet.org/), a Web3 blockchain and node running company focused on delivering value for investors &builders on multi-chain ecosystem. Nodefleet provides top quality engineering and quality infrastructure around all of its products.

Reach out to the team about TestNet directly on [Discord](https://discord.com/channels/553741558869131266/564836328202567725) tagging **Lowell | nodefleet.org#7301** (148983981134577665) and **Steven94 | nodefleet.org** (357688204566069248).

## Helper Examples

### View All Validator Tendermint Versions

_Important: You need to expose `NODE_FLEET_PASSWORD`_

```bash
#!/bin/bash
declare -a nodes=("node1.tendermint.testnet.pokt.network/status"
                  "node2.tendermint.testnet.pokt.network/status"
                  "node3.tendermint.testnet.pokt.network/status"
                  "node4.tendermint.testnet.pokt.network/status"
                  "node5.tendermint.testnet.pokt.network/status"
                  "node6.tendermint.testnet.pokt.network/status")

for node in "${nodes[@]}"
do
    url="https://testnet:${NODE_FLEET_PASSWORD}@$node"
    curl -s $url | jq '.result.node_info.version'
done
```

### View all Validator Binary Versions

```bash
#!/bin/bash
declare -a nodes=("node1.testnet.pokt.network/v1"
                  "node2.testnet.pokt.network/v1"
                  "node3.testnet.pokt.network/v1"
                  "node4.testnet.pokt.network/v1"
                  "node5.testnet.pokt.network/v1"
                  "node6.testnet.pokt.network/v1")

for node in "${nodes[@]}"
do
    url="https://$node"
    curl -X GET $url
    echo ""
done
```

### View All Validator Heights

```bash
#!/bin/bash
declare -a nodes=("node1.testnet.pokt.network/v1/query/height"
                  "node2.testnet.pokt.network/v1/query/height"
                  "node3.testnet.pokt.network/v1/query/height"
                  "node4.testnet.pokt.network/v1/query/height"
                  "node5.testnet.pokt.network/v1/query/height"
                  "node6.testnet.pokt.network/v1/query/height")

for node in "${nodes[@]}"
do
    url="https://$node"
    curl -X POST $url
done
```
