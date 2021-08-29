package heightcache

import (
	"github.com/pokt-network/pocket-core/store/types"
	"sort"
)

var _ types.Iterator = &MemoryHeightIterator{}

type MemoryHeightIterator struct {
	dataset    map[string]string
	sortedKeys []string
	curIdx     int
	startIdx   int
	endIdx     int
	start      string
	end        string
	ascending  bool
}

func NewMemoryHeightIterator(dataset map[string]string, start string, end string, sortedKeys []string, ascending bool) *MemoryHeightIterator {
	if start != "" || end != "" {
		if start != "" && end != "" && start > end { // start has to be smaller than end!
			return &MemoryHeightIterator{endIdx: -1, startIdx: 1}
		}
	}
	if start > end {
		tmp := start
		start = end
		end = tmp
	}
	if len(sortedKeys) == 0 {
		sortedKeys = make([]string, 0, len(dataset))
		for k, _ := range dataset {
			sortedKeys = append(sortedKeys, k)
		}
		sort.Strings(sortedKeys)
	}
	startIdx := 0
	if start != "" { // this is a risky assumption -- what's the diff between string([]bytes{}) and (string[]bytes(nil)) ? those are considered smallest and largest by iavl.
		for ; startIdx < len(sortedKeys)-1; startIdx++ {
			if sortedKeys[startIdx] >= start {
				break
			}
		}
	}
	endIdx := len(sortedKeys) - 1
	if end != "" {
		for ; endIdx > 0 && endIdx > startIdx; endIdx-- {
			if sortedKeys[endIdx] <= end {
				break
			}
		}
	}
	curIdx := startIdx
	if !ascending {
		curIdx = endIdx
	}
	// the start string, if not null, should be somewhere
	// the end string, if not null, should be somewhere _after_ the start string
	// curIdx is calculated to be just before the start string
	// start and end strings should be calculated if nil
	// start and end indices of course, too
	return &MemoryHeightIterator{
		dataset:    dataset,
		sortedKeys: sortedKeys,
		curIdx:     curIdx,
		startIdx:   startIdx,
		endIdx:     endIdx,
		start:      start,
		end:        end,
		ascending:  ascending,
	}
}

func (m *MemoryHeightIterator) Domain() (start []byte, end []byte) {
	return []byte(m.start), []byte(m.end)
}

func (m *MemoryHeightIterator) Valid() bool {
	if m.endIdx < m.startIdx || m.curIdx > m.endIdx {
		return false
	}
	if (m.end != "" && m.sortedKeys[m.curIdx] >= m.end) || (m.start != "" && m.sortedKeys[m.curIdx] < m.start) {
		return false
	}
	if m.sortedKeys == nil || m.dataset == nil {
		return false // we closed!!
	}
	if m.curIdx < 0 || m.curIdx > len(m.sortedKeys)-1 {
		return false // out of range!
	}
	return true

}

func (m *MemoryHeightIterator) Next() {
	if !m.Valid() {
		panic("Invalid Iterator.")
	}
	if m.ascending {
		m.curIdx++
	} else {
		m.curIdx--
	}
}

func (m *MemoryHeightIterator) Key() (key []byte) {
	if !m.Valid() {
		panic("Invalid Iterator.")
	}
	return []byte(m.sortedKeys[m.curIdx])
}
func (m *MemoryHeightIterator) key() (key string) {
	if !m.Valid() {
		panic("Invalid Iterator.")
	}
	return m.sortedKeys[m.curIdx]
}
func (m *MemoryHeightIterator) Value() (value []byte) {
	if !m.Valid() {
		panic("Invalid Iterator.")
	}
	return []byte(m.dataset[m.key()])
}

func (m *MemoryHeightIterator) Error() error {
	panic("implement me")
}

func (m *MemoryHeightIterator) Close() { // free the MEMORY
	m.sortedKeys = nil
	m.dataset = nil
}
