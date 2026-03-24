# Changelog


## &ensp; [` 📦 v0.5.0  `](https://github.com/kh4f/wike/compare/v0.4.0...v0.5.0)

### &emsp; 🧨 BREAKING CHANGES
- **Renamed config file**: the configuration file is now `config.json` instead of `settings.json` and all "settings" references are replaced with "config". [🡥](https://github.com/kh4f/wike/commit/3bd6b43)
- **Simplified key identifiers**: `trigger.kb` values no longer require the `VK_` prefix, making configuration shorter and more user-friendly. [🡥](https://github.com/kh4f/wike/commit/1c16018)
- **Short mouse button names**: mouse button values in config now use short identifiers (`L`, `R`, `M`, `X1`, `X2`, `WHEEL`, `U`) instead of longer names. [🡥](https://github.com/kh4f/wike/commit/e9d14ec)

### &emsp; 🎁 Features
- **Default rule values**: config fields `name`, `enabled`, and `consume` are now optional with sensible defaults: `name` = "Rule UNK", `enabled` = true, `consume` = true. [🡥](https://github.com/kh4f/wike/commit/fc0a175)
- **App icon and version info**: added embedded icon and version information to the Windows executable. [🡥](https://github.com/kh4f/wike/commit/29458ae)

### &emsp; 🩹 Fixes
- **Consistent executable naming**: the output executable is now consistently named `Wike.exe` for clarity and adherence to naming conventions. [🡥](https://github.com/kh4f/wike/commit/b74cc03)
- **Reliable foreground launching**: improved `Action.launch` so that if a window is already open but in the background, it is now properly brought to the front and focused. [🡥](https://github.com/kh4f/wike/commit/37804db)

##### &emsp;&emsp; [Full Changelog](https://github.com/kh4f/wike/compare/v0.4.0...v0.5.0) &ensp;•&ensp; Mar 24, 2026


## &ensp; [` 📦 v0.4.0  `](https://github.com/kh4f/wike/compare/v0.3.0...v0.4.0)

### &emsp; 🎁 Features
- **Comprehensive key mappings**: expanded `VKCodeMap` to include a complete set of virtual key codes for improved key handling. [🡥](https://github.com/kh4f/wike/commit/29859f8)

### &emsp; 🩹 Fixes
- **Avoid mutating shared default config**: default config is now cloned on each request, preventing side effects from shared references. [🡥](https://github.com/kh4f/wike/commit/7148816)
- **Improve error handling for config operations**: enhanced error handling during config parsing, saving, and reloading, with clear error messages and no silent failures. [🡥](https://github.com/kh4f/wike/commit/c3304f0)

##### &emsp;&emsp; [Full Changelog](https://github.com/kh4f/wike/compare/v0.3.0...v0.4.0) &ensp;•&ensp; Mar 19, 2026


## &ensp; [` 📦 v0.3.0  `](https://github.com/kh4f/wike/compare/v0.2.0...v0.3.0)

### &emsp; 🧨 BREAKING CHANGES
- **Config trigger field renames**: `Trigger.mouse` → `Trigger.m`, `Trigger.key` → `Trigger.kb`, `Trigger.event` → `Trigger.state`. [🡥](https://github.com/kh4f/wike/commit/9601917)
- **Unified app launch/focus**: Use `Action.launch` for both opening and focusing apps; `Action.open` and manual `windowClass` are no longer supported. [🡥](https://github.com/kh4f/wike/commit/f2bf4b6)

### &emsp; 🎁 Features
- **Multiple bindings per rule**: each rule can now define multiple independent trigger-action pairs via a `bindings` array. [🡥](https://github.com/kh4f/wike/commit/3101bd5)

### &emsp; 🩹 Fixes
- **Ignore simulated events**: hooks now skip processing of injected mouse and keyboard events, preventing unintended rule triggers. [🡥](https://github.com/kh4f/wike/commit/000d18c)
- **Negative region coordinates**: regions with negative x or y values are now properly offset from the right screen edge. [🡥](https://github.com/kh4f/wike/commit/2f17b23)

##### &emsp;&emsp; [Full Changelog](https://github.com/kh4f/wike/compare/v0.2.0...v0.3.0) &ensp;•&ensp; Mar 19, 2026


## &ensp; [` 📦 v0.2.0  `](https://github.com/kh4f/wike/compare/v0.1.0...v0.2.0)

### &emsp; 🧨 BREAKING CHANGES
- **Renamed project to Wike**: the project, binary, and repository are now named `Wike` instead of `Twike`. [🡥](https://github.com/kh4f/wike/commit/9041d8b)

### &emsp; 🎁 Features
- **Trigger event state support**: you can now specify `event` (down/up) in the `Trigger` config for more precise activation control. [🡥](https://github.com/kh4f/wike/commit/40d9669)
- **Config file loading & auto-reload**: configuration is now loaded from `config.json`, with changes saved and auto-reloaded every 5 seconds if modified. [🡥](https://github.com/kh4f/wike/commit/9df9575)

### &emsp; ⚙️ Internal
- **Rewrote from Rust to Go**: the entire application was migrated from Rust to Go, with all core logic and hooks reimplemented. [🡥](https://github.com/kh4f/wike/commit/3cdd0e4)

##### &emsp;&emsp; [Full Changelog](https://github.com/kh4f/wike/compare/v0.1.0...v0.2.0) &ensp;•&ensp; Mar 18, 2026


## &ensp; [` 📦 v0.1.0  `](https://github.com/kh4f/wike/commits/v0.1.0)

### &emsp; 🎁 Features
- **Rule-based automation**: introduced a flexible configuration system to trigger commands, keypresses, or window actions based on input events. [🡥](https://github.com/kh4f/wike/commit/3a99c6c) [🡥](https://github.com/kh4f/wike/commit/e54fece) [🡥](https://github.com/kh4f/wike/commit/8b30ff8) [🡥](https://github.com/kh4f/wike/commit/e23821f)
- **Volume scroll**: implemented system volume control by scrolling the mouse wheel on the right edge of the screen. [🡥](https://github.com/kh4f/wike/commit/f85cb79) [🡥](https://github.com/kh4f/wike/commit/6f3c06f)
- **Screen region support**: supported defining interaction zones with negative coordinates for edge-relative positioning. [🡥](https://github.com/kh4f/wike/commit/957844f)
- **Low-level hooks**: established global Windows hooks for capturing keyboard and mouse input. [🡥](https://github.com/kh4f/wike/commit/5e93834)

##### &emsp;&emsp; [Full Changelog](https://github.com/kh4f/wike/commits/v0.1.0) &ensp;•&ensp; Mar 8, 2026