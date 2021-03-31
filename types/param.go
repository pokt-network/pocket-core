package types

import (
	"errors"
	"reflect"

	"github.com/pokt-network/pocket-core/codec"

	"github.com/pokt-network/pocket-core/store/prefix"
)

const (
	paramsKey  = "params"
	paramsTKey = "transient_params"
)

var (
	ParamsKey  = NewKVStoreKey(paramsKey)
	ParamsTKey = NewTransientStoreKey(paramsTKey)
)

// Individual parameter store for each keeper
// Transient store persists for a block, so we use it for
// recording whether the parameter has been changed or not
type Subspace struct {
	cdc   *codec.Codec
	key   StoreKey // []byte -> []byte, stores parameter
	tkey  StoreKey // []byte -> bool, stores parameter change
	name  []byte
	table KeyTable
}

// NewSubspace constructs a store with namestore
func NewSubspace(name string) (res Subspace) {
	res = Subspace{
		cdc:  cdc,
		key:  ParamsKey,
		tkey: ParamsTKey,
		name: []byte(name),
		table: KeyTable{
			m: make(map[string]attribute),
		},
	}
	return
}

func (s *Subspace) SetCodec(cdc *codec.Codec) {
	s.cdc = cdc
}

// WithKeyTable initializes KeyTable and returns modified Subspace
func (s Subspace) WithKeyTable(table KeyTable) Subspace {
	if table.m == nil {
		panic("SetKeyTable() called with nil KeyTable")
	}
	if len(s.table.m) != 0 {
		panic("SetKeyTable() called on already initialized Subspace")
	}

	for k, v := range table.m {
		s.table.m[k] = v
	}

	// Allocate additional capicity for Subspace.name
	// So we don't have to allocate extra space each time appending to the key
	name := s.name
	s.name = make([]byte, len(name), len(name)+table.maxKeyLength())
	copy(s.name, name)

	return s
}

// Returns a KVStore identical with ctx.KVStore(s.key).Prefix()
func (s Subspace) kvStore(ctx Ctx) KVStore {
	// append here is safe, appends within a function won't cause
	// weird side effects when its singlethreaded
	return prefix.NewStore(ctx.KVStore(s.key), append(s.name, '/'))
}

// Returns a transient store for modification
func (s Subspace) transientStore(ctx Ctx) KVStore {
	// append here is safe, appends within a function won't cause
	// weird side effects when its singlethreaded
	return prefix.NewStore(ctx.TransientStore(s.tkey), append(s.name, '/'))
}

func concatKeys(key, subkey []byte) (res []byte) {
	res = make([]byte, len(key)+1+len(subkey))
	copy(res, key)
	res[len(key)] = '/'
	copy(res[len(key)+1:], subkey)
	return
}

// Get parameter from store
func (s Subspace) Get(ctx Ctx, key []byte, ptr interface{}) {
	store := s.kvStore(ctx)
	bz, err := store.Get(key)
	if err != nil {
		ctx.Logger().Error("error getting a value from a key in the subspace, could be an empty subspace:", err.Error())
		return
	}
	err = s.cdc.UnmarshalJSON(bz, ptr)
	if err != nil {
		ctx.Logger().Error("error unmarshalling from the subspace, could be an empty subspace", err.Error())
		return
	}
}

func (s Subspace) GetAllParamKeys(ctx Ctx) (keys []string) {
	store := s.kvStore(ctx)
	it, _ := store.Iterator(nil, nil)
	for ; it.Valid(); it.Next() {
		keys = append(keys, string(it.Key()))
	}
	return
}

// GetIfExists do not modify ptr if the stored parameter is nil
func (s Subspace) GetIfExists(ctx Ctx, key []byte, ptr interface{}) {
	store := s.kvStore(ctx)
	bz, _ := store.Get(key)
	if bz == nil {
		return
	}
	err := s.cdc.UnmarshalJSON(bz, ptr)
	if err != nil {
		panic(err)
	}
}

func (s Subspace) GetIfExistsRaw(ctx Ctx, key []byte) []byte {
	store := s.kvStore(ctx)
	bz, _ := store.Get(key)
	if bz == nil {
		return nil
	}
	return bz
}

// GetWithSubkey returns a parameter with a given key and a subkey.
func (s Subspace) GetWithSubkey(ctx Ctx, key, subkey []byte, ptr interface{}) {
	s.Get(ctx, concatKeys(key, subkey), ptr)
}

// GetWithSubkeyIfExists  returns a parameter with a given key and a subkey but does not
// modify ptr if the stored parameter is nil.
func (s Subspace) GetWithSubkeyIfExists(ctx Ctx, key, subkey []byte, ptr interface{}) {
	s.GetIfExists(ctx, concatKeys(key, subkey), ptr)
}

// Get raw bytes of parameter from store
func (s Subspace) GetRaw(ctx Ctx, key []byte) ([]byte, error) {
	store := s.kvStore(ctx)
	return store.Get(key)
}

// Check if the parameter is set in the store
func (s Subspace) Has(ctx Ctx, key []byte) (bool, error) {
	store := s.kvStore(ctx)
	return store.Has(key)
}

// Returns true if the parameter is set in the block
func (s Subspace) Modified(ctx Ctx, key []byte) (bool, error) {
	tstore := s.transientStore(ctx)
	return tstore.Has(key)
}

func (s Subspace) checkType(store KVStore, key []byte, param interface{}) {
	attr, ok := s.table.m[string(key)]
	if !ok {
		panic("Parameter not registered")
	}

	ty := attr.ty
	pty := reflect.TypeOf(param)
	if pty.Kind() == reflect.Ptr {
		pty = pty.Elem()
	}

	if pty.Kind() == reflect.Interface {
		if pty != ty && !pty.Implements(ty) {
			panic("Type mismatch with registered table")
		}
	}
	if pty != ty {
		panic("Type mismatch with registered table")
	}
}

// Set stores the parameter. It returns error if stored parameter has different type from input.
// It also set to the transient store to record change.
func (s Subspace) Set(ctx Ctx, key []byte, param interface{}) {
	store := s.kvStore(ctx)
	s.checkType(store, key, param)
	bz, err := s.cdc.MarshalJSON(param)
	if err != nil {
		panic(err)
	}
	_ = store.Set(key, bz)
	tstore := s.transientStore(ctx)
	_ = tstore.Set(key, []byte{})
}

// Update stores raw parameter bytes. It returns error if the stored parameter
// has a different type from the input. It also sets to the transient store to
// record change.
func (s Subspace) Update(ctx Ctx, key []byte, param []byte) error {
	attr, ok := s.table.m[string(key)]
	if !ok {
		panic("Parameter not registered")
	}

	ty := attr.ty
	dest := reflect.New(ty).Interface()
	s.GetIfExists(ctx, key, dest)
	err := s.cdc.UnmarshalJSON(param, dest)
	if err != nil {
		return err
	}

	s.Set(ctx, key, dest)
	tStore := s.transientStore(ctx)
	_ = tStore.Set(key, []byte{})

	return nil
}

// SetWithSubkey set a parameter with a key and subkey
// Checks parameter type only over the key
func (s Subspace) SetWithSubkey(ctx Ctx, key []byte, subkey []byte, param interface{}) {
	store := s.kvStore(ctx)

	s.checkType(store, key, param)

	newkey := concatKeys(key, subkey)

	bz, err := s.cdc.MarshalJSON(param)
	if err != nil {
		panic(err)
	}
	_ = store.Set(newkey, bz)

	tstore := s.transientStore(ctx)
	_ = tstore.Set(newkey, []byte{})
}

// UpdateWithSubkey stores raw parameter bytes  with a key and subkey. It checks
// the parameter type only over the key.
func (s Subspace) UpdateWithSubkey(ctx Ctx, key []byte, subkey []byte, param []byte) error {
	concatkey := concatKeys(key, subkey)

	attr, ok := s.table.m[string(concatkey)]
	if !ok {
		return errors.New("parameter not registered")
	}

	ty := attr.ty
	dest := reflect.New(ty).Interface()
	s.GetWithSubkeyIfExists(ctx, key, subkey, dest)
	err := s.cdc.UnmarshalJSON(param, dest)
	if err != nil {
		return err
	}

	s.SetWithSubkey(ctx, key, subkey, dest)
	tStore := s.transientStore(ctx)
	_ = tStore.Set(concatkey, []byte{})

	return nil
}

// Get to ParamSet
func (s Subspace) GetParamSet(ctx Ctx, ps ParamSet) {
	for _, pair := range ps.ParamSetPairs() {
		s.Get(ctx, pair.Key, pair.Value)
	}
}

// Set from ParamSet
func (s Subspace) SetParamSet(ctx Ctx, ps ParamSet) {
	for _, pair := range ps.ParamSetPairs() {
		// pair.Field is a pointer to the field, so indirecting the ptr.
		// go-amino automatically handles it but just for sure,
		// since SetStruct is meant to be used in InitGenesis
		// so this method will not be called frequently
		v := reflect.Indirect(reflect.ValueOf(pair.Value)).Interface()
		s.Set(ctx, pair.Key, v)
	}
}

// Returns name of Subspace
func (s Subspace) Name() string {
	return string(s.name)
}

// Wrapper of Subspace, provides immutable functions only
type ReadOnlySubspace struct {
	s Subspace
}

// Exposes Get
func (ros ReadOnlySubspace) Get(ctx Ctx, key []byte, ptr interface{}) {
	ros.s.Get(ctx, key, ptr)
}

// Exposes GetRaw
func (ros ReadOnlySubspace) GetRaw(ctx Ctx, key []byte) ([]byte, error) {
	return ros.s.GetRaw(ctx, key)
}

// Exposes Has
func (ros ReadOnlySubspace) Has(ctx Ctx, key []byte) (bool, error) {
	return ros.s.Has(ctx, key)
}

// Exposes Modified
func (ros ReadOnlySubspace) Modified(ctx Ctx, key []byte) (bool, error) {
	return ros.s.Modified(ctx, key)
}

// Exposes Space
func (ros ReadOnlySubspace) Name() string {
	return ros.s.Name()
}

type attribute struct {
	ty reflect.Type
}

// KeyTable subspaces appropriate type for each parameter key
type KeyTable struct {
	m map[string]attribute
}

// Constructs new table
func NewKeyTable(keytypes ...interface{}) (res KeyTable) {
	if len(keytypes)%2 != 0 {
		panic("odd number arguments in NewTypeKeyTable")
	}

	res = KeyTable{
		m: make(map[string]attribute),
	}

	for i := 0; i < len(keytypes); i += 2 {
		res = res.RegisterType(keytypes[i].([]byte), keytypes[i+1])
	}

	return
}

func isAlphaNumeric(key []byte) bool {
	for _, b := range key {
		if !((48 <= b && b <= 57) || // numeric
			(65 <= b && b <= 90) || // upper case
			(97 <= b && b <= 122)) { // lower case
			return false
		}
	}
	return true
}

// Register single key-type pair
func (t KeyTable) RegisterType(key []byte, ty interface{}) KeyTable {
	if len(key) == 0 {
		panic("cannot register empty key")
	}
	if !isAlphaNumeric(key) {
		panic("non alphanumeric parameter key")
	}
	keystr := string(key)
	if _, ok := t.m[keystr]; ok {
		panic("duplicate parameter key")
	}

	rty := reflect.TypeOf(ty)

	// Indirect rty if it is ptr
	if rty.Kind() == reflect.Ptr {
		rty = rty.Elem()
	}

	t.m[keystr] = attribute{
		ty: rty,
	}

	return t
}

// Register multiple pairs from ParamSet
func (t KeyTable) RegisterParamSet(ps ParamSet) KeyTable {
	for _, kvp := range ps.ParamSetPairs() {
		t = t.RegisterType(kvp.Key, kvp.Value)
	}
	return t
}

func (t KeyTable) maxKeyLength() (res int) {
	for k := range t.m {
		l := len(k)
		if l > res {
			res = l
		}
	}
	return
}

// Used for associating paramsubspace key and field of param structs
type ParamSetPair struct {
	Key   []byte      `json:"key"`
	Value interface{} `json:"value"`
}

// Slice of KeyFieldPair
type ParamSetPairs []ParamSetPair

// Interface for structs containing parameters for a module
type ParamSet interface {
	ParamSetPairs() ParamSetPairs
}
