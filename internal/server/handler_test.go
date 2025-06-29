package server

import (
	"fmt"
	"html/template"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewHandler(t *testing.T) {
	tests := []struct {
		name     string
		target   string
		wantCode int
	}{
		{
			name:     "unknown",
			target:   "/unknown",
			wantCode: 404,
		},
		{
			name:     "wordle-empty",
			target:   wordlePath,
			wantCode: 200,
		},
		{
			name:     "wordle-ok",
			target:   wordlePath + "?g0=words&s0=ccccc",
			wantCode: 200,
		},
		{
			name:     "wordle-bad",
			target:   wordlePath + "?g0=word&s0=c",
			wantCode: 400,
		},
		{
			name:     "spelling-bee-empty",
			target:   spellingBeePath,
			wantCode: 200,
		},
		{
			name:     "spelling-bee-ok",
			target:   spellingBeePath + "?" + centralLetterParam + "=a&" + otherLettersParam + "=bcdefg",
			wantCode: 200,
		},
		{
			name:     "spelling-bee-bad",
			target:   spellingBeePath + "?" + centralLetterParam + "=az&" + otherLettersParam + "=bcdefgh",
			wantCode: 400,
		},
		{
			name:     "spelling-bee-404",
			target:   spellingBeePath + "-404",
			wantCode: 404,
		},
		{
			name:     "letter-boxed-empty",
			target:   letterBoxedPath,
			wantCode: 200,
		},
		{
			name:     "letter-boxed-ok",
			target:   letterBoxedPath + "?" + letterBoxedLettersParam + "=abcdefghijkl",
			wantCode: 200,
		},
		{
			name:     "letter-boxed-bad-count",
			target:   letterBoxedPath + "?" + letterBoxedLettersParam + "=hello",
			wantCode: 400,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var wordsText string
			h := NewHandler(wordsText)
			w := httptest.NewRecorder()
			r := httptest.NewRequest("GET", test.target, nil)
			h.ServeHTTP(w, r)
			if want, got := test.wantCode, w.Result().StatusCode; want != got {
				t.Errorf("wanted %v, got %v (body: %q)", want, got, w.Body.String())
			}
		})
	}
}

func TestHandleBadRequest(t *testing.T) {
	message := "my-message"
	err := fmt.Errorf("my-error")
	w := httptest.NewRecorder()
	handleBadRequest(w, message, err)
	if want, got := 400, w.Result().StatusCode; want != got {
		t.Errorf("status codes: wanted %v, got %v", want, got)
	}
	body := strings.TrimSpace(w.Body.String())
	if want, got := "my-message: my-error", body; want != got {
		t.Errorf("response bodies: wanted: %q got: %q", want, got)
	}
}

func TestHandleTemplateError(t *testing.T) {
	tests := []struct {
		name          string
		data          any
		wantBodyStart string
		wantCode      int
	}{
		{"ok", struct{ A int }{A: 9}, "9", 200},
		{"bad", nil, "rendering", 400},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpl := template.Must(newTemplate().Parse("{{.A}}"))
			w := httptest.NewRecorder()
			handleTemplate(tmpl, w, test.data)
			if want, got := test.wantCode, w.Result().StatusCode; want != got {
				t.Errorf("status codes: wanted %v, got %v", want, got)
			}
			got := w.Body.String()
			if !strings.HasPrefix(got, test.wantBodyStart) {
				t.Errorf("response body: wanted prefix: %q got: %q", test.wantBodyStart, got)
			}
		})
	}
}

func TestHandleCheaterHtmx(t *testing.T) {
	tests := []struct {
		name    string
		headers map[string]string
		want    string
	}{
		{"all", nil, "OuterInner"},
		{"inner only", map[string]string{htmxHeader: "true"}, "Inner"},
	}
	pageTypes := []pageType{wordle_type, spelling_bee_type, letter_boxed_type}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, pt := range pageTypes {
				t.Run(fmt.Sprint(pt), func(t *testing.T) {
					r := httptest.NewRequest("GET", "/", nil)
					for k, v := range test.headers {
						r.Header.Add(k, v)
					}
					w := httptest.NewRecorder()
					name := pt.HtmxTemplateName()
					text := `Outer{{block "` + name + `" .}}Inner{{end}}`
					tmpl := template.Must(template.New("").Parse(text))
					h := Handler{
						tmpl: tmpl,
					}
					fn := h.handleCheater(pt)
					fn(w, r)
					statusCode := w.Code
					want, got := test.want, w.Body.String()
					switch {
					case statusCode != 200:
						t.Errorf("bad status: %v", statusCode)
					case want != got:
						t.Errorf("response bodies: \n wanted: %q \n    got: %q", want, got)
					}
				})
			}
		})
	}
}
