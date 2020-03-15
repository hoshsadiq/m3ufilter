package xmltv

import (
	"encoding/xml"
	"github.com/hoshsadiq/m3ufilter/m3u/filter"
	"io"
	"strings"
	"time"
)

// These structs are copied from the xteve project, with some slight adjustments

const TimeFormat = "20060102150405 -0700"

type Time time.Time

// XMLTV
type XMLTV struct {
	XMLName xml.Name `xml:"tv"`

	GeneratorName *string `xml:"generator-info-name,attr,omitempty"`
	GeneratorUrl  *string `xml:"generator-info-url,attr,omitempty"`

	SourceInfoName *string `xml:"source-info-name,attr,omitempty"`
	SourceInfoUrl  *string `xml:"source-info-url,attr,omitempty"`
	SourceUrl      *string `xml:"source-data-url,attr,omitempty"`

	Date *Time `xml:"date,attr,omitempty"`

	Channels   []*Channel   `xml:"channel"`
	Programmes []*Programme `xml:"programme"`
}

func (x *XMLTV) SetGenerator(name string, url string) {
	x.GeneratorName = &name
	x.GeneratorUrl = &url
}

func (x *XMLTV) SetSource(infoName string, infoUrl string, url string) {
	x.SourceInfoName = &infoName
	x.SourceInfoUrl = &infoUrl
	x.SourceUrl = &url
}

// Channels
type Channel struct {
	ID           string        `xml:"id,attr"`
	DisplayNames []DisplayName `xml:"display-name"`
	Icon         Icon          `xml:"icon,omitempty"`
	Url          Url           `xml:"url,omitempty"`
}

// DisplayName
type DisplayName struct {
	Lang  string `xml:"lang,attr,omitempty"`
	Value string `xml:",chardata"`
}

type Url string

type Icon struct {
	Src IconSrc `xml:"src,attr"`
}

type IconSrc string

// Programmes
type Programme struct {
	Channel string `xml:"channel,attr"`
	Start   *Time  `xml:"start,attr"`
	Stop    *Time  `xml:"stop,attr"`

	Title       []*Title       `xml:"title"`
	SubTitle    []*SubTitle    `xml:"sub-title,omitempty"`
	Description []*Description `xml:"desc"`

	Date *Time `xml:"date,omitempty"`

	Category        []*Category      `xml:"category"`
	Country         []*Country       `xml:"country"`
	EpisodeNum      []*EpisodeNum    `xml:"episode-num"`
	Poster          []Poster         `xml:"icon"`
	Language        []*Language      `xml:"language"`
	Video           *Video           `xml:"video,omitempty"`
	PreviouslyShown *PreviouslyShown `xml:"previously-shown"`
	New             *New             `xml:"new"`
	Live            *Live            `xml:"live"`
	Url             *Url             `xml:"url,omitempty"`

	// todo missing fields:
	//  attrs:
	//   - pdc-start
	//   - vps-start
	//   - showview
	//   - videoplus
	//   - clumpidx
	//  elements:
	//   - credits
	//   - keyword
	//   - orig-language
	//   - length
	//   - icon
	//   - audio
	//   - premiere
	//   - last-chance
	//   - rating
	//   - star-rating
	//   - review
}

// todo ideally some of these lang + value structs become map[lang]value instead of lists
type Title struct {
	Lang  string `xml:"lang,attr,omitempty"`
	Value string `xml:",chardata"`
}

type SubTitle struct {
	Lang  string `xml:"lang,attr,omitempty"`
	Value string `xml:",chardata"`
}

type Description struct {
	Lang  string `xml:"lang,attr,omitempty"`
	Value string `xml:",chardata"`
}

type Category struct {
	Lang  string `xml:"lang,attr,omitempty"`
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

func (d *Time) UnmarshalXMLAttr(attr xml.Attr) (err error) {
	datetime, err := time.Parse(TimeFormat, attr.Value)
	if err != nil {
		return err
	}

	*d = Time(datetime)

	return
}

func (d Time) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  name,
		Value: time.Time(d).Format(TimeFormat),
	}, nil
}

func (d *IconSrc) UnmarshalXMLAttr(attr xml.Attr) (err error) {
	*d = IconSrc(filter.EnsureUniqueUrls(attr.Value))

	return nil
}

func (d IconSrc) MarshalXMLAttr(name xml.Name) (xml.Attr, error) {
	return xml.Attr{
		Name:  name,
		Value: string(d),
	}, nil
}

func Load(r io.Reader, xmltv *XMLTV) (err error) {
	err = xml.NewDecoder(r).Decode(xmltv)
	return err
}

func Dump(w io.Writer, xmltv *XMLTV, prettify bool) error {
	enc := xml.NewEncoder(w)

	_, err := w.Write([]byte(strings.TrimSpace(xml.Header)))
	if err != nil {
		return err
	}

	if prettify {
		enc.Indent("", "  ")

		_, err = w.Write([]byte("\n"))
		if err != nil {
			return err
		}
	}

	return enc.Encode(xmltv)
}
