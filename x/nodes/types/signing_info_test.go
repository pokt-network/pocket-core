package types

import (
	"fmt"
	"github.com/pokt-network/pocket-core/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"math/rand"
	"testing"
	"time"
)

func TestValidatorSigningInfo_String(t *testing.T) {
	type fields struct {
		Address             types.Address
		StartHeight         int64
		IndexOffset         int64
		JailedUntil         time.Time
		MissedBlocksCounter int64
		JailedBlocksCounter int64
	}
	var pub ed25519.PubKeyEd25519
	_, err := rand.Read(pub[:])
	if err != nil {
		_ = err
	}
	ca := types.Address(pub.Address())

	until := time.Now()

	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{"Validator SigningInfo String", fields{
			Address:             ca,
			StartHeight:         0,
			IndexOffset:         0,
			JailedUntil:         until,
			MissedBlocksCounter: 1,
			JailedBlocksCounter: 1,
		}, fmt.Sprintf(`Validator Signing Info:
  Address:               %s
  Start Height:          %d
  Entropy Offset:        %d
  Jailed Until:          %v
  Missed Blocks Counter: %d
  Jailed Blocks Counter: %d`,
			ca, 0, 0, until, 1, 1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := ValidatorSigningInfo{
				Address:             tt.fields.Address,
				StartHeight:         tt.fields.StartHeight,
				Index:               tt.fields.IndexOffset,
				JailedUntil:         tt.fields.JailedUntil,
				MissedBlocksCounter: tt.fields.MissedBlocksCounter,
				JailedBlocksCounter: tt.fields.JailedBlocksCounter,
			}
			if got := i.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
