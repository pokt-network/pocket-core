package cli

import (
	"encoding/json"
	"fmt"
	types2 "github.com/pokt-network/pocket-core/x/apps/types"
	"github.com/pokt-network/posmint/types"
	"strconv"
	"strings"

	"github.com/pokt-network/pocket-core/app"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/spf13/cobra"
	core_types "github.com/tendermint/tendermint/rpc/core/types"
)

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.AddCommand(queryBlock)
	queryCmd.AddCommand(queryHeight)
	queryCmd.AddCommand(queryTx)
	queryCmd.AddCommand(queryAccountTxs)
	queryCmd.AddCommand(queryBlockTxs)
	queryCmd.AddCommand(queryNodes)
	queryCmd.AddCommand(queryBalance)
	queryCmd.AddCommand(queryAccount)
	queryCmd.AddCommand(queryNode)
	queryCmd.AddCommand(queryApps)
	queryCmd.AddCommand(queryApp)
	queryCmd.AddCommand(queryNodeParams)
	queryCmd.AddCommand(queryAppParams)
	queryCmd.AddCommand(queryNodeReceipts)
	queryCmd.AddCommand(queryNodeReceipt)
	queryCmd.AddCommand(queryPocketParams)
	queryCmd.AddCommand(queryPocketSupportedChains)
	queryCmd.AddCommand(querySupply)
	queryCmd.AddCommand(queryUpgrade)
	queryCmd.AddCommand(queryACL)
	queryCmd.AddCommand(queryDAOOwner)
}

var queryCmd = &cobra.Command{
	Use:   "query",
	Short: "query the blockchain",
	Long: `The query namespace handles all queryable interactions,
From getting Blocks, transactions, height; to getting params`,
}

var queryBlock = &cobra.Command{
	Use:   "block <height>",
	Short: "Get block at height",
	Long:  `Retrieves the block structure at the specified height.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var height *int64
		if len(args) == 0 {
			height = nil
		} else {
			var err error
			parsed, err := strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
			convert := int64(parsed)
			height = &convert
		}
		res, err := app.QueryBlock(height)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(string(res))
	},
}

var queryTx = &cobra.Command{
	Use:   "tx <hash>",
	Short: "Get the transaction by the hash",
	Long:  `Retrieves the transaction by the hash`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		res, err := app.QueryTx(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)
	},
}

func validatePagePerPageProveReceivedArgs(args []string) (page int, perPage int, prove bool, received bool) {
	page = 0
	perPage = 0
	prove = false
	received = false
	if len(args) >= 2 {
		parsedPage, err := strconv.Atoi(args[1])
		if err == nil {
			page = parsedPage
		}
	}
	if len(args) >= 3 {
		parsedPerPage, err := strconv.Atoi(args[2])
		if err == nil {
			perPage = parsedPerPage
		}
	}
	if len(args) >= 4 {
		parsedProve, err := strconv.ParseBool(args[3])
		if err == nil {
			prove = parsedProve
		}
	}
	if len(args) == 5 {
		parsedReceived, err := strconv.ParseBool(args[4])
		if err == nil {
			received = parsedReceived
		}
	}
	return page, perPage, prove, received
}

var queryAccountTxs = &cobra.Command{
	Use:   "account-txs <address> <page> <per_page> <prove> <received>",
	Short: "Get the transactions sent by the address, paginated by page and per_page",
	Long:  `Retrieves the transactions sent by the address`,
	Args:  cobra.RangeArgs(1, 5),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		page, perPage, prove, received := validatePagePerPageProveReceivedArgs(args)
		var res *core_types.ResultTxSearch
		var err error
		if received == true {
			res, err = app.QueryRecipientTxs(args[0], page, perPage, prove)
		} else {
			res, err = app.QueryAccountTxs(args[0], page, perPage, prove)
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		jsonRes, _ := json.Marshal(res)
		fmt.Println(fmt.Sprintf("%s", jsonRes))
	},
}

var queryBlockTxs = &cobra.Command{
	Use:   "block-txs <height> <page> <per_page> <prove>",
	Short: "Get the transactions at a certain block height, paginated by page and per_page",
	Long:  `Retrieves the transactions in the block height`,
	Args:  cobra.RangeArgs(1, 4),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		page, perPage, prove, _ := validatePagePerPageProveReceivedArgs(args)
		height, parsingErr := strconv.ParseInt(args[0], 10, 64)
		if parsingErr != nil {
			fmt.Println(parsingErr)
			return
		}
		res, err := app.QueryBlockTxs(height, page, perPage, prove)
		if err != nil {
			fmt.Println(err)
			return
		}
		jsonRes, _ := json.Marshal(res)
		fmt.Println(fmt.Sprintf("%s", jsonRes))
	},
}

var queryHeight = &cobra.Command{
	Use:   "height",
	Short: "Get current height",
	Long:  `Retrieves the current height`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		res, err := app.QueryHeight()
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Block Height: %d\n", res)
	},
}

var queryBalance = &cobra.Command{
	Use:   "balance <accAddr> <height>",
	Short: "Gets account balance",
	Long:  `Retrieves the balance of the specified <accAddr> at the specified <height>.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var height int
		if len(args) == 1 {
			height = 0 // latest
		} else {
			var err error
			height, err = strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		res, err := app.QueryBalance(args[0], int64(height))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Account Balance: %v\n", res)
	},
}

var queryAccount = &cobra.Command{
	Use:   "account <accAddr> <height>",
	Short: "Gets an account",
	Long:  `Retrieves the account structure for a specific address.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var height int
		if len(args) == 1 {
			height = 0 // latest
		} else {
			var err error
			height, err = strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		res, err := app.QueryAccount(args[0], int64(height))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%v\n", res)
	},
}

var nodeStakingStatus string
var nodeJailedStatus string
var blockchain string
var nodePage int
var nodeLimit int

func init() {
	queryNodes.Flags().StringVar(&nodeStakingStatus, "staking-status", "", "the staking status of the node")
	queryNodes.Flags().StringVar(&nodeJailedStatus, "jailed-status", "", "the jailed status of the node")
	queryNodes.Flags().StringVar(&blockchain, "blockchain", "", "the network identifier these nodes support")
	queryNodes.Flags().IntVar(&nodePage, "nodePage", 1, "mark the nodePage you want")
	queryNodes.Flags().IntVar(&nodeLimit, "nodeLimit", 10000, "reduce the amount of results")
}

var queryNodes = &cobra.Command{
	Use:   "nodes --staking-status <staked or unstaking> --jailed-status <jailed or unjailed> --blockchain <network id> --nodePage=<nodePage> --nodeLimit=<nodeLimit> <height>",
	Short: "Gets nodes",
	Long:  `Retrieves the list of all nodes known at the specified <height>.`,
	// Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var height int
		if len(args) == 0 {
			height = 0 // latest
		} else {
			var err error
			height, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		var err error
		opts := nodeTypes.QueryValidatorsParams{
			Blockchain: blockchain,
			Page:       nodePage,
			Limit:      nodeLimit,
		}
		if nodeStakingStatus != "" {
			switch strings.ToLower(nodeStakingStatus) {
			case "staked":
				opts.StakingStatus = types.Staked
			case "unstaking":
				opts.StakingStatus = types.Unstaking
			default:
				fmt.Println(fmt.Errorf("unkown staking status <staked or unstaking>"))
			}
		}
		if nodeJailedStatus != "" {
			switch strings.ToLower(nodeStakingStatus) {
			case "jailed":
				opts.JailedStatus = 1
			case "unjailed":
				opts.JailedStatus = 2
			default:
				fmt.Println(fmt.Errorf("unkown jailed status <jailed or unjailed>"))
			}
		}
		res, err := app.QueryNodes(int64(height), opts)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Nodes\n%s\n", res.String())
	},
}

var queryNode = &cobra.Command{
	Use:   "node <address> <height>",
	Short: "Gets node from address",
	Long:  `Retrieves the node at the specified <height>.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var height int
		if len(args) == 1 {
			height = 0 // latest
		} else {
			var err error
			height, err = strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		res, err := app.QueryNode(args[0], int64(height))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res.String())
	},
}

var queryNodeParams = &cobra.Command{
	Use:   "node-params <height>",
	Short: "Gets node parameters",
	Long:  `Retrieves the node parameters at the specified <height>.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var height int
		if len(args) == 0 {
			height = 0 // latest
		} else {
			var err error
			height, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		res, err := app.QueryNodeParams(int64(height))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res.String())
	},
}

var appStakingStatus string
var appPage, appLimit int

func init() {
	queryApps.Flags().StringVar(&nodeStakingStatus, "staking-status", "", "the staking status of the node")
	queryApps.Flags().IntVar(&nodePage, "appPage", 1, "mark the page you want")
	queryApps.Flags().IntVar(&nodeLimit, "appLimit", 10000, "reduce the amount of results")
}

var queryApps = &cobra.Command{
	Use:   "apps --staking-status=<nodeStakingStatus> --nodePage=<nodePage> --nodeLimit=<nodeLimit> <height>",
	Short: "Gets apps",
	Long:  `Retrieves the list of all applications known at the specified <height>`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var height int
		if len(args) == 0 {
			height = 0 // latest
		} else {
			var err error
			height, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		opts := types2.QueryApplicationsWithOpts{
			Blockchain: blockchain,
			Page:       appPage,
			Limit:      appLimit,
		}
		if appStakingStatus != "" {
			switch strings.ToLower(nodeStakingStatus) {
			case "staked":
				opts.StakingStatus = types.Staked
			case "unstaking":
				opts.StakingStatus = types.Unstaking
			default:
				fmt.Println(fmt.Errorf("unkown staking status <staked or unstaking>"))
			}
		}
		res, err := app.QueryApps(int64(height), opts)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Apps:\n%s\n", res.String())
	},
}

var queryApp = &cobra.Command{
	Use:   "app <address> <height>",
	Short: "Gets app from address",
	Long:  `Retrieves the app at the specified <height>.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var height int
		if len(args) == 1 {
			height = 0 // latest
		} else {
			var err error
			height, err = strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		res, err := app.QueryApp(args[0], int64(height))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res.String())
	},
}

var queryAppParams = &cobra.Command{
	Use:   "app-params <height>",
	Short: "Gets app parameters",
	Long:  `Retrieves the app parameters at the specified <height>.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var height int
		if len(args) == 0 {
			height = 0 // latest
		} else {
			var err error
			height, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		res, err := app.QueryAppParams(int64(height))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res.String())
	},
}

var queryNodeReceipts = &cobra.Command{
	Use:   "node-receipts <nodeAddr> <height>",
	Short: "Gets node receipts for work completed",
	Long:  `Retrieves the list of all verified proof of work submitted by <nodeAddr> at <height>.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var height int
		if len(args) == 1 {
			height = 0 // latest
		} else {
			var err error
			height, err = strconv.Atoi(args[1])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		res, err := app.QueryReceipts(args[0], int64(height))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("MerkleProofs:")
		for _, p := range res {
			fmt.Printf("%v\n", p)
		}
	},
}

var queryNodeReceipt = &cobra.Command{
	Use:   "node-receipt <nodeAddr> <appPubKey> <receiptType> <networkId> <sessionHeight> <height>`",
	Short: "Gets node receipt for work completed",
	Long:  `Gets node receipt for verified proof of work submitted for a specific session`,
	Args:  cobra.MinimumNArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var height int
		if len(args) == 5 {
			height = 0 // latest
		} else {
			var err error
			height, err = strconv.Atoi(args[4])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		sessionheight, err := strconv.Atoi(args[4])
		if err != nil {
			fmt.Println(err)
			return
		}
		res, err := app.QueryReceipt(args[3], args[1], args[0], args[2], int64(sessionheight), int64(height))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%v\n", res)
	},
}

var queryPocketParams = &cobra.Command{
	Use:   "pocket-params <height>",
	Short: "Gets pocket parameters",
	Long:  `Retrieves the pocket parameters at the specified <height>.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var height int
		if len(args) == 0 {
			height = 0 // latest
		} else {
			var err error
			height, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		res, err := app.QueryPocketParams(int64(height))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res.String())
	},
}

var queryPocketSupportedChains = &cobra.Command{
	Use:   "supported-networks <height>",
	Short: "Gets pocket supported networks",
	Long:  `Retrieves the list Network Identifiers supported by the network at the specified <height>`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var height int
		if len(args) == 0 {
			height = 0 // latest
		} else {
			var err error
			height, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		res, err := app.QueryPocketSupportedBlockchains(int64(height))
		if err != nil {
			fmt.Println(err)
			return
		}
		for i, chain := range res {
			fmt.Printf("(%d)\t%s\n", i, chain)
		}
	},
}

var querySupply = &cobra.Command{
	Use:   "supply <height>",
	Short: "Gets the supply at <height>",
	Long:  `Retrieves the list of node params specified in the <height>`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var height int
		if len(args) == 0 {
			height = 0 // latest
		} else {
			var err error
			height, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		nodesStake, total, err := app.QueryTotalNodeCoins(int64(height))
		if err != nil {
			fmt.Println(err)
			return
		}
		appsStaked, err := app.QueryTotalAppCoins(int64(height))
		if err != nil {
			fmt.Println(err)
			return
		}
		dao, err := app.QueryDaoBalance(int64(height))
		if err != nil {
			fmt.Println(err)
			return
		}
		totalStaked := nodesStake.Add(appsStaked).Add(dao)
		totalUnstaked := total.Sub(totalStaked)
		fmt.Printf("Nodes Staked:\t%v\nApps Staked:\t%v\n"+
			"Dao Supply:\t%v\nTotal Staked:\t%v\nTotalUnstaked:\t%v\nTotal Supply:\t%v\n\n",
			nodesStake, appsStaked, dao, totalStaked, totalUnstaked, total,
		)
	},
}

var queryDAOOwner = &cobra.Command{
	Use:   "daoOwner <height>",
	Short: "Gets the owner of the dao",
	Long:  `Retrieves the owner of the DAO (the account that can send/burn coins from the dao)`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var height int
		if len(args) == 0 {
			height = 0 // latest
		} else {
			var err error
			height, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		res, err := app.QueryDaoOwner(int64(height))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%v\n", res)
	},
}

var queryACL = &cobra.Command{
	Use:   "acl <height>",
	Short: "Gets the gov acl",
	Long:  `Retrieves the access control list of governance params (which account can change the param)`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var height int
		if len(args) == 0 {
			height = 0 // latest
		} else {
			var err error
			height, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		res, err := app.QueryACL(int64(height))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%v\n", res)
	},
}

var queryUpgrade = &cobra.Command{
	Use:   "upgrade <height>",
	Short: "Gets the latest gov upgrade",
	Long:  `Retrieves the latest protocol upgrade by governance`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var height int
		if len(args) == 0 {
			height = 0 // latest
		} else {
			var err error
			height, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		res, err := app.QueryUpgrade(int64(height))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("%v\n", res)
	},
}
