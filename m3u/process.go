package m3u

import (
	"github.com/grafov/m3u8"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
)

var log = logger.Get()

func ProcessPlaylist(pl *m3u8.MediaPlaylist, providerConfig *config.Provider) (*m3u8.MediaPlaylist, error) {
	p, e := m3u8.NewMediaPlaylist(pl.Count(), pl.Count())
	if e != nil {
		return nil, e
	}

	for _, segment := range pl.Segments {
		if segment == nil {
			continue
		}

		segment = segment.Clone()
		replace(segment, providerConfig.Replacements)

		if !shouldIncludeSegment(segment, providerConfig.Filters) {
			continue
		}

		setSegmentValues(segment, providerConfig.Setters)
		//renameGroups(segment, providerConfig.Groups)

		e = p.AppendSegment(segment)
		if e != nil {
			return nil, e
		}
	}

	err := p.SetWinSize(p.Count())
	if err != nil {
		return nil, err
	}

	return p, nil
}
