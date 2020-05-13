package cli

import (
	"encoding/json"
	"fmt"
	"log"
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
	Use:   "transfer <action (dao_burn or dao_transfer)> <amount> <fromAddr> <toAddr> <chainID>",
	Short: "Transfer from DAO",
	Long: `If authorized, move funds from the DAO.
Actions: [burn, transfer]`,
	Args: cobra.ExactArgs(4),
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
		res, err := DAOTx(fromAddr, toAddr, pass, types.NewInt(int64(amount)), action, args[3])
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
	Use:   "change_param <fromAddr> <chainID> <paramKey module/param> <paramValue (jsonObj)>",
	Short: "Edit a param in the network",
	Long: `If authorized, submit a tx to change any param from any module.
Will prompt the user for the <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		fmt.Println("Enter Password: ")
		var i interface{}
		err := json.Unmarshal([]byte(args[3]), &i)
		if err != nil {
			log.Fatal(err)
		}
		res, err := ChangeParam(args[0], args[2], i, app.Credentials(), args[1])
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
	Use:   "upgrade <fromAddr> <atHeight> <version>, <chainID>",
	Short: "Upgrade the protocol",
	Long: `If authorized, upgrade the protocol.
Will prompt the user for the <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, tmRPCPort, tmPeersPort)
		i, err := strconv.Atoi(args[1])
		if err != nil {
			log.Fatal(err)
		}
		u := govTypes.Upgrade{
			Height:  int64(i),
			Version: args[2],
		}
		fmt.Println("Enter Password: ")
		res, err := Upgrade(args[0], u, app.Credentials(), args[3])
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
