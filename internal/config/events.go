package config

const (
	WM_LBUTTONDOWN = 0x0201
	WM_LBUTTONUP   = 0x0202
	WM_RBUTTONDOWN = 0x0204
	WM_RBUTTONUP   = 0x0205
	WM_MBUTTONDOWN = 0x0207
	WM_MBUTTONUP   = 0x0208
	WM_XBUTTONDOWN = 0x020B
	WM_XBUTTONUP   = 0x020C
	WM_MOUSEWHEEL  = 0x020A
	WM_MOUSEMOVE   = 0x0200
	XBUTTON1       = 0x10000
	XBUTTON2       = 0x20000
)

type State string

const (
	StateDown    State = "DOWN"
	StateUp      State = "UP"
	StateMove    State = "MOVE"
	StateUnknown State = "UNKNOWN"
)

type MouseEvent struct {
	Btn   string
	State State
}

func ParseMouseEvent(wParam uintptr, mouseData uint32) MouseEvent {
	switch wParam {
	case WM_LBUTTONDOWN:
		return MouseEvent{"LMB", StateDown}
	case WM_LBUTTONUP:
		return MouseEvent{"LMB", StateUp}
	case WM_RBUTTONDOWN:
		return MouseEvent{"RMB", StateDown}
	case WM_RBUTTONUP:
		return MouseEvent{"RMB", StateUp}
	case WM_MBUTTONDOWN:
		return MouseEvent{"MMB", StateDown}
	case WM_MBUTTONUP:
		return MouseEvent{"MMB", StateUp}
	case WM_XBUTTONDOWN:
		if mouseData == XBUTTON1 {
			return MouseEvent{"X1MB", StateDown}
		}
		return MouseEvent{"X2MB", StateDown}
	case WM_XBUTTONUP:
		if mouseData == XBUTTON1 {
			return MouseEvent{"X1MB", StateUp}
		}
		return MouseEvent{"X2MB", StateUp}
	case WM_MOUSEMOVE:
		return MouseEvent{"UNKNOWN", StateMove}
	case WM_MOUSEWHEEL:
		delta := int16(mouseData >> 16)
		if delta > 0 {
			return MouseEvent{"WHEEL", StateUp}
		}
		return MouseEvent{"WHEEL", StateDown}
	default:
		return MouseEvent{"UMB", StateUnknown}
	}
}

type KbEvent struct {
	Key   string
	Event State
}

const LLKHF_UP = 0x80

func ParseKbEvent(vkCode uint16, flags uint32) KbEvent {
	keyID, found := RevVKCodeMap[vkCode]
	if !found {
		keyID = "UNKNOWN"
	}

	kbEvent := StateDown
	if (flags & LLKHF_UP) != 0 {
		kbEvent = StateUp
	}

	return KbEvent{Key: keyID, Event: kbEvent}
}
