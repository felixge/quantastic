package config

// @TODO it is rather annoying that goyaml lowercases field names as this
// breaks the camelCase convention of Go. It would be nice if we could deal
// with this in a way that doesn't require yaml tags for multi word fields.

import (
	pkglog "github.com/felixge/log"
	pkgapi "github.com/felixge/quantastic/api"
	pkgserver "github.com/felixge/quantastic/server"
	pkghttp "github.com/felixge/quantastic/server/http"
	pkgmite "github.com/felixge/quantastic/services/mite"
	"io/ioutil"
	"launchpad.net/goyaml"
	gohttp "net/http"
)

func Load(path string, config interface{}) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	return goyaml.Unmarshal(data, config)
}

type Server struct {
	Log  Log
	Api  Api
	Http Http
}

func NewServer(config Server) (*pkgserver.Server, error) {
	log, err := NewLog(config.Log)
	if err != nil {
		return nil, err
	}

	api, err := NewApi(config.Api, log)
	if err != nil {
		return nil, err
	}

	httpHandler, err := NewHttpHandler(config.Http, api, log)
	if err != nil {
		return nil, err
	}

	return pkgserver.NewServer(pkgserver.Config{
		Log:         log,
		HttpHandler: httpHandler,
		HttpAddr:    config.Http.Addr,
	}), nil
}

func NewLog(config Log) (*pkglog.Logger, error) {
	logLevel, err := pkglog.ParseLevel(config.Level)
	if err != nil {
		return nil, err
	}

	log := pkglog.NewLogger(pkglog.DefaultConfig)
	log.Handle(logLevel, pkglog.DefaultWriter)
	return log, nil
}

type Log struct {
	Level string
}

func NewApi(config Api, log pkglog.Interface) (*pkgapi.Api, error) {
	mite, err := NewMite(config.Mite)
	if err != nil {
		return nil, err
	}
	return pkgapi.NewApi(pkgapi.Config{
		Log:  log,
		Mite: mite,
	}), nil
}

type Api struct {
	Mite Mite
}

func NewMite(config Mite) (*pkgmite.Client, error) {
	return pkgmite.NewClient(config.Url, config.ApiKey)
}

type Mite struct {
	Url    string
	ApiKey string `yaml:"apiKey"`
}

func NewHttpHandler(config Http, api *pkgapi.Api, log pkglog.Interface) (*pkghttp.Handler, error) {
	return pkghttp.NewHandler(pkghttp.Config{
		Log:       log,
		Api:       api,
		Static:    gohttp.Dir(config.StaticDir),
		Templates: gohttp.Dir(config.TemplatesDir),
		BaseUrl:   config.BaseUrl,
	}), nil
}

type Http struct {
	Addr         string
	BaseUrl      string `yaml:"baseUrl"`
	StaticDir    string `yaml:"staticDir"`
	TemplatesDir string `yaml:"templatesDir"`
}
