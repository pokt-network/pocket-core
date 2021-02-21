package cli

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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
	remoteCLIURL    string
	persistentPeers string
	seeds           string
	simulateRelay   bool
	keybase         bool
	mainnet         bool
	testnet         bool
	profileApp      bool
	madvdontneed    bool
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
	rootCmd.PersistentFlags().BoolVar(&madvdontneed, "madvdontneed", true, "if enabled, run with GODEBUG=madvdontneed=1, --madvdontneed=true/false")
	startCmd.Flags().BoolVar(&simulateRelay, "simulateRelay", false, "would you like to be able to test your relays")
	startCmd.Flags().BoolVar(&keybase, "keybase", true, "run with keybase, if disabled allows you to stake for the current validator only. providing a keybase is still neccesary for staking for apps & sending transactions")
	startCmd.Flags().BoolVar(&mainnet, "mainnet", false, "run with mainnet genesis")
	startCmd.Flags().BoolVar(&testnet, "testnet", false, "run with testnet genesis")
	startCmd.Flags().BoolVar(&profileApp, "profileApp", false, "expose cpu & memory profiling")
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(resetCmd)
	rootCmd.AddCommand(version)
	rootCmd.AddCommand(stopCmd)
}

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start --keybase=<true or false>",
	Short: "starts pocket-core daemon",
	Long:  `Starts the Pocket node, picks up the config from the assigned <datadir>`,
	Run: func(cmd *cobra.Command, args []string) {

		//Get the GODEBUG env variable
		godebug := os.Getenv("GODEBUG")
		//Check if the --madvdontneed=true

		//Check if madvdontneed env variable is present or flag is not used
		if strings.Contains(godebug, "madvdontneed=1") || !madvdontneed {
			//start normally
			start(cmd, args)

		} else {
			//flag --madvdontneed=true so we add the env variable and start pocket as a subprocess
			env := append(os.Environ(), "GODEBUG="+"madvdontneed=1,"+godebug)
			comd := exec.Command(os.Args[0], os.Args[1:]...)
			comd.Env = env
			comd.Stdin = os.Stdin
			comd.Stdout = os.Stdout

			if err := comd.Start(); err != nil {
				log.Fatalf("couldn't start child process command: %v", err)
			}

			<-waitOnStopSignals()

			if err := comd.Wait(); err != nil {
				log.Fatalf("couldn't wait for child process to terminate: %v", err)
			}
		}
	},
}

func waitOnStopSignals() chan os.Signal {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig,
		syscall.SIGTERM,
		syscall.SIGINT,
		syscall.SIGQUIT,
		os.Kill,
		os.Interrupt,
	)
	return sig
}

func start(cmd *cobra.Command, args []string) {
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
	go rpc.StartRPC(app.GlobalConfig.PocketConfig.RPCPort, app.GlobalConfig.PocketConfig.RPCTimeout, simulateRelay, profileApp)

	defer func() {
		sig := <-waitOnStopSignals()

		app.ShutdownPocketCore()
		err := tmNode.Stop()
		if err != nil {
			fmt.Println(err)
			return
		}
		message := fmt.Sprintf("Exit signal %s received\n", sig)
		fmt.Println(message)
		os.Exit(0)
	}()
}

// resetCmd represents the reset command
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

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop pocket-core",
	Long:  `Stop pocket-core`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		res, err := QuerySecuredRPC(GetStopPath, []byte{}, app.GetAuthTokenFromFile())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(res)
	},
}
