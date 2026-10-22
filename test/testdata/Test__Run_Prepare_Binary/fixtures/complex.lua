-- <<<< TUI
-- whatever
hl.monitor({
	output = "eDP-1",
	mode = "2880x1920@120",
	position = "0x0",
	scale = 2,
	vrr = true,
	disabled = false,
})
hl.monitor({ output = "DP-11", disabled = true })
hl.workspace_rule({ workspace = "1", monitor = "eDP-1", default = true })
hl.monitor({ output = "DP-12", disabled=true })
hl.monitor({ output = "DP-13",   disabled   =   true   })
hl.monitor({
	output = "DP-14",
	disabled = true,
})
hl.monitor({ output = "DP-15", disabled = trueish })

hl.config({
	decoration = {
		blur = {
			enabled = false, -- Disable blur for performance
		},
	},
	animations = {
		enabled = false, -- Disable animations for performance
	},
	misc = {
		vrr = 3,
	},
})
hl.exec_cmd("brightnessctl set 100")
-- comment
