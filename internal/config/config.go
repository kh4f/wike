package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
	"wike/internal/logger"

	"gopkg.in/yaml.v3"
)

const fileName = "config.yml"

var (
	Current = Config{
		Rules: []Rule{
			{
				Name:    "Caps Lock → F13",
				Trigger: &Trigger{Kb: ptr("CAPITAL")},
				Action:  &Action{Kb: []string{"F13"}},
			},
		},
	}
	modTime time.Time
)

func Load() error {
	filePath := configFilePath()
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
	logger.Printf("Config loaded from %s:\n%s%+v\n", filePath, Current.toYAML(), Current)
	return nil
}

func Save() error {
	filePath := configFilePath()
	err := os.WriteFile(filePath, []byte(Current.toYAML()), 0644)
	if err == nil {
		info, statErr := os.Stat(filePath)
		if statErr != nil {
			return fmt.Errorf("stat config after save: %w", statErr)
		}
		modTime = info.ModTime()
		logger.Printf("Config saved to %s:\n%s%+v\n", filePath, Current.toYAML(), Current)
		return nil
	}
	return fmt.Errorf("write config: %w", err)
}

func ReloadIfModified() {
	filePath := configFilePath()
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		logger.Println("Config stat error:", err)
		return
	}

	if !info.ModTime().Equal(modTime) {
		logger.Println("Config modified, reloading...")
		if err := Load(); err != nil {
			logger.Println("Config reload error:", err)
		}
	}
}

func (s *Config) toYAML() string {
	data, _ := yaml.Marshal(s)
	return string(data)
}

func configFilePath() string {
	exePath, err := os.Executable()
	if err != nil {
		return fileName
	}
	return filepath.Join(filepath.Dir(exePath), fileName)
}

func ptr[T any](v T) *T {
	return &v
}
