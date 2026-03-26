package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"time"

	"golang.org/x/sys/windows"
)

const (
	version  = "0.5.0"
	pipeName = `\\.\pipe\wike`
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
			path, _ := windows.UTF16PtrFromString(pipeName)
			handle, _ := windows.CreateFile(
				path,
				windows.GENERIC_READ,
				0,
				nil,
				windows.OPEN_EXISTING,
				windows.FILE_ATTRIBUTE_NORMAL,
				0,
			)
			if handle == windows.InvalidHandle {
				fmt.Println("Daemon is not running")
				time.Sleep(time.Second)
				continue
			}

			clearConsole()
			showLogs()

			pipe := os.NewFile(uintptr(handle), pipeName)

			go func() {
				io.Copy(os.Stdout, pipe)
			}()

			reader.ReadString('\n')
			pipe.Close()
		case "2":
			return
		}
	}
}

func printMenu() {
	fmt.Printf(`%s

Actions:
  1) Monitor events
  2) Exit

> `, banner)
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
