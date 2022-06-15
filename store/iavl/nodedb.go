package iavl

import (
	"bytes"
	"container/list"
	"encoding/binary"
	"fmt"
	"github.com/tendermint/go-amino"
	"io"
	"math"
	"sort"
	"sync"

	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/crypto/tmhash"
	dbm "github.com/tendermint/tm-db"
)

const (
	int64Size = 8
	hashSize  = tmhash.Size
)

var (
	// All node keys are prefixed with the byte 'n'. This ensures no collision is
	// possible with the other keys, and makes them easier to traverse. They are indexed by the node hash.
	nodeKeyFormat = NewKeyFormat('n', hashSize) // n<hash>

	// Orphans are keyed in the database by their expected lifetime.
	// The first number represents the *last* version at which the orphan needs
	// to exist, while the second number represents the *earliest* version at
	// which it is expected to exist - which starts out by being the version
	// of the node being orphaned.
	orphanKeyFormat = NewKeyFormat('o', int64Size, int64Size, hashSize) // o<last-version><first-version><hash>

	// Root nodes are indexed separately by their version
	rootKeyFormat = NewKeyFormat('r', int64Size) // r<version>
)

type nodeDB struct {
	mtx            sync.Mutex       // Read/write lock.
	db             dbm.DB           // Persistent node storage.
	batch          dbm.Batch        // Batched writing buffer.
	opts           *Options         // Options to customize for pruning/writing
	versionReaders map[int64]uint32 // Number of active version readers

	latestVersion  int64
	nodeCache      map[string]*list.Element // Node cache.
	nodeCacheSize  int                      // Node cache size limit in elements.
	nodeCacheQueue *list.List               // LRU queue of cache elements. Used for deletion.
}

func newNodeDB(db dbm.DB, cacheSize int, opts *Options) *nodeDB {
	if opts == nil {
		opts = DefaultOptions()
	}
	return &nodeDB{
		db:             db,
		batch:          db.NewBatch(),
		opts:           opts,
		latestVersion:  0, // initially invalid
		nodeCache:      make(map[string]*list.Element),
		nodeCacheSize:  cacheSize,
		nodeCacheQueue: list.New(),
		versionReaders: make(map[int64]uint32, 8),
	}
}

// GetNode gets a node from memory or disk. If it is an inner node, it does not
// load its children.
func (ndb *nodeDB) GetNode(hash []byte) *Node {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()

	if len(hash) == 0 {
		panic("nodeDB.GetNode() requires hash")
	}

	// Check the cache.
	if elem, ok := ndb.nodeCache[string(hash)]; ok {
		// Already exists. Move to back of nodeCacheQueue.
		ndb.nodeCacheQueue.MoveToBack(elem)
		return elem.Value.(*Node)
	}

	// Doesn't exist, load.
	buf, err := ndb.db.Get(ndb.nodeKey(hash))
	if err != nil {
		panic(fmt.Sprintf("can't get node %X: %v", hash, err))
	}
	if buf == nil {
		panic(fmt.Sprintf("Value missing for hash %x corresponding to nodeKey %x", hash, ndb.nodeKey(hash)))
	}

	node, err := MakeNode(buf)
	if err != nil {
		panic(fmt.Sprintf("Error reading Node. bytes: %x, error: %v", buf, err))
	}

	node.hash = hash
	node.persisted = true
	ndb.cacheNode(node)

	return node
}

// SaveNode saves a node to disk.
func (ndb *nodeDB) SaveNode(node *Node) {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()

	if node.hash == nil {
		panic("Expected to find node.hash, but none found.")
	}
	if node.persisted {
		panic("Shouldn't be calling save on an already persisted node.")
	}

	// Save node bytes to db.
	var buf bytes.Buffer
	buf.Grow(node.aminoSize())

	if err := node.writeBytes(&buf); err != nil {
		panic(err)
	}

	ndb.batch.Set(ndb.nodeKey(node.hash), buf.Bytes())
	debug("BATCH SAVE %X %p\n", node.hash, node)
	node.persisted = true
	ndb.cacheNode(node)
}

// Has checks if a hash exists in the database.
func (ndb *nodeDB) Has(hash []byte) (bool, error) {
	key := ndb.nodeKey(hash)

	if ldb, ok := ndb.db.(*dbm.GoLevelDB); ok {
		exists, err := ldb.DB().Has(key, nil)
		if err != nil {
			return false, err
		}
		return exists, nil
	}
	value, err := ndb.db.Get(key)
	if err != nil {
		return false, err
	}

	return value != nil, nil
}

// SaveBranch saves the given node and all of its descendants.
// NOTE: This function clears leftNode/rigthNode recursively and
// calls _hash() on the given node.
// TODO refactor, maybe use hashWithCount() but provide a callback.
func (ndb *nodeDB) SaveBranch(node *Node) []byte {
	if node.persisted {
		return node.hash
	}

	if node.leftNode != nil {
		node.leftHash = ndb.SaveBranch(node.leftNode)
	}
	if node.rightNode != nil {
		node.rightHash = ndb.SaveBranch(node.rightNode)
	}

	node._hash()
	ndb.SaveNode(node)

	node.leftNode = nil
	node.rightNode = nil

	return node.hash
}

// DeleteVersion deletes a tree version from disk.
func (ndb *nodeDB) DeleteVersion(version int64, checkLatestVersion bool) error {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()

	if ndb.versionReaders[version] > 0 {
		return errors.Errorf("unable to delete version %v, it has %v active readers", version, ndb.versionReaders[version])
	}

	ndb.deleteOrphans(version)
	ndb.deleteRoot(version, checkLatestVersion)
	return nil
}

// DeleteVersionsFrom permanently deletes all tree versions from the given version upwards.
func (ndb *nodeDB) DeleteVersionsFrom(version int64) error {
	latest := ndb.getLatestVersion()
	if latest < version {
		return nil
	}
	root, err := ndb.getRoot(latest)
	if err != nil {
		return err
	}
	if root == nil {
		return errors.Errorf("root for version %v not found", latest)
	}

	for v, r := range ndb.versionReaders {
		if v >= version && r != 0 {
			return errors.Errorf("unable to delete version %v with %v active readers", v, r)
		}
	}

	// First, delete all active nodes in the current (latest) version whose node version is after
	// the given version.
	err = ndb.deleteNodesFrom(version, root)
	if err != nil {
		return err
	}

	// Next, delete orphans:
	// - Remove orphan entries *and referred nodes* with fromVersion >= version
	// - Remove orphan entries with toVersion >= version-1 (since orphans at latest are not orphans)
	ndb.traverseOrphans(func(key, hash []byte) {
		var fromVersion, toVersion int64
		orphanKeyFormat.Scan(key, &toVersion, &fromVersion)

		if fromVersion >= version {
			ndb.batch.Delete(key)
			ndb.batch.Delete(ndb.nodeKey(hash))
			ndb.uncacheNode(hash)
		} else if toVersion >= version-1 {
			ndb.batch.Delete(key)
		}
	})

	// Finally, delete the version root entries
	ndb.traverseRange(rootKeyFormat.Key(version), rootKeyFormat.Key(math.MaxInt64), func(k, v []byte) {
		ndb.batch.Delete(k)
	})

	return nil
}

// deleteNodesFrom deletes the given node and any descendants that have versions after the given
// (inclusive). It is mainly used via LoadVersionForOverwriting, to delete the current version.
func (ndb *nodeDB) deleteNodesFrom(version int64, hash []byte) error {
	if len(hash) == 0 {
		return nil
	}

	node := ndb.GetNode(hash)
	if node.leftHash != nil {
		if err := ndb.deleteNodesFrom(version, node.leftHash); err != nil {
			return err
		}
	}
	if node.rightHash != nil {
		if err := ndb.deleteNodesFrom(version, node.rightHash); err != nil {
			return err
		}
	}

	if node.version >= version {
		ndb.batch.Delete(ndb.nodeKey(hash))
		ndb.uncacheNode(hash)
	}

	return nil
}

// Saves orphaned nodes to disk under a special prefix.
// version: the new version being saved.
// orphans: the orphan nodes created since version-1
func (ndb *nodeDB) SaveOrphans(version int64, orphans map[string]int64) {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()

	toVersion := ndb.getPreviousVersion(version)
	for hash, fromVersion := range orphans {
		debug("SAVEORPHAN %v-%v %X\n", fromVersion, toVersion, hash)
		ndb.saveOrphan([]byte(hash), fromVersion, toVersion)
	}
}

// Saves a single orphan to disk.
func (ndb *nodeDB) saveOrphan(hash []byte, fromVersion, toVersion int64) {
	if fromVersion > toVersion {
		panic(fmt.Sprintf("Orphan expires before it comes alive.  %d > %d", fromVersion, toVersion))
	}
	key := ndb.orphanKey(fromVersion, toVersion, hash)
	ndb.batch.Set(key, hash)
}

// deleteOrphans deletes orphaned nodes from disk, and the associated orphan
// entries.
func (ndb *nodeDB) deleteOrphans(version int64) {
	// Will be zero if there is no previous version.
	predecessor := ndb.getPreviousVersion(version)

	// Traverse orphans with a lifetime ending at the version specified.
	// TODO optimize.
	ndb.traverseOrphansVersion(version, func(key, hash []byte) {
		var fromVersion, toVersion int64

		// See comment on `orphanKeyFmt`. Note that here, `version` and
		// `toVersion` are always equal.
		orphanKeyFormat.Scan(key, &toVersion, &fromVersion)

		// Remove orphan key and reverse-lookup key.
		ndb.batch.Delete(key)

		// If there is no predecessor, or the predecessor is earlier than the
		// beginning of the lifetime (ie: negative lifetime), or the lifetime
		// spans a single version and that version is the one being deleted, we
		// can delete the orphan.  Otherwise, we shorten its lifetime, by
		// moving its endpoint to the previous version.
		if predecessor < fromVersion || fromVersion == toVersion {
			debug("DELETE predecessor:%v fromVersion:%v toVersion:%v %X\n", predecessor, fromVersion, toVersion, hash)
			ndb.batch.Delete(ndb.nodeKey(hash))
			ndb.uncacheNode(hash)
		} else {
			debug("MOVE predecessor:%v fromVersion:%v toVersion:%v %X\n", predecessor, fromVersion, toVersion, hash)
			ndb.saveOrphan(hash, fromVersion, predecessor)
		}
	})
}

func (ndb *nodeDB) nodeKey(hash []byte) []byte {
	return nodeKeyFormat.KeyBytes(hash)
}

func (ndb *nodeDB) orphanKey(fromVersion, toVersion int64, hash []byte) []byte {
	return orphanKeyFormat.Key(toVersion, fromVersion, hash)
}

func (ndb *nodeDB) rootKey(version int64) []byte {
	return rootKeyFormat.Key(version)
}

func (ndb *nodeDB) getLatestVersion() int64 {
	if ndb.latestVersion == 0 {
		ndb.latestVersion = ndb.getPreviousVersion(1<<63 - 1)
	}
	return ndb.latestVersion
}

func (ndb *nodeDB) updateLatestVersion(version int64) {
	if ndb.latestVersion < version {
		ndb.latestVersion = version
	}
}

func (ndb *nodeDB) resetLatestVersion(version int64) {
	ndb.latestVersion = version
}

func (ndb *nodeDB) getPreviousVersion(version int64) int64 {
	itr, err := ndb.db.ReverseIterator(
		rootKeyFormat.Key(1),
		rootKeyFormat.Key(version),
	)
	if err != nil {
		panic(err)
	}
	defer itr.Close()

	pversion := int64(-1)
	for ; itr.Valid(); itr.Next() {
		k := itr.Key()
		rootKeyFormat.Scan(k, &pversion)
		return pversion
	}

	return 0
}

// deleteRoot deletes the root entry from disk, but not the node it points to.
func (ndb *nodeDB) deleteRoot(version int64, checkLatestVersion bool) {
	if checkLatestVersion && version == ndb.getLatestVersion() {
		panic("Tried to delete latest version")
	}
	ndb.batch.Delete(ndb.rootKey(version))
}

func (ndb *nodeDB) traverseOrphans(fn func(k, v []byte)) {
	ndb.traversePrefix(orphanKeyFormat.Key(), fn)
}

// Traverse orphans ending at a certain version.
func (ndb *nodeDB) traverseOrphansVersion(version int64, fn func(k, v []byte)) {
	ndb.traversePrefix(orphanKeyFormat.Key(version), fn)
}

// Traverse all keys.
func (ndb *nodeDB) traverse(fn func(key, value []byte)) {
	ndb.traverseRange(nil, nil, fn)
}

// Traverse all keys between a given range (excluding end).
func (ndb *nodeDB) traverseRange(start []byte, end []byte, fn func(k, v []byte)) {
	itr, err := ndb.db.Iterator(start, end)
	if err != nil {
		panic(err)
	}
	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		fn(itr.Key(), itr.Value())
	}
}

// Traverse all keys with a certain prefix.
func (ndb *nodeDB) traversePrefix(prefix []byte, fn func(k, v []byte)) {
	itr, err := dbm.IteratePrefix(ndb.db, prefix)
	if err != nil {
		panic(err)
	}
	defer itr.Close()

	for ; itr.Valid(); itr.Next() {
		fn(itr.Key(), itr.Value())
	}
}

func (ndb *nodeDB) uncacheNode(hash []byte) {
	if elem, ok := ndb.nodeCache[string(hash)]; ok {
		ndb.nodeCacheQueue.Remove(elem)
		delete(ndb.nodeCache, string(hash))
	}
}

// Add a node to the cache and pop the least recently used node if we've
// reached the cache size limit.
func (ndb *nodeDB) cacheNode(node *Node) {
	elem := ndb.nodeCacheQueue.PushBack(node)
	ndb.nodeCache[string(node.hash)] = elem

	if ndb.nodeCacheQueue.Len() > ndb.nodeCacheSize {
		oldest := ndb.nodeCacheQueue.Front()
		hash := ndb.nodeCacheQueue.Remove(oldest).(*Node).hash
		delete(ndb.nodeCache, string(hash))
	}
}

// Write to disk.
func (ndb *nodeDB) Commit() error {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()

	var err error
	if ndb.opts.Sync {
		err = ndb.batch.WriteSync()
	} else {
		err = ndb.batch.Write()
	}
	if err != nil {
		return errors.Wrap(err, "failed to write batch")
	}

	ndb.batch.Close()
	ndb.batch = ndb.db.NewBatch()

	return nil
}

func (ndb *nodeDB) getRoot(version int64) ([]byte, error) {
	return ndb.db.Get(ndb.rootKey(version))
}

func (ndb *nodeDB) getRoots() (map[int64][]byte, error) {
	roots := map[int64][]byte{}

	ndb.traversePrefix(rootKeyFormat.Key(), func(k, v []byte) {
		var version int64
		rootKeyFormat.Scan(k, &version)
		roots[version] = v
	})
	return roots, nil
}

// SaveRoot creates an entry on disk for the given root, so that it can be
// loaded later.
func (ndb *nodeDB) SaveRoot(root *Node, version int64) error {
	if len(root.hash) == 0 {
		panic("SaveRoot: root hash should not be empty")
	}
	return ndb.saveRoot(root.hash, version)
}

// SaveEmptyRoot creates an entry on disk for an empty root.
func (ndb *nodeDB) SaveEmptyRoot(version int64) error {
	return ndb.saveRoot([]byte{}, version)
}

func (ndb *nodeDB) saveRoot(hash []byte, version int64) error {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()

	if version != ndb.getLatestVersion()+1 {
		return fmt.Errorf("must save consecutive versions; expected %d, got %d", ndb.getLatestVersion()+1, version)
	}

	ndb.batch.Set(ndb.rootKey(version), hash)
	ndb.updateLatestVersion(version)
	return nil
}

func (ndb *nodeDB) incrVersionReaders(version int64) {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()
	ndb.versionReaders[version]++
}

func (ndb *nodeDB) decrVersionReaders(version int64) {
	ndb.mtx.Lock()
	defer ndb.mtx.Unlock()
	if ndb.versionReaders[version] > 0 {
		ndb.versionReaders[version]--
	}
}

////////////////// Utility and test functions /////////////////////////////////

func (ndb *nodeDB) leafNodes() []*Node {
	leaves := []*Node{}

	ndb.traverseNodes(func(hash []byte, node *Node) {
		if node.isLeaf() {
			leaves = append(leaves, node)
		}
	})
	return leaves
}

func (ndb *nodeDB) nodes() []*Node {
	nodes := []*Node{}

	ndb.traverseNodes(func(hash []byte, node *Node) {
		nodes = append(nodes, node)
	})
	return nodes
}

func (ndb *nodeDB) orphans() [][]byte {
	orphans := [][]byte{}

	ndb.traverseOrphans(func(k, v []byte) {
		orphans = append(orphans, v)
	})
	return orphans
}

func (ndb *nodeDB) roots() map[int64][]byte {
	roots, _ := ndb.getRoots()
	return roots
}

// Not efficient.
// NOTE: DB cannot implement Size() because
// mutations are not always synchronous.
func (ndb *nodeDB) size() int {
	size := 0
	ndb.traverse(func(k, v []byte) {
		size++
	})
	return size
}

func (ndb *nodeDB) traverseNodes(fn func(hash []byte, node *Node)) {
	nodes := []*Node{}

	ndb.traversePrefix(nodeKeyFormat.Key(), func(key, value []byte) {
		node, err := MakeNode(value)
		if err != nil {
			panic(fmt.Sprintf("Couldn't decode node from database: %v", err))
		}
		nodeKeyFormat.Scan(key, &node.hash)
		nodes = append(nodes, node)
	})

	sort.Slice(nodes, func(i, j int) bool {
		return bytes.Compare(nodes[i].key, nodes[j].key) < 0
	})

	for _, n := range nodes {
		fn(n.hash, n)
	}
}

func (ndb *nodeDB) String() string {
	var str string
	index := 0

	ndb.traversePrefix(rootKeyFormat.Key(), func(key, value []byte) {
		str += fmt.Sprintf("%s: %x\n", string(key), value)
	})
	str += "\n"

	ndb.traverseOrphans(func(key, value []byte) {
		str += fmt.Sprintf("%s: %x\n", string(key), value)
	})
	str += "\n"

	ndb.traverseNodes(func(hash []byte, node *Node) {
		switch {
		case len(hash) == 0:
			str += "<nil>\n"
		case node == nil:
			str += fmt.Sprintf("%s%40x: <nil>\n", nodeKeyFormat.Prefix(), hash)
		case node.value == nil && node.height > 0:
			str += fmt.Sprintf("%s%40x: %s   %-16s h=%d version=%d\n",
				nodeKeyFormat.Prefix(), hash, node.key, "", node.height, node.version)
		default:
			str += fmt.Sprintf("%s%40x: %s = %-16s h=%d version=%d\n",
				nodeKeyFormat.Prefix(), hash, node.key, node.value, node.height, node.version)
		}
		index++
	})
	return "-" + "\n" + str + "-"
}

// Node represents a node in a Tree.
type Node struct {
	key       []byte
	value     []byte
	hash      []byte
	leftHash  []byte
	rightHash []byte
	version   int64
	size      int64
	leftNode  *Node
	rightNode *Node
	height    int8
	persisted bool
}

// NewNode returns a new node from a key, value and version.
func NewNode(key []byte, value []byte, version int64) *Node {
	return &Node{
		key:     key,
		value:   value,
		height:  0,
		size:    1,
		version: version,
	}
}

// MakeNode constructs an *Node from an encoded byte slice.
//
// The new node doesn't have its hash saved or set. The caller must set it
// afterwards.
func MakeNode(buf []byte) (*Node, error) {

	// Read node header (height, size, version, key).
	height, n, cause := amino.DecodeInt8(buf)
	if cause != nil {
		return nil, errors.Wrap(cause, "decoding node.height")
	}
	buf = buf[n:]

	size, n, cause := amino.DecodeVarint(buf)
	if cause != nil {
		return nil, errors.Wrap(cause, "decoding node.size")
	}
	buf = buf[n:]

	ver, n, cause := amino.DecodeVarint(buf)
	if cause != nil {
		return nil, errors.Wrap(cause, "decoding node.version")
	}
	buf = buf[n:]

	key, n, cause := amino.DecodeByteSlice(buf)
	if cause != nil {
		return nil, errors.Wrap(cause, "decoding node.key")
	}
	buf = buf[n:]

	node := &Node{
		height:  height,
		size:    size,
		version: ver,
		key:     key,
	}

	// Read node body.

	if node.isLeaf() {
		val, _, cause := amino.DecodeByteSlice(buf)
		if cause != nil {
			return nil, errors.Wrap(cause, "decoding node.value")
		}
		node.value = val
	} else { // Read children.
		leftHash, n, cause := amino.DecodeByteSlice(buf)
		if cause != nil {
			return nil, errors.Wrap(cause, "deocding node.leftHash")
		}
		buf = buf[n:]

		rightHash, _, cause := amino.DecodeByteSlice(buf)
		if cause != nil {
			return nil, errors.Wrap(cause, "decoding node.rightHash")
		}
		node.leftHash = leftHash
		node.rightHash = rightHash
	}
	return node, nil
}

// String returns a string representation of the node.
func (node *Node) String() string {
	hashstr := "<no hash>"
	if len(node.hash) > 0 {
		hashstr = fmt.Sprintf("%X", node.hash)
	}
	return fmt.Sprintf("Node{%s:%s@%d %X;%X}#%s",
		node.key,
		node.value,
		node.version,
		node.leftHash, node.rightHash,
		hashstr)
}

// clone creates a shallow copy of a node with its hash set to nil.
func (node *Node) clone(version int64) *Node {
	if node.isLeaf() {
		panic("Attempt to copy a leaf node")
	}
	return &Node{
		key:       node.key,
		height:    node.height,
		version:   version,
		size:      node.size,
		hash:      nil,
		leftHash:  node.leftHash,
		leftNode:  node.leftNode,
		rightHash: node.rightHash,
		rightNode: node.rightNode,
		persisted: false,
	}
}

func (node *Node) isLeaf() bool {
	return node.height == 0
}

// Check if the node has a descendant with the given key.
func (node *Node) has(t *ImmutableTree, key []byte) (has bool) {
	if bytes.Equal(node.key, key) {
		return true
	}
	if node.isLeaf() {
		return false
	}
	if bytes.Compare(key, node.key) < 0 {
		return node.getLeftNode(t).has(t, key)
	}
	return node.getRightNode(t).has(t, key)
}

// Get a key under the node.
func (node *Node) get(t *ImmutableTree, key []byte) (index int64, value []byte) {
	if node.isLeaf() {
		switch bytes.Compare(node.key, key) {
		case -1:
			return 1, nil
		case 1:
			return 0, nil
		default:
			return 0, node.value
		}
	}

	if bytes.Compare(key, node.key) < 0 {
		return node.getLeftNode(t).get(t, key)
	}
	rightNode := node.getRightNode(t)
	index, value = rightNode.get(t, key)
	index += node.size - rightNode.size
	return index, value
}

func (node *Node) getByIndex(t *ImmutableTree, index int64) (key []byte, value []byte) {
	if node.isLeaf() {
		if index == 0 {
			return node.key, node.value
		}
		return nil, nil
	}
	// TODO: could improve this by storing the
	// sizes as well as left/right hash.
	leftNode := node.getLeftNode(t)

	if index < leftNode.size {
		return leftNode.getByIndex(t, index)
	}
	return node.getRightNode(t).getByIndex(t, index-leftNode.size)
}

// Computes the hash of the node without computing its descendants. Must be
// called on nodes which have descendant node hashes already computed.
func (node *Node) _hash() []byte {
	if node.hash != nil {
		return node.hash
	}

	h := tmhash.New()
	buf := new(bytes.Buffer)
	if err := node.writeHashBytes(buf); err != nil {
		panic(err)
	}
	_, err := h.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}
	node.hash = h.Sum(nil)

	return node.hash
}

// Hash the node and its descendants recursively. This usually mutates all
// descendant nodes. Returns the node hash and number of nodes hashed.
func (node *Node) hashWithCount() ([]byte, int64) {
	if node.hash != nil {
		return node.hash, 0
	}

	h := tmhash.New()
	buf := new(bytes.Buffer)
	hashCount, err := node.writeHashBytesRecursively(buf)
	if err != nil {
		panic(err)
	}
	_, err = h.Write(buf.Bytes())
	if err != nil {
		panic(err)
	}
	node.hash = h.Sum(nil)

	return node.hash, hashCount + 1
}

// validate validates the node contents
func (node *Node) validate() error {
	if node == nil {
		return errors.New("node cannot be nil")
	}
	if node.key == nil {
		return errors.New("key cannot be nil")
	}
	if node.version <= 0 {
		return errors.New("version must be greater than 0")
	}
	if node.height < 0 {
		return errors.New("height cannot be less than 0")
	}
	if node.size < 1 {
		return errors.New("size must be at least 1")
	}

	if node.height == 0 {
		// Leaf nodes
		if node.value == nil {
			return errors.New("value cannot be nil for leaf node")
		}
		if node.leftHash != nil || node.leftNode != nil || node.rightHash != nil || node.rightNode != nil {
			return errors.New("leaf node cannot have children")
		}
		if node.size != 1 {
			return errors.New("leaf nodes must have size 1")
		}
	} else {
		// Inner nodes
		if node.value != nil {
			return errors.New("value must be nil for non-leaf node")
		}
		if node.leftHash == nil && node.rightHash == nil {
			return errors.New("inner node must have children")
		}
	}
	return nil
}

// Writes the node's hash to the given io.Writer. This function expects
// child hashes to be already set.
func (node *Node) writeHashBytes(w io.Writer) error {
	err := amino.EncodeInt8(w, node.height)
	if err != nil {
		return errors.Wrap(err, "writing height")
	}
	err = amino.EncodeVarint(w, node.size)
	if err != nil {
		return errors.Wrap(err, "writing size")
	}
	err = amino.EncodeVarint(w, node.version)
	if err != nil {
		return errors.Wrap(err, "writing version")
	}

	// Key is not written for inner nodes, unlike writeBytes.

	if node.isLeaf() {
		err = amino.EncodeByteSlice(w, node.key)
		if err != nil {
			return errors.Wrap(err, "writing key")
		}
		// Indirection needed to provide proofs without values.
		// (e.g. ProofLeafNode.ValueHash)
		valueHash := tmhash.Sum(node.value)
		err = amino.EncodeByteSlice(w, valueHash)
		if err != nil {
			return errors.Wrap(err, "writing value")
		}
	} else {
		if node.leftHash == nil || node.rightHash == nil {
			panic("Found an empty child hash")
		}
		err = amino.EncodeByteSlice(w, node.leftHash)
		if err != nil {
			return errors.Wrap(err, "writing left hash")
		}
		err = amino.EncodeByteSlice(w, node.rightHash)
		if err != nil {
			return errors.Wrap(err, "writing right hash")
		}
	}

	return nil
}

// Writes the node's hash to the given io.Writer.
// This function has the side-effect of calling hashWithCount.
func (node *Node) writeHashBytesRecursively(w io.Writer) (hashCount int64, err error) {
	if node.leftNode != nil {
		leftHash, leftCount := node.leftNode.hashWithCount()
		node.leftHash = leftHash
		hashCount += leftCount
	}
	if node.rightNode != nil {
		rightHash, rightCount := node.rightNode.hashWithCount()
		node.rightHash = rightHash
		hashCount += rightCount
	}
	err = node.writeHashBytes(w)

	return
}

func (node *Node) aminoSize() int {
	n := 1 +
		amino.VarintSize(node.size) +
		amino.VarintSize(node.version) +
		amino.ByteSliceSize(node.key)
	if node.isLeaf() {
		n += amino.ByteSliceSize(node.value)
	} else {
		n += amino.ByteSliceSize(node.leftHash) +
			amino.ByteSliceSize(node.rightHash)
	}
	return n
}

// Writes the node as a serialized byte slice to the supplied io.Writer.
func (node *Node) writeBytes(w io.Writer) error {
	cause := amino.EncodeInt8(w, node.height)
	if cause != nil {
		return errors.Wrap(cause, "writing height")
	}
	cause = amino.EncodeVarint(w, node.size)
	if cause != nil {
		return errors.Wrap(cause, "writing size")
	}
	cause = amino.EncodeVarint(w, node.version)
	if cause != nil {
		return errors.Wrap(cause, "writing version")
	}

	// Unlike writeHashBytes, key is written for inner nodes.
	cause = amino.EncodeByteSlice(w, node.key)
	if cause != nil {
		return errors.Wrap(cause, "writing key")
	}

	if node.isLeaf() {
		cause = amino.EncodeByteSlice(w, node.value)
		if cause != nil {
			return errors.Wrap(cause, "writing value")
		}
	} else {
		if node.leftHash == nil {
			panic("node.leftHash was nil in writeBytes")
		}
		cause = amino.EncodeByteSlice(w, node.leftHash)
		if cause != nil {
			return errors.Wrap(cause, "writing left hash")
		}

		if node.rightHash == nil {
			panic("node.rightHash was nil in writeBytes")
		}
		cause = amino.EncodeByteSlice(w, node.rightHash)
		if cause != nil {
			return errors.Wrap(cause, "writing right hash")
		}
	}
	return nil
}

func (node *Node) getLeftNode(t *ImmutableTree) *Node {
	if node.leftNode != nil {
		return node.leftNode
	}
	return t.ndb.GetNode(node.leftHash)
}

func (node *Node) getRightNode(t *ImmutableTree) *Node {
	if node.rightNode != nil {
		return node.rightNode
	}
	return t.ndb.GetNode(node.rightHash)
}

// NOTE: mutates height and size
func (node *Node) calcHeightAndSize(t *ImmutableTree) {
	node.height = maxInt8(node.getLeftNode(t).height, node.getRightNode(t).height) + 1
	node.size = node.getLeftNode(t).size + node.getRightNode(t).size
}

func (node *Node) calcBalance(t *ImmutableTree) int {
	return int(node.getLeftNode(t).height) - int(node.getRightNode(t).height)
}

// traverse is a wrapper over traverseInRange when we want the whole tree
func (node *Node) traverse(t *ImmutableTree, ascending bool, cb func(*Node) bool) bool {
	return node.traverseInRange(t, nil, nil, ascending, false, 0, false, func(node *Node, depth uint8) bool {
		return cb(node)
	})
}

// traversePost is a wrapper over traverseInRange when we want the whole tree post-order
func (node *Node) traversePost(t *ImmutableTree, ascending bool, cb func(*Node) bool) bool {
	return node.traverseInRange(t, nil, nil, ascending, false, 0, true, func(node *Node, depth uint8) bool {
		return cb(node)
	})
}

func (node *Node) traverseInRange(t *ImmutableTree, start, end []byte, ascending bool, inclusive bool, depth uint8, post bool, cb func(*Node, uint8) bool) bool {
	if node == nil {
		return false
	}
	afterStart := start == nil || bytes.Compare(start, node.key) < 0
	startOrAfter := start == nil || bytes.Compare(start, node.key) <= 0
	beforeEnd := end == nil || bytes.Compare(node.key, end) < 0
	if inclusive {
		beforeEnd = end == nil || bytes.Compare(node.key, end) <= 0
	}

	// Run callback per inner/leaf node.
	stop := false
	if !post && (!node.isLeaf() || (startOrAfter && beforeEnd)) {
		stop = cb(node, depth)
		if stop {
			return stop
		}
	}

	if !node.isLeaf() {
		if ascending {
			// check lower nodes, then higher
			if afterStart {
				stop = node.getLeftNode(t).traverseInRange(t, start, end, ascending, inclusive, depth+1, post, cb)
			}
			if stop {
				return stop
			}
			if beforeEnd {
				stop = node.getRightNode(t).traverseInRange(t, start, end, ascending, inclusive, depth+1, post, cb)
			}
		} else {
			// check the higher nodes first
			if beforeEnd {
				stop = node.getRightNode(t).traverseInRange(t, start, end, ascending, inclusive, depth+1, post, cb)
			}
			if stop {
				return stop
			}
			if afterStart {
				stop = node.getLeftNode(t).traverseInRange(t, start, end, ascending, inclusive, depth+1, post, cb)
			}
		}
	}
	if stop {
		return stop
	}

	if post && (!node.isLeaf() || (startOrAfter && beforeEnd)) {
		stop = cb(node, depth)
		if stop {
			return stop
		}
	}

	return stop
}

// Only used in testing...
func (node *Node) lmd(t *ImmutableTree) *Node {
	if node.isLeaf() {
		return node
	}
	return node.getLeftNode(t).lmd(t)
}

// Provides a fixed-width lexicographically sortable []byte key format
type KeyFormat struct {
	prefix byte
	layout []int
	length int
}

// Create a []byte key format based on a single byte prefix and fixed width key segments each of whose length is
// specified by by the corresponding element of layout.
//
// For example, to store keys that could index some objects by a version number and their SHA256 hash using the form:
// 'c<version uint64><hash [32]byte>' then you would define the KeyFormat with:
//
//  var keyFormat = NewKeyFormat('c', 8, 32)
//
// Then you can create a key with:
//
//  func ObjectKey(version uint64, objectBytes []byte) []byte {
//  	hasher := sha256.New()
//  	hasher.Sum(nil)
//  	return keyFormat.Key(version, hasher.Sum(nil))
//  }
func NewKeyFormat(prefix byte, layout ...int) *KeyFormat {
	// For prefix byte
	length := 1
	for _, l := range layout {
		length += l
	}
	return &KeyFormat{
		prefix: prefix,
		layout: layout,
		length: length,
	}
}

// Format the byte segments into the key format - will panic if the segment lengths do not match the layout.
func (kf *KeyFormat) KeyBytes(segments ...[]byte) []byte {
	key := make([]byte, kf.length)
	key[0] = kf.prefix
	n := 1
	for i, s := range segments {
		l := kf.layout[i]
		if len(s) > l {
			panic(fmt.Errorf("length of segment %X provided to KeyFormat.KeyBytes() is longer than the %d bytes "+
				"required by layout for segment %d", s, l, i))
		}
		n += l
		// Big endian so pad on left if not given the full width for this segment
		copy(key[n-len(s):n], s)
	}
	return key[:n]
}

// Format the args passed into the key format - will panic if the arguments passed do not match the length
// of the segment to which they correspond. When called with no arguments returns the raw prefix (useful as a start
// element of the entire keys space when sorted lexicographically).
func (kf *KeyFormat) Key(args ...interface{}) []byte {
	if len(args) > len(kf.layout) {
		panic(fmt.Errorf("keyFormat.Key() is provided with %d args but format only has %d segments",
			len(args), len(kf.layout)))
	}
	segments := make([][]byte, len(args))
	for i, a := range args {
		segments[i] = format(a)
	}
	return kf.KeyBytes(segments...)
}

// Reads out the bytes associated with each segment of the key format from key.
func (kf *KeyFormat) ScanBytes(key []byte) [][]byte {
	segments := make([][]byte, len(kf.layout))
	n := 1
	for i, l := range kf.layout {
		n += l
		if n > len(key) {
			return segments[:i]
		}
		segments[i] = key[n-l : n]
	}
	return segments
}

// Extracts the segments into the values pointed to by each of args. Each arg must be a pointer to int64, uint64, or
// []byte, and the width of the args must match layout.
func (kf *KeyFormat) Scan(key []byte, args ...interface{}) {
	segments := kf.ScanBytes(key)
	if len(args) > len(segments) {
		panic(fmt.Errorf("keyFormat.Scan() is provided with %d args but format only has %d segments in key %X",
			len(args), len(segments), key))
	}
	for i, a := range args {
		scan(a, segments[i])
	}
}

// Return the prefix as a string.
func (kf *KeyFormat) Prefix() string {
	return string([]byte{kf.prefix})
}

func scan(a interface{}, value []byte) {
	switch v := a.(type) {
	case *int64:
		// Negative values will be mapped correctly when read in as uint64 and then type converted
		*v = int64(binary.BigEndian.Uint64(value))
	case *uint64:
		*v = binary.BigEndian.Uint64(value)
	case *[]byte:
		*v = value
	default:
		panic(fmt.Errorf("keyFormat scan() does not support scanning value of type %T: %v", a, a))
	}
}

func format(a interface{}) []byte {
	switch v := a.(type) {
	case uint64:
		return formatUint64(v)
	case int64:
		return formatUint64(uint64(v))
	// Provide formatting from int,uint as a convenience to avoid casting arguments
	case uint:
		return formatUint64(uint64(v))
	case int:
		return formatUint64(uint64(v))
	case []byte:
		return v
	default:
		panic(fmt.Errorf("keyFormat format() does not support formatting value of type %T: %v", a, a))
	}
}

func formatUint64(v uint64) []byte {
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, v)
	return bs
}

func maxInt8(a, b int8) int8 {
	if a > b {
		return a
	}
	return b
}
