package config

var Cfg Config

var defCfg = Config{
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

func defaultConfig() Config {
	cfg := Config{
		Rules: make([]Rule, len(defCfg.Rules)),
	}
	copy(cfg.Rules, defCfg.Rules)
	return cfg
}

func ptr[T any](v T) *T {
	return &v
}
