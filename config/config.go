package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type EpgProvider struct {
	Uri           string
	RemoveUnknown bool `yaml:"remove_unknown"`
	Pretty        bool
}

type Config struct {
	filepath  string
	Core      *Core
	Providers []*Provider
	Epg       []*EpgProvider
}

type Core struct {
	ServerListen     string `yaml:"server_listen"`
	AutoReloadConfig bool   `yaml:"auto_reload_config"`
	Output           string
	UpdateSchedule   string   `yaml:"update_schedule"`
	GroupOrder       []string `yaml:"group_order"`

	groupOrderMap map[string]int
}

type Provider struct {
	Uri               string
	IgnoreParseErrors bool `yaml:"ignore_parse_errors"`
	Filters           []string
	Setters           []*Setter
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

func New(filepath string) *Config {
	config = &Config{
		filepath: filepath,
		Core: &Core{
			AutoReloadConfig: true,
			UpdateSchedule:   "* */24 * * *",
			Output:           "m3u",
		},
	}

	config.Load()

	return config
}

func Get() *Config {
	return config
}

func (c *Config) Load() {
	yamlFile, err := ioutil.ReadFile(c.filepath)
	if err != nil {
		log.Fatalf("could not read config file %s, err = %v", c.filepath, err)
	}

	err = yaml.Unmarshal([]byte(yamlFile), &c)
	if err != nil {
		log.Fatalf("could not parse config file %s, err = %v", c.filepath, err)
	}
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
