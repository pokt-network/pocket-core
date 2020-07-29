package types

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/codec"
	posCrypto "github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
	"gopkg.in/yaml.v2"
	"os"
)

var (
	_ sdk.Tx = (*StdTx)(nil)
)

// StdTx is a standard way to wrap a Msg with Fee and Sigs.
// NOTE: the first signature is the fee payer (Sigs must not be nil).
type StdTx struct {
	Msg       sdk.Msg      `json:"msg" yaml:"msg"`
	Fee       sdk.Coins    `json:"fee" yaml:"fee"`
	Signature StdSignature `json:"signature" yaml:"signature"`
	Memo      string       `json:"memo" yaml:"memo"`
	Entropy   int64        `json:"entropy" yaml:"entropy"`
}

func NewStdTx(msgs sdk.Msg, fee sdk.Coins, sigs StdSignature, memo string, entropy int64) StdTx {
	return StdTx{
		Msg:       msgs,
		Fee:       fee,
		Signature: sigs,
		Memo:      memo,
		Entropy:   entropy,
	}
}

// GetMsg returns the all the transaction's messages.
func (tx StdTx) GetMsg() sdk.Msg { return tx.Msg }

// ValidateBasic does a simple and lightweight validation check that doesn't
// require access to any other information.
func (tx StdTx) ValidateBasic() sdk.Error {
	stdSigs := tx.GetSignature()
	if tx.Fee.IsValid() == false {
		return sdk.ErrInsufficientFee(fmt.Sprintf("invalid fee %s amount provided", tx.Fee.String()))
	}
	if len(stdSigs.Signature) == 0 {
		return sdk.ErrUnauthorized("empty signature")
	}
	return nil
}

// CountSubKeys counts the total number of keys for a multi-sig public key.
func CountSubKeys(pub crypto.PubKey) int {
	v, ok := pub.(multisig.PubKeyMultisigThreshold)
	if !ok {
		return 1
	}

	numKeys := 0
	for _, subkey := range v.PubKeys {
		numKeys += CountSubKeys(subkey)
	}

	return numKeys
}

// GetSigner returns the addresses that must sign the transaction.
// Addresses are returned in a deterministic order.
// They are accumulated from the GetSigner method for each Msg
// in the order they appear in tx.GetMsg().
// Duplicate addresses will be omitted.
func (tx StdTx) GetSigner() sdk.Address {
	return tx.GetMsg().GetSigner()
}

// GetMemo returns the memo
func (tx StdTx) GetMemo() string { return tx.Memo }

// GetSignature returns the signature of signers who signed the Msg.
// GetSignature returns the signature of signers who signed the Msg.
// CONTRACT: Length returned is same as length of
// pubkeys returned from MsgKeySigners, and the order
// matches.
// CONTRACT: If the signature is missing (ie the Msg is
// invalid), then the corresponding signature is
// .Empty().
func (tx StdTx) GetSignature() StdSignature { return tx.Signature }

// StdSignDoc is replay-prevention structure.
// It includes the result of msg.GetSignBytes(),
// as well as the ChainID (prevent cross chain replay)
// and the Entropy numbers for each signature (prevent
// inchain replay and enforce tx ordering per account).
type StdSignDoc struct {
	ChainID string          `json:"chain_id" yaml:"chain_id"`
	Fee     json.RawMessage `json:"fee" yaml:"fee"`
	Memo    string          `json:"memo" yaml:"memo"`
	Msg     json.RawMessage `json:"msg" yaml:"msg"`
	Entropy int64           `json:"entropy" yaml:"entropy"`
}

// StdSignBytes returns the bytes to sign for a transaction.
func StdSignBytes(chainID string, entropy int64, fee sdk.Coins, msg sdk.Msg, memo string) ([]byte, error) {
	msgsBytes := msg.GetSignBytes()
	var feeBytes json.RawMessage
	feeBytes, err := fee.MarshalJSON()
	if err != nil {
		return nil, fmt.Errorf("could not marshal fee to json for StdSignBytes function: %v", err.Error())
	}
	bz, err := ModuleCdc.MarshalJSON(StdSignDoc{
		ChainID: chainID,
		Fee:     feeBytes,
		Memo:    memo,
		Msg:     msgsBytes,
		Entropy: entropy,
	})
	if err != nil {
		return nil, fmt.Errorf("could not marshal bytes to json for StdSignDoc function: %v", err.Error())
	}
	return sdk.MustSortJSON(bz), nil
}

// StdSignature represents a sig
type StdSignature struct {
	posCrypto.PublicKey `json:"pub_key" yaml:"pub_key"` // technically optional if the public key is in the world state
	Signature           []byte                          `json:"signature" yaml:"signature"`
}

// DefaultTxDecoder logic for standard transaction decoding
func DefaultTxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, sdk.Error) {
		var tx = StdTx{}
		if len(txBytes) == 0 {
			return nil, sdk.ErrTxDecode("txBytes are empty")
		}
		// StdTx.Msg is an interface. The concrete types
		// are registered by MakeTxCodec
		err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx)
		if err != nil {
			return nil, sdk.ErrTxDecode("error decoding transaction").TraceSDK(err.Error())
		}
		return tx, nil
	}
}

// DefaultTxEncoder logic for standard transaction encoding
func DefaultTxEncoder(cdc *codec.Codec) sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		return cdc.MarshalBinaryLengthPrefixed(tx)
	}
}

// MarshalYAML returns the YAML representation of the signature.
func (ss StdSignature) MarshalYAML() (interface{}, error) {
	var (
		bz     []byte
		pubkey string
		err    error
	)
	if ss.PublicKey != nil {
		pubkey = ss.PublicKey.RawString()
	}
	bz, err = yaml.Marshal(struct {
		PubKey    string
		Signature string
	}{
		PubKey:    pubkey,
		Signature: fmt.Sprintf("%s", ss.Signature),
	})
	if err != nil {
		return nil, err
	}
	return string(bz), err
}

func NewTestTx(ctx sdk.Ctx, msgs sdk.Msg, priv posCrypto.PrivateKey, entropy int64, fee sdk.Coins) sdk.Tx {
	signBytes, err := StdSignBytes(ctx.ChainID(), entropy, fee, msgs, "")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	sig, err := priv.Sign(signBytes)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s := StdSignature{PublicKey: priv.PublicKey(), Signature: sig}
	tx := NewStdTx(msgs, fee, s, "", entropy)
	return tx
}
