package config

type Config struct {
	Core      *Core
	Providers []*Provider
}

type Core struct {
	Listen        string
	SyncTitleName bool
}

type Provider struct {
	Uri          string
	Replacements *Replacement
	Filters      []string
	Setters      []*Setter
}

type Setter struct {
	Name       string
	Attributes map[string]string
	Filters    []string
}

type Replacement struct {
	Name       []*Replacer
	Attributes map[string][]*Replacer
}

type Replacer struct {
	Find    string
	Replace string
}
