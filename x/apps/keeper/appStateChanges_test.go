package keeper

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/crypto"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/apps/types"
)

func TestAppStateChange_ValidateApplicaitonBeginUnstaking(t *testing.T) {
	tests := []struct {
		name        string
		application types.Application
		hasError    bool
		want        interface{}
	}{
		{
			name:        "validates application",
			application: getStakedApplication(),
			want:        nil,
		},
		{
			name:        "errors if application not staked",
			application: getUnstakedApplication(),
			want:        types.ErrApplicationStatus("apps"),
			hasError:    true,
		},
		{
			name:        "validates application",
			application: getStakedApplication(),
			hasError:    true,
			want:        "should not happen: application trying to begin unstaking has less than the minimum stake",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			switch tt.hasError {
			case true:
				tt.application.StakedTokens = sdk.NewInt(-1)
				_ = keeper.ValidateApplicationBeginUnstaking(context, tt.application)
			default:
				if got := keeper.ValidateApplicationBeginUnstaking(context, tt.application); !reflect.DeepEqual(got, tt.want) {
					t.Errorf("AppStateChange.ValidateApplicationBeginUnstaking() = got %v, want %v", got, tt.want)
				}
			}
		})
	}
}

func TestAppStateChange_ValidateApplicationStaking(t *testing.T) {
	tests := []struct {
		name            string
		application     types.Application
		panics          bool
		amount          sdk.BigInt
		stakedAppsCount int
		isAfterUpgrade  bool
		want            interface{}
	}{
		{
			name:            "validates application",
			stakedAppsCount: 0,
			application:     getUnstakedApplication(),
			amount:          sdk.NewInt(1000000),
			want:            nil,
		},
		{
			name:            "errors if below minimum stake",
			application:     getUnstakedApplication(),
			stakedAppsCount: 0,
			amount:          sdk.NewInt(0),
			want:            types.ErrMinimumStake("apps"),
		},
		{
			name:            "errors bank does not have enough coins",
			application:     getUnstakedApplication(),
			stakedAppsCount: 0,
			amount:          sdk.NewInt(1000000000000000000),
			want:            types.ErrNotEnoughCoins("apps"),
		},
		{
			name:            "errors if max applications hit",
			application:     getUnstakedApplication(),
			stakedAppsCount: 5,
			amount:          sdk.NewInt(1000000),
			want:            types.ErrMaxApplications("apps"),
			isAfterUpgrade:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			p := keeper.GetParams(context)
			p.MaxApplications = 5
			keeper.SetParams(context, p)
			for i := 0; i < tt.stakedAppsCount; i++ {
				pk := getRandomPubKey()
				keeper.SetStakedApplication(context, types.Application{
					Address:      sdk.Address(pk.Address()),
					PublicKey:    pk,
					Jailed:       false,
					Status:       2,
					Chains:       []string{"0021"},
					StakedTokens: sdk.NewInt(10000000),
				})
			}
			if tt.isAfterUpgrade {
				codec.UpgradeHeight = -1
			}
			addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
			sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, tt.application.Address, sdk.NewInt(100000000000))
			if got := keeper.ValidateApplicationStaking(context, tt.application, tt.amount); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AppStateChange.ValidateApplicationStaking() = got %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppStateChange_JailApplication(t *testing.T) {
	jailedApp := getStakedApplication()
	jailedApp.Jailed = true
	tests := []struct {
		name        string
		application types.Application
		hasError    bool
		want        interface{}
	}{
		{
			name:        "jails application",
			application: getStakedApplication(),
			want:        true,
		},
		{
			name:        "already jailed app ",
			application: jailedApp,
			hasError:    true,
			want:        "cannot jail already jailed application, application:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetApplication(context, tt.application)
			keeper.SetStakedApplication(context, tt.application)

			switch tt.hasError {
			case true:
				keeper.JailApplication(context, tt.application.GetAddress())
			default:
				keeper.JailApplication(context, tt.application.GetAddress())
				if got, _ := keeper.GetApplication(context, tt.application.GetAddress()); got.Jailed != tt.want {
					t.Errorf("AppStateChange.ValidateApplicationBeginUnstaking() = got %v, want %v", tt.application.Jailed, tt.want)
				}
			}

		})
	}
}

func TestAppStateChange_UnjailApplication(t *testing.T) {
	jailedApp := getStakedApplication()
	jailedApp.Jailed = true
	tests := []struct {
		name        string
		application types.Application
		hasError    bool
		want        interface{}
	}{
		{
			name:        "unjails application",
			application: jailedApp,
			want:        false,
		},
		{
			name:        "already jailed app ",
			application: getStakedApplication(),
			hasError:    true,
			want:        "cannot unjail already unjailed application, application:",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetApplication(context, tt.application)
			keeper.SetStakedApplication(context, tt.application)

			switch tt.hasError {
			case true:
				keeper.UnjailApplication(context, tt.application.GetAddress())
			default:
				keeper.UnjailApplication(context, tt.application.GetAddress())
				if got, _ := keeper.GetApplication(context, tt.application.GetAddress()); got.Jailed != tt.want {
					t.Errorf("AppStateChange.ValidateApplicationBeginUnstaking() = got %v, want %v", tt.application.Jailed, tt.want)
				}
			}

		})
	}
}

func TestAppStateChange_StakeApplication(t *testing.T) {
	tests := []struct {
		name        string
		application types.Application
		amount      sdk.BigInt
	}{
		{
			name:        "name registers apps",
			application: getUnstakedApplication(),
			amount:      sdk.NewInt(100000000000),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
			sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, tt.application.Address, sdk.NewInt(100000000000))
			_ = keeper.StakeApplication(context, tt.application, tt.amount)
			got, found := keeper.GetApplication(context, tt.application.Address)
			if !found {
				t.Errorf("AppStateChanges.RegisterApplication() = Did not register app")
			}
			if !got.StakedTokens.Equal(tt.amount.Add(sdk.NewInt(100000000000))) {
				t.Errorf("AppStateChanges.RegisterApplication() = Did not register app %v", got.StakedTokens)
			}

		})

	}
}

func TestAppStateChange_EditAndValidateStakeApplication(t *testing.T) {
	stakeAmount := sdk.NewInt(100000000000)
	accountAmount := sdk.NewInt(1000000000000).Add(stakeAmount)
	bumpStakeAmount := sdk.NewInt(1000000000000)
	newChains := []string{"0021"}
	app := getUnstakedApplication()
	app.StakedTokens = sdk.ZeroInt()
	// updatedStakeAmount
	updateStakeAmountApp := app
	updateStakeAmountApp.StakedTokens = bumpStakeAmount
	// updatedStakeAmountFail
	updateStakeAmountAppFail := app
	updateStakeAmountAppFail.StakedTokens = stakeAmount.Sub(sdk.OneInt())
	// updatedStakeAmountNotEnoughCoins
	notEnoughCoinsAccount := stakeAmount
	// updateChains
	updateChainsApp := app
	updateChainsApp.StakedTokens = stakeAmount
	updateChainsApp.Chains = newChains
	//same app no change no fail
	updateNothingApp := app
	updateNothingApp.StakedTokens = stakeAmount
	tests := []struct {
		name          string
		accountAmount sdk.BigInt
		origApp       types.Application
		amount        sdk.BigInt
		want          types.Application
		err           sdk.Error
	}{
		{
			name:          "edit stake amount of existing application",
			accountAmount: accountAmount,
			origApp:       app,
			amount:        stakeAmount,
			want:          updateStakeAmountApp,
		},
		{
			name:          "FAIL edit stake amount of existing application",
			accountAmount: accountAmount,
			origApp:       app,
			amount:        stakeAmount,
			want:          updateStakeAmountAppFail,
			err:           types.ErrMinimumEditStake("apps"),
		},
		{
			name:          "edit stake the chains of the application",
			accountAmount: accountAmount,
			origApp:       app,
			amount:        stakeAmount,
			want:          updateChainsApp,
		},
		{
			name:          "FAIL not enough coins to bump stake amount of existing application",
			accountAmount: notEnoughCoinsAccount,
			origApp:       app,
			amount:        stakeAmount,
			want:          updateStakeAmountApp,
			err:           types.ErrNotEnoughCoins("apps"),
		},
		{
			name:          "update nothing for the application",
			accountAmount: accountAmount,
			origApp:       app,
			amount:        stakeAmount,
			want:          updateNothingApp,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// test setup
			codec.UpgradeHeight = -1
			context, _, keeper := createTestInput(t, true)
			coins := sdk.NewCoins(sdk.NewCoin(keeper.StakeDenom(context), tt.accountAmount))
			err := keeper.AccountKeeper.MintCoins(context, types.StakedPoolName, coins)
			if err != nil {
				t.Fail()
			}
			err = keeper.AccountKeeper.SendCoinsFromModuleToAccount(context, types.StakedPoolName, tt.origApp.Address, coins)
			if err != nil {
				t.Fail()
			}
			err = keeper.StakeApplication(context, tt.origApp, tt.amount)
			if err != nil {
				t.Fail()
			}
			// test begins here
			err = keeper.ValidateApplicationStaking(context, tt.want, tt.want.StakedTokens)
			if err != nil {
				if tt.err.Error() != err.Error() {
					t.Fatalf("Got error %s wanted error %s", err, tt.err)
				}
				return
			}
			// edit stake
			_ = keeper.StakeApplication(context, tt.want, tt.want.StakedTokens)
			tt.want.MaxRelays = keeper.CalculateAppRelays(context, tt.want)
			tt.want.Status = sdk.Staked
			// see if the changes stuck
			got, _ := keeper.GetApplication(context, tt.origApp.Address)
			if !got.Equals(tt.want) {
				t.Fatalf("Got app %s\nWanted app %s", got.String(), tt.want.String())
			}
		})

	}
}

func TestAppStateChange_BeginUnstakingApplication(t *testing.T) {
	tests := []struct {
		name        string
		application types.Application
		want        sdk.StakeStatus
	}{
		{
			name:        "name registers apps",
			application: getStakedApplication(),
			want:        sdk.Unstaking,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			addMintedCoinsToModule(t, context, &keeper, types.StakedPoolName)
			sendFromModuleToAccount(t, context, &keeper, types.StakedPoolName, tt.application.Address, sdk.NewInt(100000000000))
			keeper.BeginUnstakingApplication(context, tt.application)
			got, found := keeper.GetApplication(context, tt.application.Address)
			if !found {
				t.Errorf("AppStateChanges.RegisterApplication() = Did not register app")
			}
			if got.Status != tt.want {
				t.Errorf("AppStateChanges.RegisterApplication() = Did not register app %v", got.StakedTokens)
			}
		})
	}
}

// Fund an account by minting tokens in a module account and sending it
// to the given account.
func fundAccount(
	t *testing.T,
	ctx *sdk.Context,
	k *Keeper,
	address sdk.Address,
	amount sdk.BigInt,
) {
	minter := types.StakedPoolName
	addMintedCoinsToModule(t, ctx, k, minter)
	sendFromModuleToAccount(t, ctx, k, minter, address, amount)
}

func transferApp(
	t *testing.T,
	ctx *sdk.Context,
	k *Keeper,
	transferFrom, transferTo crypto.PublicKey,
) sdk.Error {
	curApp, err := k.ValidateApplicationTransfer(
		ctx,
		transferFrom,
		types.MsgStake{
			PubKey: transferTo,
			Chains: nil,
			Value:  sdk.ZeroInt(),
		},
	)
	if err != nil {
		return err
	}
	k.TransferApplication(ctx, curApp, transferTo)
	return nil
}

func TestAppStateChange_Transfer(t *testing.T) {
	originalUpgradeHeight := codec.UpgradeHeight
	originalTestMode := codec.TestMode
	originalFeatKey := codec.UpgradeFeatureMap[codec.AppTransferKey]
	t.Cleanup(func() {
		codec.UpgradeHeight = originalUpgradeHeight
		codec.TestMode = originalTestMode
		codec.UpgradeFeatureMap[codec.AppTransferKey] = originalFeatKey
	})
	codec.UpgradeHeight = -1

	ctx, _, keeper := createTestInput(t, true)

	// Create four wallets and app-stake three of them
	apps := make([]types.Application, 3)
	pubKeys := make([]crypto.PublicKey, 4)
	addrs := make([]sdk.Address, 4)
	for i := range apps {
		apps[i] = createNewApplication()
		pubKeys[i] = apps[i].PublicKey
		addrs[i] = sdk.Address(pubKeys[i].Address())
		amount := sdk.NewInt(int64(i) * 10000000)
		fundAccount(t, &ctx, &keeper, apps[i].Address, amount)
		assert.Nil(t, keeper.StakeApplication(ctx, apps[i], amount))
	}
	pubKeys[3] = getRandomPubKey()
	addrs[3] = sdk.Address(pubKeys[3].Address())

	// apps[0]: staked
	// apps[1]: unstaking
	// apps[2]: staked and jailed
	// apps[3]: not an application (see pubKeys[3]) and will be used for the transfer
	keeper.BeginUnstakingApplication(ctx, apps[1])
	keeper.JailApplication(ctx, apps[2].Address)

	stakedApps := keeper.GetApplications(ctx, uint16(len(apps)*2))
	assert.Equal(t, len(apps), len(stakedApps))

	// transfer: apps[0]-->apps[3](new) fails before ugprade
	err := transferApp(t, &ctx, &keeper, pubKeys[0], pubKeys[3])
	assert.NotNil(t, err)
	assert.Equal(t, keeper.Codespace(), err.Codespace())
	assert.Equal(t, types.CodeInvalidStatus, err.Code())

	// upgrade!
	codec.UpgradeFeatureMap[codec.AppTransferKey] = -1

	// transfer: apps[0]-->apps[3](new) - success
	err = transferApp(t, &ctx, &keeper, pubKeys[0], pubKeys[3])
	assert.Nil(t, err)

	// transfer: apps[3]-->apps[2](unstaking) - fail
	err = transferApp(t, &ctx, &keeper, pubKeys[3], pubKeys[2])
	assert.NotNil(t, err)
	assert.Equal(t, keeper.Codespace(), err.Codespace())
	assert.Equal(t, types.CodeInvalidStatus, err.Code())

	// transfer: apps[3]-->apps[0](previous owner) success
	err = transferApp(t, &ctx, &keeper, pubKeys[3], pubKeys[0])
	assert.Nil(t, err)

	// transfer an unstaking app - fail
	err = transferApp(t, &ctx, &keeper, pubKeys[1], getRandomPubKey())
	assert.NotNil(t, err)
	assert.Equal(t, keeper.Codespace(), err.Codespace())
	assert.Equal(t, types.CodeInvalidStatus, err.Code())

	// transfer a new app - fail
	err = transferApp(t, &ctx, &keeper, getRandomPubKey(), getRandomPubKey())
	assert.NotNil(t, err)
	assert.Equal(t, keeper.Codespace(), err.Codespace())
	assert.Equal(t, types.CodeInvalidStatus, err.Code())

	// verify the state
	// app[0]: staked
	// app[1]: unstaking
	// app[2]: staked and jailed
	stakedApps = keeper.GetApplications(ctx, uint16(len(apps)*2))
	assert.Equal(t, 3, len(stakedApps))
	stakedApp, found := keeper.GetApplication(ctx, addrs[0])
	assert.True(t, found)
	assert.True(t, stakedApp.IsStaked())
	stakedApp, found = keeper.GetApplication(ctx, addrs[1])
	assert.True(t, found)
	assert.True(t, stakedApp.IsUnstaking())
	stakedApp, found = keeper.GetApplication(ctx, addrs[2])
	assert.True(t, found)
	assert.True(t, stakedApp.IsStaked())
	assert.True(t, stakedApp.IsJailed())

	// transfer a jailed app - success
	err = transferApp(t, &ctx, &keeper, pubKeys[2], pubKeys[3])
	assert.Nil(t, err)
	stakedApp, found = keeper.GetApplication(ctx, addrs[2])
	assert.False(t, found)
	stakedApp, found = keeper.GetApplication(ctx, addrs[3])
	assert.True(t, found)
	assert.True(t, stakedApp.IsStaked())
	assert.True(t, stakedApp.IsJailed())
}
