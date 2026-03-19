package app

import (
	"time"
	"wike/internal/config"
	"wike/internal/win"
)

func Run() {
	win.InitScreenSize()
	config.Cfg.Load()

	win.InstallHooks()

	go reloadConfigLoop()

	win.RunMessageLoop()
}

func reloadConfigLoop() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		config.Cfg.ReloadIfModified()
	}
}
