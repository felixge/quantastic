package quantastic

import (
	"github.com/felixge/log"
	"github.com/felixge/quantastic/config"
	"github.com/felixge/quantastic/server/http"
	"github.com/felixge/quantastic/services/mite"
	gohttp "net/http"
)

func NewServer(config config.Config) (s *Server, err error) {
	s = &Server{Config: config}
	var logLevel log.Level
	logLevel, err = log.ParseLevel(config.Log.Level)
	if err != nil {
		return
	}

	s.Log = log.NewLogger(log.DefaultConfig)
	s.Log.Handle(logLevel, log.DefaultWriter)

	s.Mite, err = mite.NewClient(config.Mite.Url, config.Mite.ApiKey)
	if err != nil {
		return
	}
	s.Http = http.NewHandler(http.Config{
		Log:       s.Log,
		Mite:      s.Mite,
		Static:    gohttp.Dir(config.Http.StaticDir),
		Templates: gohttp.Dir(config.Http.TemplatesDir),
		BaseUrl:   config.Http.BaseUrl,
	})
	return
}

type Server struct {
	Config config.Config
	Log    *log.Logger
	Mite   *mite.Client
	Http   *http.Handler
}

func (s *Server) Run() error {
	defer s.Log.Flush()

	s.Log.Info("Starting server")
	return gohttp.ListenAndServe(s.Config.Http.Addr, s.Http)
}
