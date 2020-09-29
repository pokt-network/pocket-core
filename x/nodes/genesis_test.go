package nodes

import (
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/keeper"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	abci "github.com/tendermint/tendermint/abci/types"
	"reflect"
	"testing"
	"time"
)

func TestExportGenesis(t *testing.T) {
	type args struct {
		ctx    sdk.Context
		keeper keeper.Keeper
	}

	context, _, kpr := createTestInput(t, true)

	tests := []struct {
		name string
		args args
		want types.GenesisState
	}{
		{"Test Export Genesis", args{
			ctx:    context,
			keeper: kpr,
		}, getGenesisStateForTest(context, kpr, false)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ExportGenesis(tt.args.ctx, tt.args.keeper); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExportGenesis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInitGenesis(t *testing.T) {
	type args struct {
		ctx          sdk.Context
		keeper       keeper.Keeper
		supplyKeeper types.AuthKeeper
		data         types.GenesisState
	}

	context, _, kpr := createTestInput(t, true)

	validator := getStakedValidator()
	consAddress := validator.GetAddress()
	kpr.SetPreviousProposer(context, consAddress)

	tests := []struct {
		name    string
		args    args
		wantRes []abci.ValidatorUpdate
	}{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := InitGenesis(tt.args.ctx, tt.args.keeper, tt.args.supplyKeeper, tt.args.data); !reflect.DeepEqual(gotRes, tt.wantRes) {
				t.Errorf("InitGenesis() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func TestValidateGenesis(t *testing.T) {
	type args struct {
		data types.GenesisState
	}

	ctx, _, k := createTestInput(t, true)
	datafortest := getGenesisStateForTest(ctx, k, true)
	datafortest2 := getGenesisStateForTest(ctx, k, true)
	datafortest3 := getGenesisStateForTest(ctx, k, true)
	datafortest4 := getGenesisStateForTest(ctx, k, true)
	datafortest5 := getGenesisStateForTest(ctx, k, true)
	datafortest6 := getGenesisStateForTest(ctx, k, true)
	datafortest7 := getGenesisStateForTest(ctx, k, true)
	datafortest8 := getGenesisStateForTest(ctx, k, false)

	datafortest2.Params.SlashFractionDowntime = sdk.NewDec(-3)
	datafortest3.Params.SlashFractionDoubleSign = sdk.NewDec(-3)
	datafortest4.Params.MinSignedPerWindow = sdk.NewDec(-3)
	datafortest5.Params.MaxEvidenceAge = 30 * time.Second
	datafortest6.Params.DowntimeJailDuration = 30 * time.Second
	datafortest7.Params.SignedBlocksWindow = 9

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test ValidateGenesis", args{data: datafortest}, false},
		{"Test ValidateGenesis 2", args{data: datafortest2}, true},
		{"Test ValidateGenesis 3", args{data: datafortest3}, true},
		{"Test ValidateGenesis 4", args{data: datafortest4}, true},
		{"Test ValidateGenesis 5", args{data: datafortest5}, true},
		{"Test ValidateGenesis 6", args{data: datafortest6}, true},
		{"Test ValidateGenesis 7", args{data: datafortest7}, true},
		{"Test ValidateGenesis 8", args{data: datafortest8}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateGenesis(tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("ValidateGenesis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_validateGenesisStateValidators(t *testing.T) {
	type args struct {
		validators   []types.Validator
		minimumStake sdk.BigInt
	}

	ctx, _, k := createTestInput(t, true)

	testdata := getGenesisStateForTest(ctx, k, true)

	val1 := getStakedValidator()
	val1.Jailed = true
	val2 := val1

	valList := []types.Validator{val1, val2}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test ValidateGenesisStateValidators", args{
			validators:   testdata.Validators,
			minimumStake: sdk.OneInt(),
		}, false},
		{"Test ValidateGenesisStateValidators 2 duplicatedValidator", args{
			validators:   valList,
			minimumStake: sdk.OneInt(),
		}, true},
		{"Test ValidateGenesisStateValidators 3 jailed staked", args{
			validators:   []types.Validator{val1},
			minimumStake: sdk.OneInt(),
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := validateGenesisStateValidators(tt.args.validators, tt.args.minimumStake); (err != nil) != tt.wantErr {
				t.Errorf("validateGenesisStateValidators() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
