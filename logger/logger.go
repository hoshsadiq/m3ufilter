package logger

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/sirupsen/logrus"
)

var logger *logrus.Logger
var basePkg = getPackage()
var appPath string

func Setup(baseAppPath string) {
	appPath = baseAppPath
	if runtime.GOOS == "windows" {
		appPath = strings.ReplaceAll(appPath, string(os.PathSeparator), "/")
	}
}

func Get() *logrus.Logger {
	if logger == nil {
		logger = logrus.New()

		logger.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
			CallerPrettyfier: func(f *runtime.Frame) (fn string, file string) {
				file = strings.Split(f.File, appPath+"/")[1]
				fn = f.Function
				if fn != "main.main" {
					fn = strings.Split(f.Function, basePkg+"/")[1]
				}

				return fmt.Sprintf("%s()", fn), fmt.Sprintf("%s:%d", file, f.Line)
			},
		})

		logger.SetReportCaller(true)

		//logger.SetLevel(logrus.DebugLevel)
	}

	return logger
}

func getPackage() (pkg string) {
	pc, _, _, _ := runtime.Caller(1)
	parts := strings.Split(runtime.FuncForPC(pc).Name(), "/")
	return strings.Join(parts[0:len(parts)-1], "/")
}
