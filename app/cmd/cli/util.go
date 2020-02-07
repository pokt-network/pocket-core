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
	Short: "utilities",
	Long:  ``,
}

var generateChainCmd = &cobra.Command{
	Use:   "generate-chain <ticker> <netid> <version> <client> <interface>",
	Short: "generate a chain identifier",
	Long:  `Creates a Network Identifier hash, used as a parameter for both node and App stake.`,
	Args:  cobra.MinimumNArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		res, err := app.GenerateChain(args[0], args[1], args[2], args[3], args[4])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Pocket Network Identifier: %s\n", res)
	},
}
