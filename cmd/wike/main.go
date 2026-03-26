package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	version         = "0.5.0"
	daemonFileName  = "WikeDaemon.exe"
	logPipeName     = `\\.\pipe\wike-events`
	controlPipeName = `\\.\pipe\wike-control`
)

var banner = fmt.Sprintf(`Wike v%s`, version)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		clearConsole()
		printMenu()

		choice, _ := reader.ReadString('\n')

		switch strings.TrimSpace(choice) {
		case "1":
			if daemonRunning() {
				stopDaemon()
			} else {
				startDaemon()
			}
			time.Sleep(500 * time.Millisecond)
		case "2":
			handle := openPipe(logPipeName, windows.GENERIC_READ)
			if handle == windows.InvalidHandle {
				fmt.Println("Daemon not running")
				time.Sleep(500 * time.Millisecond)
				continue
			}

			clearConsole()
			showLogs()

			pipe := os.NewFile(uintptr(handle), logPipeName)

			go func() {
				io.Copy(os.Stdout, pipe)
			}()

			reader.ReadString('\n')
			pipe.Close()
		case "3":
			return
		}
	}
}

func printMenu() {
	daemonAction := "Start daemon"
	if daemonRunning() {
		daemonAction = "Stop daemon"
	}

	fmt.Printf(`%s

Actions:
  1) %s
  2) Monitor events
  3) Exit

> `, banner, daemonAction)
}

func showLogs() {
	fmt.Printf(`%s

Monitoring events...
(Enter to go back)

`, banner)
}

func clearConsole() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func openPipe(name string, access uint32) windows.Handle {
	path, _ := windows.UTF16PtrFromString(name)
	handle, _ := windows.CreateFile(
		path,
		access,
		0,
		nil,
		windows.OPEN_EXISTING,
		windows.FILE_ATTRIBUTE_NORMAL,
		0,
	)
	return handle
}

func daemonRunning() bool {
	snapshot, err := windows.CreateToolhelp32Snapshot(windows.TH32CS_SNAPPROCESS, 0)
	if err != nil {
		return false
	}
	defer windows.CloseHandle(snapshot)

	var entry windows.ProcessEntry32
	entry.Size = uint32(unsafe.Sizeof(entry))

	if err := windows.Process32First(snapshot, &entry); err != nil {
		return false
	}

	for {
		if windows.UTF16ToString(entry.ExeFile[:]) == daemonFileName {
			return true
		}

		if err := windows.Process32Next(snapshot, &entry); err != nil {
			return false
		}
	}
}

func startDaemon() {
	if daemonRunning() {
		fmt.Println("Daemon already running")
		time.Sleep(500 * time.Millisecond)
		return
	}

	exePath, _ := os.Executable()
	daemonPath := filepath.Join(filepath.Dir(exePath), daemonFileName)
	cmd := exec.Command(daemonPath)

	if err := cmd.Start(); err != nil {
		fmt.Println("Failed to start daemon:", err)
	} else {
		fmt.Println("Daemon started")
	}

	time.Sleep(500 * time.Millisecond)
}

func stopDaemon() {
	handle := openPipe(controlPipeName, windows.GENERIC_WRITE)
	if handle == windows.InvalidHandle {
		fmt.Println("Daemon not running")
		time.Sleep(500 * time.Millisecond)
		return
	}

	pipe := os.NewFile(uintptr(handle), controlPipeName)
	pipe.WriteString("stop\n")
	pipe.Close()

	fmt.Println("Daemon stopped")
	time.Sleep(500 * time.Millisecond)
}
