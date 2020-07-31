package cli

import (
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	"github.com/spf13/cobra"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/state"
	"os"
	"strconv"
)

func init() {
	rootCmd.AddCommand(utilCmd)
	utilCmd.AddCommand(chainsGenCmd)
	utilCmd.AddCommand(chainsDelCmd)
	utilCmd.AddCommand(decodeTxCmd)
	utilCmd.AddCommand(unsafeRollbackCmd)
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

var decodeTxCmd = &cobra.Command{
	Use:   "decode-tx <tx>",
	Short: "Decodes a given transaction encoded in Amino base64 bytes",
	Long:  `Decodes a given transaction encoded in Amino base64 bytes`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		txStr := args[0]
		stdTx := app.UnmarshalTxStr(txStr)
		fmt.Printf(
			"Type:\t\t%s\nMsg:\t\t%v\nFee:\t\t%s\nEntropy:\t%d\nMemo:\t\t%s\nSigner\t\t%s\nSig:\t\t%s\n",
			stdTx.Msg.Type(), stdTx.Msg, stdTx.Fee.String(), stdTx.Entropy, stdTx.Memo, stdTx.Msg.GetSigner().String(),
			stdTx.Signature.RawString())
	},
}

func init() {
	unsafeRollbackCmd.Flags().BoolVar(&blocks, "blocks", false, "rollback blocks as well as the state")
}

var (
	blocks       bool
)

var unsafeRollbackCmd = &cobra.Command{
	Use:   "unsafe-rollback <height",
	Short: "Rollbacks the blockchain, the state, and app to a previous height",
	Long:  "Rollbacks the blockchain, the state, and app to a previous height",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		height, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println("error parsing height: ", err)
			return
		}
		db, err := app.OpenDB(app.GlobalConfig.TendermintConfig.RootDir)
		if err != nil {
			fmt.Println("error loading application database: ", err)
			return
		}
		loggerFile, _ := os.Open(os.DevNull)
		a := app.NewPocketBaseApp(log.NewTMLogger(loggerFile), db)
		// initialize stores
		a.MountKVStores(a.Keys)
		a.MountTransientStores(a.Tkeys)
		// rollback the txIndexer
		err = state.RollbackTxIndexer(&app.GlobalConfig.TendermintConfig, int64(height))
		if err != nil {
			fmt.Println("error rolling back txIndexer: ", err)
			return
		}
		// rollback the app store
		err = a.Store().RollbackVersion(int64(height))
		if err != nil {
			fmt.Println("error rolling back app: ", err)
			return
		}
		if blocks {
			// rollback block store and state
			err = state.UnsafeRollbackData(&app.GlobalConfig.TendermintConfig, true, int64(height))
			if err != nil {
				fmt.Println("error rolling back block and state: ", err)
				return
			}
		}
	},
}
