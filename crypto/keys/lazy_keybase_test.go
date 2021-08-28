package keys

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	pocketCrypto "github.com/pokt-network/pocket-core/crypto"
	"github.com/pokt-network/pocket-core/types"
	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
)

func NewTestCaseDir(t *testing.T) (string, func()) {
	dir, err := ioutil.TempDir("", t.Name()+"_")
	require.NoError(t, err)
	return dir, func() { os.RemoveAll(dir) }
}

func TestNew(t *testing.T) {
	dir, cleanup := NewTestCaseDir(t)
	defer cleanup()
	kb := New("keybasename", dir)
	lazykb, ok := kb.(*lazyKeybase)
	require.True(t, ok)
	require.Equal(t, lazykb.name, "keybasename")
	require.Equal(t, lazykb.dir, dir)
}

func Test_lazyKeybase_CloseDB(t *testing.T) {
	type fields struct {
		name     string
		dir      string
		coinbase KeyPair
	}

	dir, cleanup := NewTestCaseDir(t)
	defer cleanup()

	tests := []struct {
		name   string
		fields fields
	}{
		{"Test Keybase CloseDB", fields{
			name:     "keybasename",
			dir:      dir,
			coinbase: KeyPair{},
		}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lkb := lazyKeybase{
				name:     tt.fields.name,
				dir:      tt.fields.dir,
				coinbase: tt.fields.coinbase,
			}
			lkb.CloseDB()
		})
	}
}

func Test_lazyKeybase_Create(t *testing.T) {
	type fields struct {
		name     string
		dir      string
		coinbase KeyPair
	}
	type args struct {
		encryptPassphrase string
	}

	dir, cleanup := NewTestCaseDir(t)
	defer cleanup()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    KeyPair
		wantErr bool
	}{
		{"Test Lazykeybase Create", fields{
			name:     "keybase",
			dir:      dir,
			coinbase: KeyPair{},
		}, args{encryptPassphrase: "ENCRYPTIONPASSPHRASE"},
			KeyPair{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lkb := lazyKeybase{
				name:     tt.fields.name,
				dir:      tt.fields.dir,
				coinbase: tt.fields.coinbase,
			}
			got, err := lkb.Create(tt.args.encryptPassphrase)
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !assert.NotNil(t, got) {
				t.Errorf("Returned Nil")
			}
			if !assert.NotNil(t, got.GetAddress()) {
				t.Errorf("keypair with errors ")
			}

		})
	}
}

func Test_lazyKeybase_Delete(t *testing.T) {
	type fields struct {
		name     string
		dir      string
		coinbase KeyPair
	}

	dir, cleanup := NewTestCaseDir(t)
	defer cleanup()

	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{"Test keybase delete", fields{
			name:     "keybasename",
			dir:      dir,
			coinbase: KeyPair{},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lkb := lazyKeybase{
				name:     tt.fields.name,
				dir:      tt.fields.dir,
				coinbase: tt.fields.coinbase,
			}
			wkp, err := lkb.Create("ENCRYPTIONPASSPHRASE")
			if err != nil {
				t.Fatalf("Creation Failed")
			}
			if err := lkb.Delete(wkp.GetAddress(), "ENCRYPTIONPASSPHRASE"); (err != nil) != tt.wantErr {
				t.Fatalf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_lazyKeybase_ExportPrivKeyEncryptedArmor(t *testing.T) {
	type fields struct {
		name     string
		dir      string
		coinbase KeyPair
	}
	type args struct {
		address           types.Address
		decryptPassphrase string
		encryptPassphrase string
	}

	dir, cleanup := NewTestCaseDir(t)
	defer cleanup()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Test ExportPrivKeyEncryptedArmor", fields{
			name:     "keybasename",
			dir:      dir,
			coinbase: KeyPair{},
		}, args{
			address:           nil,
			decryptPassphrase: "ENCRYPTIONPASSPHRASE",
			encryptPassphrase: "ENCRYPTIONPASSPHRASE",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lkb := lazyKeybase{
				name:     tt.fields.name,
				dir:      tt.fields.dir,
				coinbase: tt.fields.coinbase,
			}
			wkp, err := lkb.Create("ENCRYPTIONPASSPHRASE")
			if err != nil {
				t.Fatalf("Creation Failed")
			}
			gotArmor, err := lkb.ExportPrivKeyEncryptedArmor(wkp.GetAddress(), tt.args.decryptPassphrase, tt.args.encryptPassphrase, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("ExportPrivKeyEncryptedArmor() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotArmor == "" {
				t.Errorf("ExportPrivKeyEncryptedArmor() gotArmor = %v", gotArmor)
			}
		})
	}
}

func Test_lazyKeybase_ExportPrivateKeyObject(t *testing.T) {
	type fields struct {
		name     string
		dir      string
		coinbase KeyPair
	}
	type args struct {
		address    types.Address
		passphrase string
	}

	dir, cleanup := NewTestCaseDir(t)
	defer cleanup()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Test ExportPrivateKeyObject", fields{
			name:     "keybasename",
			dir:      dir,
			coinbase: KeyPair{},
		}, args{
			address:    nil,
			passphrase: "ENCRYPTIONPASSPHRASE",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lkb := lazyKeybase{
				name:     tt.fields.name,
				dir:      tt.fields.dir,
				coinbase: tt.fields.coinbase,
			}
			wkp, err := lkb.Create("ENCRYPTIONPASSPHRASE")
			if err != nil {
				t.Fatalf("Creation Failed")
			}
			got, err := lkb.ExportPrivateKeyObject(wkp.GetAddress(), tt.args.passphrase)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExportPrivateKeyObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.NotNil(t, got) {
				t.Errorf("ExportPrivateKeyObject() got = %v", got)
			}
		})
	}
}

func Test_lazyKeybase_Get(t *testing.T) {
	type fields struct {
		name     string
		dir      string
		coinbase KeyPair
	}
	type args struct {
		address types.Address
	}

	dir, cleanup := NewTestCaseDir(t)
	defer cleanup()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    KeyPair
		wantErr bool
	}{
		{"Test Lazykeybase Get", fields{
			name:     "keybasename",
			dir:      dir,
			coinbase: KeyPair{},
		}, args{address: nil}, KeyPair{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lkb := lazyKeybase{
				name:     tt.fields.name,
				dir:      tt.fields.dir,
				coinbase: tt.fields.coinbase,
			}

			wkp, err := lkb.Create("ENCRYPTIONPASSPHRASE")
			if err != nil {
				t.Fatalf("Creation Failed")
			}
			got, err := lkb.Get(wkp.GetAddress())
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.NotNil(t, got) && assert.NotEmpty(t, got) {
				t.Errorf("Get() got = %v", got)
			}
		})
	}
}

func Test_lazyKeybase_GetCoinbase(t *testing.T) {
	type fields struct {
		name     string
		dir      string
		coinbase KeyPair
	}

	dir, cleanup := NewTestCaseDir(t)
	defer cleanup()

	tests := []struct {
		name    string
		fields  fields
		want    KeyPair
		wantErr bool
	}{
		{"Test Lazykeybase GetCoinbase", fields{
			name:     "keybasename",
			dir:      dir,
			coinbase: KeyPair{},
		}, KeyPair{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kb := &lazyKeybase{
				name:     tt.fields.name,
				dir:      tt.fields.dir,
				coinbase: tt.fields.coinbase,
			}

			_, err := kb.Create("ENCRYPTIONPASSPHRASE")
			if err != nil {
				t.Fatalf(err.Error())
			}
			got, err := kb.GetCoinbase()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCoinbase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.NotNil(t, got) {
				t.Errorf("GetCoinbase() got = %v", got)
			}
		})
	}
}

func Test_lazyKeybase_ImportPrivKey(t *testing.T) {
	type fields struct {
		name     string
		dir      string
		coinbase KeyPair
	}

	type args struct {
		armor             string
		decryptPassphrase string
		encryptPassphrase string
	}

	dir, cleanup := NewTestCaseDir(t)
	defer cleanup()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    KeyPair
		wantErr bool
	}{
		{"Test Lazykeybase Importprivkey", fields{
			name:     "keybasename",
			dir:      dir,
			coinbase: KeyPair{},
		}, args{
			armor:             "",
			decryptPassphrase: "ENCRYPTIONPASSPHRASE",
			encryptPassphrase: "ENCRYPTIONPASSPHRASE",
		}, KeyPair{},
			false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lkb := lazyKeybase{
				name:     tt.fields.name,
				dir:      tt.fields.dir,
				coinbase: tt.fields.coinbase,
			}

			wkp, err := lkb.Create("ENCRYPTIONPASSPHRASE")
			if err != nil {
				t.Fatalf("Creation Failed")
			}
			armor, err := lkb.ExportPrivKeyEncryptedArmor(wkp.GetAddress(), tt.args.decryptPassphrase, tt.args.encryptPassphrase, "")
			if err != nil {
				t.Fatalf("Creation Failed")
			}
			_ = lkb.Delete(wkp.GetAddress(), "ENCRYPTIONPASSPHRASE")

			got, err := lkb.ImportPrivKey(armor, tt.args.decryptPassphrase, tt.args.encryptPassphrase)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImportPrivKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.NotNil(t, got) {
				t.Errorf("ImportPrivKey() got = %v", got)
			}
		})
	}
}

func Test_lazyKeybase_ImportPrivateKeyObject(t *testing.T) {
	type fields struct {
		name     string
		dir      string
		coinbase KeyPair
	}
	type args struct {
		privateKey        [64]byte
		encryptPassphrase string
	}

	dir, cleanup := NewTestCaseDir(t)
	defer cleanup()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    KeyPair
		wantErr bool
	}{
		{"Test LazyKeybase ImportPrivateKeyObject", fields{
			name:     "keybasename",
			dir:      dir,
			coinbase: KeyPair{},
		}, args{
			privateKey:        [64]byte{},
			encryptPassphrase: "ENCRYPTIONPASSPHRASE",
		}, KeyPair{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lkb := lazyKeybase{
				name:     tt.fields.name,
				dir:      tt.fields.dir,
				coinbase: tt.fields.coinbase,
			}

			wkp, err := lkb.Create("ENCRYPTIONPASSPHRASE")
			if err != nil {
				t.Fatalf("Creation Failed")
			}
			exported, err := lkb.ExportPrivateKeyObject(wkp.GetAddress(), "ENCRYPTIONPASSPHRASE")
			if err != nil {
				t.Fatalf("Creation Failed")
			}
			_ = lkb.Delete(wkp.GetAddress(), "ENCRYPTIONPASSPHRASE")

			got, err := lkb.ImportPrivateKeyObject(exported.(pocketCrypto.Ed25519PrivateKey), tt.args.encryptPassphrase)
			if (err != nil) != tt.wantErr {
				t.Errorf("ImportPrivateKeyObject() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.NotNil(t, got) {
				t.Errorf("ImportPrivateKeyObject() got = %v", got)
			}
		})
	}
}

func Test_lazyKeybase_List(t *testing.T) {
	type fields struct {
		name     string
		dir      string
		coinbase KeyPair
	}

	dir, cleanup := NewTestCaseDir(t)
	defer cleanup()

	tests := []struct {
		name    string
		fields  fields
		want    []KeyPair
		wantErr bool
	}{
		{"Test lazyKeyBase List", fields{
			name:     "keybasename",
			dir:      dir,
			coinbase: KeyPair{},
		}, []KeyPair{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lkb := lazyKeybase{
				name:     tt.fields.name,
				dir:      tt.fields.dir,
				coinbase: tt.fields.coinbase,
			}

			_, err := lkb.Create("ENCRYPTIONPASSPHRASE")
			if err != nil {
				t.Fatalf(err.Error())
			}
			got, err := lkb.List()
			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.NotEmpty(t, got) {
				t.Errorf("List() got = %v", got)
			}
		})
	}
}

func Test_lazyKeybase_SetCoinbase(t *testing.T) {
	type fields struct {
		name     string
		dir      string
		coinbase KeyPair
	}
	type args struct {
		address types.Address
	}

	dir, cleanup := NewTestCaseDir(t)
	defer cleanup()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Test SetCoinbase", fields{
			name:     "keybasename",
			dir:      dir,
			coinbase: KeyPair{},
		}, args{address: nil}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			kb := &lazyKeybase{
				name:     tt.fields.name,
				dir:      tt.fields.dir,
				coinbase: tt.fields.coinbase,
			}

			wkb, _ := kb.Create("ENCRYPTIONPASSPHRASE")

			if err := kb.SetCoinbase(wkb.GetAddress()); (err != nil) != tt.wantErr {
				t.Errorf("SetCoinbase() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_lazyKeybase_Sign(t *testing.T) {
	type fields struct {
		name     string
		dir      string
		coinbase KeyPair
	}
	type args struct {
		address    types.Address
		passphrase string
		msg        []byte
	}

	dir, cleanup := NewTestCaseDir(t)
	defer cleanup()

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		want1   pocketCrypto.PublicKey
		wantErr bool
	}{
		{"Test Sign", fields{
			name:     "keybasename",
			dir:      dir,
			coinbase: KeyPair{},
		}, args{
			address:    nil,
			passphrase: "ENCRYPTIONPASSPHRASE",
			msg:        []byte("test"),
		}, []byte{}, nil, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lkb := lazyKeybase{
				name:     tt.fields.name,
				dir:      tt.fields.dir,
				coinbase: tt.fields.coinbase,
			}

			wkb, _ := lkb.Create("ENCRYPTIONPASSPHRASE")

			got, got1, err := lkb.Sign(wkb.GetAddress(), tt.args.passphrase, tt.args.msg)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !assert.NotNil(t, got) {
				t.Errorf("Sign() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, wkb.PublicKey) {
				t.Errorf("Sign() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}

func Test_lazyKeybase_Update(t *testing.T) {
	type fields struct {
		name     string
		dir      string
		coinbase KeyPair
	}
	type args struct {
		address types.Address
		oldpass string
		newpass string
	}

	dir, cleanup := NewTestCaseDir(t)
	defer cleanup()

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"Test lazyKeybase Update", fields{
			name:     "keybasename",
			dir:      dir,
			coinbase: KeyPair{},
		}, args{
			address: nil,
			oldpass: "ENCRYPTIONPASSPHRASE",
			newpass: "ENCRYPTIONPASSPHRASE2",
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lkb := lazyKeybase{
				name:     tt.fields.name,
				dir:      tt.fields.dir,
				coinbase: tt.fields.coinbase,
			}

			wkb, _ := lkb.Create("ENCRYPTIONPASSPHRASE")

			if err := lkb.Update(wkb.GetAddress(), tt.args.oldpass, tt.args.newpass); (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if err := lkb.Delete(wkb.GetAddress(), tt.args.newpass); (err != nil) != tt.wantErr {
				t.Errorf("Remove() error = %v, wantErr %v", err, tt.wantErr)
			}

		})
	}
}
