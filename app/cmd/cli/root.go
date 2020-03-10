package cli

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/app/cmd/rpc"
	"github.com/spf13/cobra"
)

var (
	datadir         string
	tmNode          string
	persistentPeers string
	seeds           string
	tmRPCPort       string
	tmPeersPort     string
	pocketRPCPort   string
	blockTime       int
)

const defaultSeeds = `610CF8A6E8CEFBADED845F1C1DC3B10A670BE26B@node1.testnet.pokt.network:26656, E6946760D9833F49DA39AAE9500537BEF6F33A7A@node2.testnet.pokt.network:26656, 7674A47CC977326F1DF6CB92C7B5A2AD36557EA2@node3.testnet.pokt.network:26656, C7B7B7665D20A7172D0C0AA58237E425F333560A@node4.testnet.pokt.network:26656, F6DC0B244C93232283CD1D8443363946D0A3D77A@node5.testnet.pokt.network:26656, 86209713BEFECA0807714BCDD5B79E81073FAF8F@node6.testnet.pokt.network:26656, 915A58AE437D2C2D6F35AC11B79F42972267700D@node7.testnet.pokt.network:26656, B3D86CD8AB4AA0CB9861CB795D8D154E685A94CF@node8.testnet.pokt.network:26656, 17CA63E4FF7535A40512C550DD0267E519CAFC1A@node9.testnet.pokt.network:26656, F99386C6D7CD42A486C63CCD80F5FBEA68759CD7@node10.testnet.pokt.network:26656`

var CLIVersion = fmt.Sprintf("%s%s", app.Tag, app.Version)

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
	rootCmd.Flags().BoolP("toggle", "t", false, "help message for toggle")
	rootCmd.PersistentFlags().StringVar(&datadir, "datadir", "", "data directory (default is $HOME/.pocket/")
	rootCmd.PersistentFlags().StringVar(&tmNode, "node", "", "takes a remote endpoint in the form <protocol>://<host>:<port>")
	rootCmd.PersistentFlags().StringVar(&persistentPeers, "persistent_peers", "", "a comma separated list of PeerURLs: <ID>@<IP>:<PORT>")
	rootCmd.PersistentFlags().StringVar(&seeds, "seeds", defaultSeeds, "a comma separated list of PeerURLs: <ID>@<IP>:<PORT>")
	rootCmd.PersistentFlags().StringVar(&tmRPCPort, "tmRPCPort", "26657", "the port for tendermint rpc")
	rootCmd.PersistentFlags().StringVar(&tmPeersPort, "tmPeersPort", "26656", "the port for tendermint p2p")
	rootCmd.PersistentFlags().StringVar(&pocketRPCPort, "pocketRPCPort", "8081", "the port for pocket rpc")
	rootCmd.PersistentFlags().IntVar(&blockTime, "blockTime", 1, "how often should the network create blocks")
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(version)
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "starts pocket-core daemon",
	Long:  `Starts the Pocket node, picks up the config from the assigned <datadir>`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		go rpc.StartRPC(pocketRPCPort)
		tmNode := app.InitApp(datadir, tmNode, strings.ToLower(persistentPeers), strings.ToLower(seeds), tmRPCPort, tmPeersPort, blockTime)
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
	Short: "reset pocket-core",
	Long:  `Reset the Pocket node`,
	Run:   app.ResetWorldState,
}

var version = &cobra.Command{
	Use:   "version",
	Short: "Get current version",
	Long:  `Returns the version`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("Version: %s\n", CLIVersion)
	},
}
