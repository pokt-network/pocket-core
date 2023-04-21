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
This is just a basic sample, please refer to the full config details to extend and add/modify more config if u need.
```json
{
  "data_dir": "/home/app/.pocket/mesh",
  "rpc_port": "8081",
  "chains_name": "chains/chains.json",
  "client_rpc_timeout": 30000,
  "chains_rpc_timeout": 30000,
  "auth_token_file": "key/auth.json",
  "relay_cache_file": "data/relays.pkt",
  "worker_strategy": "balanced",
  "servicer_max_workers": 50,
  "servicer_max_workers_capacity": 50000,
  "servicer_workers_idle_timeout": 10000,
  "servicer_private_key_file": "key/key.json",
  "servicer_rpc_timeout": 60000,
  "servicer_auth_token_file": "key/auth.json",
  "node_check_interval": 60,
  "session_cache_clean_up_interval": 1800,
  "pocket_prometheus_port": "8083",
  "prometheus_max_open_files": 3,
  "metrics_moniker": "my-mesh-node-uid",
  "metrics_worker_strategy": "lazy",
  "metrics_max_workers": 50,
  "metrics_max_workers_capacity": 50000,
  "metrics_workers_idle_timeout": 10000,
  "metrics_report_interval": 10,
  "metrics_attach_servicer_label": false
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
6. Create a chains_name_map.json or point your `remote_chains_name_map` to https://poktscan-v1.nyc3.cdn.digitaloceanspaces.com/pokt-chains-map.json that is constantly updated.
   *. You can create a map like `{"0021": "Ethereum"}` or `{"0021": {"label": "Ethereum"}}` both will work.
   *. If u wish to use your own endpoint for remote, ensure it is a GET that return the JSON as expected.
7. Start your mesh node: `pocket start-mesh --datadir </your/path>`
8Call your mesh node at `/v1/private/mesh/health?authtoken=<token>` to check it is alive and how many nodes/servicers it loaded from your setup.

### Config file details

| Key                                        	          | Type   	 | Default            	        | Description                                                                                                                                                               	                                                                 |
|-------------------------------------------------------|----------|-----------------------------|---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| data_dir                                   	          | String 	 | -                  	        | Where the data will be                                                                                                                                                    	                                                                 |
| rpc_port                                   	          | String 	 | 8081               	        | Listening port                                                                                                                                                            	                                                                 |
| chains_name                                	          | String 	 | chains.json        	        | Chains file path. This should be a filename path relative to --datadir                                                                                                    	                                                                 |
| client_rpc_timeout                         	          | Number 	 | 30000              	        | Mesh Client RPC timeout                                                                                                                                                   	                                                                 |
| log_level                                  	          | String 	 | *:info, *:error    	        | Logger namespace:level. Allow multiple values split by comma                                                                                                              	                                                                 |
| log_chain_request                                  	  | Bool 	   | false    	                  | When logger is set to debug, will attach chain request payload.                                                                                                              	                                                              |
| log_chain_response                                  	 | Bool 	   | false    	                  | When logger is set to debug, will attach chain response payload.                                                                                                              	                                                             |
| user_agent                                 	          | String 	 | -                  	        | HTTP Header User-Agent value used on every sent request to Pocket Node.                                                                                                   	                                                                 |
| auth_token_file                            	          | String 	 | auth/mesh.json     	        | Auth Token file for mesh private endpoints. This should be a filename path relative to --datadir                                                                          	                                                                 |
| json_sort_relay_responses                  	          | Bool   	 | true               	        | Turn on JSON payload sorting for responses.                                                                                                                               	                                                                 |
| chains_name_map                  	                    | String   | -                           | (optional) Local file with a chain ID / NAME map used to enhance metrics.                                                                                                                            	                                      |
| remote_chains_name_map                  	             | String   	 | -               	           | (optional) Remote file with a chain ID / NAME map used to enhance metrics. This have precedense over local one.                                                                                                                           	 |
| chain_rpc_timeout                          	          | Number 	 | 30000              	        | Chains RPC timeout                                                                                                                                                        	                                                                 |
| chain_rpc_max_idle_connections                        | Number 	 | 2500              	         | Chains RPC max idle connections                                                                                                                                                                                                             |
| chain_rpc_max_conns_per_host                          | Number 	 | 2500              	         | Chains RPC max connections per host                                                                                                                                                        	                                                |
| chain_rpc_max_idle_conns_per_host                     | Number 	 | 2500              	         | Chains RPC max idle connection per host                                                                                                                                                        	                                            |
| relay_cache_file                           	          | String 	 | data/relays.pkt    	        | Relays cache database. This database is used to persist relays in case mesh node is restarted, avoiding lost relays. This should be a filename path relative to --datadir 	                                                                 |
| relay_cache_background_sync_interval       	          | Number 	 | 3600               	        | Time in milliseconds. Read More: https://pkg.go.dev/github.com/akrylysov/pogreb#Options                                                                                   	                                                                 |
| relay_cache_background_compaction_interval 	          | Number 	 | 18000              	        | Time in milliseconds. Read More: https://pkg.go.dev/github.com/akrylysov/pogreb#Options                                                                                   	                                                                 |
| keys_hot_reload_interval                   	          | Number 	 | 180000             	        | Interval in milliseconds to reload keys. Set 0 to disable.                                                                                                                	                                                                 |
| chains_hot_reload_interval                 	          | Number 	 | 180000             	        | Interval in milliseconds to reload chains. Set 0 to disable.                                                                                                              	                                                                 |
| servicer_private_key_file                  	          | String 	 | key/key.json       	        | Pocket Node / Servicer key file. This should be a filename path relative to --datadir                                                                                     	                                                                 |
| servicer_rpc_timeout                       	          | Number 	 | 30000              	        | Pocket Node RPC calls timeout. Time in milliseconds.                                                                                                                      	                                                                 |
| servicer_rpc_max_idle_connections                     | Number 	 | 2500              	         | Servicer RPC max idle connections                                                                                                                                                                                                           |
| servicer_rpc_max_conns_per_host                       | Number 	 | 2500              	         | Servicer RPC max connections per host                                                                                                                                                        	                                              |
| servicer_rpc_max_idle_conns_per_host                  | Number 	 | 2500              	         | Servicer RPC max idle connection per host                                                                                                                                                        	                                          |
| servicer_auth_token_file                   	          | String 	 | auth/servicer.json 	        | Auth Token file for call Pocket Node mesh endpoints. This should be a filename path relative to --datadir                                                                 	                                                                 |
| servicer_retry_max_times                   	          | Number 	 | 10                 	        | How many time will a Pocket Node RPC call be retried.                                                                                                                     	                                                                 |
| servicer_retry_wait_min                    	          | Number 	 | 5                  	        | How much is the min time to wait until retry a Pocket Node RPC call. Time in milliseconds.                                                                                	                                                                 |
| servicer_retry_wait_max                    	          | Number 	 | 180                	        | How much is the max time to wait until retry a Pocket Node RPC call. Time in milliseconds.                                                                                	                                                                 |
| servicer_worker_strategy                            	 | String 	 | balanced           	        | balanced \| eager \| lazy - Read more: https://github.com/alitto/pond#resizing-strategies                                                                                 	                                                                 |
| servicer_max_workers                                	 | Number 	 | 50                 	        | Max amount of workers for each Pocket Node                                                                                                                                	                                                                 |
| servicer_max_workers_capacity                       	 | Number 	 | 50000                       | Max amount of tasks in queue without block it.                                                                                                                            	                                                                 |
| servicer_workers_idle_timeout                       	 | Number 	 | 10000              	        | Worker idle timeout. Avoid values lowers than default one.                                                                                                                	                                                                 |
| node_check_interval                        	          | Number 	 | 60                 	        | Pocket node check interval time. Time in seconds.                                                                                                                         	                                                                 |
| session_cache_clean_up_interval            	          | Number 	 | 1800               	        | In memory cache clean up interval time. Time in seconds.                                                                                                                  	                                                                 |
| pocket_prometheus_port                     	          | String 	 | 8083               	        | Prometheus metrics listening port.                                                                                                                                        	                                                                 |
| prometheus_max_open_files                  	          | Number 	 | 3                  	        | Prometheus max open files.                                                                                                                                                	                                                                 |
| metrics_moniker                     	                 | String 	 | geo-mesh-node            	 | Metrics identifier, help full to identify a mesh instance from another. Also useful to collect multi region metrics on federate prometheus.                                                                          	                      |
| metrics_worker_strategy                    	          | String 	 | lazy               	        | balanced \| eager \| lazy - Read more: <br>https://github.com/alitto/pond#resizing-strategies                                                                             	                                                                 |
| metrics_max_workers                        	          | Number 	 | 50                 	        | Max amount of workers for each Metrics of each Pocket Node                                                                                                                	                                                                 |
| metrics_max_workers_capacity               	          | Number 	 | 50000                       | Max amount of tasks in queue without block it.                                                                                                                            	                                                                 |
| metrics_workers_idle_timeout               	          | Number 	 | 10000              	        | Worker idle timeout. Avoid values lowers than default one.                                                                                                                	                                                                 |
| metrics_report_interval                    	          | Number 	 | 10                 	        | Report interval for each Pocket node metric. Time in seconds.                                                                                                             	                                                                 |
| metrics_attach_servicer_label                    	    | Bool 	  | 10                 	        | Add servicer address to metric entries. This add more cardinality on your metrics.                                                                                                             	                                            |

### Metrics

* `moniker` helpful label to identify individual mesh instances across metrics
* `stat_type` could be `metric` or `node`
  * `metric` values of the worker pool to dispatch metrics to prometheus
  * `node` values of the worker used to notify relays to node
* `is_notify` indicate if the metric is result of the call from chain to mesh (aka relay) OR mesh to node (aka notify) 
* `status_type` indicate where in the code the error happen, expected values are: 
  * `internal` unexpected error like parsing/read values like json
  * `notify` error handled in the process of notify the node about a relay 
  * `chain` error handled on the call to the blockchain
* `servicer_address` label is optional, you can add or remove from metrics in case you want to reduce the cardinality of the metrics.
To do it set `metrics_attach_servicer_label` to `false` on config.json

| Name | Type | Labels | Description |
|---|---|---|---|
| pocketcore_geo_mesh_workers_running | gauge | moniker, stat_type, node_name | Number of running worker goroutines |
| pocketcore_geo_mesh_workers_idle | gauge | moniker, stat_type, node_name | Number of idle worker goroutines |
| pocketcore_geo_mesh_tasks_submitted_total | gauge | moniker, stat_type, node_name | Number of tasks submitted |
| pocketcore_geo_mesh_tasks_waiting_total | gauge | moniker, stat_type, node_name | Number of tasks waiting in the queue |
| pocketcore_geo_mesh_tasks_successful_total | gauge | moniker, stat_type, node_name | Number of tasks that completed successfully |
| pocketcore_geo_mesh_tasks_failed_total | gauge | moniker, stat_type, node_name | Number of tasks that completed with panic |
| pocketcore_geo_mesh_tasks_completed_total | gauge | moniker, stat_type, node_name | Number of tasks that completed either successfully or with panic |
| pocketcore_geo_mesh_min_workers | gauge | moniker, stat_type, node_name | Number min workers of node pool |
| pocketcore_geo_mesh_max_workers | gauge | moniker, stat_type, node_name | Number max workers of node pool |
| pocketcore_geo_mesh_max_capacity | gauge | moniker, stat_type, node_name | Number max capacity of node pool |
| pocketcore_geo_mesh_relay_count | counter | moniker, chain_id, chain_name, is_notify, servicer_address (optional) | Number of relays executed |
| pocketcore_geo_mesh_relay_time | histogram | moniker, chain_id, chain_name, is_notify, servicer_address (optional) | Relay duration in milliseconds |
| pocketcore_geo_mesh_error_count | counter | moniker, chain_id, chain_name, is_notify, status_type, status_code, servicer_address (optional) | Number of errors resulting from relays (mesh or chain) |

#### Metrics Grafana Dashboard

This dashboard allow you to monitor your mesh nodes across regions, so to properly work every metric need to contains an extra label called `region`
You can add it on your scrape job with the following rule:
```yml
scrape_configs:
  - job_name: <jobname>
    scrape_interval: <scrape interval>
    static_configs:
    - targets: [<target>]
      labels:
        regions: '<YOUR REGION>'
```

```json
{
  "annotations": {
    "list": [
      {
        "builtIn": 1,
        "datasource": {
          "type": "grafana",
          "uid": "-- Grafana --"
        },
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "target": {
          "limit": 100,
          "matchAny": false,
          "tags": [],
          "type": "dashboard"
        },
        "type": "dashboard"
      }
    ]
  },
  "description": "",
  "editable": true,
  "fiscalYearStartMonth": 0,
  "graphTooltip": 0,
  "id": 30,
  "links": [],
  "liveNow": false,
  "panels": [
    {
      "collapsed": false,
      "gridPos": {
        "h": 1,
        "w": 24,
        "x": 0,
        "y": 0
      },
      "id": 33,
      "panels": [],
      "title": "GeoMesh Statistics",
      "type": "row"
    },
    {
      "datasource": {
        "type": "prometheus"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "dark-red",
                "value": null
              },
              {
                "color": "green",
                "value": ""
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 12,
        "w": 12,
        "x": 0,
        "y": 1
      },
      "id": 49,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "prometheus"
          },
          "editorMode": "code",
          "expr": "sum by(region) (rate(pocketcore_geo_mesh_relay_count{region=~\"$region\", is_notify=\"false\", moniker=~\"$mesh_moniker\"}[$__rate_interval]))",
          "legendFormat": "__auto",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Relay Rate By Region",
      "type": "timeseries"
    },
    {
      "datasource": {
        "type": "prometheus"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 12,
        "w": 12,
        "x": 12,
        "y": 1
      },
      "id": 35,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "prometheus"
          },
          "editorMode": "code",
          "expr": "sum by(chain_name, region) (rate(pocketcore_geo_mesh_relay_count{chain_name=~\"$chain_name\", region=~\"$region\", is_notify=\"false\", moniker=~\"$mesh_moniker\"}[$__rate_interval]))",
          "legendFormat": "{{region}}:{{chain_name}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Relay Rate By Region/Chain",
      "type": "timeseries"
    },
    {
      "datasource": {
        "type": "prometheus"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            }
          },
          "mappings": []
        },
        "overrides": []
      },
      "gridPos": {
        "h": 9,
        "w": 5,
        "x": 0,
        "y": 13
      },
      "id": 56,
      "options": {
        "legend": {
          "displayMode": "list",
          "placement": "right",
          "showLegend": true
        },
        "pieType": "pie",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "prometheus"
          },
          "editorMode": "code",
          "expr": "sum by(region) (increase((pocketcore_geo_mesh_relay_count{chain_name=~\"$chain_name\", region=~\"$region\", is_notify=\"false\", moniker=~\"$mesh_moniker\"}[$__rate_interval])))",
          "legendFormat": "{{region}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Relay Breakdown By Region",
      "type": "piechart"
    },
    {
      "datasource": {
        "type": "prometheus"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            }
          },
          "mappings": []
        },
        "overrides": []
      },
      "gridPos": {
        "h": 9,
        "w": 7,
        "x": 5,
        "y": 13
      },
      "id": 57,
      "options": {
        "legend": {
          "displayMode": "list",
          "placement": "right",
          "showLegend": true
        },
        "pieType": "pie",
        "reduceOptions": {
          "calcs": [
            "lastNotNull"
          ],
          "fields": "",
          "values": false
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "prometheus"
          },
          "editorMode": "code",
          "expr": "sum by(region, chain_name) (increase((pocketcore_geo_mesh_relay_count{chain_name=~\"$chain_name\", region=~\"$region\", is_notify=\"false\", moniker=~\"$mesh_moniker\"}[$__rate_interval])))",
          "legendFormat": "{{region}}:{{chain_name}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Relay Breakdown by Region/Chain",
      "type": "piechart"
    },
    {
      "datasource": {
        "type": "prometheus"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 9,
        "w": 12,
        "x": 12,
        "y": 13
      },
      "id": 46,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "prometheus"
          },
          "editorMode": "code",
          "expr": "sum by (region) (rate(pocketcore_geo_mesh_tasks_completed_total{region=~\"$region\", stat_type=\"node\", moniker=~\"$mesh_moniker\"}[$__rate_interval]))",
          "legendFormat": "{{region}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Relay Notifier Rate By Region",
      "type": "timeseries"
    },
    {
      "datasource": {
        "type": "prometheus"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 15,
        "w": 12,
        "x": 0,
        "y": 22
      },
      "id": 54,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "prometheus"
          },
          "editorMode": "code",
          "expr": "sum by(region) (delta(pocketcore_geo_mesh_relay_count{region=~\"$region\", is_notify=\"false\", moniker=~\"$mesh_moniker\"}[$__rate_interval]))",
          "legendFormat": "__auto",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Relay 6h Delta By Region",
      "type": "timeseries"
    },
    {
      "datasource": {
        "type": "prometheus"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "off"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 15,
        "w": 12,
        "x": 12,
        "y": 22
      },
      "id": 55,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "prometheus"
          },
          "editorMode": "code",
          "expr": "sum by(chain_name, region) (delta(pocketcore_geo_mesh_relay_count{chain_name=~\"$chain_name\", region=~\"$region\", is_notify=\"false\", moniker=~\"$mesh_moniker\"}[6h]))",
          "legendFormat": "{{region}}:{{chain_name}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Relay 6h Delta By Region/Chain",
      "type": "timeseries"
    },
    {
      "datasource": {
        "type": "prometheus"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "area"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 100
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 16,
        "w": 12,
        "x": 0,
        "y": 37
      },
      "id": 40,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "prometheus"
          },
          "editorMode": "code",
          "expr": "sum by(region) (rate(pocketcore_geo_mesh_relay_time_sum{chain_name=~\"$chain_name\", region=~\"$region\", is_notify=\"false\", moniker=~\"$mesh_moniker\"}[$__rate_interval])) / sum by(region) (rate(pocketcore_geo_mesh_relay_time_count{chain_name=~\"$chain_name\", region=~\"$region\", is_notify=\"false\", moniker=~\"$mesh_moniker\"}[$__rate_interval]))",
          "legendFormat": "__auto",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Latency By Region",
      "type": "timeseries"
    },
    {
      "datasource": {
        "type": "prometheus"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "area"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 100
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 16,
        "w": 12,
        "x": 12,
        "y": 37
      },
      "id": 41,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "prometheus"
          },
          "editorMode": "code",
          "expr": "sum by(region, chain_name) (rate(pocketcore_geo_mesh_relay_time_sum{chain_name=~\"$chain_name\", region=~\"$region\", is_notify=\"false\", moniker=~\"$mesh_moniker\"}[$__rate_interval])) / sum by(region, chain_name) (rate(pocketcore_geo_mesh_relay_time_count{chain_name=~\"$chain_name\", region=~\"$region\", is_notify=\"false\", moniker=~\"$mesh_moniker\"}[$__rate_interval]))",
          "legendFormat": "{{region}}:{{chain_name}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Latency by Region/Chain",
      "type": "timeseries"
    },
    {
      "datasource": {
        "type": "prometheus"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "area"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 1
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 17,
        "w": 12,
        "x": 0,
        "y": 53
      },
      "id": 45,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "prometheus"
          },
          "editorMode": "code",
          "expr": "sum by(region, status_code, status_type) (rate(pocketcore_geo_mesh_error_count{chain_name=~\"$chain_name\", region=~\"$region\", is_notify=\"false\", moniker=~\"$mesh_moniker\"}[$__rate_interval]))",
          "legendFormat": "{{region}}:{{status_type}}:{{status_code}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Error Rate By Region",
      "type": "timeseries"
    },
    {
      "datasource": {
        "type": "prometheus"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "area"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 1
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 17,
        "w": 12,
        "x": 12,
        "y": 53
      },
      "id": 53,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "prometheus"
          },
          "editorMode": "code",
          "expr": "sum by(chain_name, region, status_code, status_type) (rate(pocketcore_geo_mesh_error_count{chain_name=~\"$chain_name\", region=~\"$region\", is_notify=\"false\", moniker=~\"$mesh_moniker\"}[$__rate_interval]))",
          "legendFormat": "{{region}}:{{chain_name}}:{{status_type}}:{{status_code}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Error Rate By Region/Chain",
      "transformations": [],
      "type": "timeseries"
    },
    {
      "datasource": {
        "type": "prometheus"
      },
      "description": "This is the servicer queue utilization per full node.",
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "palette-classic"
          },
          "custom": {
            "axisCenteredZero": false,
            "axisColorMode": "text",
            "axisLabel": "",
            "axisPlacement": "auto",
            "barAlignment": 0,
            "drawStyle": "line",
            "fillOpacity": 0,
            "gradientMode": "none",
            "hideFrom": {
              "legend": false,
              "tooltip": false,
              "viz": false
            },
            "lineInterpolation": "linear",
            "lineWidth": 1,
            "pointSize": 5,
            "scaleDistribution": {
              "type": "linear"
            },
            "showPoints": "auto",
            "spanNulls": false,
            "stacking": {
              "group": "A",
              "mode": "none"
            },
            "thresholdsStyle": {
              "mode": "area"
            }
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 50
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 13,
        "w": 12,
        "x": 0,
        "y": 70
      },
      "id": 50,
      "options": {
        "legend": {
          "calcs": [],
          "displayMode": "list",
          "placement": "bottom",
          "showLegend": true
        },
        "tooltip": {
          "mode": "single",
          "sort": "none"
        }
      },
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "prometheus"
          },
          "editorMode": "code",
          "expr": "(pocketcore_geo_mesh_tasks_waiting_total{moniker=~\"$mesh_moniker\"} / on(stat_type, node_name, moniker, region) pocketcore_geo_mesh_max_capacity{moniker=~\"$mesh_moniker\"}) * 100",
          "legendFormat": "{{region}}:{{moniker}}:{{stat_type}}",
          "range": true,
          "refId": "A"
        }
      ],
      "title": "Moniker Worker Utilization %",
      "transformations": [],
      "type": "timeseries"
    },
    {
      "datasource": {
        "type": "prometheus"
      },
      "fieldConfig": {
        "defaults": {
          "color": {
            "mode": "thresholds"
          },
          "custom": {
            "align": "left",
            "cellOptions": {
              "type": "auto"
            },
            "filterable": false,
            "inspect": false
          },
          "mappings": [],
          "thresholds": {
            "mode": "absolute",
            "steps": [
              {
                "color": "green",
                "value": null
              },
              {
                "color": "red",
                "value": 80
              }
            ]
          }
        },
        "overrides": []
      },
      "gridPos": {
        "h": 13,
        "w": 12,
        "x": 12,
        "y": 70
      },
      "id": 52,
      "options": {
        "footer": {
          "countRows": false,
          "fields": "",
          "reducer": [
            "sum"
          ],
          "show": false
        },
        "frameIndex": 0,
        "showHeader": true,
        "sortBy": [
          {
            "desc": false,
            "displayName": "Value"
          }
        ]
      },
      "pluginVersion": "9.4.7",
      "targets": [
        {
          "datasource": {
            "type": "prometheus",
            "uid": "prometheus"
          },
          "editorMode": "code",
          "expr": "pocketcore_geo_mesh_max_workers{moniker=~\"$mesh_moniker\"}",
          "legendFormat": "__auto",
          "range": true,
          "refId": "A"
        },
        {
          "datasource": {
            "type": "prometheus",
            "uid": "prometheus"
          },
          "editorMode": "code",
          "expr": "pocketcore_geo_mesh_min_workers{moniker=~\"$mesh_moniker\"}",
          "hide": false,
          "legendFormat": "__auto",
          "range": true,
          "refId": "B"
        },
        {
          "datasource": {
            "type": "prometheus",
            "uid": "prometheus"
          },
          "editorMode": "code",
          "expr": "pocketcore_geo_mesh_max_capacity{moniker=~\"$mesh_moniker\"}",
          "hide": false,
          "legendFormat": "__auto",
          "range": true,
          "refId": "C"
        },
        {
          "datasource": {
            "type": "prometheus",
            "uid": "prometheus"
          },
          "editorMode": "code",
          "expr": "pocketcore_geo_mesh_tasks_waiting_total{moniker=~\"$mesh_moniker\"}",
          "hide": false,
          "legendFormat": "__auto",
          "range": true,
          "refId": "D"
        },
        {
          "datasource": {
            "type": "prometheus",
            "uid": "prometheus"
          },
          "editorMode": "code",
          "expr": "pocketcore_geo_mesh_tasks_failed_total{moniker=~\"$mesh_moniker\"}",
          "hide": false,
          "legendFormat": "__auto",
          "range": true,
          "refId": "E"
        }
      ],
      "title": "Metric Worker Table",
      "transformations": [
        {
          "id": "labelsToFields",
          "options": {
            "keepLabels": [
              "node_name",
              "__name__",
              "stat_type",
              "moniker"
            ],
            "mode": "columns"
          }
        },
        {
          "id": "merge",
          "options": {}
        },
        {
          "id": "groupBy",
          "options": {
            "fields": {
              "Time": {
                "aggregations": []
              },
              "Value": {
                "aggregations": [
                  "lastNotNull"
                ],
                "operation": "aggregate"
              },
              "__name__": {
                "aggregations": [],
                "operation": "groupby"
              },
              "moniker": {
                "aggregations": [],
                "operation": "groupby"
              },
              "node_name": {
                "aggregations": [],
                "operation": "groupby"
              },
              "stat_type": {
                "aggregations": [],
                "operation": "groupby"
              }
            }
          }
        },
        {
          "id": "organize",
          "options": {
            "excludeByName": {},
            "indexByName": {},
            "renameByName": {
              "Value (lastNotNull)": "Value",
              "__name__": "Metric",
              "moniker": "Monkier",
              "node_name": "Full Node Host",
              "stat_type": "Stat type"
            }
          }
        },
        {
          "id": "sortBy",
          "options": {
            "fields": {},
            "sort": [
              {
                "field": "Monkier"
              }
            ]
          }
        },
        {
          "id": "sortBy",
          "options": {
            "fields": {},
            "sort": [
              {
                "field": "Value"
              }
            ]
          }
        },
        {
          "id": "sortBy",
          "options": {
            "fields": {},
            "sort": [
              {
                "field": "Stat type"
              }
            ]
          }
        }
      ],
      "type": "table"
    }
  ],
  "refresh": "30s",
  "revision": 1,
  "schemaVersion": 38,
  "style": "dark",
  "tags": [],
  "templating": {
    "list": [
      {
        "current": {
          "selected": true,
          "text": [
            "All"
          ],
          "value": [
            "$__all"
          ]
        },
        "datasource": {
          "type": "prometheus"
        },
        "definition": "label_values(region)",
        "hide": 0,
        "includeAll": true,
        "label": "Region",
        "multi": true,
        "name": "region",
        "options": [],
        "query": {
          "query": "label_values(region)",
          "refId": "StandardVariableQuery"
        },
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "sort": 1,
        "type": "query"
      },
      {
        "current": {
          "selected": true,
          "text": [
            "All"
          ],
          "value": [
            "$__all"
          ]
        },
        "datasource": {
          "type": "prometheus"
        },
        "definition": "label_values(chain_name)",
        "hide": 0,
        "includeAll": true,
        "label": "Blockchain",
        "multi": true,
        "name": "chain_name",
        "options": [],
        "query": {
          "query": "label_values(chain_name)",
          "refId": "StandardVariableQuery"
        },
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "sort": 3,
        "type": "query"
      },
      {
        "current": {
          "selected": true,
          "text": [
            "All"
          ],
          "value": [
            "$__all"
          ]
        },
        "datasource": {
          "type": "prometheus"
        },
        "definition": "label_values(moniker)",
        "hide": 0,
        "includeAll": true,
        "label": "Lean Node",
        "multi": true,
        "name": "mesh_moniker",
        "options": [],
        "query": {
          "query": "label_values(moniker)",
          "refId": "StandardVariableQuery"
        },
        "refresh": 1,
        "regex": "",
        "skipUrlSync": false,
        "sort": 0,
        "type": "query"
      }
    ]
  },
  "time": {
    "from": "now-12h",
    "to": "now"
  },
  "timepicker": {},
  "timezone": "utc",
  "title": "Pokt Monitoring",
  "uid": "A_kISSYVz",
  "version": 4,
  "weekStart": ""
}
```

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
