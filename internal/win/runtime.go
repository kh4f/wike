package win

import (
	"wike/internal/shared"

	"golang.org/x/sys/windows"
)

const (
	WH_KEYBOARD_LL = 13
	WH_MOUSE_LL    = 14
	SM_CXSCREEN    = 0
	SM_CYSCREEN    = 1
)

var (
	getSystemMetrics  = user32.NewProc("GetSystemMetrics").Call
	setWindowsHookExW = user32.NewProc("SetWindowsHookExW").Call
	getMessageW       = user32.NewProc("GetMessageW").Call
)

func InstallHooks() {
	setWindowsHookExW(uintptr(WH_MOUSE_LL), windows.NewCallback(mHook), 0, 0)
	setWindowsHookExW(uintptr(WH_KEYBOARD_LL), windows.NewCallback(kHook), 0, 0)
}

func RunMessageLoop() {
	getMessageW(0, 0, 0, 0)
}

func InitScreenSize() {
	w, _, _ := getSystemMetrics(uintptr(SM_CXSCREEN))
	h, _, _ := getSystemMetrics(uintptr(SM_CYSCREEN))
	shared.ScreenWidth, shared.ScreenHeight = int16(w), int16(h)
}

func utf16(s string) *uint16 {
	ptr, err := windows.UTF16PtrFromString(s)
	if err != nil {
		panic(err)
	}
	return ptr
}
