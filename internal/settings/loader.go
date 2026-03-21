package settings

import (
	"encoding/json"
	"fmt"
	"os"
)

const configPath = "config.json"

var modTime int64

func (c *Settings) Load() error {
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			*c = defaultSettings()
			return c.Save()
		}
		return fmt.Errorf("read config: %w", err)
	}

	*c = Settings{}
	if err := json.Unmarshal(data, c); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	info, err := os.Stat(configPath)
	if err != nil {
		return fmt.Errorf("stat config: %w", err)
	}
	modTime = info.ModTime().Unix()
	fmt.Println("Config loaded:", c.toJSON())
	fmt.Printf("%+v\n", c)
	return nil
}

func (c *Settings) Save() error {
	err := os.WriteFile(configPath, []byte(c.toJSON()), 0644)
	if err == nil {
		info, statErr := os.Stat(configPath)
		if statErr != nil {
			return fmt.Errorf("stat config after save: %w", statErr)
		}
		modTime = info.ModTime().Unix()
		fmt.Println("Config saved:", c.toJSON())
		fmt.Printf("%+v\n", c)
		return nil
	}
	return fmt.Errorf("write config: %w", err)
}

func (c *Settings) ReloadIfModified() {
	info, err := os.Stat(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Println("Config stat error:", err)
		return
	}

	if info.ModTime().Unix() != modTime {
		fmt.Println("Config modified, reloading...")
		if err := c.Load(); err != nil {
			fmt.Println("Config reload error:", err)
		}
	}
}

func (c *Settings) toJSON() string {
	data, _ := json.MarshalIndent(c, "", "  ")
	return string(data)
}
