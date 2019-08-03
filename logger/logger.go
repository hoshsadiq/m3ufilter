package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

var logger *logrus.Logger
var callerinfo = getPackage() + "/"

func Get() *logrus.Logger {
	if logger == nil {
		logger = logrus.New()

		logger.SetFormatter(&logrus.TextFormatter{
			DisableColors: true,
			FullTimestamp: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				repopath := strings.Split(f.File, callerinfo)[1]
				return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", repopath, f.Line)
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
