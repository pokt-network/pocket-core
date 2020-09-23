package types

import (
	"errors"
	"fmt"
	"github.com/pokt-network/pocket-core/codec"
	"github.com/pokt-network/pocket-core/crypto"
	tmCrypto "github.com/tendermint/tendermint/crypto"
	"time"

	"gopkg.in/yaml.v2"

	sdk "github.com/pokt-network/pocket-core/types"
	"github.com/pokt-network/pocket-core/x/auth/exported"
)

//-----------------------------------------------------------------------------
// BaseAccount
var _ exported.Account = (*BaseAccount)(nil)
var _ codec.ProtoMarshaler = &BaseAccount{}

// BaseAccount - a base account structure.
type BaseAccount struct {
	Address sdk.Address      `json:"address" yaml:"address"`
	Coins   sdk.Coins        `json:"coins" yaml:"coins"`
	PubKey  crypto.PublicKey `json:"public_key" yaml:"public_key"`
}

func (acc *BaseAccount) Reset() {
	*acc = BaseAccount{}
}

func (acc *BaseAccount) ProtoMessage() {
	p := acc.ToProto()
	p.ProtoMessage()
}

func (acc *BaseAccount) Marshal() ([]byte, error) {
	p := acc.ToProto()
	return p.Marshal()
}

func (acc *BaseAccount) MarshalTo(data []byte) (n int, err error) {
	p := acc.ToProto()
	return p.MarshalTo(data)
}

func (acc *BaseAccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	p := acc.ToProto()
	return p.MarshalToSizedBuffer(dAtA)
}

func (acc *BaseAccount) Size() int {
	p := acc.ToProto()
	return p.Size()
}

func (acc *BaseAccount) Unmarshal(data []byte) error {
	var bae ProtoBaseAccount
	err := bae.Unmarshal(data)
	if err != nil {
		return err
	}
	ba, err := bae.FromProto()
	if err != nil {
		return err
	}
	*acc = ba
	return nil
}

type Accounts []exported.Account

// NewBaseAccountWithAddress - returns a new base account with a given address
func NewBaseAccountWithAddress(addr sdk.Address) BaseAccount {
	return BaseAccount{
		Address: addr,
	}
}

// String implements fmt.Stringer
func (acc *BaseAccount) String() string {
	var pubkey string
	if acc.PubKey != nil {
		pubkey = acc.PubKey.RawString()
	}
	return fmt.Sprintf(`Account:
  Address:       %s
  Pubkey:        %s
  Coins:         %s`,
		acc.Address, pubkey, acc.Coins,
	)
}

// GetAddress - Implements sdk.Account.
func (acc BaseAccount) GetAddress() sdk.Address {
	return acc.Address
}

// SetAddress - Implements sdk.Account.
func (acc *BaseAccount) SetAddress(addr sdk.Address) error {
	if len(acc.Address) != 0 {
		return errors.New("cannot override BaseAccount address")
	}
	acc.Address = addr
	return nil
}

// GetPubKey - Implements sdk.Account.
func (acc BaseAccount) GetPubKey() crypto.PublicKey {
	return acc.PubKey
}

// SetPubKey - Implements sdk.Account.
func (acc *BaseAccount) SetPubKey(pubKey crypto.PublicKey) error {
	acc.PubKey = pubKey
	return nil
}

// GetCoins - Implements sdk.Account.
func (acc *BaseAccount) GetCoins() sdk.Coins {
	return acc.Coins
}

// SetCoins - Implements sdk.Account.
func (acc *BaseAccount) SetCoins(coins sdk.Coins) error {
	acc.Coins = coins
	return nil
}

// SpendableCoins returns the total set of spendable coins. For a base account,
// this is simply the base coins.
func (acc *BaseAccount) SpendableCoins(_ time.Time) sdk.Coins {
	return acc.GetCoins()
}

// MarshalYAML returns the YAML representation of an account.
func (acc BaseAccount) MarshalYAML() (interface{}, error) {
	var bs []byte
	var err error
	var pubkey string

	if acc.PubKey != nil {
		pubkey = acc.PubKey.RawString()
	}

	bs, err = yaml.Marshal(marshalBaseAccount{
		Address: acc.Address,
		Coins:   acc.Coins,
		PubKey:  pubkey,
	})
	if err != nil {
		return nil, err
	}

	return string(bs), err
}

func (acc BaseAccount) ToProto() ProtoBaseAccount {
	var pk []byte
	if acc.PubKey != nil {
		pk = acc.PubKey.RawBytes()
	}
	return ProtoBaseAccount{
		Address: acc.Address,
		Coins:   acc.Coins,
		PubKey:  pk,
	}
}

type marshalBaseAccount struct {
	Address sdk.Address
	Coins   sdk.Coins
	PubKey  string
}

// multisig account

var _ exported.Account = (*MultiSigAccount)(nil)
var _ codec.ProtoMarshaler = &MultiSigAccount{}

type MultiSigAccount struct {
	Address   sdk.Address              `json:"address"`
	PublicKey crypto.PublicKeyMultiSig `json:"public_key_multi_sig"`
	Coins     sdk.Coins                `json:"coins"`
}

func (m MultiSigAccount) GetAddress() sdk.Address {
	return m.Address
}

func (m *MultiSigAccount) SetAddress(_ sdk.Address) error {
	if m.Address != nil && len(m.Address) != 0 {
		return sdk.ErrInternal(fmt.Sprintf("address already set: %s", m.Address))
	}
	if m.PublicKey == nil {
		return sdk.ErrInternal("cannot have a nil public key for a multisig account")
	}
	m.Address = sdk.Address(m.PublicKey.Address())
	return nil
}

func (m MultiSigAccount) GetPubKey() crypto.PublicKey {
	return m.PublicKey
}

func (m MultiSigAccount) GetMultiPubKey() crypto.PublicKeyMultiSig {
	return m.PublicKey
}

func (m MultiSigAccount) SetPubKey(pk crypto.PublicKey) error {
	p, ok := pk.(crypto.PublicKeyMultiSig)
	if !ok {
		return sdk.ErrInternal("the public key must be of interface type: PublicKeyMultiSig")
	}
	m.PublicKey = p
	return nil
}

func (m MultiSigAccount) GetCoins() sdk.Coins {
	return m.Coins
}

func (m *MultiSigAccount) SetCoins(c sdk.Coins) error {
	m.Coins = c
	return nil
}

func (m MultiSigAccount) SpendableCoins(blockTime time.Time) sdk.Coins {
	return m.GetCoins()
}

func (m MultiSigAccount) String() string {
	return fmt.Sprintf(`
  Address:       %s
  Pubkey:        %s
  Coins:         %s`,
		m.Address, m.PublicKey, m.Coins,
	)
}

func (m *MultiSigAccount) Reset() {
	*m = MultiSigAccount{}
}

func (m MultiSigAccount) ProtoMessage() {
	p := m.ToProto()
	p.ProtoMessage()
}

func (m MultiSigAccount) Marshal() ([]byte, error) {
	p := m.ToProto()
	return p.Marshal()
}

func (m MultiSigAccount) MarshalTo(data []byte) (n int, err error) {
	p := m.ToProto()
	return p.MarshalTo(data)
}

func (m MultiSigAccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	p := m.ToProto()
	return p.MarshalToSizedBuffer(dAtA)
}

func (m MultiSigAccount) Size() int {
	p := m.ToProto()
	return p.Size()
}

func (m *MultiSigAccount) Unmarshal(data []byte) error {
	var pms ProtoMultiSigAccount
	err := pms.Unmarshal(data)
	if err != nil {
		return err
	}
	msa, err := pms.FromProto()
	if err != nil {
		return err
	}
	*m = msa
	return nil
}

func (m MultiSigAccount) ToProto() ProtoMultiSigAccount {
	return ProtoMultiSigAccount{
		Address: m.Address,
		PubKey:  m.PublicKey.RawBytes(),
		Coins:   m.Coins,
	}
}

func (pms ProtoMultiSigAccount) FromProto() (MultiSigAccount, error) {
	pk, err := crypto.PublicKeyMultiSignature{}.NewPublicKey(pms.PubKey)
	if err != nil {
		return MultiSigAccount{}, err
	}
	pkms, ok := pk.(crypto.PublicKeyMultiSignature)
	if !ok {
		return MultiSigAccount{}, fmt.Errorf("%s", "multisig account must have multipublickey type")
	}
	return MultiSigAccount{
		Address:   pms.Address,
		PublicKey: pkms,
		Coins:     pms.Coins,
	}, nil
}

var _ exported.ModuleAccountI = (*ModuleAccount)(nil)
var _ codec.ProtoMarshaler = &ModuleAccount{}

// ModuleAccount defines an account for modules that holds coins on a pool
type ModuleAccount struct {
	*BaseAccount
	Name        string   `json:"name" yaml:"name"`               // name of the module
	Permissions []string `json:"permissions" yaml:"permissions"` // permissions of module account
}

// NewModuleAddress creates an Address from the hash of the module's name
func NewModuleAddress(name string) sdk.Address {
	return sdk.Address(tmCrypto.AddressHash([]byte(name)))
}

func NewEmptyModuleAccount(name string, permissions ...string) *ModuleAccount {
	moduleAddress := NewModuleAddress(name)
	baseAcc := NewBaseAccountWithAddress(moduleAddress)

	if err := validatePermissions(permissions...); err != nil {
		fmt.Println(fmt.Errorf("invalid permissions for module account %s with permissions %v\n leaving permissionless", name, permissions))
		return &ModuleAccount{
			BaseAccount: &baseAcc,
			Name:        name,
			Permissions: []string{},
		}
	}

	return &ModuleAccount{
		BaseAccount: &baseAcc,
		Name:        name,
		Permissions: permissions,
	}
}

func (ma *ModuleAccount) Reset() {
	*ma = ModuleAccount{}
}

func (ma *ModuleAccount) ProtoMessage() {
	p := ma.ToProto()
	p.ProtoMessage()
}

func (ma *ModuleAccount) Marshal() ([]byte, error) {
	p := ma.ToProto()
	return p.Marshal()
}

func (ma *ModuleAccount) MarshalTo(data []byte) (n int, err error) {
	p := ma.ToProto()
	return p.MarshalTo(data)
}

func (ma *ModuleAccount) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	p := ma.ToProto()
	return p.MarshalToSizedBuffer(dAtA)
}

func (ma *ModuleAccount) Size() int {
	p := ma.ToProto()
	return p.Size()
}

func (ma *ModuleAccount) Unmarshal(data []byte) error {
	var mae ProtoModuleAccount
	err := mae.Unmarshal(data)
	if err != nil {
		return err
	}
	m, err := mae.FromProto()
	if err != nil {
		return err
	}
	*ma = m
	return nil
}

// HasPermission returns whether or not the module account has permission.
func (ma ModuleAccount) HasPermission(permission string) bool {
	for _, perm := range ma.Permissions {
		if perm == permission {
			return true
		}
	}
	return false
}

// GetName returns the the name of the holder's module
func (ma ModuleAccount) GetName() string {
	return ma.Name
}

// GetPermissions returns permissions granted to the module account
func (ma ModuleAccount) GetPermissions() []string {
	return ma.Permissions
}

// SetPubKey - Implements Account
func (ma ModuleAccount) SetPubKey(pubKey crypto.PublicKey) error {
	return fmt.Errorf("not supported for module accounts")
}

// String follows stringer interface
func (ma ModuleAccount) String() string {
	b, err := yaml.Marshal(ma)
	if err != nil {
		fmt.Println("couldn't convert module account to yaml string: " + err.Error())
		return ""
	}
	return string(b)
}

// MarshalYAML returns the YAML representation of a ModuleAccount.
func (ma ModuleAccount) MarshalYAML() (interface{}, error) {
	bs, err := yaml.Marshal(struct {
		Address     sdk.Address
		Coins       sdk.Coins
		PubKey      string
		Name        string
		Permissions []string
	}{
		Address:     ma.Address,
		Coins:       ma.Coins,
		PubKey:      "",
		Name:        ma.Name,
		Permissions: ma.Permissions,
	})

	if err != nil {
		return nil, err
	}

	return string(bs), nil
}

func (ma ModuleAccount) ToProto() ProtoModuleAccount {
	ba := ma.BaseAccount.ToProto()
	return ProtoModuleAccount{
		ProtoBaseAccount: ba,
		Name:             ma.Name,
		Permissions:      ma.Permissions,
	}
}

func (m *ProtoBaseAccount) FromProto() (ba BaseAccount, err error) {
	var pk crypto.PublicKey
	if m.PubKey != nil {
		pk, err = crypto.NewPublicKeyBz(m.PubKey)
		if err != nil {
			return BaseAccount{}, err
		}
	}
	return BaseAccount{
		Address: m.Address,
		Coins:   m.Coins,
		PubKey:  pk,
	}, nil
}

func (m *ProtoModuleAccount) FromProto() (ModuleAccount, error) {
	ba, err := m.ProtoBaseAccount.FromProto()
	if err != nil {
		return ModuleAccount{}, nil
	}
	return ModuleAccount{
		BaseAccount: &ba,
		Name:        m.Name,
		Permissions: m.Permissions,
	}, nil
}
