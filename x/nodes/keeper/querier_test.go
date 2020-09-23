package keeper

import (
	"reflect"
	"testing"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/go-amino"
	abci "github.com/tendermint/tendermint/abci/types"
)

func Test_NewQuerier(t *testing.T) {
	type args struct {
		ctx  sdk.Context
		req  abci.RequestQuery
		path []string
		k    Keeper
	}
	context, _, keeper := createTestInput(t, true)
	jsondata, _ := amino.MarshalJSON(types.QueryValidatorsParams{
		Page:  1,
		Limit: 1,
	})
	jsondataValidatorAddr, _ := amino.MarshalJSON(getRandomValidatorAddress())

	conAddress := sdk.Address(getRandomPubKey().Address())
	jsondatasigninfo, _ := amino.MarshalJSON(types.QuerySigningInfoParams{
		Address: conAddress,
	})

	expectedValidatosPage := types.ValidatorsPage{Result: []types.Validator{}, Total: 1, Page: 1}
	jsonresponse, _ := amino.MarshalJSONIndent(expectedValidatosPage, "", "  ")
	jsonresponseSigningInfos, _ := amino.MarshalJSONIndent([]types.ValidatorSigningInfo{}, "", "  ")
	jsonresponsestakedPool, _ := amino.MarshalJSONIndent(sdk.ZeroInt(), "", "  ")
	jsonresponseInt, _ := amino.MarshalJSONIndent("0", "", "  ")
	jsonresponseForParams, _ := amino.MarshalJSONIndent(keeper.GetParams(context), "", "  ")
	tests := []struct {
		name  string
		args  args
		want  []byte
		want1 sdk.Error
	}{
		{
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

	conAddress := sdk.Address(getRandomPubKey().Address())
	context, _, keeper := createTestInput(t, true)
	jsondata, _ := amino.MarshalJSON(types.QuerySigningInfoParams{Address: conAddress})
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
	jsonresponse, _ := amino.MarshalJSONIndent(sdk.ZeroInt(), "", "  ")

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

func Test_queryValidator(t *testing.T) {
	type args struct {
		ctx sdk.Context
		req abci.RequestQuery
		k   Keeper
	}

	ranAddress := getRandomValidatorAddress()
	context, _, keeper := createTestInput(t, true)
	jsondata, _ := amino.MarshalJSON(types.QueryValidatorParams{
		Address: ranAddress,
	})
	//jsonresponse,_:= keeper.GetValidator(context,ranAddress)

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
	stakedValidator := getStakedValidator()
	validators := types.Validators{stakedValidator}
	type args struct {
		ctx sdk.Context
		req abci.RequestQuery
		k   Keeper
	}
	context, _, keeper := createTestInput(t, true)
	for _, validator := range validators {
		keeper.SetValidator(context, validator)
	}
	jsondata, _ := amino.MarshalJSON(types.QueryValidatorsParams{
		Page:  1,
		Limit: 1,
	})
	expectedValidatosPage := types.ValidatorsPage{Result: []types.Validator{stakedValidator}, Total: 1, Page: 1}
	jsonresponse, _ := amino.MarshalJSONIndent(expectedValidatosPage, "", "  ")

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

func TestQueryAccount(t *testing.T) {
	context, accs, keeper := createTestInput(t, true)
	jsondata, _ := keeper.Cdc.MarshalJSON(types.QueryAccountParams{
		Address: accs[0].GetAddress(),
	})
	req := abci.RequestQuery{
		Data: jsondata,
		Path: "account",
	}
	res, err := queryAccount(context, req, keeper)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	var acc auth.BaseAccount
	er := keeper.Cdc.UnmarshalJSON(res, &acc)
	assert.Nil(t, er)
	assert.Equal(t, accs[0], &acc)
}
