Feature: Pocket Query Commands

  Scenario: To show existing commands within the pocket query section
	Given the user has a pocket client
	When the user runs the command "query"
	Then the user should be able to see standard output containing "Usage:"
	And the pocket client should have exited without error

  Scenario: To query an existing account
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query account {{address}}" against validator 0
	Then the user should be able to see standard output containing "/v1/query/account"
	And the user should be able to see standard output containing "address"
	And the pocket client should have exited without error

  Scenario: To query an existing account, wrong address
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query account {{address}} 0" against validator -1
	Then the user should be able to see standard output containing "encoding/hex: invalid byte:"
	And the pocket client should have exited without error

  Scenario: Pocket Query Account address is invalid But Has valid Bytes
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query account 2BBCA5DC9792C72AC3A2616910C4AAAA 0"
	Then the user should be able to see standard output containing "Incorrect address length"
	And the pocket client should have exited without error

  Scenario: To query an existing account, wrong height
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query account {{address}} 5600000" against validator 0
	Then the user should be able to see standard output containing "wanted to load target 5600000 but only found up to"
	And the pocket client should have exited without error

  Scenario: To query an non-existing Account
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query account c7b6a62385f8999c0cb63ad9cb464ee730c597b2 0"
	Then the user should be able to see standard output containing "null"
	And the pocket client should have exited without error

  Scenario: To query an existing account, incomplete command
	Given the user has a pocket client
	When the user runs the command "query account"
	Then the user should be able to see standard error containing "Usage:"
	And the pocket client should have exited with error

  Scenario: To query an existing account txs
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query account-txs c7b6a62385f8999c0cb63ad9cb464ee730c597b1 0 1 0 0" against validator 0
	Then the user should be able to see standard output containing "page_count"
	And the user should be able to see standard output containing "total_txs"
	And the pocket client should have exited without error

  Scenario: To query an existing account txs, wrong address
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query account-txs c7b6a62385f8999c0cb63ad9cb464ee730c597bw 0 1 0 0" against validator 0
	Then the user should be able to see standard output containing "encoding/hex: invalid byte:"
	And the pocket client should have exited without error

  Scenario: To query an existing account txs, incomplete command
	Given the user has a pocket client
	When the user runs the command "query account-txs"
	Then the user should be able to see standard error containing "Usage:"
	And the pocket client should have exited with error

  Scenario: To query current acl
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query acl"
	Then the user should be able to see standard output containing "/v1/query/acl"
	And the user should be able to see standard output containing "gov/non_map_acl"
	And the pocket client should have exited without error

  Scenario Outline: Unsuccessful Query of an App
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query app {{address}} 0" against validator <index>
	Then the user should be able to see standard output containing <message>
	And the pocket client should have exited without error
	Examples:
	  | index | message                      |
	  | -1    | "invalid byte"               |
	  | 0     | "application does not exist" |

  Scenario: To query a non-existing app from address, address has incorrect length but only valid bytes
	Given the user has a pocket client
	When the user runs the command "query app 4920ce1d787123456aeff366c79e8aa2 0"
	Then the user should be able to see standard output containing "Incorrect address length"
	And the pocket client should have exited without error

  Scenario: To query an existing app from address, incomplete command
	Given the user has a pocket client
	When the user runs the command "query app"
	Then the user should be able to see standard error containing "Usage:"
	And the pocket client should have exited with error

  Scenario: To query an existing app parameters
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query app-params"
	Then the user should be able to see standard output containing "/v1/query/appparams"
	And the user should be able to see standard output containing "base_relays_per_pokt"
	And the pocket client should have exited without error

  Scenario Outline: Querying App Parameters, Invalid Height
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command <command_with_bad_height>
	Then the user should be able to see standard output containing <error_msg>
	And the pocket client should have exited without error
	Examples:
	  | command_with_bad_height  | error_msg             |
	  | "query app-params 56000" | "error loading store" |
	  | "query app-params ^"     | "invalid syntax"      |

  Scenario: To query the list of existing apps
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query apps"
	Then the user should be able to see standard output containing "/v1/query/apps"
	And the user should be able to see standard output containing "result"
	And the pocket client should have exited without error

  Scenario: To query an address balance in the network
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query balance {{address}} 0" against validator 0
	Then the user should be able to see standard output containing "/v1/query/balance"
	And the pocket client should have exited without error

  Scenario: To query an address balance in the network, wrong address
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query balance {{address}} 0" against validator -1
	Then the user should be able to see standard output containing "invalid byte"
	And the pocket client should have exited without error

  Scenario: To query an address balance in the network, wrong height
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query balance {{address}} -3" against validator 0
	Then the user should be able to see standard error containing "unknown shorthand flag: '3' in -3"
	And the pocket client should have exited with error

  Scenario: To query an address balance in the network, incomplete command
	Given the user has a pocket client
	When the user runs the command "query balance"
	Then the user should be able to see standard error containing "Usage:"
	And the pocket client should have exited with error

  Scenario: To query a block at height, wrong height
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query block 0"
	Then the user should be able to see standard output containing "height must be greater than 0"
	And the pocket client should have exited without error

  Scenario: To query a block's at height, given an invalid number height
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query block 1000"
	Then the user should be able to see standard output containing "must be less than or equal to the current blockchain height"
	And the pocket client should have exited without error

  Scenario: To query a block's at height, given an invalid char in height
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query block 0-0"
	Then the user should be able to see standard output containing "invalid syntax"
	And the pocket client should have exited without error

  Scenario: To query block-txs at height
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query block-txs 0 1 1 1"
	Then the user should be able to see standard output containing "/v1/query/blocktxs"
	And the user should be able to see standard output containing "total_txs"
	And the pocket client should have exited without error

  Scenario: To query a block-txs at height, incomplete
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query block-txs"
	Then the user should be able to see standard error containing "Usage:"
	And the pocket client should have exited with error

  Scenario: To query dao owner
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query daoOwner" against validator 0
	Then the user should be able to see standard output containing "/v1/query/daoowner"

  Scenario: To query the chains height
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query height"
	Then the user should be able to see standard output containing "height"
	And the pocket client should have exited without error

  Scenario: To query an existing node from address
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query node {{address}}" against validator 0
	Then the user should be able to see standard output containing "/v1/query/node"
	And the user should be able to see standard output containing "address"
	And the user should be able to see standard output containing "chains"
	And the user should be able to see standard output containing "jailed"
	And the pocket client should have exited without error

  Scenario: To query an existing node from address, wrong address
	Given the user has a pocket client
	When the user runs the command "query node tsst"
	Then the user should be able to see standard output containing "invalid byte"
	And the pocket client should have exited without error

  Scenario: To query an existing node from address, wrong address with only valid bytes
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query node d8cbb00bf6ea51971448eae2fe8d8321ffadbf4f"
	Then the user should be able to see standard output containing "validator not found"
	And the pocket client should have exited without error

  Scenario: To query an existing node from address, incomplete command
	Given the user has a pocket client
	When the user runs the command "query node"
	Then the user should be able to see standard error containing "Usage:"
	And the pocket client should have exited with error

  Scenario: To query an existing node claim, incomplete
	Given the user has a pocket client
	When the user runs the command "query node-claim"
	Then the user should be able to see standard error containing "Usage:"
	And the pocket client should have exited with error

  Scenario: To query an existing node claims
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query node-claims {{address}} 0" against validator 0
	Then the user should be able to see standard output containing "/v1/query/nodeclaims"
	And the user should be able to see standard output containing "page"
	And the user should be able to see standard output containing "result"
	And the user should be able to see standard output containing "total_pages"
	And the pocket client should have exited without error

  Scenario: To query an existing node parameters
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query node-claims" against validator 0
	Then the user should be able to see standard output containing "/v1/query/nodeclaims"
	And the user should be able to see standard output containing "result"
	And the pocket client should have exited without error

  Scenario: To query an existing node parameters, wrong address
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query node-claims {{address}}" against validator -1
	Then the user should be able to see standard output containing "invalid byte"
	And the pocket client should have exited without error

  Scenario: To query the list of existing nodes
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query nodes" against validator 0
	Then the user should be able to see standard output containing "/v1/query/nodes"
	And the user should be able to see standard output containing "result"
	And the pocket client should have exited without error

  Scenario: To query the list of existing nodes, Returns error when given wrong height
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query nodes 60000"
	Then the user should be able to see standard output containing "error loading store"
	And the pocket client should have exited without error

  Scenario: To query the list of existing nodes
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query nodes" against validator 0
	Then the user should be able to see standard output containing "result"
	And the user should be able to see standard output containing "address"
	And the user should be able to see standard output containing "jailed"
	And the user should be able to see standard output containing "chains"
	And the user should be able to see standard output containing "service_url"
	And the pocket client should have exited without error

  Scenario: To query single param
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query param pocketcore/MinimumNumberOfProofs 0" against validator 0
	Then the user should be able to see standard output containing "param_key"
	And the user should be able to see standard output containing "pocketcore/MinimumNumberOfProofs"
	And the pocket client should have exited without error

  Scenario: To query the list of params
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query params"
	Then the user should be able to see standard output containing "param_key"
	And the user should be able to see standard output containing "param_value"
	And the pocket client should have exited without error

  Scenario: To query the list of pocketcore params
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query pocket-params"
	Then the user should be able to see standard output containing "claim_expiration"
	And the pocket client should have exited without error

  Scenario: To query the list of pocketcore params, Return error with incorrect height
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query pocket-params 9999999"
	Then the user should be able to see standard output containing "error loading store"
	And the pocket client should have exited without error

  Scenario: To query a specific node's signing info
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query signing-info {{address}} 0" against validator 0
	Then the user should be able to see standard output containing "/v1/query/signinginfo"
	And the user should be able to see standard output containing "result"
	And the pocket client should have exited without error

  Scenario: To get supplies
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query supply"
	Then the user should be able to see standard output containing "/v1/query/supply"
	And the user should be able to see standard output containing "app_staked"
	And the pocket client should have exited without error

  Scenario: To get supplies, Returns error code with wrong height
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query supply 999999"
	Then the user should be able to see standard output containing "error loading store"
	And the pocket client should have exited without error

  Scenario: To get existing supported networks
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query supported-networks"
	Then the user should be able to see standard output containing "/v1/query/supportedchains"
	And the user should be able to see standard output containing "0011"
	And the pocket client should have exited without error

  Scenario: To get existing supported networks, Returns an error code with invalid height
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query supported-networks 999999"
	Then the user should be able to see standard output containing "error loading store"
	And the pocket client should have exited without error

  Scenario Outline: To get transactions based on the hash, wrong hash
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command <query_command>
	Then the user should be able to see standard output containing <error_msg>
	And the pocket client should have exited without error
	Examples:
	  | query_command                                                                  | error_msg               |
	  | "query tx 01a"                                                                 | "odd length hex string" |
	  | "query tx 23197E4D46009879F28F978A90627C7DFEAB64B4777AFCC24E2B9C3D72B4DADA22"  | "not found"             |
	  | "query tx 23197E4D4)(&09879F28F978A90627C7DFEAB64B4777AFCC24E2B9C3D72B4DADA22" | "invalid byte"          |

  Scenario: To get transactions based on the hash, incomplete commands
	Given the user has a pocket client
	When the user runs the command "query tx"
	Then the user should be able to see standard error containing "Usage:"
	And the pocket client should have exited with error

  Scenario: To query the latest upgrade
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query upgrade"
	Then the user should be able to see standard output containing "/v1/query/upgrade"
	And the user should be able to see standard output containing "Height"
	And the pocket client should have exited without error

  Scenario: To query the latest upgrade, invalid height
	Given the user is running the network "single_node_network"
	And the user has a pocket client
	When the user runs the command "query upgrade 50000"
	Then the user should be able to see standard output containing "error loading store"
	And the pocket client should have exited without error
