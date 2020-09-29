package nodes

import (
	"fmt"
	"os"
	"time"

	"github.com/pokt-network/pocket-core/x/nodes/keeper"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis sets up the module based on the genesis state
// First TM block is at height 1, so state updates applied from
// genesis.json are in block 0.
func InitGenesis(ctx sdk.Ctx, keeper keeper.Keeper, supplyKeeper types.AuthKeeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	// zero out a staked tokens variable for traking the number of staked tokens
	stakedTokens := sdk.ZeroInt()
	// set the context
	ctx = ctx.WithBlockHeight(1 - sdk.ValidatorUpdateDelay)
	// set the parameters from the data
	keeper.SetParams(ctx, data.Params)
	// set the 'previous state total power' from the data
	keeper.SetPrevStateValidatorsPower(ctx, data.PrevStateTotalPower)
	// for each validator in validators, setup based on genesis file
	for _, validator := range data.Validators {
		// if the validator is unstaked, then panic because, we shouldn't have unstaked validators in the genesis file
		if validator.IsUnstaked() {
			keeper.Logger(ctx).Error(fmt.Errorf("%v the validators must be staked or unstaking at genesis", validator).Error())
			os.Exit(1)
		}
		// set the validators from the data
		keeper.SetValidator(ctx, validator)
		keeper.SetStakedValidatorByChains(ctx, validator)
		// ensure there's a signing info entry for the validator (used in slashing)
		_, found := keeper.GetValidatorSigningInfo(ctx, validator.GetAddress())
		if !found {
			signingInfo := types.ValidatorSigningInfo{
				Address:     validator.GetAddress(),
				StartHeight: ctx.BlockHeight(),
				JailedUntil: time.Unix(0, 0),
			}
			keeper.SetValidatorSigningInfo(ctx, validator.GetAddress(), signingInfo)
		}
		// if the validator is staked then add their tokens to the staked pool
		if validator.IsStaked() {
			stakedTokens = stakedTokens.Add(validator.GetTokens())
		}
	}
	// take the staked amount and create the corresponding coins object
	stakedCoins := sdk.NewCoins(sdk.NewCoin(keeper.StakeDenom(ctx), stakedTokens))
	// check if the staked pool accounts exists
	stakedPool := keeper.GetStakedPool(ctx)
	// if the stakedPool is nil
	if stakedPool == nil {
		keeper.Logger(ctx).Error(fmt.Errorf("%s module account has not been set", types.StakedPoolName).Error())
		os.Exit(1)
	}
	// add coins if not provided on genesis (there's an option to provide the coins in genesis)
	if stakedPool.GetCoins().IsZero() {
		if err := stakedPool.SetCoins(stakedCoins); err != nil {
			keeper.Logger(ctx).Error(fmt.Errorf("unable to set set coins for %s module account", types.StakedPoolName).Error())
			os.Exit(1)
		}
		supplyKeeper.SetModuleAccount(ctx, stakedPool)
	} else {
		// if it is provided in the genesis file then ensure the two are equal
		if !stakedPool.GetCoins().IsEqual(stakedCoins) {
			keeper.Logger(ctx).Error(fmt.Sprintf("%s module account total does not equal the amount in each validator account", types.StakedPoolName))
			os.Exit(1)
		}
	}
	// add coins to the total supply
	keeper.AccountKeeper.SetSupply(ctx, keeper.AccountKeeper.GetSupply(ctx).Inflate(stakedCoins))
	// don't need to run Tendermint updates if we exported
	if data.Exported {
		for _, lv := range data.PrevStateValidatorPowers {
			// set the staked validator powers from the previous state
			keeper.SetPrevStateValPower(ctx, lv.Address, lv.Power)
			validator, found := keeper.GetValidator(ctx, lv.Address)
			if !found {
				keeper.Logger(ctx).Error(fmt.Sprintf("%s validator not found from exported genesis", lv.Address))
				continue
			}
			update := validator.ABCIValidatorUpdate()
			update.Power = lv.Power // keep the next-val-set offset, use the prevState power for the first block
			res = append(res, update)
		}
	} else {
		// run tendermint updates
		res = keeper.UpdateTendermintValidators(ctx)
	}
	// update signing information from genesis state
	for addr, info := range data.SigningInfos {
		address, err := sdk.AddressFromHex(addr)
		if err != nil {
			keeper.Logger(ctx).Error(fmt.Sprintf("unable to convert address from hex in genesis signing info for addr: %s err: %v", addr, err))
			os.Exit(1)
		}
		keeper.SetValidatorSigningInfo(ctx, address, info)
	}
	// update missed block information from genesis state
	for addr, array := range data.MissedBlocks {
		address, err := sdk.AddressFromHex(addr)
		if err != nil {
			keeper.Logger(ctx).Error(fmt.Sprintf("unable to convert address from hex in genesis missed blocks for addr: %s err: %v", addr, err))
			os.Exit(1)
		}
		for _, missed := range array {
			keeper.SetValidatorMissedAt(ctx, address, missed.Index, missed.Missed)
		}
	}
	// set the params set in the keeper
	keeper.Paramstore.SetParamSet(ctx, &data.Params)
	if data.PreviousProposer != nil {
		keeper.SetPreviousProposer(ctx, data.PreviousProposer)
	}
	return res
}

// ExportGenesis returns a GenesisState for a given context and keeper
func ExportGenesis(ctx sdk.Ctx, keeper keeper.Keeper) types.GenesisState {
	params := keeper.GetParams(ctx)
	prevStateTotalPower := keeper.PrevStateValidatorsPower(ctx)
	validators := keeper.GetAllValidators(ctx)
	var prevStateValidatorPowers []types.PrevStatePowerMapping
	keeper.IterateAndExecuteOverPrevStateValsByPower(ctx, func(addr sdk.Address, power int64) (stop bool) {
		prevStateValidatorPowers = append(prevStateValidatorPowers, types.PrevStatePowerMapping{Address: addr, Power: power})
		return false
	})
	signingInfos := make(map[string]types.ValidatorSigningInfo)
	missedBlocks := make(map[string][]types.MissedBlock)
	keeper.IterateAndExecuteOverValSigningInfo(ctx, func(address sdk.Address, info types.ValidatorSigningInfo) (stop bool) {
		addrstring := address.String()
		info.Index = 0 // reset the index offset
		signingInfos[addrstring] = info
		return false
	})
	prevProposer := keeper.GetPreviousProposer(ctx)

	return types.GenesisState{
		Params:                   params,
		PrevStateTotalPower:      prevStateTotalPower,
		PrevStateValidatorPowers: prevStateValidatorPowers,
		Validators:               validators,
		Exported:                 true,
		SigningInfos:             signingInfos,
		MissedBlocks:             missedBlocks,
		PreviousProposer:         prevProposer,
	}
}

// ValidateGenesis validates the provided staking genesis state to ensure the
// expected invariants holds. (i.e. params in correct bounds, no duplicate validators)
func ValidateGenesis(data types.GenesisState) error {
	err := validateGenesisStateValidators(data.Validators, sdk.NewInt(data.Params.StakeMinimum))
	if err != nil {
		return err
	}
	err = data.Params.Validate()
	if err != nil {
		return err
	}
	downtime := data.Params.SlashFractionDowntime
	if downtime.IsNegative() || downtime.GT(sdk.OneDec()) {
		return fmt.Errorf("Slashing fraction downtime should be less than or equal to one and greater than zero, is %s", downtime.String())
	}

	dblSign := data.Params.SlashFractionDoubleSign
	if dblSign.IsNegative() || dblSign.GT(sdk.OneDec()) {
		return fmt.Errorf("Slashing fraction double sign should be less than or equal to one and greater than zero, is %s", dblSign.String())
	}

	minSign := data.Params.MinSignedPerWindow
	if minSign.IsNegative() || minSign.GT(sdk.OneDec()) {
		return fmt.Errorf("Min signed per window should be less than or equal to one and greater than zero, is %s", minSign.String())
	}

	maxEvidence := data.Params.MaxEvidenceAge
	if maxEvidence < 1*time.Minute {
		return fmt.Errorf("Max evidence age must be at least 1 minute, is %s", maxEvidence.String())
	}

	downtimeJail := data.Params.DowntimeJailDuration
	if downtimeJail < 1*time.Minute {
		return fmt.Errorf("Downtime unblond duration must be at least 1 minute, is %s", downtimeJail.String())
	}

	signedWindow := data.Params.SignedBlocksWindow
	if signedWindow < 10 {
		return fmt.Errorf("Signed blocks window must be at least 10, is %d", signedWindow)
	}
	return nil
}

func validateGenesisStateValidators(validators []types.Validator, minimumStake sdk.BigInt) (err error) {
	addrMap := make(map[string]bool, len(validators))
	for i := 0; i < len(validators); i++ {
		val := validators[i]
		strKey := val.PublicKey.RawString()
		if _, ok := addrMap[strKey]; ok {
			return fmt.Errorf("duplicate validator in genesis state: address %v", val.Address)
		}
		if val.StakedTokens.IsZero() && !val.IsUnstaked() {
			return fmt.Errorf("staked/unstaked genesis validator cannot have zero stake, validator: %v", val)
		}
		addrMap[strKey] = true
		if !val.IsUnstaked() && val.StakedTokens.LT(minimumStake) {
			return fmt.Errorf("validator has less than minimum stake: %v", val)
		}
		if err := types.ValidateServiceURL(val.ServiceURL); err != nil {
			return types.ErrInvalidServiceURL(types.ModuleName, err)
		}
		for _, chain := range val.Chains {
			err := types.ValidateNetworkIdentifier(chain)
			if err != nil {
				return err
			}
		}
	}
	return
}
