package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const path = "config.json"

var modTime int64

type Config struct {
	Rules []Rule `json:"rules"`
}

func (c *Config) Load() error {
	data, err := os.ReadFile(path)
	if err != nil {
		*c = defCfg
		return c.Save()
	}
	*c = Config{}
	json.Unmarshal(data, c)

	info, _ := os.Stat(path)
	modTime = info.ModTime().Unix()
	fmt.Println("Config loaded:", c.toJSON())
	fmt.Printf("%+v\n", c)
	return nil
}

func (c *Config) Save() error {
	err := os.WriteFile(path, []byte(c.toJSON()), 0644)
	if err == nil {
		info, _ := os.Stat(path)
		modTime = info.ModTime().Unix()
		fmt.Println("Config saved:", c.toJSON())
		fmt.Printf("%+v\n", c)
	}
	return err
}

func (c *Config) ReloadIfModified() {
	info, _ := os.Stat(path)
	if info.ModTime().Unix() != modTime {
		fmt.Println("Config modified, reloading...")
		c.Load()
	}
}

func (c *Config) toJSON() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}

var Cfg Config

var defCfg = Config{
	Rules: []Rule{
		{
			Name:    "Default Rule",
			Enabled: true,
			Trigger: &Trigger{Kb: Ptr("VK_CAPITAL")},
			Action:  &Action{Kb: []string{"VK_F13"}},
			Consume: Ptr(true),
		},
	},
}

type Rule struct {
	Name     string    `json:"name"`
	Enabled  bool      `json:"enabled"`
	Region   *Region   `json:"region,omitempty"`
	Trigger  *Trigger  `json:"trigger,omitempty"`
	Action   *Action   `json:"action,omitempty"`
	Bindings []Binding `json:"bindings,omitempty"`
	Consume  *bool     `json:"consume,omitempty"`
}

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
	EventDown    State = "DOWN"
	EventUp      State = "UP"
	EventMove    State = "MOVE"
	EventUnknown State = "UNKNOWN"
)

type MouseEvent struct {
	Btn   string
	State State
}

func ParseMouseEvent(wParam uintptr, mouseData uint32) MouseEvent {
	switch wParam {
	case WM_LBUTTONDOWN:
		return MouseEvent{"LMB", EventDown}
	case WM_LBUTTONUP:
		return MouseEvent{"LMB", EventUp}
	case WM_RBUTTONDOWN:
		return MouseEvent{"RMB", EventDown}
	case WM_RBUTTONUP:
		return MouseEvent{"RMB", EventUp}
	case WM_MBUTTONDOWN:
		return MouseEvent{"MMB", EventDown}
	case WM_MBUTTONUP:
		return MouseEvent{"MMB", EventUp}
	case WM_XBUTTONDOWN:
		if mouseData == XBUTTON1 {
			return MouseEvent{"X1MB", EventDown}
		}
		return MouseEvent{"X2MB", EventDown}
	case WM_XBUTTONUP:
		if mouseData == XBUTTON1 {
			return MouseEvent{"X1MB", EventUp}
		}
		return MouseEvent{"X2MB", EventUp}
	case WM_MOUSEMOVE:
		return MouseEvent{"UNKNOWN", EventMove}
	case WM_MOUSEWHEEL:
		delta := int16(mouseData >> 16)
		if delta > 0 {
			return MouseEvent{"WHEEL", EventUp}
		}
		return MouseEvent{"WHEEL", EventDown}
	default:
		return MouseEvent{"UMB", EventUnknown}
	}
}

const (
	VK_CAPITAL     = 0x14
	VK_F13         = 0x7C
	VK_VOLUME_DOWN = 0xAE
	VK_VOLUME_UP   = 0xAF
)

var VKCodeMap = map[string]uint16{
	"VK_F13":         VK_F13,
	"VK_CAPITAL":     VK_CAPITAL,
	"VK_VOLUME_UP":   VK_VOLUME_UP,
	"VK_VOLUME_DOWN": VK_VOLUME_DOWN,
}

var RevVKCodeMap = func() map[uint16]string {
	m := make(map[uint16]string)
	for k, v := range VKCodeMap {
		m[v] = k
	}
	return m
}()

type Trigger struct {
	M     *string `json:"m,omitempty"`
	Kb    *string `json:"kb,omitempty"`
	State *State  `json:"state,omitempty"`
}

type KbEvent struct {
	Key   string
	Event State
}

const (
	LLKHF_UP = 0x80
)

type KBDLLHOOKSTRUCT struct {
	VkCode      uint32
	ScanCode    uint32
	Flags       uint32
	Time        uint32
	DwExtraInfo uintptr
}

func ParseKbEvent(info *KBDLLHOOKSTRUCT) KbEvent {
	keyID, found := RevVKCodeMap[uint16(info.VkCode)]
	if !found {
		keyID = "UNKNOWN"
	}
	kbEvent := EventDown
	if (info.Flags & LLKHF_UP) != 0 {
		kbEvent = EventUp
	}
	return KbEvent{Key: keyID, Event: kbEvent}
}

type Binding struct {
	Trigger *Trigger `json:"trigger"`
	Action  *Action  `json:"action"`
}

type Action struct {
	Kb     []string `json:"kb,omitempty"`
	Cmd    *string  `json:"cmd,omitempty"`
	Launch *string  `json:"launch,omitempty"`
}

type Region struct {
	X int32 `json:"x"`
	Y int32 `json:"y"`
	W int32 `json:"w"`
	H int32 `json:"h"`
}

func (r *Region) Contains(pt POINT) bool {
	rx, ry := r.X, r.Y
	if r.X < 0 {
		rx += int32(ScreenW)
	}
	if r.Y < 0 {
		ry += int32(ScreenH)
	}
	return pt.X >= rx && pt.X < rx+r.W &&
		pt.Y >= ry && pt.Y < ry+r.H
}

func Ptr[T any](v T) *T {
	return &v
}

type POINT struct {
	X int32
	Y int32
}

var (
	ScreenW int16
	ScreenH int16
)
