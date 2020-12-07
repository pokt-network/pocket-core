package main

import (
	"github.com/pokt-network/pocket-core/app/cmd/cli"
	"os"
	"os/exec"
	"strings"
)

func main() {
	//Get the GODEBUG env variable
	godebug := os.Getenv("GODEBUG")
	//Check if the --madvdontneed=true
	flagPresent := containsFlag(os.Args[1:], "--madvdontneed=true")

	//Check if madvdontneed env variable is present or flag is not used
	if strings.Contains(godebug, "madvdontneed=1") || !flagPresent {
		//start normally
		cli.Execute()
	} else {
		//flag --madvdontneed=true so we add the env variable and start pocket as a subprocess
		env := append(os.Environ(), "GODEBUG="+"madvdontneed=1,"+godebug)
		cmd := exec.Command(os.Args[0], os.Args[1:]...)
		cmd.Env = env
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
}

func containsFlag(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
