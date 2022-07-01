package memdb

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	db "github.com/tendermint/tm-db"
	"math/rand"
	"testing"
	"time"
)

func Test(t *testing.T) {
	memDB := NewPocketMemDB()
	memDB2 := db.NewGoLevelMemDB()
	kv := map[string][]byte{"1": []byte("1Value"), "2": []byte("2Value"), "3": []byte("3Value"), "5": []byte("5Value")}
	for k, v := range kv {
		if err := memDB.Set([]byte(k), v); err != nil {
			panic(k)
		}
		if err := memDB2.Set([]byte(k), v); err != nil {
			panic(err)
		}
	}
	itA, err := memDB.Iterator([]byte("1"), []byte("6"))
	assert.Nil(t, err)
	itB, err := memDB2.Iterator([]byte("1"), []byte("6"))
	assert.Nil(t, err)
	for ; itA.Valid(); itA.Next() {
		if string(itA.Key()) == "3" {
			memDB.Set([]byte("5"), []byte("newValue"))
			memDB.Delete([]byte("3"))
		}
		fmt.Println(string(itA.Key()), string(itA.Value()))
	}
	fmt.Println("")
	for ; itB.Valid(); itB.Next() {
		if string(itB.Key()) == "3" {
			memDB2.Set([]byte("5"), []byte("newValue"))
			memDB2.Delete([]byte("3"))
		}
		fmt.Println(string(itB.Key()), string(itB.Value()))
	}
}

func dbSetup(_ *testing.T, memDB, memDB2 db.DB) {
	kv := make(map[string][]byte)
	fmt.Println("initializing the data to go in each db")
	for i := 0; i < (rand.Intn(10000) + 100); i++ {
		rand.Seed(int64(time.Now().Nanosecond()))
		key := randString(rand.Intn(100) + 1)
		rand.Seed(int64(time.Now().Nanosecond()))
		value := randString(rand.Intn(100) + 1)
		kv[key] = []byte(value)
	}
	fmt.Println("setting data in the db")
	for k, v := range kv {
		if err := memDB.Set([]byte(k), v); err != nil {
			panic(k)
		}
		if err := memDB2.Set([]byte(k), v); err != nil {
			panic(err)
		}
	}
}

func TestMemDBParityIterator(t *testing.T) {
	memDB := NewPocketMemDB()
	memDB2 := db.NewGoLevelMemDB()
	fmt.Println("running the iterator scenarios")
	dbSetup(t, memDB, memDB2)
	for i := 0; i < 2000; i++ {
		var iteratorCondition string
		fmt.Println("iteration", i)
		var err error
		itA, itB := db.Iterator(nil), db.Iterator(nil)
		rand.Seed(int64(time.Now().Nanosecond()))
		start := []byte(randString(rand.Intn(100) + 1))
		rand.Seed(int64(time.Now().Nanosecond()))
		end := []byte(randString(rand.Intn(100) + 1))
		if len(start) <= 20 {
			start = nil
		}
		if len(end) <= 20 {
			end = nil
		}
		if rand.Intn(2) == 0 {
			iteratorCondition += fmt.Sprintf("Normal iterator with %v %v", string(start), string(end))
			itA, err = memDB.Iterator([]byte(start), []byte(end))
			assert.Nil(t, err)
			itB, err = memDB2.Iterator([]byte(start), []byte(end))
			assert.Nil(t, err)
		} else {
			iteratorCondition += fmt.Sprintf("Reverse iterator with %v %v", string(start), string(end))
			itA, err = memDB.ReverseIterator([]byte(start), []byte(end))
			assert.Nil(t, err)
			itB, err = memDB2.ReverseIterator([]byte(start), []byte(end))
			assert.Nil(t, err)
		}
		keysA := make([][]byte, 0)
		valuesA := make([][]byte, 0)
		keysB := make([][]byte, 0)
		valuesB := make([][]byte, 0)
		for ; itA.Valid(); itA.Next() {
			keysA = append(keysA, itA.Key())
			valuesA = append(valuesA, itA.Value())
		}
		for ; itB.Valid(); itB.Next() {
			keysB = append(keysB, itB.Key())
			valuesB = append(valuesB, itB.Value())
		}
		itA.Close()
		itB.Close()
		assert.Equal(t, len(keysA), len(keysB), "keysA and keysB are not equal len: ", len(keysA), len(keysB), iteratorCondition)
		assert.Equal(t, len(valuesA), len(valuesB), "valuesA and valuesB are not equal len: ", len(keysA), len(keysB), iteratorCondition)
		assert.Equal(t, len(keysA), len(keysB))
		assert.Equal(t, keysA, keysB, "keysA and keysB are not equal slices", keysA, keysB, iteratorCondition)
		assert.Equal(t, valuesA, valuesB, "valuesA and valuesB are not equal slices", valuesA, valuesB, iteratorCondition)
		if len(keysA) != 0 {
			op(t, memDB, memDB2, keysB, valuesB)
		}
	}
}

func op(t *testing.T, dbA, dbB db.DB, keys [][]byte, values [][]byte) {
	rand.Seed(int64(time.Now().Nanosecond()))
	switch rand.Intn(6) {
	case 0:
		fmt.Println("has")
		key := keys[rand.Intn(len(keys))]
		found, err := dbA.Has(key)
		assert.Nil(t, err)
		assert.True(t, found)
		found, err = dbB.Has(key)
		assert.Nil(t, err)
		assert.True(t, found)
	case 1:
		fmt.Println("get")
		index := rand.Intn(len(keys))
		key := keys[index]
		val, err := dbA.Get(key)
		assert.Nil(t, err)
		assert.Equal(t, val, values[index])
		val, err = dbB.Get(key)
		assert.Nil(t, err)
		assert.Equal(t, val, values[index])
	case 2:
		fmt.Println("Set")
		index := rand.Intn(len(keys))
		newValue := []byte(randString(rand.Intn(1000)))
		key := keys[index]
		err := dbA.Set(key, newValue)
		assert.Nil(t, err)
		val, err := dbA.Get(key)
		assert.Nil(t, err)
		assert.Equal(t, val, newValue)
		values[index] = newValue

		err = dbB.Set(key, newValue)
		assert.Nil(t, err)
		val, err = dbB.Get(key)
		assert.Nil(t, err)
		assert.Equal(t, val, newValue)
		values[index] = newValue
	case 3:
		fmt.Println("delete")
		index := rand.Intn(len(keys))
		key := keys[index]
		err := dbA.Delete(key)
		assert.Nil(t, err)
		found, err := dbA.Has(key)
		assert.Nil(t, err)
		assert.True(t, !found)
		val, err := dbA.Get(key)
		assert.Nil(t, err)
		assert.Nil(t, val)

		err = dbB.Delete(key)
		assert.Nil(t, err)
		found, err = dbB.Has(key)
		assert.Nil(t, err)
		assert.True(t, !found)
		val, err = dbB.Get(key)
		assert.Nil(t, err)
		assert.Nil(t, val)
	case 4:
		fmt.Println("batch")
		b := dbA.NewBatch()
		b2 := dbB.NewBatch()
		notValue := []byte(randString(rand.Intn(1000)))
		newValue := []byte(randString(rand.Intn(1000)))
		index := rand.Intn(len(keys))
		key := keys[index]
		b.Set(key, notValue)
		b2.Set(key, notValue)
		b.Delete(key)
		b2.Delete(key)
		b.Set(key, newValue)
		b2.Set(key, newValue)
		err := b.Write()
		b.Close()
		assert.Nil(t, err)
		err = b2.Write()
		b2.Close()
		assert.Nil(t, err)
		val, err := dbA.Get(key)
		assert.Nil(t, err)
		assert.Equal(t, val, newValue)
		values[index] = newValue
		val, err = dbB.Get(key)
		assert.Nil(t, err)
		assert.Equal(t, val, newValue)
		values[index] = newValue
	}
}

func randString(n int) string {
	var runes = []rune("1234567890abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}
