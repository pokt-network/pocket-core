package nodes

import (
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/pokt-network/posmint/codec"
	sdk "github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/rpc/client"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"reflect"
	"testing"
)

func TestQueryAccountBalance(t *testing.T) {
	type args struct {
		cdc    *codec.Codec
		tmNode client.Client
		addr   sdk.Address
		height int64
	}
	tests := []struct {
		name    string
		args    args
		want    sdk.Int
		wantErr bool
	}{
		{"TestQueryAccountBalance", args{
			cdc:    makeTestCodec(),
			tmNode: GetTestTendermintClient(), // todo these tests seem to return an error?
			addr:   getRandomValidatorAddress(),
			height: 0,
		}, sdk.ZeroInt(), true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryAccountBalance(tt.args.cdc, tt.args.tmNode, tt.args.addr, tt.args.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryAccountBalance() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryAccountBalance() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryBlock(t *testing.T) {
	type args struct {
		tmNode client.Client
		height *int64
	}
	zero := int64(0)
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{"Test QueryBlock", args{
			tmNode: GetTestTendermintClient(),
			height: &zero,
		}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryBlock(tt.args.tmNode, tt.args.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryBlock() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryBlock() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryChainHeight(t *testing.T) {
	type args struct {
		tmNode client.Client
	}
	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{"Test QueryChainHeight", args{tmNode: GetTestTendermintClient()},
			-1, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryChainHeight(tt.args.tmNode)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryChainHeight() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("QueryChainHeight() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryDAO(t *testing.T) {
	type args struct {
		cdc    *codec.Codec
		tmNode client.Client
		height int64
	}
	tests := []struct {
		name         string
		args         args
		wantDaoCoins sdk.Int
		wantErr      bool
	}{
		{"Test QueryDAo", args{
			cdc:    makeTestCodec(),
			tmNode: GetTestTendermintClient(),
			height: 0,
		}, sdk.Int{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDaoCoins, err := QueryDAO(tt.args.cdc, tt.args.tmNode, tt.args.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryDAO() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotDaoCoins, tt.wantDaoCoins) {
				t.Errorf("QueryDAO() gotDaoCoins = %v, want %v", gotDaoCoins, tt.wantDaoCoins)
			}
		})
	}
}

func TestQueryNodeStatus(t *testing.T) {
	type args struct {
		tmNode client.Client
	}
	tests := []struct {
		name    string
		args    args
		want    *ctypes.ResultStatus
		wantErr bool
	}{
		{"Test Query NodeStatus", args{tmNode: GetTestTendermintClient()},
			nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryNodeStatus(tt.args.tmNode)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryNodeStatus() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryNodeStatus() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryPOSParams(t *testing.T) {
	type args struct {
		cdc    *codec.Codec
		tmNode client.Client
		height int64
	}
	tests := []struct {
		name    string
		args    args
		want    types.Params
		wantErr bool
	}{
		{"Test QueryPOSParams", args{
			cdc:    makeTestCodec(),
			tmNode: GetTestTendermintClient(),
			height: 0,
		}, types.Params{
			UnstakingTime:           0,
			MaxValidators:           0,
			StakeDenom:              "",
			StakeMinimum:            0,
			DAOAllocation:           0,
			SessionBlockFrequency:   0,
			ProposerAllocation:      0,
			MaxEvidenceAge:          0,
			SignedBlocksWindow:      0,
			MinSignedPerWindow:      sdk.Dec{},
			DowntimeJailDuration:    0,
			SlashFractionDoubleSign: sdk.Dec{},
			SlashFractionDowntime:   sdk.Dec{},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryPOSParams(tt.args.cdc, tt.args.tmNode, tt.args.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryPOSParams() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryPOSParams() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuerySigningInfo(t *testing.T) {
	type args struct {
		cdc      *codec.Codec
		tmNode   client.Client
		height   int64
		consAddr sdk.Address
	}
	tests := []struct {
		name    string
		args    args
		want    types.ValidatorSigningInfo
		wantErr bool
	}{
		{"Test QuerySigningInfo", args{
			cdc:      makeTestCodec(),
			tmNode:   GetTestTendermintClient(),
			height:   0,
			consAddr: sdk.Address(getRandomPubKey().Address()),
		}, types.ValidatorSigningInfo{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QuerySigningInfo(tt.args.cdc, tt.args.tmNode, tt.args.height, tt.args.consAddr)
			if (err != nil) != tt.wantErr {
				t.Errorf("QuerySigningInfo() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QuerySigningInfo() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryStakedValidators(t *testing.T) {
	type args struct {
		cdc    *codec.Codec
		tmNode client.Client
		height int64
	}
	tests := []struct {
		name    string
		args    args
		want    types.Validators
		wantErr bool
	}{
		{"Test QueryStakedValidators", args{
			cdc:    makeTestCodec(),
			tmNode: GetTestTendermintClient(),
			height: 0,
		}, types.Validators{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryStakedValidators(tt.args.cdc, tt.args.tmNode, tt.args.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryStakedValidators() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryStakedValidators() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQuerySupply(t *testing.T) {
	type args struct {
		cdc    *codec.Codec
		tmNode client.Client
		height int64
	}
	tests := []struct {
		name              string
		args              args
		wantStakedCoins   sdk.Int
		wantUnstakedCoins sdk.Int
		wantErr           bool
	}{
		{"Test QuerySupply", args{
			cdc:    makeTestCodec(),
			tmNode: GetTestTendermintClient(),
			height: 0,
		}, sdk.Int{}, sdk.Int{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotStakedCoins, gotUnstakedCoins, err := QuerySupply(tt.args.cdc, tt.args.tmNode, tt.args.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("QuerySupply() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(gotStakedCoins, tt.wantStakedCoins) {
				t.Errorf("QuerySupply() gotStakedCoins = %v, want %v", gotStakedCoins, tt.wantStakedCoins)
			}
			if !reflect.DeepEqual(gotUnstakedCoins, tt.wantUnstakedCoins) {
				t.Errorf("QuerySupply() gotUnstakedCoins = %v, want %v", gotUnstakedCoins, tt.wantUnstakedCoins)
			}
		})
	}
}

func TestQueryTransaction(t *testing.T) {
	type args struct {
		tmNode client.Client
		hash   string
	}
	tests := []struct {
		name    string
		args    args
		want    *ctypes.ResultTx
		wantErr bool
	}{
		{"Test QueryTransaction", args{
			tmNode: GetTestTendermintClient(),
			hash:   "",
		}, nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryTransaction(tt.args.tmNode, tt.args.hash)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryTransaction() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryUnstakedValidators(t *testing.T) {
	type args struct {
		cdc    *codec.Codec
		tmNode client.Client
		height int64
	}
	tests := []struct {
		name    string
		args    args
		want    types.Validators
		wantErr bool
	}{
		{"Test QueryUnstakedValidators", args{
			cdc:    makeTestCodec(),
			tmNode: GetTestTendermintClient(),
			height: 0,
		}, types.Validators{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryUnstakedValidators(tt.args.cdc, tt.args.tmNode, tt.args.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryUnstakedValidators() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryUnstakedValidators() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryUnstakingValidators(t *testing.T) {
	type args struct {
		cdc    *codec.Codec
		tmNode client.Client
		height int64
	}
	tests := []struct {
		name    string
		args    args
		want    types.Validators
		wantErr bool
	}{
		{"Test QueryUnstakingValidators", args{
			cdc:    makeTestCodec(),
			tmNode: GetTestTendermintClient(),
			height: 0,
		}, types.Validators{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryUnstakingValidators(tt.args.cdc, tt.args.tmNode, tt.args.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryUnstakingValidators() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryUnstakingValidators() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryValidator(t *testing.T) {
	type args struct {
		cdc    *codec.Codec
		tmNode client.Client
		addr   sdk.Address
		height int64
	}
	tests := []struct {
		name    string
		args    args
		want    types.Validator
		wantErr bool
	}{
		{"Test QueryValidator", args{
			cdc:    makeTestCodec(),
			tmNode: GetTestTendermintClient(),
			addr:   getRandomValidatorAddress(),
			height: 0,
		}, types.Validator{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryValidator(tt.args.cdc, tt.args.tmNode, tt.args.addr, tt.args.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryValidator() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryValidator() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestQueryValidators(t *testing.T) {
	type args struct {
		cdc    *codec.Codec
		tmNode client.Client
		height int64
	}
	tests := []struct {
		name    string
		args    args
		want    types.Validators
		wantErr bool
	}{
		{"Test QueryValidators", args{
			cdc:    makeTestCodec(),
			tmNode: GetTestTendermintClient(),
			height: 0,
		}, types.Validators{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := QueryValidators(tt.args.cdc, tt.args.tmNode, tt.args.height)
			if (err != nil) != tt.wantErr {
				t.Errorf("QueryValidators() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("QueryValidators() got = %v, want %v", got, tt.want)
			}
		})
	}
}
