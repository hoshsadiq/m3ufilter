package main

import (
	"flag"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
	"github.com/hoshsadiq/m3ufilter/m3u"
	"github.com/hoshsadiq/m3ufilter/server"
	"github.com/hoshsadiq/m3ufilter/writer"
	"github.com/mitchellh/go-homedir"
	"os"
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
	appPath    = filepath.Dir(filepath.Dir(filepath.Dir(b)))
)

func main() {
	configFile := flag.String("config", "~/.m3u.conf", "Config file location")
	playlistOutput := flag.String("playlist", "", "Where to output the playlist data. Ignored when using -server flag. Defaults to stdout")
	logOutput := flag.String("log", "", "Where to output logs. Defaults to stderr")
	versionFlag := flag.Bool("version", false, "show version and exit")
	flag.Parse()

	if *versionFlag {
		config.ShowVersion()
	}

	path, e := homedir.Expand(*configFile)
	if e != nil {
		panic(e)
	}

	logger.Setup(appPath)
	run(path, fd(*playlistOutput, false), fd(*logOutput, true))
}

func run(configFilename string, stdout *os.File, stderr *os.File) {
	log := logger.Get()
	log.SetOutput(stderr)

	conf := config.New(configFilename)
	m3u.InitClient(conf)
	if conf.Core.ServerListen != "" {
		server.Serve(conf)
	} else {
		playlist, _ /*todo epg*/, _ := m3u.ProcessConfig(conf)
		writer.WriteOutput(conf.Core.Output, stdout, playlist)
	}
}

func fd(filename string, defaultStderr bool) *os.File {
	if filename == "-" {
		return os.Stdout
	}
	if filename == "" && defaultStderr {
		return os.Stderr
	}
	if filename == "" {
		return os.Stdout
	}

	fd, err := os.Create(filename)
	if err != nil {
		panic(err)
	}

	return fd
}
