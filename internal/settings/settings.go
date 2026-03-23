package settings

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
)

const filePath = "settings.json"

var (
	Current = Settings{
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
		return fmt.Errorf("read settings: %w", err)
	}

	Current = Settings{}
	if err := json.Unmarshal(data, &Current); err != nil {
		return fmt.Errorf("parse settings: %w", err)
	}

	info, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("stat settings: %w", err)
	}
	modTime = info.ModTime()
	fmt.Printf("Settings loaded: %s\n %+v", Current.toJSON(), Current)
	return nil
}

func Save() error {
	err := os.WriteFile(filePath, []byte(Current.toJSON()), 0644)
	if err == nil {
		info, statErr := os.Stat(filePath)
		if statErr != nil {
			return fmt.Errorf("stat settings after save: %w", statErr)
		}
		modTime = info.ModTime()
		fmt.Printf("Settings saved: %s\n %+v", Current.toJSON(), Current)
		return nil
	}
	return fmt.Errorf("write settings: %w", err)
}

func ReloadIfModified() {
	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return
		}
		fmt.Println("Settings stat error:", err)
		return
	}

	if !info.ModTime().Equal(modTime) {
		fmt.Println("Settings modified, reloading...")
		if err := Load(); err != nil {
			fmt.Println("Settings reload error:", err)
		}
	}
}

func (s *Settings) toJSON() string {
	data, _ := json.MarshalIndent(s, "", "  ")
	return string(data)
}

func ptr[T any](v T) *T {
	return &v
}
