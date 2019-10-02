package m3u

import (
	"github.com/hoshsadiq/m3ufilter/config"
)

func setSegmentValues(ms *Stream, setters []*config.Setter) {
	var newValue string
	var err error

	for _, setter := range setters {
		if shouldIncludeStream(ms, setter.Filters, false) {
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

			if setter.Shift != "" {
				newValue, err = evaluateStr(ms, setter.Shift)
				if err != nil {
					log.Errorln(err)
				}
				if newValue != ms.Shift {
					log.Tracef("id %s replaced with %s; expr = %v", ms.Shift, newValue, setter.Shift)
				}

				ms.Shift = newValue
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

			if setter.ChNo != "" {
				newValue, err = evaluateStr(ms, setter.ChNo)
				if err != nil {
					log.Errorln(err)
				}
				if newValue != ms.ChNo {
					log.Tracef("title %s replaced with %s; expr = %v", ms.ChNo, newValue, setter.ChNo)
				}

				ms.ChNo = newValue
			}
		}
	}
}
