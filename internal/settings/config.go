package settings

var Current Settings

var defSettings = Settings{
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

func defaultSettings() Settings {
	cfg := Settings{
		Rules: make([]Rule, len(defSettings.Rules)),
	}
	copy(cfg.Rules, defSettings.Rules)
	return cfg
}

func ptr[T any](v T) *T {
	return &v
}
