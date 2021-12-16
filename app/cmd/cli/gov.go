package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/types"
	govTypes "github.com/pokt-network/pocket-core/x/gov/types"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(govCmd)
	govCmd.AddCommand(govDAOTransfer)
	govCmd.AddCommand(govDAOBurn)
	govCmd.AddCommand(govChangeParam)
	govCmd.AddCommand(govUpgrade)
	govCmd.AddCommand(govFeatureEnable)
}

var govCmd = &cobra.Command{
	Use:   "gov",
	Short: "governance management",
	Long: `The gov namespace handles all governance related interactions,
from DAOTransfer, change parameters; to performing protocol Upgrades. `,
}

func init() {
	govDAOTransfer.Flags().StringVar(&pwd, "pwd", "", "defines the passphrase used by the cmd non empty usage bypass interactive prompt ")
	govDAOBurn.Flags().StringVar(&pwd, "pwd", "", "defines the passphrase used by the cmd non empty usage bypass interactive prompt ")
	govChangeParam.Flags().StringVar(&pwd, "pwd", "", "defines the passphrase used by the cmd non empty usage bypass interactive prompt ")
	govUpgrade.Flags().StringVar(&pwd, "pwd", "", "defines the passphrase used by the cmd non empty usage bypass interactive prompt ")
}

var govDAOTransfer = &cobra.Command{
	Use:   "transfer <amount> <fromAddr> <toAddr> <networkID> <fees>",
	Short: "Transfer from DAO",
	Long: `If authorized, move funds from the DAO.
Actions: [burn, transfer]`,
	Args: cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		toAddr := args[2]
		fromAddr := args[1]
		amount, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		fees, err := strconv.Atoi(args[4])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Enter Password: ")
		pass := app.Credentials(pwd)
		res, err := DAOTx(fromAddr, toAddr, pass, types.NewInt(int64(amount)), "dao_transfer", args[3], int64(fees), false)
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

var govDAOBurn = &cobra.Command{
	Use:   "burn <amount> <fromAddr> <toAddr> <networkID> <fees>",
	Short: "Burn from DAO",
	Long: `If authorized, burn funds from the DAO.
Actions: [burn, transfer]`,
	Args: cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		var toAddr string
		if len(args) == 4 {
			toAddr = args[2]
		}
		fromAddr := args[1]
		amount, err := strconv.Atoi(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		fees, err := strconv.Atoi(args[4])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Enter Password: ")
		pass := app.Credentials(pwd)
		res, err := DAOTx(fromAddr, toAddr, pass, types.NewInt(int64(amount)), "dao_burn", args[3], int64(fees), false)
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
var govChangeParam = &cobra.Command{
	Use:   "change_param <fromAddr> <networkID> <paramKey module/param> <paramValue (jsonObj)> <fees>",
	Short: "Edit a param in the network",
	Long: `If authorized, submit a tx to change any param from any module.
Will prompt the user for the <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		fmt.Println("Enter Password: ")
		fees, err := strconv.Atoi(args[4])
		if err != nil {
			fmt.Println(err)
			return
		}

		res, err := ChangeParam(args[0], args[2], []byte(args[3]), app.Credentials(pwd), args[1], int64(fees), false)
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

var govUpgrade = &cobra.Command{
	Use:   "upgrade <fromAddr> <atHeight> <version> <networkID> <fees>",
	Short: "Upgrade the protocol",
	Long: `If authorized, upgrade the protocol.
Will prompt the user for the <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		i, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal(err)
		}
		u := govTypes.Upgrade{
			Height:  int64(i),
			Version: dropTag(args[2]),
		}
		fees, err := strconv.Atoi(args[4])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Enter Password: ")
		res, err := Upgrade(args[0], u, app.Credentials(pwd), args[3], int64(fees), false)
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

func dropTag(version string) string {
	if !strings.Contains(version, "-") {
		return version
	}
	s := strings.Split(version, "-")
	return s[1]
}

const FeatureUpgradeKey = "FEATURE"
const FeatureUpgradeHeight = int64(1)

var govFeatureEnable = &cobra.Command{
	Use:   "enable <fromAddr> <atHeight> <key> <networkID> <fees>",
	Short: "enable a protocol feature",
	Long: `If authorized, enable the protocol feature with the key.
Will prompt the user for the <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		height := args[1]
		key := args[2]
		fstring := fmt.Sprintf("%s:%s", key, height)

		u := govTypes.Upgrade{
			Height:   FeatureUpgradeHeight,
			Version:  FeatureUpgradeKey,
			Features: []string{fstring},
		}
		fees, err := strconv.Atoi(args[4])
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("Enter Password: ")
		res, err := Upgrade(args[0], u, app.Credentials(pwd), args[3], int64(fees), false)
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
