package config

import (
	"encoding/json"
	"fmt"
	"os"
)

const path = "config.json"

var modTime int64

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
