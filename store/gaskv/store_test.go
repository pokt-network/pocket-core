package gaskv_test

import (
	"fmt"
	"testing"

	dbm "github.com/tendermint/tm-db"

	"github.com/pokt-network/pocket-core/store/dbadapter"
	"github.com/pokt-network/pocket-core/store/gaskv"
	"github.com/pokt-network/pocket-core/store/types"

	"github.com/stretchr/testify/require"
)

func bz(s string) []byte { return []byte(s) }

func keyFmt(i int) []byte { return bz(fmt.Sprintf("key%0.8d", i)) }
func valFmt(i int) []byte { return bz(fmt.Sprintf("value%0.8d", i)) }

func TestGasKVStoreBasic(t *testing.T) {
	mem := dbadapter.Store{DB: dbm.NewMemDB()}
	meter := types.NewGasMeter(10000)
	st := gaskv.NewStore(mem, meter, types.KVGasConfig())
	sg, _ := st.Get(keyFmt(1))
	require.Empty(t, sg, "Expected `key1` to be empty")
	_ = st.Set(keyFmt(1), valFmt(1))
	sg, _ = st.Get(keyFmt(1))
	require.Equal(t, valFmt(1), sg)
	_ = st.Delete(keyFmt(1))
	sg, _ = st.Get(keyFmt(1))
	require.Empty(t, sg, "Expected `key1` to be empty")
	require.Equal(t, meter.GasConsumed(), types.Gas(6429))
}

func TestGasKVStoreIterator(t *testing.T) {
	mem := dbadapter.Store{DB: dbm.NewMemDB()}
	meter := types.NewGasMeter(10000)
	st := gaskv.NewStore(mem, meter, types.KVGasConfig())
	sg, _ := st.Get(keyFmt(1))
	sg2, _ := st.Get(keyFmt(2))
	require.Empty(t, sg, "Expected `key1` to be empty")
	require.Empty(t, sg2, "Expected `key2` to be empty")
	_ = st.Set(keyFmt(1), valFmt(1))
	_ = st.Set(keyFmt(2), valFmt(2))
	iterator, _ := st.Iterator(nil, nil)
	ka := iterator.Key()
	require.Equal(t, ka, keyFmt(1))
	va := iterator.Value()
	require.Equal(t, va, valFmt(1))
	iterator.Next()
	kb := iterator.Key()
	require.Equal(t, kb, keyFmt(2))
	vb := iterator.Value()
	require.Equal(t, vb, valFmt(2))
	iterator.Next()
	require.False(t, iterator.Valid())
	require.Panics(t, iterator.Next)
	require.Equal(t, meter.GasConsumed(), types.Gas(6987))
}

func TestGasKVStoreOutOfGasSet(t *testing.T) {
	mem := dbadapter.Store{DB: dbm.NewMemDB()}
	meter := types.NewGasMeter(0)
	st := gaskv.NewStore(mem, meter, types.KVGasConfig())
	require.Panics(t, func() { _ = st.Set(keyFmt(1), valFmt(1)) }, "Expected out-of-gas")
}

func TestGasKVStoreOutOfGasIterator(t *testing.T) {
	mem := dbadapter.Store{DB: dbm.NewMemDB()}
	meter := types.NewGasMeter(20000)
	st := gaskv.NewStore(mem, meter, types.KVGasConfig())
	_ = st.Set(keyFmt(1), valFmt(1))
	iterator, _ := st.Iterator(nil, nil)
	iterator.Next()
	require.Panics(t, func() { iterator.Value() }, "Expected out-of-gas")
}
