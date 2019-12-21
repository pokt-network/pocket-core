package cli

import (
	"errors"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/posmint/types"
	"github.com/spf13/viper"
	"github.com/tendermint/tendermint/p2p"
	"github.com/tendermint/tendermint/privval"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
)

const (
	pubKeytype  = "tendermint/PubKeyEd25519"
	privKeyType = "tendermint/PrivKeyEd25519"
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
		setCodec()
		initDataDirectory()
		initKeybase()
		setGenesisFile()
		initChains()
		// remove below
		kps, _ := (*app.GetKeybase()).List()
		kp := kps[0]
		j, _ := app.Cdc.MarshalJSON(types.ValAddress(kp.PubKey.Address()))
		j3, _ := app.Cdc.MarshalJSON(types.ConsAddress(kp.PubKey.Address()))
		j2, _ := types.Bech32ifyConsPub(kp.PubKey)
		fmt.Println("VAL ADDR: -> " + string(j))
		fmt.Println("Cons ADDR: -> " + string(j3))
		fmt.Println("PUBKEY -> " + j2)
		// remove above
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

func initKeyfiles(passphrase string) {
	kb := app.GetKeybase()
	keypairs, err := (*kb).List()
	if err != nil {
		panic(err)
	}
	coinbaseKeypair := keypairs[0]
	res, err := (*kb).ExportPrivateKeyObject(coinbaseKeypair.GetAddress(), passphrase)
	if err != nil {
		panic(err)
	}
	privValKey := privval.FilePVKey{
		Address: res.PubKey().Address(),
		PubKey:  res.PubKey(),
		PrivKey: res,
	}
	privValState := privval.FilePVLastSignState{}
	nodeKey := p2p.NodeKey{
		PrivKey: res,
	}
	pvkBz, err := app.Cdc.MarshalJSONIndent(privValKey, "", "  ")
	if err != nil {
		panic(err)
	}
	nkBz, err := app.Cdc.MarshalJSONIndent(nodeKey, "", "  ")
	if err != nil {
		panic(err)
	}
	pvsBz, err := app.Cdc.MarshalJSONIndent(privValState, "", "  ")
	if err != nil {
		panic(err)
	}
	pvFile, err := os.OpenFile(datadir+fs+"priv_val_key.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	_, err = pvFile.Write(pvkBz)
	if err != nil {
		panic(err)
	}
	pvStateFile, err := os.OpenFile(datadir+fs+"priv_val_state.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	_, err = pvStateFile.Write(pvsBz)
	if err != nil {
		panic(err)
	}
	nkFile, err := os.OpenFile(datadir+fs+"node_key.json", os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		panic(err)
	}
	_, err = nkFile.Write(nkBz)
	if err != nil {
		panic(err)
	}
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
		datadir = home + fs + ".pocket"
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

func setCodec() {
	app.SetCodec()
}

func initChains() {
	app.InitHostedChains(datadir + fs + "config")
}

func setGenesisFile() {
	app.SetGenesisFilepath(datadir + fs + "config" + fs + "genesis.json")
}

func initKeybase() {
	if err := app.InitKeybase(datadir); err != nil {
		fmt.Println("Initializing keybase: enter coinbase passphrase")
		password := credentials()
		err := app.CreateKeybase(datadir, password)
		if err != nil {
			panic(err)
		}
		initKeyfiles(password)
	}
}

func initTendermint() {
	tmNode := app.InitTendermintNode(datadir, "", "node_key.json", "priv_val_key.json",
		"priv_val_state.json", persistentPeers, seeds, "0.0.0.0:46656")
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
