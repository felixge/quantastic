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
	Mite Mite
}

type Mite struct {
	Url    string
	ApiKey string
}
