package main

import (
	"github.com/jinzhu/copier"
	"github.com/jpillora/opts"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

const ConfigFile = "recipemark.yaml"
const AppName = "recipemark"

func main() {
	// Overly complex configuration setup
	defaultConfig := Config{
		Source:      "./recipes",
		Listen:      "127.0.0.1:3000",
		Destination: "./site",
		Assets:      "./assets",
	}
	g := Global{}
	build := Build{}
	serve := Serve{}
	subCommands := []interface{}{&build, &serve}

	if _, err := os.Stat(ConfigFile); !os.IsNotExist(err) {
		copyConfig(readConfigFile(ConfigFile, defaultConfig), subCommands)
	} else {
		copyConfig(defaultConfig, subCommands)
	}

	opth := opts.New(&g).Name(AppName).
		AddCommand(opts.New(&build).Summary("Build the site")).
		AddCommand(opts.New(&serve).Summary("Serve the site"))

	parsedOpts := opth.Parse()
	if g.Config != "" {
		copyConfig(readConfigFile(g.Config, defaultConfig), subCommands)
		parsedOpts = opth.Parse()
	}

	// Actually run the subcommand
	err := parsedOpts.Run()
	if err != nil {
		log.Fatalf("failed: %v", err)
	}
}

func copyConfig(c Config, dest []interface{}) {
	for _, sub := range dest {
		err := copier.Copy(sub, c)
		if err != nil {
			log.Fatalf("Failed to clone configuration file: %v", err)
		}
	}
}

func readConfigFile(filename string, defaults Config) Config {
	// defaults
	c := defaults
	f, err := os.Open(filename)
	if err != nil {
		log.Fatalf("Failed to open configuration file '%s': %v", filename, err)
	}
	err = yaml.NewDecoder(f).Decode(&c)
	if err != nil {
		log.Fatalf("Failed to parse configuration file '%s': %v", filename, err)
	}
	return c
}
