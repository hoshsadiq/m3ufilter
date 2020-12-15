package config

import (
	"github.com/hoshsadiq/m3ufilter/logger"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var log = logger.Get()

type ChannelIdRename struct {
	From string
	To   string
}

type EpgProvider struct {
	Uri              string
	ChannelIdRenames []ChannelIdRename `yaml:"channel_id_renames"`
}

type Config struct {
	filepath     string
	Core         *Core
	Providers    []*Provider
	EpgProviders []*EpgProvider `yaml:"epg_providers"`
}

type Core struct {
	ServerListen         string   `yaml:"server_listen"`
	AutoReloadConfig     bool     `yaml:"auto_reload_config"`
	UpdateSchedule       string   `yaml:"update_schedule"`
	GroupOrder           []string `yaml:"group_order"`
	PrettyOutputXml      bool
	HttpTimeout          uint8 `yaml:"http_timeout"` // in seconds
	HttpMaxRetryAttempts int   `yaml:"http_max_retry_attempts"`

	groupOrderMap map[string]int
}

type CheckStreamsAction string

const (
	InvalidStreamRemove CheckStreamsAction = "remove"
	InvalidStreamNoop                      = "noop"
)

type CheckStreams struct {
	Enabled bool
	Method  string
	Action  CheckStreamsAction
}

func (c *CheckStreams) UnmarshalYAML(unmarshal func(interface{}) error) (err error) {
	// first we try the old mechanism for backwards compatibility
	err = unmarshal(&c.Enabled)
	if err == nil {
		log.Warnf("using a boolean value for provider.check_streams is deprecated, this will be removed in the future. Please upgrade to new method. See the docs for information.")
		c.Method = "head"
		c.Action = "remove"
		return
	}

	cs := struct {
		Enabled bool
		Method  string
		Action  CheckStreamsAction
	}{
		Enabled: false,
		Action:  "remove",
		Method:  "head",
	}

	err = unmarshal(&cs)
	if err != nil {
		return err
	}

	c.Enabled = cs.Enabled
	c.Action = cs.Action
	c.Method = cs.Method

	return
}

type Provider struct {
	Uri          string
	Csv          string
	CheckStreams CheckStreams `yaml:"check_streams"`
}

var config *Config

func New(filepath string) (*Config, error) {
	config = &Config{
		filepath: filepath,
		Core: &Core{
			AutoReloadConfig:     true,
			UpdateSchedule:       "0 */24 * * *",
			HttpTimeout:          60,
			HttpMaxRetryAttempts: 5,
		},
	}

	err := config.Load()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (c *Config) Load() error {
	yamlFile, err := ioutil.ReadFile(c.filepath)
	if err != nil {
		log.Errorf("could not read config file %s, err = %v", c.filepath, err)
		return err
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		log.Errorf("could not parse config file %s, err = %v", c.filepath, err)
		return err
	}

	return nil
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
