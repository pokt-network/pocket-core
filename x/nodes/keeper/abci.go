package keeper

import (
	"fmt"
	sdk "github.com/pokt-network/pocket-core/types"
	abci "github.com/tendermint/tendermint/abci/types"
	tmtypes "github.com/tendermint/tendermint/types"
)

// BeginBlocker - Called at the beggining of every block
// 1) allocate tokens to block producer
// 2) mint any custom awards for each validator
// 3) set new proposer
// 4) check block sigs and byzantine evidence to slash
func BeginBlocker(ctx sdk.Ctx, req abci.RequestBeginBlock, k Keeper) {
	// reward the proposer with fees
	if ctx.BlockHeight() > 1 {
		previousProposer := k.GetPreviousProposer(ctx)
		k.blockReward(ctx, previousProposer)
	}
	// record the new proposer for when we payout on the next block
	addr := sdk.Address(req.Header.ProposerAddress)
	k.SetPreviousProposer(ctx, addr)
	// Iterate over all the validators which *should* have signed this block
	// store whether or not they have actually signed it and slash/unstake any
	// which have missed too many blocks in a row (downtime slashing)
	signedBlocksWindow := k.SignedBlocksWindow(ctx)
	minSignedPerWindow := k.MinSignedPerWindow(ctx)
	downtimeJailDuration := k.DowntimeJailDuration(ctx)
	slashFractionDowntime := k.SlashFractionDowntime(ctx)

	for _, voteInfo := range req.LastCommitInfo.GetVotes() {
		k.handleValidatorSignature(ctx, voteInfo.Validator.Address, voteInfo.Validator.Power, voteInfo.SignedLastBlock, signedBlocksWindow, minSignedPerWindow, downtimeJailDuration, slashFractionDowntime)
		// remove those who are part of the tendermint validator set (jailed validators will never be a part of the set)
	}
	// Iterate through any newly discovered evidence of infraction
	// slash any validators (and since-unstaked stake within the unstaking period)
	// who contributed to valid infractions
	maxEvidenceAgeMin := int(k.GetParams(ctx).MaxEvidenceAge.Minutes())
	for _, evidence := range req.ByzantineValidators {
		switch evidence.Type {
		case tmtypes.ABCIEvidenceTypeDuplicateVote:
			if ctx.IsAfterUpgradeHeight() {
				evidenceAgeInBlocks := maxEvidenceAgeMin / 15
				if evidenceAgeInBlocks == 0 {
					// minimum of 1 block ago
					evidenceAgeInBlocks = 1
				}
				if ctx.BlockHeight()-evidence.Height <= int64(evidenceAgeInBlocks) {
					k.handleDoubleSign(ctx, evidence.Validator.Address, evidence.Height, evidence.Time, evidence.Validator.Power)
				} else {
					k.Logger(ctx).Error(fmt.Sprintf("disregarded %s evidence, max evidence age exceeded, height of evidence=%d", evidence.Type, evidence.Height))
				}
			} else {
				k.handleDoubleSign(ctx, evidence.Validator.Address, evidence.Height, evidence.Time, evidence.Validator.Power)
			}
		default:
			k.Logger(ctx).Error(fmt.Sprintf("ignored unknown evidence type: %s", evidence.Type))
		}
	}
}

// EndBlocker - Called at the end of every block, update validator set
func EndBlocker(ctx sdk.Ctx, k Keeper) []abci.ValidatorUpdate {
	// increment jailed blocks counter
	k.IncrementJailedValidators(ctx)
	// NOTE: UpdateTendermintValidators has to come before unstakeAllMatureValidators.
	validatorUpdates := k.UpdateTendermintValidators(ctx)
	// Unstake all mature validators from the unstakeing queue.
	k.unstakeAllMatureValidators(ctx)
	return validatorUpdates
}
