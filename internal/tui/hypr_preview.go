package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type HyprPreviewPane struct {
	monitors []*MonitorSpec
	clamp    bool
	width    int
	colors   *ColorsManager
}

func NewHyprPreviewPane(monitors []*MonitorSpec, colors *ColorsManager) *HyprPreviewPane {
	return &HyprPreviewPane{
		monitors: monitors,
		width:    0,
		clamp:    true,
		colors:   colors,
	}
}

func (h *HyprPreviewPane) SetWidth(width int) {
	h.width = width
}

func (h *HyprPreviewPane) SetClamp(clamp bool) {
	h.clamp = clamp
}

func (h *HyprPreviewPane) View() string {
	var sections []string

	if len(h.monitors) == 0 {
		sections = append(sections, h.colors.TitleStyle().Render(
			"Hyprland Config Preview\n\nNo monitors to configure"))
	} else {
		sections = append(sections, h.colors.TitleStyle().Render("Hyprland Config Preview"))
	}

	for _, monitor := range h.monitors {
		hyprCommand := monitor.ToLua()
		if h.clamp {
			hyprCommand = ClampTextTo(compactLua(monitor, h.width), h.width)
		}
		sections = append(sections, h.colors.MutedStyle().Render(hyprCommand))
	}
	return lipgloss.JoinVertical(lipgloss.Left, sections...)
}

// compactLua fits the hl.monitor call into maxLen: it first drops the
// hl.monitor({ output = ... }) boilerplate and then trims the middle of the description,
// keeping both the monitor identity and its settings visible.
func compactLua(monitor *MonitorSpec, maxLen int) string {
	const minKept = 8
	line := monitor.ToLua()
	if maxLen < 0 || len(line) <= maxLen {
		return line
	}
	line = luaFieldsOnly(monitor)
	overflow := len(line) - maxLen
	if overflow <= 0 || monitor.Description == "" {
		return line
	}
	desc := []rune(monitor.Description)
	keep := max(len(desc)-overflow-len(shortenedMarker), minKept)
	if keep >= len(desc) {
		return line
	}
	head := (keep + 1) / 2
	shortened := *monitor
	shortened.Description = string(desc[:head]) + shortenedMarker + string(desc[len(desc)-(keep-head):])
	return luaFieldsOnly(&shortened)
}

func luaFieldsOnly(monitor *MonitorSpec) string {
	return strings.TrimSuffix(strings.TrimPrefix(monitor.ToLua(), "hl.monitor({ output = "), " })")
}

const shortenedMarker = "..."

func ClampTextTo(text string, maxLen int) string {
	if maxLen < 0 {
		return text
	}
	if len(text) <= maxLen {
		return text
	}
	dots := "..."
	lenWithoutDots := maxLen - len(dots)
	return dots + text[len(text)-lenWithoutDots:]
}
