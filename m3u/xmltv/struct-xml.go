package xmltv

import (
	"encoding/xml"
	"io"
)

const timefmt = "20060102150405 -0700"

// XMLTV : XMLTV Datei
type XMLTV struct {
	Generator string   `xml:"generator-info-name,attr"`
	Source    string   `xml:"source-info-name,attr"`
	XMLName   xml.Name `xml:"tv"`

	Channel []*Channel `xml:"channel"`
	Program []*Program `xml:"programme"`
}

// Channel : Kanäle
type Channel struct {
	ID          string        `xml:"id,attr"`
	DisplayName []DisplayName `xml:"display-name"`
	Icon        Icon          `xml:"icon"`
}

// DisplayName : Kanalname
type DisplayName struct {
	Value string `xml:",chardata"`
}

// Icon : Senderlogo
type Icon struct {
	Src string `xml:"src,attr"`
}

// Program : Programme
type Program struct {
	Channel string `xml:"channel,attr"`
	Start   string `xml:"start,attr"`
	Stop    string `xml:"stop,attr"`

	Title           []*Title         `xml:"title"`
	SubTitle        []*SubTitle      `xml:"sub-title"`
	Desc            []*Desc          `xml:"desc"`
	Category        []*Category      `xml:"category"`
	Country         []*Country       `xml:"country"`
	EpisodeNum      []*EpisodeNum    `xml:"episode-num"`
	Poster          []Poster         `xml:"icon"`
	Language        []*Language      `xml:"language"`
	Video           Video            `xml:"video"`
	Date            string           `xml:"date"`
	PreviouslyShown *PreviouslyShown `xml:"previously-shown"`
	New             *New             `xml:"new"`
	Live            *Live            `xml:"live"`
}

// todo ideally some of these lang + value structs become map[lang]value instead of lists
type Title struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

type SubTitle struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

type Desc struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

type Category struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

type Language struct {
	Value string `xml:",chardata"`
}

type Country struct {
	Lang  string `xml:"lang,attr"`
	Value string `xml:",chardata"`
}

type EpisodeNum struct {
	System string `xml:"system,attr"`
	Value  string `xml:",chardata"`
}

type Poster struct {
	Height string `xml:"height,attr"`
	Src    string `xml:"src,attr"`
	Value  string `xml:",chardata"`
	Width  string `xml:"width,attr"`
}

type Video struct {
	Aspect  string `xml:"aspect,omitempty"`
	Colour  string `xml:"colour,omitempty"`
	Present string `xml:"present,omitempty"`
	Quality string `xml:"quality,omitempty"`
}

type PreviouslyShown struct {
	Start string `xml:"start,attr"`
}

type New struct {
	Value string `xml:",chardata"`
}

type Live struct {
	Value string `xml:",chardata"`
}

func Load(r io.Reader, xmltv *XMLTV) error {
	return xml.NewDecoder(r).Decode(xmltv)
}
