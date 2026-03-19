package config

type Config struct {
	Rules []Rule `json:"rules"`
}

type Rule struct {
	Name     string    `json:"name"`
	Enabled  bool      `json:"enabled"`
	Region   *Region   `json:"region,omitempty"`
	Trigger  *Trigger  `json:"trigger,omitempty"`
	Action   *Action   `json:"action,omitempty"`
	Bindings []Binding `json:"bindings,omitempty"`
	Consume  *bool     `json:"consume,omitempty"`
}

type Trigger struct {
	M     *string `json:"m,omitempty"`
	Kb    *string `json:"kb,omitempty"`
	State *State  `json:"state,omitempty"`
}

type Binding struct {
	Trigger *Trigger `json:"trigger"`
	Action  *Action  `json:"action"`
}

type Action struct {
	Kb     []string `json:"kb,omitempty"`
	Cmd    *string  `json:"cmd,omitempty"`
	Launch *string  `json:"launch,omitempty"`
}
