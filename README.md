
<img src="https://github.com/user-attachments/assets/0effc242-3d3d-4d39-a183-0a567c4da3a9" width="90" style="margin-right:10px" align=left alt="hyprdynamicmonitors logo">
<H1>HyprDynamicMonitors</H1><br>


[![Documentation](https://img.shields.io/badge/docs-fiffeek.github.io-blue)](https://fiffeek.github.io/hyprdynamicmonitors/)
![Release](https://img.shields.io/github/v/release/fiffeek/hyprdynamicmonitors.svg?style=flat)
![Downloads](https://img.shields.io/github/downloads/fiffeek/hyprdynamicmonitors/total.svg?style=flat)
![License](https://img.shields.io/github/license/fiffeek/hyprdynamicmonitors.svg?style=flat)
![Code Size](https://img.shields.io/github/languages/code-size/fiffeek/hyprdynamicmonitors.svg?style=flat)
![Go Version](https://img.shields.io/github/go-mod/go-version/fiffeek/hyprdynamicmonitors?style=flat)
![AUR version](https://img.shields.io/aur/version/hyprdynamicmonitors-bin?style=flat&label=AUR)
![Go Report Card](https://goreportcard.com/badge/github.com/fiffeek/hyprdynamicmonitors)
![Build](https://img.shields.io/github/actions/workflow/status/fiffeek/hyprdynamicmonitors/test.yaml?branch=main&style=flat)
[![Codecov](https://codecov.io/github/fiffeek/hyprdynamicmonitors/graph/badge.svg?token=HZONCOKJTE)](https://codecov.io/github/fiffeek/hyprdynamicmonitors)

An event-driven service with an interactive TUI that automatically manages Hyprland monitor configurations based on connected displays, power and lid state.

> [!IMPORTANT]
> **Hyprland Lua only.** This fork targets Hyprland's Lua configuration (`~/.config/hypr/hyprland.lua`).
> It renders `hl.monitor({ ... })` calls instead of hyprlang `monitor=` lines, and it no longer uses
> `hyprctl keyword` (rejected by Hyprland under a Lua config). If you still use `hyprland.conf`, use upstream
> [fiffeek/hyprdynamicmonitors](https://github.com/fiffeek/hyprdynamicmonitors). See [Hyprland Lua](#hyprland-lua) for how to migrate.


## Preview


![TUI theming](./preview/output/themes-showcase.gif)

<details>
<summary>TUI demo in default theme</summary>

![TUI demo](./preview/output/demo.gif)

</details>

<details>
<summary>Daemon Response to Power State Events</summary>

![power events demo](./preview/output/power_events.gif)

</details>

<details>
<summary>Daemon and TUI Response to Laptop Lid Events</summary>

![lid tui events demo](./preview/output/lid_tui.gif)

</details>

<details>
<summary>Dynamic Profile Modes depending on the user context</summary>

A full guide is available at [Dynamic Profile Modes](https://fiffeek.github.io/hyprdynamicmonitors/docs/guides/ad-hoc).

![Dynamic Profile Modes](./preview/output/ad-hoc.gif)

</details>

## Table Of Contents

<!--ts-->
* [HyprDynamicMonitors](#hyprdynamicmonitors)
   * [Preview](#preview)
   * [Table Of Contents](#table-of-contents)
   * [Features](#features)
   * [Quick Start](#quick-start)
      * [Installation](#installation)
      * [Setup](#setup)
   * [Hyprland Lua](#hyprland-lua)
      * [Migrating an existing template](#migrating-an-existing-template)
      * [Live apply](#live-apply)
         * [Profile-specific keybinds](#profile-specific-keybinds)
   * [Documentation](#documentation)
   * [Examples](#examples)
   * [Runtime Requirements](#runtime-requirements)
   * [Alternative Software](#alternative-software)
   * [Support](#support)
   * [License](#license)
<!--te-->

## Features

- **Event-driven architecture** responding to monitor, power and lid state changes in real-time
- **Interactive TUI** for visual monitor configuration and profile management
- **Profile-based configuration** with different settings for different monitor setups
- **Template support** for dynamic configuration generation
- **Hot reloading** automatically detects and applies configuration changes without restart (optional)
- **Power state awareness** built-in AC/battery detection for laptop users (optional)
- **Lid state awareness** built-in lid state detection for laptop users (optional)
- **Desktop notifications** for configuration changes (optional)
- **Theming support** for built-in themes and dynamic generation from wallpapers

## Quick Start

### Installation

```bash
# Binary release
export DESTDIR="$HOME/.local/bin"  # optional, defaults to ~/.local/bin/
curl -o- https://raw.githubusercontent.com/fiffeek/hyprdynamicmonitors/refs/heads/main/scripts/install.sh | bash

# AUR (Arch Linux)
$aurHelper -S hyprdynamicmonitors-bin

# Nix, for flakes/modules see: https://hyprdynamicmonitors.filipmikina.com/docs/advanced/systemd#nix
nix run github:fiffeek/hyprdynamicmonitors
```

See the [Installation Guide](https://fiffeek.github.io/hyprdynamicmonitors/docs/quickstart/installation) for more options.

### Setup

The easiest way to get started is using the TUI (see the [Getting started with the TUI](https://hyprdynamicmonitors.filipmikina.com/docs/quickstart/tui/) guide):

```bash
# Launch the interactive TUI
hyprdynamicmonitors tui

# Configure your monitors visually
# Press 'Tab' to switch to Profile view
# Press 'n' to create a new profile
# For getting started see: https://hyprdynamicmonitors.filipmikina.com/docs/quickstart/tui/
```

Then add to your `~/.config/hypr/hyprland.lua` (alternatively, see [Running with systemd](https://fiffeek.github.io/hyprdynamicmonitors/docs/advanced/systemd) guide):

```lua
-- Load the generated monitor configuration (~/.config/hypr/hyprdynamicmonitors.lua).
-- Put this AFTER any other hl.monitor(...) calls so the generated rules win.
require("hyprdynamicmonitors")

-- Run the daemon for automatic profile switching
hl.on("hyprland.start", function()
	hl.exec_cmd("hyprdynamicmonitors run")
end)

-- or use systemd instead
-- see: https://fiffeek.github.io/hyprdynamicmonitors/docs/advanced/systemd
-- systemctl --user enable --now hyprdynamicmonitors.service
```

The generated file is written to `$HOME/.config/hypr/hyprdynamicmonitors.lua` by default (`[general] destination`).
Avoid naming it `monitors.lua`: `require("monitors")` would be shadowed by (or shadow) a `monitors/init.lua` module.

Ensure you're running `hyprdynamicmonitors prepare` prior to running Hyprland (it rewrites `disabled = true` to `disabled = false` in the destination so Hyprland never starts without an active display; see [the guide and the why](https://fiffeek.github.io/hyprdynamicmonitors/docs/advanced/prepare)):
```bash
# e.g. from a tty:
hyprdynamicmonitors prepare && Hyprland
# or use systemd
# systemctl --user enable hyprdynamicmonitors-prepare.service
```

For detailed setup instructions, see the [Quick Start Guide](https://fiffeek.github.io/hyprdynamicmonitors/docs/category/quick-start).

> [!CAUTION]
> For production environments prefer running with `systemd` ([guide](https://fiffeek.github.io/hyprdynamicmonitors/docs/advanced/systemd)).

## Hyprland Lua

Profiles still point at either a static file (`config_file_type = "static"`, symlinked to the destination) or a
Go [`text/template`](https://pkg.go.dev/text/template) (`config_file_type = "template"`). The template data and
functions are unchanged (`.MonitorsByTag`, `.RequiredMonitors`, `.ExtraMonitors`, `.Monitors`, `isOnAC`,
`isOnBattery`, `isLidClosed`, ...), but the **output must be Lua**. By convention templates are named `*.lua.tmpl`
(`hyprdynamicmonitors freeze` creates `hyprconfigs/<profile>.lua.tmpl`) and static files `*.lua`.

### Migrating an existing template

| hyprlang (old) | Lua (new) |
|---|---|
| `monitor=desc:{{$ext.Description}},2560x1440@60,0x0,1` | `hl.monitor({ output = "desc:{{$ext.Description}}", mode = "2560x1440@60", position = "0x0", scale = 1 })` |
| `monitor={{$laptop.Name}},2880x1920@120,0x0,2,vrr,1,transform,0` | `hl.monitor({ output = "{{$laptop.Name}}", mode = "2880x1920@120", position = "0x0", scale = 2, vrr = true, transform = 0 })` |
| `monitor={{.Name}},disable` | `hl.monitor({ output = "{{.Name}}", disabled = true })` |
| `monitor=,preferred,auto,1` | `hl.monitor({ output = "", mode = "preferred", position = "auto", scale = 1 })` |
| `workspace = 2,monitor:{{$ext.Name}},default:true` | `hl.workspace_rule({ workspace = "2", monitor = "{{$ext.Name}}", default = true })` |
| `exec = hyprctl dispatch moveworkspacetomonitor 1 eDP-1` | `hl.dispatch(hl.dsp.workspace.move({ workspace = "1", monitor = "eDP-1" }))` |
| `decoration { blur { enabled = false } }` | `hl.config({ decoration = { blur = { enabled = false } } })` |
| `# comment` / `# <<<<< TUI AUTO START` | `-- comment` / `-- <<<<< TUI AUTO START` |

A complete template:

```lua
{{- $laptop := index .MonitorsByTag "laptop" -}}
{{- $external := index .MonitorsByTag "external" -}}
-- Generated for power state: {{.PowerState}}
hl.monitor({
	output = "{{$laptop.Name}}",
	mode = "2880x1920@{{if isOnBattery}}60{{else}}120{{end}}",
	position = "0x0",
	scale = 2,
})
{{- if $external }}
hl.monitor({ output = "desc:{{$external.Description}}", mode = "2560x1440@60", position = "auto-right", scale = 1 })
hl.workspace_rule({ workspace = "2", monitor = "desc:{{$external.Description}}" })
{{- end }}
```

Notes:
- Descriptions no longer need the hyprlang `##` escaping of `#`; they sit inside a Lua string. A description containing
  a double quote `"` (or a backslash) would need escaping (`\"`) inside the string.
- Supported `hl.monitor` fields: `output`, `mode`, `position`, `scale`, `transform`, `vrr`, `bitdepth`, `cm`,
  `sdrbrightness`, `sdrsaturation`, `mirror`, `disabled`, ... (see `HL.MonitorSpec` in Hyprland's
  `/usr/share/hypr/stubs/hl.meta.lua`).
- The TUI only rewrites the region between `-- <<<<< TUI AUTO START` and `-- <<<<< TUI AUTO END`; keep it at the end
  of the file so its monitor rules win.
- You can syntax-check a static file with `luac -p file.lua`.

### Live apply

```toml
[general]
destination = "$HOME/.config/hypr/hyprdynamicmonitors.lua"
live_apply = true # default
```

With `live_apply = true` the daemon (and the TUI) runs the destination inside Hyprland right after writing it, the socket
equivalent of `hyprctl eval 'dofile("<destination>")'`, so the change takes effect immediately. Because the file is
re-run on every apply, it should only contain idempotent calls: `hl.monitor`, `hl.workspace_rule`, `hl.config`,
`hl.window_rule` *with a `name`* (named rules are updated, unnamed ones are added again). Keep `hl.bind`,
`hl.on(...)` and autostart `hl.exec_cmd(...)` in `hyprland.lua` - in the destination they would be registered or run again
on every apply. Set `live_apply = false` to rely on Hyprland reloading its config instead.

#### Profile-specific keybinds

A plain `hl.bind` in a profile stacks a duplicate on every apply (one key press then fires the action several times;
toggles cancel out). If a profile really needs its own bind, unbind the combo first so re-running replaces it:

```lua
-- Replace any existing bind on this combo instead of stacking a duplicate on every apply.
local function rebind(keys, dispatcher, opts)
	hl.unbind(keys)
	return hl.bind(keys, dispatcher, opts)
end

rebind("SUPER + P", hl.dsp.exec_cmd("some-profile-specific-thing"))
```

- Spell the combo exactly like the bind you are replacing. Spacing does not matter, but modifier order does:
  `hl.unbind("CTRL + SUPER + P")` does not remove a bind made as `"SUPER + CTRL + P"`.
- This overwrites globally, including binds from `hyprland.lua`, and switching to another profile does not restore
  the original - it only comes back on a full Hyprland reload. Either `rebind` the combo in every profile, or keep one
  bind in `hyprland.lua` and decide what it does when pressed (bind a Lua function).
- Don't use the handle returned by `hl.bind` (`:remove()`) for this: it removes every bind on the combo, so the first
  apply keeps your base bind and later applies silently delete it.

## Documentation

**Full documentation is available at [fiffeek.github.io/hyprdynamicmonitors](https://fiffeek.github.io/hyprdynamicmonitors/)**

Key topics:

| Topic | Description |
|-------|-------------|
| [Quick Start](https://fiffeek.github.io/hyprdynamicmonitors/docs/category/quick-start) | Get up and running quickly |
| [TUI Guide](https://fiffeek.github.io/hyprdynamicmonitors/docs/quickstart/tui) | Interactive monitor configuration |
| [Configuration](https://fiffeek.github.io/hyprdynamicmonitors/docs/category/configuration) | Profiles, monitors, and power management |
| [Templates](https://fiffeek.github.io/hyprdynamicmonitors/docs/advanced/templates) | Dynamic configuration generation |
| [Running with systemd](https://fiffeek.github.io/hyprdynamicmonitors/docs/advanced/systemd) | Production deployment |
| [CLI Commands](https://fiffeek.github.io/hyprdynamicmonitors/docs/usage/commands) | Command reference |
| [FAQ](https://fiffeek.github.io/hyprdynamicmonitors/docs/faq) | Common questions |
| [Power Events](https://fiffeek.github.io/hyprdynamicmonitors/docs/guides/power-events) | Guide to running with Power Events |
| [Lid Events](https://fiffeek.github.io/hyprdynamicmonitors/docs/guides/lid-events) | Guide to running with Lid Events |
| [Dynamic Profile Modes](https://fiffeek.github.io/hyprdynamicmonitors/docs/guides/ad-hoc) | Guide to making dynamic profile modes |
| [Prepare Command](https://fiffeek.github.io/hyprdynamicmonitors/docs/advanced/prepare) | Do I need `hyprdynamicmonitors prepare`? |
| [Make it yours!](https://fiffeek.github.io/hyprdynamicmonitors/docs/configuration/theming) | How to apply/generate a theme |

## Examples

See the [`examples/`](examples) directory for complete configuration examples:

| Example | Description |
|---------|-------------|
| [Basic Setup](examples/basic) | Simple laptop configuration |
| [Full Configuration](examples/full) | All available options |
| [Power States](examples/power-states) | AC/battery-aware profiles |
| [Lid States](examples/lid-states) | Laptop lid detection |
| [Template Variables](examples/template-variables) | Dynamic templates |
| [Disable Monitors](examples/disable-monitors) | Managing unexpected displays |

## Runtime Requirements

- Hyprland with a Lua config (`hyprland.lua`) and IPC support
- UPower (required on laptops, unless `--disable-power-events` is passed, for power state monitoring)
- D-Bus access (required if power events, lid state or notifications are enabled)

## Alternative Software

Similar tools worth checking out:
- [kanshi](https://sr.ht/~emersion/kanshi/) - Generic Wayland output management
- [shikane](https://github.com/hw0lff/shikane) - Another Wayland output manager
- [nwg-displays](https://github.com/nwg-piotr/nwg-displays) - GUI-based display configuration tool
- [hyprmon](https://github.com/erans/hyprmon) - TUI-based display configuration tool

HyprDynamicMonitors is Hyprland-specific but offers deeper integration, an interactive TUI, template system, and power state awareness. See [Introduction](https://fiffeek.github.io/hyprdynamicmonitors/docs/) for a detailed comparison.

## Support

If you find HyprDynamicMonitors useful and want to support its development, consider buying me a coffee on [Ko-fi](https://ko-fi.com/filipmikina). Your support helps keep the project maintained and enables new features. Thank you!

## License

See [LICENSE](LICENSE) file.
