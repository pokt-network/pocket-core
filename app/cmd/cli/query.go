package cli

import (
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	appTypes "github.com/pokt-network/pocket-core/x/apps/types"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

func init() {
	rootCmd.AddCommand(queryCmd)
	queryCmd.AddCommand(queryBlock)
	queryCmd.AddCommand(queryHeight)
	queryCmd.AddCommand(queryTx)
	queryCmd.AddCommand(queryNodes)
	queryCmd.AddCommand(queryBalance)
	queryCmd.AddCommand(queryAccount)
	queryCmd.AddCommand(queryNode)
	queryCmd.AddCommand(queryApps)
	queryCmd.AddCommand(queryApp)
	queryCmd.AddCommand(queryNodeParams)
	queryCmd.AddCommand(queryAppParams)
	queryCmd.AddCommand(queryNodeProofs)
	queryCmd.AddCommand(queryNodeProof)
	queryCmd.AddCommand(queryPocketParams)
	queryCmd.AddCommand(queryPocketSupportedChains)
	queryCmd.AddCommand(querySupply)
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
	Use:   "tx <height>",
	Short: "Get the transaction by the hash",
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

var queryNodeProofs = &cobra.Command{
	Use:   "node-proofs <nodeAddr> <height>",
	Short: "Gets node proofs",
	Long:  `Returns the list of all Relay Batch proofs submitted by <nodeAddr> at <height>.`,
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
		res, err := app.QueryProofs(args[0], int64(height))
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

var queryNodeProof = &cobra.Command{
	Use:   "node-proof <nodeAddr> <appPubKey> <networkId> <sessionHeight> <height>`",
	Short: "Gets node proof",
	Long:  `Gets node proof for specific session`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		var height int
		if len(args) == 4 {
			height = 0 // latest
		} else {
			var err error
			height, err = strconv.Atoi(args[4])
			if err != nil {
				fmt.Println(err)
				return
			}
		}
		sessionheight, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Println(err)
			return
		}
		res, err := app.QueryProof(args[2], args[1], args[0], int64(sessionheight), int64(height))
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
