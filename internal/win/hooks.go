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
	if nCode < 0 {
		return callNextHook(nCode, wParam, lParam)
	}

	info := (*MSLLHOOKSTRUCT)(unsafe.Pointer(lParam))
	if (info.Flags & LLKHF_INJECTED) != 0 {
		return callNextHook(nCode, wParam, lParam)
	}

	mouseEvent := config.ParseMouseEvent(wParam, info.MouseData)
	if mouseEvent.State != config.StateMove {
		fmt.Printf("Mouse event: %+v (wParam=0x%X; mouseData=0x%X; pt=%d:%d)\n", mouseEvent, wParam, info.MouseData, info.Pt.X, info.Pt.Y)
	}

	for _, rule := range config.Current.Rules {
		if !rule.Enabled || (rule.Region != nil && !rule.Region.Contains(info.Pt)) {
			continue
		}

		matched := false
		for _, binding := range rule.BindingsWithPrimary() {
			if binding.Trigger == nil ||
				binding.Action == nil ||
				binding.Trigger.M == nil ||
				*binding.Trigger.M != mouseEvent.Btn {
				continue
			}

			matched = true
			if shouldExecuteBinding(binding.Trigger.State, mouseEvent.State) {
				fmt.Printf("Rule triggered: '%s'\n", rule.Name)
				executeAction(*binding.Action)
			}
		}

		if matched && rule.Consume {
			fmt.Println("Event consumed")
			return 1
		}
	}

	return callNextHook(nCode, wParam, lParam)
}

func kHook(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode < 0 {
		return callNextHook(nCode, wParam, lParam)
	}

	info := (*KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
	if (info.Flags & LLKHF_INJECTED) != 0 {
		return callNextHook(nCode, wParam, lParam)
	}

	var pt shared.Point
	getCursorPos(uintptr(unsafe.Pointer(&pt)))
	kbEvent := config.ParseKbEvent(uint16(info.VkCode), info.Flags)

	fmt.Printf("Kb event: %+v (wParam=0x%X; vkCode=0x%X; pt=%d:%d)\n", kbEvent, wParam, info.VkCode, pt.X, pt.Y)

	for _, rule := range config.Current.Rules {
		if !rule.Enabled || (rule.Region != nil && !rule.Region.Contains(pt)) {
			continue
		}

		matched := false
		for _, binding := range rule.BindingsWithPrimary() {
			if binding.Trigger == nil ||
				binding.Action == nil ||
				binding.Trigger.Kb == nil ||
				*binding.Trigger.Kb != kbEvent.Key {
				continue
			}

			matched = true
			if shouldExecuteBinding(binding.Trigger.State, kbEvent.Event) {
				fmt.Printf("Rule triggered: '%s'\n", rule.Name)
				executeAction(*binding.Action)
			}
		}

		if matched && rule.Consume {
			fmt.Println("Event consumed")
			return 1
		}
	}

	return callNextHook(nCode, wParam, lParam)
}

func shouldExecuteBinding(triggerState *config.State, eventState config.State) bool {
	if triggerState == nil {
		return eventState == config.StateDown
	}
	return *triggerState == eventState
}

func callNextHook(nCode int, wParam uintptr, lParam uintptr) uintptr {
	ret, _, _ := callNextHookEx(0, uintptr(nCode), wParam, lParam)
	return ret
}
