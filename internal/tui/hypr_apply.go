package tui

import (
	"context"
	"errors"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/fiffeek/hyprdynamicmonitors/internal/config"
	"github.com/fiffeek/hyprdynamicmonitors/internal/generators"
	"github.com/fiffeek/hyprdynamicmonitors/internal/hypr"
	"github.com/fiffeek/hyprdynamicmonitors/internal/matchers"
	"github.com/fiffeek/hyprdynamicmonitors/internal/power"
	"github.com/fiffeek/hyprdynamicmonitors/internal/profilemaker"
	"github.com/sirupsen/logrus"
)

type HyprApply struct {
	profileMaker *profilemaker.Service
	generator    *generators.ConfigGenerator
}

func NewHyprApply(profileMaker *profilemaker.Service, generator *generators.ConfigGenerator) *HyprApply {
	return &HyprApply{
		profileMaker: profileMaker,
		generator:    generator,
	}
}

func (h *HyprApply) ApplyCurrent(monitors []*MonitorSpec) tea.Cmd {
	blocks := make([]string, 0, len(monitors))
	for _, monitor := range monitors {
		blocks = append(blocks, monitor.ToLua())
	}

	err := hypr.Eval(context.Background(), strings.Join(blocks, " "))
	if err != nil {
		logrus.WithError(err).Error("cant apply hypr settings")
	}

	return OperationStatusCmd(OperationNameEphemeralApply, err)
}

func (h *HyprApply) CreateProfile(monitors []*MonitorSpec, name, file string) tea.Cmd {
	hyprMonitors, err := ConvertToHyprMonitors(monitors)
	if err != nil {
		return OperationStatusCmd(OperationNameCreateProfile, err)
	}
	err = h.profileMaker.FreezeGivenAs(name, file, hyprMonitors)
	return OperationStatusCmd(OperationNameCreateProfile, err)
}

func (h *HyprApply) EditProfile(monitors []*MonitorSpec, name string) tea.Cmd {
	hyprMonitors, err := ConvertToHyprMonitors(monitors)
	if err != nil {
		return OperationStatusCmd(OperationNameEditProfile, err)
	}
	err = h.profileMaker.EditExisting(name, hyprMonitors)
	return OperationStatusCmd(OperationNameEditProfile, err)
}

func (h *HyprApply) GenerateThroughHDM(cfg *config.Config, profile *matchers.MatchedProfile,
	monitors []*MonitorSpec, powerState power.PowerState, lidState power.LidState,
) tea.Cmd {
	if profile == nil {
		return OperationStatusCmd(OperationNameHydrate, errors.New("profile is nil"))
	}
	hyprMonitors, err := ConvertToHyprMonitors(monitors)
	if err != nil {
		return OperationStatusCmd(OperationNameHydrate, err)
	}
	destination := *cfg.Get().General.Destination
	changed, err := h.generator.GenerateConfig(cfg.Get(), profile, hyprMonitors, powerState, lidState, destination, false)
	if err == nil && changed && *cfg.Get().General.LiveApply {
		err = hypr.ApplyLuaFile(context.Background(), destination)
	}
	return OperationStatusCmd(OperationNameHydrate, err)
}
