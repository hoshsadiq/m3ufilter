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
			UpdateSchedule:   "* */24 * * *",
			Output:           "m3u",
			GroupOrder: []string{
				"Entertainment",
				"Family/Kids",
			},
			HttpTimeout:          60,
			HttpMaxRetryAttempts: 5,
			groupOrderMap:        nil,
			Canonicalise: Canonicalise{
				Enable:         true,
				DefaultCountry: "uk",
			},
		},
		Providers: []*Provider{
			{
				Uri:               "file://playlist.m3u",
				IgnoreParseErrors: true,
				CheckStreams:      CheckStreams{Enabled: false},
				Filters: []string{
					"Group in [\n  \"Documentaries\",\n  \"Entertainment\",\n  \"Kids\",\n  \"Movies\",\n  \"Music\",\n  \"News\",\n]\n",
				},
				Setters: []*Setter{
					{
						Name:    `replace(Name, " +", " ")`,
						Filters: nil,
					},
				},
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

func TestConfigLoadNewCheckStreams(t *testing.T) {
	filepath := "./testdata/config_new_check_streams.yaml"
	expectedConfig := &Config{
		filepath: filepath,
		Core: &Core{
			ServerListen:         "",
			AutoReloadConfig:     true,
			UpdateSchedule:       "* */24 * * *",
			Output:               "m3u",
			GroupOrder:           []string{"Entertainment"},
			HttpTimeout:          60,
			HttpMaxRetryAttempts: 5,
			groupOrderMap:        nil,
			Canonicalise: Canonicalise{
				Enable:         true,
				DefaultCountry: "uk",
			},
		},
		Providers: []*Provider{
			{
				Uri:               "file://playlist.m3u",
				IgnoreParseErrors: false,
				CheckStreams:      CheckStreams{Enabled: true, Method: "get", Action: "none"},
				Filters:           []string{`Group in ["Documentaries"]`},
				Setters:           nil,
			},
		},
	}

	validateConfig(t, filepath, expectedConfig)
}

func TestConfigLoadOldCheckStreams(t *testing.T) {
	filepath := "./testdata/config_old_check_streams.yaml"
	expectedConfig := &Config{
		filepath: filepath,
		Core: &Core{
			ServerListen:         "",
			AutoReloadConfig:     true,
			UpdateSchedule:       "* */24 * * *",
			Output:               "m3u",
			GroupOrder:           []string{"Entertainment"},
			HttpTimeout:          60,
			HttpMaxRetryAttempts: 5,
			groupOrderMap:        nil,
			Canonicalise: Canonicalise{
				Enable:         true,
				DefaultCountry: "uk",
			},
		},
		Providers: []*Provider{
			{
				Uri:               "file://playlist.m3u",
				IgnoreParseErrors: false,
				CheckStreams:      CheckStreams{Enabled: true, Method: "head", Action: "remove"},
				Filters:           []string{`Group in ["Documentaries"]`},
				Setters:           nil,
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
