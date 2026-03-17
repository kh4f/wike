package config

import "wike/internal/shared"

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
		r.X += int32(shared.ScreenW)
	}
	if y < 0 {
		r.Y += int32(shared.ScreenH)
	}
	return r
}

func (r Region) Contains(pt shared.POINT) bool {
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

var Cfg = Config{
	Rules: []Rule{
		{
			Name:    "Useful CapsLock",
			Enabled: true,
			Trigger: Trigger{Key: shared.Ptr("VK_CAPITAL")},
			Action:  Action{Keys: []string{"VK_F13"}},
			Consume: shared.Ptr(true),
		},
	},
}
