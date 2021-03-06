package common

import (
	"github.com/BurntSushi/toml"
	"io/ioutil"
)

type config struct {
	LogFile  string
	Database string
	Host     string
	Port     int
}

var ConfigRoot *config

func LoadConfig() {
	buf, err := ioutil.ReadFile("Config.toml")
	if err != nil {
		panic("Could not read Config.toml")
	}

	var conf config
	if _, err := toml.Decode(string(buf), &conf); err != nil {
		panic("Could not parse Config.toml: " + err.Error())
	}

	ConfigRoot = &conf
}
