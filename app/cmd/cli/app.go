package cli

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/spf13/cobra"

	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/crypto/keys/mintkey"
	"github.com/pokt-network/pocket-core/types"
)

func init() {
	rootCmd.AddCommand(appCmd)
	appCmd.AddCommand(appStakeCmd)
	appCmd.AddCommand(appUnstakeCmd)
	appCmd.AddCommand(appTransferCmd)
	appCmd.AddCommand(createAATCmd)
}

var appCmd = &cobra.Command{
	Use:   "apps",
	Short: "application management",
	Long: `The apps namespace handles all application related interactions,
from staking and unstaking; to generating AATs.`,
}

func init() {
	appStakeCmd.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")
	appUnstakeCmd.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")
	createAATCmd.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")
}

var appStakeCmd = &cobra.Command{
	Use:   "stake <fromAddr> <amount> <relayChainIDs> <networkID> <fee> ",
	Short: "Stake an app into the network",
	Long: `Stake the app into the network, giving it network throughput for the selected chains.
Will prompt the user for the <fromAddr> account passphrase. After the 0.6.X upgrade, if the app is already staked, this transaction acts as an *update* transaction.
A app can updated relayChainIDs, and raise the stake/max_relays amount with this transaction.
If the app is currently staked at X and you submit an update with new stake Y. Only Y-X will be subtracted from an account
If no changes are desired for the parameter, just enter the current param value just as before`,
	Args: cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		fromAddr := args[0]
		amount, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		reg, err := regexp.Compile("[^,a-fA-F0-9]+")
		if err != nil {
			log.Fatal(err)
		}
		fee, err := strconv.Atoi(args[4])
		if err != nil {
			fmt.Println(err)
			return
		}
		rawChains := reg.ReplaceAllString(args[2], "")
		chains := strings.Split(rawChains, ",")
		fmt.Println("Enter passphrase: ")
		res, err := StakeApp(chains, fromAddr, app.Credentials(pwd), args[3], types.NewInt(int64(amount)), int64(fee), false)
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
	Use:   "unstake <fromAddr> <networkID> <fee>",
	Short: "Unstake an app from the network",
	Long: `Unstake an app from the network, changing it's status to Unstaking.
Prompts the user for the <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		fee, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Enter Password: ")
		res, err := UnstakeApp(args[0], app.Credentials(pwd), args[1], int64(fee), false)
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

var appTransferCmd = &cobra.Command{
	Use:   "transfer <fromAddr> <newAppPubKey> <networkID> <fee> [memo]",
	Short: "Transfer the ownership of a staked app from one to another",
	Long: `Submits a transaction to transfer the ownership of a staked app from
<fromAddr> to a new account specified as <newAppPubKey> without unstaking
any app.  In other words, this edits the address of a staked app.  To run this
command, you must have the private key of the current staked app <fromAddr>
`,
	Args: cobra.MinimumNArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)

		currentAppAddr := args[0]
		newAppPubKey := args[1]
		networkId := args[2]
		feeStr := args[3]
		memo := ""
		if len(args) >= 5 {
			memo = args[4]
		}

		fee, err := strconv.ParseInt(feeStr, 10, 64)
		if err != nil {
			fmt.Println("Invalid fee:", err)
			return
		}

		fmt.Printf("Enter passphrase to unlock %s: ", currentAppAddr)
		passphrase := app.Credentials(pwd)

		rawTx, err := TransferApp(
			currentAppAddr,
			newAppPubKey,
			passphrase,
			networkId,
			fee,
			memo,
		)
		if err != nil {
			fmt.Println("Failed to build a transaction:", err)
			return
		}

		rawTxBytes, err := json.Marshal(rawTx)
		if err != nil {
			fmt.Println(err)
			return
		}

		resp, err := QueryRPC(SendRawTxPath, rawTxBytes)
		if err != nil {
			fmt.Println("Failed to submit a transaction:", err)
			return
		}

		fmt.Println(resp)
	},
}

var createAATCmd = &cobra.Command{
	Use:   "create-aat <appAddr> <clientPubKey>",
	Short: "Creates an application authentication token",
	Long: `Creates a signed Application Authentication Token.
Creates a signed AAT (= Application Authentication Token) where the version is
hardcoded as "0.0.1" that is the only version supported by the protocol.

This command prompts you to input the <appAddr> account passphrase.
When you send a relay request with AAT, <appAddr> needs to be a staked
application.

Please read doc/specs/application-auth-token.md for additional details.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
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
		kp, err := kb.Get(addr)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Enter passphrase: ")
		cred := app.Credentials(pwd)

		// Retrieve the priv & public keys
		pubKey := kp.PublicKey.RawBytes()
		pubKeyHexEncoded := hex.EncodeToString(pubKey)
		privKeyArmored := kp.PrivKeyArmor
		privKey, err := mintkey.UnarmorDecryptPrivKey(privKeyArmored, cred)
		if err != nil {
			return
		}

		aat, err := app.GenerateAAT(pubKeyHexEncoded, args[1], privKey)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(string(aat))
	},
}
