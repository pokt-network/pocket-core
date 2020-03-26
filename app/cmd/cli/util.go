package cli

import (
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(utilCmd)
	utilCmd.AddCommand(generateChainCmd)
	utilCmd.AddCommand(exportAppStateCmd)
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
	Args:  cobra.MinimumNArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		var client, inter string
		app.SetTMNode(tmNode)
		switch len(args) {
		case 4:
			client = args[3]
		case 5:
			client = args[3]
			inter = args[4]
		}
		res, err := app.GenerateChain(args[0], args[1], args[2], client, inter)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Pocket Network Identifier: %s\n", res)
	},
}

var exportAppStateCmd = &cobra.Command{
	Use:   "export-state <output-file>",
	Short: "export current state genesis",
	Long:  `Export the current app state in a genesis file`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		err := app.ExportState(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Successfully exported state")
	},
}
