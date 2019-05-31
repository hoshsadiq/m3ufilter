package config

type Config struct {
	Core      *Core
	Providers []*Provider
}

type Core struct {
	ServerListen   string `yaml:"server_listen"`
	Output         string
	UpdateSchedule string `yaml:"update_schedule"`
	GroupOrder     []string `yaml:"group_order"`

	groupOrderMap map[string]int
}

type Provider struct {
	Uri     string
	Filters []string
	Setters []*Setter
}

type Setter struct {
	ChNo    string `yaml:"chno"`
	Id      string
	Name    string
	Logo    string
	Group   string
	Shift   string
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

var config *Config

func Get() *Config {
	if config == nil {
		config = &Config{}
	}

	return config
}

func (c *Config) GetGroupOrder() map[string]int {
	if c.Core.groupOrderMap == nil {
		c.Core.groupOrderMap = map[string]int{}

		for order, groupTitle := range c.Core.GroupOrder {
			c.Core.groupOrderMap[groupTitle] = order
		}
	}

	return c.Core.groupOrderMap
}
