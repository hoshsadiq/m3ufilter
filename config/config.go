package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

type EpgProvider struct {
	Uri           string
}

type Config struct {
	filepath     string
	Core         *Core
	Providers    []*Provider
	EpgProviders []*EpgProvider `yaml:"epg_providers"`
}

type Canonicalise struct {
	Enable         bool
	DefaultCountry string `yaml:"default_country"`
}

type Core struct {
	ServerListen     string `yaml:"server_listen"`
	AutoReloadConfig bool   `yaml:"auto_reload_config"`
	Output           string
	UpdateSchedule   string       `yaml:"update_schedule"`
	Canonicalise     Canonicalise `yaml:"canonicalise"`
	GroupOrder       []string     `yaml:"group_order"`

	groupOrderMap map[string]int
}

type Provider struct {
	Uri               string
	IgnoreParseErrors bool `yaml:"ignore_parse_errors"`
	CheckStreams      bool `yaml:"check_streams"`
	Filters           []string
	Setters           []*Setter
}

type Attributes struct {
	ChNo  string `yaml:"chno"`
	Id    string `yaml:"tvg-id"`
	Logo  string `yaml:"tvg-logo"`
	Group string `yaml:"group-title"`
	Shift string `yaml:"tvg-shift"`
}

type Setter struct {
	Name       string
	Attributes Attributes
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

var config *Config

func New(filepath string) *Config {
	config = &Config{
		filepath: filepath,
		Core: &Core{
			AutoReloadConfig: true,
			UpdateSchedule:   "* */24 * * *",
			Output:           "m3u",
			Canonicalise: Canonicalise{
				Enable:         true,
				DefaultCountry: "uk",
			},
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
