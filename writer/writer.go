package writer

import (
	"fmt"
	"github.com/grafov/m3u8"
	"github.com/hoshsadiq/m3ufilter/logger"
	"io"
)

var log = logger.Get()

func WriteOutput(Output string, w io.Writer, mediaPlaylist *m3u8.MediaPlaylist) {
	switch Output {
	case "m3u":
		writeM3U(w, mediaPlaylist)
	case "csv":
		writeCsv(w, mediaPlaylist)
	default:
		panic(fmt.Errorf("output type unknown expected m3u|csv, got %s", Output))
	}
}