package util

import (
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth"
)

// CompleteAndBroadcastTxCLI implements a utility function that facilitates
// sending a series of messages in a signed transaction given a TxBuilder and a
// QueryContext. It ensures that the account exists, has a proper number and
// sequence set. In addition, it builds and signs a transaction with the
// supplied messages. Finally, it broadcasts the signed transaction to a node.
func CompleteAndBroadcastTxCLI(txBldr auth.TxBuilder, cliCtx CLIContext, msgs sdk.ProtoMsg, legacyCodec bool) (*sdk.TxResponse, error) {
	//txBldr, err := PrepareTxBuilder(txBldr, cliCtx)
	//if err != nil {
	//	return nil, err
	//} TODO removed safety check for auto transactions

	// build and sign the transaction

	if cliCtx.PrivateKey != nil {
		txBytes, err := txBldr.BuildAndSign(cliCtx.FromAddress, cliCtx.PrivateKey, msgs, legacyCodec)
		if err != nil {
			return nil, err
		}
		// broadcast to a Tendermint node
		tx, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			return nil, err
		}

		return &tx, nil
	} else {
		txBytes, err := txBldr.BuildAndSignWithKeyBase(cliCtx.FromAddress, cliCtx.Passphrase, msgs, legacyCodec)
		if err != nil {
			return nil, err
		}
		// broadcast to a Tendermint node
		tx, err := cliCtx.BroadcastTx(txBytes)
		if err != nil {
			return nil, err
		}

		return &tx, nil

	}

}

// PrepareTxBuilder populates a TxBuilder in preparation for the build of a Tx.
func PrepareTxBuilder(txBldr auth.TxBuilder, cliCtx CLIContext) (auth.TxBuilder, error) {
	from := cliCtx.GetFromAddress()
	if err := cliCtx.EnsureExists(from); err != nil {
		return txBldr, err
	}
	return txBldr, nil
}

// Paginate returns the correct starting and ending index for a paginated query,
// given that client provides a desired page and limit of objects and the handler
// provides the total number of objects. If the start page is invalid, non-positive
// values are returned signaling the request is invalid.
//
// NOTE: The start page is assumed to be 1-indexed.
func Paginate(numObjs, page, limit, defLimit int) (start, end int) {
	if page == 0 {
		// invalid start page
		return -1, -1
	} else if limit == 0 {
		limit = defLimit
	}

	start = (page - 1) * limit
	end = limit + start

	if end >= numObjs {
		end = numObjs
	}

	if start >= numObjs {
		// page is out of bounds
		return -1, -1
	}

	return start, end
}
