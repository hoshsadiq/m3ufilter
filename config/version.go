package config

import (
	"fmt"
	"os"
)

// all these are set at compile time
var (
	Version   = ""
	GitCommit = ""
	BuildDate = ""
	GoVersion = ""
	Platform  = ""
)

func ShowVersion() {
	fmt.Printf("Version: %s\n", Version)
	fmt.Printf("GitCommit: %s\n", GitCommit)
	fmt.Printf("BuildDate: %s\n", BuildDate)
	fmt.Printf("GoVersion: %s\n", GoVersion)
	fmt.Printf("Platform: %s\n", Platform)
	os.Exit(0)
}
