package fixtures

import (
	"fmt"
	"testing"
)

func TestGetNodes(t *testing.T) {
	res, err := GetNodes()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(res)
}

func TestGetDevelopers(t *testing.T) {
	res, err := GetDevelopers()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(res)
}
