package config

type Config struct {
	Core      *Core
	Providers []*Provider
}

type Core struct {
	ServerListen        string
	SyncTitleName bool `yaml:"sync_title_name"`
	Output        string
}

type Provider struct {
	Uri          string
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
