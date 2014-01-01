package http

import (
	"github.com/felixge/log"
	"github.com/felixge/quantastic/services/mite"
	gohttp "net/http"
)

type Config struct {
	Mite      *mite.Client
	Static    gohttp.FileSystem
	Templates gohttp.FileSystem
	Log       log.Interface
	BaseUrl   string
}

func NewHandler(c Config) *Handler {
	return &Handler{
		log:       c.Log,
		mite:      c.Mite,
		static:    gohttp.FileServer(c.Static),
		templates: newTemplates(c.Templates, c.Log),
		baseUrl:   c.BaseUrl,
	}
}

type Handler struct {
	log       log.Interface
	mite      *mite.Client
	static    gohttp.Handler
	templates *templates
	baseUrl   string
}

func (h *Handler) absoluteUrl(relativeUrl string) string {
	return h.baseUrl+relativeUrl
}

func (h *Handler) ServeHTTP(w gohttp.ResponseWriter, r *gohttp.Request) {
	h.log.Info("%s %s", r.Method, r.URL.String())
	switch r.URL.Path {
	case "/":
		w.Header().Set("Location", h.absoluteUrl("/time/week"))
		w.WriteHeader(gohttp.StatusSeeOther)
	case "/time/week":
		h.serveTimeIndex(w, r)
	case "/css":
		h.templates.Render(w, r, "css", nil)
	default:
		h.static.ServeHTTP(w, r)
	}
}

func (h *Handler) serveTimeIndex(w gohttp.ResponseWriter, r *gohttp.Request) {
	//customers, err := h.mite.Customers(nil)
	//if err != nil {
	//h.serveInternalError(w, r, "Could not load customers.", "err=%s", err)
	//return
	//}

	//h.log.Debug("customers: %#v", customers)
	//h.templates.Render(w, r, "time/index", customers)
	h.templates.Render(w, r, "time/index", nil)
}

func (h *Handler) serveInternalError(w gohttp.ResponseWriter, r *gohttp.Request, msg string, details string, args ...interface{}) {
	logArgs := []interface{}{msg + " " + details}
	h.log.Error(append(logArgs, args...))
	h.templates.Render(w, r, "_errors/internal", msg)
}
