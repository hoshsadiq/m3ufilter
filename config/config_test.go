package config

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestConfigLoadExample(t *testing.T) {
	filepath := "./testdata/config.yaml"
	expectedConfig := &Config{
		filepath: filepath,
		Core: &Core{
			ServerListen:     "",
			AutoReloadConfig: true,
			UpdateSchedule:   "0 */24 * * *",
			GroupOrder: []string{
				"Entertainment",
				"Family/Kids",
			},
			HttpTimeout:          60,
			HttpMaxRetryAttempts: 5,
			groupOrderMap:        nil,
		},
		Providers: []*Provider{
			{
				Uri:          "file://playlist.m3u",
				Csv:          "file://playlist.csv",
				CheckStreams: CheckStreams{Enabled: false},
			},
		},
		EpgProviders: []*EpgProvider{
			{
				Uri: "file://myepg.xml",
			},
		},
	}

	validateConfig(t, filepath, expectedConfig)
}

func TestConfigLoadCheckStreams(t *testing.T) {
	filepath := "./testdata/config_check_streams.yaml"
	expectedConfig := &Config{
		filepath: filepath,
		Core: &Core{
			ServerListen:         "",
			AutoReloadConfig:     true,
			UpdateSchedule:       "0 */24 * * *",
			GroupOrder:           []string{"Entertainment"},
			HttpTimeout:          60,
			HttpMaxRetryAttempts: 5,
			groupOrderMap:        nil,
		},
		Providers: []*Provider{
			{
				Uri:          "file://playlist.m3u",
				Csv:          "file://playlist.csv",
				CheckStreams: CheckStreams{Enabled: true, Method: "get", Action: "none"},
			},
		},
	}

	validateConfig(t, filepath, expectedConfig)
}

func validateConfig(t *testing.T, filepath string, expectedConfig *Config) {
	conf, err := New(filepath)
	if err != nil {
		t.Fatalf("Failed to loading config file, err = %s", err)
	}

	if !reflect.DeepEqual(conf, expectedConfig) {
		t.Errorf("return config is not as expected \n\t\texpected: %s\n\t\t     got: %s", getConfigAsJson(t, expectedConfig), getConfigAsJson(t, conf))
	}
}

func getConfigAsJson(t *testing.T, conf *Config) string {
	result, err := json.Marshal(conf)
	if err != nil {
		t.Fatalf("failed to convert config to json; err = %s", err)
		return ""
	}

	return string(result)
}
