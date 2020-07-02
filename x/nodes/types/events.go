package types

// pos module event types
const (
	EventTypeCompleteUnstaking       = "complete_unstaking"
	EventTypeCreateValidator         = "create_validator"
	EventTypeStake                   = "stake"
	EventTypeBeginUnstake            = "begin_unstake"
	EventTypeWaitingToBeginUnstaking = "waiting_to_begin_unstaking"
	EventTypeUnstake                 = "unstake"
	EventTypeProposerReward          = "proposer_reward"
	EventTypeDAOAllocation           = "dao_allocation"
	EventTypeSlash                   = "slash"
	EventTypeJail                    = "jail"
	EventTypeLiveness                = "liveness"
	AttributeKeyAddress              = "address"
	AttributeKeyHeight               = "height"
	AttributeKeyPower                = "power"
	AttributeKeyReason               = "reason"
	AttributeKeyJailed               = "jailed"
	AttributeKeyMissedBlocks         = "missed_blocks"
	AttributeValueDoubleSign         = "double_sign"
	AttributeValueMissingSignature   = "missing_signature"
	AttributeKeyValidator            = "validator"
	AttributeValueCategory           = ModuleName
)
