package win

import (
	"fmt"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	SW_RESTORE = 9
	SW_SHOW    = 5
)

var (
	getClassNameW       = user32.NewProc("GetClassNameW").Call
	getWindowTextW      = user32.NewProc("GetWindowTextW").Call
	isWindowVisible     = user32.NewProc("IsWindowVisible").Call
	showWindow          = user32.NewProc("ShowWindow").Call
	setForegroundWindow = user32.NewProc("SetForegroundWindow").Call
	shellExecuteW       = shell32.NewProc("ShellExecuteW").Call
)

func openOrFocus(proc string) {
	proc = normalizePath(proc)

	fmt.Printf("Trying to find windows for process: %q\n", proc)

	found := 0
	focused := 0

	err := windows.EnumWindows(syscall.NewCallback(func(hwnd windows.HWND, _ uintptr) uintptr {
		pid, exe, fullProcPath := getProcessInfo(hwnd)
		if !strings.HasSuffix(normalizePath(fullProcPath), proc) {
			return 1
		}

		className := getClassName(hwnd)
		title := getWindowText(hwnd)
		found++
		if shouldFocusWindow(hwnd, className, title) {
			fmt.Printf("Focusing window: class=%q | title=%q | exe=%q | proc=%q | pid=%d | hwnd=0x%X\n",
				className,
				title,
				exe,
				fullProcPath,
				pid,
				uintptr(hwnd),
			)
			focusWindow(hwnd)
			focused++
		}
		return 1
	}), nil)
	if err != nil {
		fmt.Println("EnumWindows err:", err)
	}

	fmt.Printf("\nMatched windows: %d, focused: %d\n", found, focused)

	if focused == 0 {
		fmt.Println("No windows matched the proc filter, opening: ", proc)
		openWindow(proc)
		return
	}
}

func getClassName(hwnd windows.HWND) string {
	buf := make([]uint16, 256)
	n, _, _ := getClassNameW(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)
	if n == 0 {
		return ""
	}
	return windows.UTF16ToString(buf[:n])
}

func getWindowText(hwnd windows.HWND) string {
	buf := make([]uint16, 512)
	n, _, _ := getWindowTextW(
		uintptr(hwnd),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(len(buf)),
	)
	if n == 0 {
		return ""
	}
	return windows.UTF16ToString(buf[:n])
}

func getProcessInfo(hwnd windows.HWND) (uint32, string, string) {
	var pid uint32
	windows.GetWindowThreadProcessId(hwnd, &pid)
	if pid == 0 {
		return 0, "", ""
	}

	h, err := windows.OpenProcess(windows.PROCESS_QUERY_LIMITED_INFORMATION, false, pid)
	if err != nil || h == 0 {
		return pid, "", ""
	}
	defer windows.CloseHandle(h)

	buf := make([]uint16, windows.MAX_PATH*4)
	size := uint32(len(buf))
	if err := windows.QueryFullProcessImageName(h, 0, &buf[0], &size); err != nil {
		return pid, "", ""
	}

	fullPath := windows.UTF16ToString(buf[:size])
	return pid, filepath.Base(fullPath), fullPath
}

func shouldFocusWindow(hwnd windows.HWND, className, title string) bool {
	if hwnd == 0 || strings.TrimSpace(title) == "" {
		return false
	}

	visible, _, _ := isWindowVisible(uintptr(hwnd))
	if visible == 0 {
		return false
	}

	switch className {
	case "IME", "MSCTFIME UI", "Progman":
		return false
	}

	return true
}

func openWindow(path string) {
	shellExecuteW(
		0,
		uintptr(unsafe.Pointer(utf16("open"))),
		uintptr(unsafe.Pointer(utf16(path))),
		0,
		0,
		uintptr(SW_SHOW),
	)
}

func focusWindow(hwnd windows.HWND) {
	if hwnd != 0 {
		isWindow := windows.IsWindow(hwnd)
		if isWindow {
			showWindow(uintptr(hwnd), uintptr(SW_RESTORE))
			setForegroundWindow(uintptr(hwnd))
		}
	}
}

func normalizePath(p string) string {
	p = strings.TrimSpace(p)
	p = strings.ReplaceAll(p, "/", `\`)
	return strings.ToLower(p)
}
