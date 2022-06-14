## Pocket Network Persistence Replacement

### Background:
The IAVL store (Cosmos’ app level storage tree) is the legacy implementation that serves as both the merkle root integrity structure
and the general purpose database structure for reads. Though the structure is seemingly efficient at generating merkle roots for the
ABCI AppHash, Pocket Network observes a considerable downgrade in performance for historical queries, block processing, and overall
disk size management. The idea behind Pocket Network’s Persistence Replacement is to preserve the efficient state commitment root
generation while providing an alternative implementation for DB reads. The idea is to prune N-X IAVL versions in order to enable merkle
tree generation, while maintaining a completely separate storage layer for efficient storage, historical queries, and block processing.
The idea and design is very similar to Cosmos’ [ADR-40](https://github.com/cosmos/cosmos-sdk/discussions/8297).

Required Goals:
- [ ] Simpler Codebase
- [ ] Non-Consensus Breaking
- [ ] Comparable syncing performance
- [ ] Comparable disk usage

Additional Goals:
- [ ] Faster historical queries
- [ ] Faster current queries
- [ ] Smaller disk usage
- [ ] Better syncing performance
- [ ] Slower disk usage rate

### Approaches Considered / Tested / Combined for the query database:
1. Cache with warmup (load last N heights into memory).
- A full mem go-leveldb compatible cache was created just for this test. The real challenge of creating an in memory goleveldb
is that it must match the exact behavior when *deleting items while iterating*. This behavior is often undefined / undocumented so
matching the behavior was riddled with trial and error testing. A weakness of the cache design is the only way to ensure the exact
behavior is to do a seek during each Next() call which hinders the OLog(n) performance of the BTree iterator implementation. The idea
with the cache is to write current heights and read from the last N heights. During commit() you persist the height to disk and get a
little speed up due to the batch write functionality of leveldb. In the end, little to no performance was seen with this cacheDB and was
omitted due to long ‘warmup times’ and code complexity.
- NOTE: Go-leveldb comes with a memdb that is append only - don’t get your hopes up
- NOTE: TenderminttDB comes with a memdb that is not iterator compatible with goleveldb - don’t get your hopes up

2. Postgres implementation
- The idea was to implement a V1 standard postgresDB in order to take advantage of PgSQL’s considerable scalability and lookup capabilities.
Though the idea is logically sound, the issue is the way the Pocket Core V0 ABCI app is designed. Pocket Core inherits the Cosmos SDK architecture
which uses Protobuf (previously Amino) encoding and serves arbitrary bytes to the storage layer. This architecture leaves a big challenge for
the Postgres implementation: A custom key interpreter / ORM needs to be designed in order to appropriately assign the key /values to the
proper tables. With careful design, this overhead may be minimized from a performance standpoint, but it will never be zero. The other major
challenge with a Postgres implementation is that the current setups of Node Runners will need to be migrated to support a database engine as
opposed to a file system db. In early prototyping, block processing performance was heavily degraded for the Postgres implementation, however the
datadir size shrunk considerably. The design was ultimately discarded due to the above challenges and disappointing performance.

3. Straight K/V (AKA Waves)
- A straightforward, naive approach to just store K/V in a golevel db for each height without any deduplication. Performance was by far the best with
this branch, seeing block production times and query performance improve by a minimum of 3X the RC-0.6.3 implementation. The issue with this branch,
is that without any de-duplication logic, the storage disk space requirement was dramatically increased (almost 2x). Unfortunately, this rendered ‘waves’
an unacceptable solution.

4. Deduplicated Waves Using Index Store
- Spec

  Two Datastores:

  - Link Store (Indexing) - index keys recorded at every possible height that ‘links’ to the data blob
  - Data Store (Data Bytes) - hash keys that contains the actual data value

  - Link Store:
    - KEY: height uint32 + prefixByte byte + dataKey
    - VAL: 16 Byte Hash Of LinkKEY

  - Data Store
    - KEY: 16 Byte Hash Of LinkKEY
    - VAL: Data Blob

- Results:
Comparable datastore size, comparable performance for block processing, faster historical queries, and far simpler code complexity.
This is likely an ‘acceptable’ solution but not optimal and somewhat underwhelming given the overhaul.

5. Deduplication (2) Using 'Heights' Index
- Spec

  Two Datastores:

    - Link Store (Indexing) - index keys recorded at every possible height that ‘links’ to the data blob
    - Data Store (Data Bytes) - hash keys that contains the actual data value

    - Link Store:
        - KEY: prefixByte byte + dataKey
        - VAL: Array of heights Ex. (0, 1, 2, 100)

    - Data Store
        - KEY: 16 Byte Hash Of LinkKEY
        - VAL: Data Blob

- Details
```
GET:
// first lookup
<prefix> + <dataKey> -> array of heights (choose latest compared to query)
// second lookup
<prefix> + <dataKey> + <latestHeight> -> data blob

SET/DEL
// first write
<prefix> + <dataKey> -> APPEND TO (array of heights 1,2,3,101)
// second write
<prefix> + <dataKey> + <newHeight> -> data blob (or DEL for deleted)

ITER (Issue: iterating over orphans)
// step 1 iterate over first write ^
<prefix> + <dataKey>
// step 2 upon first lookup (Next() call)
get latest height relative to query
// step 3 is ensure the value is not nil (deleted)
if nil (skip) -> Next() again
```

- Results
  Approximately 10x slower than current implementation with 40% disk savings. Deemed unacceptable due to performance

6. Dedup (3) AKA Slim  ✅
- Specification
  No index, use reverse iterator 'seek' to get lexicographically [ELEN](https://www.zanopha.com/docs/elen.pdf) sorted
  heightKeys from the datastore. Load last N heights into memory and do nothing upon store Commit() besides increment height.
  Use specially optimized functions: DeleteHeight() and PrepareHeight() which should take advantage of golangs
  built in delete(map,key) and copy(dst,src) functions to increment the in-memory store.

- Details
```
GET:
// get the latest version by taking advantage of goleveldb.Seek() and use a reverse iterator to find the next appropriate lexicographically sorted height key
REVERSE ITERATOR -> <prefix> + <dataKey> + <latestHeight> -> data blob

SET/DEL (2 sets per write)
// exists store (cached) <the main reason we do this is to have a space to iterate over>
<exists>+<prefix>+<dataKey> -> (true if set, false if del)
// ‘main’ set
<prefix> + <dataKey> + <newHeight> -> data blob (or DEL for deleted)

ITER (Issue: iterating over orphans)
// step 1 iterate over exists store
<prefix> + <dataKey>
// step 2 do the GET() logic to find the latest version key
latestKey = <prefix> + <dataKey> + <latestHeightRespectiveToQuery> -> data blob
// step 3 is ensure the value is not nil (deleted)
if DEL or future height-> skip
```

### Result Data:

**Block Processing (Sampled 15 heights starting on 44376)**

*Before:*

&nbsp;&nbsp;&nbsp;&nbsp;Total: 154s

&nbsp;&nbsp;&nbsp;&nbsp;Avg:      10s

*After:*

&nbsp;&nbsp;&nbsp;&nbsp;Total: 135s

&nbsp;&nbsp;&nbsp;&nbsp;Avg:       9s

**Historical Querying (20 historical queries while syncing (query nodes 4000<0-9> --nodeLimit=20000)**

*Before:*

&nbsp;&nbsp;&nbsp;&nbsp;Total: 60.66s

&nbsp;&nbsp;&nbsp;&nbsp;Avg:    3.03s

*After:*

&nbsp;&nbsp;&nbsp;&nbsp;Total: 39.37s:

&nbsp;&nbsp;&nbsp;&nbsp;Avg:    1.97s

**Disk Usage (Height 44335)**

*Before:*

&nbsp;&nbsp;&nbsp;&nbsp;*du -h data/*

&nbsp;&nbsp;&nbsp;&nbsp;12G    data//txindexer.db

&nbsp;&nbsp;&nbsp;&nbsp;63G    data//application.db

&nbsp;&nbsp;&nbsp;&nbsp;…

&nbsp;&nbsp;&nbsp;&nbsp;TOTAL:  131G    data/

*After:*

&nbsp;&nbsp;&nbsp;&nbsp;*du -h data/*

&nbsp;&nbsp;&nbsp;&nbsp;6.6G	/txindexer.db

&nbsp;&nbsp;&nbsp;&nbsp;13G	/application.db

&nbsp;&nbsp;&nbsp;&nbsp;…

&nbsp;&nbsp;&nbsp;&nbsp;TOTAL:  77G	/data

**41% Total reduction and 79% reduction in applicationDB**

*It is important to note that since the sample is taken on height 44335 and not current height 65K+,
we’ll likely see even more reduction onward - as the applicationDB is the fastest growing in the current iteration.*

### Conclusion:

- [x] Simpler Codebase
- [x] Non-Consensus Breaking
- [x] Comparable syncing performance
- [x] Comparable disk usage
- [x] Faster historical queries
- [x] Faster current queries
- [x] Smaller disk usage
- [ ] Better syncing performance
- [x] Slower disk usage rate

### Release:
- Only one option: full archival (but storage requirements significantly reduced)
- Requires a snapshot for others to download as syncing is still very time consuming
- Non-consensus breaking change and uses leveldb, requires a config.json change

### Learnings:
- Retrofitting the behaviors of IAVL is the biggest challenge in this project
- Working with 200GB datadirs for testing is difficult
- The tradeoffs between all of the *goals* continuously conflict with each other
- Developing against mainnet is the only real tool that exists for debugging the current store and the appHash error is no help
- When you develop a persistence module for a blockchain, you should invest the most engineering effort on getting the storage absolutely optimized before launching








