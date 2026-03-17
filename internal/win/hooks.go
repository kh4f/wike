package win

import (
	"fmt"
	"os/exec"
	"unsafe"
	"wike/internal/config"
	"wike/internal/shared"

	"golang.org/x/sys/windows"
)

var (
	shell32             = windows.NewLazySystemDLL("shell32.dll")
	SetWindowsHookExW   = user32.NewProc("SetWindowsHookExW").Call
	CallNextHookEx      = user32.NewProc("CallNextHookEx").Call
	GetMessageW         = user32.NewProc("GetMessageW").Call
	GetCursorPos        = user32.NewProc("GetCursorPos").Call
	FindWindowW         = user32.NewProc("FindWindowW").Call
	IsWindow            = user32.NewProc("IsWindow").Call
	ShowWindow          = user32.NewProc("ShowWindow").Call
	SetForegroundWindow = user32.NewProc("SetForegroundWindow").Call
	ShellExecuteW       = shell32.NewProc("ShellExecuteW").Call
)

const (
	SW_RESTORE = 9
	SW_SHOW    = 5
)

func mHook(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 {
		info := (*MSLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		x, y := info.Pt.X, info.Pt.Y
		mEvent := mouseEventMap[wParam]
		if mEvent != "Move" {
			fmt.Printf("Mouse event: wParam=0x%X (%s) at (%d, %d)\n", wParam, mEvent, x, y)
		}

		for _, rule := range config.Cfg.Rules {
			if !rule.Enabled ||
				rule.Trigger.Mouse == nil ||
				*rule.Trigger.Mouse != mEvent ||
				rule.Region != nil && !rule.Region.Contains(info.Pt) {
				continue
			}

			fmt.Printf("Rule triggered: '%s' at (%d, %d)\n", rule.Name, x, y)

			executeAction(rule.Action)

			if rule.Consume != nil && *rule.Consume {
				return 1
			}
		}
	}

	ret, _, _ := CallNextHookEx(0, uintptr(nCode), wParam, lParam)
	return ret
}

func kHook(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 {
		info := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		var pt shared.POINT
		GetCursorPos(uintptr(unsafe.Pointer(&pt)))
		fmt.Printf("Key event: wParam=0x%X vkCode=0x%X at (%d, %d)\n", wParam, info.VkCode, pt.X, pt.Y)

		for _, rule := range config.Cfg.Rules {
			if !rule.Enabled ||
				rule.Trigger.Key == nil ||
				virtualKeyMap[*rule.Trigger.Key] != uint16(info.VkCode) ||
				rule.Region != nil && !rule.Region.Contains(pt) {
				continue
			}

			fmt.Printf("Rule triggered: '%s' at (%d, %d)\n", rule.Name, pt.X, pt.Y)

			executeAction(rule.Action)

			if rule.Consume != nil && *rule.Consume {
				return 1
			}
		}
	}

	ret, _, _ := CallNextHookEx(0, uintptr(nCode), wParam, lParam)
	return ret
}

func executeAction(action config.Action) {
	if len(action.Keys) > 0 {
		var vkeys []uint16
		for _, keyName := range action.Keys {
			if code, ok := virtualKeyMap[keyName]; ok {
				vkeys = append(vkeys, code)
			}
		}
		if len(vkeys) > 0 {
			pressKeys(vkeys)
		}
	}

	if action.Cmd != nil {
		exec.Command(*action.Cmd).Run()
	}

	if action.Open != nil {
		openOrFocus(*action.Open)
	}
}

func openOrFocus(open config.Open) {
	if open.WindowClass != nil {
		classPtr := utf16(*open.WindowClass)
		hwnd, _, _ := FindWindowW(uintptr(unsafe.Pointer(classPtr)), 0)

		if hwnd != 0 {
			isWindow, _, _ := IsWindow(hwnd)
			if isWindow != 0 {
				ShowWindow(hwnd, uintptr(SW_RESTORE))
				SetForegroundWindow(hwnd)
				return
			}
		}
	}

	ShellExecuteW(
		0,
		uintptr(unsafe.Pointer(utf16("open"))),
		uintptr(unsafe.Pointer(utf16(open.Target))),
		0,
		0,
		uintptr(SW_SHOW),
	)
}
