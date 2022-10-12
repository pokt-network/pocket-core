### Summary

This client allows all node-runners to take advantage of relay traffic in locations around the world without having a full Pocket node at these locations.

### Nomenclature

* Servicer Node: Current implementation of the Pocket Node (with all the features, even Lean Node)
* Mesh Node: Lightweight Pocket Node Proxy

### What it does?

This allows to receive relays in name of a **Servicer** in a different geolocation, serve it with the minimum validation and later notify the **Servicer** node.
At the end of the session process the **Servicer** node will post a single claim transaction for all the relays.

To achieve this you need to deploy your Servicer Node and your Mesh Node behind a proxy like you are currently doing to set up your SSL.
Also, you need to set both of them behind a Global DNS provider. The Global DNS will provide the IP address of the closest **Servicer** OR **Mesh** node, based on the request location.

### Features

* Relay Approach: First Free
  * this mean that mesh node will serve first relay for free and check the servicer notify response to understand it should keep servicing or not.
* Proxy any request to servicer
* Monitor Servicer health using /v1/health and in memory cron jobs
* Keep sessions in local cache
* Keep relays in local cache until they are notified
  * this ensures that even if the mesh node is restarted and the relay was not yet notified, they will be after node startup again
* Auto clean up old serve session from cache
* Servicer relays notification has a retry mechanism just in case the Servicer node is not responsive for a while.
  * this will only retry in some scenarios.
* Implements a worker queue for the notification to keep the process simple and secure event under crash circumstances.
* Handle minimum mesh node side validations using response get from Health monitor about Height, Starting and Catching Up (same of pocket node)
* Handle minimum relay validations about payload format (same of pocket node)

### Included Branches:

* [Lean Node](https://github.com/pokt-network/pocket-core/tree/ethereal-wombat)
* [Fix #1457 & Memory enhance](https://github.com/pokt-network/pocket-core/pull/1485)
* [Feature #1456](https://github.com/pokt-network/pocket-core/pull/1483)

### Hardware Requirements
(we are expecting your feedback on this!)
* CPU: 1 vcpu or less
* Memory: 200mb or less

### Software Requirements
* Reverse Proxy (SSL)
* Global DNS

### How to use it?

#### Pre-requisites:

This guide assume you already have a Servicer properly setup and running. If not, please refer to our docs to understand [how](cli/default.md)

1. Global DNS that handle your domain and forward to proper region node (servicer or mesh)
2. Chains on each region you want to deploy (servicer or mesh)
3. Servicer in one region
4. Mesh in the other N regions

#### Prepare Servicer:

1. (optional) Create auth.json if it does not exist on your pocket node.
```json
{
  "Value": "<SOME VALUE HERE>",
  "Issued": "2022-09-29T00:00:00.000000000-00:00"
}
```
2. Update your `config.json`
   1. add `mesh_node` option as `true` into the section `pocket_config`
   2. change `generate_token_on_start` option to `false` into the section `pocket_config`
3. If your proxy has all the endpoints closed except `/v1` and `/v1/client`, please add:
   1. `/v1/private/mesh/relay` - allow mesh node to notify about relays done
   2. `/v1/health` - return node status: version, height, starting, catching_up
4. Start your node as you were doing.

#### Setup Mesh Node:

1. Set up your proxy in the same way for a Servicer
2. Create the following `config.json` file inside the `--datadir` directory
```json
{
	"data_dir": "/home/app/.pocket/mesh",
    "rpc_port": "8081", // mesh node listening port
    "chains_name": "chains/chains.json", // chains for mesh node. This should be a filename path relative to --datadir
    "rpc_timeout": 30000, // chains rpc timeout
    "log_level": "*:info", // log level, you can try with *:error or even *:debug (this print a lot)
    "user_agent": "mesh-node", // https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/User-Agent
    "auth_token_file": "key/auth.json", // authtoken for mesh private endpoints. This should be a filename path relative to --datadir
    "json_sort_relay_responses": true,
    "pocket_prometheus_port": "8083",
    "prometheus_max_open_files": 3,
    "relay_cache_file": "data/relays.pkt",
    "session_cache_file": "data/session.pkt",
    // Worker options match with: https://github.com/alitto/pond#resizing-strategies
    // These are used for the relay report queue
    "worker_strategy": "balanced", // Kind of worker strategy, could be: balanced | eager | lazy
    "max_workers": 20,
    "max_workers_capacity": 100,
    "workers_idle_timeout": 1000,
    "servicer_private_key_file": "key/key.json", // servicer private key to sign proof message on relay response. This should be a filename path relative to --datadir
    "servicer_rpc_timeout": 30000, // servicer rpc timeout
    "servicer_auth_token_file": "key/auth.json", // authtoken used to call servicer. This should be a filename path relative to --datadir
    // Servicer relay notification has a retry mechanism, refer to: https://github.com/hashicorp/go-retryablehttp
    "servicer_retry_max_times": 10,
    "servicer_retry_wait_min": 10,
    "servicer_retry_wait_max": 180000
}
```
3. Create Servicer private key file with following format into the path you set on `config.json`:
```json
[
  {
	"priv_key": "aaabbbbcccccddd", // servicer private key
	"servicer_url": "http://localhost:8081" // servicer url/ip where mesh node can reach the servicer node to check health, proxy requests and notify relays
  },
  {
	... // add as much servicers as you need to handle with one single geo-mesh process.
  }
]

```
4. Create auth.json files into the path you set on `config.json` for `auth_token_file` and `servicer_auth_token_file`:
```json
{
  "Value": "<SOME VALUE HERE>",
  "Issued": "2022-09-29T00:00:00.000000000-00:00"
}
```
5. Create chains.json files into the path you set on `config.json`
   * IMPORTANT: this need to handle all the chains that servicer support.
   * NOTE: If you want to support a subset of chains in a region, you will need set the chains here too but point them to the closest chain you have.
6. Start your mesh node: `pocket start-mesh --datadir </your/path>`
7. Call your mesh node at `/v1/mesh/health` to check it is alive

### How to Test?

You can test the Mesh node as any other kind of node. The Mesh node support the --simulateRelay parameter as Servicer does, so you can use it.

Also, you can use [LocalNet Repository](https://github.com/pokt-scan/pocket-localnet) to deploy a local network and test all this together locally.

### TODO/Enhancements:
* Validate chains with servicer /v1/private/chains endpoint
* Support "Always Bill" approach
  * Allow mesh node to have a Remote CLI url (this should be a Pocket Node, close to mesh, not need to be stake) that will be used
  by the mesh node to run any kind of validations about relays before service them, grab session/response in cache to reduce next call response times, and finally resolve the request.
  * On this we may need a bit more of assistance from Pocket Core team to run the best possible validations for each relay using a cache from mesh side.
* Community feedback, issues, etc.

### Resources:

* Dockerhub Image:
* New external libraries:
  * Worker Pool: [pond](https://github.com/alitto/pond)
  * INotify: [inotify](https://github.com/fsnotify/fsnotify)
  * Http Retry: [httpretryable](https://github.com/hashicorp/go-retryablehttp)
  * Fastest Key/Value Cache: [pogreb](https://github.com/akrylysov/pogreb)
