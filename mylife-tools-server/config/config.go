package config

import (
	config "github.com/gookit/config/v2"
	yaml "github.com/gookit/config/v2/yaml"

	log "mylife-tools-server/log"
)

var logger = log.CreateLogger("mylife:server:config")
var conf *config.Config

func init() {
	conf = config.NewWithOptions("mylife-config", config.ParseEnv, config.Readonly)

	// add driver for support yaml content
	conf.AddDriver(yaml.Driver)

	err := conf.LoadFiles("config.yaml")
	if err != nil {
		panic(err)
	}

	logger.WithField("data", conf.Data()).Info("Config loaded")
}

func BindStructure(key string, value any) {
	err := conf.Structure(key, value)
	if err != nil {
		panic(err)
	}

	logger.WithFields(log.Fields{"key": key, "value": value}).Debug("Config fetched")
}

func GetString(key string) string {
	value := conf.MustString(key)

	logger.WithFields(log.Fields{"key": key, "value": value}).Debug("Config fetched")
	return value
}
