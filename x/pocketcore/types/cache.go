package types

import (
	"encoding/hex"
	"fmt"
	sdk "github.com/pokt-network/posmint/types"
	db "github.com/tendermint/tm-db"
	"github.com/willf/bloom"
	"log"
	"sync"
)

var (
	// cache for session objects
	globalSessionCache *CacheStorage
	// cache for evidence objects
	globalEvidenceCache *CacheStorage
	// sync.once to perform initialization
	cacheOnce sync.Once
)

// "CacheStorage" - Contains an LRU cache and a database instance w/ mutex
type CacheStorage struct {
	Cache *Cache     // lru cache
	DB    db.DB      // persisted
	l     sync.Mutex // lock TODO NOTE: this mutex doesn't guarantee data consistency for gets and sets (max relays could go over)
}

type CacheObject interface {
	Marshal() ([]byte, error)
	Unmarshal(b []byte) (CacheObject, error)
	Key() ([]byte, error)
}

// "Init" - Initializes a cache storage object
func (cs *CacheStorage) Init(dir, name string, dbType db.DBBackendType, maxEntries int) {
	// init the lru cache with a max entries
	var err error
	cs.Cache, err = New(maxEntries)
	if err != nil {
		log.Fatal(fmt.Errorf("could not initialize cache storage: " + err.Error()))
	}
	// intialize the db
	cs.DB = db.NewDB(name, dbType, dir)
}

// "Get" - Returns the value from a key
func (cs *CacheStorage) Get(key []byte, object CacheObject) (interface{}, bool) {
	cs.l.Lock()
	defer cs.l.Unlock()
	// get the object using hex string of key
	if res, ok := cs.Cache.Get(hex.EncodeToString(key)); ok {
		return res, true
	}
	// not in cache, so search database
	bz := cs.DB.Get(key)
	if len(bz) == 0 {
		return nil, false
	}
	res, err := object.Unmarshal(bz)
	if err != nil {
		fmt.Printf("Error in CacheStorage.Get(): %s\n", err.Error())
		return nil, true
	}
	// add to cache
	cs.Cache.Add(hex.EncodeToString(key), res)
	return res, true
}

// "Set" - Sets the KV pair in cache and db
func (cs *CacheStorage) Set(key []byte, val CacheObject) {
	if cs.Cache.Len() == cs.Cache.Cap() {
		err := cs.FlushToDB()
		if err != nil {
			fmt.Printf("cache storage cannot be flushed to database (in set): %s", err.Error())
		}
	}
	cs.l.Lock()
	defer cs.l.Unlock()
	// add to cache
	cs.Cache.Add(hex.EncodeToString(key), val)
}

// "Delete" - Deletes the item from stores
func (cs *CacheStorage) Delete(key []byte) {
	cs.l.Lock()
	defer cs.l.Unlock()
	// remove from cache
	cs.Cache.Remove(hex.EncodeToString(key))
	// remove from db
	cs.DB.Delete(key)
}

func (cs *CacheStorage) FlushToDB() error {
	cs.l.Lock()
	defer cs.l.Unlock()
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
		bz, err := co.Marshal()
		if err != nil {
			return fmt.Errorf("error flushing database, marshalling value for DB: %s", err.Error())
		}
		kBz, err := hex.DecodeString(key)
		if err != nil {
			return fmt.Errorf("error flushing database, couldn't hex decode key: %s", err.Error())
		}
		// set to DB
		cs.DB.Set(kBz, bz)
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
	iter := cs.DB.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		cs.DB.Delete(iter.Key())
	}
}

// "Iterator" - Returns an iterator for all of the items in the stores
func (cs *CacheStorage) Iterator() db.Iterator {
	_ = cs.FlushToDB()
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
	s, err := session.Unmarshal(si.Iterator.Value())
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
	return SessionIt{
		Iterator: globalSessionCache.Iterator(),
	}
}

// "GetEvidence" - Retrieves the evidence object from the storage
func GetEvidence(header SessionHeader, evidenceType EvidenceType, max sdk.Int) (evidence Evidence, err error) {
	// generate the key for the evidence
	key, err := KeyForEvidence(header, evidenceType)
	if err != nil {
		return
	}
	// get the bytes from the storage
	val, found := globalEvidenceCache.Get(key, evidence)
	if !found && max.Equal(sdk.ZeroInt()) {
		return Evidence{}, fmt.Errorf("evidence not found")
	}
	if !found {
		bloomFilter := bloom.NewWithEstimates(uint(sdk.NewUintFromBigInt(max.BigInt()).Uint64()), .01)
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
		err = fmt.Errorf("could not unmarshal into session from cache with header %v", header)
		return
	}
	return
}

// "SetEvidence" - Sets an evidence object in the storage
func SetEvidence(evidence Evidence) {
	// generate the key for the evidence
	key, err := KeyForEvidence(evidence.SessionHeader, evidence.EvidenceType)
	if err != nil {
		return
	}
	globalEvidenceCache.Set(key, evidence)
}

// "DeleteEvidence" - Delete the evidence from the stores
func DeleteEvidence(header SessionHeader, evidenceType EvidenceType) error {
	// generate key for evidence
	key, err := KeyForEvidence(header, evidenceType)
	if err != nil {
		return err
	}
	// delete from cache
	globalEvidenceCache.Delete(key)
	return nil
}

// "ClearEvidence" - Clear stores of all evidence
func ClearEvidence() {
	if globalEvidenceCache != nil {
		globalEvidenceCache.Clear()
	}
}

// "EvidenceIt" - An evidence iterator instance of the globalEvidenceCache
type EvidenceIt struct {
	db.Iterator
}

// "Value" - Returns the evidence object value of the iterator
func (ei *EvidenceIt) Value() (evidence Evidence) {
	// unmarshal the value (bz) into an evidence object
	e, err := evidence.Unmarshal(ei.Iterator.Value())
	if err != nil {
		log.Fatal(fmt.Errorf("can't unmarshal evidence iterator value into evidence: %s", err.Error()))
	}
	evidence, ok := e.(Evidence)
	if !ok {
		log.Fatal("can't unmarshal evidence iterator value into evidence: cache object is not evidence")
	}
	return
}

// "EvidenceIterator" - Returns a globalEvidenceCache iterator instance
func EvidenceIterator() EvidenceIt {
	return EvidenceIt{
		Iterator: globalEvidenceCache.Iterator(),
	}
}

// "GetProof" - Returns the Proof object from a specific piece of evidence at a certain index
func GetProof(header SessionHeader, evidenceType EvidenceType, index int64) Proof {
	// retrieve the evidence
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

// "SetProof" - Sets a proof object in the evidence, using the header and evidence type
func SetProof(header SessionHeader, evidenceType EvidenceType, p Proof, max sdk.Int) {
	// retireve the evidence
	evidence, err := GetEvidence(header, evidenceType, max)
	// if not found generate the evidence object
	if err != nil {
		log.Fatalf("could not set proof object: %s", err.Error())
	}
	// add proof
	evidence.AddProof(p)
	// set evidence back
	SetEvidence(evidence)
}

func IsUniqueProof(p Proof, evidence Evidence) bool {
	return !evidence.Bloom.Test(p.Hash())
}

// "GetTotalProofs" - Returns the total number of proofs for a piece of evidence
func GetTotalProofs(h SessionHeader, et EvidenceType, maxPossibleRelays sdk.Int) (Evidence, int64) {
	// retrieve the evidence
	evidence, err := GetEvidence(h, et, maxPossibleRelays)
	if err != nil {
		log.Fatalf("could not get total proofs for evidence: %s", err.Error())
	}
	// return number of proofs
	return evidence, evidence.NumOfProofs
}
