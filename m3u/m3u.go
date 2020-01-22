package m3u

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/m3u/xmltv"
	"io"
	"strings"
	"time"
)

var groupOrder map[string]int

type Streams []*Stream

func (s Streams) Len() int {
	return len(s)
}

func (s Streams) Less(i, j int) bool {
	iOrder, ok := groupOrder[s[i].Group]
	if !ok {
		return true
	}

	jOrder, ok := groupOrder[s[j].Group]
	if !ok {
		return false
	}

	if iOrder == jOrder {
		return strings.Compare(s[i].meta.canonicalName, s[j].meta.canonicalName) < 0
	}

	return iOrder < jOrder
}

func (s Streams) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

type streamMeta struct {
	originalName string
	originalId   string

	canonicalName string

	definition string
	country    string

	showCountry    bool
	showDefinition bool

	epgChannel *xmltv.Channel
}

func GetMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

type Stream struct {
	Duration string
	Name     string
	Uri      string
	CUID     string `yaml:"CUID"`

	// these are attributes
	ChNo    string `yaml:"chno"`
	Id      string `yaml:"tvg-id"`
	TvgName string `yaml:"tvg-name"`
	Shift   string `yaml:"tvg-shift"`
	Logo    string `yaml:"tvg-logo"`
	Group   string `yaml:"group-title"`

	meta streamMeta
}

func (s Stream) GetName() string {
	name := s.Name
	if s.meta.showCountry {
		name += " " + s.meta.country
	}
	if s.meta.showDefinition {
		name += " " + s.meta.definition
	}

	return name
}

func decode(conf *config.Config, reader io.Reader, providerConfig *config.Provider, epg *xmltv.XMLTV) (Streams, error) {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(reader)
	if err != nil {
		log.Infof("Failed to read from reader to decode m3u due to err = %v", err)
		return nil, err
	}

	groupOrder = conf.GetGroupOrder()

	var extinfLine string
	var urlLine string
	var eof bool
	var epgChannel *xmltv.Channel
	streams := Streams{}

	lines := 0
	start := time.Now()
	for !eof {
		// todo we need to support additional tags - https://github.com/kodi-pvr/pvr.iptvsimple#supported-m3u-and-xmltv-elements
		for !eof && !strings.HasPrefix(extinfLine, "#EXTINF:") {
			extinfLine, eof = getLine(buf)
		}
		if eof {
			break
		}

		urlLine, eof = getLine(buf)
		if eof {
			break
		}

		lines++
		stream, err := parseExtinfLine(extinfLine, urlLine)
		if err != nil {
			if providerConfig.IgnoreParseErrors {
				continue
			}
			return nil, err
		}

		if !shouldIncludeStream(stream, providerConfig.Filters, providerConfig.CheckStreams) {
			continue
		}

		if epg != nil {
			epgChannel = getEpgChannel(stream, epg)
		}

		setSegmentValues(stream, epgChannel, providerConfig.Setters)

		streams = append(streams, stream)
	}
	end := time.Since(start).Truncate(time.Millisecond)

	log.Infof("Matched %d valid streams out of %d. Took %s", len(streams), lines, end)

	return streams, nil
}

func getEpgChannel(stream *Stream, xmltv *xmltv.XMLTV) *xmltv.Channel {
	var convertedDisplayName string
	var displayNameLower string

	streamTvgIdLower := strings.ToLower(stream.Id)
	streamTvgNameLower := strings.ToLower(stream.TvgName)
	streamChannelName := strings.ToLower(stream.Name)

	for _, xmltvChannel := range xmltv.Channels {
		if streamTvgIdLower == strings.ToLower(xmltvChannel.ID) {
			return xmltvChannel
		}
	}
	for _, xmltvChannel := range xmltv.Channels {
		for _, displayName := range xmltvChannel.DisplayNames {
			displayNameLower = strings.ToLower(displayName.Value)
			convertedDisplayName = strings.Replace(displayNameLower, " ", "_", -1)

			if convertedDisplayName == streamTvgNameLower || displayNameLower == streamTvgNameLower {
				return xmltvChannel
			}
		}
	}
	for _, xmltvChannel := range xmltv.Channels {
		for _, displayName := range xmltvChannel.DisplayNames {
			displayNameLower = strings.ToLower(displayName.Value)

			if streamChannelName == displayNameLower {
				return xmltvChannel
			}
		}
	}

	return nil
}

func getLine(buf *bytes.Buffer) (string, bool) {
	var eof bool
	var line string
	var err error
	for !eof {
		line, err = buf.ReadString('\n')
		if err == io.EOF {
			eof = true
		} else if err != nil {
			log.Fatalf("unknown error: %v", err)
		}

		if len(line) < 1 || line == "\r" {
			continue
		}
		break
	}
	return line, eof
}

func parseExtinfLine(attrLine string, urlLine string) (*Stream, error) {
	attrLine = strings.TrimSpace(attrLine)
	urlLine = strings.TrimSpace(urlLine)

	stream := &Stream{Uri: urlLine}
	stream.CUID = GetMD5Hash(urlLine)

	state := "duration"
	key := ""
	value := ""
	quote := "\""
	escapeNext := false
	for i, c := range attrLine {
		if i < 8 {
			continue
		}

		if escapeNext {
			if state == "duration" {
				stream.Duration += string(c)
			} else if state == "keyname" {
				key += string(c)
			} else if state == "quotes" {
				value += string(c)
			}

			escapeNext = false
			continue
		}

		if c == '\\' {
			escapeNext = true
			continue
		}

		if state == "quotes" {
			if string(c) != quote {
				value += string(c)
			} else {
				switch strings.ToLower(key) {
				case "tvg-chno":
					stream.ChNo = value
				case "tvg-id":
					stream.Id = value
				case "tvg-shift":
					stream.Shift = value
				case "tvg-name":
					stream.TvgName = value
				case "tvg-logo":
					stream.Logo = value
				case "group-title":
					stream.Group = value
				}

				key = ""
				value = ""
				state = "start"
			}
			continue
		} else if state == "name" {
			stream.Name += string(c)
			continue
		}

		if c == '"' || c == '\'' {
			if state != "value" {
				return nil, errors.New(fmt.Sprintf("Unexpected character '%s' found, expected '=' for key %s on position %d in line: %s", string(c), key, i, attrLine))
			}
			state = "quotes"
			quote = string(c)
			continue
		}

		if c == ',' {
			state = "name"
			continue
		}

		if state == "keyname" {
			if c == ' ' || c == '\t' {
				key = ""
				state = "start"
			} else if c == '=' {
				state = "value"
			} else {
				key += string(c)
			}
			continue
		}

		if state == "duration" {
			if (c >= 48 && c <= 57) || c == '.' || c == '-' {
				stream.Duration += string(c)
				continue
			}
		}

		if c != ' ' && c != '\t' {
			state = "keyname"
			key += string(c)
		}
	}

	if state == "keyname" && value == "" {
		return nil, errors.New(fmt.Sprintf("Key %s started but no value assigned on line: %s", key, attrLine))
	}

	if state == "quotes" {
		return nil, errors.New(fmt.Sprintf("Unclosed quote on line: %s", attrLine))
	}

	stream.Name = strings.TrimSpace(stream.Name)

	return stream, nil
}
