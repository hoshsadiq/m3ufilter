package m3u

import (
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/m3u/xmltv"
	"regexp"
	"strings"
)

func setSegmentValues(ms *Stream, epgChannel *xmltv.Channel, setters []*config.Setter) {
	var newValue string
	var err error

	ms.meta.country = findCountry(ms)
	ms.meta.definition = findDefinition(ms)
	ms.meta.canonicalName = canonicaliseName(ms.Name)
	ms.meta.originalName = ms.Name
	ms.meta.originalId = ms.Id
	ms.meta.epgChannel = epgChannel

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
					log.Tracef("shift %s replaced with %s; expr = %v", ms.Shift, newValue, setter.Shift)
				}

				ms.Shift = newValue
			}

			if setter.Logo != "" {
				newValue, err = evaluateStr(ms, setter.Logo)
				if err != nil {
					log.Errorln(err)
				}
				if newValue != ms.Logo {
					log.Tracef("logo %s replaced with %s; expr = %v", ms.Logo, newValue, setter.Logo)
				}

				ms.Logo = newValue
			}

			if setter.Group != "" {
				newValue, err = evaluateStr(ms, setter.Group)
				if err != nil {
					log.Errorln(err)
				}
				if newValue != ms.Group {
					log.Tracef("group %s replaced with %s; expr = %v", ms.Group, newValue, setter.Group)
				}

				ms.Group = newValue
			}

			if setter.ChNo != "" {
				newValue, err = evaluateStr(ms, setter.ChNo)
				if err != nil {
					log.Errorln(err)
				}
				if newValue != ms.ChNo {
					log.Tracef("chno %s replaced with %s; expr = %v", ms.ChNo, newValue, setter.ChNo)
				}

				ms.ChNo = newValue
			}
		}
	}
}

func findCountry(stream *Stream) string {
	if stream.Id != "" && strings.Count(stream.Id, ".") == 1 {
		return strings.ToUpper(strings.Split(stream.Id, ".")[1])
	}

	regex := `(?i)\b(` + countries + `)\b`
	r := regexp.MustCompile(regex)
	matches := r.FindStringSubmatch(stream.Name)
	if matches != nil {
		country := strings.ToUpper(matches[0])

		if val, ok := countryOverrides[country]; ok {
			country = val
		}

		return country
	}

	return ""
}

func findDefinition(stream *Stream) string {
	regex := `(?i)\b(` + definitions + `)\b`
	r := regexp.MustCompile(regex)
	matches := r.FindStringSubmatch(stream.Name)
	if matches != nil {
		definition := strings.ToUpper(matches[0])
		if val, ok := definitionOverrides[definition]; ok {
			definition = val
		}

		return definition
	}

	return ""
}

func canonicaliseName(name string) string {
	name = strings.Replace(name, ":", "", -1)
	name = strings.Replace(name, "|", "", -1)
	name = regexWordCallback(name, countries, removeWord)
	name = regexWordCallback(name, definitions, removeWord)
	name = regexWordCallback(name, "TV", removeWord)
	// todo this still isn't correct
	//if !cache.Regexp("(?i)^Channel \\d+$").Match([]byte(name)) {
	//	name = regexWordCallback(name, "Channel", removeWord)
	//}

	name = strings.Title(name)
	name = strings.ToLower(name)
	return strings.TrimSpace(name)
}
