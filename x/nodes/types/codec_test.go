package types

import (
	"github.com/pokt-network/posmint/codec"
	"testing"
)

func TestRegisterCodec(t *testing.T) {
	type args struct {
		cdc *codec.Codec
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}
