package cli

import (
	"encoding/hex"
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/posmint/types"
	"github.com/spf13/cobra"
	"strconv"
	"strings"
)

func init() {
	rootCmd.AddCommand(appCmd)
	appCmd.AddCommand(appStakeCmd)
	appCmd.AddCommand(appUnstakeCmd)
	appCmd.AddCommand(createAATCmd)
}

var appCmd = &cobra.Command{
	Use:   "apps",
	Short: "Functions for app management",
	Long:  ``,
}

var appStakeCmd = &cobra.Command{
	Use:   "stake <fromAddr> <amount> <chains>",
	Short: "Stake an app in the network",
	Long:  `Stakes the app into the network, making it have network throughjput. Prompts the user for the <fromAddr> account passphrase.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		fromAddr := args[0]
		amount, err := strconv.Atoi(args[1])
		if err != nil {
			panic(err)
		}
		chains := strings.Split(args[2], ",")
		fmt.Println("Enter Password: ")
		res, err := app.StakeApp(chains, fromAddr, app.Credentials(), types.NewInt(int64(amount)))
		if err != nil {
			panic(err)
		}
		fmt.Printf("Transaction Submitted %s", res.TxHash)
	},
}

var appUnstakeCmd = &cobra.Command{
	Use:   "unstake <fromAddr>",
	Short: "Unstake an app in the network",
	Long:  `Unstakes an app from the network, changing it's status to Unstaking. Prompts the user for the <fromAddr> account passphrase.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		res, err := app.UnstakeApp(args[0], app.Credentials())
		if err != nil {
			panic(err)
		}
		fmt.Printf("Transaction Submitted %s", res.TxHash)
	},
}

var createAATCmd = &cobra.Command{
	Use:   "create-aat <appAddr> <clientPubKey>",
	Short: "Creates an application authentication token",
	Long: `Creates a signed application authentication token (version 0.0.1 of the AAT spec), that can be embedded into application software for Relay servicing. 
Will prompt the user for the <appAddr> account passphrase. 
Read the Application Authentication Token documentation for more. 
NOTE: USE THIS METHOD AT YOUR OWN RISK. READ THE APPLICATION SECURITY GUIDELINES IN ORDER TO UNDERSTAND WHAT'S THE RECOMMENDED AAT CONFIGURATION FOR YOUR APPLICATION:`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		kb := app.GetKeybase()
		if kb == nil {
			panic(app.UninitializedKeybaseError)
		}
		addr, err := types.AccAddressFromHex(args[0])
		if err != nil {
			panic(err)
		}
		res, err := (*kb).Get(addr)
		if err != nil {
			panic(err)
		}
		fmt.Println("Enter Password: ")
		aatBytes, err := app.GenerateAAT(hex.EncodeToString(res.PubKey.Bytes()), args[1], app.Credentials())
		fmt.Println(string(aatBytes))
	},
}
