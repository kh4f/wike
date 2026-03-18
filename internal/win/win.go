package win

import (
	"time"
	"wike/internal/config"
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
	user32            = windows.NewLazySystemDLL("user32.dll")
	GetSystemMetrics  = user32.NewProc("GetSystemMetrics").Call
	SetWindowsHookExW = user32.NewProc("SetWindowsHookExW").Call
	GetMessageW       = user32.NewProc("GetMessageW").Call
)

func RunMessageLoop() {
	initScreenSize()
	config.Cfg.Load()

	SetWindowsHookExW(uintptr(WH_MOUSE_LL), windows.NewCallback(mHook), 0, 0)
	SetWindowsHookExW(uintptr(WH_KEYBOARD_LL), windows.NewCallback(kHook), 0, 0)

	go func() {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for range ticker.C {
			config.Cfg.ReloadIfModified()
		}
	}()

	GetMessageW(0, 0, 0, 0)
}

func initScreenSize() {
	w, _, _ := GetSystemMetrics(uintptr(SM_CXSCREEN))
	h, _, _ := GetSystemMetrics(uintptr(SM_CYSCREEN))
	shared.ScreenW, shared.ScreenH = int16(w), int16(h)
}

func utf16(s string) *uint16 {
	ptr, err := windows.UTF16PtrFromString(s)
	if err != nil {
		panic(err)
	}
	return ptr
}
