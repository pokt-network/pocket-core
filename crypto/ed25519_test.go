package crypto

import (
	"reflect"
	"testing"
)

func TestEd25519PrivateKey_PrivateKeyFromBytes(t *testing.T) {
	type args struct {
		b []byte
	}
	privkey := getRandomPrivateKey(t)
	privkey2 := Ed25519PrivateKey{}
	_ = cdc.UnmarshalBinaryBare(privkey.Bytes(), &privkey2)

	tests := []struct {
		name    string
		ed      Ed25519PrivateKey
		args    args
		want    PrivateKey
		wantErr bool
	}{
		{
			name:    "Default",
			ed:      Ed25519PrivateKey{},
			args:    args{privkey.RawBytes()},
			want:    privkey,
			wantErr: false,
		},
		{
			name:    "Compatibility with amino decoding",
			ed:      Ed25519PrivateKey{},
			args:    args{privkey.RawBytes()},
			want:    privkey2,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.ed.PrivateKeyFromBytes(tt.args.b)
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
