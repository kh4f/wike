<div align="center">
	<img alt="logo" src="assets/logo.png">
	<br>
	<a href="https://github.com/kh4f/wike/releases"><img src="https://img.shields.io/github/v/tag/kh4f/wike?label=%F0%9F%8F%B7%EF%B8%8F%20Release&style=flat-square&color=EAE2DC&labelColor=303145" alt="version"/></a>&nbsp;
	<a href="https://github.com/kh4f/wike/issues?q=is%3Aissue+is%3Aopen+label%3Abug"><img src="https://img.shields.io/github/issues/kh4f/wike/bug?label=%F0%9F%90%9B%20Bugs&style=flat-square&color=EAE2DC&labelColor=303145" alt="bugs"></a>&nbsp;
	<a href="https://github.com/kh4f/wike/blob/master/LICENSE"><img src="https://img.shields.io/github/license/kh4f/wike?style=flat-square&label=%F0%9F%9B%A1%EF%B8%8F%20License&color=EAE2DC&labelColor=303144" alt="license"></a>&nbsp;
	<br><br>
	A fast, lightweight and flexible <b>hotkey manager</b> for Windows
	<br><br>
	<b>
		<a href="#-features">Features</a>&nbsp; •&nbsp;
		<a href="#-installation">Installation</a>&nbsp; •&nbsp;
		<a href="#%EF%B8%8F-usage">Usage</a>
	</b>
</div>

## 🔥 Features

- Keyboard & mouse remapping
- Region‑aware hotkeys
- Multi‑key actions & app launching
- Simple JSON configuration

## 📥 Installation

Download and extract the [latest release](https://github.com/kh4f/wike/releases/latest).

## 🕹️ Usage

Wike is configured through a single `config.json` file.

Below is a compact example demonstrating the main features:

```jsonc
{
	"rules": [
		{
			"name": "Caps Lock → F13",
			"enabled": true,
			// trigger when Caps Lock is pressed
			"trigger": { "kb": "VK_CAPITAL" },
			// simulate pressing F13
			"action": { "kb": [ "VK_F13" ] },
			// prevent the original Caps Lock event
			"consume": true
		},
		{
			"name": "Volume Scroll",
			"enabled": true,
			// screen region where the rule is active
			// negative x/y are relative to the right/bottom edges
			"region": { "x": -1, "y": -500, "w": 1, "h": 500 },
			// define multiple bindings within a single rule
			"bindings": [
				{
					"trigger": { "m": "WHEEL", "state": "UP" },
					"action": { "kb": [ "VK_VOLUME_UP" ] }
				},
				{
					"trigger": { "m": "WHEEL", "state": "DOWN" },
					"action": { "kb": [ "VK_VOLUME_DOWN" ] }
				}
			],
			"consume": true
		},
		{
			"name": "Toggle PowerToys Always on Top",
			"enabled": true,
			// the right edge of the screen
			"region": { "x": -1, "y": 0, "w": 1, "h": 1080 },
			"trigger": { "kb": "VK_PAUSE" },
			// send a key combination (Win+Ctrl+Shift+F1)
			"action": { "kb": [ "VK_LWIN", "VK_LCONTROL", "VK_LSHIFT", "VK_F1" ] },
			"consume": true
		},
		{
			"name": "Quick Explorer",
			"enabled": true,
			// small region on the right side of the taskbar
			"region": { "x": -660, "y": -2, "w": 240, "h": 5 },
			"trigger": { "m": "LMB" },
			// launch the app (or focus it if already running)
			"action": { "launch": "explorer.exe" },
			"consume": true
		}
	]
}
```

Notes:
- Supported keyboard keys: [Virtual-Key Codes](https://learn.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes)
- Supported mouse inputs: `LMB`, `RMB`, `MMB`, `X1`, `X2`, `WHEEL` (with `state` `UP`/`DOWN`)

> [!WARNING]
> Wike is in early development — expect breaking changes

</br>

<div align="center">
  <b>MIT License © 2026 <a href="https://github.com/kh4f">kh4f</a></b>
</div>
