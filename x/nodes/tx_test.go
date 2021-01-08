package nodes

import (
	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/crypto/keys"
	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/tendermint/tendermint/rpc/client"
	"reflect"
	"testing"
)

func TestRawTx(t *testing.T) {
	type args struct {
		cdc      *codec.Codec
		tmNode   client.Client
		fromAddr sdk.Address
		txBytes  []byte
	}
	tests := []struct {
		name    string
		args    args
		want    sdk.TxResponse
		wantErr bool
	}{
		{"Test RawTx", args{
			cdc:      makeTestCodec(),
			tmNode:   GetTestTendermintClient(),
			fromAddr: getRandomValidatorAddress(),
			txBytes:  []byte{0x51, 0x41, 0x33},
		}, sdk.TxResponse{}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RawTx(tt.args.cdc, tt.args.tmNode, tt.args.fromAddr, tt.args.txBytes)
			if (err != nil) != tt.wantErr {
				t.Errorf("RawTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RawTx() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSend(t *testing.T) {
	type args struct {
		cdc        *codec.Codec
		tmNode     client.Client
		keybase    keys.Keybase
		fromAddr   sdk.Address
		toAddr     sdk.Address
		passphrase string
		amount     sdk.BigInt
	}
	tests := []struct {
		name    string
		args    args
		want    *sdk.TxResponse
		wantErr bool
	}{
		{"Test Send", args{
			cdc:        makeTestCodec(),
			tmNode:     GetTestTendermintClient(),
			keybase:    nil,
			fromAddr:   nil,
			toAddr:     nil,
			passphrase: "",
			amount:     sdk.BigInt{},
		}, nil,
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Send(tt.args.cdc, tt.args.tmNode, tt.args.keybase, tt.args.fromAddr, tt.args.toAddr, tt.args.passphrase, tt.args.amount, false)
			if (err != nil) != tt.wantErr {
				t.Errorf("Send() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Send() got = %v, want %v", got, tt.want)
			}
		})
	}
}

//func Test_newTx(t *testing.T) {
//	type args struct {
//		cdc        *codec.Codec
//		fromAddr   sdk.Address
//		tmNode     client.Client
//		keybase    keys.Keybase
//		passphrase string
//	}
//	tests := []struct {
//		name          string
//		args          args
//		wantTxBuilder auth.TxBuilder
//		wantCliCtx    util.CLIContext
//	}{
//		{"test newTx", args{
//			cdc:        makeTestCodec(),
//			fromAddr:   sdk.Address(getRandomPubKey().Address()),
//			tmNode:     mock.Client{
//				ABCIClient:     mock.ABCIMock{},
//				SignClient:     mock.ABCIMock{},
//				HistoryClient:  nil,
//				StatusClient:   nil,
//				EventsClient:   nil,
//				EvidenceClient: nil,
//				MempoolClient:  nil,
//				Service:        nil,
//			},
//			keybase:    nil,
//			passphrase: "",
//		},auth.TxBuilder{},
//		util.CLIContext{}},
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			gotTxBuilder, gotCliCtx := newTx(tt.args.cdc, tt.args.fromAddr, tt.args.tmNode, tt.args.keybase, tt.args.passphrase)
//			if !reflect.DeepEqual(gotTxBuilder, tt.wantTxBuilder) {
//				t.Errorf("newTx() gotTxBuilder = %v, want %v", gotTxBuilder, tt.wantTxBuilder)
//			}
//			if !reflect.DeepEqual(gotCliCtx, tt.wantCliCtx) {
//				t.Errorf("newTx() gotCliCtx = %v, want %v", gotCliCtx, tt.wantCliCtx)
//			}
//		})
//	}
//}
