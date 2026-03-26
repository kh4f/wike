package main

import (
	"fmt"
	"os"
	"time"
	"wike/internal/config"
	"wike/internal/logger"
	"wike/internal/win"

	"golang.org/x/sys/windows"
)

const pipeName = `\\.\pipe\wike`

func main() {
	go serveLogs()
	runApp()
}

func runApp() {
	win.InitScreenSize()
	win.InstallHooks()

	if err := config.Load(); err != nil {
		logger.Println("Config load error:", err)
	}
	go reloadConfigLoop()

	win.RunMessageLoop()
}

func reloadConfigLoop() {
	for {
		time.Sleep(time.Second * 5)
		config.ReloadIfModified()
	}
}

func serveLogs() {
	path, _ := windows.UTF16PtrFromString(pipeName)
	clientID := 0

	for {
		handle, _ := windows.CreateNamedPipe(
			path,
			windows.PIPE_ACCESS_OUTBOUND,
			windows.PIPE_TYPE_BYTE|windows.PIPE_WAIT,
			windows.PIPE_UNLIMITED_INSTANCES,
			4096,
			4096,
			0,
			nil,
		)
		pipe := os.NewFile(uintptr(handle), pipeName)

		windows.ConnectNamedPipe(windows.Handle(pipe.Fd()), nil)
		clientID++
		logger.AddOutput(fmt.Sprintf("pipe-%d", clientID), pipe)
	}
}
