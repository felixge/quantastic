package http

import (
	"fmt"
	"github.com/felixge/log"
	"github.com/tav/golly/httputil"
	"html/template"
	"io/ioutil"
	_http "net/http"
)

var extensions = map[string]string{
	"text/html": ".html",
}

func newTemplates(fs _http.FileSystem, log log.Interface) *templates {
	return &templates{fs: fs, log: log}
}

type templates struct {
	fs  _http.FileSystem
	log log.Interface
}

func (t *templates) Render(w _http.ResponseWriter, r *_http.Request, path string, data interface{}) {
	ext := preferredExtension(r)
	path += ext

	tmpl := template.New("content")
	if err := t.load(tmpl, path); err != nil {
		t.renderCouldNotLoadTemplate(w, err, path)
		return
	}

	tmpl = tmpl.New("layout")
	layoutPath := "_layouts/default" + ext
	if err := t.load(tmpl, layoutPath); err != nil {
		t.renderCouldNotLoadTemplate(w, err, layoutPath)
		return
	}

	if err := tmpl.Execute(w, data); err != nil {
		t.renderCouldNotLoadTemplate(w, err, path)
		return
	}
}

func (t *templates) RenderInternalError(w _http.ResponseWriter, err error, text string) {
	w.WriteHeader(_http.StatusInternalServerError)
	fmt.Fprintf(w, text)
}

func (t *templates) load(tmpl *template.Template, path string) error {
	file, err := t.fs.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	tmplData, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if _, err := tmpl.Parse(string(tmplData)); err != nil {
		return err
	}
	return nil
}

const couldNotLoadTemplate = "Could not load template."

func (t *templates) renderCouldNotLoadTemplate(w _http.ResponseWriter, err error, path string) {
	w.WriteHeader(_http.StatusInternalServerError)
	fmt.Fprintf(w, couldNotLoadTemplate)
	t.log.Error("%s path=%s err=%s", couldNotLoadTemplate, path, err)
}

func preferredExtension(r *_http.Request) string {
	types := make([]string, 0, len(extensions))
	for key, _ := range extensions {
		types = append(types, key)
	}

	acceptable := httputil.Parse(r, "Accept")
	preferred := acceptable.FindPreferred(types...)
	if len(preferred) == 0 {
		return ".html"
	}
	return extensions[preferred[0]]
}
