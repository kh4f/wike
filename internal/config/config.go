package config

var Cfg Config

var defCfg = Config{
	Rules: []Rule{
		{
			Name:    "Default Rule",
			Enabled: true,
			Trigger: &Trigger{Kb: ptr("VK_CAPITAL")},
			Action:  &Action{Kb: []string{"VK_F13"}},
			Consume: ptr(true),
		},
	},
}

func ptr[T any](v T) *T {
	return &v
}
