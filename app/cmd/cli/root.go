package cli

import (
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	fs              = string(filepath.Separator)
	datadir         string
	persistentPeers = "" // todo pull from file
	seeds           = "" // todo pull from file
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "pocket",
	Short: "Pocket provides a trustless API Layer, allowing easy access to any blockchain.",
	Long: `Pocket is a distributed network that relays data requests and responses to and from any blockchain system. 
Pocket verifies all relayed data and proportionally rewards the participating nodes with native cryptographic tokens.
`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	startCmd.PersistentFlags().StringVar(&datadir, "data_dir", "", "data directory (default is $HOME/.pocket/")
	rootCmd.AddCommand(startCmd)
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the pocket-core client",
	Long:  `Starts the Pocket node, picks up the config from the assigned <datadir>`,
	Run: func(cmd *cobra.Command, args []string) {
		// We trap kill signals (2,3,15,9)
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel,
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGQUIT,
			os.Kill,
			os.Interrupt)
		// setup the codec
		app.SetCodec()
		// setup the data directory
		datadir = app.InitDataDirectory(datadir)
		// setup the keybase
		pswrd := app.InitKeyfiles(datadir)
		// setup the genesis.json
		app.InitGenesis(datadir + fs + "config")
		// setup the chains.json
		app.InitHostedChains(datadir + fs + "config")
		// setup coinbase password
		if pswrd == "" {
			fmt.Println("Pocket core needs your passphrase to start")
			pswrd = app.Credentials()
			err := app.ConfirmCoinbasePassword(pswrd)
			if err != nil {
				panic("Coinbase Password could not be verified: " + err.Error())
			}
		}
		app.SetCoinbasePassphrase(pswrd)
		// init the tendermint node
		app.InitTendermint(datadir, persistentPeers, seeds)
		// catch end signal
		defer func() {
			sig := <-signalChannel
			tmNode, _ := app.GetTendermintNode()
			err := tmNode.Stop()
			if err != nil {
				panic("unable to stop Tendermint node: " + err.Error())
			}
			message := fmt.Sprintf("Exit signal %s received\n", sig)
			fmt.Println(message)
			os.Exit(3)
		}()
	},
}
