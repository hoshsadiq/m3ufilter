package writer

import (
	csvw "encoding/csv"
	"github.com/hoshsadiq/m3ufilter/m3u"
	"github.com/hoshsadiq/m3ufilter/m3u/csv"
	"io"
)

var csvHeaders = []string{
	csv.HeaderCuid,
	csv.HeaderSearchName,
	csv.HeaderNumber,
	csv.HeaderId,
	csv.HeaderGroup,
	csv.HeaderName,
	csv.HeaderShift,
	csv.HeaderLogo,
	csv.HeaderUri,
}

func WriteCsv(w io.Writer, streams m3u.Streams) {
	writer := csvw.NewWriter(w)
	defer writer.Flush()

	err := writer.Write(csvHeaders)
	if err != nil {
		log.Errorf("Could not write csv header, err = %s", err)
	}

	for _, stream := range streams {
		printPlaylist(stream, writer)
	}
}

func printPlaylist(pl *m3u.Stream, w *csvw.Writer) {
	row := []string{
		pl.CUID,
		pl.Name,
		pl.ChNo,
		pl.Id,
		pl.Group,
		pl.Name,
		pl.Shift,
		pl.Logo,
		pl.Uri,
	}

	err := w.Write(row)
	if err != nil {
		log.Errorf("Could not write csv row, row = %v, err = %s", row, err)
	}
}
