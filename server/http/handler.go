package http

import (
	"github.com/felixge/log"
	"github.com/felixge/quantastic/api"
	"github.com/felixge/quantastic/model"
	gohttp "net/http"
	"time"
)

type Config struct {
	Log       log.Interface
	Api       *api.Api
	Static    gohttp.FileSystem
	Templates gohttp.FileSystem
	BaseUrl   string
}

func NewHandler(c Config) *Handler {
	return &Handler{
		log:       c.Log,
		api:       c.Api,
		static:    gohttp.FileServer(c.Static),
		templates: newTemplates(c.Templates, c.Log),
		baseUrl:   c.BaseUrl,
	}
}

type Handler struct {
	log       log.Interface
	api       *api.Api
	static    gohttp.Handler
	templates *templates
	baseUrl   string
}

func (h *Handler) absoluteUrl(relativeUrl string) string {
	return h.baseUrl + relativeUrl
}

func (h *Handler) ServeHTTP(w gohttp.ResponseWriter, r *gohttp.Request) {
	h.log.Info("%s %s", r.Method, r.URL.String())
	switch r.URL.Path {
	case "/":
		w.Header().Set("Location", h.absoluteUrl("/time/day"))
		w.WriteHeader(gohttp.StatusSeeOther)
	case "/time/day":
		h.serveTime(w, r)
	case "/prototype":
		h.templates.Render(w, r, "prototype", nil)
	case "/css":
		h.templates.Render(w, r, "css", nil)
	default:
		h.static.ServeHTTP(w, r)
	}
}

func (h *Handler) serveTime(w gohttp.ResponseWriter, r *gohttp.Request) {
	timeRange := model.Day(time.Now())
	if err := h.api.ReadTimeRange(timeRange); err != nil {
		h.serveInternalError(w, r, "Could not read time range.", "err=%s", err)
		return
	}
	h.templates.Render(w, r, "time/index", timeRange)
}

func (h *Handler) serveInternalError(w gohttp.ResponseWriter, r *gohttp.Request, msg string, details string, args ...interface{}) {
	logArgs := []interface{}{msg + " " + details}
	h.log.Error(append(logArgs, args...))
	h.templates.Render(w, r, "_errors/internal", msg)
}
