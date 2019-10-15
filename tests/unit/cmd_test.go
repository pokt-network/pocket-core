package unit

import (
	"os/exec"
	"testing"

	utils "github.com/pokt-network/pocket-core/tests/unit/utils"
)

func TestBuildClient(t *testing.T) {
	// Build the pocket_core binary in workdir
	build := exec.Command("go", "build", "../../cmd/pocket_core")
	err := build.Run()
	if err != nil {
		t.Fatalf("Error building pocket_core binary")
	}
}

func TestStartClientTerminateSignal(t *testing.T) {
	// Starts pocket core and sends kill signal signal (15)
	args := []string{"./pocket_core"}
	utils.StartKillPocketCore(args, 15, "terminated", 500, t)
}

func TestStartClientInterruptSignal(t *testing.T) {
	// Starts pocket core and sends kill signal signal (2)
	args := []string{"./pocket_core"}
	utils.StartKillPocketCore(args, 2, "interrupt", 500, t)

}
