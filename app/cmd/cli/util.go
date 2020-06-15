package cli

import (
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(utilCmd)
	utilCmd.AddCommand(chainsGenCmd)
	utilCmd.AddCommand(chainsDelCmd)
}

var utilCmd = &cobra.Command{
	Use:   "util",
	Short: "utility functions",
	Long:  `The util namespace handles all utility functions`,
}

var chainsGenCmd = &cobra.Command{
	Use:   "generate-chains",
	Short: "Generates chains file",
	Long:  `Generate the chains file for network identifiers`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		c := app.NewHostedChains(true)
		fmt.Println(app.GlobalConfig.PocketConfig.ChainsName + " contains: \n")
		for _, chain := range c.M {
			fmt.Println(chain.ID + " @ " + chain.URL)
		}
		fmt.Println("If incorrect: please remove the chains.json with the " + chainsDelCmd.NameAndAliases() + " command")
	},
}

var chainsDelCmd = &cobra.Command{
	Use:   "delete-chains",
	Short: "Delete chains file",
	Long:  `Delete the chains file for network identifiers`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		app.DeleteHostedChains()
		fmt.Println("successfully deleted " + app.GlobalConfig.PocketConfig.ChainsName)
	},
}
