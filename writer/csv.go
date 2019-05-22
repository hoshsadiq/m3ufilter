package writer

import (
	"encoding/csv"
	"github.com/grafov/m3u8"
	"github.com/hoshsadiq/m3ufilter/util"
	"io"
)

var csvHeaders = []string{
	"tvg-id",
	"group-title",
	"tvg-name",
	"tvg-logo",
	"channel-name",
}

func writeCsv(w io.Writer, pl *m3u8.MediaPlaylist) {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	err := writer.Write(csvHeaders)
	if err != nil {
		log.Errorf("Could not write csv header, err = %s", err)
	}

	for _, ms := range pl.Segments {
		if ms == nil { // todo why do we get a nil value after the last item?
			continue
		}

		row := []string{
			util.GetAttr(ms, "tvg-id").Value,
			util.GetAttr(ms, "group-title").Value,
			util.GetAttr(ms, "tvg-name").Value,
			util.GetAttr(ms, "tvg-logo").Value,
			util.GetAttr(ms, "channel-name").Value,
		}

		err = writer.Write(row)
		if err != nil {
			log.Errorf("Could not write csv row, row = %v, err = %s", row, err)
		}
	}
}
