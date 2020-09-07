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

	ms.meta.originalName = ms.Name
	ms.meta.originalId = ms.Id
	ms.meta.epgChannel = epgChannel

	for _, setter := range setters {
		if shouldIncludeStream(ms, setter.Filters, config.CheckStreams{Enabled: false}) {
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

				// todo if we change the id, we need to accommodate the epg as well
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

	ms.meta.country = findCountry(ms)
	ms.meta.definition = findDefinition(ms)
	ms.meta.canonicalName = canonicaliseName(ms.Name)
}

func addDisplayNameToChannel(epgChannel *xmltv.Channel, newValue string) {
	shouldAdd := true
	for _, dn := range epgChannel.DisplayNames {
		if dn.Value == newValue {
			shouldAdd = false
			break
		}
	}

	if shouldAdd {
		epgChannel.DisplayNames = append(epgChannel.DisplayNames, xmltv.DisplayName{Value: newValue})
	}
}

func findCountry(stream *Stream) string {
	var country string
	country = attemptGetCountry(stream.Id, stream.Name)
	if country == "" {
		country = attemptGetCountry(stream.meta.originalId, stream.meta.originalName)
	}
	return country
}

func attemptGetCountry(id, name string) string {
	if id != "" && strings.Count(id, ".") == 1 {
		return strings.ToUpper(strings.Split(id, ".")[1])
	}

	regex := `(?i)\b(` + countries + `)\b`
	r := regexp.MustCompile(regex)
	matches := r.FindStringSubmatch(name)
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
	var definition string
	definition = attemptFindDefinition(stream.Name)
	if definition == "" {
		definition = attemptFindDefinition(stream.meta.originalName)
	}

	return definition
}

func attemptFindDefinition(name string) string {
	regex := `(?i)\b(` + definitions + `)\b`
	r := regexp.MustCompile(regex)
	matches := r.FindStringSubmatch(name)
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
	name = regexWordCallback(name, stopWords, removeWord)
	// todo this still isn't correct
	//if !cache.Regexp("(?i)^Channel \\d+$").Match([]byte(name)) {
	//	name = regexWordCallback(name, "Channel", removeWord)
	//}

	name = regexWordCallback(name, " +", func(s string) string { return " " })
	name = strings.ToLower(name)
	name = strings.Title(name)
	return strings.TrimSpace(name)
}
