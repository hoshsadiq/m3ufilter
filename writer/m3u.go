package writer

import (
	"github.com/hoshsadiq/m3ufilter/m3u"
	"io"
	"strconv"
	"strings"
)

func writeM3U(w io.Writer, streams []*m3u.Stream) {
	_, err := w.Write([]byte("#EXTM3U"))
	if err != nil {
		log.Fatalf("unable to write new streams, err = %v", err)
	}

	for _, stream := range streams {
		_, err := w.Write(getStreamExtinf(stream))
		if err != nil {
			log.Fatalf("unable to write new streams, err = %v", err)
		}
	}
}

func getStreamExtinf(stream *m3u.Stream) []byte {
	var str strings.Builder
	str.WriteString("\n")
	str.WriteString("#EXTINF:")
	str.WriteString(stream.Duration)

	str.WriteString(` tvg-id=`)
	str.WriteString(strconv.Quote(stream.Id))

	str.WriteString(` tvg-name=`)
	str.WriteString(strconv.Quote(stream.Name))

	str.WriteString(` tvg-logo=`)
	str.WriteString(strconv.Quote(stream.Logo))

	str.WriteString(` group-title=`)
	str.WriteString(strconv.Quote(stream.Group))

	str.WriteRune(',')
	str.WriteString(stream.Name)
	str.WriteString("\n")
	str.WriteString(stream.Uri)

	return []byte(str.String())
}
