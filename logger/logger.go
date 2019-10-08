package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
)

var logger *logrus.Logger
var basePkg = getPackage()

func Get() *logrus.Logger {
	if logger == nil {
		logger = logrus.New()

		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			CallerPrettyfier: func(f *runtime.Frame) (string, string) {
				repopath := strings.Split(f.File, basePkg+"/")[1]
				fnc := strings.Split(f.Function, basePkg+"/")[1]
				return fmt.Sprintf("%s()", fnc), fmt.Sprintf("%s:%d", repopath, f.Line)
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
