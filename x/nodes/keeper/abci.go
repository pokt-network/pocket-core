package keeper

import (
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// 1) allocate tokens to block producer
// 2) mint any custom awards for each validator
// 3) set new proposer
// 4) check block sigs and byzantine evidence to slash
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k Keeper) {
	// reward the proposer with fees
	if ctx.BlockHeight() > 1 {
		previousProposer := k.GetPreviousProposer(ctx)
		k.blockReward(ctx, previousProposer)
	}
	// mint any custom validator awards
	k.mintNodeRelayRewards(ctx)
	// burn any custom validator slashes
	k.burnValidators(ctx)
	// record the new proposer for when we payout on the next block
	addr := sdk.Address(req.Header.ProposerAddress)
	k.SetPreviousProposer(ctx, addr)
	// Iterate over all the validators which *should* have signed this block
	// store whether or not they have actually signed it and slash/unstake any
	// which have missed too many blocks in a row (downtime slashing)
	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		k.handleValidatorSignature(ctx, voteInfo.Validator.Address, voteInfo.Validator.Power, voteInfo.SignedLastBlock)
	}
	// Iterate through any newly discovered evidence of infraction
	// slash any validators (and since-unstaked stake within the unstaking period)
	// who contributed to valid infractions
	for _, evidence := range req.ByzantineValidators {
		switch evidence.Type {
		case tmtypes.ABCIEvidenceTypeDuplicateVote:
			k.handleDoubleSign(ctx, evidence.Validator.Address, evidence.Height, evidence.Time, evidence.Validator.Power)
		default:
			k.Logger(ctx).Error(fmt.Sprintf("ignored unknown evidence type: %s", evidence.Type))
		}
	}
}

// Called every block, update validator set
func EndBlocker(ctx sdk.Context, k Keeper) []abci.ValidatorUpdate {
	// Calculate validator set changes.
	// NOTE: UpdateTendermintValidators has to come before unstakeAllMatureValidators.
	validatorUpdates := k.UpdateTendermintValidators(ctx)
	// Unstake all mature validators from the unstakeing queue.
	k.unstakeAllMatureValidators(ctx)
	return validatorUpdates
}
