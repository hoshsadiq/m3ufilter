package m3u

import (
	"github.com/hoshsadiq/m3ufilter/config"
	"reflect"
	"sort"
	"testing"
)

func TestGroupOrder(t *testing.T) {
	conf = &config.Config{
		Core: &config.Core{
			GroupOrder: []string{
				"Group1",
				"Group2",
				"Group3",
				"Group4",
				"Group5",
			},
		},
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
