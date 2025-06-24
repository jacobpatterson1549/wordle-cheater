package server

import (
	"bytes"
	"embed"
	"fmt"
	"net/http"
	"text/template"
)

//go:embed main.html main.css wordle.html spelling_bee.html
var _siteFS embed.FS

type handler struct {
	wordsText string
	mux       http.Handler
	tmpl      *template.Template
}

func NewHandler(wordsText string) http.Handler {
	mux := http.NewServeMux()

	inc := func(i int) int {
		return i + 1
	}
	funcs := template.FuncMap{
		"inc": inc,
	}
	tmpl := template.Must(template.New("main.html").
		Funcs(funcs).
		ParseFS(_siteFS, "*.html", "*.css"))

	h := handler{
		wordsText: wordsText,
		mux:       mux,
		tmpl:      tmpl,
	}

	mux.Handle("GET /{$}", h.wordleCheater())
	mux.Handle("GET /spelling-bee", h.spellingBeeCheater())
	return h
}

func (h handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	withContentEncoding(h.mux).
		ServeHTTP(w, r)
}

func (h handler) wordleCheater() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		c, err := RunWordleCheater(q, h.wordsText)
		if err != nil {
			handleBadRequest(w, "creating wordle cheater", err)
			return
		}
		p := Page{
			Name:    "wordle",
			Title:   "Wordle Cheater",
			Cheater: *c,
		}
		h.handleTemplate(w, p)
	}
}

func (h handler) spellingBeeCheater() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		c, err := RunSpellingBeeCheater(q, h.wordsText)
		if err != nil {
			handleBadRequest(w, "creating wordle cheater", err)
			return
		}
		p := Page{
			Name:    "spelling_bee",
			Title:   "Spelling Bee Cheater",
			Cheater: c,
		}
		h.handleTemplate(w, p)
	}
}

func handleBadRequest(w http.ResponseWriter, message string, err error) {
	message = fmt.Sprintf("%v: %v", message, err)
	http.Error(w, message, http.StatusBadRequest)
}

func (h handler) handleTemplate(w http.ResponseWriter, data any) {
	buf := new(bytes.Buffer)
	if err := h.tmpl.Execute(buf, data); err != nil {
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
