package testutils

import (
	"bufio"
	"context"
	"net"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func SetupHyprEnvVars(t *testing.T) (string, string) {
	tempDir := t.TempDir()
	signature := "test_signature"
	hyprDir := filepath.Join(tempDir, "hypr", signature)
	//nolint:gosec
	err := os.MkdirAll(hyprDir, 0o755)
	require.NoError(t, err, "failed to create hypr directory")

	originalXDG := os.Getenv("XDG_RUNTIME_DIR")
	originalSig := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	t.Cleanup(func() {
		_ = os.Setenv("XDG_RUNTIME_DIR", originalXDG)
		_ = os.Setenv("HYPRLAND_INSTANCE_SIGNATURE", originalSig)
	})

	err = os.Setenv("XDG_RUNTIME_DIR", tempDir)
	require.NoError(t, err, "failed to set env var")
	err = os.Setenv("HYPRLAND_INSTANCE_SIGNATURE", signature)
	require.NoError(t, err, "failed to set env var")
	return tempDir, signature
}

func SetupHyprSocket(ctx context.Context, t *testing.T, xdgRuntimeDir, signature string,
	hyprSocketFun func(string, string) string,
) (net.Listener, func()) {
	socketPath := hyprSocketFun(xdgRuntimeDir, signature)
	lc := &net.ListenConfig{}
	listener, err := lc.Listen(ctx, "unix", socketPath)
	require.NoError(t, err, "failed to create a test socket %s", socketPath)
	return listener, func() {
		err = listener.Close()
		require.NoError(t, err, "cant close hypr socket")
	}
}

// SetupFakeHyprIPCWriter serves expectedCommands in order, one per connection.
// Lua evals ("/eval ...", sent when the daemon live-applies a config) can arrive at any
// point; they are answered with "ok" without consuming an expected command, and keep being
// answered after the returned channel closes, until the listener is closed.
func SetupFakeHyprIPCWriter(t *testing.T, listener net.Listener, responseData [][]byte,
	expectedCommands []string, exitOnError bool,
) chan struct{} {
	serverDone := make(chan struct{})
	go func() {
		i := 0
		if len(expectedCommands) == 0 {
			close(serverDone)
		}
		for {
			conn, err := listener.Accept()
			if err != nil {
				if i < len(expectedCommands) {
					if !exitOnError {
						assert.Nil(t, err, "failed to accept connection")
					}
					close(serverDone)
				}
				return
			}
			if !respondToHyprIPC(t, conn, expectedCommands, responseData, &i) {
				continue
			}
			if i == len(expectedCommands) {
				close(serverDone)
			}
		}
	}()
	return serverDone
}

// respondToHyprIPC answers a single connection and reports whether it consumed an expected command.
func respondToHyprIPC(t *testing.T, conn net.Conn, expectedCommands []string, responseData [][]byte, i *int) bool {
	defer func() { _ = conn.Close() }()

	t.Log("listener accepted a connection")

	s := bufio.NewScanner(conn)
	if !s.Scan() {
		t.Errorf("no data to read")
		return false
	}
	command := s.Text()

	if strings.HasPrefix(command, "/eval ") {
		t.Logf("answering eval: %s", command)
		_, err := conn.Write([]byte("ok"))
		assert.Nil(t, err, "failed to write eval response")
		return false
	}

	if *i >= len(expectedCommands) {
		t.Errorf("unexpected command after all expected ones were served: %s", command)
		return false
	}

	assert.Equal(t, expectedCommands[*i], command, "wrong command")
	_, err := conn.Write(responseData[*i])
	assert.Nil(t, err, "failed to write response")
	t.Logf("wrote response to the client %s", string(responseData[*i]))
	*i++
	return true
}

func SetupFakeHyprEventsServer(ctx context.Context, t *testing.T, listener net.Listener, events []string) chan struct{} {
	return SetupFakeHyprEventsServerWithSleep(ctx, t, listener, events, nil, 10*time.Millisecond)
}

func SetupFakeHyprEventsServerWithSleep(ctx context.Context, t *testing.T, listener net.Listener,
	events []string, initialSleep *time.Duration, sleepBetweenEvents time.Duration,
) chan struct{} {
	serverDone := make(chan struct{})
	go func() {
		defer func() {
			Logf(t, "Fake Hypr events server goroutine exiting")
			close(serverDone)
		}()

		Logf(t, "Fake Hypr events server: waiting for connection")
		conn, err := listener.Accept()
		if err != nil {
			t.Errorf("Failed to accept connection: %v", err)
			return
		}
		defer func() {
			Logf(t, "Fake Hypr events server: closing connection")
			if conn != nil {
				_ = conn.Close()
			}
		}()

		Logf(t, "Fake Hypr events server: accepted connection, ready to send events")

		for i, event := range events {
			if initialSleep != nil {
				time.Sleep(*initialSleep)
			}

			select {
			case <-ctx.Done():
				Logf(t, "Fake Hypr events server: context cancelled while sending events")
				return
			default:
			}

			if _, err := conn.Write([]byte(event + "\n")); err != nil {
				t.Errorf("Failed to write event %d: %v", i, err)
				return
			}
			Logf(t, "Wrote event %d on the event socket: %s", i, event)
			time.Sleep(sleepBetweenEvents)
		}

		Logf(t, "Fake Hypr events server: finished sending %d events, waiting for context cancellation", len(events))
		<-ctx.Done()
		Logf(t, "Fake Hypr events server: context cancelled (%v), shutting down", context.Cause(ctx))
	}()
	return serverDone
}
