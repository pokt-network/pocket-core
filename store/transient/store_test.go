package transient

import (
	"testing"

	"github.com/stretchr/testify/require"
)

var k, v = []byte("hello"), []byte("world")

func TestTransientStore(t *testing.T) {
	tstore := NewStore()

	tg, _ := tstore.Get(k)
	require.Nil(t, tg)

	_ = tstore.Set(k, v)

	tg, _ = tstore.Get(k)
	require.Equal(t, v, tg)

	tstore.Commit()

	tg, _ = tstore.Get(k)
	require.Nil(t, tg)
}
