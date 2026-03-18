package win

import (
	"fmt"
	"os/exec"
	"unsafe"
	"wike/internal/config"
	"wike/internal/shared"

	"golang.org/x/sys/windows"
)

const (
	LLKHF_UP   = 0x80
	SW_RESTORE = 9
	SW_SHOW    = 5
)

var (
	shell32             = windows.NewLazySystemDLL("shell32.dll")
	CallNextHookEx      = user32.NewProc("CallNextHookEx").Call
	GetCursorPos        = user32.NewProc("GetCursorPos").Call
	FindWindowW         = user32.NewProc("FindWindowW").Call
	IsWindow            = user32.NewProc("IsWindow").Call
	ShowWindow          = user32.NewProc("ShowWindow").Call
	SetForegroundWindow = user32.NewProc("SetForegroundWindow").Call
	ShellExecuteW       = shell32.NewProc("ShellExecuteW").Call
)

var vKCodeMap = map[string]uint16{
	"VK_F13":     0x7C,
	"VK_CAPITAL": 0x14,
}

var revVKCodeMap = func() map[uint16]string {
	m := make(map[uint16]string)
	for k, v := range vKCodeMap {
		m[v] = k
	}
	return m
}()

type MSLLHOOKSTRUCT struct {
	Pt          shared.POINT
	MouseData   uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

type KEYBDINPUT struct {
	WVk         uint16
	WScan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

func mHook(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 {
		info := (*MSLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		x, y := info.Pt.X, info.Pt.Y
		mouseEvent := config.ParseMouseEvent(wParam, info.MouseData)

		if mouseEvent.Event != config.EventMove {
			fmt.Printf("Mouse %s: %s (0x%X; %d:%d)\n", mouseEvent.Event, mouseEvent.Button, wParam, x, y)
		}

		for _, rule := range config.Cfg.Rules {
			if !rule.Enabled ||
				rule.Trigger.Mouse == nil ||
				*rule.Trigger.Mouse != mouseEvent.Button ||
				*rule.Trigger.Event != mouseEvent.Event ||
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
		keyID, found := revVKCodeMap[uint16(info.VkCode)]
		if !found {
			keyID = "unknown"
		}
		keyEvent := config.EventDown
		if (info.Flags & LLKHF_UP) != 0 {
			keyEvent = config.EventUp
		}

		fmt.Printf("Key %s: %s (0x%X; 0x%X; %d:%d)\n", keyEvent, keyID, wParam, info.VkCode, pt.X, pt.Y)

		for _, rule := range config.Cfg.Rules {
			if !rule.Enabled ||
				rule.Trigger.Key == nil ||
				*rule.Trigger.Key != keyID ||
				rule.Region != nil && !rule.Region.Contains(pt) ||
				rule.Trigger.Event == nil && keyEvent == config.EventUp ||
				rule.Trigger.Event != nil && *rule.Trigger.Event != keyEvent {
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
			if code, ok := vKCodeMap[keyName]; ok {
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
