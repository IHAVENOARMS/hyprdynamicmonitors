package hypr_test

import (
	"testing"

	"github.com/fiffeek/hyprdynamicmonitors/internal/hypr"
	"github.com/stretchr/testify/assert"
)

func TestLuaString(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"DP-1", `"DP-1"`},
		{`say "hi"`, `"say \"hi\""`},
		{`back\slash`, `"back\\slash"`},
		{"Monitor #1", `"Monitor #1"`},
		{"tab\there\n", `"tab\there\n"`},
		{"bell\x07", `"bell\007"`},
	}
	for _, tt := range tests {
		assert.Equal(t, tt.want, hypr.LuaString(tt.in))
	}
}

func TestMonitorSpec_ToLua(t *testing.T) {
	id := 1
	tests := []struct {
		name       string
		monitor    *hypr.MonitorSpec
		wantBlock  string
		wantInline string
	}{
		{
			name: "enabled with description",
			monitor: &hypr.MonitorSpec{
				ID: &id, Name: "DP-1", Description: "Dell Inc. DELL U2719DC D6LPTS2",
				Width: 2560, Height: 1440, RefreshRate: 59.951, X: -2560, Y: 0,
				Scale: 1, Transform: 3, SdrBrightness: 1, SdrSaturation: 1,
			},
			wantBlock: `hl.monitor({
	output = "desc:Dell Inc. DELL U2719DC D6LPTS2",
	mode = "2560x1440@59.951",
	position = "-2560x0",
	scale = 1,
	transform = 3,
	vrr = false,
})`,
			wantInline: `hl.monitor({ output = "desc:Dell Inc. DELL U2719DC D6LPTS2", mode = "2560x1440@59.951", ` +
				`position = "-2560x0", scale = 1, transform = 3, vrr = false })`,
		},
		{
			name: "falls back to name, hdr, mirror, bitdepth",
			monitor: &hypr.MonitorSpec{
				ID: &id, Name: "eDP-1", Width: 2880, Height: 1800, RefreshRate: 120,
				Scale: 1.5, Vrr: true, TenBitdepth: true, ColorPreset: "hdr",
				SdrBrightness: 1.2, SdrSaturation: 0.95, Mirror: "DP-2",
			},
			wantBlock: `hl.monitor({
	output = "eDP-1",
	mode = "2880x1800@120",
	position = "0x0",
	scale = 1.5,
	transform = 0,
	vrr = true,
	bitdepth = 10,
	cm = "hdr",
	sdrbrightness = 1.2,
	sdrsaturation = 0.95,
	mirror = "DP-2",
})`,
		},
		{
			name: "sdr values ignored without hdr, default presets omitted",
			monitor: &hypr.MonitorSpec{
				ID: &id, Name: "eDP-1", Width: 1920, Height: 1080, RefreshRate: 60,
				Scale: 1, ColorPreset: "srgb", SdrBrightness: 1.2, Mirror: "none",
			},
			wantBlock: `hl.monitor({
	output = "eDP-1",
	mode = "1920x1080@60",
	position = "0x0",
	scale = 1,
	transform = 0,
	vrr = false,
})`,
		},
		{
			name: "disabled is a single line",
			monitor: &hypr.MonitorSpec{
				ID: &id, Name: "eDP-1", Description: `Weird "Quoted" #1`, Disabled: true,
				Width: 1920, Height: 1080,
			},
			wantBlock:  `hl.monitor({ output = "desc:Weird \"Quoted\" #1", disabled = true })`,
			wantInline: `hl.monitor({ output = "desc:Weird \"Quoted\" #1", disabled = true })`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantBlock, tt.monitor.ToLua())
			if tt.wantInline != "" {
				assert.Equal(t, tt.wantInline, tt.monitor.ToLuaInline())
			}
		})
	}
}
