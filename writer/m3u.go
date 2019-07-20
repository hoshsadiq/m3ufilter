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
		log.Fatalf("unable to write extm3u, err = %v", err)
	}

	for _, stream := range streams {
		extinf := getStreamExtinf(stream)
		_, err := w.Write(extinf)
		if err != nil {
			log.Fatalf("unable to write new streams, err = %v, extinf = %v", err, extinf)
		}
	}
}

func getStreamExtinf(stream *m3u.Stream) []byte {
	b := strings.Builder{}
	b.WriteString("\n")
	b.WriteString("#EXTINF:")
	b.WriteString(stream.Duration)

	if stream.ChNo != "" {
		b.WriteString(` tvg-chno=`)
		b.WriteString(stream.ChNo)
	}

	b.WriteString(` tvg-id=`)
	b.WriteString(strconv.Quote(stream.Id))

	if stream.Shift != "" {
		b.WriteString(` tvg-shift=`)
		b.WriteString(stream.Shift)
	}

	b.WriteString(` tvg-name=`)
	b.WriteString(strconv.Quote(stream.Name))

	b.WriteString(` group-title=`)
	b.WriteString(strconv.Quote(stream.Group))

	b.WriteString(` tvg-logo=`)
	b.WriteString(strconv.Quote(stream.Logo))

	b.WriteRune(',')
	b.WriteString(stream.Name)
	b.WriteString("\n")
	b.WriteString(stream.Uri)

	return []byte(b.String())
}
