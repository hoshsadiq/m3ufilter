package main

import (
	"flag"
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
	"github.com/hoshsadiq/m3ufilter/m3u"
	"github.com/hoshsadiq/m3ufilter/server"
	"github.com/hoshsadiq/m3ufilter/writer"
	"github.com/mitchellh/go-homedir"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	_ "net/http/pprof"
	"os"
)

func main() {
	configFile := flag.String("config", "~/.m3u.conf", "Config file location")
	playlistOutput := flag.String("playlist", "", "Where to output the playlist data. Ignored when using -server flag. Defaults to stdout")
	logOutput := flag.String("log", "", "Where to output logs. Defaults to stderr")
	flag.Parse()

	path, e := homedir.Expand(*configFile)
	if e != nil {
		panic(e)
	}

	run(path, fd(*playlistOutput, false), fd(*logOutput, true))
}

func run(configFilename string, stdout *os.File, stderr *os.File) {
	log := logger.Get()
	log.SetOutput(stderr)

	yamlFile, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Fatalf("could not read config file %s, err = %v", configFilename, err)
	}

	conf := config.Get()
	err = yaml.Unmarshal([]byte(yamlFile), &conf)
	if err != nil {
		log.Fatalf("could not parse config file %s, err = %v", configFilename, err)
	}

	if conf.Core.ServerListen != "" {
		server.Serve(conf)
	} else {
		playlist := m3u.GetPlaylist(conf)
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

	fd, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

	return fd
}
