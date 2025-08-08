package server

import (
	"bytes"
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

//go:embed main.html main.css wordle.html spelling_bee.html letter_boxed.html instructions.html
var _siteFS embed.FS

type Handler struct {
	wordsText string
	mux       http.Handler
	tmpl      *template.Template
}

const (
	wordlePath      = "/"
	spellingBeePath = "/spelling-bee"
	letterBoxedPath = "/letter-boxed"
)

func NewHandler(wordsText string) *Handler {
	mux := http.NewServeMux()

	inc := func(i int) int {
		return i + 1
	}
	arr := func(s ...string) []string {
		return s
	}
	funcs := template.FuncMap{
		"inc": inc,
		"arr": arr,
	}
	tmpl := template.Must(newTemplate().
		Funcs(funcs).
		ParseFS(_siteFS, "*.html", "*.css"))

	h := Handler{
		wordsText: wordsText,
		mux:       mux,
		tmpl:      tmpl,
	}

	mux.HandleFunc("GET "+wordlePath+"{$}", h.handle(wordlePage))
	mux.HandleFunc("GET "+spellingBeePath, h.handle(spellingBeePage))
	mux.HandleFunc("GET "+letterBoxedPath, h.handle(letterBoxedPage))
	return &h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	withContentEncoding(h.mux).
		ServeHTTP(w, r)
}

func (h *Handler) handle(p page) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		d, err := p.newDisplay(q, h.wordsText)
		if err != nil {
			handleBadRequest(w, "creating cheater", err)
			return
		}
		d.NoJS = q.Has("NoJS")
		tmplName := r.Header.Get("Hx-Target")
		if tmplName == "main-template" {
			tmplName = p.tmplName
		}
		tmpl := resolveTemplate(h.tmpl, tmplName)
		handleTemplate(tmpl, w, d)
	}
}

func resolveTemplate(tmpl *template.Template, tmplName string) *template.Template {
	if len(tmplName) > 0 {
		tmpl = tmpl.Lookup(tmplName)
	}
	return tmpl
}

func handleTemplate(tmpl *template.Template, w http.ResponseWriter, data any) {
	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, data); err != nil {
		handleBadRequest(w, "rendering template", err)
		return
	}
	buf.WriteTo(w)
}

func handleBadRequest(w http.ResponseWriter, message string, err error) {
	message = fmt.Sprintf("%v: %v", message, err)
	http.Error(w, message, http.StatusBadRequest)
}

func newTemplate() *template.Template {
	tmpl := template.New("main.html")
	tmpl.Option("missingkey=error")
	return tmpl
}
