## Feature Description:

### Query Node signing Info

New cmd/rpc that enables easy querying for the signing information of the nodes. Some active and skillful node runners
have already played with this, parsing the information available from the state (using the v1/query/state rpc call). Now
on this release we are providing an easier way to query this useful information.

### How To Check node's signing-info

- Use the command ```pocket query signing-info <node-address> ```
- You should see an output similar to this one :
    ```
    {
    "page": 1,
    "result": [
        {
            "address": <node-address>,
            "index_offset": 2,
            "jailed_blocks_counter": 0,
            "jailed_until": "2020-12-19T00:49:39.004489669Z",
            "missed_blocks_counter": 2,
            "start_height": 0
        }
    ],
    "total_pages": 1
    }
    ```
- The query returns useful information for the node operator:

  ```jailed_blocks_counter``` : The amount of blocks the node has been in jail.

  ```jailed_until```: The time in jail until the node is able to submit an unjail tx.

  ```missed_block_counter``` : The amount of blocks this node has missed in the current block windows (currently set
  at ```10``` blocks)

	- To get the latest value for the param use the command ```pocket query param pos/SignedBlocksWindow```
	  you should see an output like this:
	    ```
		{
			"param_key": "pos/SignedBlocksWindow",
			"param_value": "10"
		}
		```

- The command also accepts an optional param for the height enabling the users to check node behaviour on past blocks
  ```pocket query signing-info <node-address> <height> ```


- This feature is also available as an RPC call :
    ```
    curl --location --request POST 'http://localhost:8081/v1/query/signinginfo' \
    --header 'Content-Type: application/json' \
    --data-raw '              {
                  "address": "<node-address>",
                  "height": <height>,
                  "page": 1,
                  "per_page": 1}'
    ```
  *TIP:* for RPC calls, address and height parameters can be omitted to get all staked nodes signing-infos

# What else I  can do with this query?

## How to check if my node is about to get jailed?

- Use the command ```pocket query signing-info <node-address> ```
- You should see an output similar to this one :
    ```
    {
    "page": 1,
    "result": [
        {
            "address": <node-address>,
            "index_offset": 2,
            "jailed_blocks_counter": 0,
            "jailed_until": "2020-12-19T00:49:39.004489669Z",
            "missed_blocks_counter": 2,
            "start_height": 0
        }
    ],
    "total_pages": 1
    }
    ```
- Look for the property **"missed_blocks_counter"**, if that value is larger than 0, means your node is missing blocks.
- For any given block window, currently defined at ```10``` blocks (```pos/SignedBlocksWindow```), your node is required
  to sign ```60%``` of the blocks on that window (```pos/MinSignedPerWindow```).
- That set of values currently means your node will need to sign at least ```6``` blocks out of ```10``` blocks in the
  window.
- If the **"missed_blocks_counter"** is at ```4``` you are about to get jailed on the next block.
- The **"missed_blocks_counter"** will reset every block window (```10``` blocks)

## How to check if my node is near the max jailed blocks

First, lets define ```max_jailed_blocks```, the current value for this parameter as defined on genesis is ```37,960```,
this is the maximun amount of block a node can be in jail.

If a node is left jailed this amount of blocks, it will become ```Force Unstaked```.

To be ```Force Unstaked``` means that **your node will receive a slash equivalent to the total amount that it is staked
for, effectively removing it from the network and burning staked the tokens**.

Force unstaking only happens when you get under the minimun stake amount (```pos/StakeMinimum```) by getting slashed for
bad behaviour or if your node is left on jail for the max jailed blocks (```pos/MaxJailedBlocks```).

Now lets see how to check if your node is at risk:

- Use the command ```pocket query signing-info <node-address> ```
- You should see an output similar to this one :
    ```
    {
    "page": 1,
    "result": [
        {
            "address": <node-address>,
            "index_offset": 2,
            "jailed_blocks_counter": 22476,
            "jailed_until": "2020-06-19T00:49:39.004489669Z",
            "missed_blocks_counter": 0,
            "start_height": 0
        }
    ],
    "total_pages": 1
    }
    ```
- Look for the property **"jailed_blocks_counter"** and compare the value against the ```max_jailed_blocks``` param (
  currently ```37,960``` )
- To get the latest value for the param use the command ```pocket query param pos/MaxJailedBlocks```
  you should see an output like this:
    ```
    {
        "param_key": "pos/MaxJailedBlocks",
        "param_value": "37960"
    }
    ```
- If your node  **"jailed_blocks_counter"** is close to that value, is ```highly recommended``` that you send and unjail
  tx immediately to solve this situation, or you may risk losing you staked tokens.

