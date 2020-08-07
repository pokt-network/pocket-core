package crypto

import (
	"reflect"
	"testing"
)

func TestSecp256k1PrivateKey_PrivateKeyFromBytes(t *testing.T) {
	type args struct {
		b []byte
	}

	privkey := getRandomPrivateKeySecp(t)
	privkey2 := Secp256k1PrivateKey{}
	_ = cdc.UnmarshalBinaryBare(privkey.Bytes(), &privkey2)

	tests := []struct {
		name    string
		se      Secp256k1PrivateKey
		args    args
		want    PrivateKey
		wantErr bool
	}{
		{
			name:    "Default",
			se:      Secp256k1PrivateKey{},
			args:    args{privkey.RawBytes()},
			want:    privkey,
			wantErr: false,
		},
		{
			name:    "Compatibility with amino decoding",
			se:      Secp256k1PrivateKey{},
			args:    args{privkey.RawBytes()},
			want:    privkey2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.se.PrivateKeyFromBytes(tt.args.b)
			if (err != nil) != tt.wantErr {
				t.Errorf("PrivateKeyFromBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrivateKeyFromBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}
