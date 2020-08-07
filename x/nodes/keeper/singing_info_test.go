package keeper

import (
	"testing"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/nodes/types"
	"github.com/stretchr/testify/assert"
)

func TestGetMissedArray(t *testing.T) {
	validator := getStakedValidator()
	consAddr := validator.GetAddress()

	tests := []struct {
		name     string
		expected bool
		address  sdk.Address
	}{
		{
			name:     "gets missed block array",
			address:  consAddr,
			expected: true,
		},
		{
			name:     "gets missed block array",
			address:  consAddr,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetValidatorMissedAt(context, test.address, 1, test.expected)
			missed := keeper.valMissedAt(context, test.address, 1)
			assert.Equal(t, missed, test.expected, "found does not match")
		})
	}
}

func TestClearMissedArray(t *testing.T) {
	validator := getStakedValidator()
	consAddr := validator.GetAddress()

	tests := []struct {
		name     string
		expected bool
		address  sdk.Address
	}{
		{
			name:     "gets missed block array",
			address:  consAddr,
			expected: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			context, _, keeper := createTestInput(t, true)
			keeper.SetValidatorMissedAt(context, test.address, 1, true)
			keeper.clearValidatorMissed(context, test.address)
			missed := keeper.valMissedAt(context, test.address, 1)
			assert.Equal(t, missed, test.expected, "found does not match")
		})
	}
}

func TestKeeper_IterateAndExecuteOverMissedArray(t *testing.T) {
	type fields struct {
		Keeper Keeper
	}
	type args struct {
		ctx     sdk.Context
		address sdk.Address
		handler func(index int64, missed bool) (stop bool)
	}

	context, _, keeper := createTestInput(t, true)

	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{"Test IterateAndExecuteOverMissedArray", fields{Keeper: keeper},
			args{
				ctx:     context,
				address: sdk.Address(getRandomPubKey().Address()),
				handler: func(index int64, missed bool) (stop bool) {
					localMissedBlocks := []types.MissedBlock{}

					_ = append(localMissedBlocks, types.MissedBlock{Index: index, Missed: missed})
					return false
				},
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			k := tt.fields.Keeper

			k.IterateAndExecuteOverMissedArray(tt.args.ctx, tt.args.address, tt.args.handler)

		})
	}
}
