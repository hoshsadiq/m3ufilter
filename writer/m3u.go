package writer

import (
	"github.com/grafov/m3u8"
	"io"
)

func writeM3U(w io.Writer, mediaPlaylist *m3u8.MediaPlaylist) {
	_, err := w.Write(mediaPlaylist.Encode().Bytes())
	if err != nil {
		log.Errorf("unable to write new playlist, err = %v", mediaPlaylist, err)
	}
}
