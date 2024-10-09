# Pocket TestNet <!-- omit in toc -->

## Table of Contents <!-- omit in toc -->

- [RPC \& Tendermint Public Endpoints](#rpc--tendermint-public-endpoints)
  - [BreezyTm (StakeNodes)](#breezytm-stakenodes)
  - [Ian (Cryptonode.tools)](#ian-cryptonodetools)
- [Example Queries](#example-queries)
  - [Example Query - Net Info](#example-query---net-info)
  - [Example Query - Status](#example-query---status)
  - [Example Query - State](#example-query---state)
  - [Example Query - Binary Version](#example-query---binary-version)
- [TestNet Seeds](#testnet-seeds)

## RPC & Tendermint Public Endpoints

### BreezyTm (StakeNodes)

`breezytm | Stakenodes (277262895459336194)` has made his Validator's RPC + Tendermint
endpoints available behind a public load balancer here:

- **RPC**: [rpc.testnet.pokt.network/lb/6d6f727365](https://rpc.testnet.pokt.network/lb/6d6f727365/)
- **Tendermint**: [rpc.testnet.pokt.network/lb/6d6f727365](https://tendermint.testnet.pokt.network/lb/6d6f727365)

### Ian (Cryptonode.tools)

`Ian | cryptonode.tools (693644362575511573)` also made his Validator's TM endpoint available here

- **Tendermint**: [https://morse-tendermint.chains-eu6.cryptonode.tools/](https://morse-tendermint.chains-eu6.cryptonode.tools/)

## Example Queries

### Example Query - Net Info

```bash
curl https://tendermint.testnet.pokt.network/lb/6d6f727365/net_info
```

### Example Query - Status

```bash
curl https://tendermint.testnet.pokt.network/lb/6d6f727365/status
```

### Example Query - State

```bash
curl -X POST https://rpc.testnet.pokt.network/lb/6d6f727365/v1/query/state | tee query_state.json | jq
```

### Example Query - Binary Version

```bash
curl https://rpc.testnet.pokt.network/lb/6d6f727365/v1
```

## TestNet Seeds

The following seeds can be used to sync with TestNet. Copy-paste the following list of seeds into the `config.json` file on the `seeds` variable:

```bash
b3d86cd8ab4aa0cb9861cb795d8d154e685a94cf@seed1.testnet.pokt.network:26663,5b0107a5252f6a037eed7f5c24a7d916e4dd93bd@testnet_seed_4.cryptonode.tools:16646
```
