package keeper

import (
	"github.com/pokt-network/pocket-core/x/nodes/types"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
	"reflect"
	"testing"
)

func Test_NewQuerier(t *testing.T) {
	type args struct {
		ctx  sdk.Context
		req  abci.RequestQuery
		path []string
		k    Keeper
	}
	context, _, keeper := createTestInput(t, true)
	jsondata, _ := amino.MarshalJSON(types.QueryUnstakingValidatorsParams{
		Page:  1,
		Limit: 1,
	})
	jsondataValidatorAddr, _ := amino.MarshalJSON(getRandomValidatorAddress())

	conAddress := sdk.ConsAddress(getRandomPubKey().Address())
	jsondatasigninfo, _ := amino.MarshalJSON(types.QuerySigningInfoParams{
		ConsAddress: conAddress,
	})

	jsonresponse, _ := amino.MarshalJSONIndent([]types.Validator{}, "", "  ")
	jsonresponseSigningInfos, _ := amino.MarshalJSONIndent([]types.ValidatorSigningInfo{}, "", "  ")
	jsonresponsestakedPool, _ := amino.MarshalJSONIndent(types.StakingPool(types.NewPool(sdk.ZeroInt())), "", "  ")
	jsonresponseDAO, _ := amino.MarshalJSONIndent(types.NewPool(sdk.ZeroInt()), "", "  ")
	jsonresponseInt, _ := amino.MarshalJSONIndent("0", "", "  ")
	jsonresponseForParams, _ := amino.MarshalJSONIndent(keeper.GetParams(context), "", "  ")
	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 sdk.Error
	}{
		{
			name: "Test queryUnstakingValidators",
			args: args{
				ctx:  context,
				req:  abci.RequestQuery{Data: jsondata, Path: "unstaking_validators"},
				path: []string{types.QueryUnstakingValidators},
				k:    keeper,
			},
			want:  jsonresponse,
			want1: nil,
		},
		{
			name: "Test queryUnstakedValidators",
			args: args{
				ctx:  context,
				req:  abci.RequestQuery{Data: jsondata, Path: "unstaking_validators"},
				path: []string{types.QueryUnstakedValidators},
				k:    keeper,
			},
			want:  jsonresponse,
			want1: nil,
		}, {
			name: "Test QuerySigningInfo",
			args: args{
				ctx:  context,
				req:  abci.RequestQuery{Data: jsondatasigninfo, Path: "signingInfo"},
				path: []string{types.QuerySigningInfo},
				k:    keeper,
			},
			want:  nil,
			want1: types.ErrNoSigningInfoFound("pos", conAddress),
		}, {
			name: "Test QuerySigningInfos",
			args: args{
				ctx:  context,
				req:  abci.RequestQuery{Data: jsondata, Path: "signingInfos"},
				path: []string{types.QuerySigningInfos},
				k:    keeper,
			},
			want:  jsonresponseSigningInfos,
			want1: nil,
		}, {
			name: "Test QueryStakedPool",
			args: args{
				ctx:  context,
				req:  abci.RequestQuery{Data: jsondata, Path: "stakedPool"},
				path: []string{types.QueryStakedPool},
				k:    keeper,
			},
			want:  jsonresponsestakedPool,
			want1: nil,
		}, {
			name: "Test QueryUnstakedPool",
			args: args{
				ctx:  context,
				req:  abci.RequestQuery{Data: jsondata, Path: "unstakedPool"},
				path: []string{types.QueryUnstakedPool},
				k:    keeper,
			},
			want:  jsonresponsestakedPool,
			want1: nil,
		}, {
			name: "Test QueryDAO",
			args: args{
				ctx:  context,
				req:  abci.RequestQuery{Data: jsondata, Path: "dao"},
				path: []string{types.QueryDAO},
				k:    keeper,
			},
			want:  jsonresponseDAO,
			want1: nil,
		}, {
			name: "Test QueryAccountBalance",
			args: args{
				ctx:  context,
				req:  abci.RequestQuery{Data: jsondataValidatorAddr, Path: "account_balance"},
				path: []string{types.QueryAccountBalance},
				k:    keeper,
			},
			want:  jsonresponseInt,
			want1: nil,
		},
		{
			name: "Test queryStakedValidators",
			args: args{
				ctx:  context,
				req:  abci.RequestQuery{Data: jsondata, Path: "unstaking_validators"},
				path: []string{types.QueryStakedValidators},
				k:    keeper,
			},
			want:  jsonresponse,
			want1: nil,
		},
		{
			name: "Test queryParams",
			args: args{
				ctx:  context,
				req:  abci.RequestQuery{Data: jsondata, Path: "unstaking_validators"},
				path: []string{types.QueryParameters},
				k:    keeper,
			},
			want:  jsonresponseForParams,
			want1: nil,
		},
		{
			name: "Test queryValidators",
			args: args{
				ctx:  context,
				req:  abci.RequestQuery{Data: jsondata, Path: "unstaking_validators"},
				path: []string{types.QueryValidators},
				k:    keeper,
			},
			want:  jsonresponse,
			want1: nil,
		},
		{
			name: "Test query application",
			args: args{
				ctx:  context,
				req:  abci.RequestQuery{Data: jsondata, Path: "unstaking_validators"},
				path: []string{types.QueryValidator},
				k:    keeper,
			},
			want:  []byte(nil),
			want1: types.ErrNoValidatorFound(types.DefaultCodespace),
		},
		{
			name: "Test default querier",
			args: args{
				ctx:  context,
				req:  abci.RequestQuery{Data: jsondata, Path: "unstaking_validators"},
				path: []string{"query"},
				k:    keeper,
			},
			want:  []byte(nil),
			want1: sdk.ErrUnknownRequest("unknown staking query endpoint"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fn := NewQuerier(tt.args.k)
			got, got1 := fn(tt.args.ctx, tt.args.path, tt.args.req)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryUnstakingValidators() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("queryUnstakingValidators() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
func Test_queryAccountBalance(t *testing.T) {
	type args struct {
		ctx sdk.Context
		req abci.RequestQuery
		k   Keeper
	}

	context, _, keeper := createTestInput(t, true)
	jsondata, _ := amino.MarshalJSON(getRandomValidatorAddress())
	jsonresponse, _ := amino.MarshalJSONIndent("0", "", "  ")

	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 sdk.Error
	}{
		{"Test QueryAccBalance", args{
			ctx: context,
			req: abci.RequestQuery{
				Path: "account_balance",
				Data: jsondata,
			},
			k: keeper,
		}, jsonresponse,
			nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := queryAccountBalance(tt.args.ctx, tt.args.req, tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryAccountBalance() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("queryAccountBalance() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_queryDAO(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   Keeper
	}
	context, _, keeper := createTestInput(t, true)
	jsonresponse, _ := amino.MarshalJSONIndent(types.NewPool(sdk.ZeroInt()), "", "  ")

	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 sdk.Error
	}{
		{"Test QueryDao", args{
			ctx: context,
			k:   keeper,
		}, jsonresponse, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := queryDAO(tt.args.ctx, tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryDAO() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("queryDAO() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_queryParameters(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   Keeper
	}

	context, _, keeper := createTestInput(t, true)
	jsonresponse, _ := amino.MarshalJSONIndent(keeper.GetParams(context), "", "  ")

	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 sdk.Error
	}{
		{"Test Queryparameters", args{
			ctx: context,
			k:   keeper,
		}, jsonresponse, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := queryParameters(tt.args.ctx, tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryParameters() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("queryParameters() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_querySigningInfo(t *testing.T) {
	type args struct {
		ctx sdk.Context
		req abci.RequestQuery
		k   Keeper
	}

	conAddress := sdk.ConsAddress(getRandomPubKey().Address())
	context, _, keeper := createTestInput(t, true)
	jsondata, _ := amino.MarshalJSON(types.QuerySigningInfoParams{ConsAddress: conAddress})
	var jsonresponse []byte

	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 sdk.Error
	}{
		{"Test QuerySigningInfo - No SigningInfor for random key", args{
			ctx: context,
			req: abci.RequestQuery{
				Data: jsondata,
				Path: "signingInfo",
			},
			k: keeper,
		}, jsonresponse, types.ErrNoSigningInfoFound("pos", conAddress)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := querySigningInfo(tt.args.ctx, tt.args.req, tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("querySigningInfo() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("querySigningInfo() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_querySigningInfos(t *testing.T) {
	type args struct {
		ctx sdk.Context
		req abci.RequestQuery
		k   Keeper
	}

	context, _, keeper := createTestInput(t, true)
	jsondata, _ := amino.MarshalJSON(types.QuerySigningInfosParams{
		Page:  1,
		Limit: 1,
	})
	jsonresponse, _ := amino.MarshalJSONIndent([]types.ValidatorSigningInfo{}, "", "  ")

	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 sdk.Error
	}{
		{"Test QuerySigningInfos", args{
			ctx: context,
			req: abci.RequestQuery{
				Data: jsondata,
				Path: "signingInfos",
			},
			k: keeper,
		}, jsonresponse, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := querySigningInfos(tt.args.ctx, tt.args.req, tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("querySigningInfos() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("querySigningInfos() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_queryStakedPool(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   Keeper
	}

	context, _, keeper := createTestInput(t, true)
	jsonresponse, _ := amino.MarshalJSONIndent(types.StakingPool(types.NewPool(sdk.ZeroInt())), "", "  ")

	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 sdk.Error
	}{
		{"test QueryStakedPool", args{
			ctx: context,
			k:   keeper,
		}, jsonresponse, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := queryStakedPool(tt.args.ctx, tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryStakedPool() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("queryStakedPool() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_queryStakedValidators(t *testing.T) {
	type args struct {
		ctx sdk.Context
		req abci.RequestQuery
		k   Keeper
	}

	context, _, keeper := createTestInput(t, true)
	jsondata, _ := amino.MarshalJSON(types.QueryStakedValidatorsParams{
		Page:  1,
		Limit: 1,
	})
	jsonresponse, _ := amino.MarshalJSONIndent([]types.Validator{}, "", "  ")

	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 sdk.Error
	}{
		{"Test queryStakedValidators", args{
			ctx: context,
			req: abci.RequestQuery{
				Data: jsondata,
				Path: "staked_validators",
			},
			k: keeper,
		}, jsonresponse, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := queryStakedValidators(tt.args.ctx, tt.args.req, tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryStakedValidators() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("queryStakedValidators() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_queryUnstakedPool(t *testing.T) {
	type args struct {
		ctx sdk.Context
		k   Keeper
	}

	context, _, keeper := createTestInput(t, true)
	jsonresponse, _ := amino.MarshalJSONIndent(types.StakingPool(types.NewPool(sdk.ZeroInt())), "", "  ")

	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 sdk.Error
	}{
		{"Test QueryUnstakedPool", args{
			ctx: context,
			k:   keeper,
		}, jsonresponse, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := queryUnstakedPool(tt.args.ctx, tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryUnstakedPool() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("queryUnstakedPool() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_queryUnstakedValidators(t *testing.T) {
	type args struct {
		ctx sdk.Context
		req abci.RequestQuery
		k   Keeper
	}

	context, _, keeper := createTestInput(t, true)
	jsondata, _ := amino.MarshalJSON(types.QueryValidatorsParams{
		Page:  1,
		Limit: 1,
	})
	jsonresponse, _ := amino.MarshalJSONIndent([]types.Validator{}, "", "  ")

	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 sdk.Error
	}{
		{"Test queryUnstakedValidators", args{
			ctx: context,
			req: abci.RequestQuery{
				Data: jsondata,
				Path: "unstaked_validators",
			},
			k: keeper,
		}, jsonresponse, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := queryUnstakedValidators(tt.args.ctx, tt.args.req, tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryUnstakedValidators() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("queryUnstakedValidators() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_queryUnstakingValidators(t *testing.T) {
	type args struct {
		ctx sdk.Context
		req abci.RequestQuery
		k   Keeper
	}

	context, _, keeper := createTestInput(t, true)
	jsondata, _ := amino.MarshalJSON(types.QueryUnstakingValidatorsParams{
		Page:  1,
		Limit: 1,
	})
	jsonresponse, _ := amino.MarshalJSONIndent([]types.Validator{}, "", "  ")

	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 sdk.Error
	}{
		{"Test queryUnstakinValidators", args{
			ctx: context,
			req: abci.RequestQuery{Data: jsondata, Path: "unstaking_validators"},
			k:   keeper,
		}, jsonresponse, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := queryUnstakingValidators(tt.args.ctx, tt.args.req, tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryUnstakingValidators() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("queryUnstakingValidators() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_queryValidator(t *testing.T) {
	type args struct {
		ctx sdk.Context
		req abci.RequestQuery
		k   Keeper
	}

	ranValAddress := getRandomValidatorAddress()
	context, _, keeper := createTestInput(t, true)
	jsondata, _ := amino.MarshalJSON(types.QueryValidatorParams{
		Address: ranValAddress,
	})
	//jsonresponse,_:= keeper.GetValidator(context,ranValAddress)

	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 sdk.Error
	}{
		{"Test Query Validator", args{
			ctx: context,
			req: abci.RequestQuery{
				Data: jsondata,
				Path: "validator",
			},
			k: keeper,
		}, nil, types.ErrNoValidatorFound(types.DefaultCodespace)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := queryValidator(tt.args.ctx, tt.args.req, tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryValidator() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("queryValidator() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_queryValidators(t *testing.T) {
	type args struct {
		ctx sdk.Context
		req abci.RequestQuery
		k   Keeper
	}

	context, _, keeper := createTestInput(t, true)
	jsondata, _ := amino.MarshalJSON(types.QueryValidatorsParams{
		Page:  1,
		Limit: 1,
	})
	jsonresponse, _ := amino.MarshalJSONIndent([]types.Validator{}, "", "  ")

	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 sdk.Error
	}{
		{"Test Query Validators", args{
			ctx: context,
			req: abci.RequestQuery{
				Data: jsondata,
				Path: "validators",
			},
			k: keeper,
		}, jsonresponse, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := queryValidators(tt.args.ctx, tt.args.req, tt.args.k)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("queryValidators() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("queryValidators() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
