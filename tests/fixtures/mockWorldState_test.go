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

func TestGetApplications(t *testing.T) {
	res, err := GetApplications()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(res)
}
