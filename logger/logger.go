package logger

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger
var basePkg = getPackage()

func Get() *logrus.Logger {
	if logger == nil {
		logger = logrus.New()

		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				paths := strings.Split(f.File, "/")
				return fmt.Sprintf("%s()", f.Function[len(basePkg)+1:]),
					fmt.Sprintf("%s:%d", paths[5]+"/"+paths[6], f.Line)
			},
		})

		logger.SetReportCaller(true)

		//logger.SetLevel(logrus.DebugLevel)
	}

	return logger
}

func getPackage() string {
	pc, _, _, _ := runtime.Caller(1)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	pkage := strings.Join(parts[0:len(parts)-1], "/")
	return pkage
}
