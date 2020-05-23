package m3u

import (
	"encoding/json"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestGroupOrder(t *testing.T) {
	groupOrder = map[string]int{
		"Group1": 1,
		"Group2": 2,
		"Group3": 3,
		"Group4": 4,
		"Group5": 5,
	}

	streams := Streams{
		&Stream{Group: "Group2"},
		&Stream{Group: "Group4"},
		&Stream{Group: "Group4"},
		&Stream{Group: "Group4"},
		&Stream{Group: "Group4"},
		&Stream{Group: "Group4"},
		&Stream{Group: "Group1"},
		&Stream{Group: "Group3"},
	}

	sort.Sort(streams)

	expectedStreamsOrder := Streams{
		&Stream{Group: "Group1"},
		&Stream{Group: "Group2"},
		&Stream{Group: "Group3"},
		&Stream{Group: "Group4"},
		&Stream{Group: "Group4"},
		&Stream{Group: "Group4"},
		&Stream{Group: "Group4"},
		&Stream{Group: "Group4"},
	}
	if !reflect.DeepEqual(expectedStreamsOrder, streams) {
		t.Errorf("Grouping was not ordered corrected. Expected %v, got %v", expectedStreamsOrder, streams)
	}
}

func TestExtinfLineParser(t *testing.T) {
	tests := []struct {
		attrLine       string
		urlLine        string
		expectedStream *Stream
		errorText      string
	}{
		{
			attrLine: `#EXTINF:-1 tvg-id="channel.uk" tvg-name="Channel Name" tvg-logo="http://imgur.com/img.png" group-title="Channel Group",Stream Name`,
			urlLine:  "http://somestreamer.com/mystream.mp4",
			expectedStream: &Stream{
				Duration: "-1",
				Group:    "Channel Group",
				TvgName:  "Channel Name",
				ChNo:     "",
				Shift:    "",
				Uri:      "http://somestreamer.com/mystream.mp4",
				Logo:     "http://imgur.com/img.png",
				Name:     "Stream Name",
				Id:       "channel.uk",
				CUID:     "7b5d2b7ac72fcd316cce413671ca0ada",
			},
		},
		{
			attrLine: `#EXTINF:-1 tvg-id="channel.uk" tvg-name="Channel\" Name" tvg-logo="http://imgur.com/img.png" group-title="Channel Group",Stream Name`,
			urlLine:  "http://somestreamer.com/mystream.mp4",
			expectedStream: &Stream{
				Duration: "-1",
				Group:    "Channel Group",
				TvgName:  `Channel" Name`,
				ChNo:     "",
				Shift:    "",
				Uri:      "http://somestreamer.com/mystream.mp4",
				Logo:     "http://imgur.com/img.png",
				Name:     "Stream Name",
				Id:       "channel.uk",
				CUID:     "7b5d2b7ac72fcd316cce413671ca0ada",
			},
		},
		{
			attrLine:  `#EXTINF:-1 tvg-id="channel.uk tvg-name="Channel Name" tvg-logo="http://imgur.com/img.png" group-title="Channel Group",Stream Name`,
			urlLine:   "http://somestreamer.com/mystream.mp4",
			errorText: `Unexpected character '"' found, expected '=' for key Name on position 52 in line: #EXTINF:-1 tvg-id="channel.uk tvg-name="Channel Name" tvg-logo="http://imgur.com/img.png" group-title="Channel Group",Stream Name`,
		},
	}

	for i, test := range tests {
		stream, err := parseExtinfLine(test.attrLine, test.urlLine)
		if err != nil {
			if test.errorText == "" {
				t.Errorf("test %d, did not expect an error, got: %s", i, err)
			} else if test.errorText != err.Error() {
				t.Errorf("test %d, expected err = %s, got %s", i, test.errorText, err)
			}
		} else {
			if test.errorText != "" {
				t.Errorf("test %d, expected err = %s, got nil", i, test.errorText)
			}

			if !reflect.DeepEqual(test.expectedStream, stream) {
				t.Errorf("test %d, expected stream = %v, got %v", i, test.expectedStream, stream)
			}
		}

	}
}

func TestDecoder(t *testing.T) {
	conf := &config.Config{
		Core: &config.Core{
			GroupOrder: []string{},
		},
	}

	logger.Get().SetLevel(logrus.WarnLevel)

	_ = filepath.Walk("testdata/m3u-decode", func(path string, info os.FileInfo, err error) error {
		//var testData interface{}
		var testData simpleTest

		ext := filepath.Ext(path)
		if !info.IsDir() && (ext == ".yaml" || ext == ".yml") {
			runTest(path, t, testData, ext, conf)
		}

		return nil
	})
}

func runTest(path string, t *testing.T, testData simpleTest, ext string, conf *config.Config) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	err = yaml.Unmarshal(yamlFile, &testData)
	if err != nil {
		t.Fatal(err)
	}

	m3ufile := strings.TrimSuffix(path, ext) + ".m3u"
	f, err := os.Open(m3ufile)
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			t.Errorf("could not close file %s, err = %v", m3ufile, err)
		}
	}()

	streams, err := decode(conf, f, &testData.Provider)
	if testData.ExpectedError != "__no_error__" && (err == nil || err.Error() != testData.ExpectedError) {
		t.Errorf("Test %s failed. Expected err %s, but got %s", path, testData.ExpectedError, err)
		return
	}

	if !reflect.DeepEqual(streams, testData.Streams) {
		expectedStreams, err := json.Marshal(testData.Streams)
		if err != nil {
			t.Fatalf("Failed to get json for expected streams; err = %v", err)
		}
		actualStreams, err := json.Marshal(streams)
		if err != nil {
			t.Fatalf("Failed to get json for actual streams; err = %v", err)
		}

		t.Logf("Test %s failed.", path)
		t.Logf("  Expected streans: %s", expectedStreams)
		t.Logf("  Got:              %s", actualStreams)
		t.Fail()
	}
}

type simpleTest struct {
	ExpectedError string `yaml:"expected_error"`
	Streams       Streams
	Provider      config.Provider
}
