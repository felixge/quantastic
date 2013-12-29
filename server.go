package quantastic

import (
	"github.com/felixge/log"
	"github.com/felixge/quantastic/config"
	"github.com/felixge/quantastic/http"
	"github.com/felixge/quantastic/mite"
	_http "net/http"
)

func NewServer(configPath string) (s *Server, err error) {
	s = &Server{}
	if err = config.Load("config.xml", &s.Config); err != nil {
		return
	}

	var logLevel log.Level
	logLevel, err = log.ParseLevel(s.Config.Log.Level)
	if err != nil {
		return
	}

	s.Log = log.NewLogger(log.DefaultConfig)
	s.Log.Handle(logLevel, log.DefaultWriter)

	s.Mite, err = mite.NewClient(s.Config.Mite.Url, s.Config.Mite.ApiKey)
	if err != nil {
		return
	}
	s.Http = http.NewHandler(http.Config{
		Log:       s.Log,
		Mite:      s.Mite,
		Static:    _http.Dir(s.Config.Http.StaticDir),
		Templates: _http.Dir(s.Config.Http.TemplatesDir),
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
	return _http.ListenAndServe(s.Config.Http.Addr, s.Http)
}
