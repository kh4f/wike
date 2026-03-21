package app

import (
	"time"
	"wike/internal/settings"
	"wike/internal/win"
)

func main() {
	win.InitScreenSize()
	config.Settings.Load()

	win.InstallHooks()

	go reloadConfigLoop()

	win.RunMessageLoop()
}

func reloadConfigLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		config.Settings.ReloadIfModified()
	}
}