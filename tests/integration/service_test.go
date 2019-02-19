package integration

import "testing"

func TestRegister(t *testing.T) {
	resp, err := requestFromFile("register", Client)
	if err != nil {
		t.Log(assumptions)
		t.Fatalf(err.Error())
	}
	t.Log(resp)
}

func TestUnRegister(t *testing.T) {
	resp, err := requestFromFile("unregister", Client)
	if err != nil {
		t.Log(assumptions)
		t.Fatalf(err.Error())
	}
	t.Log(resp)
}
