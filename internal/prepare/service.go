// Package prepare defines utility service to prepare prior to hyprdynamicmonitors daemon run
package prepare

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/fiffeek/hyprdynamicmonitors/internal/config"
	"github.com/fiffeek/hyprdynamicmonitors/internal/utils"
	"github.com/sirupsen/logrus"
)

type Service struct {
	cfg                 *config.Config
	monitorDisableRegex *regexp.Regexp
}

func NewService(cfg *config.Config) *Service {
	monitorDisableRegex := regexp.MustCompile(`\bdisabled\s*=\s*true\b`)
	return &Service{
		cfg,
		monitorDisableRegex,
	}
}

// TruncateDestination re-enables every monitor the destination disables
// (disabled = true -> disabled = false), so Hyprland never starts without an active display.
// Lua comment lines are left untouched.
func (s *Service) TruncateDestination() error {
	file := s.cfg.Get().General.Destination
	_, err := os.Stat(*file)
	if err != nil {
		logrus.WithFields(logrus.Fields{"destination": *file}).Info("file does not exist")
		return nil
	}

	contents, err := os.ReadFile(*file)
	if err != nil {
		return fmt.Errorf("cant read the %s destination file: %w", *file, err)
	}

	lines := strings.Split(string(contents), "\n")
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "--") || !s.monitorDisableRegex.MatchString(line) {
			continue
		}
		lines[i] = s.monitorDisableRegex.ReplaceAllString(line, "disabled = false")
		logrus.Infof("Line %s will be re-enabled", line)
	}

	newContents := strings.Join(lines, "\n")
	if err := utils.WriteAtomic(*file, []byte(newContents)); err != nil {
		return fmt.Errorf("cant write to %s destination file: %w", *file, err)
	}

	return nil
}
