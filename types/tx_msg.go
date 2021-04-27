package types

import "github.com/golang/protobuf/proto" // nolint

type Msg interface {
	// Return the message type.
	// Must be alphanumeric or empty.
	Route() string

	// Returns a human-readable string for the message, intended for utilization
	// within tags
	Type() string

	// ValidateBasic does a simple validation check that
	// doesn't require access to any other information.
	ValidateBasic() Error

	// Get the canonical byte representation of the ProtoMsg.
	GetSignBytes() []byte

	// Signers returns the addrs of signers that must sign.
	// CONTRACT: All signatures must be present to be valid.
	// CONTRACT: Returns addrs in some deterministic order.
	GetSigner() Address

	// Returns the recipient of the tx, if no recipient returns nil
	GetRecipient() Address

	// Returns an BigInt for the ProtoMsg
	GetFee() BigInt
}

var _ Msg = ProtoMsg(nil)

// Transactions messages must fulfill the ProtoMsg
type ProtoMsg interface {
	proto.Message
	// Return the message type.
	// Must be alphanumeric or empty.
	Route() string

	// Returns a human-readable string for the message, intended for utilization
	// within tags
	Type() string

	// ValidateBasic does a simple validation check that
	// doesn't require access to any other information.
	ValidateBasic() Error

	// Get the canonical byte representation of the ProtoMsg.
	GetSignBytes() []byte

	// Signers returns the addrs of signers that must sign.
	// CONTRACT: All signatures must be present to be valid.
	// CONTRACT: Returns addrs in some deterministic order.
	GetSigner() Address

	// Returns the recipient of the tx, if no recipient returns nil
	GetRecipient() Address

	// Returns an BigInt for the ProtoMsg
	GetFee() BigInt
}

//__________________________________________________________

// Transactions objects must fulfill the Tx
type Tx interface {
	// Gets the all the transaction's messages.
	GetMsg() Msg

	// ValidateBasic does a simple and lightweight validation check that doesn't
	// require access to any other information.
	ValidateBasic() Error
}

//__________________________________________________________

// TxDecoder unmarshals transaction bytes
type TxDecoder func(txBytes []byte, blockHeight int64) (Tx, Error)

// TxEncoder marshals transaction to bytes
type TxEncoder func(tx Tx, blockHeight int64) ([]byte, error)
