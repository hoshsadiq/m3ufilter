package main

import (
	"flag"
	"github.com/hoshsadiq/m3ufilter/app"
	"github.com/mitchellh/go-homedir"
)

func main() {
	configFile := flag.String("config", "~/.m3u.conf", "Config file location")
	flag.Parse()

	path, e := homedir.Expand(*configFile)
	if e != nil {
		panic(e)
	}

	app.Run(path)
}
