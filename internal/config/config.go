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
	return nil
}

func (c *Config) Save() error {
	err := os.WriteFile(path, []byte(c.toJSON()), 0644)
	if err == nil {
		info, _ := os.Stat(path)
		modTime = info.ModTime().Unix()
		fmt.Println("Config saved:", c.toJSON())
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
			Trigger: Trigger{Key: shared.Ptr("VK_CAPITAL")},
			Action:  Action{Keys: []string{"VK_F13"}},
			Consume: shared.Ptr(true),
		},
	},
}

type Rule struct {
	Name    string  `json:"name"`
	Enabled bool    `json:"enabled"`
	Region  *Region `json:"region,omitempty"`
	Trigger Trigger `json:"trigger"`
	Action  Action  `json:"action"`
	Consume *bool   `json:"consume,omitempty"`
}

type Trigger struct {
	Mouse *string `json:"mouse,omitempty"`
	Key   *string `json:"key,omitempty"`
	Event *string `json:"event,omitempty"`
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
