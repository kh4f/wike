package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
	"wike/internal/config"
	"wike/internal/logger"
	"wike/internal/win"

	"golang.org/x/sys/windows"
)

const (
	logPipeName     = `\\.\pipe\wike-events`
	controlPipeName = `\\.\pipe\wike-control`
)

func main() {
	go serveControl()
	go serveLogs()
	runApp()
}

func runApp() {
	logger.Println("Starting daemon...")
	win.InitScreenSize()

	if err := config.Load(); err != nil {
		logger.Println("Config load error:", err)
	}

	win.InstallHooks()
	go reloadConfigLoop()

	logger.Println("Daemon ready")
	win.RunMessageLoop()
}

func reloadConfigLoop() {
	for {
		time.Sleep(time.Second * 5)
		config.ReloadIfModified()
	}
}

func serveLogs() {
	path, _ := windows.UTF16PtrFromString(logPipeName)
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
		pipe := os.NewFile(uintptr(handle), logPipeName)

		windows.ConnectNamedPipe(windows.Handle(pipe.Fd()), nil)
		clientID++
		logger.AddOutput(fmt.Sprintf("pipe-%d", clientID), pipe)
		logger.Println("log monitor connected")
	}
}

func serveControl() {
	path, _ := windows.UTF16PtrFromString(controlPipeName)

	for {
		handle, _ := windows.CreateNamedPipe(
			path,
			windows.PIPE_ACCESS_INBOUND,
			windows.PIPE_TYPE_BYTE|windows.PIPE_WAIT,
			1,
			4096,
			4096,
			0,
			nil,
		)
		pipe := os.NewFile(uintptr(handle), controlPipeName)

		windows.ConnectNamedPipe(windows.Handle(pipe.Fd()), nil)

		command, _ := bufio.NewReader(pipe).ReadString('\n')
		pipe.Close()

		if strings.TrimSpace(command) == "stop" {
			os.Exit(0)
		}
	}
}
