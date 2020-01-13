package types

import (
	"fmt"
	"github.com/pokt-network/posmint/types"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"math/rand"
	"testing"
	"time"
)

func TestValidatorSigningInfo_String(t *testing.T) {
	type fields struct {
		Address             types.ConsAddress
		StartHeight         int64
		IndexOffset         int64
		JailedUntil         time.Time
		Tombstoned          bool
		MissedBlocksCounter int64
	}
	var pub ed25519.PubKeyEd25519
	rand.Read(pub[:])
	ca := types.ConsAddress(pub.Address())

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
			Tombstoned:          false,
			MissedBlocksCounter: 1,
		}, fmt.Sprintf(`Validator Signing Info:
  Address:               %s
  Start Height:          %d
  Entropy Offset:        %d
  Jailed Until:          %v
  Tombstoned:            %t
  Missed Blocks Counter: %d`,
			ca, 0, 0, until,
			false, 1)},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := ValidatorSigningInfo{
				Address:             tt.fields.Address,
				StartHeight:         tt.fields.StartHeight,
				IndexOffset:         tt.fields.IndexOffset,
				JailedUntil:         tt.fields.JailedUntil,
				Tombstoned:          tt.fields.Tombstoned,
				MissedBlocksCounter: tt.fields.MissedBlocksCounter,
			}
			if got := i.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}
