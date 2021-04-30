package keeper

import (
	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/apps/types"
	"reflect"
	"testing"
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

func TestAppStateChange_ValidateApplicaitonStaking(t *testing.T) {
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
