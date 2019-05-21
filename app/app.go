package app

import (
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
	"github.com/hoshsadiq/m3ufilter/m3u"
	"github.com/hoshsadiq/m3ufilter/server"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
)

func Run(configFilename string, useServer bool, stdout *os.File, stderr *os.File) {
	log := logger.Get()
	log.SetOutput(stderr)

	yamlFile, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Fatalf("could not read config file %s, err = %v", configFilename, err)
	}

	var conf *config.Config
	err = yaml.Unmarshal([]byte(yamlFile), &conf)
	if err != nil {
		log.Fatalf("could not parse config file %s, err = %v", configFilename, err)
	}

	if useServer {
		server.Serve(conf)
	} else {
		m3u.GetPlaylist(stdout, conf)
	}
}
