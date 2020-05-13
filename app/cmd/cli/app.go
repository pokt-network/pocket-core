package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/posmint/types"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(appCmd)
	appCmd.AddCommand(appStakeCmd)
	appCmd.AddCommand(appUnstakeCmd)
	appCmd.AddCommand(createAATCmd)
}

var appCmd = &cobra.Command{
	Use:   "apps",
	Short: "application management",
	Long: `The apps namespace handles all applicaiton related interactions,
from staking and unstaking; to generating AATs.`,
}

var appStakeCmd = &cobra.Command{
	Use:   "stake <fromAddr> <amount> <chains> <chainID>",
	Short: "Stake an app into the network",
	Long: `Stake the app into the network, making it have network throughput.
Will prompt the user for the <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		fromAddr := args[0]
		amount, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		reg, err := regexp.Compile("[^,a-zA-Z0-9]+")
		if err != nil {
			log.Fatal(err)
		}
		rawChains := reg.ReplaceAllString(args[2], "")
		chains := strings.Split(rawChains, ",")
		fmt.Println("Enter passphrase: ")
		res, err := StakeApp(chains, fromAddr, args[3], app.Credentials(), types.NewInt(int64(amount)))
		if err != nil {
			fmt.Println(err)
			return
		}
		j, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
			return
		}
		resp, err := QueryRPC(SendRawTxPath, j)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(resp)
	},
}

var appUnstakeCmd = &cobra.Command{
	Use:   "unstake <fromAddr> <chainID>",
	Short: "Unstake an app from the network",
	Long: `Unstake an app from the network, changing it's status to Unstaking.
Prompts the user for the <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		fmt.Println("Enter Password: ")
		res, err := UnstakeApp(args[0], app.Credentials(), args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		j, err := json.Marshal(res)
		if err != nil {
			fmt.Println(err)
			return
		}
		resp, err := QueryRPC(SendRawTxPath, j)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println(resp)
	},
}

var createAATCmd = &cobra.Command{
	Use:   "create-aat <appAddr> <clientPubKey>",
	Short: "Creates an application authentication token",
	Long: `Creates a signed application authentication token (version 0.0.1 of the AAT spec), that can be embedded into application software for Relay servicing.
Will prompt the user for the <appAddr> account passphrase.
Read the Application Authentication Token documentation for more information.
NOTE: USE THIS METHOD AT YOUR OWN RISK. READ THE APPLICATION SECURITY GUIDELINES IN ORDER TO UNDERSTAND WHAT'S THE RECOMMENDED AAT CONFIGURATION FOR YOUR APPLICATION.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Println(app.UninitializedKeybaseError)
			return
		}
		addr, err := types.AddressFromHex(args[0])
		if err != nil {
			fmt.Printf("Address Error %s", err)
			return
		}
		res, err := kb.Get(addr)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Enter passphrase: ")
		aatBytes, err := app.PCA.GenerateAAT(hex.EncodeToString(res.PublicKey.RawBytes()), args[1], app.Credentials())
		fmt.Println(string(aatBytes))
	},
}
