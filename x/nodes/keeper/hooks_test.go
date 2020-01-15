package keeper

import (
	sdk "github.com/pokt-network/posmint/types"
	"github.com/stretchr/testify/mock"
	"testing"
)

func TestHooks_BeforeValidatorRegistered(t *testing.T) {
	tests := []struct {
		name string
		args *POSHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(POSHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("BeforeValidatorRegistered", context, sdk.Address{}).Return(mock.Anything)
			keeper.BeforeValidatorRegistered(context, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_AfterValidatorRegistered(t *testing.T) {
	tests := []struct {
		name string
		args *POSHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(POSHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("AfterValidatorRegistered", context, sdk.Address{}).Return(mock.Anything)
			keeper.AfterValidatorRegistered(context, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_BeforeValidatorRemoved(t *testing.T) {
	tests := []struct {
		name string
		args *POSHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(POSHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("BeforeValidatorRemoved", context, sdk.Address{}, sdk.Address{}).Return(mock.Anything)
			keeper.BeforeValidatorRemoved(context, sdk.Address{}, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_AfterValidatorRemoved(t *testing.T) {
	tests := []struct {
		name string
		args *POSHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(POSHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("AfterValidatorRemoved", context, sdk.Address{}, sdk.Address{}).Return(mock.Anything)
			keeper.AfterValidatorRemoved(context, sdk.Address{}, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_BeforeValidatorStaked(t *testing.T) {
	tests := []struct {
		name string
		args *POSHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(POSHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("BeforeValidatorStaked", context, sdk.Address{}, sdk.Address{}).Return(mock.Anything)
			keeper.BeforeValidatorStaked(context, sdk.Address{}, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_AfterValidatorStaked(t *testing.T) {
	tests := []struct {
		name string
		args *POSHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(POSHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("AfterValidatorStaked", context, sdk.Address{}, sdk.Address{}).Return(mock.Anything)
			keeper.AfterValidatorStaked(context, sdk.Address{}, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_BeforeValidatorBeginUnstaking(t *testing.T) {
	tests := []struct {
		name string
		args *POSHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(POSHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("BeforeValidatorBeginUnstaking", context, sdk.Address{}, sdk.Address{}).Return(mock.Anything)
			keeper.BeforeValidatorBeginUnstaking(context, sdk.Address{}, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_AfterValidatorBeginUnstaking(t *testing.T) {
	tests := []struct {
		name string
		args *POSHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(POSHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("AfterValidatorBeginUnstaking", context, sdk.Address{}, sdk.Address{}).Return(mock.Anything)
			keeper.AfterValidatorBeginUnstaking(context, sdk.Address{}, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_BeforeValidatorUnstaked(t *testing.T) {
	tests := []struct {
		name string
		args *POSHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(POSHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("BeforeValidatorUnstaked", context, sdk.Address{}, sdk.Address{}).Return(mock.Anything)
			keeper.BeforeValidatorUnstaked(context, sdk.Address{}, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_AfterValidatorUnstaked(t *testing.T) {
	tests := []struct {
		name string
		args *POSHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(POSHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("AfterValidatorUnstaked", context, sdk.Address{}, sdk.Address{}).Return(mock.Anything)
			keeper.AfterValidatorUnstaked(context, sdk.Address{}, sdk.Address{})
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_BeforeValidatorSlashed(t *testing.T) {
	tests := []struct {
		name string
		args *POSHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(POSHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("BeforeValidatorSlashed", context, sdk.Address{}, sdk.NewDec(1)).Return(mock.Anything)
			keeper.BeforeValidatorSlashed(context, sdk.Address{}, sdk.NewDec(1))
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
func TestHooks_AfterValidatorSlashed(t *testing.T) {
	tests := []struct {
		name string
		args *POSHooks
		want bool
	}{
		{
			name: "calls on hook",
			args: new(POSHooks),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			_ = keeper.SetHooks(tt.args)
			tt.args.On("AfterValidatorSlashed", context, sdk.Address{}, sdk.NewDec(1)).Return(mock.Anything)
			keeper.AfterValidatorSlashed(context, sdk.Address{}, sdk.NewDec(1))
			if len(tt.args.Calls) != 1 {
				t.Errorf("hook was not called %v", len(tt.args.Calls))
			}
		})
	}
}
