package m3u

import (
	"errors"
	"fmt"
	"github.com/grafov/m3u8"
)

func GetAttr(ms *m3u8.MediaSegment, attrKey string) (*m3u8.Attribute, error) {
	for _, attr := range ms.Attributes {
		if attr.Key == attrKey {
			return attr, nil
		}
	}

	return nil, errors.New(fmt.Sprintf("attribute %s not found", attrKey))
}

func SetAttr(ms *m3u8.MediaSegment, attrKey string, newValue string) {
	attr, err := GetAttr(ms, attrKey)
	if err == nil {
		if newValue != attr.Value {
			log.Tracef("attr %v has new value = %s", attr, newValue)
		}
		attr.Value = newValue
	} else {
		attr = &m3u8.Attribute{Key: attrKey, Value: newValue}
		log.Tracef("adding new attr %v", attr)
		ms.Attributes = append(ms.Attributes, attr)
	}
}
