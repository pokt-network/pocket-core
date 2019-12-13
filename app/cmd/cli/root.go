package cli

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pokt-network/pocket-core/app"
	"github.com/spf13/viper"
	"log"
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
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(startCmd)
	initCmd.PersistentFlags().StringVar(&datadir, "data_dir", "", "data directory (default is $HOME/.pocket/")
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
		fmt.Println("Pocket core needs your passphrase to start")
		app.SetCoinbasePassphrase(credentials())
		initTendermint()
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

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		initDataDirectory()
		initKeybase()
		setGenesisFile()
		initChains()
	},
}

func initDataDirectory() {
	// check for empty data_dir
	if datadir == "" {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// set the default data directory
		datadir = home + fs + ".pocket" + fs
	}

	// setup config file
	viper.AddConfigPath(datadir + ".pocket")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
	// create the folder if not already created
	err := os.MkdirAll(datadir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}

func initChains() {
	app.InitHostedChains(datadir + fs + "config" + fs + "chains.json")
}

func setGenesisFile() {
	app.SetGenesisFilepath(datadir + fs + "config" + fs + "genesis.json")
}

func initKeybase() {
	if _, err := app.GetKeybase(); err != nil {
		fmt.Println("Initializing keybase: enter coinbase passphrase")
		password := credentials()
		app.InitKeybase(datadir, password)
	}
}

func initTendermint() {
	tmNode := app.InitTendermintNode(datadir, "", "", "", "", persistentPeers, seeds, "")
	if err := tmNode.Start(); err != nil {
		panic(err)
	}
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
		tmNode, _ := app.GetTendermintNode()
		err := tmNode.Stop()
		if err != nil {
			panic(err)
		}
		message := fmt.Sprintf("Exit signal %s received\n", sig)
		fmt.Println(message)
		os.Exit(3)
	}()
}

var (
	InitError = errors.New(" -> must run init command before any other")
)

func NewBeforeInitError(err error) error {
	return errors.New(err.Error() + InitError.Error())
}
