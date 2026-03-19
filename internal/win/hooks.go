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
	shell32        = windows.NewLazySystemDLL("shell32.dll")
	CallNextHookEx = user32.NewProc("CallNextHookEx").Call
	GetCursorPos   = user32.NewProc("GetCursorPos").Call
)

type MSLLHOOKSTRUCT struct {
	Pt          shared.POINT
	MouseData   uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

func mHook(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 {
		info := (*MSLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		x, y := info.Pt.X, info.Pt.Y
		mouseEvent := config.ParseMouseEvent(wParam, info.MouseData)

		if mouseEvent.State != config.EventMove {
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
					b.Trigger.State == nil && mouseEvent.State == config.EventDown {
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

	ret, _, _ := CallNextHookEx(0, uintptr(nCode), wParam, lParam)
	return ret
}

func kHook(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 {
		info := (*config.KBDLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		var pt shared.POINT
		GetCursorPos(uintptr(unsafe.Pointer(&pt)))
		kbEvent := config.ParseKbEvent(info)

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
					b.Trigger.State == nil && kbEvent.Event == config.EventDown {
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

	ret, _, _ := CallNextHookEx(0, uintptr(nCode), wParam, lParam)
	return ret
}

func executeAction(action config.Action) {
	fmt.Printf("Executing action: %+v\n", action)

	if len(action.Kb) > 0 {
		sendKeys(action.Kb, true, true)
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
