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
	"VK_F13":         0x7C,
	"VK_CAPITAL":     0x14,
	"VK_VOLUME_UP":   0xAF,
	"VK_VOLUME_DOWN": 0xAE,
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

func processMouseBinding(binding config.Binding, region *config.Region, ruleName string, mouseEvent *config.MouseEvent, pt *shared.POINT) bool {
	if binding.Trigger == nil || binding.Action == nil {
		return false
	}
	if binding.Trigger.Mouse == nil || *binding.Trigger.Mouse != mouseEvent.Button ||
		binding.Trigger.Event != nil && *binding.Trigger.Event != mouseEvent.Event {
		return false
	}
	if region != nil && !region.Contains(*pt) {
		return false
	}

	fmt.Printf("Rule triggered: '%s' (%d:%d)\n", ruleName, pt.X, pt.Y)
	executeAction(*binding.Action)
	return true
}

func processKeyBinding(binding config.Binding, region *config.Region, ruleName string, keyEvent config.TriggerEvent, keyID string, pt *shared.POINT) bool {
	if binding.Trigger == nil || binding.Action == nil {
		return false
	}
	if binding.Trigger.Key == nil || *binding.Trigger.Key != keyID ||
		binding.Trigger.Event != nil && *binding.Trigger.Event != keyEvent {
		return false
	}
	if region != nil && !region.Contains(*pt) {
		return false
	}

	fmt.Printf("Rule triggered: '%s' (%d:%d)\n", ruleName, pt.X, pt.Y)
	executeAction(*binding.Action)
	return true
}

func mHook(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 {
		info := (*MSLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		x, y := info.Pt.X, info.Pt.Y
		mouseEvent := config.ParseMouseEvent(wParam, info.MouseData)

		if mouseEvent.Event != config.EventMove {
			fmt.Printf("Mouse %s: %s (wParam=0x%X; mouseData=0x%X; pt=%d:%d)\n", mouseEvent.Event, mouseEvent.Button, wParam, info.MouseData, x, y)
		}

		for _, rule := range config.Cfg.Rules {
			if !rule.Enabled {
				continue
			}

			bindings := []config.Binding{}
			bindings = append(bindings, config.Binding{
				Trigger: rule.Trigger,
				Action:  rule.Action,
			})
			if rule.Bindings != nil {
				bindings = append(bindings, *rule.Bindings...)
			}

			isRuleMatched := false

			for _, b := range bindings {
				if processMouseBinding(b, rule.Region, rule.Name, &mouseEvent, &info.Pt) {
					isRuleMatched = true
				}
			}

			if isRuleMatched && rule.Consume != nil && *rule.Consume {
				fmt.Println("Event consumed")
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

		fmt.Printf("Key %s: %s (wParam=0x%X; VkCode=0x%X; pt=%d:%d)\n", keyEvent, keyID, wParam, info.VkCode, pt.X, pt.Y)

		for _, rule := range config.Cfg.Rules {
			if !rule.Enabled {
				continue
			}

			bindings := []config.Binding{}
			bindings = append(bindings, config.Binding{
				Trigger: rule.Trigger,
				Action:  rule.Action,
			})
			if rule.Bindings != nil {
				bindings = append(bindings, *rule.Bindings...)
			}

			isRuleMatched := false

			for _, b := range bindings {
				if processKeyBinding(b, rule.Region, rule.Name, keyEvent, keyID, &pt) {
					isRuleMatched = true
				}
			}

			if isRuleMatched && rule.Consume != nil && *rule.Consume {
				fmt.Println("Event consumed")
				return 1
			}
		}
	}

	ret, _, _ := CallNextHookEx(0, uintptr(nCode), wParam, lParam)
	return ret
}

func executeAction(action config.Action) {
	fmt.Printf("Executing action: %+v\n", action)

	if len(action.Keys) > 0 {
		var vkeys []uint16
		for _, keyName := range action.Keys {
			if code, ok := vKCodeMap[keyName]; ok {
				vkeys = append(vkeys, code)
			}
		}
		if len(vkeys) > 0 {
			fmt.Printf("Simulating key press: %v\n", action.Keys)
			pressKeys(vkeys)
		}
	}

	if action.Cmd != nil {
		fmt.Printf("Executing command: %s\n", *action.Cmd)
		exec.Command(*action.Cmd).Run()
	}

	if action.Launch != nil {
		fmt.Printf("Launching application: %s\n", *action.Launch)
		openOrFocus(*action.Launch)
	}
}
