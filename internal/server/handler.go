package server

import (
	"bytes"
	"embed"
	"fmt"
	"net/http"
	"text/template"
)

var (
	handler http.Handler
	tmpl    *template.Template

	//go:embed main.html main.css wordle.html
	_siteFS embed.FS
)

func init() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", wordleCheater)
	handler = mux

	inc := func(i int) int {
		return i + 1
	}
	funcs := template.FuncMap{
		"inc": inc,
	}
	tmpl = template.Must(template.New("main.html").
		Funcs(funcs).
		ParseFS(_siteFS, "*.html", "*.css"))
	tmpl.Funcs(funcs)
}

func ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h := handler
	h = withContentEncoding(h)
	h.ServeHTTP(w, r)
}

func wordleCheater(w http.ResponseWriter, r *http.Request) {
	c, err := RunWordleCheater(r.URL.Query())
	if err != nil {
		handleError(w, "creating wordle cheater", err)
		return
	}
	p := Page{
		Name:    "wordle",
		Title:   "Wordle Cheater",
		Cheater: *c,
	}
	handleTemplate(w, tmpl, "wordle.html", p)
}

func handleError(w http.ResponseWriter, message string, err error) {
	message = fmt.Sprintf("%v: %v", message, err)
	http.Error(w, message, http.StatusBadRequest)
}

func handleTemplate(w http.ResponseWriter, tmpl *template.Template, name string, data any) {
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		handleError(w, "rendering template", err)
		return
	}
	buf.WriteTo(w)
}

type Page struct {
	Name    string
	Title   string
	Cheater any
}
