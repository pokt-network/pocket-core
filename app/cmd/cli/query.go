package cli

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pokt-network/pocket-core/app"
	appTypes "github.com/pokt-network/pocket-core/x/apps/types"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/spf13/cobra"
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
	Long:  ``,
}

var queryBlock = &cobra.Command{
	Use:   "block <height>",
	Short: "Get block at height",
	Long:  `Returns the block structure at the specified height.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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
	Args:  cobra.ExactArgs(1),
	Long:  `Returns the transaction by the hash`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		res, err := app.QueryTx(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)
	},
}

var queryAccountTxs = &cobra.Command{
	Use:   "account-txs <address> <page> <per_page>",
	Short: "Get the transactions sent by the address, paginated by page and per_page",
	Args:  cobra.RangeArgs(1, 3),
	Long:  `Returns the transactions sent by the address`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		page := 1
		perPage := 30
		if len(args) == 2 {
			parsedPage, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				page = int(parsedPage)
			}
		}
		if len(args) == 3 {
			parsedPerPage, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				perPage = int(parsedPerPage)
			}
		}
		res, err := app.QueryAccountTxs(args[0], page, perPage)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)
	},
}

var queryBlockTxs = &cobra.Command{
	Use:   "block-txs <height> <page> <per_page>",
	Short: "Get the transactions at a certain block height, paginated by page and per_page",
	Args:  cobra.RangeArgs(1, 3),
	Long:  `Returns the transactions in the block height`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		page := 1
		perPage := 30
		if len(args) == 2 {
			parsedPage, err := strconv.ParseInt(args[1], 10, 64)
			if err != nil {
				page = int(parsedPage)
			}
		}
		if len(args) == 3 {
			parsedPerPage, err := strconv.ParseInt(args[2], 10, 64)
			if err != nil {
				perPage = int(parsedPerPage)
			}
		}
		height, parsingErr := strconv.ParseInt(args[0], 10, 64)
		if parsingErr != nil {
			fmt.Println(parsingErr)
			return
		}
		res, err := app.QueryBlockTxs(height, page, perPage)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)
	},
}

var queryHeight = &cobra.Command{
	Use:   "height",
	Short: "Get current height",
	Long:  `Returns the current height`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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
	Args:  cobra.MinimumNArgs(1),
	Long:  `Returns the balance of the specified <accAddr> at the specified <height>.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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
	Args:  cobra.MinimumNArgs(1),
	Long:  `Returns the account structure for a specific address.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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

func init() {
	queryNodes.Flags().StringVar(&nodeStakingStatus, "staking-status", "", "the staking status of the node")
}

var queryNodes = &cobra.Command{
	Use:   "nodes --staking-status=<nodeStakingStatus> <height>",
	Short: "Gets nodes",
	Long:  `Returns the list of all nodes known at the specified <height>.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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
		var res nodeTypes.Validators
		var err error
		switch strings.ToLower(nodeStakingStatus) {
		case "":
			// no status passed
			res, err = app.QueryAllNodes(int64(height))
		case "staked":
			// staked nodes
			res, err = app.QueryStakedNodes(int64(height))
		case "unstaked":
			// unstaked nodes
			res, err = app.QueryUnstakedNodes(int64(height))
		case "unstaking":
			// unstaking nodes
			res, err = app.QueryUnstakingNodes(int64(height))
		default:
			fmt.Println("invalid staking status, can be staked, unstaked, unstaking, or empty")
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		if res == nil {
			fmt.Println("nil nodes result")
			return
		}
		fmt.Printf("Nodes\n%s\n", res.String())
	},
}

var queryNode = &cobra.Command{
	Use:   "node <address> <height>",
	Short: "Gets node from address",
	Args:  cobra.MinimumNArgs(1),
	Long:  `Returns the node at the specified <height>.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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
	Long:  `Returns the node parameters at the specified <height>.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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

func init() {
	queryApps.Flags().StringVar(&nodeStakingStatus, "staking-status", "", "the staking status of the node")
}

var queryApps = &cobra.Command{
	Use:   "apps --staking-status=<nodeStakingStatus> <height>",
	Short: "Gets apps",
	Long:  `Returns the list of all applications known at the specified <height>`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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
		var res appTypes.Applications
		var err error
		switch strings.ToLower(appStakingStatus) {
		case "":
			// no status passed
			res, err = app.QueryAllApps(int64(height))
		case "staked":
			// staked nodes
			res, err = app.QueryStakedApps(int64(height))
		case "unstaked":
			// unstaked nodes
			res, err = app.QueryUnstakedApps(int64(height))
		case "unstaking":
			// unstaking nodes
			res, err = app.QueryUnstakingApps(int64(height))
		default:
			fmt.Printf("invalid staking status, can be staked, unstaked, unstaking, or empty")
			return
		}
		if err != nil {
			fmt.Println(err)
			return
		}
		if res == nil {
			fmt.Println("nil Apps result")
			return
		}
		fmt.Printf("Apps:\n%s\n", res.String())
	},
}

var queryApp = &cobra.Command{
	Use:   "app <address> <height>",
	Short: "Gets app from address",
	Args:  cobra.MinimumNArgs(1),
	Long:  `Returns the app at the specified <height>.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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
	Long:  `Returns the app parameters at the specified <height>.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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
	Args:  cobra.MinimumNArgs(1),
	Long:  `Returns the list of all verified proof of work submitted by <nodeAddr> at <height>.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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
	Long:  `Returns the pocket parameters at the specified <height>.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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
	Long:  `Returns the list Network Identifiers supported by the network at the specified <height>`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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
	Short: "Returns",
	Long:  `Returns the list of node params specified in the <height>`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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
		nodesStake, nodesUnstaked, err := app.QueryTotalNodeCoins(int64(height))
		if err != nil {
			fmt.Println(err)
			return
		}
		appsStaked, appsUnstaked, err := app.QueryTotalAppCoins(int64(height))
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
		totalUnstaked := nodesUnstaked.Add(appsUnstaked).Sub(nodesStake).Sub(appsStaked)
		total := totalStaked.Add(totalUnstaked)
		fmt.Printf("Nodes Staked:\t%v\nApps Staked:\t%v\n"+
			"Dao Supply:\t%v\nTotal Staked:\t%v\nTotalUnstaked:\t%v\nTotal Supply:\t%v\n\n",
			nodesStake, appsStaked, dao, totalStaked, totalUnstaked, total,
		)
	},
}

var queryDAOOwner = &cobra.Command{
	Use:   "daoOwner <height>",
	Short: "Gets the owner of the dao",
	Long:  `Returns the owner of the DAO (the account that can send/burn coins from the dao)`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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
	Long:  `Returns the access control list of governance params (which account can change the param)`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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
	Long:  `Returns the latest protocol upgrade by governance`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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
