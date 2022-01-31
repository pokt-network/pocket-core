package cli

import (
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/types"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(nodesCmd)
	nodesCmd.AddCommand(nodeStakeCmd)
	nodesCmd.AddCommand(nodeUnstakeCmd)
	nodesCmd.AddCommand(nodeUnjailCmd)
}

var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "node management",
	Long: `The node namespace handles all node related interactions,
from staking and unstaking; to unjailing.`,
}

func init() {
	nodeStakeCmd.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")
	nodeUnstakeCmd.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")
	nodeUnjailCmd.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")
}

var nodeStakeCmd = &cobra.Command{
	Use:   "stake <operatorAddress||signerAddress> <amount> <RelayChainIDs> <serviceURI> <outputAddress||signerAddress> <networkID> <fee> <isBefore8.0>",
	Short: "Stake a node in the network, the signer may be the operator or the output address. The signer must specify the public key of the output or operator",
	Long: `Stake the node into the network, making it available for service.
Will prompt the user for the <signerAddress> account passphrase. After the 0.6.X upgrade, if the node is already staked, this transaction acts as an *update* transaction.
A node can updated relayChainIDs, serviceURI, and raise the stake amount with this transaction.
If the node is currently staked at X and you submit an update with new stake Y. Only Y-X will be subtracted from an account
If no changes are desired for the parameter, just enter the current param value just as before.
The signer may be the operator or the output address.`,
	Args: cobra.ExactArgs(8),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		operator := args[0]
		output := args[4]
		amount, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		am := types.NewInt(int64(amount))
		if am.LTE(types.NewInt(15100000000)) {
			fmt.Println("The amount you are staking for is below the recommendation of 15100 POKT, would you still like to continue? y|n")
			if !app.Confirmation() {
				return
			}
		}
		reg, err := regexp.Compile("[^,a-fA-F0-9]+")
		if err != nil {
			log.Fatal(err)
		}
		rawChains := reg.ReplaceAllString(args[2], "")
		chains := strings.Split(rawChains, ",")
		serviceURI := args[3]
		fee, err := strconv.Atoi(args[6])
		if err != nil {
			fmt.Println(err)
			return
		}
		isBefore8, err := strconv.ParseBool(args[7])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Enter Passphrase: ")
		res, err := StakeNode(chains, serviceURI, operator, output, app.Credentials(pwd), args[5], types.NewInt(int64(amount)), int64(fee), isBefore8)
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

var nodeUnstakeCmd = &cobra.Command{
	Use:   "unstake <operatorAddr> <fromAddr> <networkID> <fee> <isBefore8.0>",
	Short: "Unstake a node in the network",
	Long: `Unstake a node from the network, changing it's status to Unstaking.
Will prompt the user for the <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(5),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		fee, err := strconv.Atoi(args[3])
		if err != nil {
			fmt.Println(err)
			return
		}
		isBefore8, err := strconv.ParseBool(args[4])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Enter Password: ")
		res, err := UnstakeNode(args[0], args[1], app.Credentials(pwd), args[2], int64(fee), isBefore8)
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

var nodeUnjailCmd = &cobra.Command{
	Use:   "unjail <operatorAddr> <fromAddr> <networkID> <fee> <isBefore8.0>",
	Short: "Unjails a node in the network",
	Long: `Unjails a node from the network, allowing it to participate in service and consensus again.
Will prompt the user for the <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		fee, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		isBefore8, err := strconv.ParseBool(args[3])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Enter Password: ")
		res, err := UnjailNode(args[0], args[1], app.Credentials(pwd), args[1], int64(fee), isBefore8)
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
