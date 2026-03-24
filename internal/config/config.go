package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const filePath = "config.json"

var (
	Current = Config{
		Rules: []Rule{
			{
				Name:    "Caps Lock to F13",
				Enabled: true,
				Trigger: &Trigger{Kb: ptr("VK_CAPITAL")},
				Action:  &Action{Kb: []string{"VK_F13"}},
				Consume: ptr(true),
			},
		},
	}
	modTime time.Time
)

func Load() error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return Save()
		}
		return fmt.Errorf("read config: %w", err)
	}

	Current = Config{}
	if err := json.Unmarshal(data, &Current); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("stat config: %w", err)
	}
	modTime = info.ModTime()
	fmt.Printf("Config loaded: %s\n %+v", Current.toJSON(), Current)
	return nil
}

func Save() error {
	err := os.WriteFile(filePath, []byte(Current.toJSON()), 0644)
	if err == nil {
		info, statErr := os.Stat(filePath)
		if statErr != nil {
			return fmt.Errorf("stat config after save: %w", statErr)
		}
		modTime = info.ModTime()
		fmt.Printf("Config saved: %s\n %+v", Current.toJSON(), Current)
		return nil
	}
	return fmt.Errorf("write config: %w", err)
}

func ReloadIfModified() {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Println("Config stat error:", err)
		return
	}

	if !info.ModTime().Equal(modTime) {
		fmt.Println("Config modified, reloading...")
		if err := Load(); err != nil {
			fmt.Println("Config reload error:", err)
		}
	}
}

func (s *Config) toJSON() string {
	data, _ := json.MarshalIndent(s, "", "  ")
	return string(data)
}

func ptr[T any](v T) *T {
	return &v
}
