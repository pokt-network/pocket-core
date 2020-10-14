package cli

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/app/cmd/rpc"
	"github.com/spf13/cobra"
)

var (
	datadir         string
	tmNode          string
	remoteCLIURL    string
	persistentPeers string
	seeds           string
	simulateRelay   bool
	keybase         bool
	mainnet         bool
	testnet         bool
)

var CLIVersion = app.AppVersion

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
		log.Fatal(err)
	}
}

func init() {
	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "help message for toggle")
	rootCmd.PersistentFlags().StringVar(&datadir, "datadir", "", "data directory (default is $HOME/.pocket/")
	rootCmd.PersistentFlags().StringVar(&tmNode, "node", "", "takes a remote endpoint in the form <protocol>://<host>:<port>")
	rootCmd.PersistentFlags().StringVar(&remoteCLIURL, "remoteCLIURL", "", "takes a remote endpoint in the form of <protocol>://<host> (uses RPC Port)")
	rootCmd.PersistentFlags().StringVar(&persistentPeers, "persistent_peers", "", "a comma separated list of PeerURLs: '<ID>@<IP>:<PORT>,<ID2>@<IP2>:<PORT>...<IDn>@<IPn>:<PORT>'")
	rootCmd.PersistentFlags().StringVar(&seeds, "seeds", "", "a comma separated list of PeerURLs: '<ID>@<IP>:<PORT>,<ID2>@<IP2>:<PORT>...<IDn>@<IPn>:<PORT>'")
	startCmd.Flags().BoolVar(&simulateRelay, "simulateRelay", false, "would you like to be able to test your relays")
	startCmd.Flags().BoolVar(&keybase, "keybase", true, "run with keybase, if disabled allows you to stake for the current validator only. providing a keybase is still neccesary for staking for apps & sending transactions")
	startCmd.Flags().BoolVar(&mainnet, "mainnet", false, "run with mainnet genesis")
	startCmd.Flags().BoolVar(&testnet, "testnet", false, "run with testnet genesis")
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(version)
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start --keybase=<true or false>",
	Short: "starts pocket-core daemon",
	Long:  `Starts the Pocket node, picks up the config from the assigned <datadir>`,
	Run: func(cmd *cobra.Command, args []string) {
		var genesisType app.GenesisType
		if mainnet && testnet {
			fmt.Println("cannot run with mainnet and testnet genesis simultaneously, please choose one")
			return
		}
		if mainnet {
			genesisType = app.MainnetGenesisType
		}
		if testnet {
			genesisType = app.TestnetGenesisType
		}
		tmNode := app.InitApp(datadir, tmNode, persistentPeers, seeds, remoteCLIURL, keybase, genesisType)
		go rpc.StartRPC(app.GlobalConfig.PocketConfig.RPCPort, app.GlobalConfig.PocketConfig.RPCTimeout, simulateRelay)
		// trap kill signals (2,3,15,9)
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel,
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGQUIT,
			os.Kill, //nolint
			os.Interrupt)

		defer func() {
			sig := <-signalChannel
			app.ShutdownPocketCore()
			err := tmNode.Stop()
			if err != nil {
				fmt.Println(err)
				return
			}
			message := fmt.Sprintf("Exit signal %s received\n", sig)
			fmt.Println(message)
			os.Exit(3)
		}()
	},
}

// startCmd represents the start command
var resetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset pocket-core",
	Long:  `Reset the Pocket node daemon`,
	Run:   app.ResetWorldState,
}

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop pocket-core",
	Long:  `Stop the Pocket node daemon`,
	Run:   app.ShutdownPocketCore,
}

var version = &cobra.Command{
	Use:   "version",
	Short: "Get current version",
	Long:  `Retrieves the version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("AppVersion: %s\n", CLIVersion)
	},
}
