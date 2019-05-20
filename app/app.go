package app

import (
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/logger"
	"github.com/hoshsadiq/m3ufilter/server"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

func Run(configFilename string) {
	log := logger.Get()

	yamlFile, err := ioutil.ReadFile(configFilename)
	if err != nil {
		log.Fatalf("could not read config file %s, err = %v", configFilename, err)
	}

	var conf *config.Config
	err = yaml.Unmarshal([]byte(yamlFile), &conf)
	if err != nil {
		log.Fatalf("could not parse config file %s, err = %v", configFilename, err)
	}

	server.Serve(conf)
}
