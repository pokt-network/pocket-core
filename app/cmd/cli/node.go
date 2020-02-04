package cli

import (
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/posmint/types"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
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
	Long:  ``,
}

var nodeStakeCmd = &cobra.Command{
	Use:   "stake <fromAddr> <amount> <chains> <serviceURI>",
	Short: "Stake a node in the network",
	Long:  `Stake the node into the network, making it available for service. Prompts the user for the <fromAddr> account passphrase.`,
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		fromAddr := args[0]
		amount, err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}
		chains := strings.Split(args[2], ",")
		serviceURI := args[3]
		fmt.Println("Enter Password: ")
		res, err := app.StakeNode(chains, serviceURI, fromAddr, app.Credentials(), types.NewInt(int64(amount)))
		if err != nil {
			panic(err)
		}
		fmt.Printf("Transaction Submitted: %s\n", res.TxHash)
	},
}

var nodeUnstakeCmd = &cobra.Command{
	Use:   "unstake <fromAddr>",
	Short: "Unstake a node in the network",
	Long:  `Unstake a node from the network, changing it's status to Unstaking. Prompts the user for the <fromAddr> account passphrase.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Enter Password: ")
		res, err := app.UnstakeNode(args[0], app.Credentials())
		if err != nil {
			panic(err)
		}
		fmt.Printf("Transaction Submitted: %s\n", res.TxHash)
	},
}

var nodeUnjailCmd = &cobra.Command{
	Use:   "unjail <fromAddr>",
	Short: "Unjails a node in the network",
	Long:  `Unjails a node from the network, allowing it to participate in service and consensus again. Prompts the user for the <fromAddr> account passphrase.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Enter Password: ")
		res, err := app.UnjailNode(args[0], app.Credentials())
		if err != nil {
			panic(err)
		}
		fmt.Printf("Transaction Submitted: %s\n", res.TxHash)
	},
}
