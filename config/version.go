package config

import (
	"fmt"
	"os"
	"strings"
)

// all these are set at compile time
var (
	Version   = ""
	GitCommit = ""
	BuildDate = ""
	GoVersion = ""
	Platform  = ""
)

var epgGeneratorName = "M3U Filter"
var epgGeneratorUrl = "" // todo actually we want to retrieve this from http package

func ShowVersion() {
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("GitCommit: %s\n", GitCommit)
	fmt.Printf("BuildDate: %s\n", BuildDate)
	fmt.Printf("GoVersion: %s\n", GoVersion)
	fmt.Printf("Platform: %s\n", Platform)
	os.Exit(0)
}

func EpgGeneratorName() string {
	return strings.TrimSpace(fmt.Sprintf("%s %s", epgGeneratorName, Version))
}
func EpgGeneratorUrl() string {
	return epgGeneratorUrl
}
