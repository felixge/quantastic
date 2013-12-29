package config

import (
	"encoding/xml"
	"io/ioutil"
)

func Load(path string, config *Config) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, config)
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
	ApiKey string
}

type Http struct {
	Addr         string
	StaticDir    string
	TemplatesDir string
}
