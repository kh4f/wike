package win

import (
	"fmt"
	"unsafe"
	"wike/internal/config"
	"wike/internal/shared"
)

var (
	callNextHookEx = user32.NewProc("CallNextHookEx").Call
	getCursorPos   = user32.NewProc("GetCursorPos").Call
)

const LLKHF_INJECTED = 0x10

type MSLLHOOKSTRUCT struct {
	Pt          shared.Point
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

func mHook(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 {
		info := (*MSLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		if (info.Flags & LLKHF_INJECTED) != 0 {
			return callNextHook(nCode, wParam, lParam)
		}

		x, y := info.Pt.X, info.Pt.Y
		mouseEvent := config.ParseMouseEvent(wParam, info.MouseData)

		if mouseEvent.State != config.StateMove {
			fmt.Printf("Mouse event: %+v (wParam=0x%X; mouseData=0x%X; pt=%d:%d)\n", mouseEvent, wParam, info.MouseData, x, y)
		}

		for _, rule := range config.Cfg.Rules {
			if !rule.Enabled || rule.Region != nil && !rule.Region.Contains(info.Pt) {
				continue
			}

			bindings := []config.Binding{}
			bindings = append(bindings, config.Binding{
				Trigger: rule.Trigger,
				Action:  rule.Action,
			})
			bindings = append(bindings, rule.Bindings...)

			isRuleMatched := false

			for _, b := range bindings {
				if b.Trigger == nil || b.Action == nil || b.Trigger.M == nil || *b.Trigger.M != mouseEvent.Btn {
					continue
				}

				isRuleMatched = true

				if b.Trigger.State != nil && *b.Trigger.State == mouseEvent.State ||
					b.Trigger.State == nil && mouseEvent.State == config.StateDown {
					fmt.Printf("Rule triggered: '%s'\n", rule.Name)
					executeAction(*b.Action)
				}
			}

			if isRuleMatched && rule.Consume != nil && *rule.Consume {
				fmt.Println("Event consumed")
				return 1
			}
		}
	}

	return callNextHook(nCode, wParam, lParam)
}

func kHook(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 {
		info := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		if (info.Flags & LLKHF_INJECTED) != 0 {
			return callNextHook(nCode, wParam, lParam)
		}

		var pt shared.Point
		getCursorPos(uintptr(unsafe.Pointer(&pt)))
		kbEvent := config.ParseKbEvent(uint16(info.VkCode), info.Flags)

		fmt.Printf("Key event: %+v (wParam=0x%X; VkCode=0x%X; pt=%d:%d)\n", kbEvent, wParam, info.VkCode, pt.X, pt.Y)

		for _, rule := range config.Cfg.Rules {
			if !rule.Enabled {
				continue
			}

			bindings := []config.Binding{}
			bindings = append(bindings, config.Binding{
				Trigger: rule.Trigger,
				Action:  rule.Action,
			})
			bindings = append(bindings, rule.Bindings...)

			shouldConsume := false

			for _, b := range bindings {
				if b.Trigger == nil || b.Action == nil || b.Trigger.Kb == nil || *b.Trigger.Kb != kbEvent.Key {
					continue
				}

				shouldConsume = true

				if b.Trigger.State != nil && *b.Trigger.State == kbEvent.Event ||
					b.Trigger.State == nil && kbEvent.Event == config.StateDown {
					fmt.Printf("Rule triggered: '%s'\n", rule.Name)
					executeAction(*b.Action)
				}
			}

			if shouldConsume && rule.Consume != nil && *rule.Consume {
				fmt.Println("Event consumed")
				return 1
			}
		}
	}

	return callNextHook(nCode, wParam, lParam)
}

func callNextHook(nCode int, wParam uintptr, lParam uintptr) uintptr {
	ret, _, _ := callNextHookEx(0, uintptr(nCode), wParam, lParam)
	return ret
}
