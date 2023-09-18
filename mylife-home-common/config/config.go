package config

import (
	config "github.com/gookit/config/v2"
	yaml "github.com/gookit/config/v2/yaml"

	log "mylife-home-common/log"
)

var logger = log.CreateLogger("mylife:home:config")
var conf *config.Config

func init() {
	conf = config.NewWithOptions("mylife-home-config", config.ParseEnv, config.Readonly)

	// add driver for support yaml content
	conf.AddDriver(yaml.Driver)

	err := conf.LoadFiles("config.yaml")
	if err != nil {
		panic(err)
	}

	logger.Infof("Config loaded: %+v", conf.Data())
}

func BindStructure(key string, value any) {
	err := conf.Structure(key, value)
	if err != nil {
		panic(err)
	}

	logger.Debugf("Config '%s' fetched: %+v", key, value)
}

func GetString(key string) string {
	value := conf.MustString(key)

	logger.Debugf("Config '%s' fetched: %s", key, value)
	return value
}
