package main

import (
	"flag"
	"github.com/hoshsadiq/m3ufilter/app"
	"github.com/mitchellh/go-homedir"
	"os"
)

func main() {
	configFile := flag.String("config", "~/.m3u.conf", "Config file location")
	useServer := flag.Bool("server", false, "Run a server to retrieve the playlist as a URL")
	playlistOutput := flag.String("playlist-output", "", "Where to output the playlist data. Ignored when using -server flag. Defaults to stdout")
	logOutput := flag.String("log-output", "", "Where to output logs. Defaults to stderr")
	flag.Parse()

	path, e := homedir.Expand(*configFile)
	if e != nil {
		panic(e)
	}

	app.Run(path, *useServer, fd(*playlistOutput, false), fd(*logOutput, true))
}

func fd(filename string, useStderr bool) *os.File {
	if filename == "-" {
		return os.Stdout
	}
	if filename == "" && useStderr {
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
