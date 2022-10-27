package main

import (
	"flag"
	"onboarding"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "configPath", "config.yaml", "The path to the CSV file.")
}

func main() {
	flag.Parse()
	onboarding.LoadPeople(configPath)
}
