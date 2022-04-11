package cli

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/types"
	"github.com/spf13/cobra"
	"log"
	"regexp"
	"strconv"
	"strings"
)

func init() {
	nodesCmd.AddCommand(nodeStakeCmd)
	nodeStakeCmd.AddCommand(custodialStakeCmd)
	nodeStakeCmd.AddCommand(nonCustodialstakeCmd)

	custodialStakeCmd.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")
	nonCustodialstakeCmd.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")

}

var nodeStakeCmd = &cobra.Command{
	Use:   "stake",
	Short: "Stake a node in the network",
	Long:  "Stake the node into the network, making it available for service.",
}

var custodialStakeCmd = &cobra.Command{
	Use:   "custodial <fromAddr> <amount> <RelayChainIDs> <serviceURI> <networkID> <fee> <isBefore8.0>",
	Short: "Stake a node in the network. Custodial stake uses the same address as operator/output for rewards/return of staked funds.",
	Long: `Stake the node into the network, making it available for service.
Will prompt the user for the <fromAddr> account passphrase. If the node is already staked, this transaction acts as an *update* transaction.
A node can updated relayChainIDs, serviceURI, and raise the stake amount with this transaction.
If the node is currently staked at X and you submit an update with new stake Y. Only Y-X will be subtracted from an account
If no changes are desired for the parameter, just enter the current param value just as before`,
	Args: cobra.ExactArgs(7),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		fromAddr := args[0]
		amount, err := strconv.Atoi(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		am := types.NewInt(int64(amount))
		if am.LTE(types.NewInt(15100000000)) {
			fmt.Println("The amount you are staking for is below the recommendation of 15100 POKT, would you still like to continue? y|n")
			if !app.Confirmation(pwd) {
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
		fee, err := strconv.Atoi(args[5])
		if err != nil {
			fmt.Println(err)
			return
		}
		isBefore8, err := strconv.ParseBool(args[6])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Enter Passphrase: ")
		res, err := LegacyStakeNode(chains, serviceURI, fromAddr, app.Credentials(pwd), args[4], types.NewInt(int64(amount)), int64(fee), isBefore8)
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

var nonCustodialstakeCmd = &cobra.Command{
	Use:   "non-custodial <operatorPublicKey> <outputAddress> <amount> <RelayChainIDs> <serviceURI> <networkID> <fee> <isBefore8.0>",
	Short: "Stake a node in the network, non-custodial stake allows a different output address for rewards/return of staked funds. The signer may be the operator or the output address. The signer must specify the public key of the operator",
	Long: `Stake the node into the network, making it available for service.
Will prompt the user for the signer account passphrase, fund and fees are collected from signer account. If both accounts are present signer priority is first output then operator. If the node is already staked, this transaction acts as an *update* transaction.
A node can updated relayChainIDs, serviceURI, and raise the stake amount with this transaction.
If the node is currently staked at X and you submit an update with new stake Y. Only Y-X will be subtracted from an account
If no changes are desired for the parameter, just enter the current param value just as before.
The signer may be the operator or the output address.`,
	Args: cobra.ExactArgs(8),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		operatorPubKey := args[0]
		output := args[1]
		amount, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		am := types.NewInt(int64(amount))
		if am.LTE(types.NewInt(15100000000)) {
			fmt.Println("The amount you are staking for is below the recommendation of 15100 POKT, would you still like to continue? y|n")
			if !app.Confirmation("") {
				return
			}
		}
		reg, err := regexp.Compile("[^,a-fA-F0-9]+")
		if err != nil {
			log.Fatal(err)
		}
		rawChains := reg.ReplaceAllString(args[3], "")
		chains := strings.Split(rawChains, ",")
		serviceURI := args[4]
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
		res, err := StakeNode(chains, serviceURI, operatorPubKey, output, app.Credentials(pwd), args[5], types.NewInt(int64(amount)), int64(fee), isBefore8)
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
