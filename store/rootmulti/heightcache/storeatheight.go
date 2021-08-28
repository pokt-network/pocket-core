package heightcache

type StoreAtHeight struct {
	height int64
	hash   string
	data   map[string]string
}

func NewStoreAtHeight() *StoreAtHeight {
	return &StoreAtHeight{
		height: -1,
		data:   map[string]string{},
	}
}
