package types

import (
	"errors"
	"fmt"
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

// BaseAccount - a base account structure.
type BaseAccount struct {
	Address sdk.Address      `json:"address" yaml:"address"`
	Coins   sdk.Coins        `json:"coins" yaml:"coins"`
	PubKey  crypto.PublicKey `json:"public_key" yaml:"public_key"`
}

type Accounts []exported.Account

// NewBaseAccount creates a new BaseAccount object
func NewBaseAccount(address sdk.Address, coins sdk.Coins,
	pubKey crypto.PublicKey) *BaseAccount {
	return &BaseAccount{
		Address: address,
		Coins:   coins,
		PubKey:  pubKey,
	}
}

// String implements fmt.Stringer
func (acc BaseAccount) String() string {
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

// NewBaseAccountWithAddress - returns a new base account with a given address
func NewBaseAccountWithAddress(addr sdk.Address) BaseAccount {
	return BaseAccount{
		Address: addr,
	}
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

type marshalBaseAccount struct {
	Address sdk.Address
	Coins   sdk.Coins
	PubKey  string
}

// multisig account

var _ exported.Account = (*MultiSigAccount)(nil)

type MultiSigAccount struct {
	Address   sdk.Address              `json:"address"`
	PublicKey crypto.PublicKeyMultiSig `json:"public_key_multi_sig"`
	Coins     sdk.Coins                `json:"coins"`
}

func NewMultiSigAccount(publicKey crypto.PublicKeyMultiSig, coins sdk.Coins) *MultiSigAccount {
	return &MultiSigAccount{
		Address:   sdk.Address(publicKey.Address()),
		PublicKey: publicKey,
		Coins:     coins,
	}
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

var _ exported.ModuleAccountI = (*ModuleAccount)(nil)

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

// NewModuleAccount creates a new ModuleAccount instance
func NewModuleAccount(ba *BaseAccount, name string, permissions ...string) *ModuleAccount {
	if err := validatePermissions(permissions...); err != nil {
		fmt.Println(fmt.Errorf("invalid permissions for module account %s with permissions %v\n leaving permissionless", name, permissions))
		return &ModuleAccount{
			BaseAccount: ba,
			Name:        name,
			Permissions: []string{},
		}
	}

	return &ModuleAccount{
		BaseAccount: ba,
		Name:        name,
		Permissions: permissions,
	}
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
