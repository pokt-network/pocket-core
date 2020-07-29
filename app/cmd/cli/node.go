package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/types"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(nodesCmd)
	nodesCmd.AddCommand(nodeStakeCmd)
	nodesCmd.AddCommand(nodeUnstakeCmd)
	nodesCmd.AddCommand(nodeUnjailCmd)
}

var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "node management",
	Long: `The node namespace handles all node related interactions,
from staking and unstaking; to unjailing.`,
}

var nodeStakeCmd = &cobra.Command{
	Use:   "stake <fromAddr> <amount> <chains> <serviceURI> <chainID> <fees>",
	Short: "Stake a node in the network",
	Long: `Stake the node into the network, making it available for service.
Will prompt the user for the <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(6),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		fromAddr := args[0]
		amount, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		reg, err := regexp.Compile("[^,a-fA-F0-9]+")
		if err != nil {
			log.Fatal(err)
		}
		rawChains := reg.ReplaceAllString(args[2], "")
		chains := strings.Split(rawChains, ",")
		serviceURI := args[3]
		fees, err := strconv.Atoi(args[5])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Enter Passphrase: ")
		res, err := StakeNode(chains, serviceURI, fromAddr, app.Credentials(), args[4], types.NewInt(int64(amount)), int64(fees))
		if err != nil {
			fmt.Println(err)
			return
		}
		j, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
			return
		}
		resp, err := QueryRPC(SendRawTxPath, j)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(resp)
	},
}

var nodeUnstakeCmd = &cobra.Command{
	Use:   "unstake <fromAddr> <chainID> <fees>",
	Short: "Unstake a node in the network",
	Long: `Unstake a node from the network, changing it's status to Unstaking.
Will prompt the user for the <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		fees, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Enter Password: ")
		res, err := UnstakeNode(args[0], app.Credentials(), args[1], int64(fees))
		if err != nil {
			fmt.Println(err)
			return
		}
		j, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
			return
		}
		resp, err := QueryRPC(SendRawTxPath, j)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(resp)
	},
}

var nodeUnjailCmd = &cobra.Command{
	Use:   "unjail <fromAddr> <chainID> <fees>",
	Short: "Unjails a node in the network",
	Long: `Unjails a node from the network, allowing it to participate in service and consensus again.
Will prompt the user for the <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		fees, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Enter Password: ")
		res, err := UnjailNode(args[0], app.Credentials(), args[1], int64(fees))
		if err != nil {
			fmt.Println(err)
			return
		}
		j, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
			return
		}
		resp, err := QueryRPC(SendRawTxPath, j)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(resp)
	},
}
