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

* Relay Approach:
  * Validate servicer_health + minimum validations like pokt node does and then:
    * If not session in cache, query it to servicer_url/v1/private/mesh/session and then process relay
    * If session in cache and MaxRelays is not hit, process the relay.
* Proxy any request to servicer
* Monitor Pocket node. It will check the following:
  * Chains
  * Servicer addresses
  * Health (starting, catching up and height)
* Keep sessions in memory cache
* Keep relays in local cache (persistent) until they are notified
  * this ensures that even if the mesh node is restarted and the relay was not yet notified, they will be after node startup again
* Auto clean up old serve session from cache
* Servicer relays notification has a retry mechanism just in case the Servicer node is not responsive for a while.
  * this will only retry in some scenarios like http code greater than 401 and few code from pocket core
* Implements a worker queue (per pocket node) for the notification to keep the process simple and secure event under crash circumstances.
* Supports many nodes/servicer at once (like Lean Pokt or even multiple Lean Pokt at once)
* Handle minimum mesh node side validations using response get from Health monitor about Height, Starting and Catching Up (same of pocket node)
* Handle minimum relay validations about payload format (same of pocket node)
* Run connectivity chains on startup or reload:
  * Check: /v1/private/mesh/relay
  * Check: /v1/private/mesh/session
  * Check /v1/private/mesh/check
* Reload chains
* Reload keys
* Expose metrics of servicer relays
* Expose metrics of workers

### Included Branches:

* RC-0.9.2 [Lean Node](https://github.com/pokt-network/pocket-core/tree/staging)
* [Fix #1457 & Memory enhance](https://github.com/pokt-network/pocket-core/pull/1485)
* [Feature #1456](https://github.com/pokt-network/pocket-core/pull/1483)
  * This two PR was already merge for the next big RC of Pocket but due to the v0 was set on a sort of maintenance mode this change will not be release by them.

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
3. If your proxy has all the endpoints closed except `/v1` and `/v1/client`, please add: `/v1/private/mesh/<health|relay|session|check>`
4. Start your node as you were doing.

#### Setup Mesh Node:

1. Set up your proxy in the same way for a Servicer
2. Create the following `config.json` file inside the `--datadir` directory
```json
{
  "data_dir": "/home/app/.pocket/mesh",
  "rpc_port": "8081",
  "chains_name": "chains/chains.json",
  "client_rpc_timeout": 30000,
  "chains_rpc_timeout": 30000,
  "log_level": "*:info, *:error",
  "log_chain_request": false,
  "log_chain_response": false,
  "user_agent": "mesh-node",
  "auth_token_file": "key/auth.json",
  "json_sort_relay_responses": true,
  "relay_cache_file": "data/relays.pkt",
  "relay_cache_background_sync_interval": 3600,
  "relay_cache_background_compaction_interval": 18000,
  "keys_hot_reload_interval": 180000,
  "chains_hot_reload_interval": 180000,
  "worker_strategy": "balanced",
  "servicer_max_workers": 50,
  "servicer_max_workers_capacity": 50000,
  "servicer_workers_idle_timeout": 10000,
  "servicer_private_key_file": "key/key.json",
  "servicer_rpc_timeout": 60000,
  "servicer_auth_token_file": "key/auth.json",
  "servicer_retry_max_times": 10,
  "servicer_retry_wait_min": 10,
  "servicer_retry_wait_max": 180000,
  "node_check_interval": 60,
  "session_cache_clean_up_interval": 1800,
  "pocket_prometheus_port": "8083",
  "prometheus_max_open_files": 3,
  "metrics_worker_strategy": "lazy",
  "metrics_max_workers": 50,
  "metrics_max_workers_capacity": 50000,
  "metrics_workers_idle_timeout": 10000,
  "metrics_report_interval": 10
}
```
3. Create Servicer private key file with one of the  following formats into the path you set on `config.json`:
Fallback Format:
```json
[
  {
    // servicer private key
    "priv_key": "aaabbbbcccccddd",
    // servicer url/ip where mesh node can reach the servicer node to check health, proxy requests and notify relays
    // NOTE: name will be hostname parsed from servicer_url
    "servicer_url": "http://localhost:8081"
  },
  {
    ... // add as much servicers as you need to handle with one single geo-mesh process.
  }
]
```
New Format:
```json
[
  {
    // in case u do not set, will be the hostname of the url provided.
    "name": "optional name to identify it by name/uid",
    // servicer url/ip where mesh node can reach the servicer node to check health, proxy requests and notify relays
    "url": "http://localhost:8081",
    // servicers private keys
    "keys": ["<key1>", "<keyN>"]
  },
  {
    ...
    // add as much node/servicers as you need to handle with one single geo-mesh process.
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
7. Call your mesh node at `/v1/private/mesh/health?authtoken=<token>` to check it is alive and how many nodes/servicers it loaded from your setup.

### Config file details

| Key                                        	          | Type   	 | Default            	  | Description                                                                                                                                                               	                        |
|-------------------------------------------------------|----------|-----------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| data_dir                                   	          | String 	 | -                  	  | Where the data will be                                                                                                                                                    	                        |
| rpc_port                                   	          | String 	 | 8081               	  | Listening port                                                                                                                                                            	                        |
| chains_name                                	          | String 	 | chains.json        	  | Chains file path. This should be a filename path relative to --datadir                                                                                                    	                        |
| client_rpc_timeout                         	          | Number 	 | 30000              	  | Mesh Client RPC timeout                                                                                                                                                   	                        |
| log_level                                  	          | String 	 | *:info, *:error    	  | Logger namespace:level. Allow multiple values split by comma                                                                                                              	                        |
| log_chain_request                                  	  | Bool 	   | false    	            | When logger is set to debug, will attach chain request payload.                                                                                                              	                     |
| log_chain_response                                  	 | Bool 	   | false    	            | When logger is set to debug, will attach chain response payload.                                                                                                              	                    |
| user_agent                                 	          | String 	 | -                  	  | HTTP Header User-Agent value used on every sent request to Pocket Node.                                                                                                   	                        |
| auth_token_file                            	          | String 	 | auth/mesh.json     	  | Auth Token file for mesh private endpoints. This should be a filename path relative to --datadir                                                                          	                        |
| json_sort_relay_responses                  	          | Bool   	 | true               	  | Turn on JSON payload sorting for responses.                                                                                                                               	                        |
| chain_rpc_timeout                          	          | Number 	 | 30000              	  | Chains RPC timeout                                                                                                                                                        	                        |
| chain_rpc_max_idle_connections                        | Number 	 | 2500              	   | Chains RPC max idle connections                                                                                                                                                                    |
| chain_rpc_max_conns_per_host                          | Number 	 | 2500              	   | Chains RPC max connections per host                                                                                                                                                        	       |
| chain_rpc_max_idle_conns_per_host                     | Number 	 | 2500              	   | Chains RPC max idle connection per host                                                                                                                                                        	   |
| relay_cache_file                           	          | String 	 | data/relays.pkt    	  | Relays cache database. This database is used to persist relays in case mesh node is restarted, avoiding lost relays. This should be a filename path relative to --datadir 	                        |
| relay_cache_background_sync_interval       	          | Number 	 | 3600               	  | Time in milliseconds. Read More: https://pkg.go.dev/github.com/akrylysov/pogreb#Options                                                                                   	                        |
| relay_cache_background_compaction_interval 	          | Number 	 | 18000              	  | Time in milliseconds. Read More: https://pkg.go.dev/github.com/akrylysov/pogreb#Options                                                                                   	                        |
| keys_hot_reload_interval                   	          | Number 	 | 180000             	  | Interval in milliseconds to reload keys. Set 0 to disable.                                                                                                                	                        |
| chains_hot_reload_interval                 	          | Number 	 | 180000             	  | Interval in milliseconds to reload chains. Set 0 to disable.                                                                                                              	                        |
| servicer_private_key_file                  	          | String 	 | key/key.json       	  | Pocket Node / Servicer key file. This should be a filename path relative to --datadir                                                                                     	                        |
| servicer_rpc_timeout                       	          | Number 	 | 30000              	  | Pocket Node RPC calls timeout. Time in milliseconds.                                                                                                                      	                        |
| servicer_rpc_max_idle_connections                     | Number 	 | 2500              	   | Servicer RPC max idle connections                                                                                                                                                                  |
| servicer_rpc_max_conns_per_host                       | Number 	 | 2500              	   | Servicer RPC max connections per host                                                                                                                                                        	     |
| servicer_rpc_max_idle_conns_per_host                  | Number 	 | 2500              	   | Servicer RPC max idle connection per host                                                                                                                                                        	 |
| servicer_auth_token_file                   	          | String 	 | auth/servicer.json 	  | Auth Token file for call Pocket Node mesh endpoints. This should be a filename path relative to --datadir                                                                 	                        |
| servicer_retry_max_times                   	          | Number 	 | 10                 	  | How many time will a Pocket Node RPC call be retried.                                                                                                                     	                        |
| servicer_retry_wait_min                    	          | Number 	 | 5                  	  | How much is the min time to wait until retry a Pocket Node RPC call. Time in milliseconds.                                                                                	                        |
| servicer_retry_wait_max                    	          | Number 	 | 180                	  | How much is the max time to wait until retry a Pocket Node RPC call. Time in milliseconds.                                                                                	                        |
| servicer_worker_strategy                            	 | String 	 | balanced           	  | balanced \| eager \| lazy - Read more: https://github.com/alitto/pond#resizing-strategies                                                                                 	                        |
| servicer_max_workers                                	 | Number 	 | 50                 	  | Max amount of workers for each Pocket Node                                                                                                                                	                        |
| servicer_max_workers_capacity                       	 | Number 	 | 50000               	 | Max amount of tasks in queue without block it.                                                                                                                            	                        |
| servicer_workers_idle_timeout                       	 | Number 	 | 10000              	  | Worker idle timeout. Avoid values lowers than default one.                                                                                                                	                        |
| node_check_interval                        	          | Number 	 | 60                 	  | Pocket node check interval time. Time in seconds.                                                                                                                         	                        |
| session_cache_clean_up_interval            	          | Number 	 | 1800               	  | In memory cache clean up interval time. Time in seconds.                                                                                                                  	                        |
| pocket_prometheus_port                     	          | String 	 | 8083               	  | Prometheus metrics listening port.                                                                                                                                        	                        |
| prometheus_max_open_files                  	          | Number 	 | 3                  	  | Prometheus max open files.                                                                                                                                                	                        |
| metrics_worker_strategy                    	          | String 	 | lazy               	  | balanced \| eager \| lazy - Read more: <br>https://github.com/alitto/pond#resizing-strategies                                                                             	                        |
| metrics_max_workers                        	          | Number 	 | 50                 	  | Max amount of workers for each Metrics of each Pocket Node                                                                                                                	                        |
| metrics_max_workers_capacity               	          | Number 	 | 50000               	 | Max amount of tasks in queue without block it.                                                                                                                            	                        |
| metrics_workers_idle_timeout               	          | Number 	 | 10000              	  | Worker idle timeout. Avoid values lowers than default one.                                                                                                                	                        |
| metrics_report_interval                    	          | Number 	 | 10                 	  | Report interval for each Pocket node metric. Time in seconds.                                                                                                             	                        |

### How to Test?

You can test the Mesh node as any other kind of node. The Mesh node support the --simulateRelay parameter as Servicer does, so you can use it.

Also, you can use [LocalNet Repository](https://github.com/pokt-scan/pocket-localnet) to deploy a local network and test all this together locally.

### TODO/Enhancements:
* Validate chains with servicer /v1/private/chains endpoint
* Community feedback, issues, etc.

### Resources:

* Dockerhub Image:
* New external libraries:
  * Worker Pool: [pond](https://github.com/alitto/pond)
  * INotify: [inotify](https://github.com/fsnotify/fsnotify)
  * Http Retry: [httpretryable](https://github.com/hashicorp/go-retryablehttp)
  * Fastest Key/Value Cache: [pogreb](https://github.com/akrylysov/pogreb)
