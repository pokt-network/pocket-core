package unit

import (
	"bufio"
	"fmt"
	"os/exec"
	"strings"
	"syscall"
	"testing"
	"time"
)

func TestStartClientTerminateSignal(t *testing.T) {
	// Starts pocket core and sends kill signal signal (15)
	StartKillPocketCore("../../cmd/pocket_core/pocket_core", 15, "terminated", 500, t)
}

func TestStartClientInterruptSignal(t *testing.T) {
	// Starts pocket core and sends kill signal signal (2)
	StartKillPocketCore("../../cmd/pocket_core/pocket_core", 2, "interrupt", 500, t)

}

func StartKillPocketCore(command string, killSignal syscall.Signal, textSignal string, millisecondsTimeout time.Duration, t *testing.T) {
	// Test when the pocket_core start and fails with no chains
	// Run the pocket-core command
	cmd := exec.Command(command)

	// We assume that we have the wrong signal until the opposite is confirmed
	correctSignal := false
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

	// Configure pipes for analyzing the process output later in a goroutine
	stdout, err_pipe := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err_pipe != nil {
		t.Fatalf("Error creating StdoutPipe for Cmd")
	}

	// Scanners for analizing the stdout/stderr process outputs
	scanner := bufio.NewScanner(stdout)
	scanner_err := bufio.NewScanner(stderr)

	defer func() {
		cmd.Wait()
		for scanner.Scan() {
			fmt.Printf("\t%s\n", scanner.Text())

			// Checks if the output of pocket-core is receiving the kill signal
			textSearch := fmt.Sprintf(`msg="Shutting down Pocket Core: Exit signal %s received`, textSignal)
			if strings.Contains(scanner.Text(), textSearch) {
				correctSignal = true
			}
		}

		for scanner_err.Scan() {
			fmt.Printf("\tError: %s\n", scanner_err.Text())
		}

		if correctSignal == true {
			t.Logf("Graceful shutdown completed with signal %s", textSignal)
		} else {
			t.Fatalf("Graceful shutdown failed")
		}
	}()

	// Run pocket core command in background
	err := cmd.Start()
	if err != nil {
		t.Fatalf("cmd.Start() failed with %s\n", err)
	}

	pgid, _ := syscall.Getpgid(cmd.Process.Pid)

	// Wait for the process to finish or kill it after a timeout (whichever happens first):
	done := make(chan error, 1)

	go func() {
		done <- cmd.Wait()
	}()

	select {
	// Send the kill signal after millisecondsTimeout
	case <-time.After(millisecondsTimeout * time.Millisecond):
		if err := syscall.Kill(-pgid, killSignal); err != nil {
			t.Fatalf("Failed to kill process")
		}
	case err := <-done:
		if err != nil {
			t.Fatalf("process finished with error")
		}
	}

}
