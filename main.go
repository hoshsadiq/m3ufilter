package main

import (
	"bufio"
	"fmt"
	"github.com/grafov/m3u8"
	"github.com/maja42/goval"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
)

var evaluator = goval.NewEvaluator()

func main() {
	yamlFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(err)
	}

	var config *Config
	err = yaml.Unmarshal([]byte(yamlFile), &config)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	t := &http.Transport{}
	t.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))
	c := &http.Client{Transport: t}
	for _, provider := range config.Providers {
		resp, err := c.Get(provider.Uri)
		if err != nil {
			panic(err)
		}
		defer resp.Body.Close()

		p, listType, err := m3u8.DecodeFrom(bufio.NewReader(resp.Body), false)
		if err != nil {
			panic(err)
		}

		switch listType {
		case m3u8.MEDIA:
			mediapl := p.(*m3u8.MediaPlaylist)
			newp, err := filterSegments(mediapl, provider)
			if err != nil {
				panic(err)
			}
			fmt.Printf("%+v\n", newp)
		default:
			fmt.Printf("Found a non-media thing. Code needs to be updated!")
		}
	}
}

func filterSegments(pl *m3u8.MediaPlaylist, providerConfig *Provider) (*m3u8.MediaPlaylist, error) {
	p, e := m3u8.NewMediaPlaylist(pl.Count(), pl.Count())
	if e != nil {
		return nil, e
	}

	for _, segment := range pl.Segments {
		if segment == nil {
			continue
		}

		segment = segment.Clone()
		doReplace(segment, providerConfig.Replacements)

		if !shouldIncludeSegment(segment, providerConfig.Filters) {
			continue
		}

		setSegmentValues(segment, providerConfig.Setters)
		//renameGroups(segment, providerConfig.Groups)

		e = p.AppendSegment(segment)
		if e != nil {
			return nil, e
		}
	}

	err := p.SetWinSize(p.Count())
	if err != nil {
		return nil, err
	}

	return p, nil
}
