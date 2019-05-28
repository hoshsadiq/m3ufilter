package writer

import (
	"github.com/grafov/m3u8"
	"io"
)

func writeM3U(w io.Writer, playlists []*m3u8.MediaPlaylist) {
	for _, pl := range playlists {
		_, err := w.Write(pl.Encode().Bytes())
		if err != nil {
			log.Errorf("unable to write new playlist, err = %v", err)
		}
	}
}
