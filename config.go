package main

import (
	"github.com/amaxlab/go-lib/config"
)

type KC868Configuration struct {
	Host string
	Port int
}

type Configuration struct {
	Debug       bool
	Port        int
	KC868Config *KC868Configuration
}

func NewConfiguration() *Configuration {
	loader := config.NewConfigLoader()

	configuration := Configuration{}
	configuration.Debug = loader.Bool("debug", false)
	configuration.Port = loader.Int("port", 8080)

	KC868Config := &KC868Configuration{}
	KC868Config.Host = loader.String("kc868-host", "192.168.0.1")
	KC868Config.Port = loader.Int("kc868-port", 4196)

	configuration.KC868Config = KC868Config

	return &configuration
}
