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

func (m *MonitorSpec) luaFields() []luaField {
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
	if m.TenBitdepth {
		fields = append(fields, luaField{"bitdepth", "10"})
	}
	if m.HasNonDefaultColorPreset() {
		fields = append(fields, luaField{"cm", LuaString(m.ColorPreset)})
	}
	if m.HDR() && m.SdrBrightness != 1.0 {
		fields = append(fields, luaField{"sdrbrightness", luaNumber(m.SdrBrightness)})
	}
	if m.HDR() && m.SdrSaturation != 1.0 {
		fields = append(fields, luaField{"sdrsaturation", luaNumber(m.SdrSaturation)})
	}
	if m.HasMirror() {
		fields = append(fields, luaField{"mirror", LuaString(m.Mirror)})
	}
	return fields
}

// ToLuaInline renders the monitor as a single-line hl.monitor call.
func (m *MonitorSpec) ToLuaInline() string {
	parts := []string{}
	for _, f := range m.luaFields() {
		parts = append(parts, f.key+" = "+f.value)
	}
	return "hl.monitor({ " + strings.Join(parts, ", ") + " })"
}

// ToLua renders the monitor as an hl.monitor call; disabled monitors stay on
// a single line, enabled ones get one field per line.
func (m *MonitorSpec) ToLua() string {
	if m.Disabled {
		return m.ToLuaInline()
	}
	var b strings.Builder
	b.WriteString("hl.monitor({\n")
	for _, f := range m.luaFields() {
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
