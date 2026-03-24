package win

import (
	"fmt"
	"unsafe"
	"wike/internal/config"
)

const (
	INPUT_KEYBOARD  = 1
	KEYEVENTF_KEYUP = 0x0002
)

var sendInput = user32.NewProc("SendInput").Call

type INPUT struct {
	Type uint32
	Ki   KEYBDINPUT
	_    [4]byte
}

type KEYBDINPUT struct {
	WVk         uint16
	WScan       uint16
	DwFlags     uint32
	Time        uint32
	DwExtraInfo uintptr
}

func sendKeys(keys []string, press bool, release bool) {
	var vkeys []uint16
	for _, keyName := range keys {
		if code, ok := config.VKCodeMap[keyName]; ok {
			vkeys = append(vkeys, code)
		}
	}
	if len(vkeys) == 0 {
		return
	}

	var inputs []INPUT

	if press {
		fmt.Printf("Simulating key press: %v\n", keys)
		for _, k := range vkeys {
			inputs = append(inputs, createInput(k, true))
		}
	}

	if release {
		fmt.Printf("Simulating key release: %v\n", keys)
		for i := len(vkeys) - 1; i >= 0; i-- {
			inputs = append(inputs, createInput(vkeys[i], false))
		}
	}

	sendInput(
		uintptr(len(inputs)),
		uintptr(unsafe.Pointer(&inputs[0])),
		uintptr(unsafe.Sizeof(inputs[0])),
	)
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
