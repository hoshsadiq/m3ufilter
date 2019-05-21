package m3u

import (
	"github.com/grafov/m3u8"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
)

var log = logger.Get()

func ProcessPlaylist(pl *m3u8.MediaPlaylist, providerConfig *config.Provider, syncTitleName bool) (*m3u8.MediaPlaylist, error) {
	p, err := m3u8.NewMediaPlaylist(pl.Count(), pl.Count())
	if err != nil {
		return nil, err
	}

	for _, segment := range pl.Segments {
		if segment == nil {
			continue
		}

		segment = segment.Clone()
		//segment.Title, err = strconv.Unquote("\"" + segment.Title + "\"")
		//if err != nil {
		//	log.Errorf("error unquoting %s", segment.Title)
		//}

		replace(segment, providerConfig.Replacements, syncTitleName)

		if !shouldIncludeSegment(segment, providerConfig.Filters) {
			continue
		}

		setSegmentValues(segment, providerConfig.Setters, syncTitleName)
		//renameGroups(segment, providerConfig.Groups)

		err = p.AppendSegment(segment)
		if err != nil {
			return nil, err
		}
	}

	err = p.SetWinSize(p.Count())
	if err != nil {
		return nil, err
	}

	return p, nil
}
