package dedup

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"
)

func Test(t *testing.T) {
	sTimes := make([]int64, 10000)
	const toAdd = "5"
	a := []byte("1,2,3,4")
	for i := 0; i < 10000; i++ {
		ti := time.Now()
		m := UnmarshalString(a)
		r := AddToString(m, toAdd)
		GetFromString(r)
		MarshalString(m)
		sTimes[i] = time.Since(ti).Microseconds()
	}
	jTimes := make([]int64, 10000)
	b, _ := json.Marshal([]string{"1,2,3,4"})
	for i := 0; i < 10000; i++ {
		ti := time.Now()
		m := make([]string, 0)
		_ = json.Unmarshal(b, &m)
		AddToSlice(m, toAdd)
		json.Marshal(m)
		jTimes[i] = time.Since(ti).Microseconds()
	}
	totalA := int64(0)
	for _, t := range sTimes {
		totalA += t
	}

	totalB := int64(0)
	for _, t := range jTimes {
		totalB += t
	}

	m := make([]string, 0)
	_ = json.Unmarshal(b, &m)
	AddToSlice(m, toAdd)
	bz, _ := json.Marshal(m)

	ma := UnmarshalString(a)
	r := AddToString(ma, toAdd)
	GetFromString(r)
	mar := MarshalString(ma)

	fmt.Println(float64(totalA) / float64(10000))
	fmt.Println(float64(totalB) / float64(10000))
	bz2 := make([]byte, 0)
	buf1 := bytes.NewBuffer(bz2)
	g := gob.NewEncoder(buf1)
	g.Encode(mar)

	bz3 := make([]byte, 0)
	buf2 := bytes.NewBuffer(bz3)
	g2 := gob.NewEncoder(buf2)
	g2.Encode(bz)
	fmt.Println(len(buf1.Bytes()))
	fmt.Println(len(buf2.Bytes()))
}

const (
	delim = ","
)

func AddToString(originalString, stringToAdd string) string {
	return originalString + delim + stringToAdd
}

func GetFromString(originalString string) []string {
	return strings.Split(originalString, delim)
}

func MarshalString(s string) []byte {
	return []byte(s)
}

func UnmarshalString(b []byte) string {
	return string(b)
}

func AddToSlice(originalSlice []string, stringToAdd string) []string {
	return append(originalSlice, stringToAdd)
}
