package nodes

import (
	"encoding/json"
	"reflect"
	"testing"

	"github.com/pokt-network/pocket-core/codec"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/keeper"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

func TestAppModuleBasic_DefaultGenesis(t *testing.T) {
	tests := []struct {
		name string
		want json.RawMessage
	}{
		{"Test DefaultGenesis", types.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := AppModuleBasic{}
			if got := ap.DefaultGenesis(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DefaultGenesis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppModuleBasic_Name(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := AppModuleBasic{}
			if got := ap.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppModuleBasic_RegisterCodec(t *testing.T) {
	type args struct {
		cdc *codec.Codec
	}
	tests := []struct {
		name string
		args args
	}{
		{"Test RegisterCodec", args{cdc: makeTestCodec()}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := AppModuleBasic{}
			ap.RegisterCodec(tt.args.cdc)
		})
	}
}

func TestAppModuleBasic_ValidateGenesis(t *testing.T) {
	type args struct {
		bz json.RawMessage
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"Test ValidateGenesis", args{bz: types.ModuleCdc.MustMarshalJSON(types.DefaultGenesisState())}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := AppModuleBasic{}
			if err := ap.ValidateGenesis(tt.args.bz); (err != nil) != tt.wantErr {
				t.Errorf("ValidateGenesis() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAppModule_BeginBlock(t *testing.T) {
	type fields struct {
		AppModuleBasic AppModuleBasic
		keeper         keeper.Keeper
		accountKeeper  types.AuthKeeper
		supplyKeeper   types.AuthKeeper
	}
	type args struct {
		ctx sdk.Context
		req abci.RequestBeginBlock
	}

	ctx, _, k := createTestInput(t, true)
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test BeginBlock", fields{
			AppModuleBasic: AppModuleBasic{},
			keeper:         k,
			accountKeeper:  nil,
			supplyKeeper:   nil,
		}, args{
			ctx: ctx,
			req: abci.RequestBeginBlock{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := AppModule{
				AppModuleBasic: tt.fields.AppModuleBasic,
				keeper:         tt.fields.keeper,
			}
			am.BeginBlock(tt.args.ctx, tt.args.req)
		})
	}
}

func TestAppModule_EndBlock(t *testing.T) {
	type fields struct {
		AppModuleBasic AppModuleBasic
		keeper         keeper.Keeper
		accountKeeper  types.AuthKeeper
		supplyKeeper   types.AuthKeeper
	}
	type args struct {
		ctx sdk.Context
		in1 abci.RequestEndBlock
	}

	ctx, _, k := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
		want   []abci.ValidatorUpdate
	}{
		{"Test EndBlock", fields{
			AppModuleBasic: AppModuleBasic{},
			keeper:         k,
			accountKeeper:  nil,
			supplyKeeper:   nil,
		}, args{
			ctx: ctx,
			in1: abci.RequestEndBlock{},
		}, []abci.ValidatorUpdate{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := AppModule{
				AppModuleBasic: tt.fields.AppModuleBasic,
				keeper:         tt.fields.keeper,
			}
			if got := am.EndBlock(tt.args.ctx, tt.args.in1); !reflect.DeepEqual(len(got), len(tt.want)) {
				t.Errorf("EndBlock() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppModule_ExportGenesis(t *testing.T) {
	type fields struct {
		AppModuleBasic AppModuleBasic
		keeper         keeper.Keeper
		accountKeeper  types.AuthKeeper
		supplyKeeper   types.AuthKeeper
	}
	context, _, k := createTestInput(t, true)

	k.SetPreviousProposer(context, sdk.GetAddress(getRandomPubKey()))
	type args struct {
		ctx sdk.Context
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   json.RawMessage
	}{
		{"Test Export Genesis", fields{
			AppModuleBasic: AppModuleBasic{},
			keeper:         k,
			accountKeeper:  nil,
			supplyKeeper:   nil,
		}, args{ctx: context}, types.ModuleCdc.MustMarshalJSON(ExportGenesis(context, k))},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := AppModule{
				AppModuleBasic: tt.fields.AppModuleBasic,
				keeper:         tt.fields.keeper,
			}
			if got := am.ExportGenesis(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExportGenesis() = %s, want %s", got, tt.want)
			}
		})
	}
}

func TestAppModule_InitGenesis(t *testing.T) {
	type fields struct {
		AppModuleBasic AppModuleBasic
		keeper         keeper.Keeper
	}
	type args struct {
		ctx  sdk.Context
		data json.RawMessage
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   []abci.ValidatorUpdate
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := AppModule{
				AppModuleBasic: tt.fields.AppModuleBasic,
				keeper:         tt.fields.keeper,
			}
			if got := am.InitGenesis(tt.args.ctx, tt.args.data); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InitGenesis() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppModule_Name(t *testing.T) {
	type fields struct {
		AppModuleBasic AppModuleBasic
		keeper         keeper.Keeper
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := AppModule{
				AppModuleBasic: tt.fields.AppModuleBasic,
				keeper:         tt.fields.keeper,
			}
			if got := ap.Name(); got != tt.want {
				t.Errorf("Name() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppModule_NewHandler(t *testing.T) {
	type fields struct {
		AppModuleBasic AppModuleBasic
		keeper         keeper.Keeper
		accountKeeper  types.AuthKeeper
		supplyKeeper   types.AuthKeeper
	}

	_, _, k := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		want   sdk.Handler
	}{
		{"Test NewHandler", fields{
			AppModuleBasic: AppModuleBasic{},
			keeper:         k,
			accountKeeper:  nil,
			supplyKeeper:   nil,
		}, NewHandler(k)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := AppModule{
				AppModuleBasic: tt.fields.AppModuleBasic,
				keeper:         tt.fields.keeper,
			}
			am.NewHandler()
		})
	}
}

func TestAppModule_NewQuerierHandler(t *testing.T) {
	type fields struct {
		AppModuleBasic AppModuleBasic
		keeper         keeper.Keeper
		accountKeeper  types.AuthKeeper
		supplyKeeper   types.AuthKeeper
	}
	tests := []struct {
		name   string
		fields fields
		want   sdk.Querier
	}{
		{"Test Querier Handler", fields{
			AppModuleBasic: AppModuleBasic{},
			keeper:         keeper.Keeper{},
			accountKeeper:  nil,
			supplyKeeper:   nil,
		}, keeper.NewQuerier(keeper.Keeper{})},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			am := AppModule{
				AppModuleBasic: tt.fields.AppModuleBasic,
				keeper:         tt.fields.keeper,
			}
			am.NewQuerierHandler()
		})
	}
}

func TestAppModule_QuerierRoute(t *testing.T) {
	type fields struct {
		AppModuleBasic AppModuleBasic
		keeper         keeper.Keeper
		accountKeeper  types.AuthKeeper
		supplyKeeper   types.AuthKeeper
	}

	_, _, k := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Test QuerierRoute", fields{
			AppModuleBasic: AppModuleBasic{},
			keeper:         k,
			accountKeeper:  nil,
			supplyKeeper:   nil,
		}, types.ModuleName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := AppModule{
				AppModuleBasic: tt.fields.AppModuleBasic,
				keeper:         tt.fields.keeper,
			}
			if got := ap.QuerierRoute(); got != tt.want {
				t.Errorf("QuerierRoute() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAppModule_Route(t *testing.T) {
	type fields struct {
		AppModuleBasic AppModuleBasic
		keeper         keeper.Keeper
		accountKeeper  types.AuthKeeper
		supplyKeeper   types.AuthKeeper
	}
	_, _, keeper := createTestInput(t, true)
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"test Route", fields{
			AppModuleBasic: AppModuleBasic{},
			keeper:         keeper,
			accountKeeper:  nil,
			supplyKeeper:   nil,
		}, types.ModuleName},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ap := AppModule{
				AppModuleBasic: tt.fields.AppModuleBasic,
				keeper:         tt.fields.keeper,
			}
			if got := ap.Route(); got != tt.want {
				t.Errorf("Route() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewAppModule(t *testing.T) {
	type args struct {
		keeper keeper.Keeper
	}
	tests := []struct {
		name string
		args args
		want AppModule
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewAppModule(tt.args.keeper); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewAppModule() = %v, want %v", got, tt.want)
			}
		})
	}
}
