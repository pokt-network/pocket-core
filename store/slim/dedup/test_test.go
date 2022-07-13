package dedup

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"math"
	"testing"
)

func TestT(t *testing.T) {
	k := HeightKey(1, "params", []byte("hello"))
	fmt.Println(string(k))
	fmt.Println(FromHeightKey(string(k)))
}

func getRealSizeOf(v interface{}) int {
	b := new(bytes.Buffer)
	if err := gob.NewEncoder(b).Encode(v); err != nil {
		panic(err)
	}
	return b.Len()
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}
