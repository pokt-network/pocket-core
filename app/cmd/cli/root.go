package cli

import (
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/app/cmd/rpc"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	datadir         string
	tmNode          string
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
	startCmd.PersistentFlags().StringVar(&datadir, "node", "", "takes a remote endpoint in the form <protocol>://<host>:<port>")
	rootCmd.AddCommand(startCmd)
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts the pocket-core client",
	Long:  `Starts the Pocket node, picks up the config from the assigned <datadir>`,
	Run: func(cmd *cobra.Command, args []string) {
		go rpc.StartRPC("8081")
		// setup the codec
		app.MakeCodec()
		// setup the data directory
		app.InitDataDirectory(datadir)
		// setup the keybase
		pswrd := app.InitKeyfiles()
		// setup the genesis.json
		app.InitGenesis()
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
		// set tendermint node
		app.SetTendermintNode(tmNode)
		// init the tendermint node
		tmNode := app.InitTendermint(persistentPeers, seeds)
		// We trap kill signals (2,3,15,9)
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel,
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGQUIT,
			os.Kill,
			os.Interrupt)

		defer func() {
			sig := <-signalChannel
			err := tmNode.Stop()
			if err != nil {
				panic(err)
			}
			message := fmt.Sprintf("Exit signal %s received\n", sig)
			fmt.Println(message)
			os.Exit(3)
		}()
	},
}
