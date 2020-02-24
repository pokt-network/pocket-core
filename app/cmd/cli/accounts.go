package cli

import (
	"encoding/hex"
	"fmt"
	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/posmint/crypto"
	"github.com/pokt-network/posmint/crypto/keys"
	"github.com/pokt-network/posmint/types"
	"github.com/spf13/cobra"
	"path/filepath"
	"strconv"
)

func init() {
	rootCmd.AddCommand(accountsCmd)
	accountsCmd.AddCommand(createCmd)
	accountsCmd.AddCommand(getCoinbase)
	accountsCmd.AddCommand(deleteCmd)
	accountsCmd.AddCommand(listCmd)
	accountsCmd.AddCommand(showCmd)
	accountsCmd.AddCommand(updatePassphraseCmd)
	accountsCmd.AddCommand(signCmd)
	accountsCmd.AddCommand(importArmoredCmd)
	accountsCmd.AddCommand(importCmd)
	accountsCmd.AddCommand(exportCmd)
	accountsCmd.AddCommand(exportRawCmd)
	accountsCmd.AddCommand(sendTxCmd)
	accountsCmd.AddCommand(sendRawTxCmd)
}

// accountsCmd represents the accounts namespace command
var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "account management",
	Long: `The accounts namespace handles all account related interactions,
from creating and deleting accounts, to importing and exporting accounts.`,
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new account",
	Long:  `Creates and persists a new account in the Keybase. Will prompt the user for a passphrase to encrypt the generated keypair.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		kb := keys.New(app.KeybaseName, app.InitDataDirectory(datadir)+string(filepath.Separator)+app.KBDirectoryName)
		fmt.Print("Enter Passphrase: \n")
		kp, err := kb.Create(app.Credentials())
		if err != nil {
			fmt.Printf("Account generation Failed, %s", err)
			return
		}
		fmt.Printf("Account generated successfully:\nAddress: %s\n", kp.GetAddress())
	},
}

var getCoinbase = &cobra.Command{
	Use:   "get-coinbase",
	Short: "Gets the coinbase account from the keybase",
	Long:  `Retrieves the coinbase account from the pocket core keybase`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Printf(app.UninitializedKeybaseError.Error())
			return
		}
		coinbase, err := kb.GetCoinbase()
		if err != nil {
			fmt.Printf("Could not get coinbase, %s", err)
			return
		}
		fmt.Printf("Coinbase Account:%s\n", coinbase.GetAddress())
	},
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete  <address>",
	Short: "DeleteInvoice an account",
	Long:  `Deletes a keypair from the keybase. Will prompt the user for the account passphrase`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Printf(app.UninitializedKeybaseError.Error())
			return

		}
		addr, err := types.AddressFromHex(args[0])
		if err != nil {
			fmt.Printf("Address Error %s", err)
			return
		}
		fmt.Print("Enter Password: \n")
		err = kb.Delete(addr, app.Credentials())
		if err != nil {
			fmt.Printf("Error Deleting Account, check your credentials")
			return
		}
		fmt.Println("Account deleted successfully")
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long: `Lists all the account addresses stored in the keybase. Example output:
(0) b3746D30F2A579a2efe7F2F6E8E06277a78054C1
(1) ab514F27e98DE7E3ecE3789b511dA955C3F09Bbc`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Printf(app.UninitializedKeybaseError.Error())
			return
		}
		kp, err := kb.List()
		if err != nil {
			fmt.Printf("Error retrieving accounts from keybase, %s", err)
			return
		}
		for i, key := range kp {
			fmt.Printf("(%d) %s\n", i, key.GetAddress().String())
		}
	},
}

// showCmd represents the show command
var showCmd = &cobra.Command{
	Use:   "show  <address>",
	Short: "Shows a pubkey for address",
	Long:  `Lists an account address and public key`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Printf(app.UninitializedKeybaseError.Error())
			return
		}
		addr, err := types.AddressFromHex(args[0])
		if err != nil {
			fmt.Printf("Address Error, %s", err)
			return
		}
		kp, err := kb.Get(addr)
		if err != nil {
			fmt.Printf("Error Getting pubkey For Address, %s", err)
			return
		}
		fmt.Printf("Address:\t%s\nPublic Key:\t%s\n",
			kp.GetAddress().String(),
			hex.EncodeToString(kp.PublicKey.RawBytes()))
	},
}

// updatePassphraseCmd represents the updatePassphrase command
var updatePassphraseCmd = &cobra.Command{
	Use:   "update-passphrase <address>",
	Short: "Update account passphrase",
	Long:  `Updates the passphrase for the indicated account. Will prompt the user for the current account passphrase and the new account passphrase.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Printf(app.UninitializedKeybaseError.Error())
			return
		}
		addr, err := types.AddressFromHex(args[0])
		if err != nil {
			fmt.Printf("Address Error, %s", err)
			return
		}
		fmt.Println("Enter Password: ")
		oldpass := app.Credentials()
		fmt.Println("Enter New Password: ")
		newpass := app.Credentials()
		err = kb.Update(addr, oldpass, newpass)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Successfully updated account: " + addr.String())
	},
}

// signCmd represents the sign command
var signCmd = &cobra.Command{
	Use:   "sign <address> <msg>",
	Short: "Sign a message with an account",
	Long:  `Digitally signs the specified <msg> using the specified <address> account credentials. Will prompt the user for the account passphrase.`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Printf(app.UninitializedKeybaseError.Error())
			return
		}
		addr, err := types.AddressFromHex(args[0])
		if err != nil {
			fmt.Printf("Address Error %s", err)
			return
		}
		msg, err := hex.DecodeString(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Enter Password: ")
		sig, _, err := kb.Sign(addr, app.Credentials(), msg)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Original Message:\t%s\nSignature:\t%s\n", args[1], hex.EncodeToString(sig))
	},
}

var importArmoredCmd = &cobra.Command{
	Use:   "import-armored <armor>",
	Short: "Import keypair using armor",
	Long:  `Imports an account using the Encrypted ASCII armored <armor> string. Will prompt the user for a decryption passphrase of the <armor> string and for an encryption passphrase to store in the Keybase.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Printf(app.UninitializedKeybaseError.Error())
			return
		}
		fmt.Println("Enter decrypt pass")
		dPass := app.Credentials()
		fmt.Println("Enter encrypt pass")
		ePass := app.Credentials()
		kp, err := kb.ImportPrivKey(args[0], dPass, ePass)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Account imported successfully:\n%s", kp.GetAddress().String())
	},
}

var exportCmd = &cobra.Command{
	Use:   "export <address>",
	Short: "Export an account",
	Args:  cobra.ExactArgs(1),
	Long: `Exports the account with <address>, encrypted and ASCII armored.
Will prompt the user for the account passphrase and for an encryption passphrase for the exported account.
`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Printf(app.UninitializedKeybaseError.Error())
			return
		}
		addr, err := types.AddressFromHex(args[0])
		if err != nil {
			fmt.Printf("Address Error %s", err)
			return
		}
		fmt.Println("Enter Decrypt Passphrase")
		dPass := app.Credentials()
		fmt.Println("Enter Encrypt Passphrase")
		ePass := app.Credentials()
		pk, err := kb.ExportPrivKeyEncryptedArmor(addr, dPass, ePass)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Exported Armor Private Key:\n%s\n", pk)
	},
}

// exportRawCmd represents the exportRaw command
var exportRawCmd = &cobra.Command{
	Use:   "export-raw <address>",
	Short: "Export Plaintext Privkey",
	Args:  cobra.ExactArgs(1),
	Long: `Exports the raw private key in hex format. Will prompt the user for the account passphrase.
NOTE: THIS METHOD IS NOT RECOMMENDED FOR SECURITY REASONS, USE AT YOUR OWN RISK.*`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Printf(app.UninitializedKeybaseError.Error())
			return
		}
		addr, err := types.AddressFromHex(args[0])
		if err != nil {
			fmt.Printf("Address Error %s", err)
			return
		}
		fmt.Println("Enter Decrypt Passphrase")
		dPass := app.Credentials()
		pk, err := kb.ExportPrivateKeyObject(addr, dPass)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Exported Raw Private Key:\n%s\n", hex.EncodeToString(pk.RawBytes()))
	},
}

// sendTxCmd represents the sendTx command
var sendTxCmd = &cobra.Command{
	Use:   "send-tx <fromAddr> <toAddr> <amount>",
	Short: "Send POKT",
	Long:  `Sends <amount> POKT <fromAddr> to <toAddr>. Prompts the user for <fromAddr> account passphrase.`,
	Args:  cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		amount, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Println("Enter Password: ")
		res, err := app.SendTransaction(args[0], args[1], app.Credentials(), types.NewInt(int64(amount)))
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Transaction Submitted: %s", res.TxHash)
	},
}

// sendRawTxCmd represents the sendTx command
var sendRawTxCmd = &cobra.Command{
	Use:   "send-raw-tx <fromAddr> <txBytes>",
	Short: "Send raw transaction from signed bytes",
	Long:  `Sends presigned transaction through the tendermint node`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		bz, err := hex.DecodeString(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		res, err := app.SendRawTx(args[0], bz)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Transaction Submitted: %s", res.TxHash)
	},
}

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import-raw <private-key-hex>",
	Short: "import-raw <private-key-hex>",
	Args:  cobra.ExactArgs(1),
	Long: `Imports an account using the provided <private-key-hex>

Will prompt the user for a passphrase to encrypt the generated keypair.
`,
	Run: func(cmd *cobra.Command, args []string) {
		app.SetTMNode(tmNode)
		pkBytes, err := hex.DecodeString(args[0])
		kb := keys.New(app.KeybaseName, app.InitDataDirectory(datadir)+string(filepath.Separator)+app.KBDirectoryName)
		fmt.Println("Enter Encrypt Passphrase")
		ePass := app.Credentials()
		var pk [crypto.Ed25519PrivKeySize]byte
		copy(pk[:], pkBytes)
		kp, err := kb.ImportPrivateKeyObject(pk, ePass)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("Account imported successfully:\n%s\n", kp.GetAddress().String())
	},
}
