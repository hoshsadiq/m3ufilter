package logger

import (
	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger

func Get() *logrus.Logger {
	if logger == nil {
		logger = logrus.New()

		logger.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
		})

		//logger.SetReportCaller(true)

		//logger.SetLevel(logrus.DebugLevel)
	}

	return logger
}
