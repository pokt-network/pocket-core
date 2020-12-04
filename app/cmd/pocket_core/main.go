package main

import (
	"github.com/pokt-network/pocket-core/app/cmd/cli"
	"os"
	"os/exec"
	"strings"
)

func main() {
	godebug := os.Getenv("GODEBUG")
	//Check for the madvdontneed variable
	if strings.Contains(godebug, "madvdontneed=1") {
		cli.Execute()
	} else {
		env := append(os.Environ(), "GODEBUG="+"madvdontneed=1,"+godebug)
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Env = env
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
}
