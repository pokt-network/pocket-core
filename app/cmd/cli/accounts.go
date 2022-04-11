package cli

import (
	"bufio"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/pokt-network/pocket-core/app/cmd/rpc"

	"github.com/pokt-network/pocket-core/app"
	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/crypto/keys"
	"github.com/pokt-network/pocket-core/types"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(accountsCmd)
	accountsCmd.AddCommand(createCmd)
	accountsCmd.AddCommand(getValidator)
	accountsCmd.AddCommand(setValidator)
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
	accountsCmd.AddCommand(newMultiPublicKey)
	accountsCmd.AddCommand(signMS)
	accountsCmd.AddCommand(signNexMS)
	accountsCmd.AddCommand(buildMultisig)
	accountsCmd.AddCommand(unsafeDeleteCmd)
}

// accountsCmd represents the accounts namespace command
var accountsCmd = &cobra.Command{
	Use:   "accounts",
	Short: "account management",
	Long: `The accounts namespace handles all account related interactions,
from creating and deleting accounts; to importing and exporting accounts.`,
}

var pwd, oldPwd, decryptPwd, encryptPwd string

func init() {
	buildMultisig.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")
	createCmd.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")
	deleteCmd.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")
	sendTxCmd.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")
	setValidator.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")
	signCmd.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")
	signMS.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")
	signNexMS.Flags().StringVar(&pwd, "pwd", "", "passphrase used by the cmd, non empty usage bypass interactive prompt")

	exportCmd.Flags().StringVar(&decryptPwd, "pwd-decrypt", "", "decrypt passphrase used by the cmd, non empty usage bypass interactive prompt")
	exportCmd.Flags().StringVar(&encryptPwd, "pwd-encrypt", "", "encrypt passphrase used by the cmd, non empty usage bypass interactive prompt")

	exportRawCmd.Flags().StringVar(&decryptPwd, "pwd-decrypt", "", "decrypt passphrase used by the cmd, non empty usage bypass interactive prompt")

	importArmoredCmd.Flags().StringVar(&decryptPwd, "pwd-decrypt", "", "decrypt passphrase used by the cmd, non empty usage bypass interactive prompt")
	importArmoredCmd.Flags().StringVar(&encryptPwd, "pwd-encrypt", "", "encrypt passphrase used by the cmd, non empty usage bypass interactive prompt")

	importCmd.Flags().StringVar(&encryptPwd, "pwd-encrypt", "", "encrypt passphrase used by the cmd, non empty usage bypass interactive prompt")

	updatePassphraseCmd.Flags().StringVar(&pwd, "pwd-new", "", "new passphrase used by the cmd, non empty usage bypass interactive prompt")
	updatePassphraseCmd.Flags().StringVar(&oldPwd, "pwd-old", "", "old passphrase used by the cmd, non empty usage bypass interactive prompt")
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new account",
	Long: `Creates and persists a new account in the Keybase.
Will prompt the user for a passphrase to encrypt the generated keypair.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		kb := keys.New(app.GlobalConfig.PocketConfig.KeybaseName, app.GlobalConfig.PocketConfig.DataDir)
		fmt.Print("Enter Passphrase: \n")
		pass := app.Credentials(pwd)
		fmt.Print("Enter passphrase again: \n")
		confirmedpass := app.Credentials(pwd)
		if pass == confirmedpass {
			kp, err := kb.Create(confirmedpass)
			if err != nil {
				fmt.Printf("Account generation Failed, %s", err)
				return
			}
			fmt.Printf("Account generated successfully:\nAddress: %s\n", kp.GetAddress())
		} else {
			fmt.Println("Account generation Failed, Passphrases do not match")
			return
		}

	},
}

var getValidator = &cobra.Command{
	Use:   "get-validator",
	Short: "Retrieves the main validator from the priv_val file",
	Long:  `Retrieves the main validator from the priv_val file`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Println(app.UninitializedKeybaseError.Error())
			return
		}
		val := app.GetPrivValFile()
		fmt.Printf("Validator Address:%s\n", strings.ToLower(val.Address.String()))
	},
}

var setValidator = &cobra.Command{
	Use:   "set-validator <address>",
	Short: "Sets the main validator account for tendermint",
	Long:  `Sets the main validator account that will be used across all Tendermint functions`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		addr, err := types.AddressFromHex(args[0])
		if err != nil {
			fmt.Printf("Address Error %s", err)
			return
		}
		fmt.Println("Enter the password:")
		app.SetValidator(addr, app.Credentials(pwd))
	},
}

// deleteCmd represents the delete command
var deleteCmd = &cobra.Command{
	Use:   "delete <address>",
	Short: "Remove an account",
	Long: `Deletes a keypair from the keybase.
Will prompt the user for the account passphrase`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Println(app.UninitializedKeybaseError.Error())
			return

		}
		addr, err := types.AddressFromHex(args[0])
		if err != nil {
			fmt.Printf("Address Error %s", err)
			return
		}
		fmt.Print("Enter passphrase: \n")
		err = kb.Delete(addr, app.Credentials(pwd))
		if err != nil {
			fmt.Printf("Error Deleting Account, check your credentials")
			return
		}
		fmt.Println("Account deleted successfully")
	},
}

// unsafeDeleteCmd represents the unsafe delete command (no passphrase)
var unsafeDeleteCmd = &cobra.Command{
	Use:   "unsafe-delete <address>",
	Short: "Remove an account without passphrase",
	Long:  `Deletes a keypair from the keybase without passphrase verification`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Println(app.UninitializedKeybaseError.Error())
			return
		}
		addr, err := types.AddressFromHex(args[0])
		if err != nil {
			fmt.Printf("Address error %s", err)
			return
		}
		fmt.Printf("Are you sure you would like to delete account %s \n", addr.String())
		if !app.Confirmation("") {
			return
		}
		err = kb.UnsafeDelete(addr)
		if err != nil {
			fmt.Printf("Error deleting account: %s", err.Error())
			return
		}
		fmt.Println("Account deleted successfully")
	},
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all accounts",
	Long: `Lists all the account addresses stored in the keybase.
Example output:
	(0) b3746D30F2A579a2efe7F2F6E8E06277a78054C1
	(1) ab514F27e98DE7E3ecE3789b511dA955C3F09Bbc`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Println(app.UninitializedKeybaseError.Error())
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
	Use:   "show <address>",
	Short: "Shows a pubkey for address",
	Long: `Lists an account address and public key.
Example output:
  Address: 		a8cb9e1c0d98fa3a4e1772ada19b8c7f191e61d7
  Public Key: ccc15d61fa80c707cb55ccd80b61720abbac13ca56f7896057e889521462052d`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Println(app.UninitializedKeybaseError.Error())
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
	Long: `Updates the passphrase for the indicated account.
Will prompt the user for the current account passphrase and the new account passphrase.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Println(app.UninitializedKeybaseError.Error())
			return
		}
		addr, err := types.AddressFromHex(args[0])
		if err != nil {
			fmt.Printf("Address Error, %s", err)
			return
		}
		fmt.Println("Enter passphrase: ")
		oldpass := app.Credentials(oldPwd)
		fmt.Println("Enter new passphrase: ")
		newpass := app.Credentials(pwd)
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
	Long: `Digitally signs the specified <msg> using the specified <address> account credentials.
Will prompt the user for the account passphrase.`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Println(app.UninitializedKeybaseError.Error())
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
		fmt.Println("Enter passphrase: ")
		sig, _, err := kb.Sign(addr, app.Credentials(pwd), msg)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Original Message:\t%s\nSignature:\t%s\n", args[1], hex.EncodeToString(sig))
	},
}

var importArmoredCmd = &cobra.Command{
	Use:   "import-armored <armoredJSONFile>",
	Short: "Import keypair using armor",
	Long: `Imports an account using the Encrypted ASCII armored file.
Will prompt the user for a decryption passphrase of the armored ASCII file and for an encryption passphrase to store in the Keybase.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		kb := keys.New(app.GlobalConfig.PocketConfig.KeybaseName, app.GlobalConfig.PocketConfig.DataDir)

		fmt.Println("Enter decrypt pass")
		dPass := app.Credentials(decryptPwd)
		fmt.Println("Enter encrypt pass")
		ePass := app.Credentials(encryptPwd)

		b, err := ioutil.ReadFile(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		kp, err := kb.ImportPrivKey(string(b), dPass, ePass)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("Account imported successfully:\n%s", kp.GetAddress().String())
	},
}
var filePath string

func init() {
	exportCmd.Flags().StringVar(&filePath, "path", "", "the /path/to/export/location where you want to save the file")
}

var exportCmd = &cobra.Command{
	Use:   "export [--path <path>] <address> ",
	Short: "Export an account",
	Long: `Exports the account with <address>, to a file encrypted and ASCII armored in a location specified with --path , if you dont provide a path it will store it on the folder where its running.
Will prompt the user for the account passphrase and for an encryption passphrase for the exported account. Also prompt for an optional hint for the password`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Println(app.UninitializedKeybaseError.Error())
			return
		}
		addr, err := types.AddressFromHex(args[0])
		if err != nil {
			fmt.Printf("Address Error %s", err)
			return
		}
		fmt.Println("Enter Decrypt Passphrase")
		dPass := app.Credentials(decryptPwd)
		fmt.Println("Enter Encrypt Passphrase")
		ePass := app.Credentials(encryptPwd)
		fmt.Println("Enter an optional Hint for remembering the Passphrase")
		reader := bufio.NewReader(os.Stdin)
		hint, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println(err)
			return
		}

		hint = strings.TrimSuffix(hint, "\n")

		pk, err := kb.ExportPrivKeyEncryptedArmor(addr, dPass, ePass, hint)
		if err != nil {
			fmt.Println(err)
			return
		}
		var fname string

		if strings.HasSuffix(filePath, "/") || filePath == "" {
			fname = "pocket-account-" + addr.String() + ".json"

		} else {
			fname = "/pocket-account-" + addr.String() + ".json"
		}

		err = ioutil.WriteFile(filePath+fname, []byte(pk), 0644)
		if err != nil {
			fmt.Println(err)
			fmt.Println("Check --path, if exporting to a folder the folder must exist")
			return
		}

		fmt.Printf("Exported Armor Private Key:\n%s\n", pk)
		fmt.Println("Export Completed")
	},
}

// exportRawCmd represents the exportRaw command
var exportRawCmd = &cobra.Command{
	Use:   "export-raw <address>",
	Short: "Export Plaintext Privkey",
	Long: `Exports the raw private key in hex format.
Will prompt the user for the account passphrase.
NOTE: THIS METHOD IS NOT RECOMMENDED FOR SECURITY REASONS, USE AT YOUR OWN RISK.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		kb := app.MustGetKeybase()
		if kb == nil {
			fmt.Println(app.UninitializedKeybaseError.Error())
			return
		}
		addr, err := types.AddressFromHex(args[0])
		if err != nil {
			fmt.Printf("Address Error %s", err)
			return
		}
		fmt.Println("Enter Decrypt Passphrase")
		dPass := app.Credentials(decryptPwd)
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
	Use:   "send-tx <fromAddr> <toAddr> <amount> <networkID> <fee> <memo>",
	Short: "Send uPOKT",
	Long: `Sends <amount> uPOKT <fromAddr> to <toAddr> with the specified <memo>.
Prompts the user for <fromAddr> account passphrase.`,
	Args: cobra.ExactArgs(6),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		amount, err := strconv.Atoi(args[2])
		if err != nil {
			fmt.Println(err)
			return
		}
		fees, err := strconv.Atoi(args[4])
		if err != nil {
			fmt.Println(err)
			return
		}
		memo := args[5]
		fmt.Printf("Adding Memo: %v\n", memo)
		fmt.Println("Enter passphrase: ")
		res, err := SendTransaction(args[0], args[1], app.Credentials(pwd), args[3], types.NewInt(int64(amount)), int64(fees), memo, false)
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

// sendRawTxCmd represents the sendTx command
var sendRawTxCmd = &cobra.Command{
	Use:   "send-raw-tx <fromAddr> <txBytes>",
	Short: "Send raw transaction from signed bytes",
	Long:  `Sends presigned transaction through the tendermint node`,
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		bz, err := hex.DecodeString(args[1])
		if err != nil {
			fmt.Println(err)
			return
		}
		p := rpc.SendRawTxParams{
			Addr:        args[0],
			RawHexBytes: hex.EncodeToString(bz),
		}
		j, err := json.Marshal(p)
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

// importCmd represents the import command
var importCmd = &cobra.Command{
	Use:   "import-raw <private-key-hex>",
	Short: "import-raw <private-key-hex>",
	Long: `Imports an account using the provided <private-key-hex>
Will prompt the user for a passphrase to encrypt the generated keypair.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		pkBytes, err := hex.DecodeString(args[0])
		if err != nil {
			fmt.Println(err)
			return
		}
		kb := keys.New(app.GlobalConfig.PocketConfig.KeybaseName, app.GlobalConfig.PocketConfig.DataDir)
		fmt.Println("Enter Encrypt Passphrase")
		ePass := app.Credentials(encryptPwd)
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

var newMultiPublicKey = &cobra.Command{
	Use:   "create-multi-public <ordered-comma-separated-hex-pubkeys>",
	Short: "create a multisig public key",
	Long:  `create a multisig public key with a comma separated list of hex encoded public keys`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		rawPKs := strings.Split(strings.TrimSpace(args[0]), ",")
		var pks []crypto.PublicKey
		for _, pk := range rawPKs {
			p, err := crypto.NewPublicKey(pk)
			if err != nil {
				fmt.Println(fmt.Errorf("error in public key creation: %v", err))
				return
			}
			pks = append(pks, p)
		}
		multiSigPubKey := crypto.PublicKeyMultiSignature{PublicKeys: pks}
		fmt.Printf("Sucessfully generated Multisig Public Key:\n%s\nWith Address:\n%s\n", multiSigPubKey.String(), multiSigPubKey.Address())
	},
}

var buildMultisig = &cobra.Command{
	Use:   "build-MS-Tx <signer-address> <json-message> <ordered-comma-separated-hex-pubkeys> <networkID> <fees>",
	Short: "Build and sign a multisig tx",
	Args:  cobra.ExactArgs(5),
	Long:  `Build and sign a multisignature transaction from scratch: result is hex encoded std tx object.`,
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		msg := args[1]
		rawPKs := strings.Split(strings.TrimSpace(args[2]), ",")
		var pks []crypto.PublicKey
		for _, pk := range rawPKs {
			p, err := crypto.NewPublicKey(pk)
			if err != nil {
				fmt.Println(fmt.Errorf("error creating the public key: %v", err))
				continue
			}
			pks = append(pks, p)
		}

		multiSigPubKey := crypto.PublicKeyMultiSignature{PublicKeys: pks}
		fmt.Println("Enter passphrase: ")
		fees, err := strconv.Atoi(args[4])
		if err != nil {
			fmt.Println(err)
			return
		}
		bz, err := app.BuildMultisig(args[0], msg, app.Credentials(pwd), args[3], multiSigPubKey, int64(fees), false)
		if err != nil {
			fmt.Println(fmt.Errorf("error building the multisig: %v", err))
		}
		fmt.Println("Multisig transaction: \n" + hex.EncodeToString(bz))
	},
}

var signMS = &cobra.Command{
	Use:   "sign-ms-tx <signer-address> <hex-amino-stdtx> <hex-pubkeys> <networkID> ",
	Short: "sign a multisig tx",
	Long:  `sign a multisignature transaction using public keys, and the transaciton object, result is hex encoded std tx object`,
	Args:  cobra.ExactArgs(4),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		msg := args[1]
		rawPKs := strings.Split(strings.TrimSpace(args[2]), ",")
		var pks []crypto.PublicKey
		for _, pk := range rawPKs {
			p, err := crypto.NewPublicKey(pk)
			if err != nil {
				fmt.Println(fmt.Errorf("error generating public key: %v", err))
				continue
			}
			pks = append(pks, p)
		}
		fmt.Println("Enter passphrase: ")
		bz, err := app.SignMultisigOutOfOrder(args[0], msg, app.Credentials(pwd), args[3], pks, false)
		if err != nil {
			fmt.Println(fmt.Errorf("error signing multisig: %v", err))
		}
		fmt.Println("Multisig transaction: \n" + hex.EncodeToString(bz))
	},
}

var signNexMS = &cobra.Command{
	Use:   "sign-ms-next <signer-address> <hex-stdtx> <networkID> ",
	Short: "Sign a multisig tx",
	Long: `Sign a multisignature transaction using the transaciton object, result is hex encoded std tx object
NOTE: you MUST be the next signer (in order of public keys in the ms public key object) or the signature will be invalid.`,
	Args: cobra.ExactArgs(3),
	Run: func(cmd *cobra.Command, args []string) {
		app.InitConfig(datadir, tmNode, persistentPeers, seeds, remoteCLIURL)
		msg := args[1]
		fmt.Println("Enter password: ")
		bz, err := app.SignMultisigNext(args[0], msg, app.Credentials(pwd), args[2], false)
		if err != nil {
			fmt.Println(fmt.Errorf("error signing the multisig: %v", err))
		}
		fmt.Println("Multisig transaction: \n" + hex.EncodeToString(bz))
	},
}
