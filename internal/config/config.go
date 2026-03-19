package config

import (
	"encoding/json"
	"fmt"
	"os"
	"wike/internal/shared"
)

const path = "config.json"

var modTime int64

type Config struct {
	Rules []Rule `json:"rules"`
}

func (c *Config) Load() error {
	data, err := os.ReadFile(path)
	if err != nil {
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

var Cfg = Config{
	Rules: []Rule{
		{
			Name:    "Useful CapsLock",
			Enabled: true,
			Trigger: &Trigger{Key: shared.Ptr("VK_CAPITAL")},
			Action:  &Action{Keys: []string{"VK_F13"}},
			Consume: shared.Ptr(true),
		},
	},
}

type Rule struct {
	Name     string     `json:"name"`
	Enabled  bool       `json:"enabled"`
	Region   *Region    `json:"region,omitempty"`
	Trigger  *Trigger   `json:"trigger,omitempty"`
	Action   *Action    `json:"action,omitempty"`
	Bindings *[]Binding `json:"bindings,omitempty"`
	Consume  *bool      `json:"consume,omitempty"`
}

type MouseButton string

const (
	LMB   MouseButton = "left"
	RMB   MouseButton = "right"
	MMB   MouseButton = "middle"
	X1B   MouseButton = "x1"
	X2B   MouseButton = "x2"
	Wheel MouseButton = "wheel"
	UMB   MouseButton = "unknown"
)

type TriggerEvent string

const (
	EventDown    TriggerEvent = "down"
	EventUp      TriggerEvent = "up"
	EventMove    TriggerEvent = "move"
	EventUnknown TriggerEvent = "unknown"
)

const (
	WM_LBUTTONDOWN = 0x0201
	WM_LBUTTONUP   = 0x0202
	WM_RBUTTONDOWN = 0x0204
	WM_RBUTTONUP   = 0x0205
	WM_MBUTTONDOWN = 0x0207
	WM_MBUTTONUP   = 0x0208
	WM_MOUSEMOVE   = 0x0200
	WM_MOUSEWHEEL  = 0x020A
	WM_XBUTTONDOWN = 0x020B
	WM_XBUTTONUP   = 0x020C
	XBUTTON1       = 0x10000
	XBUTTON2       = 0x20000
)

type MouseEvent struct {
	Button MouseButton
	Event  TriggerEvent
}

type Binding struct {
	Trigger *Trigger `json:"trigger"`
	Action  *Action  `json:"action"`
}

type Trigger struct {
	Mouse *MouseButton  `json:"mouse,omitempty"`
	Key   *string       `json:"key,omitempty"`
	Event *TriggerEvent `json:"event,omitempty"`
}

type Action struct {
	Keys   []string `json:"keys,omitempty"`
	Cmd    *string  `json:"cmd,omitempty"`
	Launch *string  `json:"launch,omitempty"`
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
		r.X += int32(shared.ScreenW)
	}
	if y < 0 {
		r.Y += int32(shared.ScreenH)
	}
	return r
}

func (r *Region) Contains(pt shared.POINT) bool {
	return pt.X >= r.X && pt.X < r.X+r.W &&
		pt.Y >= r.Y && pt.Y < r.Y+r.H
}

func ParseMouseEvent(wParam uintptr, mouseData uint32) MouseEvent {
	switch wParam {
	case WM_LBUTTONDOWN:
		return MouseEvent{LMB, EventDown}
	case WM_LBUTTONUP:
		return MouseEvent{LMB, EventUp}
	case WM_RBUTTONDOWN:
		return MouseEvent{RMB, EventDown}
	case WM_RBUTTONUP:
		return MouseEvent{RMB, EventUp}
	case WM_MBUTTONDOWN:
		return MouseEvent{MMB, EventDown}
	case WM_MBUTTONUP:
		return MouseEvent{MMB, EventUp}
	case WM_XBUTTONDOWN:
		if mouseData == XBUTTON1 {
			return MouseEvent{X1B, EventDown}
		}
		return MouseEvent{X2B, EventDown}
	case WM_XBUTTONUP:
		if mouseData == XBUTTON1 {
			return MouseEvent{X1B, EventUp}
		}
		return MouseEvent{X2B, EventUp}
	case WM_MOUSEMOVE:
		return MouseEvent{UMB, EventMove}
	case WM_MOUSEWHEEL:
		delta := int16(mouseData >> 16)
		if delta > 0 {
			return MouseEvent{Wheel, EventUp}
		}
		return MouseEvent{Wheel, EventDown}
	default:
		return MouseEvent{UMB, EventUnknown}
	}
}
