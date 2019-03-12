package integration

import "testing"

const(
	register = "register"
	unregister = "unregister"
)

func TestRegister(t *testing.T) {
	resp, err := requestFromFile(register)
	if err != nil {
		t.Log(assumptions)
		t.Fatalf(err.Error())
	}
	t.Log(resp)
}

func TestUnRegister(t *testing.T) {
	resp, err := requestFromFile(unregister)
	if err != nil {
		t.Log(assumptions)
		t.Fatalf(err.Error())
	}
	t.Log(resp)
}
