package main

import (
	"time"
	"wike/internal/config"
	"wike/internal/win"
)

func main() {
	win.InitScreenSize()
	win.InstallHooks()
	config.Load()

	go reloadConfigLoop()
	win.RunMessageLoop()
}

func reloadConfigLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		config.ReloadIfModified()
	}
}
