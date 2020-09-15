package types

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/codec/types"
	posCrypto "github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

var (
	_ StdTxI        = (*StdTx)(nil)
	_ StdTxI        = (*LegacyStdTx)(nil)
	_ StdSignatureI = (*StdSignature)(nil)
	_ StdSignatureI = (*LegacyStdSignature)(nil)
)

type StdTxI interface {
	sdk.Tx
	GetMemo() string
	GetSignature() StdSignatureI
	GetSigner() sdk.Address
	GetEntropy() int64
	GetFee() sdk.Coins
	WithSignature(i StdSignatureI) (StdTxI, error)
}

type StdSignatureI interface {
	GetSignature() []byte
	GetPublicKey() string // hex string
}

// StdTx is a standard way to wrap a Msg with Fee and Sigs.
// NOTE: the first signature is the fee payer (Sigs must not be nil).
func NewTx(msgs sdk.Msg, fee sdk.Coins, sig StdSignature, memo string, entropy int64, afterUpgradeHeight bool) sdk.Tx {
	if afterUpgradeHeight {
		any, _ := types.NewAnyWithValue(msgs)
		return StdTx{
			Msg:       *any,
			Fee:       fee,
			Signature: sig,
			Memo:      memo,
			Entropy:   entropy,
		}
	} else {
		var pk posCrypto.PublicKey
		if sig.PublicKey != "" {
			pk, _ = posCrypto.NewPublicKey(sig.PublicKey)
		}
		return LegacyStdTx{
			Msg: msgs,
			Fee: fee,
			Signature: LegacyStdSignature{
				PublicKey: pk,
				Signature: sig.Signature,
			},
			Memo:    memo,
			Entropy: entropy,
		}
	}
}

// GetMsg returns the all the transaction's messages.
func (tx StdTx) GetMsg() sdk.LegacyMsg {
	var res sdk.Msg
	err := ModuleCdc.ProtoCodec().UnpackAny(&tx.Msg, &res)
	if err != nil {
		panic("unable to retrive msg: " + err.Error())
	}
	return res
}

// ValidateBasic does a simple and lightweight validation check that doesn't
// require access to any other information.
func (tx StdTx) ValidateBasic() sdk.Error {
	stdSigs := tx.GetSignature()
	if !tx.Fee.IsValid() {
		return sdk.ErrInsufficientFee(fmt.Sprintf("invalid fee %s amount provided", tx.Fee.String()))
	}
	if len(stdSigs.GetSignature()) == 0 {
		return sdk.ErrUnauthorized("empty signature")
	}
	return nil
}

func (tx *StdTx) MarshalJSON() ([]byte, error) {
	pk, err := posCrypto.NewPublicKey(tx.GetSignature().GetPublicKey())
	if err != nil {
		return nil, err
	}
	res := LegacyStdTx{
		Msg: tx.GetMsg(),
		Fee: tx.Fee,
		Signature: LegacyStdSignature{
			PublicKey: pk,
			Signature: tx.GetSignature().GetSignature(),
		},
		Memo:    tx.GetMemo(),
		Entropy: tx.Entropy,
	}

	return json.Marshal(&res)
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

func (tx StdTx) GetSigner() sdk.Address {
	return tx.GetMsg().GetSigner()
}

// GetMemo returns the memo
func (tx StdTx) GetMemo() string { return tx.Memo }

func (tx StdTx) GetSignature() StdSignatureI { return &tx.Signature }

func (tx StdTx) GetEntropy() int64 {
	return tx.Entropy
}

func (tx StdTx) GetFee() sdk.Coins {
	return tx.Fee
}

func (tx StdTx) WithSignature(i StdSignatureI) (StdTxI, error) {
	sig, ok := i.(*StdSignature)
	if !ok {
		return nil, fmt.Errorf("the signature passed did not correspond to LegacyStdSignature")
	}
	tx.Signature = *sig
	return tx, nil
}

// StdSignBytes returns the bytes to sign for a transaction.
func StdSignBytes(chainID string, entropy int64, fee sdk.Coins, msg sdk.LegacyMsg, memo string) ([]byte, error) {
	msgsBytes := msg.GetSignBytes()
	var feeBytes sdk.Raw
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

// Legacy Amino Code Below
// ---------------------------------------------------------------------------------------------------------------------

type LegacyStdTx struct {
	Msg       sdk.LegacyMsg      `json:"msg" yaml:"msg"`
	Fee       sdk.Coins          `json:"fee" yaml:"fee"`
	Signature LegacyStdSignature `json:"signature" yaml:"signature"`
	Memo      string             `json:"memo" yaml:"memo"`
	Entropy   int64              `json:"entropy" yaml:"entropy"`
}

func (tx LegacyStdTx) WithSignature(i StdSignatureI) (StdTxI, error) {
	sig, ok := i.(LegacyStdSignature)
	if !ok {
		return nil, fmt.Errorf("the signature passed did not correspond to LegacyStdSignature")
	}
	tx.Signature = sig
	return tx, nil
}

func (tx LegacyStdTx) GetEntropy() int64 {
	return tx.Entropy
}

func (tx LegacyStdTx) GetFee() sdk.Coins {
	return tx.Fee
}

func (tx LegacyStdTx) GetMemo() string {
	return tx.Memo
}

func (tx LegacyStdTx) GetSignature() StdSignatureI {
	return tx.Signature
}

func (tx LegacyStdTx) GetSigner() sdk.Address {
	return tx.GetMsg().GetSigner()
}

// StdSignature represents a sig
type LegacyStdSignature struct {
	posCrypto.PublicKey `json:"pub_key" yaml:"pub_key"` // technically optional if the public key is in the world state
	Signature           []byte                          `json:"signature" yaml:"signature"`
}

func (ss LegacyStdSignature) GetSignature() []byte {
	return ss.Signature
}

func (ss LegacyStdSignature) GetPublicKey() string {
	if ss.PublicKey != nil {
		return ss.PublicKey.RawString()
	}
	return ""
}

// GetMsg returns the all the transaction's messages.
func (tx LegacyStdTx) GetMsg() sdk.LegacyMsg { return tx.Msg }

// ValidateBasic does a simple and lightweight validation check that doesn't
// require access to any other information.
func (tx LegacyStdTx) ValidateBasic() sdk.Error {
	if !tx.Fee.IsValid() {
		return sdk.ErrInsufficientFee(fmt.Sprintf("invalid fee %s amount provided", tx.Fee.String()))
	}
	if len(tx.Signature.Signature) == 0 {
		return sdk.ErrUnauthorized("empty signature")
	}
	return nil
}

// DefaultTxDecoder logic for standard transaction decoding
func DefaultTxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte) (sdk.Tx, sdk.Error) {
		if cdc.IsAfterUpgrade() {
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
		} else {
			var tx = LegacyStdTx{}
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
}

// DefaultTxEncoder logic for standard transaction encoding
func DefaultTxEncoder(cdc *codec.Codec) sdk.TxEncoder {
	return func(tx sdk.Tx) ([]byte, error) {
		if cdc.IsAfterUpgrade() {
			t, ok := tx.(StdTx)
			if !ok {
				log.Fatal("tx must be of type stdTx")
			}
			return cdc.MarshalBinaryLengthPrefixed(&t)
		}
		t, ok := tx.(LegacyStdTx)
		if !ok {
			log.Fatal("tx must be of type LegacyStdTx")
		}
		return cdc.MarshalBinaryLengthPrefixed(&t)
	}
}

// MarshalYAML returns the YAML representation of the signature.
func (ss LegacyStdSignature) MarshalYAML() (interface{}, error) {
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
	s := StdSignature{PublicKey: priv.PublicKey().RawString(), Signature: sig}
	tx := NewTx(msgs, fee, s, "", entropy, ctx.IsAfterUpgradeHeight())
	return tx
}
