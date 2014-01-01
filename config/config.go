package config

// @TODO it is rather annoying that goyaml lowercases field names as this
// breaks the camelCase convention of Go. It would be nice if we could deal
// with this in a way that doesn't require yaml tags for multi word fields.

import (
	"io/ioutil"
	"launchpad.net/goyaml"
)

func Load(path string, config *Config) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return goyaml.Unmarshal(data, config)
}

type Config struct {
	Log  Log
	Http Http
	Mite Mite
}

type Log struct {
	Level string
}

type Mite struct {
	Url    string
	ApiKey string `yaml:"apiKey"`
}

type Http struct {
	Addr         string
	BaseUrl      string `yaml:"baseUrl"`
	StaticDir    string `yaml:"staticDir"`
	TemplatesDir string `yaml:"templatesDir"`
}
