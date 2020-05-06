package types

import (
	"bytes"
	"encoding/binary"
	"math"
	"math/rand"
	"reflect"

	"golang.org/x/crypto/blake2b"
)

// "HashSum" - A structure to represent the hash and the sum at an index in the merkle sum index tree
type HashSum struct {
	Hash []byte `json:"hash"`
	Sum  uint64 `json:"sum"`
}

// "MerkleProof" - A structure used to verify a leaf of the tree.
type MerkleProof struct {
	Index    int       `json:"index"`
	HashSums []HashSum `json:"hash_sums"`
}

// "MerkleProofs" - Two merkle proof objects (needed for replay attack)
type MerkleProofs [2]MerkleProof

// "Validate" - Verifies the Proof from the leaf/cousin node data, the merkle root, and the Proof object
func (mp MerkleProofs) Validate(root HashSum, leaf, cousin Proof, totalRelays int64) (isValid bool, isReplay bool) {
	// check if levels and total relays is valid
	numOfLevels, valid := levelsIsValid(len(mp[0].HashSums), len(mp[1].HashSums), totalRelays)
	if !valid {
		return false, false
	}
	// verifier is the opposing verification piece to the merkle proofs, with its counterpart, we will verify the tree
	var verifier [2]HashSum
	// convert leaf to hashsum
	verifier[0].Hash = hash(leaf.Hash())
	verifier[0].Sum = sumFromHash(verifier[0].Hash)
	// convert cousin to hashsum
	verifier[1].Hash = hash(cousin.Hash())
	verifier[1].Sum = sumFromHash(verifier[1].Hash)
	// replay attack check -> params (leaf, sibling, cousin, cousinSibling, leafIndex, size of the tree)
	if isReplayAttack(verifier[0], mp[0].HashSums[0], verifier[1], mp[1].HashSums[0], int64(mp[0].Index), totalRelays) {
		return false, true
	}
	// execute the for loop for each level
	for i := 0; i < numOfLevels; i++ {
		if mp[0].Index%2 == 1 { // odd index
			// child sum should be greater than sibling sum
			if verifier[0].Sum <= mp[0].HashSums[i].Sum {
				return false, false
			}
			// calculate the parent sum and store it where the child used to be
			verifier[0].Sum += mp[0].HashSums[i].Sum
			// generate the parent hash and store it where the child used to be
			verifier[0].Hash = parentHash(mp[0].HashSums[i].Hash, verifier[0].Hash, verifier[0].Sum, 2*mp[0].Index-1)
		} else { // even index
			// child sum should be less than sibling sum
			if verifier[0].Sum >= mp[0].HashSums[i].Sum {
				return false, false
			}
			// calculate the parent sum and store it where the child used to be
			verifier[0].Sum += mp[0].HashSums[i].Sum
			// generate the parent hash and store it where the child used to be
			verifier[0].Hash = parentHash(verifier[0].Hash, mp[0].HashSums[i].Hash, verifier[0].Sum, 2*mp[0].Index+1)
		}
		if mp[1].Index%2 == 1 { // odd index
			// (cousin) child sum should be greater than sibling sum
			if verifier[1].Sum <= mp[1].HashSums[i].Sum {
				return false, false
			}
			// calculate the parent sum and store it where the child used to be
			verifier[1].Sum += mp[1].HashSums[i].Sum
			// generate the parent hash and store it where the child used to be
			verifier[1].Hash = parentHash(mp[1].HashSums[i].Hash, verifier[1].Hash, verifier[1].Sum, 2*mp[1].Index-1)
		} else {
			// (cousin) child sum should be less than sibling sum
			if verifier[1].Sum >= mp[1].HashSums[i].Sum {
				return false, false
			}
			// calculate the parent sum and store it where the child used to be
			verifier[1].Sum += mp[1].HashSums[i].Sum
			// generate the parent hash and store it where the child used to be
			verifier[1].Hash = parentHash(verifier[1].Hash, mp[1].HashSums[i].Hash, verifier[1].Sum, 2*mp[1].Index+1)
		}
		// half the indices as we are going up one level
		mp[0].Index /= 2
		mp[1].Index /= 2
	}
	// ensure root == verification for leaf and cousin
	return reflect.DeepEqual(root, verifier[0]) && reflect.DeepEqual(root, verifier[1]), false
}

// "GenerateProofs" - Generates the merkle Proof object from the leaf node data and the index
func GenerateProofs(p []Proof, index int) (merkleProofs MerkleProofs, cousinIndex int) {
	data, _ := sortAndStructure(p)
	dataCopy := make([]HashSum, len(data))
	// Copy from the original map to the target map
	copy(dataCopy, data)
	// calculate cousin index
	cousinIndex = getCousinIndex(len(p), index)
	// generate Proof for leaf
	merkleProofs[0] = merkleProof(data, index, &MerkleProof{})
	// reset leaf index
	merkleProofs[0].Index = index
	// generate Proof for cousin
	merkleProofs[1] = merkleProof(dataCopy, cousinIndex, &MerkleProof{})
	// reset cousin index
	merkleProofs[1].Index = cousinIndex
	// return merkleProofs object
	return
}

// "GenerateRoot" - generates the merkle root from leaf node data
func GenerateRoot(data []Proof) (r HashSum, sortedData []Proof) {
	// structure the leafs
	d, sortedProofs := sortAndStructure(data)
	// call the root function and return
	return root(d), sortedProofs
}

// "sortAndStructure" - takes Proof data, sorts, and structures them as a `balanced` merkle tree
func sortAndStructure(relayProofs []Proof) (d []HashSum, sortedProofs []Proof) {
	// we need a tree of proper length. Get the # of relayProofs
	numberOfProofs := len(relayProofs)
	properLength := nextPowerOfTwo(uint(numberOfProofs))
	// initialize the data
	data := make([]HashSum, properLength)
	// first, let's tHash the data
	for i, p := range relayProofs {
		// save the hash and sum of the Proof in the new tree slice
		data[i].Hash = hash(p.Hash()) // todo should this be hash with signature for RelayProofs? // todo remove double hash
		data[i].Sum = sumFromHash(data[i].Hash)
	}
	// for the rest, add the max uint32
	for i := numberOfProofs; i < int(properLength); i++ {
		data[i] = HashSum{
			Hash: Hash([]byte("0")),
			Sum:  uint64(math.MaxUint32),
		}
	}
	proofs := make([]Proof, int(properLength)-numberOfProofs)
	relayProofs = append(relayProofs, proofs...)
	// sort the slice based on the numerical value of the tHash data
	data, relayProofs = quickSort(data, relayProofs)
	return data, relayProofs[:numberOfProofs]
}

// "root" - Generates the root (highest level) from the hash sum data recursively
// CONTRACT: dataLength must be > 1 or this breaks
func root(data []HashSum) HashSum {
	data, atRoot := levelUp(data)
	if !atRoot {
		// if not at root continue to level up
		root(data)
	}
	// if at root return
	return data[0]
}

// "merkleProof" - recursive Proof function that generates the Proof object one level at a time
func merkleProof(data []HashSum, index int, p *MerkleProof) MerkleProof {
	if index%2 == 1 { // odd index so sibling to the left
		p.HashSums = append(p.HashSums, data[index-1])
	} else { // even index so sibling to the right
		p.HashSums = append(p.HashSums, data[index+1])
	}
	data, atRoot := levelUp(data)
	if !atRoot {
		// next level Entropy = previous index / 2 (
		merkleProof(data, index/2, p)
	}
	return *p
}

// "levelUp" - takes the previous level data and converts it to the next level data
func levelUp(data []HashSum) (nextLevelData []HashSum, atRoot bool) {
	for i, d := range data {
		// if odd element, skip
		if i%2 == 1 {
			continue
		}
		// calculate the sum
		data[i/2].Sum = d.Sum + data[i+1].Sum
		// calculate the parent hash
		data[i/2].Hash = parentHash(d.Hash, data[i+1].Hash, data[i/2].Sum, 2*i+1)
	}
	// check to see if at root
	dataLen := len(data) / 2
	if dataLen == 1 {
		return data[:dataLen], true
	}
	return data[:dataLen], false
}

// "isReplayAttack" - Check for replay attack by comparing the order and value of a leaf, the sibling, the cousin, and the cousins sibling
func isReplayAttack(leaf, sibling, cousin, cousinsSibling HashSum, leafIndex int64, treeSize int64) bool {
	// check equality among all leaves
	if bytes.Equal(leaf.Hash, sibling.Hash) ||
		bytes.Equal(leaf.Hash, cousin.Hash) ||
		bytes.Equal(leaf.Hash, cousinsSibling.Hash) ||
		bytes.Equal(sibling.Hash, cousin.Hash) ||
		bytes.Equal(sibling.Hash, cousinsSibling.Hash) ||
		bytes.Equal(cousin.Hash, sibling.Hash) {
		return true
	}
	if leafIndex == 0 {
		// if leaf is at the beginning of the tree, the order is leaf -> sibling -> cousin -> sibling cousin
		return !(leaf.Sum < sibling.Sum && sibling.Sum < cousin.Sum && cousin.Sum < cousinsSibling.Sum)
	} else if leafIndex == treeSize-1 {
		// if leaf is at the end of the tree
		if leafIndex%2 == 0 {
			// if even index at the end, the order is cousin sibling -> cousin -> leaf -> (sibling is filler value)
			return !(cousinsSibling.Sum < cousin.Sum && cousin.Sum < leaf.Sum)
		}
		// if odd index a the end, the order is sibling cousin -> cousin -> sibling -> leaf
		return !(cousinsSibling.Sum < cousin.Sum && cousin.Sum < sibling.Sum && sibling.Sum < leaf.Sum)
	} else {
		// if the leaf is not at the beginning or the end
		if leafIndex%2 == 1 {
			// leaf has odd index so order is sibling -> leaf -> cousin -> cousinSibling
			return !(sibling.Sum < leaf.Sum && leaf.Sum < cousin.Sum && cousin.Sum < cousinsSibling.Sum)
		} else {
			// odd index so order is Cousinsibling -> Cousin -> leaf -> sibling
			return !(cousinsSibling.Sum < cousin.Sum && cousin.Sum < leaf.Sum && leaf.Sum < sibling.Sum)
		}
	}
}

// "getCousinIndex" - Retrieves the index of the cousin by the leaf index
func getCousinIndex(dataLength, leafIndex int) (cousinIndex int) {
	if leafIndex == 0 {
		// beginning so return cousin of sibling as leafIndex
		return 2
	}
	end := dataLength - 1
	if leafIndex == end {
		// at end of tree so return cousin of sibling as leafIndex
		if leafIndex%2 == 0 {
			// if even at the end, cousin is one to the left
			return end - 1
		}
		// if odd leafIndex at the end, cousin is two to the left
		return end - 2
	}
	if leafIndex%2 == 1 {
		// odd leafIndex so cousin to the right
		return leafIndex + 1
	} else {
		// even leafIndex so cousin to the left
		return leafIndex - 1
	}
}

// "levelIsValid" - Ensure that the number of levels in the relayProof is valid
func levelsIsValid(leafNumOfLevels, cousinNumOfLevels int, totalRelays int64) (numOfLevels int, isValid bool) {
	// only accept merkle proofs for more than 4 relays
	if totalRelays < 5 {
		return leafNumOfLevels, false
	}
	if leafNumOfLevels != cousinNumOfLevels {
		return leafNumOfLevels, false
	}
	return leafNumOfLevels, nextPowerOfTwo(uint(totalRelays)) == uint(math.Pow(2, float64(leafNumOfLevels)))
}

// "nextPowrOfTwo" - Computes the next power of 2 given an u-integer
func nextPowerOfTwo(v uint) uint {
	v--
	v |= v >> 1
	v |= v >> 2
	v |= v >> 4
	v |= v >> 8
	v |= v >> 16
	v++
	return v
}

// "hash" - the hash function used in the merkle tree
func hash(data []byte) []byte {
	hash := blake2b.Sum256((data))
	return hash[:]
}

// "parentHash" - Compute the hash of the parent by hashing the hashes, sum and parent
func parentHash(hash1, hash2 []byte, sum uint64, parentIndex int) []byte {
	return hash(append(append(append(hash1, hash2...), uint64ToBytes(sum)...), uint64ToBytes(uint64(parentIndex))...))
}

// "sumFromHash" - get leaf sum from hash
func sumFromHash(hash []byte) uint64 {
	return binary.LittleEndian.Uint64(append(hash[:3], make([]byte, 5)...))
}

// "uint64ToBytes" - convert the uint64 to bytes
func uint64ToBytes(a uint64) (bz []byte) {
	bz = make([]byte, 8)
	binary.LittleEndian.PutUint64(bz, a)
	return
}

// "quickSort" - Sort the hash sum and the proofs by hash sum
// CONTRACT hash sum and proofs must be the same length
func quickSort(hs []HashSum, p []Proof) ([]HashSum, []Proof) {
	// get the length of the hash sums
	hsLen := len(hs)
	if hsLen <= 1 {
		return hs, p
	}
	// set left to zero and right to end index
	l, r := 0, hsLen-1
	// generate random pivot
	pivot := rand.Int() % hsLen
	// switch the two
	hs[pivot], hs[r] = hs[r], hs[pivot]
	// duplicate behavior for the proof
	p[pivot], p[r] = p[r], p[pivot]
	// loop through each item and compare
	for i := range hs {
		if hs[i].Sum < hs[r].Sum {
			// switch the two
			hs[i], hs[l] = hs[l], hs[i]
			// duplicate behavior for the proof
			p[i], p[l] = p[l], p[i]
			// increment left
			l++
		}
	}
	// switch the two
	hs[l], hs[r] = hs[r], hs[l]
	// duplicate behavior for the proof
	p[l], p[r] = p[r], p[l]
	// recursive quicksort
	quickSort(hs[:l], p[:l])
	quickSort(hs[l+1:], p[l+1:])
	return hs, p
}
