package types

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/pokt-network/pocket-core/x/nodes/exported"
	"github.com/tendermint/tendermint/config"
	"log"
	"sort"
	"sync"

	sdk "github.com/pokt-network/pocket-core/types"
	db "github.com/tendermint/tm-db"
	"github.com/willf/bloom"
)

var (
	// cache for session objects
	globalSessionCache *CacheStorage
	// cache for session validators
	GlobalSessionVals *SessionValidators
	// cache for GOBEvidence objects
	globalEvidenceCache *CacheStorage
	// sync.once to perform initialization
	cacheOnce sync.Once

	globalEvidenceSealedMap map[string]struct{}
)

// "CacheStorage" - Contains an LRU cache and a database instance w/ mutex
type CacheStorage struct {
	Cache *sdk.Cache // lru cache
	DB    db.DB      // persisted
	l     sync.Mutex // lock
}

type CacheObject interface {
	MarshalObject() ([]byte, error)
	UnmarshalObject(b []byte) (CacheObject, error)
	Key() ([]byte, error)
	Seal() CacheObject
	IsSealed() bool
}

// "Init" - Initializes a cache storage object
func (cs *CacheStorage) Init(dir, name string, options config.LevelDBOptions, maxEntries int) {
	// init the lru cache with a max entries
	cs.Cache = sdk.NewCache(maxEntries)
	// intialize the db
	var err error
	cs.DB, err = sdk.NewLevelDB(name, dir, options.ToGoLevelDBOpts())
	if err != nil {
		panic(err)
	}
}

// "Get" - Returns the value from a key
func (cs *CacheStorage) Get(key []byte, object CacheObject) (interface{}, bool) {
	cs.l.Lock()
	defer cs.l.Unlock()
	return cs.GetWithoutLock(key, object)
}

func (cs *CacheStorage) GetWithoutLock(key []byte, object CacheObject) (interface{}, bool) {
	// get the object using hex string of key
	if res, ok := cs.Cache.Get(hex.EncodeToString(key)); ok {
		return res, true
	}
	// not in cache, so search database
	bz, _ := cs.DB.Get(key)
	if len(bz) == 0 {
		return nil, false
	}
	res, err := object.UnmarshalObject(bz)
	if err != nil {
		fmt.Printf("Error in CacheStorage.Get(): %s\n", err.Error())
		return nil, true
	}
	// add to cache
	cs.Cache.Add(hex.EncodeToString(key), res)
	return res, true
}

// "Seal" - Seals the cache object so it is no longer writable in the cache store
func (cs *CacheStorage) Seal(object CacheObject) (cacheObject CacheObject, isOK bool) {
	if object.IsSealed() {
		return object, true
	}
	cs.l.Lock()
	defer cs.l.Unlock()
	// get the key from the object
	k, err := object.Key()
	if err != nil {
		return object, false
	}
	// make READONLY
	sealed := object.Seal()
	// set in db and cache
	cs.SetWithoutLockAndSealCheck(hex.EncodeToString(k), sealed)
	return sealed, true
}

// "Set" - Sets the KV pair in cache and db
func (cs *CacheStorage) Set(key []byte, val CacheObject) {
	keyString := hex.EncodeToString(key)
	cs.l.Lock()
	defer cs.l.Unlock()
	// get object to check if sealed
	res, found := cs.GetWithoutLock(key, val)
	if found {
		co, ok := res.(CacheObject)
		if !ok {
			fmt.Printf("ERROR: cannot convert object into cache object (in set)")
			return
		}
		// if evidence, check sealed map
		if ev, ok := co.(Evidence); ok {
			if _, ok := globalEvidenceSealedMap[ev.HashString()]; ok {
				return
			}
		}
	}
	cs.SetWithoutLockAndSealCheck(keyString, val)
}

// "SetWithoutLockAndSealCheck" - CONTRACT: used in a function with lock
//                                          cache must be flushed to db before any DB iterator
func (cs *CacheStorage) SetWithoutLockAndSealCheck(key string, val CacheObject) {
	// flush to db
	if cs.Cache.Len() == cs.Cache.Cap() && !cs.Cache.Contains(key) {
		err := cs.FlushToDBWithoutLock()
		if err != nil {
			fmt.Printf("ERROR: cache storage cannot be flushed to database (in set): %s", err.Error())
			return
		}
	}
	// add to cache
	cs.Cache.Add(key, val)
}

// "Delete" - Deletes the item from stores
func (cs *CacheStorage) Delete(key []byte) {
	cs.l.Lock()
	defer cs.l.Unlock()
	// remove from cache
	cs.Cache.Remove(hex.EncodeToString(key))
	// remove from db
	_ = cs.DB.Delete(key)
}

func (cs *CacheStorage) FlushToDB() error {
	cs.l.Lock()
	defer cs.l.Unlock()
	return cs.FlushToDBWithoutLock()
}

func (cs *CacheStorage) FlushToDBWithoutLock() error {
	// flush all to database
	for {
		key, val, ok := cs.Cache.RemoveOldest()
		if !ok {
			break
		}
		// value should be cache object
		co, ok := val.(CacheObject)
		if !ok {
			return fmt.Errorf("object in cache does not impement the cache object interface")
		}
		// marshal object to bytes
		bz, err := co.MarshalObject()
		if err != nil {
			return fmt.Errorf("error flushing database, marshalling value for DB: %s", err.Error())
		}
		kBz, err := hex.DecodeString(key)
		if err != nil {
			return fmt.Errorf("error flushing database, couldn't hex decode key: %s", err.Error())
		}
		// set to DB
		_ = cs.DB.Set(kBz, bz)
	}
	return nil
}

// "Clear" - Deletes all items from stores
func (cs *CacheStorage) Clear() {
	cs.l.Lock()
	defer cs.l.Unlock()
	// clear cache
	cs.Cache.Purge()
	// clear db
	iter, _ := cs.DB.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		_ = cs.DB.Delete(iter.Key())
	}
}

// "Iterator" - Returns an iterator for all of the items in the stores
func (cs *CacheStorage) Iterator() (db.Iterator, error) {
	err := cs.FlushToDB()
	if err != nil {
		fmt.Printf("unable to flush to db before iterator created in cacheStorage Iterator(): %s", err.Error())
	}
	return cs.DB.Iterator(nil, nil)
}

// "GetSession" - Returns a session (value) from the stores using a header (key)
func GetSession(header SessionHeader) (session Session, found bool) {
	// generate the key from the header
	key := header.Hash()
	// check stores
	val, found := globalSessionCache.Get(key, session)
	if !found {
		return Session{}, found
	}
	session, ok := val.(Session)
	if !ok {
		fmt.Println(fmt.Errorf("could not unmarshal into session from cache with header %v", header))
	}
	return
}

// "SetSession" - Sets a session (value) in the stores using the header (key)
func SetSession(session Session) {
	// get the key for the session
	key := session.SessionHeader.Hash()
	globalSessionCache.Set(key, session)
}

// "DeleteSession" - Deletes a session (value) from the stores
func DeleteSession(header SessionHeader) {
	// delete from stores using header.ID as key
	globalSessionCache.Delete(header.Hash())
}

// "ClearSessionCache" - Clears all items from the session cache db
func ClearSessionCache() {
	if globalSessionCache != nil {
		globalSessionCache.Clear()
	}
}

// "SessionIt" - An iterator value for the sessionCache structure
type SessionIt struct {
	db.Iterator
}

// "Value" - returns the value of the iterator (session)
func (si *SessionIt) Value() (session Session) {
	s, err := session.UnmarshalObject(si.Iterator.Value())
	if err != nil {
		log.Fatal(fmt.Errorf("can't unmarshal session iterator value into session: %s", err.Error()))
	}
	session, ok := s.(Session)
	if !ok {
		log.Fatal("can't unmarshal session iterator value into session: cache object is not a session")
	}
	return
}

// "SessionIterator" - Returns an instance iterator of the globalSessionCache
func SessionIterator() SessionIt {
	it, _ := globalSessionCache.Iterator()
	return SessionIt{
		Iterator: it,
	}
}

type SessionValidators struct {
	// 3 deep nested map structure height || chain || address -> sessionValidator (Optimizing for speed)
	M map[int64]map[string]exported.ValidatorI // used to store session end state
	S map[int64]map[string][]SessionValidator  // used to store session start, needed for sorting
	L *sync.Mutex
}

// Set a session validator under all of it's respective relayChains
func SetSessionValidator(height int64, val exported.ValidatorI, setS bool) {
	if GlobalSessionVals == nil || height < 0 {
		return // Only set after initialization
	}
	GlobalSessionVals.L.Lock()
	defer GlobalSessionVals.L.Unlock()
	// delete sequence
	if !val.IsStaked() {
		if val.IsUnstaked() {
			delete(GlobalSessionVals.M[height], val.GetAddress().String())
		} else {
			// set sequence
			GlobalSessionVals.M[height][val.GetAddress().String()] = val
		}
		if setS {
			for _, c := range val.GetChains() {
				found := false
				valIndex := 0
				for index, sv := range GlobalSessionVals.S[height][c] { // TODO quick search algo
					if sv.Address.Equals(val.GetAddress()) {
						found = true
						valIndex = index
						break
					}
				}
				if found {
					// delete from slice
					GlobalSessionVals.S[height][c] = append(GlobalSessionVals.S[height][c][:valIndex], GlobalSessionVals.S[height][c][valIndex+1:]...)
				}
			}
		}
		return
	}
	// set sequence
	GlobalSessionVals.M[height][val.GetAddress().String()] = val
	if setS {
		for _, c := range val.GetChains() {
			found := false
			valIndex := 0
			for index, sv := range GlobalSessionVals.S[height][c] { // TODO quick search algo
				if sv.Address.Equals(val.GetAddress()) {
					found = true
					valIndex = index
					break
				}
			}
			if !found {
				GlobalSessionVals.S[height][c] = append(GlobalSessionVals.S[height][c], SessionValidator{
					Address:    val.GetAddress(),
					PublicKey:  val.GetPublicKey().RawString(),
					ServiceURL: val.GetServiceURL(),
				})
				sort.Slice(GlobalSessionVals.S[height][c], func(i, j int) bool {
					return bytes.Compare(GlobalSessionVals.S[height][c][i].Address, GlobalSessionVals.S[height][c][j].Address) == -1
				})
			} else { // needed for edit stake
				GlobalSessionVals.S[height][c][valIndex] = SessionValidator{
					Address:    val.GetAddress(),
					PublicKey:  val.GetPublicKey().RawString(),
					ServiceURL: val.GetServiceURL(),
				}
			}
		}
	}
	return
}

func GetSessionValidator(height int64, address sdk.Address) (val exported.ValidatorI, found bool) {
	GlobalSessionVals.L.Lock()
	defer GlobalSessionVals.L.Unlock()
	if GlobalSessionVals.M[height] == nil {
		return
	}
	val, found = GlobalSessionVals.M[height][address.String()]
	return
}

// get a group of session validators by
func GetSessionValidators(height int64, relayChain string) (svs []SessionValidator) {
	GlobalSessionVals.L.Lock()
	defer GlobalSessionVals.L.Unlock()
	if GlobalSessionVals.S[height] == nil {
		return
	}
	return GlobalSessionVals.S[height][relayChain]
}

func InitSessionValidators(height int64, sessionStartState, sessionEndState []exported.ValidatorI) {
	// lock for thread safety
	GlobalSessionVals.L.Lock()
	defer GlobalSessionVals.L.Unlock()
	// for each session height
	// create the map [relayChain] SessionValidators slice
	if GlobalSessionVals.S[height] == nil {
		GlobalSessionVals.S[height] = make(map[string][]SessionValidator)
	}
	// for each validator in the session start
	for _, val := range sessionStartState {
		if !val.IsStaked() {
			continue
		}
		// for each validator object let's create the corresponding sessionVal object
		sv := SessionValidator{
			Address:    val.GetAddress(),
			PublicKey:  val.GetPublicKey().RawString(),
			ServiceURL: val.GetServiceURL(),
		}
		for _, chain := range val.GetChains() {
			if GlobalSessionVals.S[height][chain] == nil {
				GlobalSessionVals.S[height][chain] = make([]SessionValidator, 0)
			}
			GlobalSessionVals.S[height][chain] = append(GlobalSessionVals.S[height][chain], sv)
		}
	}
	// pre-sort starting states by chain:
	// sort set of validators by chain by address
	for _, svmp := range GlobalSessionVals.S[height] {
		sort.Slice(svmp, func(i, j int) bool {
			return bytes.Compare(svmp[i].Address, svmp[j].Address) == -1
		})
	}
	// setup the sessionEndState 'map'
	for _, val := range sessionEndState {
		if val.IsUnstaked() {
			continue // skip over unstaked validators
		}
		if GlobalSessionVals.M[height] == nil {
			GlobalSessionVals.M[height] = make(map[string]exported.ValidatorI)
		}
		GlobalSessionVals.M[height][val.GetAddress().String()] = val
	}
}

func UpdateSessionValidators(sessionHeight, blocksAgoHeight, blocksPerSession int64) error {
	if sessionHeight <= 1 {
		return nil
	}
	// heights we'll prune
	pruneHeight := sessionHeight - blocksAgoHeight
	pruneHeight = (pruneHeight/blocksPerSession)*blocksPerSession + 1
	// next session end height
	lastSessionHeight := sessionHeight - blocksPerSession
	// lock for thread safety
	GlobalSessionVals.L.Lock()
	defer GlobalSessionVals.L.Unlock()

	// sanity check
	if GlobalSessionVals.M[lastSessionHeight] == nil {
		return fmt.Errorf("critical error when updating sessionValidators, previous session endstate (M) was nil for height: %d", lastSessionHeight)
	}
	// first let's prune the old ones
	if pruneHeight >= 1 {
		delete(GlobalSessionVals.S, pruneHeight)
		delete(GlobalSessionVals.M, pruneHeight)
	}
	// then let's add a new entry for the current session start height
	GlobalSessionVals.S[sessionHeight] = make(map[string][]SessionValidator)
	GlobalSessionVals.M[sessionHeight] = make(map[string]exported.ValidatorI)
	// copy endstate map
	for addr, val := range GlobalSessionVals.M[lastSessionHeight] {
		GlobalSessionVals.M[sessionHeight][addr] = val
		if val.IsUnstaking() {
			continue // skip over unstaking validators
		}
		for _, c := range val.GetChains() {
			GlobalSessionVals.S[sessionHeight][c] = append(GlobalSessionVals.S[sessionHeight][c], SessionValidator{
				Address:    val.GetAddress(),
				PublicKey:  val.GetPublicKey().RawString(),
				ServiceURL: val.GetServiceURL(),
			})
		}
	}
	// sort every slice for each chain
	for _, svSlice := range GlobalSessionVals.S[sessionHeight] {
		sort.Slice(svSlice, func(i, j int) bool { return bytes.Compare(svSlice[i].Address, svSlice[j].Address) == -1 })
	}
	return nil
}

func (s SessionValidator) MarshalObject() ([]byte, error) {
	return ModuleCdc.ProtoMarshalBinaryBare(&s)
}

func (s SessionValidator) UnmarshalObject(b []byte) (SessionValidator, error) {
	err := ModuleCdc.ProtoUnmarshalBinaryBare(b, &s)
	return s, err
}

// "GetEvidence" - Retrieves the GOBEvidence object from the storage
func GetEvidence(header SessionHeader, evidenceType EvidenceType, max sdk.BigInt) (evidence Evidence, err error) {
	// generate the key for the GOBEvidence
	key, err := KeyForEvidence(header, evidenceType)
	if err != nil {
		return
	}
	// get the bytes from the storage
	val, found := globalEvidenceCache.Get(key, evidence)
	if !found && max.Equal(sdk.ZeroInt()) {
		return Evidence{}, fmt.Errorf("GOBEvidence not found")
	}
	if !found {
		bloomFilter := bloom.NewWithEstimates(uint(sdk.NewUintFromBigInt(max.BigInt()).Uint64()), .01)
		// add to metric
		GlobalServiceMetric().AddSessionFor(header.Chain)
		return Evidence{
			Bloom:         *bloomFilter,
			SessionHeader: header,
			NumOfProofs:   0,
			Proofs:        make([]Proof, 0),
			EvidenceType:  evidenceType,
		}, nil
	}
	evidence, ok := val.(Evidence)
	if !ok {
		err = fmt.Errorf("could not unmarshal into evidence from cache with header %v", header)
		return
	}
	if evidence.IsSealed() {
		return evidence, nil
	}
	// if hit relay limit... Seal the evidence
	if found && !max.Equal(sdk.ZeroInt()) && evidence.NumOfProofs >= max.Int64() {
		evidence, ok = SealEvidence(evidence)
		if !ok {
			err = fmt.Errorf("max relays is hit and could not seal evidence! GetEvidence() with header %v", header)
			return
		}
	}
	return
}

// "SetEvidence" - Sets an GOBEvidence object in the storage
func SetEvidence(evidence Evidence) {
	// generate the key for the evidence
	key, err := evidence.Key()
	if err != nil {
		return
	}
	globalEvidenceCache.Set(key, evidence)
}

// "DeleteEvidence" - Delete the GOBEvidence from the stores
func DeleteEvidence(header SessionHeader, evidenceType EvidenceType) error {
	// generate key for GOBEvidence
	key, err := KeyForEvidence(header, evidenceType)
	if err != nil {
		return err
	}
	// delete from cache
	globalEvidenceCache.Delete(key)
	delete(globalEvidenceSealedMap, header.HashString())
	return nil
}

// "SealEvidence" - Locks/sets the evidence from the stores
func SealEvidence(evidence Evidence) (Evidence, bool) {
	// delete from cache
	co, ok := globalEvidenceCache.Seal(evidence)
	if !ok {
		return Evidence{}, ok
	}
	e, ok := co.(Evidence)
	return e, ok
}

// "ClearEvidence" - Clear stores of all evidence
func ClearEvidence() {
	if globalEvidenceCache != nil {
		globalEvidenceCache.Clear()
		globalEvidenceSealedMap = make(map[string]struct{})
	}
}

// "EvidenceIt" - An GOBEvidence iterator instance of the globalEvidenceCache
type EvidenceIt struct {
	db.Iterator
}

// "Value" - Returns the GOBEvidence object value of the iterator
func (ei *EvidenceIt) Value() (evidence Evidence) {
	// unmarshal the value (bz) into an GOBEvidence object
	e, err := evidence.UnmarshalObject(ei.Iterator.Value())
	if err != nil {
		log.Fatal(fmt.Errorf("can't unmarshal GOBEvidence iterator value into GOBEvidence: %s", err.Error()))
	}
	evidence, ok := e.(Evidence)
	if !ok {
		log.Fatal("can't unmarshal GOBEvidence iterator value into GOBEvidence: cache object is not GOBEvidence")
	}
	return
}

// "EvidenceIterator" - Returns a globalEvidenceCache iterator instance
func EvidenceIterator() EvidenceIt {
	it, _ := globalEvidenceCache.Iterator()

	return EvidenceIt{
		Iterator: it,
	}
}

// "GetProof" - Returns the Proof object from a specific piece of GOBEvidence at a certain index
func GetProof(header SessionHeader, evidenceType EvidenceType, index int64) Proof {
	// retrieve the GOBEvidence
	evidence, err := GetEvidence(header, evidenceType, sdk.ZeroInt())
	if err != nil {
		return nil
	}
	// check for out of bounds
	if evidence.NumOfProofs-1 < index || index < 0 {
		return nil
	}
	// return the propoer proof
	return evidence.Proofs[index]
}

// "SetProof" - Sets a proof object in the GOBEvidence, using the header and GOBEvidence type
func SetProof(header SessionHeader, evidenceType EvidenceType, p Proof, max sdk.BigInt) {
	// retireve the GOBEvidence
	evidence, err := GetEvidence(header, evidenceType, max)
	// if not found generate the GOBEvidence object
	if err != nil {
		log.Fatalf("could not set proof object: %s", err.Error())
	}
	// add proof
	evidence.AddProof(p)
	// set GOBEvidence back
	SetEvidence(evidence)
}

func IsUniqueProof(p Proof, evidence Evidence) bool {
	return !evidence.Bloom.Test(p.Hash())
}

// "GetTotalProofs" - Returns the total number of proofs for a piece of GOBEvidence
func GetTotalProofs(h SessionHeader, et EvidenceType, maxPossibleRelays sdk.BigInt) (Evidence, int64) {
	// retrieve the GOBEvidence
	evidence, err := GetEvidence(h, et, maxPossibleRelays)
	if err != nil {
		log.Fatalf("could not get total proofs for GOBEvidence: %s", err.Error())
	}
	// return number of proofs
	return evidence, evidence.NumOfProofs
}
