package helper

import (
	"github.com/hoshsadiq/m3ufilter/logger"
	"io"
)

var log = logger.Get()

func Close(rw io.Closer, info string) func() {
	return func() {
		err := rw.Close()
		if err != nil {
			log.Errorf("could not close request body%s, err = %v", info, err)
		}
	}
}
