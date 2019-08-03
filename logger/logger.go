package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

var logger *logrus.Logger

func Get() *logrus.Logger {
	if logger == nil {
		logger = logrus.New()

		logger.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				repopath := strings.Split(f.File, "github.com/hoshsadiq/")[1]
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", repopath, f.Line)
			},
		})

		logger.SetReportCaller(true)

		//logger.SetLevel(logrus.DebugLevel)
	}

	return logger
}
