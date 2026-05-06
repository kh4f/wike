<div align="center">
	<picture>
		<source media="(prefers-color-scheme: dark)" srcset="https://raw.githubusercontent.com/kh4f/wike/refs/heads/assets/logo-dark.png">
		<img alt="logo" src="https://raw.githubusercontent.com/kh4f/wike/refs/heads/assets/logo-light.png">
	</picture>
	<br>
	A fast, lightweight, and flexible <b>hotkey manager</b> for Windows
	<br><br>
	<p>
		<a href="https://github.com/kh4f/wike/releases"><img src="https://img.shields.io/github/v/tag/kh4f/wike?label=%F0%9F%8F%B7%EF%B8%8F%20Release&style=flat-square&color=CCCCCC&labelColor=CE2A11" alt="version"/></a>&nbsp;
		<a href="https://github.com/kh4f/wike/issues?q=is%3Aissue+is%3Aopen+label%3Abug"><img src="https://img.shields.io/github/issues/kh4f/wike/bug?label=%F0%9F%90%9B%20Bugs&style=flat-square&color=CCCCCC&labelColor=CE2A11" alt="bugs"></a>&nbsp;
		<a href="https://github.com/kh4f/wike/blob/master/LICENSE"><img src="https://img.shields.io/github/license/kh4f/wike?style=flat-square&label=%F0%9F%9B%A1%EF%B8%8F%20License&color=CCCCCC&labelColor=CE2A11" alt="license"></a>
	</p>
	<p><b>
		<a href="#-features">Features</a>&nbsp; •&nbsp;
		<a href="#-installation">Installation</a>&nbsp; •&nbsp;
		<a href="#%EF%B8%8F-usage">Usage</a>&nbsp; •&nbsp;
		<a href="#%EF%B8%8F-configuration">Configuration</a>
	</b></p>
</div>

## 🔥 Features
- Keyboard & mouse remapping
- Region‑aware hotkeys
- Multi‑key actions & app launching
- Simple YAML-based configuration
- Launch at Windows startup

## 📥 Installation
Download and extract the [latest release](https://github.com/kh4f/wike/releases/latest).

## 🕹️ Usage
Control the app from the command line with `Wike.exe`:

```
Wike v0.5.0

Actions:
  1) Start daemon
  2) Add to startup
  3) Monitor events
  4) Exit
```
- `Start daemon` launches `WikeDaemon.exe` in the background to enable the hotkeys.
- `Add to startup` enables automatic launch on system startup.
- `Monitor events` is useful for debugging - logs all input events and triggered rules in real time.

## ⚙️ Configuration
Wike uses a single `config.yml` file for configuration. Below is a compact example demonstrating the main features:

```yaml
rules:
  - name: Caps Lock → F13 # "Rule UNK" by default
    enabled: true # true by default
    trigger: { kb: CAPITAL } # trigger when Caps Lock is pressed
    action: { kb: [F13] } # simulate pressing F13
    consume: true # prevent the original Caps Lock event (true by default)

  - name: Volume Scroll
    # define multiple bindings within a single rule
    bindings:
      - trigger: { m: WHEEL, state: UP } # mouse wheel up
        action: { kb: [VOLUME_UP] } # increase volume
      - trigger: { m: WHEEL, state: DOWN } # mouse wheel down
        action: { kb: [VOLUME_DOWN] } # decrease volume
    # screen region where the rule is active
    # negative values are relative to the right/bottom edges
    # defaults: x1: 0, y1: 0, x2: <screen width>, y2: <screen height>
    region: { x1: -1, y1: -500 }

  - name: PowerToys Always on Top
    region: { x1: -1 } # right edge of the screen
    trigger: { m: X1 } # mouse back button
    action: { kb: [LWIN, LCONTROL, LSHIFT, F1] } # send a key combination

  - name: Quick Explorer
    region: { x1: 1260, x2: 1500, y1: -2 }
    trigger: { m: L }
    action: { launch: explorer.exe } # launch an app
```

- Supported keyboard keys: [Virtual-Key Codes](https://learn.microsoft.com/en-us/windows/win32/inputdev/virtual-key-codes)
- Supported mouse inputs: `L`, `R`, `M`, `X1`, `X2`, `WHEEL` (with `state: UP/DOWN`)