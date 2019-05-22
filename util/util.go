package util

import (
	"github.com/grafov/m3u8"
	"github.com/hoshsadiq/m3ufilter/logger"
)

var log = logger.Get()

func GetAttr(ms *m3u8.MediaSegment, attrKey string) *m3u8.Attribute {
	if ms == nil {
		panic("yo bruv")
	}

	for _, attr := range ms.Attributes {
		if attr.Key == attrKey {
			return attr
		}
	}

	return &m3u8.Attribute{}
}

func SetAttr(ms *m3u8.MediaSegment, attrKey string, newValue string) {
	attr := GetAttr(ms, attrKey)
	if newValue != attr.Value {
		log.Tracef("attr %v has new value = %s", attr, newValue)
	}
	attr.Value = newValue
}
