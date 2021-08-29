package heightcache

type StoreAtHeight struct {
	height      int64
	data        map[string]string
	orderedKeys []string
}

func NewStoreAtHeight() *StoreAtHeight {
	return &StoreAtHeight{
		height: -1,
		data:   map[string]string{},
	}
}
