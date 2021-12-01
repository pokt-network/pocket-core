package auth

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/pokt-network/pocket-core/codec"
	posCrypto "github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth/keeper"
	"github.com/pokt-network/pocket-core/x/auth/types"
	"github.com/tendermint/tendermint/state/txindex"
	tmTypes "github.com/tendermint/tendermint/types"
	"os"
)

// NewAnteHandler returns an AnteHandler that checks signatures and deducts fees from the first signer.
func NewAnteHandler(ak keeper.Keeper) sdk.AnteHandler {
	return func(ctx sdk.Ctx, tx sdk.Tx, txBz []byte, txIndexer txindex.TxIndexer, simulate bool) (newCtx sdk.Ctx, res sdk.Result, signer posCrypto.PublicKey, abort bool) {
		if addr := ak.GetModuleAddress(types.FeeCollectorName); addr == nil {
			ctx.Logger().Error(fmt.Sprintf("%s module account has not been set", types.FeeCollectorName))
			os.Exit(1)
		}
		// validate the transaction
		if err := tx.ValidateBasic(); err != nil {
			return newCtx, err.Result(), nil, true
		}
		stdTx, ok := tx.(types.StdTx)
		if !ok {
			return newCtx, sdk.ErrInternal("all transactions must be convertible to inteface: ProtoStdTx").Result(), nil, true
		}
		signer, err := ValidateTransaction(ctx, ak, stdTx, ak.GetParams(ctx), txIndexer, txBz, simulate)
		if err != nil {
			return newCtx, err.Result(), signer, true
		}
		err = DeductFees(ak, ctx, stdTx, signer)
		if err != nil {
			return newCtx, err.Result(), signer, true
		}
		return ctx, sdk.Result{}, signer, false // continue...
	}
}

func ValidateTransaction(ctx sdk.Ctx, k Keeper, stdTx types.StdTx, params Params, txIndexer txindex.TxIndexer, txBz []byte, simulate bool) (signer posCrypto.PublicKey, sdkErr sdk.Error) {
	// validate the memo
	if err := ValidateMemo(stdTx, params); err != nil {
		return nil, types.ErrInvalidMemo(ModuleName, err)
	}
	// check for duplicate transaction to prevent replay attacks
	txHash := tmTypes.Tx(txBz).Hash()
	// make http call to tendermint to check txIndexer
	if txIndexer == nil {
		ctx.Logger().Error(types.ErrNilTxIndexer(ModuleName).Error())
		return nil, types.ErrNilTxIndexer(ModuleName)
	}
	res, err := (txIndexer).Get(txHash)
	if err != nil {
		ctx.Logger().Error(err.Error())
		return nil, sdk.ErrInternal(err.Error())
	}
	if res != nil {
		return nil, types.ErrDuplicateTx(ModuleName, hex.EncodeToString(txHash))
	}
	var pk posCrypto.PublicKey
	for _, signer := range stdTx.GetSigners() {
		// attempt to get the public key from the signature
		if stdTx.GetSignature().GetPublicKey() != "" {
			var err error
			pk, err = posCrypto.NewPublicKey(stdTx.GetSignature().GetPublicKey())
			if err != nil {
				return nil, sdk.ErrInvalidPubKey(err.Error())
			}
		} else {
			// public key in the signature not found so check world state
			acc := k.GetAccount(ctx, signer)
			if acc == nil {
				return nil, types.ErrAccountNotFound(ModuleName)
			}
			if pk = acc.GetPubKey(); pk == nil {
				return nil, types.ErrEmptyPublicKey(ModuleName)
			}
		}
		//patch sync fix : add Verify against after codec upgrade chainhalt height
		if !bytes.Equal(pk.Address(), signer) && ctx.BlockHeight() != codec.CodecChainHaltHeight {
			continue
		}
		// get the sign bytes from the tx
		signBytes, err := GetSignBytes(ctx.ChainID(), stdTx)
		if err != nil {
			return nil, sdk.ErrInternal(err.Error())
		}
		// get the fees from the tx
		expectedFee := sdk.NewCoins(sdk.NewCoin(sdk.DefaultStakeDenom, k.GetParams(ctx).FeeMultiplier.GetFee(stdTx.GetMsg())))
		// test for public key type
		p, ok := pk.(posCrypto.PublicKeyMultiSig)
		// if standard public key
		if !ok {
			// validate the fees for a standard public key
			if !stdTx.GetFee().IsAllGTE(expectedFee) {
				return nil, types.ErrInsufficientFee(ModuleName, expectedFee, stdTx.GetFee())
			}
			// validate signature for regular public key
			if !simulate && !pk.VerifyBytes(signBytes, stdTx.GetSignature().GetSignature()) {
				continue
			}
			return pk, nil
		}
		// validate the signature depth
		ok = ValidateSignatureDepth(params.TxSigLimit, p)
		if !ok {
			return nil, types.ErrTooManySignatures(ModuleName, params.TxSigLimit)
		}
		// validate the multi sig
		if !simulate && !pk.VerifyBytes(signBytes, stdTx.GetSignature().GetSignature()) {
			continue
		}
		return pk, nil
	}
	return nil, sdk.ErrUnauthorized("signature verification failed for the transaction")
}

func ValidateSignatureDepth(limit uint64, publicKey posCrypto.PublicKeyMultiSig) (ok bool) {
	_, ok = recSignDepth(1, limit, publicKey)
	return
}

// recSignDepth ensures that the number of signatures does not exceed the max
func recSignDepth(count, limit uint64, publicKey posCrypto.PublicKeyMultiSig) (c uint64, ok bool) {
	for _, p := range publicKey.Keys() {
		count++
		if pk, ok := p.(posCrypto.PublicKeyMultiSig); ok {
			count, ok = recSignDepth(count, limit, pk)
			if !ok {
				return count, ok
			}
		}
		if count > limit {
			return count, false
		}
	}
	return count, true
}

// GetSignerAcc returns an account for a given address that is expected to sign
// a transaction.
func GetSignerAcc(ctx sdk.Ctx, ak keeper.Keeper, addr sdk.Address) (Account, sdk.Error) {
	if acc := ak.GetAccount(ctx, addr); acc != nil {
		return acc, nil
	}
	return nil, sdk.ErrUnknownAddress(fmt.Sprintf("account %s does not exist", addr))
}

// ValidateMemo validates the memo size.
func ValidateMemo(stdTx types.StdTx, params Params) sdk.Error {
	memoLength := len(stdTx.GetMemo())
	if uint64(memoLength) > params.MaxMemoCharacters {
		return sdk.ErrMemoTooLarge(
			fmt.Sprintf(
				"maximum number of characters is %d but received %d characters",
				params.MaxMemoCharacters, memoLength,
			),
		)
	}
	return nil
}

// DeductFees deducts fees from the given account.
func DeductFees(keeper keeper.Keeper, ctx sdk.Ctx, tx types.StdTx, signer posCrypto.PublicKey) sdk.Error {
	fees := tx.GetFee()
	if !fees.IsValid() {
		return sdk.ErrInsufficientFee(fmt.Sprintf("invalid fee amount: %s", fees))
	}
	var acc Account
	var err sdk.Error

	if keeper.Cdc.IsAfterNonCustodialUpgrade(ctx.BlockHeight()) {
		acc, err = GetSignerAcc(ctx, keeper, sdk.Address(signer.Address()))
		if err != nil {
			return err
		}
	} else {
		acc, err = GetSignerAcc(ctx, keeper, tx.GetSigners()[0])
		if err != nil {
			return err
		}
	}

	coins := acc.GetCoins()
	// verify the account has enough funds to pay for fees
	_, hasNeg := coins.SafeSub(fees)
	if hasNeg {
		return types.ErrInsufficientBalance(ModuleName, acc.GetAddress(), fees)
	}
	err = keeper.SendCoinsFromAccountToModule(ctx, acc.GetAddress(), types.FeeCollectorName, fees)
	if err != nil {
		return err
	}
	return nil
}

// GetSignBytes returns a slice of bytes to sign over for a given transaction
// and an account.
func GetSignBytes(chainID string, stdTx types.StdTx) ([]byte, error) {
	return StdSignBytes(
		chainID, stdTx.GetEntropy(), stdTx.GetFee(), stdTx.GetMsg(), stdTx.GetMemo(),
	)
}
