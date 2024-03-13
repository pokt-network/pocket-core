package cli

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/spf13/cobra"

	"github.com/pokt-network/pocket-core/app"
)

func init() {
	rootCmd.AddCommand(nodesCmd)
	nodesCmd.AddCommand(nodeUnstakeCmd)
	nodesCmd.AddCommand(nodeUnjailCmd)
	nodesCmd.AddCommand(stakeNewCmd)
}

var nodesCmd = &cobra.Command{
	Use:   "nodes",
	Short: "node management",
	Long: `The node namespace handles all node related interactions, from staking and unstaking; to unjailing.

---

Operator Address (i.e. Non-Custodial Address) can do the following:
- Submit Block, Claim & Proof Txs

Output Address (i.e. Custodial Address) can do the following:
- Receive earned rewards
- Receive funds after unstaking

Both Operator and Output Addresses can do the following:
- Submit Stake, EditStake, Unstake, Unjail Txs

---
`,
}

func init() {
	nodeUnstakeCmd.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")
	nodeUnjailCmd.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")
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
		res, err := UnjailNode(args[0], args[1], app.Credentials(pwd), args[2], int64(fee), isBefore8)
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

// stakeNewCmd is an upgraded version of `nodesCmd` that captures newer
// on-chain functionality in a cleaner way
var stakeNewCmd = &cobra.Command{
	Use:   "stakeNew <OperatorPublicKey> <OutputAddress> <SignerAddress> <Stake> <ChainIDs> <ServiceURL> <RewardDelegators> <NetworkID> <Fee> [Memo]",
	Short: "Stake a node in the network",
	Long: `Stake a node in the network, promoting it to a servicer or a validator.

The command takes the following parameters.

  OperatorPublicKey Public key to use as the node's operator account
  OutputAddress     Address to use as the node's output account
  SignerAddress     Address to sign the transaction
  Stake             Amount to stake in uPOKT
  ChainIDs          Comma-separated chain IDs to host on the node
  ServiceURL        Relay endpoint of the node.  Must include the port number.
  RewardDelegators  Addresses to share rewards
  NetworkID         Network ID to submit a transaction to e.g. mainnet or testnet
  Fee               Transaction fee in uPOKT
  Memo              Optional. Text to include in the transaction.  No functional effect.

Example:
$ pocket nodes stakeNew \
    e237efc54a93ed61689959e9afa0d4bd49fa11c0b946c35e6bebaccb052ce3fc \
    fe818527cd743866c1db6bdeb18731d04891df78 \
    1164b9c95638fc201f35eca2af4c35fe0a81b6cf \
    8000000000000 \
    DEAD,BEEF \
    https://x.com:443 \
    '{"1000000000000000000000000000000000000000":1,"2000000000000000000000000000000000000000":2}' \
    mainnet \
    10000 \
    "new stake with delegators!"
`,
	Args: cobra.MinimumNArgs(9),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)

		operatorPubKey := args[0]
		outputAddr := args[1]
		signerAddr := args[2]
		stakeAmount := args[3]
		chains := args[4]
		serviceUrl := args[5]
		delegators := args[6]
		networkId := args[7]
		fee := args[8]
		memo := ""
		if len(args) >= 10 {
			memo = args[9]
		}

		fmt.Println("Enter Passphrase:")
		passphrase := app.Credentials(pwd)

		rawStakeTx, err := BuildStakeTx(
			operatorPubKey,
			outputAddr,
			stakeAmount,
			chains,
			serviceUrl,
			delegators,
			networkId,
			fee,
			memo,
			signerAddr,
			passphrase,
		)
		if err != nil {
			fmt.Println(err)
			return
		}
		txBytes, err := json.Marshal(rawStakeTx)
		if err != nil {
			fmt.Println("Fail to build a transaction:", err)
			return
		}
		resp, err := QueryRPC(SendRawTxPath, txBytes)
		if err != nil {
			fmt.Println("Fail to submit a transaction:", err)
			return
		}
		fmt.Println(resp)
	},
}
