package config

import "encoding/json"

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
	Consume  bool      `json:"consume"`
}

type Trigger struct {
	M     *MouseButton `json:"m,omitempty"`
	Kb    *string      `json:"kb,omitempty"`
	State *State       `json:"state,omitempty"`
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

func (r *Rule) UnmarshalJSON(data []byte) error {
	type ruleAlias Rule

	*r = Rule{
		Name:    "Rule UNK",
		Enabled: true,
		Consume: true,
	}

	return json.Unmarshal(data, (*ruleAlias)(r))
}

func (r Rule) BindingsWithPrimary() []Binding {
	bindings := make([]Binding, 0, len(r.Bindings)+1)
	if r.Trigger != nil || r.Action != nil {
		bindings = append(bindings, Binding{
			Trigger: r.Trigger,
			Action:  r.Action,
		})
	}
	bindings = append(bindings, r.Bindings...)
	return bindings
}
