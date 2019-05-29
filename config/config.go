package config

type Config struct {
	Core      *Core
	Providers []*Provider
}

type Core struct {
	ServerListen   string `yaml:"server_listen"`
	SyncTitleName  bool   `yaml:"sync_title_name"`
	Output         string
	UpdateSchedule string `yaml:"update_schedule"`
}

type Provider struct {
	Uri     string
	Filters []string
	Setters []*Setter
}

type Setter struct {
	Name    string
	Id      string
	Logo    string
	Group   string
	Filters []string
}

type Replacement struct {
	Name       []*Replacer
	Attributes map[string][]*Replacer
}

type Replacer struct {
	Find    string
	Replace string
}
