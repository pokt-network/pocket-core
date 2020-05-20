package types

import (
	"encoding/hex"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	sdk "github.com/pokt-network/posmint/types"
	db "github.com/tendermint/tm-db"
	"github.com/willf/bloom"
	"log"
	"sync"
	"time"
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
	Cache *lru.Cache // lru cache
	DB    db.DB      // persisted
	l     sync.Mutex // lock
}

// "Init" - Initializes a cache storage object
func (cs *CacheStorage) Init(dir, name string, dbType db.DBBackendType, maxEntries int) {
	// init the lru cache with a max entries
	var err error
	cs.Cache, err = lru.New(maxEntries)
	if err != nil {
		log.Fatal(fmt.Errorf("could not initialize cache storage: " + err.Error()))
	}
	// intialize the db
	cs.DB = db.NewDB(name, dbType, dir)
}

// "Get" - Returns the value from a key
func (cs *CacheStorage) Get(key []byte) ([]byte, bool) {
	defer timeTrack(time.Now(), "get from db/cache")
	cs.l.Lock()
	defer cs.l.Unlock()
	// get the object using hex string of key
	if res, ok := cs.Cache.Get(hex.EncodeToString(key)); ok {
		return res.([]byte), true
	}
	//// not in cache, so search database
	//bz := cs.DB.Get(key)
	//if len(bz) == 0 {
	//	return nil, false
	//}
	//var value []byte
	//err := ModuleCdc.UnmarshalBinaryLengthPrefixed(bz, &value)
	//if err != nil {
	//	return nil, false
	//}
	//// add to cache
	//cs.Cache.Add(key, value)
	return nil, false
}

// "Set" - Sets the KV pair in cache and db
func (cs *CacheStorage) Set(key []byte, val []byte) {
	defer timeTrack(time.Now(), "set to database")
	cs.l.Lock()
	defer cs.l.Unlock()
	// add to cache
	cs.Cache.Add(hex.EncodeToString(key), val)
	// add to database
	//cs.DB.Set(key, val)
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
	return cs.DB.Iterator(nil, nil)
}

// "GetSession" - Returns a session (value) from the stores using a header (key)
func GetSession(header SessionHeader) (session Session, found bool) {
	// generate the key from the header
	key := header.Hash()
	// check stores
	val, found := globalSessionCache.Get(key)
	if !found {
		return Session{}, found
	}
	// if found, unmarshal into session object
	err := ModuleCdc.UnmarshalBinaryLengthPrefixed(val, &session)
	if err != nil {
		fmt.Println(fmt.Errorf("could not unmarshal into session from cache: %s with header %v", err.Error(), header))
		return Session{}, false
	}
	return
}

// "SetSession" - Sets a session (value) in the stores using the header (key)
func SetSession(session Session) {
	// get the key for the session
	key := session.SessionHeader.Hash()
	// marshal into amino-json bz
	bz, err := ModuleCdc.MarshalJSON(session)
	if err != nil {
		fmt.Println(fmt.Errorf("could not marshal into session for cache: %s", err.Error()))
		return
	}
	// set it in stores
	globalSessionCache.Set(key, bz)
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
	err := ModuleCdc.UnmarshalBinaryLengthPrefixed(si.Iterator.Value(), &session)
	if err != nil {
		log.Fatal(fmt.Errorf("can't unmarshal session iterator value into session: %s", err.Error()))
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
	defer timeTrack(time.Now(), "getEvidence()")
	// generate the key for the evidence
	key, err := KeyForEvidence(header, evidenceType)
	if err != nil {
		return
	}
	// get the bytes from the storage
	bz, found := globalEvidenceCache.Get(key)
	if !found && max.Equal(sdk.ZeroInt()) {
		return Evidence{}, fmt.Errorf("evidence not found")
	}
	if !found {
		return Evidence{
			Bloom:         bloom.NewWithEstimates(uint(sdk.NewUintFromBigInt(max.BigInt()).Uint64()), .01),
			SessionHeader: header,
			NumOfProofs:   0,
			Proofs:        make([]Proof, 0),
			EvidenceType:  evidenceType,
		}, nil
	}
	// unmarshal into evidence obj
	err = ModuleCdc.UnmarshalBinaryLengthPrefixed(bz, &evidence)
	if err != nil {
		return Evidence{}, fmt.Errorf("could not unmarshal into evidence from cache: %s", err.Error())
	}
	return
}

// "SetEvidence" - Sets an evidence object in the storage
func SetEvidence(evidence Evidence) {
	defer timeTrack(time.Now(), "set evidence")
	// generate the key for the evidence
	key, err := KeyForEvidence(evidence.SessionHeader, evidence.EvidenceType)
	if err != nil {
		return
	}
	// marshal into bytes to store
	s := time.Now()
	bz, err := ModuleCdc.MarshalBinaryLengthPrefixed(evidence)
	fmt.Println("marshalling in setevidence ", time.Since(s))
	if err != nil {
		fmt.Println(fmt.Errorf("could not marshal into evidence for cache: %s", err.Error()))
		return
	}
	// set in storage
	globalEvidenceCache.Set(key, bz)
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
	err := ModuleCdc.UnmarshalBinaryBare(ei.Iterator.Value(), &evidence)
	if err != nil {
		log.Fatal(fmt.Errorf("can't unmarshal evidence iterator value into evidence: %s", err.Error()))
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
	defer timeTrack(time.Now(), "setProof")
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
