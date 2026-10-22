package hypr

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fiffeek/hyprdynamicmonitors/internal/dial"
	"github.com/sirupsen/logrus"
)

const evalTimeout = 5 * time.Second

// Eval runs a Lua chunk inside the running Hyprland instance (the socket
// equivalent of hyprctl eval). The chunk must be a single line.
func Eval(ctx context.Context, code string) error {
	if strings.Contains(code, "\n") {
		return errors.New("lua chunk for eval must be a single line")
	}
	signature := os.Getenv("HYPRLAND_INSTANCE_SIGNATURE")
	if signature == "" {
		return errors.New("HYPRLAND_INSTANCE_SIGNATURE environment variable not set - are you running under Hyprland?")
	}

	ctx, cancel := context.WithTimeout(ctx, evalTimeout)
	defer cancel()

	socketPath := GetHyprSocket(os.Getenv("XDG_RUNTIME_DIR"), signature)
	conn, teardown, err := dial.GetUnixSocketConnection(ctx, socketPath)
	if err != nil {
		return fmt.Errorf("cant open socket to %s: %w", socketPath, err)
	}
	defer teardown()

	if deadline, ok := ctx.Deadline(); ok {
		_ = conn.SetDeadline(deadline)
	}
	if _, err := conn.Write([]byte("/eval " + code + "\n")); err != nil {
		return fmt.Errorf("cant send eval: %w", err)
	}
	response, err := io.ReadAll(conn)
	if err != nil {
		return fmt.Errorf("cant read eval response: %w", err)
	}

	output := strings.TrimSpace(string(response))
	logrus.WithFields(logrus.Fields{"code": code, "response": output}).Debug("eval response")
	if strings.HasPrefix(output, "error") {
		return fmt.Errorf("hyprland eval failed: %s", output)
	}
	return nil
}

// ApplyLuaFile executes a rendered Lua config file inside Hyprland so its
// hl.monitor calls take effect without a full config reload.
func ApplyLuaFile(ctx context.Context, path string) error {
	return Eval(ctx, "dofile("+LuaString(path)+")")
}
