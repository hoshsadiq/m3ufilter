package m3u

import (
	"github.com/hoshsadiq/m3ufilter/config"
)

func setSegmentValues(ms *Stream, setters []*config.Setter) {
	var newValue string
	var err error

	for _, setter := range setters {
		if len(setter.Filters) == 0 || shouldIncludeSegment(ms, setter.Filters) {
			if setter.Name != "" {
				newValue, err = evaluateStr(ms, setter.Name)
				if err != nil {
					log.Errorln(err)
				}
				if newValue != ms.Name {
					log.Tracef("title %s replaced with %s; expr = %v", ms.Name, newValue, setter.Name)
				}

				ms.Name = newValue
			}

			if setter.Id != "" {
				newValue, err = evaluateStr(ms, setter.Id)
				if err != nil {
					log.Errorln(err)
				}
				if newValue != ms.Id {
					log.Tracef("id %s replaced with %s; expr = %v", ms.Id, newValue, setter.Id)
				}

				ms.Id = newValue
			}

			if setter.Logo != "" {
				newValue, err = evaluateStr(ms, setter.Logo)
				if err != nil {
					log.Errorln(err)
				}
				if newValue != ms.Logo {
					log.Tracef("title %s replaced with %s; expr = %v", ms.Logo, newValue, setter.Logo)
				}

				ms.Logo = newValue
			}

			if setter.Group != "" {
				newValue, err = evaluateStr(ms, setter.Group)
				if err != nil {
					log.Errorln(err)
				}
				if newValue != ms.Group {
					log.Tracef("title %s replaced with %s; expr = %v", ms.Group, newValue, setter.Group)
				}

				ms.Group = newValue
			}
		}
	}
}
