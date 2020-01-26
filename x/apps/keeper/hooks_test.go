package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestHooks_BeforeApplicationRegistered(t *testing.T) {
	tests := []struct {
		name string
		args *AppHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(AppHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("BeforeApplicationRegistered", context, sdk.Address{}).Return(mock.Anything)
			keeper.BeforeApplicationRegistered(context, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_AfterApplicationRegistered(t *testing.T) {
	tests := []struct {
		name string
		args *AppHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(AppHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("AfterApplicationRegistered", context, sdk.Address{}).Return(mock.Anything)
			keeper.AfterApplicationRegistered(context, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}

func TestHooks_BeforeApplicationRemoved(t *testing.T) {
	tests := []struct {
		name string
		args *AppHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(AppHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("BeforeApplicationRemoved", context, sdk.Address{}, sdk.Address{}).Return(mock.Anything)
			keeper.BeforeApplicationRemoved(context, sdk.Address{}, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}

func TestHooks_AfterApplicationRemoved(t *testing.T) {
	tests := []struct {
		name string
		args *AppHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(AppHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("AfterApplicationRemoved", context, sdk.Address{}, sdk.Address{}).Return(mock.Anything)
			keeper.AfterApplicationRemoved(context, sdk.Address{}, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}

func TestHooks_BeforeApplicationStaked(t *testing.T) {
	tests := []struct {
		name string
		args *AppHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(AppHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("BeforeApplicationStaked", context, sdk.Address{}, sdk.Address{}).Return(mock.Anything)
			keeper.BeforeApplicationStaked(context, sdk.Address{}, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_AfterApplicationStaked(t *testing.T) {
	tests := []struct {
		name string
		args *AppHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(AppHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("AfterApplicationStaked", context, sdk.Address{}, sdk.Address{}).Return(mock.Anything)
			keeper.AfterApplicationStaked(context, sdk.Address{}, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_BeforeApplicationBeginUnstaking(t *testing.T) {
	tests := []struct {
		name string
		args *AppHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(AppHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("BeforeApplicationBeginUnstaking", context, sdk.Address{}, sdk.Address{}).Return(mock.Anything)
			keeper.BeforeApplicationBeginUnstaking(context, sdk.Address{}, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_AfterApplicationBeginUnstaking(t *testing.T) {
	tests := []struct {
		name string
		args *AppHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(AppHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("AfterApplicationBeginUnstaking", context, sdk.Address{}, sdk.Address{}).Return(mock.Anything)
			keeper.AfterApplicationBeginUnstaking(context, sdk.Address{}, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_BeforeApplicationUnstaked(t *testing.T) {
	tests := []struct {
		name string
		args *AppHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(AppHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("BeforeApplicationUnstaked", context, sdk.Address{}, sdk.Address{}).Return(mock.Anything)
			keeper.BeforeApplicationUnstaked(context, sdk.Address{}, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_AfterApplicationUnstaked(t *testing.T) {
	tests := []struct {
		name string
		args *AppHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(AppHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("AfterApplicationUnstaked", context, sdk.Address{}, sdk.Address{}).Return(mock.Anything)
			keeper.AfterApplicationUnstaked(context, sdk.Address{}, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_BeforeApplicationSlashed(t *testing.T) {
	tests := []struct {
		name string
		args *AppHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(AppHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("BeforeApplicationSlashed", context, sdk.Address{}, sdk.NewDec(1)).Return(mock.Anything)
			keeper.BeforeApplicationSlashed(context, sdk.Address{}, sdk.NewDec(1))
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_AfterApplicationSlashed(t *testing.T) {
	tests := []struct {
		name string
		args *AppHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(AppHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("AfterApplicationSlashed", context, sdk.Address{}, sdk.NewDec(1)).Return(mock.Anything)
			keeper.AfterApplicationSlashed(context, sdk.Address{}, sdk.NewDec(1))
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
