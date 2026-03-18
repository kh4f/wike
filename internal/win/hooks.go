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
	LLKHF_UP = 0x80
)

var (
	shell32        = windows.NewLazySystemDLL("shell32.dll")
	CallNextHookEx = user32.NewProc("CallNextHookEx").Call
	GetCursorPos   = user32.NewProc("GetCursorPos").Call
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

	if action.Launch != nil {
		openOrFocus(*action.Launch)
	}
}
