package quantastic

import (
	"github.com/felixge/log"
	"github.com/felixge/quantastic/server/http"
	gohttp "net/http"
)

type Config struct {
	Log         log.Interface
	HttpHandler *http.Handler
	HttpAddr    string
}

func NewServer(config Config) *Server {
	return &Server{
		log:         config.Log,
		httpHandler: config.HttpHandler,
		httpAddr:    config.HttpAddr,
	}
}

type Server struct {
	log         log.Interface
	httpHandler *http.Handler
	httpAddr    string
}

func (s *Server) Run() error {
	s.log.Info("Starting server")
	return gohttp.ListenAndServe(s.httpAddr, s.httpHandler)
}
