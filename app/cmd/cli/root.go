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
	tmRPCPort       string
	tmPeersPort     string
	pocketRPCPort   string
	blockTime       int
	testnet         bool
	simulateRelay   bool
	keybase         bool
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
	rootCmd.PersistentFlags().StringVar(&tmRPCPort, "tmRPCPort", "", "the port for tendermint rpc")
	rootCmd.PersistentFlags().StringVar(&tmPeersPort, "tmPeersPort", "", "the port for tendermint p2p")
	rootCmd.PersistentFlags().StringVar(&pocketRPCPort, "pocketRPCPort", "", "the port for pocket rpc")
	rootCmd.PersistentFlags().IntVar(&blockTime, "blockTime", 1, "how often should the network create blocks")
	rootCmd.PersistentFlags().BoolVar(&testnet, "testnet", false, "would you like to connect to Pocket Network testnet")
	rootCmd.PersistentFlags().BoolVar(&simulateRelay, "simulateRelay", false, "would you like to be able to test your relays")
	rootCmd.Flags().BoolVar(&keybase, "keybase", true, "run wiith keybase, if disabled allows you to stake for the current validator only. providing a keybase is still neccesary for staking for apps & sending transactions")
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(version)
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start --keybase=<keybase>",
	Short: "starts pocket-core daemon",
	Long:  `Starts the Pocket node, picks up the config from the assigned <datadir>`,
	Run: func(cmd *cobra.Command, args []string) {
		tmNode := app.InitApp(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort, remoteCLIURL, keybase)
		go rpc.StartRPC(app.GlobalConfig.PocketConfig.RPCPort, simulateRelay)
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

var version = &cobra.Command{
	Use:   "version",
	Short: "Get current version",
	Long:  `Retrieves the version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("AppVersion: %s\n", CLIVersion)
	},
}
