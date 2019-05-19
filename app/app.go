package app

import (
	"github.com/hoshsadiq/m3ufilter/config"
	"github.com/hoshsadiq/m3ufilter/server"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
)

func Run() {
	yamlFile, err := ioutil.ReadFile("config.yaml") // todo path needs to be read from args
	if err != nil {
		panic(err)
	}

	var conf *config.Config
	err = yaml.Unmarshal([]byte(yamlFile), &conf)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	server.Serve(conf)
}
