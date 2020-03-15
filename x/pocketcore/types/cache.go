package types

import (
	"encoding/hex"
	"fmt"
	"github.com/hashicorp/golang-lru"
	db "github.com/tendermint/tm-db"
	"sync"
)

var (
	globalSessionCache  *CacheStorage
	globalEvidenceCache *CacheStorage
	cacheOnce           sync.Once
)

type CacheStorage struct {
	Cache *lru.Cache // lru cache
	DB    db.DB      // persisted
	l     sync.Mutex // lock
}

func InitCache(evidenceDir, sessionDir string, sessionDBType, evidenceDBType db.DBBackendType, maxEvidenceEntries, maxSessionEntries int) {
	cacheOnce.Do(func() {
		globalEvidenceCache = new(CacheStorage)
		globalSessionCache = new(CacheStorage)
		globalEvidenceCache.Init(evidenceDir, "evidence", evidenceDBType, maxEvidenceEntries)
		globalSessionCache.Init(sessionDir, "session", sessionDBType, maxSessionEntries)
	})
}

func (cs *CacheStorage) Init(dir, name string, dbType db.DBBackendType, maxEntries int) {
	var err error
	cs.Cache, err = lru.New(maxEntries)
	if err != nil {
		panic(err)
	}
	cs.DB = db.NewDB(name, dbType, dir)
}

func (cs *CacheStorage) Get(key []byte) ([]byte, bool) {
	cs.l.Lock()
	defer cs.l.Unlock()
	if res, ok := cs.Cache.Get(hex.EncodeToString(key)); ok {
		return res.([]byte), true
	}
	// not in cache, so search database
	bz := cs.DB.Get(key)
	if len(bz) == 0 {
		return nil, false
	}
	var value []byte
	err := ModuleCdc.UnmarshalJSON(bz, &value)
	if err != nil {
		return nil, false
	}
	// add to cache
	cs.Cache.Add(key, value)
	return value, true
}

func (cs *CacheStorage) Set(key []byte, val []byte) {
	cs.l.Lock()
	defer cs.l.Unlock()
	// add to cache
	cs.Cache.Add(hex.EncodeToString(key), val)
	cs.DB.Set(key, val)
}

func (cs *CacheStorage) Delete(key []byte) {
	cs.l.Lock()
	defer cs.l.Unlock()
	cs.Cache.Remove(hex.EncodeToString(key))
	cs.DB.Delete(key)
}

func (cs *CacheStorage) Clear() {
	cs.l.Lock()
	defer cs.l.Unlock()
	// clear cache
	cs.Cache.Purge()
	// clear db todo is this the best way to clear a db?
	iter := cs.DB.Iterator(nil, nil)
	defer iter.Close()
	for ; iter.Valid(); iter.Next() {
		cs.DB.Delete(iter.Key())
	}
}

func (cs *CacheStorage) Iterator() db.Iterator {
	return cs.DB.Iterator(nil, nil)
}

func GetSession(header SessionHeader) (session Session, found bool) {
	key := header.Hash()
	val, found := globalSessionCache.Get(key)
	if !found {
		return Session{}, found
	}
	err := ModuleCdc.UnmarshalJSON(val, &session)
	if err != nil {
		panic(fmt.Sprintf("could not unmarshal into session from cache: %s", err.Error()))
	}
	return
}

func SetSession(session Session) {
	key := session.SessionHeader.Hash()
	bz, err := ModuleCdc.MarshalJSON(session)
	if err != nil {
		panic(fmt.Sprintf("could not marshal into session for cache: %s", err.Error()))
	}
	globalSessionCache.Set(key, bz)
}

func DeleteSession(header SessionHeader) {
	key := header.Hash()
	globalSessionCache.Delete(key)
}

func ClearSessionCache() {
	globalSessionCache.Clear()
}

type SessionIt struct {
	db.Iterator
}

func (ei *SessionIt) Value() (session Session) {
	err := ModuleCdc.UnmarshalJSON(ei.Iterator.Value(), &session)
	if err != nil {
		panic(fmt.Errorf("can't unmarshal session iterator value into session: %s", err.Error()))
	}
	return
}

func SessionIterator() SessionIt {
	return SessionIt{
		Iterator: globalSessionCache.Iterator(),
	}
}

// todo save merkle tree to evidence

func GetEvidence(header SessionHeader, evidenceType EvidenceType) (evidence Evidence, found bool) {
	key := KeyForEvidence(header, evidenceType)
	bz, found := globalEvidenceCache.Get(key)
	if !found {
		return
	}
	err := ModuleCdc.UnmarshalJSON(bz, &evidence)
	if err != nil {
		panic(fmt.Sprintf("could not unmarshal into evidence from cache: %s", err.Error()))
	}
	return
}

func SetEvidence(evidence Evidence, evidenceType EvidenceType) {
	key := KeyForEvidence(evidence.SessionHeader, evidenceType)
	bz, err := ModuleCdc.MarshalJSON(evidence)
	if err != nil {
		panic(fmt.Sprintf("could not marshal into evidence for cache: %s", err.Error()))
	}
	globalEvidenceCache.Set(key, bz)
}

func DeleteEvidence(header SessionHeader, evidenceType EvidenceType) {
	// delete from cache
	key := KeyForEvidence(header, evidenceType)
	globalEvidenceCache.Delete(key)
}

func ClearEvidence() {
	globalEvidenceCache.Clear()
}

type EvidenceIt struct {
	db.Iterator
}

func (ei *EvidenceIt) Value() (evidence Evidence) {
	err := ModuleCdc.UnmarshalJSON(ei.Iterator.Value(), &evidence)
	if err != nil {
		panic(fmt.Errorf("can't unmarshal evidence iterator value into evidence: %s", err.Error()))
	}
	return evidence
}

func EvidenceIterator() EvidenceIt {
	return EvidenceIt{
		Iterator: globalEvidenceCache.Iterator(),
	}
}

func GetProof(header SessionHeader, evidenceType EvidenceType, index int64) Proof {
	evidence, found := GetEvidence(header, evidenceType)
	if !found {
		return nil
	}
	if evidence.NumOfProofs-1 < index {
		return nil
	}
	return evidence.Proofs[index]
}

func SetProof(header SessionHeader, evidenceType EvidenceType, p Proof) {
	evidence, found := GetEvidence(header, evidenceType)
	if !found {
		evidence = Evidence{
			SessionHeader: header,
			NumOfProofs:   0,
			Proofs:        make([]Proof, 0),
		}
	}
	// add proof
	evidence.AddProof(p)
	// set evidence back
	SetEvidence(evidence, evidenceType)
}

func IsUniqueProof(h SessionHeader, p Proof) bool {
	evidence, found := GetEvidence(h, p.EvidenceType())
	if !found {
		return true
	}
	// iterate over evidence to see if unique // todo efficiency (store hashes in cache/db)
	for _, proof := range evidence.Proofs {
		if proof.HashString() == p.HashString() {
			return false
		}
	}
	return true
}

func GetTotalProofs(h SessionHeader, et EvidenceType) int64 {
	evidence, found := GetEvidence(h, et)
	if !found {
		return 0
	}
	return evidence.NumOfProofs
}
