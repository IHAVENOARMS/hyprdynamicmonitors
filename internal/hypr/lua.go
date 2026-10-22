package hypr

import (
	"fmt"
	"strconv"
	"strings"
)

// LuaString quotes s as a Lua string literal.
func LuaString(s string) string {
	var b strings.Builder
	b.WriteByte('"')
	for _, r := range s {
		switch r {
		case '\\':
			b.WriteString(`\\`)
		case '"':
			b.WriteString(`\"`)
		case '\n':
			b.WriteString(`\n`)
		case '\r':
			b.WriteString(`\r`)
		case '\t':
			b.WriteString(`\t`)
		default:
			if r < 0x20 || r == 0x7f {
				fmt.Fprintf(&b, `\%03d`, r)
				continue
			}
			b.WriteRune(r)
		}
	}
	b.WriteByte('"')
	return b.String()
}

func luaNumber(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

// LuaOutput returns the hl.monitor output selector, preferring the description
// so that profiles survive port swaps.
func (m *MonitorSpec) LuaOutput() string {
	if m.Description != "" {
		return "desc:" + m.Description
	}
	return m.Name
}

type luaField struct {
	key   string
	value string
}

// luaFields lists the hl.monitor fields for the monitor. hl.monitor merges into an
// existing rule for the same output instead of replacing it, so with explicitDefaults
// every field is written (including defaults) to fully reset whatever was set before,
// e.g. disabled = true from an earlier apply. Without it, defaults are omitted, which is
// only meant for display.
func (m *MonitorSpec) luaFields(explicitDefaults bool) []luaField {
	fields := []luaField{{"output", LuaString(m.LuaOutput())}}
	if m.Disabled {
		return append(fields, luaField{"disabled", "true"})
	}

	fields = append(fields,
		luaField{"mode", LuaString(fmt.Sprintf("%dx%d@%s", m.Width, m.Height, luaNumber(m.RefreshRate)))},
		luaField{"position", LuaString(fmt.Sprintf("%dx%d", m.X, m.Y))},
		luaField{"scale", luaNumber(m.Scale)},
		luaField{"transform", strconv.Itoa(m.Transform)},
		luaField{"vrr", strconv.FormatBool(m.Vrr)},
	)
	if explicitDefaults {
		fields = append(fields, luaField{"disabled", "false"})
	}
	if m.TenBitdepth {
		fields = append(fields, luaField{"bitdepth", "10"})
	} else if explicitDefaults {
		fields = append(fields, luaField{"bitdepth", "8"})
	}
	if m.HasNonDefaultColorPreset() {
		fields = append(fields, luaField{"cm", LuaString(m.ColorPreset)})
	} else if explicitDefaults {
		preset := m.ColorPreset
		if preset == "" {
			preset = "srgb"
		}
		fields = append(fields, luaField{"cm", LuaString(preset)})
	}
	sdrBrightness, sdrSaturation := 1.0, 1.0
	if m.HDR() {
		sdrBrightness, sdrSaturation = m.SdrBrightness, m.SdrSaturation
	}
	if explicitDefaults || sdrBrightness != 1.0 {
		fields = append(fields, luaField{"sdrbrightness", luaNumber(sdrBrightness)})
	}
	if explicitDefaults || sdrSaturation != 1.0 {
		fields = append(fields, luaField{"sdrsaturation", luaNumber(sdrSaturation)})
	}
	if m.HasMirror() {
		fields = append(fields, luaField{"mirror", LuaString(m.Mirror)})
	} else if explicitDefaults {
		fields = append(fields, luaField{"mirror", `""`})
	}
	return fields
}

func joinInline(fields []luaField) string {
	parts := []string{}
	for _, f := range fields {
		parts = append(parts, f.key+" = "+f.value)
	}
	return "hl.monitor({ " + strings.Join(parts, ", ") + " })"
}

// ToLuaInline renders the monitor as a single-line hl.monitor call that fully
// resets the output's rule.
func (m *MonitorSpec) ToLuaInline() string {
	return joinInline(m.luaFields(true))
}

// ToLuaSummary renders a single-line hl.monitor call with default values omitted.
// It is for display only: applying it would keep stale fields of an existing rule.
func (m *MonitorSpec) ToLuaSummary() string {
	return joinInline(m.luaFields(false))
}

// ToLua renders the monitor as an hl.monitor call that fully resets the output's
// rule; disabled monitors stay on a single line, enabled ones get one field per line.
func (m *MonitorSpec) ToLua() string {
	if m.Disabled {
		return m.ToLuaInline()
	}
	var b strings.Builder
	b.WriteString("hl.monitor({\n")
	for _, f := range m.luaFields(true) {
		fmt.Fprintf(&b, "\t%s = %s,\n", f.key, f.value)
	}
	b.WriteString("})")
	return b.String()
}

// ToLua renders every monitor as an hl.monitor call.
func (m MonitorSpecs) ToLua() []string {
	blocks := make([]string, 0, len(m))
	for _, monitor := range m {
		blocks = append(blocks, monitor.ToLua())
	}
	return blocks
}
