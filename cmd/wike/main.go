package main

import (
	"time"
	"wike/internal/settings"
	"wike/internal/win"
)

func main() {
	win.InitScreenSize()
	win.InstallHooks()
	settings.Load()

	go reloadConfigLoop()
	win.RunMessageLoop()
}

func reloadConfigLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		settings.ReloadIfModified()
	}
}
