package cli

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/posmint/types"
	govTypes "github.com/pokt-network/posmint/x/gov/types"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(govCmd)
	govCmd.AddCommand(govDAOTransfer)
	govCmd.AddCommand(govChangeParam)
	govCmd.AddCommand(govUpgrade)
}

var govCmd = &cobra.Command{
	Use:   "gov",
	Short: "governance management",
	Long: `The gov namespace handles all governance related interactions,
from DAOTransfer, change parameters; to performing protocol Upgrades. `,
}

var govDAOTransfer = &cobra.Command{
	Use:   "transfer <action (dao_burn or dao_transfer)> <amount> <fromAddr> <toAddr> ",
	Short: "Transfer from DAO",
	Long: `If authorized, move funds from the DAO.
Actions: [burn, transfer]`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		var toAddr string
		if len(args) == 4 {
			toAddr = args[3]
		}
		fromAddr := args[2]
		action := args[0]
		amount, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Enter Password: ")
		pass := app.Credentials()
		res, err := app.DAOTx(fromAddr, toAddr, pass, types.NewInt(int64(amount)), action)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Transaction Submitted: %s\n", res.TxHash)
	},
}

var govChangeParam = &cobra.Command{
	Use:   "change_param <fromAddr> <paramKey module/param> <paramValue (jsonObj)>",
	Short: "Edit a param in the network",
	Long: `If authorized, submit a tx to change any param from any module.
Will prompt the user for the <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		fmt.Println("Enter Password: ")
		var i interface{}
		err := json.Unmarshal([]byte(args[2]), &i)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		res, err := app.ChangeParam(args[0], args[1], i, app.Credentials())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Transaction Submitted: %s\n", res.TxHash)
	},
}

var govUpgrade = &cobra.Command{
	Use:   "upgrade <fromAddr> <atHeight> <version>",
	Short: "Upgrade the protocol",
	Long: `If authorized, upgrade the protocol.
Will prompt the user for the <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		i, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		u := govTypes.Upgrade{
			Height:  int64(i),
			Version: args[2],
		}
		fmt.Println("Enter Password: ")
		res, err := app.Upgrade(args[0], u, app.Credentials())
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Transaction Submitted: %s\n", res.TxHash)
	},
}
