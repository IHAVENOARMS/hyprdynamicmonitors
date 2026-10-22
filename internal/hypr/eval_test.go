package hypr_test

import (
	"bufio"
	"context"
	"testing"
	"time"

	"github.com/fiffeek/hyprdynamicmonitors/internal/hypr"
	"github.com/fiffeek/hyprdynamicmonitors/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyLuaFile(t *testing.T) {
	tests := []struct {
		name          string
		path          string
		response      string
		expectedCmd   string
		expectedError string
	}{
		{
			name:        "sends dofile for the destination",
			path:        "/home/user/.config/hypr/hyprdynamicmonitors.lua",
			response:    "ok",
			expectedCmd: `/eval dofile("/home/user/.config/hypr/hyprdynamicmonitors.lua")`,
		},
		{
			name:        "quotes paths",
			path:        `/tmp/we"ird\path.lua`,
			response:    "ok",
			expectedCmd: `/eval dofile("/tmp/we\"ird\\path.lua")`,
		},
		{
			name:          "surfaces lua errors",
			path:          "/tmp/broken.lua",
			response:      `error: /tmp/broken.lua:1: hl.monitor: unknown field 'bogus'`,
			expectedCmd:   `/eval dofile("/tmp/broken.lua")`,
			expectedError: "unknown field 'bogus'",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()

			xdgRuntimeDir, signature := testutils.SetupHyprEnvVars(t)
			listener, teardown := testutils.SetupHyprSocket(ctx, t, xdgRuntimeDir, signature, hypr.GetHyprSocket)
			defer teardown()

			received := make(chan string, 1)
			go func() {
				conn, err := listener.Accept()
				if err != nil {
					return
				}
				defer func() { _ = conn.Close() }()
				s := bufio.NewScanner(conn)
				if s.Scan() {
					received <- s.Text()
				}
				_, _ = conn.Write([]byte(tt.response))
			}()

			err := hypr.ApplyLuaFile(ctx, tt.path)

			select {
			case cmd := <-received:
				assert.Equal(t, tt.expectedCmd, cmd)
			case <-time.After(time.Second):
				t.Fatal("server did not receive a command")
			}

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}
			assert.NoError(t, err)
		})
	}
}

func TestEval_RejectsMultiline(t *testing.T) {
	err := hypr.Eval(context.Background(), "local a = 1\nlocal b = 2")
	assert.Error(t, err)
}
