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

const (
	wordlePath      = "/"
	spellingBeePath = "/spelling-bee"
	letterBoxedPath = "/letter-boxed"
)

func NewHandler(wordsText string) http.Handler {
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
	
	mux := http.NewServeMux()
	mux.HandleFunc("GET "+wordlePath+"{$}", handle(wordlePage, wordsText, tmpl))
	mux.HandleFunc("GET "+spellingBeePath, handle(spellingBeePage, wordsText, tmpl))
	mux.HandleFunc("GET "+letterBoxedPath, handle(letterBoxedPage, wordsText, tmpl))

	return withContentEncoding(mux)
}

func handle(p page, wordsText string, tmpl *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		d, err := p.newDisplay(q, wordsText)
		if err != nil {
			handleBadRequest(w, "creating cheater", err)
			return
		}
		d.NoJS = q.Has("NoJS")
		tmplName := r.Header.Get("Hx-Target")
		if tmplName == "main-template" {
			tmplName = p.tmplName
		}
		tmpl := resolveTemplate(tmpl, tmplName)
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
