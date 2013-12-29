package http

import (
	"github.com/felixge/log"
	"github.com/felixge/quantastic/mite"
	_http "net/http"
)

type Config struct {
	Mite      *mite.Client
	Static    _http.FileSystem
	Templates _http.FileSystem
	Log       log.Interface
}

func NewHandler(c Config) *Handler {
	return &Handler{
		log:       c.Log,
		mite:      c.Mite,
		static:    _http.FileServer(c.Static),
		templates: newTemplates(c.Templates, c.Log),
	}
}

type Handler struct {
	log       log.Interface
	mite      *mite.Client
	static    _http.Handler
	templates *templates
}

func (h *Handler) ServeHTTP(w _http.ResponseWriter, r *_http.Request) {
	h.log.Info("%s %s", r.Method, r.URL.String())
	if r.URL.Path == "/" {
		h.serveTimeIndex(w, r)
		return
	}

	h.static.ServeHTTP(w, r)
}

func (h *Handler) serveTimeIndex(w _http.ResponseWriter, r *_http.Request) {
	customers, err := h.mite.Customers(nil)
	if err != nil {
		h.serveInternalError(w, r, "Could not load customers.", "err=%s", err)
		return
	}

	h.log.Debug("customers: %#v", customers)
	h.templates.Render(w, r, "time/index", customers)
}

func (h *Handler) serveInternalError(w _http.ResponseWriter, r *_http.Request, msg string, details string, args ...interface{}) {
	logArgs := []interface{}{msg+" "+details}
	h.log.Error(append(logArgs, args...))
	h.templates.Render(w, r, "_errors/internal", msg)
}
