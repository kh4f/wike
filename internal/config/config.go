package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

const filePath = "config.yml"

var (
	Current = Config{
		Rules: []Rule{
			{
				Name:    "Caps Lock to F13",
				Enabled: true,
				Trigger: &Trigger{Kb: ptr("VK_CAPITAL")},
				Action:  &Action{Kb: []string{"VK_F13"}},
				Consume: true,
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
	if err := yaml.Unmarshal(data, &Current); err != nil {
		return fmt.Errorf("parse config: %w", err)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("stat config: %w", err)
	}
	modTime = info.ModTime()
	fmt.Printf("Config loaded: %s\n %+v", Current.toYAML(), Current)
	return nil
}

func Save() error {
	err := os.WriteFile(filePath, []byte(Current.toYAML()), 0644)
	if err == nil {
		info, statErr := os.Stat(filePath)
		if statErr != nil {
			return fmt.Errorf("stat config after save: %w", statErr)
		}
		modTime = info.ModTime()
		fmt.Printf("Config saved: %s\n %+v", Current.toYAML(), Current)
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

func (s *Config) toYAML() string {
	data, _ := yaml.Marshal(s)
	return string(data)
}

func ptr[T any](v T) *T {
	return &v
}
