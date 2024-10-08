# Pocket TestNet <!-- omit in toc -->

## Table of Contents <!-- omit in toc -->

- [Validator List w/ Pocket RPC](#validator-list-w-pocket-rpc)
  - [Example Queries](#example-queries)
    - [Querying State](#querying-state)
    - [Querying Binary Version](#querying-binary-version)
- [Validators w/ Tendermint RPC](#validators-w-tendermint-rpc)
  - [RPC \& Tendermint Public Endpoints](#rpc--tendermint-public-endpoints)
  - [Querying Net Info](#querying-net-info)
  - [Querying Status](#querying-status)
- [TestNet Seeds](#testnet-seeds)

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

### RPC & Tendermint Public Endpoints

`breezytm | Stakenodes (277262895459336194)` has made his Validator endpoints available
behind a public load balancer here:

- **RPC**: [rpc.testnet.pokt.network/lb/6d6f727365/](https://rpc.testnet.pokt.network/lb/6d6f727365/)
- **Tendermint**: [rpc.testnet.pokt.network/lb/6d6f727365/](https://tendermint.testnet.pokt.network/lb/6d6f727365)

### Querying Net Info

```bash
curl https://tendermint.testnet.pokt.network/lb/6d6f727365//net_info
```

### Querying Status

```bash
curl https://tendermint.testnet.pokt.network/lb/6d6f727365//status
```

## TestNet Seeds

The following seeds can be used to sync with TestNet. Copy-paste the following list of seeds into the `config.json` file on the `seeds` variable:

```bash
d90094952a3a67a99243cca645cdd5bd55fe8d27@seed1.testnet.pokt.network:26668, 2a5258dcdbaa5ca6fd882451f5a725587427a793@seed2.testnet.pokt.network:26669, a37baa84a53f2aab1243986c1cd4eff1591e50d0@seed3.testnet.pokt.network:26668, fb18401cf435bd24a2e8bf75ea7041afcf122acf@seed4.testnet.pokt.network:26669
```
