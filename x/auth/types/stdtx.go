package types

import (
	"encoding/json"
	"fmt"
	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/codec/types"
	posCrypto "github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	types2 "github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/tendermint/tendermint/crypto"
	"github.com/tendermint/tendermint/crypto/multisig"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

// ProtoStdTx is a standard way to wrap a ProtoMsg with Fee and Sigs.
// NOTE: the first signature is the fee payer (Sigs must not be nil).
func NewTx(msgs sdk.ProtoMsg, fee sdk.Coins, sig StdSignature, memo string, entropy int64) sdk.Tx {
	return StdTx{
		Msg:       msgs,
		Fee:       fee,
		Signature: sig,
		Memo:      memo,
		Entropy:   entropy,
	}
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

// StdSignBytes returns the bytes to sign for a transaction.
func StdSignBytes(chainID string, entropy int64, fee sdk.Coins, msg sdk.Msg, memo string) ([]byte, error) {
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

var _ codec.ProtoMarshaler = &StdTx{}

type StdTx struct {
	Msg       sdk.Msg      `json:"msg" yaml:"msg"`
	Fee       sdk.Coins    `json:"fee" yaml:"fee"`
	Signature StdSignature `json:"signature" yaml:"signature"`
	Memo      string       `json:"memo" yaml:"memo"`
	Entropy   int64        `json:"entropy" yaml:"entropy"`
}

func (tx *StdTx) Reset() {
	*tx = StdTx{}
}

func (tx StdTx) String() string {
	p, _ := tx.ToProto()
	return p.String()
}

func (tx StdTx) ProtoMessage() {
	p, _ := tx.ToProto()
	p.ProtoMessage()
}

func (tx StdTx) Marshal() ([]byte, error) {
	p, err := tx.ToProto()
	if err != nil {
		return nil, err
	}
	return p.Marshal()
}

func (tx StdTx) MarshalTo(data []byte) (n int, err error) {
	p, err := tx.ToProto()
	if err != nil {
		return 0, err
	}
	return p.MarshalTo(data)
}

func (tx StdTx) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	p, err := tx.ToProto()
	if err != nil {
		return 0, err
	}
	return p.MarshalToSizedBuffer(dAtA)
}

func (tx StdTx) Size() int {
	p, _ := tx.ToProto()
	return p.Size()
}

func (tx *StdTx) Unmarshal(data []byte) error {
	var pstdtx ProtoStdTx
	err := pstdtx.Unmarshal(data)
	if err != nil {
		return err
	}
	stdTx, err := pstdtx.FromProto()
	if err != nil {
		return err
	}
	*tx = stdTx
	return nil
}

func (tx StdTx) ToProto() (ProtoStdTx, error) {
	pMsg, ok := tx.Msg.(sdk.ProtoMsg)
	if !ok {
		return ProtoStdTx{}, fmt.Errorf("unable to convert sdk.Msg to sdk.ProtoMsg: %v", tx.Msg)
	}
	any, err := types.NewAnyWithValue(pMsg)
	if err != nil {
		return ProtoStdTx{}, fmt.Errorf("unable to convert sdk.ProtoMsg into any %v", pMsg)
	}
	return ProtoStdTx{
		Msg:       *any,
		Fee:       tx.Fee,
		Signature: tx.Signature.ToProto(),
		Memo:      tx.Memo,
		Entropy:   tx.Entropy,
	}, nil
}

func (tx StdTx) WithSignature(sig StdSignature) (StdTx, error) {
	tx.Signature = sig
	return tx, nil
}

func (tx StdTx) GetEntropy() int64 {
	return tx.Entropy
}

func (tx StdTx) GetFee() sdk.Coins {
	return tx.Fee
}

func (tx StdTx) GetMemo() string {
	return tx.Memo
}

func (tx StdTx) GetSignature() StdSignature {
	return tx.Signature
}

func (tx StdTx) GetSigners() []sdk.Address {
	return tx.GetMsg().GetSigners()
}

// GetMsg returns the all the transaction's messages.
func (tx StdTx) GetMsg() sdk.Msg { return tx.Msg }

// ValidateBasic does a simple and lightweight validation check that doesn't
// require access to any other information.
func (tx StdTx) ValidateBasic() sdk.Error {
	if !tx.Fee.IsValid() {
		return sdk.ErrInsufficientFee(fmt.Sprintf("invalid fee %s amount provided", tx.Fee.String()))
	}
	if len(tx.Signature.Signature) == 0 {
		return sdk.ErrUnauthorized("empty signature")
	}
	return nil
}

func (ptx ProtoStdTx) FromProto() (StdTx, error) {
	var res sdk.ProtoMsg
	err := ModuleCdc.ProtoCodec().UnpackAny(&ptx.Msg, &res)
	if err != nil {
		return StdTx{}, err
	}
	ss, err := ptx.Signature.FromProto()
	if err != nil {
		return StdTx{}, err
	}
	return StdTx{
		Msg:       res,
		Fee:       ptx.Fee,
		Signature: ss,
		Memo:      ptx.Memo,
		Entropy:   ptx.Entropy,
	}, nil
}

func (ptx *ProtoStdTx) MarshalJSON() ([]byte, error) {
	s, err := ptx.FromProto()
	if err != nil {
		return nil, err
	}
	return json.Marshal(&s)
}

var _ codec.ProtoMarshaler = &StdSignature{}

// ProtoStdSignature represents a sig
type StdSignature struct {
	posCrypto.PublicKey `json:"pub_key" yaml:"pub_key"` // technically optional if the public key is in the world state
	Signature           []byte                          `json:"signature" yaml:"signature"`
}

func (ss *StdSignature) Reset() {
	*ss = StdSignature{}
}

func (ss StdSignature) ProtoMessage() {
	p := ss.ToProto()
	p.ProtoMessage()
}

func (ss StdSignature) Marshal() ([]byte, error) {
	p := ss.ToProto()
	return p.Marshal()
}

func (ss StdSignature) MarshalTo(data []byte) (n int, err error) {
	p := ss.ToProto()
	return p.MarshalTo(data)
}

func (ss StdSignature) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	p := ss.ToProto()
	return p.MarshalToSizedBuffer(dAtA)
}

func (ss *StdSignature) Unmarshal(data []byte) error {
	var pss ProtoStdSignature
	err := pss.Unmarshal(data)
	if err != nil {
		return err
	}
	s, err := pss.FromProto()
	if err != nil {
		return err
	}
	*ss = s
	return nil
}

func (ss StdSignature) ToProto() ProtoStdSignature {
	return ProtoStdSignature{
		PublicKey: ss.PublicKey.RawBytes(),
		Signature: ss.Signature,
	}
}

func (ss ProtoStdSignature) FromProto() (sig StdSignature, err error) {
	var pk posCrypto.PublicKey
	if ss.PublicKey != nil {
		pk, err = posCrypto.NewPublicKeyBz(ss.PublicKey)
		if err != nil {
			return StdSignature{}, err
		}
	}
	return StdSignature{
		PublicKey: pk,
		Signature: ss.Signature,
	}, nil
}

func (ss StdSignature) GetSignature() []byte {
	return ss.Signature
}

func (ss StdSignature) GetPublicKey() string {
	if ss.PublicKey != nil {
		return ss.PublicKey.RawString()
	}
	return ""
}

// DefaultTxDecoder logic for standard transaction decoding
func DefaultTxDecoder(cdc *codec.Codec) sdk.TxDecoder {
	return func(txBytes []byte, blockHeight int64) (sdk.Tx, sdk.Error) {
		var tx = StdTx{}
		if len(txBytes) == 0 {
			return nil, sdk.ErrTxDecode("txBytes are empty")
		}
		// ProtoStdTx.ProtoMsg is an interface. The concrete types
		// are registered by MakeTxCodec
		err := cdc.UnmarshalBinaryLengthPrefixed(txBytes, &tx, blockHeight)

		//replicate error on new stake msg sent before upgrade block for compatibility reasons (happened on 56550 BU)
		if !cdc.IsAfterNonCustodialUpgrade(blockHeight) {
			if _, ok := tx.Msg.(*types2.MsgStake); ok {
				return nil, sdk.ErrTxDecode("error decoding transaction: no concrete type registered for type URL /x.nodes.MsgProtoStake8 against interface *types.ProtoMsg")
			}
			if _, ok := tx.Msg.(*types2.MsgUnjail); ok {
				return nil, sdk.ErrTxDecode("error decoding transaction: no concrete type registered for type URL /x.nodes.MsgUnjail8 against interface *types.ProtoMsg")
			}
			if _, ok := tx.Msg.(*types2.MsgBeginUnstake); ok {
				return nil, sdk.ErrTxDecode("error decoding transaction: no concrete type registered for type URL /x.nodes.MsgBeginUnstake8 against interface *types.ProtoMsg")
			}
		}

		if err != nil {
			return nil, sdk.ErrTxDecode("error decoding transaction: " + err.Error()).TraceSDK(err.Error())
		}
		return tx, nil
	}
}

// DefaultTxEncoder logic for standard transaction encoding
func DefaultTxEncoder(cdc *codec.Codec) sdk.TxEncoder {
	return func(tx sdk.Tx, blockHeight int64) ([]byte, error) {
		t, ok := tx.(StdTx)
		if !ok {
			log.Fatal("tx must be of type StdTx")
		}
		return cdc.MarshalBinaryLengthPrefixed(&t, blockHeight)
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

func NewTestTx(ctx sdk.Ctx, msgs sdk.ProtoMsg, priv posCrypto.PrivateKey, entropy int64, fee sdk.Coins) sdk.Tx {
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
	tx := NewTx(msgs, fee, s, "", entropy)
	return tx
}
