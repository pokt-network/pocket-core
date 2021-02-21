package cli

import (
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestWaitOnStopSignals(t *testing.T) {
	go func() {
		time.Sleep(100 * time.Millisecond)
		syscall.Kill(os.Getpid(), syscall.SIGTERM)
	}()

	sig := <-waitOnStopSignals()
	require.Equal(t, sig, syscall.SIGTERM)
}
