package cli

import (
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(utilCmd)
	utilCmd.AddCommand(generateChainCmd)
}

// accountsCmd represents the accounts namespace command
var utilCmd = &cobra.Command{
	Use:   "util",
	Short: "Functions for pocket core util",
	Long:  ``,
}

var generateChainCmd = &cobra.Command{
	Use:   "util generate-chain <ticker> <netid> <client> <version> <interface>",
	Short: "Stake an app in the network",
	Long:  `Creates a Network Identifier hash, used as a parameter for both node and App stake.`,
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		res, err := app.GenerateChain(args[0], args[1], args[3], args[2], args[4])
		if err != nil {
			panic(err)
		}
		fmt.Printf("Pocket Network Identifier: %s", res)
	},
}
