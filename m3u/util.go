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
