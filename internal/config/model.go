package config

import "gopkg.in/yaml.v3"

type Config struct {
	Rules []Rule `yaml:"rules"`
}

type Rule struct {
	Name     string    `yaml:"name"`
	Enabled  bool      `yaml:"enabled"`
	Region   *Region   `yaml:"region,omitempty"`
	Trigger  *Trigger  `yaml:"trigger,omitempty"`
	Action   *Action   `yaml:"action,omitempty"`
	Bindings []Binding `yaml:"bindings,omitempty"`
	Consume  bool      `yaml:"consume"`
}

type Trigger struct {
	M     *MouseButton `yaml:"m,omitempty"`
	Kb    *string      `yaml:"kb,omitempty"`
	State *State       `yaml:"state,omitempty"`
}

type Binding struct {
	Trigger *Trigger `yaml:"trigger"`
	Action  *Action  `yaml:"action"`
}

type Action struct {
	Kb     []string `yaml:"kb,omitempty"`
	Cmd    *string  `yaml:"cmd,omitempty"`
	Launch *string  `yaml:"launch,omitempty"`
}

func (r *Rule) UnmarshalYAML(value *yaml.Node) error {
	type ruleAlias Rule

	*r = Rule{
		Name:    "Rule UNK",
		Enabled: true,
		Consume: true,
	}

	return value.Decode((*ruleAlias)(r))
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
