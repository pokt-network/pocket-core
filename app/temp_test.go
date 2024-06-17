package app

import (
	// "encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth"
	authTypes "github.com/pokt-network/pocket-core/x/auth/types"
	nodeTypes "github.com/pokt-network/pocket-core/x/nodes/types"
)

func TestTiger_Decoder(t *testing.T) {
	rawHexBytes := "8e08aced000573720034636f6d2e62697468756d622e7369676e65722e636f696e732e706f6b74626173652e506f6b745374645472616e73616374696f6ecfda1d1f580a49100200054c0007656e74726f70797400104c6a6176612f6c616e672f4c6f6e673b4c00036665657400104c6a6176612f7574696c2f4c6973743b4c00046d656d6f7400124c6a6176612f6c616e672f537472696e673b4c00036d73677400324c636f6d2f62697468756d622f7369676e65722f636f696e732f706f6b74626173652f506f6b745374644d6573736167653b4c0005747853696774003c4c636f6d2f62697468756d622f7369676e65722f636f696e732f706f6b74626173652f506f6b745374645472616e73616374696f6e2454785369673b78707372000e6a6176612e6c616e672e4c6f6e673b8be490cc8f23df0200014a000576616c7565787200106a6176612e6c616e672e4e756d62657286ac951d0b94e08b02000078700a65e77a4e5bd8007372001a6a6176612e7574696c2e4172726179732441727261794c697374d9a43cbecd8806d20200015b0001617400135b4c6a6176612f6c616e672f4f626a6563743b7870757200385b4c636f6d2e62697468756d622e7369676e65722e636f696e732e706f6b74626173652e506f6b745472616e73616374696f6e244665653b7a1f658018bb6da302000078700000000173720035636f6d2e62697468756d622e7369676e65722e636f696e732e706f6b74626173652e506f6b745472616e73616374696f6e24466565ab35dab790bbf7cd0200024c0006616d6f756e7471007e00034c000564656e6f6d71007e00037870740005313030303074000575706f6b7474000073720030636f6d2e62697468756d622e7369676e65722e636f696e732e706f6b74626173652e506f6b745374644d657373616765943f72b5c3d750000200024c00077479706555726c71007e00034c000576616c756571007e000378707400102f782e6e6f6465732e4d736753656e6474003c70776e7a624869304a477a366f4f7a445938313853526455544172323475732f746e786a44576c6837347a37644e6a536774506f717a55774d4441777372003a636f6d2e62697468756d622e7369676e65722e636f696e732e706f6b74626173652e506f6b745374645472616e73616374696f6e245478536967a5d184995a47a92e0200025b00097075626c69634b65797400025b425b00097369676e617475726571007e00197870757200025b42acf317f8060854e002000078700000002003bece7df964f73215601cd844704f515ff9ab080d4f425ac6d05afd50af31277571007e001b00000040e296359b4a6350d1b460df5c45ae604dde8f8ac871e499b4e06b581660e1a60d73314c8ce9943802b9405e61a2b0c00c8ee511d447cccd1e9ee777fd1d2ff09d"
	cdc := memCodec()
	decoder := auth.DefaultTxDecoder(cdc)
	txBytes, err := hex.DecodeString(rawHexBytes)
	if err != nil {
		t.Fatal(err)
	}
	tx, err := decoder(txBytes, 12000000)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(tx)
}

func TestTiger_Encode(t *testing.T) {
	cdc := memCodec() // app.Codec()
	encoder := auth.DefaultTxEncoder(cdc)
	decoder := auth.DefaultTxDecoder(cdc)
	// privKey, err := crypto.NewPrivateKey("5d86a93dee1ef5f950ccfaafd09d9c812f790c3b2c07945501f68b339118aca0e237efc54a93ed61689959e9afa0d4bd49fa11c0b946c35e6bebaccb052ce3fc")
	privKey, err := crypto.NewPrivateKey("3505756aeeaa33364451f0681c44631db9d4bf5b6acc8d571b7d76005057ebf603bece7df964f73215601cd844704f515ff9ab080d4f425ac6d05afd50af3127")
	if err != nil {
		t.Fatal(err)
	}
	fromAddr := sdk.Address(privKey.PubKey().Address())
	toAddr, err := types.AddressFromHex("f6e2eb3fb67c630d6961ef8cfb74d8d282d3e8ab")
	if err != nil {
		t.Fatal(err)
	}
	msg := &nodeTypes.MsgSend{
		Amount:      types.NewInt(50000),
		FromAddress: fromAddr,
		ToAddress:   toAddr,
	}
	builder := authTypes.NewTxBuilder(
		encoder,
		decoder,
		"mainnet",
		"", // empty memo
		types.NewCoins(types.NewCoin(types.DefaultStakeDenom, types.NewInt(10000))),
	)
	entropy := int64(749259425513723904)
	txBytes, err := builder.BuildAndSignWithEntropyForTesting(privKey, msg, entropy)
	if err != nil {
		t.Fatal(err)
	}
	// hexString := base64.StdEncoding.EncodeToString(txBytes)
	// hexString := hex.EncodeToString(txBytes)
	// fmt.Println("raw_hex_bytes", hexString)

	// txDecoded, err := decoder(txBytes, 900000)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Printf("%v", txDecoded)

	stdTx := auth.StdTx{}
	err = cdc.UnmarshalBinaryLengthPrefixed(txBytes, &stdTx, 90000)
	if err != nil {
		t.Fatal(err)
	}
	// fmt.Printf("%+v\n", stdTx)
	// fmt.Println("stdTx.Msgs: ", stdTx.Msg.Type(), stdTx.Msg.GetSignBytes())
	// fmt.Println("stdTx.Signature: ", stdTx.Signature)
	// fmt.Println("hex.hex.EncodeToString(stdTx.Signature.PublicKey): ", hex.EncodeToString(stdTx.Signature.PublicKey.Bytes()))
	// fmt.Println("hex.EncodeToString(stdTx.Signature.Signature): ", hex.EncodeToString(stdTx.Signature.Signature))

	// jsonData, err := json.Marshal(stdTx)
	jsonData, err := json.MarshalIndent(stdTx, "", "    ")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("json", string(jsonData))
	fmt.Println(stdTx, stdTx.GetMsg().Type())

	// protoStdTx := auth.ProtoStdTx{}
	// err = cdc.UnmarshalBinaryLengthPrefixed(txBytes, &protoStdTx, 90000)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// fmt.Println("protoStdTx.Signature", protoStdTx.Signature)
}

// Refernce:
// PoktTransaction = {"chain_id":"mainnet","entropy":749259425513723904,"fee":[{"amount":"10000","denom":"upokt"}],"memo":"","msg":{"type":"pos/Send","value":{"amount":"50000","from_address":"a709f36c78b4246cfaa0ecc363cd7c4917544c0a","to_address":"f6e2eb3fb67c630d6961ef8cfb74d8d282d3e8ab"}}}
// PoktStdTransaction = {"msg":{"typeUrl":"/x.nodes.MsgSend","value":"pwnzbHi0JGz6oOzDY818SRdUTAr24us/tnxjDWlh74z7dNjSgtPoqzUwMDAw"},"fee":[{"amount":"10000","denom":"upokt"}],"txSig":{"publicKey":"A77Offlk9zIVYBzYRHBPUV/5qwgNT0JaxtBa/VCvMSc=","Signature":"4pY1m0pjUNG0YN9cRa5gTd6Pishx5Jm04GtYFmDhpg1zMUyM6ZQ4ArlAXmGisMAMjuUR1EfMzR6e53f9HS/wnQ=="},"memo":"","entropy":749259425513723904}

/*

~~~ Mine ~~~

json {
    "msg": {
        "from_address": "a709f36c78b4246cfaa0ecc363cd7c4917544c0a",
        "to_address": "f6e2eb3fb67c630d6961ef8cfb74d8d282d3e8ab",
        "amount": "50000"
    },
    "fee": [
        {
            "denom": "upokt",
            "amount": "10000"
        }
    ],
    "signature": {
        "pub_key": "03bece7df964f73215601cd844704f515ff9ab080d4f425ac6d05afd50af3127",
        "signature": "sfJPdWoW8HMRrbsw47t81z7Kv/v4S7xTBDZQyh4VKFnUKCZC6FYlWALa7xa5lV8uvxumcDVDCMJaofzXXEheDw=="
    },
    "memo": "",
    "entropy": 749259425513723904

~~~ Thereis ~~~

{
        "msg": {
            "typeUrl": "/x.nodes.MsgSend",
            "value": "pwnzbHi0JGz6oOzDY818SRdUTAr24us/tnxjDWlh74z7dNjSgtPoqzUwMDAw"
        },
        "fee": [
            {
                "amount": "10000",
                "denom": "upokt"
            }
        ],
        "memo": "",
        "entropy": 749259425513723904,
        "signature": {
            "publicKey": "A77Offlk9zIVYBzYRHBPUV/5qwgNT0JaxtBa/VCvMSc=",
            "Signature": "GS89yqTf1u9Mc/WvaEo7MmRdyVkPrno1BZi9FcnoXADJLMlKOpS/ZB+hiXpOiYV6mHEk50SHDEBXQgLUOdBwvw=="
        }
    }
}
*/
