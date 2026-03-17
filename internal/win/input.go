package win

import "unsafe"

const (
	INPUT_KEYBOARD  = 1
	KEYEVENTF_KEYUP = 0x0002
)

var SendInput = user32.NewProc("SendInput").Call

type INPUT struct {
	Type uint32
	Ki   KEYBDINPUT
	_    [4]byte
}

func pressKeys(keys []uint16) {
	var inputs []INPUT

	for _, k := range keys {
		inputs = append(inputs, createInput(k, true))
	}

	for i := len(keys) - 1; i >= 0; i-- {
		inputs = append(inputs, createInput(keys[i], false))
	}

	SendInput(
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
