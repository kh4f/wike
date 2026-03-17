package main

import (
	"fmt"
	"os/exec"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	user32              = windows.NewLazySystemDLL("user32.dll")
	kernel32            = windows.NewLazySystemDLL("kernel32.dll")
	shell32             = windows.NewLazySystemDLL("shell32.dll")
	SetWindowsHookExW   = user32.NewProc("SetWindowsHookExW").Call
	CallNextHookEx      = user32.NewProc("CallNextHookEx").Call
	GetMessageW         = user32.NewProc("GetMessageW").Call
	FindWindowW         = user32.NewProc("FindWindowW").Call
	IsWindow            = user32.NewProc("IsWindow").Call
	ShowWindow          = user32.NewProc("ShowWindow").Call
	SetForegroundWindow = user32.NewProc("SetForegroundWindow").Call
	ShellExecuteW       = shell32.NewProc("ShellExecuteW").Call
	GetCursorPos        = user32.NewProc("GetCursorPos").Call
)

const (
	WH_KEYBOARD_LL = 13
	WH_MOUSE_LL    = 14
	SW_RESTORE     = 9
	SW_SHOW        = 5
)

func utf16(s string) *uint16 {
	ptr, err := windows.UTF16PtrFromString(s)
	if err != nil {
		panic(err)
	}
	return ptr
}

func ptr[T any](v T) *T {
	return &v
}

var GetSystemMetrics = user32.NewProc("GetSystemMetrics").Call

var ScreenW, ScreenH int16

const (
	SM_CXSCREEN = 0
	SM_CYSCREEN = 1
)

func init() {
	w, _, _ := GetSystemMetrics(uintptr(SM_CXSCREEN))
	h, _, _ := GetSystemMetrics(uintptr(SM_CYSCREEN))
	ScreenW = int16(w)
	ScreenH = int16(h)
}

type POINT struct {
	X int32
	Y int32
}

type MSLLHOOKSTRUCT struct {
	Pt          POINT
	MouseData   uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

const (
	INPUT_KEYBOARD  = 1
	KEYEVENTF_KEYUP = 0x0002
)

type KEYBDINPUT struct {
	WVk         uint16
	WScan       uint16
	DwFlags     uint32
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

type INPUT struct {
	Type uint32
	Ki   KEYBDINPUT
	_    [4]byte
}

const (
	WM_MOUSEMOVE   = 0x0200
	WM_LBUTTONDOWN = 0x0201
	WM_LBUTTONUP   = 0x0202
	WM_RBUTTONDOWN = 0x0204
	WM_RBUTTONUP   = 0x0205
	WM_MBUTTONDOWN = 0x0207
	WM_MBUTTONUP   = 0x0208
)

var sendInput = user32.NewProc("SendInput").Call

type Config struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Name    string  `json:"name"`
	Enabled bool    `json:"enabled"`
	Region  *Region `json:"region,omitempty"`
	Trigger Trigger `json:"trigger"`
	Action  Action  `json:"action"`
	Consume *bool   `json:"consume,omitempty"`
}

type Region struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
	W int32 `json:"w"`
	H int32 `json:"h"`
}

func NewRegion(x, y, w, h int32) *Region {
	r := &Region{X: x, Y: y, W: w, H: h}
	if x < 0 {
		r.X += int32(ScreenW)
	}
	if y < 0 {
		r.Y += int32(ScreenH)
	}
	return r
}

func (r Region) Contains(pt POINT) bool {
	return pt.X >= r.X && pt.X < r.X+r.W &&
		pt.Y >= r.Y && pt.Y < r.Y+r.H
}

type Trigger struct {
	Key   *string `json:"key,omitempty"`
	Mouse *string `json:"mouse,omitempty"`
}

type Action struct {
	Keys []string `json:"keys,omitempty"`
	Cmd  *string  `json:"cmd,omitempty"`
	Open *Open    `json:"open,omitempty"`
}

type Open struct {
	Target      string  `json:"target"`
	WindowClass *string `json:"windowClass,omitempty"`
}

var mouseEventMap = map[uintptr]string{
	WM_LBUTTONDOWN: "LeftDown",
	WM_LBUTTONUP:   "LeftUp",
	WM_RBUTTONDOWN: "RightDown",
	WM_RBUTTONUP:   "RightUp",
	WM_MBUTTONDOWN: "MiddleDown",
	WM_MBUTTONUP:   "MiddleUp",
	WM_MOUSEMOVE:   "Move",
}

var virtualKeyMap = map[string]uint16{
	"VK_F13":     0x7C,
	"VK_CAPITAL": 0x14,
}

func createInput(vKey uint16, keyDown bool) INPUT {
	flags := uint32(0)
	if !keyDown {
		flags = KEYEVENTF_KEYUP
	}

	return INPUT{
		Type: INPUT_KEYBOARD,
		Ki: KEYBDINPUT{
			WVk:         vKey,
			WScan:       0,
			DwFlags:     flags,
			Time:        0,
			DwExtraInfo: 0,
		},
	}
}

func pressKeys(keys []uint16) {
	var inputs []INPUT

	for _, k := range keys {
		inputs = append(inputs, createInput(k, true))
	}

	for i := len(keys) - 1; i >= 0; i-- {
		inputs = append(inputs, createInput(keys[i], false))
	}

	sendInput(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		uintptr(unsafe.Sizeof(inputs[0])),
	)
}

func openOrFocus(open Open) {
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

func executeAction(action Action) {
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

func mHook(nCode int, wParam uintptr, lParam uintptr) uintptr {
	if nCode >= 0 {
		info := (*MSLLHOOKSTRUCT)(unsafe.Pointer(lParam))
		x, y := info.Pt.X, info.Pt.Y
		mEvent := mouseEventMap[wParam]
		if mEvent != "Move" {
			fmt.Printf("Mouse event: wParam=0x%X (%s) at (%d, %d)\n", wParam, mEvent, x, y)
		}

		for _, rule := range cfg.Rules {
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
		var pt POINT
		GetCursorPos(uintptr(unsafe.Pointer(&pt)))
		fmt.Printf("Key event: wParam=0x%X vkCode=0x%X at (%d, %d)\n", wParam, info.VkCode, pt.X, pt.Y)

		for _, rule := range cfg.Rules {
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

func main() {
	fmt.Printf("Screen size: %d x %d\n", ScreenW, ScreenH)
	SetWindowsHookExW(uintptr(WH_MOUSE_LL), windows.NewCallback(mHook), 0, 0)
	SetWindowsHookExW(uintptr(WH_KEYBOARD_LL), windows.NewCallback(kHook), 0, 0)

	for {
		GetMessageW(0, 0, 0, 0)
	}
}

var cfg = Config{
	Rules: []Rule{
		{
			Name:    "Useful CapsLock",
			Enabled: true,
			Trigger: Trigger{Key: ptr("VK_CAPITAL")},
			Action:  Action{Keys: []string{"VK_F13"}},
			Consume: ptr(true),
		},
	},
}
