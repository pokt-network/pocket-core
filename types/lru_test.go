package types

import (
	"testing"
)

func TestCache(t *testing.T) {
	testCache := NewCache(10)
	type Foo struct {
		a int
	}
	f := Foo{1}
	testCache.Add("Key", f)

	v, ok := testCache.Get("Key")
	if !ok {
		t.FailNow()
	}

	rec, _ := v.(Foo)
	if rec != f {
		t.FailNow()
	}
	testCache.Purge()
	if _, ok := testCache.Get("Key"); ok {
		t.FailNow()
	}
}
