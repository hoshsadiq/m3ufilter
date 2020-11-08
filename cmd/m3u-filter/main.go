package main

import (
	"flag"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
	"github.com/hoshsadiq/m3ufilter/m3u"
	"github.com/hoshsadiq/m3ufilter/net"
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
	playlistOutput := flag.String("playlist", "", "Where to output the playlist data. Ignored when using server options in the config. Defaults to stdout")
	logOutput := flag.String("log", "", "Where to output logs. Defaults to stderr")
	generatingCsv := flag.Bool("csv", false, "Generate CSV instead of processing the M3U")
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
	run(path, *generatingCsv, fd(*playlistOutput, false), fd(*logOutput, true))
}

func run(configFilename string, generatingCsv bool, stdout *os.File, stderr *os.File) {
	log := logger.Get()
	log.SetOutput(stderr)

	conf, err := config.New(configFilename)
	if err != nil {
		os.Exit(1)
	}

	net.InitClient(conf)
	if generatingCsv {
		playlist, _, _ := m3u.ProcessConfig(conf, generatingCsv)
		writer.WriteCsv(stdout, playlist)
		return
	}

	if !generatingCsv && conf.Core.ServerListen != "" {
		server.Serve(conf)
	} else {
		playlist, _ /*todo epg*/, _ := m3u.ProcessConfig(conf, generatingCsv)
		writer.WriteM3U(stdout, playlist)
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
