-- Some existing config
require("looknfeel")

-- Existing settings
hl.bind("SUPER + Q", hl.dsp.exec_cmd("kitty"))
hl.bind("SUPER + C", hl.dsp.window.close())
-- <<<<< TUI AUTO START
hl.monitor({
	output = "desc:New Monitor A",
	mode = "2560x1440@120",
	position = "0x0",
	scale = 1.5,
	transform = 0,
	vrr = false,
	disabled = false,
	bitdepth = 8,
	cm = "srgb",
	sdrbrightness = 1,
	sdrsaturation = 1,
	mirror = "",
})
hl.monitor({
	output = "desc:New Monitor B",
	mode = "1920x1080@60",
	position = "2560x0",
	scale = 1,
	transform = 0,
	vrr = true,
	disabled = false,
	bitdepth = 10,
	cm = "hdr",
	sdrbrightness = 1.1,
	sdrsaturation = 0.98,
	mirror = "eDP-1",
})
hl.monitor({
	output = "monC",
	mode = "1000x1000@60",
	position = "-1000x-1000",
	scale = 1,
	transform = 0,
	vrr = false,
	disabled = false,
	bitdepth = 8,
	cm = "srgb",
	sdrbrightness = 1,
	sdrsaturation = 1,
	mirror = "",
})
hl.monitor({
	output = "desc:Dell #Whatever",
	mode = "1337x500@60",
	position = "2137x-5000",
	scale = 2,
	transform = 2,
	vrr = false,
	disabled = false,
	bitdepth = 8,
	cm = "srgb",
	sdrbrightness = 1,
	sdrsaturation = 1,
	mirror = "",
})
-- <<<<< TUI AUTO END
