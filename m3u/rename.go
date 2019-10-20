package m3u

import (
	"github.com/hoshsadiq/m3ufilter/config"
	"regexp"
	"strings"
)

func setSegmentValues(ms *Stream, setters []*config.Setter) {
	var newValue string
	var err error

	ms.meta.country = findCountry(ms)
	ms.meta.definition = findDefinition(ms)
	ms.meta.canonicalName = canonicaliseName(ms.Name)
	ms.meta.originalName = ms.Name

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

			if setter.Attributes.Id != "" {
				newValue, err = evaluateStr(ms, setter.Attributes.Id)
				if err != nil {
					log.Errorln(err)
				}
				if newValue != ms.Id {
					log.Tracef("id %s replaced with %s; expr = %v", ms.Id, newValue, setter.Attributes.Id)
				}

				ms.Id = newValue
			}

			if setter.Attributes.Shift != "" {
				newValue, err = evaluateStr(ms, setter.Attributes.Shift)
				if err != nil {
					log.Errorln(err)
				}
				if newValue != ms.Shift {
					log.Tracef("id %s replaced with %s; expr = %v", ms.Shift, newValue, setter.Attributes.Shift)
				}

				ms.Shift = newValue
			}

			if setter.Attributes.Logo != "" {
				newValue, err = evaluateStr(ms, setter.Attributes.Logo)
				if err != nil {
					log.Errorln(err)
				}
				if newValue != ms.Logo {
					log.Tracef("title %s replaced with %s; expr = %v", ms.Logo, newValue, setter.Attributes.Logo)
				}

				ms.Logo = newValue
			}

			if setter.Attributes.Group != "" {
				newValue, err = evaluateStr(ms, setter.Attributes.Group)
				if err != nil {
					log.Errorln(err)
				}
				if newValue != ms.Group {
					log.Tracef("title %s replaced with %s; expr = %v", ms.Group, newValue, setter.Attributes.Group)
				}

				ms.Group = newValue
			}

			if setter.Attributes.ChNo != "" {
				newValue, err = evaluateStr(ms, setter.Attributes.ChNo)
				if err != nil {
					log.Errorln(err)
				}
				if newValue != ms.ChNo {
					log.Tracef("title %s replaced with %s; expr = %v", ms.ChNo, newValue, setter.Attributes.ChNo)
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

	regex := `(?i)\b(` + countryReplaces + `)\b`
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
	regex := `(?i)\b(` + definitionReplaces + `)\b`
	r := regexp.MustCompile(regex)
	matches := r.FindStringSubmatch(stream.Name)
	if matches != nil {
		return strings.ToUpper(matches[0])
	}

	return ""
}

func canonicaliseName(name string) string {
	name = regexWordCallback(name, countryReplaces, removeWord)
	name = regexWordCallback(name, "TV|Channel", removeWord)
	name = regexWordCallback(name, definitionReplaces, removeWord)

	name = strings.Replace(name, ":", "", -1)
	name = strings.Replace(name, "|", "", -1)

	name = strings.Title(name)
	name = strings.ToLower(name)
	return strings.TrimSpace(name)
}
