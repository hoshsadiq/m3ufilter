package m3u

import (
	"github.com/hoshsadiq/m3ufilter/m3u/csv"
	"github.com/hoshsadiq/m3ufilter/m3u/xmltv"
	"regexp"
	"strings"
)

func setSegmentValues(ms *Stream, streamData *csv.StreamData, epgChannel *xmltv.Channel) {
	ms.meta.originalName = ms.Name
	ms.meta.originalId = ms.Id
	ms.meta.epgChannel = epgChannel

	if streamData != nil {
		ms.Id = replaceField("id", ms.Id, streamData.Id) // todo if we change the id, we need to accommodate the epg as well
		ms.Name = replaceField("name", ms.Name, streamData.Name)
		ms.Shift = replaceField("shift", ms.Shift, streamData.Shift)
		ms.Logo = replaceField("logo", ms.Logo, streamData.Logo)
		ms.Group = replaceField("group", ms.Group, streamData.Group)
		ms.ChNo = replaceField("chno", ms.ChNo, streamData.ChNo)
	}

	ms.meta.country = findCountry(ms)
	ms.meta.definition = findDefinition(ms)
	ms.meta.canonicalName = canonicaliseName(ms.Name)
}

func replaceField(name string, oldValue string, newValue string) string {
	if strings.ToLower(newValue) == "clear" {
		log.Tracef("clearing field '%s'", name)
		return ""
	}

	if newValue != "" && newValue != oldValue {
		log.Tracef("%s '%s' replaced with '%s'", name, oldValue, newValue)
		return newValue
	}

	return oldValue
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
