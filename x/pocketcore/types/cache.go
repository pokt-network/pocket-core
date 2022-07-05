package types

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/tendermint/tendermint/config"
	"log"
	"path/filepath"
	"sync"
	"syscall"

	sdk "github.com/pokt-network/pocket-core/types"
	db "github.com/tendermint/tm-db"
	"github.com/willf/bloom"
)

var (
	// sync.once to perform initialization
	ConfigOnce sync.Once
)

// "CacheStorage" - Contains an LRU cache and a database instance w/ mutex
type CacheStorage struct {
	Cache   *sdk.Cache // lru cache
	DB      db.DB      // persisted
	l       sync.Mutex // lock
	SealMap *sync.Map
}

type CacheObject interface {
	MarshalObject() ([]byte, error)
	UnmarshalObject(b []byte) (CacheObject, error)
	Key() ([]byte, error)
	IsSealable() bool
	HashString() string
}

// "Init" - Initializes a cache storage object
func (cs *CacheStorage) Init(dir, name string, options config.LevelDBOptions, maxEntries int, inMemoryDB bool) {
	// init the lru cache with a max entries
	cs.Cache = sdk.NewCache(maxEntries)
	// intialize the db
	var err error
	if inMemoryDB {
		cs.DB = db.NewGoLevelMemDBWithCapacity(maxEntries)
		return
	}
	cs.DB, err = sdk.NewLevelDB(name, dir, options.ToGoLevelDBOpts())
	if err != nil {
		if err == syscall.EWOULDBLOCK {
			message := fmt.Sprintf("can't open files needed for execution. Another instance may be running. path: %s\n", filepath.Join(dir, name+".db"))
			panic(errors.New(message))
		}
		panic(err)
	}
	cs.SealMap = &sync.Map{}
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

// "DeleteEvidence" - Remove the GOBEvidence from the store
func DeleteEvidence(header SessionHeader, evidenceType EvidenceType, evidenceStore *CacheStorage) error {
	// generate key for GOBEvidence
	key, err := KeyForEvidence(header, evidenceType)
	if err != nil {
		return err
	}
	// delete from cache
	evidenceStore.Delete(key)
	evidenceStore.SealMap.Delete(header.HashString())
	return nil
}

func (cs *CacheStorage) IsSealedWithoutLock(object CacheObject) bool {
	_, ok := cs.SealMap.Load(object.HashString())
	return ok
}

func (cs *CacheStorage) IsSealed(object CacheObject) bool {
	cs.l.Lock()
	defer cs.l.Unlock()
	return cs.IsSealedWithoutLock(object)
}

// "Seal" - Seals the cache object so it is no longer writable in the cache store
func (cs *CacheStorage) Seal(object CacheObject) (cacheObject CacheObject, isOK bool) {

	if !object.IsSealable() {
		return object, false
	}

	if cs.IsSealed(object) {
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
	cs.SealMap.Store(object.HashString(), struct{}{})
	// set in db and cache
	cs.SetWithoutLockAndSealCheck(hex.EncodeToString(k), object)
	return object, true
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
		if co.IsSealable() {
			if ok := cs.IsSealedWithoutLock(co); ok {
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

// "Remove" - Deletes the item from stores
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
func GetSession(header SessionHeader, sessionStore *CacheStorage) (session Session, found bool) {
	// generate the key from the header
	key := header.Hash()
	// check stores
	val, found := sessionStore.Get(key, session)
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
func SetSession(session Session, sessionStore *CacheStorage) {
	// get the key for the session
	key := session.SessionHeader.Hash()
	sessionStore.Set(key, session)
}

// "DeleteSession" - Deletes a session (value) from the stores
func DeleteSession(header SessionHeader, sessionStore *CacheStorage) {
	// delete from stores using header.ID as key
	sessionStore.Delete(header.Hash())
}

// "ClearSessionCache" - Clears all items from the session cache db
func ClearSessionCache(sessionStore *CacheStorage) {
	if sessionStore != nil {
		sessionStore.Clear()
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

// "SessionIterator" - Returns an instance iterator of the GlobalSessionCache
func SessionIterator(sessionStore *CacheStorage) SessionIt {
	it, _ := sessionStore.Iterator()
	return SessionIt{
		Iterator: it,
	}
}

// "GetEvidence" - Retrieves the GOBEvidence object from the storage
func GetEvidence(header SessionHeader, evidenceType EvidenceType, max sdk.BigInt, storage *CacheStorage) (evidence Evidence, err error) {
	// generate the key for the GOBEvidence
	key, err := KeyForEvidence(header, evidenceType)
	if err != nil {
		return
	}
	// get the bytes from the storage
	val, found := storage.Get(key, evidence)
	if !found && max.Equal(sdk.ZeroInt()) {
		return Evidence{}, fmt.Errorf("GOBEvidence not found")
	}
	if !found {
		bloomFilter := bloom.NewWithEstimates(uint(sdk.NewUintFromBigInt(max.BigInt()).Uint64()), .01)
		// add to metric
		addSessionMetricFunc := func() {
			GlobalServiceMetric().AddSessionFor(header.Chain, nil)
		}
		if GlobalPocketConfig.LeanPocket {
			go addSessionMetricFunc()
		} else {
			addSessionMetricFunc()
		}
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
	if storage.IsSealed(evidence) {
		return evidence, nil
	}
	// if hit relay limit... Seal the evidence
	if found && !max.Equal(sdk.ZeroInt()) && evidence.NumOfProofs >= max.Int64() {
		evidence, ok = SealEvidence(evidence, storage)
		if !ok {
			err = fmt.Errorf("max relays is hit and could not seal evidence! GetEvidence() with header %v", header)
			return
		}
	}
	return
}

// "SetEvidence" - Sets an GOBEvidence object in the storage
func SetEvidence(evidence Evidence, evidenceStore *CacheStorage) {
	// generate the key for the evidence
	key, err := evidence.Key()
	if err != nil {
		return
	}
	evidenceStore.Set(key, evidence)
}

// "SealEvidence" - Locks/sets the evidence from the stores
func SealEvidence(evidence Evidence, storage *CacheStorage) (Evidence, bool) {
	co, ok := storage.Seal(evidence)
	if !ok {
		return Evidence{}, ok
	}
	e, ok := co.(Evidence)
	return e, ok
}

// "ClearEvidence" - Clear stores of all evidence
func ClearEvidence(evidenceStore *CacheStorage) {
	if evidenceStore != nil {
		evidenceStore.Clear()
		evidenceStore.SealMap = &sync.Map{}
	}
}

// "EvidenceIt" - An GOBEvidence iterator instance of the GlobalEvidenceCache
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

// "EvidenceIterator" - Returns a GlobalEvidenceCache iterator instance
func EvidenceIterator(evidenceStore *CacheStorage) EvidenceIt {
	it, _ := evidenceStore.Iterator()
	return EvidenceIt{
		Iterator: it,
	}
}

// "GetProof" - Returns the Proof object from a specific piece of GOBEvidence at a certain index
func GetProof(header SessionHeader, evidenceType EvidenceType, index int64, evidenceStore *CacheStorage) Proof {
	// retrieve the GOBEvidence
	evidence, err := GetEvidence(header, evidenceType, sdk.ZeroInt(), evidenceStore)
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
func SetProof(header SessionHeader, evidenceType EvidenceType, p Proof, max sdk.BigInt, evidenceStore *CacheStorage) {
	// retireve the GOBEvidence
	evidence, err := GetEvidence(header, evidenceType, max, evidenceStore)
	// if not found generate the GOBEvidence object
	if err != nil {
		log.Fatalf("could not set proof object: %s", err.Error())
	}
	// add proof
	evidence.AddProof(p)
	// set GOBEvidence back
	SetEvidence(evidence, evidenceStore)
}

func IsUniqueProof(p Proof, evidence Evidence) bool {
	return !evidence.Bloom.Test(p.Hash())
}

// "GetTotalProofs" - Returns the total number of proofs for a piece of GOBEvidence
func GetTotalProofs(h SessionHeader, et EvidenceType, maxPossibleRelays sdk.BigInt, evidenceStore *CacheStorage) (Evidence, int64) {
	// retrieve the GOBEvidence
	evidence, err := GetEvidence(h, et, maxPossibleRelays, evidenceStore)
	if err != nil {
		log.Fatalf("could not get total proofs for GOBEvidence: %s", err.Error())
	}
	// return number of proofs
	return evidence, evidence.NumOfProofs
}
