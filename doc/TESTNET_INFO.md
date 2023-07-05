# Pocket TestNet information overview

## Pocket Testnet validators list with pocket RPC

```
https://node1.testnet.pokt.network/
https://node2.testnet.pokt.network/
https://node3.testnet.pokt.network/
https://node4.testnet.pokt.network/
https://node5.testnet.pokt.network/
https://node6.testnet.pokt.network/
```

Also, other queries can be done querying RPC nodes. Here an example for downloading TestNet's state backup as a local file through terminal:

```bash
curl -X POST https://node1.testnet.pokt.network/v1/query/state > query_state.json
```

Or checking version of the pocket binary on the testnet validators

```bash
curl https://node1.testnet.pokt.network/v1
```


## Pocket Testnet validators list with Tendermint RPC 

This info can be checked querying these links (Note that corresponding credentials are needed to accessing this info):

```
https://node1.tendermint.testnet.pokt.network/
https://node2.tendermint.testnet.pokt.network/
https://node3.tendermint.testnet.pokt.network/
https://node4.tendermint.testnet.pokt.network/
https://node5.tendermint.testnet.pokt.network/
https://node6.tendermint.testnet.pokt.network/
```

Ask for the TestNet credentials in Pocket's Node-Chat Discord channel [here](https://discord.com/channels/553741558869131266/564836328202567725)

For example, accessing tendermint's endpoints for network information:

```bash
curl https://node1.tendermint.testnet.pokt.network/net_info
```

Verifying status of testnet node1 via tendermint:

```bash
curl https://node1.tendermint.testnet.pokt.network/status
```


## Pocket TestNet Seeds


If you want to sync with testnet, you only need to copy/paste the list of nodes below on the `config.json` file on the variable `Seeds`

```
d90094952a3a67a99243cca645cdd5bd55fe8d27@seed1.testnet.pokt.network:26668, 2a5258dcdbaa5ca6fd882451f5a725587427a793@seed2.testnet.pokt.network:26669, a37baa84a53f2aab1243986c1cd4eff1591e50d0@seed3.testnet.pokt.network:26668, fb18401cf435bd24a2e8bf75ea7041afcf122acf@seed4.testnet.pokt.network:26669
```


## Pocket TestNet metrics dashboard

TestNet metrics can be check following these links:

**Testnet Loadbalancer metrics  (Network traffic dashboard)**
>https://monitoring.nodefleet.net/d/O23g2BeWk/testnet-loadbalancer-metrics?orgId=4&var-service=testnet1@file&var-entrypoint=All&from=now-3h&to=now&refresh=5m

**Tendermint metrics (Consensus, Blocks, Transactions dashboard information and so on)**
>https://monitoring.nodefleet.net/d/UJyurCTWz/testnet-validators-tendermint-metrics

**Node exporter metrics (Instance metrics)**
>https://monitoring.nodefleet.net/d/Gm5yJc94z/testnet-validators-telegraf-metrics

**Loki dashboard (Testnet Logs and explorer search)**
>https://monitoring.nodefleet.net/d/_j0yAcrVz/testnet-validators-loki

Note corresponding credentials will be needed to accessing this info. Ask for credentials in Pocket's Node-Chat Discord channel [here](https://discord.com/channels/553741558869131266/564836328202567725)

### Who's supporting Pocket's TestNet infrastructure?

The support for Testnet infrastructure is given by Nodefleet.org, which is a Web3 blockchain and node running company focused on delivering value for investors/builders on multi-chain ecosystems providing top quality engineering and quality infrastructure around all of its products. Feel free and ask anything about TestNet [here](https://discord.com/channels/553741558869131266/564836328202567725) tagging **Lowell | nodefleet.org** and **Steven94 | nodefleet.org**