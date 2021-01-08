package types

import (
	"errors"
	"fmt"
	"github.com/tendermint/tendermint/libs/rand"
	"strings"

	"github.com/pokt-network/pocket-core/crypto"

	crkeys "github.com/pokt-network/pocket-core/crypto/keys"
	sdk "github.com/pokt-network/pocket-core/types"
)

// TxBuilder implements a transaction context created in SDK modules.
type TxBuilder struct {
	txEncoder sdk.TxEncoder
	txDecoder sdk.TxDecoder
	keybase   crkeys.Keybase
	chainID   string
	memo      string
	fees      sdk.Coins
}

// NewTxBuilder returns a new initialized TxBuilder.
func NewTxBuilder(txEncoder sdk.TxEncoder, txDecoder sdk.TxDecoder, chainID, memo string, fees sdk.Coins) TxBuilder {
	return TxBuilder{
		txEncoder: txEncoder,
		txDecoder: txDecoder,
		keybase:   nil,
		chainID:   chainID,
		memo:      memo,
		fees:      fees,
	}
}

// TxEncoder returns the transaction encoder
func (bldr TxBuilder) TxEncoder() sdk.TxEncoder { return bldr.txEncoder }

// TxEncoder returns the transaction encoder
func (bldr TxBuilder) TxDecoder() sdk.TxDecoder { return bldr.txDecoder }

// Keybase returns the keybase
func (bldr TxBuilder) Keybase() crkeys.Keybase { return bldr.keybase }

// ChainID returns the chain id
func (bldr TxBuilder) ChainID() string { return bldr.chainID }

// Memo returns the memo message
func (bldr TxBuilder) Memo() string { return bldr.memo }

// Fees returns the fees for the transaction
func (bldr TxBuilder) Fees() sdk.Coins { return bldr.fees }

// WithTxEncoder returns a copy of the context with an updated codec.
func (bldr TxBuilder) WithTxEncoder(txEncoder sdk.TxEncoder) TxBuilder {
	bldr.txEncoder = txEncoder
	return bldr
}

// WithChainID returns a copy of the context with an updated chainID.
func (bldr TxBuilder) WithChainID(chainID string) TxBuilder {
	bldr.chainID = chainID
	return bldr
}

// WithFees returns a copy of the context with an updated fee.
func (bldr TxBuilder) WithFees(fees string) TxBuilder {
	parsedFees, err := sdk.ParseCoins(fees)
	if err != nil {
		fmt.Println(fmt.Errorf("error adding fees to the tx builder: " + err.Error()))
		return bldr
	}

	bldr.fees = parsedFees
	return bldr
}

// WithKeybase returns a copy of the context with updated keybase.
func (bldr TxBuilder) WithKeybase(keybase crkeys.Keybase) TxBuilder {
	bldr.keybase = keybase
	return bldr
}

// WithMemo returns a copy of the context with an updated memo.
func (bldr TxBuilder) WithMemo(memo string) TxBuilder {
	bldr.memo = strings.TrimSpace(memo)
	return bldr
}

// BuildAndSign builds a single message to be signed, and signs a transaction
// with the built message given a address, private key, and a set of messages.
func (bldr TxBuilder) BuildAndSign(address sdk.Address, privateKey crypto.PrivateKey, msg sdk.ProtoMsg, legacyCodec bool) ([]byte, error) {
	if bldr.chainID == "" {
		return nil, errors.New("cant build and sign transaciton: the chainID is empty")
	}
	entropy := rand.Int64()
	bytesToSign, err := StdSignBytes(bldr.chainID, entropy, bldr.fees, msg, bldr.memo)
	if err != nil {
		return nil, err
	}
	sigBytes, err := privateKey.Sign(bytesToSign)
	if err != nil {
		return nil, err
	}
	sig := StdSignature{
		Signature: sigBytes,
		PublicKey: privateKey.PublicKey(),
	}
	if legacyCodec {
		return bldr.txEncoder(NewTx(msg, bldr.fees, sig, bldr.memo, entropy), 0)
	}
	return bldr.txEncoder(NewTx(msg, bldr.fees, sig, bldr.memo, entropy), -1)
}

// BuildAndSignWithKeyBase builds a single message to be signed, and signs a transaction
// with the built message given a address, passphrase, and a set of messages.
func (bldr TxBuilder) BuildAndSignWithKeyBase(address sdk.Address, passphrase string, msg sdk.ProtoMsg, legacyCodec bool) ([]byte, error) {
	if bldr.keybase == nil {
		return nil, errors.New("cant build and sign transaciton: the keybase is nil")
	}
	if bldr.chainID == "" {
		return nil, errors.New("cant build and sign transaciton: the chainID is empty")
	}
	entropy := rand.Int64()
	bytesToSign, err := StdSignBytes(bldr.chainID, entropy, bldr.fees, msg, bldr.memo)
	if err != nil {
		return nil, err
	}
	sigBytes, pk, err := bldr.keybase.Sign(address, passphrase, bytesToSign)
	if err != nil {
		return nil, err
	}
	sig := StdSignature{
		Signature: sigBytes,
		PublicKey: pk,
	}
	if legacyCodec {
		return bldr.txEncoder(NewTx(msg, bldr.fees, sig, bldr.memo, entropy), 0)
	}
	return bldr.txEncoder(NewTx(msg, bldr.fees, sig, bldr.memo, entropy), -1)
}

func (bldr TxBuilder) SignMultisigTransaction(address sdk.Address, keys []crypto.PublicKey, passphrase string, txBytes []byte, legacyCodec bool) (signedTx []byte, err error) {
	if bldr.keybase == nil {
		return nil, errors.New("cant build and sign transaciton: the keybase is nil")
	}
	if bldr.chainID == "" {
		return nil, errors.New("cant build and sign transaciton: the chainID is empty")
	}
	t := sdk.Tx(nil)
	// decode the transaction
	if legacyCodec {
		t, err = bldr.txDecoder(txBytes, 0)
		if err != nil {
			return nil, err
		}
	} else {
		t, err = bldr.txDecoder(txBytes, -1)
		if err != nil {
			return nil, err
		}
	}
	tx := t.(StdTx)
	// get the sign bytes from the transaction
	bytesToSign, err := StdSignBytes(bldr.chainID, tx.GetEntropy(), tx.GetFee(), tx.GetMsg(), tx.GetMemo())
	if err != nil {
		return nil, err
	}
	sigBytes, pubKey, err := bldr.keybase.Sign(address, passphrase, bytesToSign)
	if err != nil {
		return nil, err
	}
	// sign using multisignature sturcture
	var ms = crypto.MultiSig(crypto.MultiSignature{})
	if tx.GetSignature().GetSignature() == nil || len(tx.GetSignature().GetSignature()) == 0 {
		ms = ms.NewMultiSignature()
	} else {
		ms = ms.Unmarshal(tx.GetSignature().GetSignature())
	}
	if len(keys) != 0 {
		ms, err = ms.AddSignature(sigBytes, pubKey, keys)
		if err != nil {
			return nil, err
		}
	} else {
		ms = ms.AddSignatureByIndex(sigBytes, len(ms.Signatures()))
	}
	sig := StdSignature{
		PublicKey: tx.Signature.PublicKey,
		Signature: ms.Marshal(),
	}
	// replace the old multi-signature with the new multi-signature (containing the additional signature)
	tx, err = tx.WithSignature(sig)
	if err != nil {
		return nil, err
	}
	if legacyCodec {
		return bldr.TxEncoder()(tx, 0)
	}
	// encode using the standard encoder
	return bldr.TxEncoder()(tx, -1)
}

func (bldr TxBuilder) BuildAndSignMultisigTransaction(address sdk.Address, publicKey crypto.PublicKeyMultiSig, m sdk.ProtoMsg, passphrase string, fees int64, legacyCodec bool) (signedTx []byte, err error) {
	if bldr.keybase == nil {
		return nil, errors.New("cant build and sign transaciton: the keybase is nil")
	}
	if bldr.chainID == "" {
		return nil, errors.New("cant build and sign transaciton: the chainID is empty")
	}
	// bulid the transaction from scratch
	entropy := rand.Int64()
	fee := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, sdk.NewInt(fees)))
	signBz, err := StdSignBytes(bldr.chainID, entropy, fee, m, bldr.memo)
	if err != nil {
		return nil, err
	}
	sigBytes, _, err := bldr.keybase.Sign(address, passphrase, signBz)
	if err != nil {
		return nil, err
	}
	// sign using multisignature structure
	var ms = crypto.MultiSig(crypto.MultiSignature{})
	ms = ms.NewMultiSignature()
	ms = ms.AddSignatureByIndex(sigBytes, 0)
	sig := StdSignature{
		PublicKey: publicKey,
		Signature: ms.Marshal(),
	}
	// create a new standard transaction object
	tx := NewTx(m, fee, sig, "", entropy)
	// encode it using the default encoder
	if legacyCodec {
		return bldr.TxEncoder()(tx, 0)
	}
	return bldr.TxEncoder()(tx, -1)
}
