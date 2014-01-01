package api

import (
	"github.com/felixge/log"
	"github.com/felixge/quantastic/services/mite"
)

type Config struct{
	Log log.Interface
	Mite *mite.Client
}

func NewApi(config Config) *Api {
	return &Api{
		mite: config.Mite,
	}
}

type Api struct {
	mite *mite.Client
}
