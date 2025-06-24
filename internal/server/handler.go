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

	//go:embed main.html main.css wordle.html spelling_bee.html
	_siteFS embed.FS
)

func init() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", wordleCheater)
	mux.HandleFunc("GET /spelling-bee", spellingBeeCheater)
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
	q := r.URL.Query()
	c, err := RunWordleCheater(q)
	if err != nil {
		handleBadRequest(w, "creating wordle cheater", err)
		return
	}
	p := Page{
		Name:    "wordle",
		Title:   "Wordle Cheater",
		Cheater: *c,
	}
	handleTemplate(w, tmpl, p)
}

func spellingBeeCheater(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	c, err := RunSpellingBeeCheater(q)
	if err != nil {
		handleBadRequest(w, "creating wordle cheater", err)
		return
	}
	p := Page{
		Name:    "spelling_bee",
		Title:   "Spelling Bee Cheater",
		Cheater: c,
	}
	handleTemplate(w, tmpl, p)
}

func handleBadRequest(w http.ResponseWriter, message string, err error) {
	message = fmt.Sprintf("%v: %v", message, err)
	http.Error(w, message, http.StatusBadRequest)
}

func handleTemplate(w http.ResponseWriter, tmpl *template.Template, data any) {
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		handleBadRequest(w, "rendering template", err)
		return
	}
	buf.WriteTo(w)
}

type Page struct {
	Name    string
	Title   string
	Cheater any
}
