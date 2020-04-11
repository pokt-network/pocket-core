package cli

import (
	"encoding/hex"
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
	Use:   "stake <fromAddr> <amount> <chains>",
	Short: "Stake an app into the network",
	Long: `Stake the app into the network, making it have network throughput.
Will prompt the user for the <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
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
		res, err := app.StakeApp(chains, fromAddr, app.Credentials(), types.NewInt(int64(amount)))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Transaction Submitted: %s\n", res.TxHash)
	},
}

var appUnstakeCmd = &cobra.Command{
	Use:   "unstake <fromAddr>",
	Short: "Unstake an app from the network",
	Long: `Unstake an app from the network, changing it's status to Unstaking.
Prompts the user for the <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		fmt.Println("Enter passphrase: ")
		res, err := app.UnstakeApp(args[0], app.Credentials())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Transaction Submitted: %s\n", res.TxHash)
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
		app.SetTMNode(tmNode)
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
		aatBytes, err := app.GenerateAAT(hex.EncodeToString(res.PublicKey.RawBytes()), args[1], app.Credentials())
		fmt.Println(string(aatBytes))
	},
}
